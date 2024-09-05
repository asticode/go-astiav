package astiav

//#include <libavcodec/bsf.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/release/5.1/libavcodec/bsf.h#L68
type BitStreamFilterContext struct {
	c *C.AVBSFContext
}

func newBSFContextFromC(c *C.AVBSFContext) *BitStreamFilterContext {
	if c == nil {
		return nil
	}
	bsfc := &BitStreamFilterContext{c: c}
	classers.set(bsfc)
	return bsfc
}

var _ Classer = (*BitStreamFilterContext)(nil)

func AllocBitStreamFilterContext(f *BitStreamFilter) (*BitStreamFilterContext, error) {
	if f == nil {
		return nil, errors.New("astiav: bit stream filter must not be nil")
	}

	var bsfc *C.AVBSFContext
	if err := newError(C.av_bsf_alloc(f.c, &bsfc)); err != nil {
		return nil, err
	}

	return newBSFContextFromC(bsfc), nil
}

func (bsfc *BitStreamFilterContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(bsfc.c))
}

func (bsfc *BitStreamFilterContext) Initialize() error {
	return newError(C.av_bsf_init(bsfc.c))
}

func (bsfc *BitStreamFilterContext) SendPacket(p *Packet) error {
	var pc *C.AVPacket
	if p != nil {
		pc = p.c
	}
	return newError(C.av_bsf_send_packet(bsfc.c, pc))
}

func (bsfc *BitStreamFilterContext) ReceivePacket(p *Packet) error {
	if p == nil {
		return errors.New("astiav: packet must not be nil")
	}
	return newError(C.av_bsf_receive_packet(bsfc.c, p.c))
}

func (bsfc *BitStreamFilterContext) Free() {
	classers.del(bsfc)
	C.av_bsf_free(&bsfc.c)
}

func (bsfc *BitStreamFilterContext) InputTimeBase() Rational {
	return newRationalFromC(bsfc.c.time_base_in)
}

func (bsfc *BitStreamFilterContext) SetInputTimeBase(r Rational) {
	bsfc.c.time_base_in = r.c
}

func (bsfc *BitStreamFilterContext) InputCodecParameters() *CodecParameters {
	return newCodecParametersFromC(bsfc.c.par_in)
}
