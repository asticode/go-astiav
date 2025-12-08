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
	lastFrames    map[MediaType]*Frame
}

func newHelperInput() *helperInput {
	return &helperInput{lastFrames: make(map[MediaType]*Frame)}
}

func (h *helper) inputFormatContext(name string, ifmt *InputFormat) (fc *FormatContext, err error) {
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

	if err = fc.OpenInput("testdata/"+name, ifmt, nil); err != nil {
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
		h.inputs[name] = newHelperInput()
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
	if fc, err = h.inputFormatContext(name, nil); err != nil {
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

func (h *helper) inputLastFrame(name string, mediaType MediaType, ifmt *InputFormat) (*Frame, error) {
	h.m.Lock()
	if i, ok := h.inputs[name]; ok {
		if len(i.lastFrames) > 0 {
			f, ok := i.lastFrames[mediaType]
			h.m.Unlock()
			if ok {
				return f, nil
			}
			return nil, fmt.Errorf("astiav_test: no last frame for media type %s", mediaType)
		}
	}
	h.m.Unlock()

	fc, err := h.inputFormatContext(name, ifmt)
	if err != nil {
		return nil, fmt.Errorf("astiav_test: getting input format context failed: %w", err)
	}

	type stream struct {
		cc *CodecContext
		s  *Stream
	}
	streams := make(map[int]*stream)
	mediaTypeFound := false
	for _, v := range fc.Streams() {
		s := &stream{s: v}
		streams[v.Index()] = s

		c := FindDecoder(v.CodecParameters().CodecID())
		if c == nil {
			return nil, fmt.Errorf("astiav_test: no decoder found for %s", v.CodecParameters().CodecID())
		}

		s.cc = AllocCodecContext(c)
		if s.cc == nil {
			return nil, errors.New("astiav_test: no codec context")
		}
		h.closer.Add(s.cc.Free)

		if err = s.s.CodecParameters().ToCodecContext(s.cc); err != nil {
			return nil, fmt.Errorf("astiav_test: updating codec context failed: %w", err)
		}

		if err = s.cc.Open(c, nil); err != nil {
			return nil, fmt.Errorf("astiav_test: opening codec context failed: %w", err)
		}

		if _, ok := h.inputs[name].lastFrames[s.cc.MediaType()]; !ok {
			h.inputs[name].lastFrames[s.cc.MediaType()] = AllocFrame()
			h.closer.Add(h.inputs[name].lastFrames[s.cc.MediaType()].Free)
		}

		if s.cc.MediaType() == mediaType {
			mediaTypeFound = true
		}
	}

	if !mediaTypeFound {
		return nil, fmt.Errorf("astiav_test: no stream for media type %s", mediaType)
	}

	var pkt1 *Packet
	if pkt1, err = h.inputFirstPacket(name); err != nil {
		return nil, fmt.Errorf("astiav_test: getting input first packet failed: %w", err)
	}

	pkt2 := AllocPacket()
	h.closer.Add(pkt2.Free)

	f := AllocFrame()
	h.closer.Add(f.Free)

	pkts := []*Packet{pkt1}
	for {
		if err = fc.ReadFrame(pkt2); err != nil {
			if errors.Is(err, ErrEof) || errors.Is(err, ErrEagain) {
				if len(pkts) == 0 {
					err = nil
					break
				}
			} else {
				return nil, fmt.Errorf("astiav_test: reading frame failed: %w", err)
			}
		} else {
			pkts = append(pkts, pkt2)
		}

		for _, pkt := range pkts {
			s, ok := streams[pkt.StreamIndex()]
			if !ok {
				continue
			}

			if err = s.cc.SendPacket(pkt); err != nil {
				return nil, fmt.Errorf("astiav_test: sending packet failed: %w", err)
			}

			for {
				if err = s.cc.ReceiveFrame(f); err != nil {
					if errors.Is(err, ErrEof) || errors.Is(err, ErrEagain) {
						err = nil
						break
					}
					return nil, fmt.Errorf("astiav_test: receiving frame failed: %w", err)
				}

				h.m.Lock()
				h.inputs[name].lastFrames[s.cc.MediaType()].Unref()
				err = h.inputs[name].lastFrames[s.cc.MediaType()].Ref(f)
				h.m.Unlock()
				if err != nil {
					return nil, fmt.Errorf("astiav_test: refing frame failed: %w", err)
				}
			}
		}

		pkts = []*Packet{}
	}
	return h.inputs[name].lastFrames[mediaType], nil
}
