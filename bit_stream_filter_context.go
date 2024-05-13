package astiav

//#cgo pkg-config: libavcodec
//#include <libavcodec/bsf.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/release/5.1/libavcodec/bsf.h#L68
type BSFContext struct {
	c *C.struct_AVBSFContext
}

func newBSFContextFromC(c *C.struct_AVBSFContext) *BSFContext {
	if c == nil {
		return nil
	}
	bsfCtx := &BSFContext{c: c}
	classers.set(bsfCtx)
	return bsfCtx
}

var _ Classer = (*BSFContext)(nil)

func AllocBitStreamContext(f *BitStreamFilter) (*BSFContext, error) {
	if f == nil {
		return nil, errors.New("astiav: bit stream filter must not be nil")
	}

	var bsfCtx *C.struct_AVBSFContext
	if err := newError(C.av_bsf_alloc(f.c, &bsfCtx)); err != nil {
		return nil, err
	}

	return newBSFContextFromC(bsfCtx), nil
}

func (bsfCtx *BSFContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(bsfCtx.c))
}

func (bsfCtx *BSFContext) Init() error {
	return newError(C.av_bsf_init(bsfCtx.c))
}

func (bsfCtx *BSFContext) SendPacket(p *Packet) error {
	if p == nil {
		return errors.New("astiav: packet must not be nil")
	}
	return newError(C.av_bsf_send_packet(bsfCtx.c, p.c))
}

func (bsfCtx *BSFContext) ReceivePacket(p *Packet) error {
	if p == nil {
		return errors.New("astiav: packet must not be nil")
	}
	return newError(C.av_bsf_receive_packet(bsfCtx.c, p.c))
}

func (bsfCtx *BSFContext) Free() {
	classers.del(bsfCtx)
	C.av_bsf_free(&bsfCtx.c)
}

func (bsfCtx *BSFContext) TimeBaseIn() Rational {
	return newRationalFromC(bsfCtx.c.time_base_in)
}

func (bsfCtx *BSFContext) SetTimeBaseIn(r Rational) {
	bsfCtx.c.time_base_in = r.c
}

func (bsfCtx *BSFContext) CodecParametersIn() *CodecParameters {
	return newCodecParametersFromC(bsfCtx.c.par_in)
}

func (bsfCtx *BSFContext) SetCodecParametersIn(cp *CodecParameters) {
	bsfCtx.c.par_in = cp.c
}
