package astiav

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/asticode/go-astikit"
)

var globalHelper = newHelper()

func TestMain(m *testing.M) {
	// Make sure to exit with the proper code
	var code int
	defer func(code *int) {
		os.Exit(*code)
	}(&code)

	// Make sure to close global helper
	defer globalHelper.close()

	// Run
	code = m.Run()
}

type helper struct {
	closer *astikit.Closer
	inputs map[string]*helperInput
	m      *sync.Mutex // Locks inputs
}

func newHelper() *helper {
	return &helper{
		closer: astikit.NewCloser(),
		inputs: make(map[string]*helperInput),
		m:      &sync.Mutex{},
	}
}

func (h *helper) close() {
	h.closer.Close()
}

type helperInput struct {
	firstPkt      *Packet
	formatContext *FormatContext
	lastFrame     *Frame
}

func (h *helper) inputFormatContext(name string) (fc *FormatContext, err error) {
	h.m.Lock()
	i, ok := h.inputs[name]
	if ok && i.formatContext != nil {
		h.m.Unlock()
		return i.formatContext, nil
	}
	h.m.Unlock()

	if fc = AllocFormatContext(); fc == nil {
		err = errors.New("astiav_test: allocated format context is nil")
		return
	}
	h.closer.Add(fc.Free)

	if err = fc.OpenInput("testdata/"+name, nil, nil); err != nil {
		err = fmt.Errorf("astiav_test: opening input failed: %w", err)
		return
	}
	h.closer.Add(fc.CloseInput)

	if err = fc.FindStreamInfo(nil); err != nil {
		err = fmt.Errorf("astiav_test: finding stream info failed: %w", err)
		return
	}

	h.m.Lock()
	if _, ok := h.inputs[name]; !ok {
		h.inputs[name] = &helperInput{}
	}
	h.inputs[name].formatContext = fc
	h.m.Unlock()
	return
}

func (h *helper) inputFirstPacket(name string) (pkt *Packet, err error) {
	h.m.Lock()
	i, ok := h.inputs[name]
	if ok && i.firstPkt != nil {
		h.m.Unlock()
		return i.firstPkt, nil
	}
	h.m.Unlock()

	var fc *FormatContext
	if fc, err = h.inputFormatContext(name); err != nil {
		err = fmt.Errorf("astiav_test: getting input format context failed")
		return
	}

	pkt = AllocPacket()
	if pkt == nil {
		err = errors.New("astiav_test: pkt is nil")
		return
	}
	h.closer.Add(pkt.Free)

	if err = fc.ReadFrame(pkt); err != nil {
		err = fmt.Errorf("astiav_test: reading frame failed: %w", err)
		return
	}

	h.m.Lock()
	h.inputs[name].firstPkt = pkt
	h.m.Unlock()
	return
}

func (h *helper) inputLastFrame(name string, mediaType MediaType) (f *Frame, err error) {
	h.m.Lock()
	i, ok := h.inputs[name]
	if ok && i.lastFrame != nil {
		h.m.Unlock()
		return i.lastFrame, nil
	}
	h.m.Unlock()

	var fc *FormatContext
	if fc, err = h.inputFormatContext(name); err != nil {
		err = fmt.Errorf("astiav_test: getting input format context failed: %w", err)
		return
	}

	var cc *CodecContext
	var cs *Stream
	for _, s := range fc.Streams() {
		if s.CodecParameters().MediaType() != mediaType {
			continue
		}

		cs = s

		c := FindDecoder(s.CodecParameters().CodecID())
		if c == nil {
			err = errors.New("astiav_test: no codec")
			return
		}

		cc = AllocCodecContext(c)
		if cc == nil {
			err = errors.New("astiav_test: no codec context")
			return
		}
		h.closer.Add(cc.Free)

		if err = cs.CodecParameters().ToCodecContext(cc); err != nil {
			err = fmt.Errorf("astiav_test: updating codec context failed: %w", err)
			return
		}

		if err = cc.Open(c, nil); err != nil {
			err = fmt.Errorf("astiav_test: opening codec context failed: %w", err)
			return
		}
		break
	}

	if cs == nil {
		err = errors.New("astiav_test: no valid video stream")
		return
	}

	var pkt1 *Packet
	if pkt1, err = h.inputFirstPacket(name); err != nil {
		err = fmt.Errorf("astiav_test: getting input first packet failed: %w", err)
		return
	}

	pkt2 := AllocPacket()
	h.closer.Add(pkt2.Free)

	f = AllocFrame()
	h.closer.Add(f.Free)

	lastFrame := AllocFrame()
	h.closer.Add(lastFrame.Free)

	pkts := []*Packet{pkt1}
	for {
		if err = fc.ReadFrame(pkt2); err != nil {
			if errors.Is(err, ErrEof) || errors.Is(err, ErrEagain) {
				if len(pkts) == 0 {
					if err = f.Ref(lastFrame); err != nil {
						err = fmt.Errorf("astiav_test: last refing frame failed: %w", err)
						return
					}
					err = nil
					break
				}
			} else {
				err = fmt.Errorf("astiav_test: reading frame failed: %w", err)
				return
			}
		} else {
			pkts = append(pkts, pkt2)
		}

		for _, pkt := range pkts {
			if pkt.StreamIndex() != cs.Index() {
				continue
			}

			if err = cc.SendPacket(pkt); err != nil {
				err = fmt.Errorf("astiav_test: sending packet failed: %w", err)
				return
			}

			for {
				if err = cc.ReceiveFrame(f); err != nil {
					if errors.Is(err, ErrEof) || errors.Is(err, ErrEagain) {
						err = nil
						break
					}
					err = fmt.Errorf("astiav_test: receiving frame failed: %w", err)
					return
				}

				if err = lastFrame.Ref(f); err != nil {
					err = fmt.Errorf("astiav_test: refing frame failed: %w", err)
					return
				}
			}
		}

		pkts = []*Packet{}
	}

	h.m.Lock()
	h.inputs[name].lastFrame = f
	h.m.Unlock()
	return
}
