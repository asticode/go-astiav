package astiav

//#include <libavcodec/bsf.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVBSFContext.html
type BitStreamFilterContext struct {
	classerHandler
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

// https://ffmpeg.org/doxygen/7.0/group__lavc__bsf.html#ga7da65af303e20c9546e15ec266b182c1
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

// https://ffmpeg.org/doxygen/7.0/structAVBSFContext.html#aa5d5018816daac804414c459ec8a1c5c
func (bsfc *BitStreamFilterContext) Class() *Class {
	if bsfc.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(bsfc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__bsf.html#ga242529d54013acf87e94273d298a5ff2
func (bsfc *BitStreamFilterContext) Initialize() error {
	return bsfc.newError(C.av_bsf_init(bsfc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__bsf.html#gaada9ea8f08d3dcf23c14564dbc88992c
func (bsfc *BitStreamFilterContext) SendPacket(p *Packet) error {
	var pc *C.AVPacket
	if p != nil {
		pc = p.c
	}
	return bsfc.newError(C.av_bsf_send_packet(bsfc.c, pc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__bsf.html#ga7fffb6c87b91250956e7a2367af56b38
func (bsfc *BitStreamFilterContext) ReceivePacket(p *Packet) error {
	if p == nil {
		return errors.New("astiav: packet must not be nil")
	}
	return bsfc.newError(C.av_bsf_receive_packet(bsfc.c, p.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__bsf.html#ga08d53431e76355f88e27763b1940df4f
func (bsfc *BitStreamFilterContext) Free() {
	if bsfc.c != nil {
		// Make sure to clone the classer before freeing the object since
		// the C free method may reset the pointer
		c := newClonedClasser(bsfc)
		C.av_bsf_free(&bsfc.c)
		// Make sure to remove from classers after freeing the object since
		// the C free method may use methods needing the classer
		if c != nil {
			classers.del(c)
		}
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVBSFContext.html#ad75adf988c00f89202099c87ea39f0db
func (bsfc *BitStreamFilterContext) InputTimeBase() Rational {
	return newRationalFromC(bsfc.c.time_base_in)
}

// https://ffmpeg.org/doxygen/7.0/structAVBSFContext.html#ad75adf988c00f89202099c87ea39f0db
func (bsfc *BitStreamFilterContext) SetInputTimeBase(r Rational) {
	bsfc.c.time_base_in = r.c
}

// https://ffmpeg.org/doxygen/7.0/structAVBSFContext.html#a702ace639b8193475cf0a12ebdebd738
func (bsfc *BitStreamFilterContext) InputCodecParameters() *CodecParameters {
	return newCodecParametersFromC(bsfc.c.par_in)
}

// https://ffmpeg.org/doxygen/7.0/structAVBSFContext.html#ab58f8c37eec197e0f30d17d60959a60d
func (bsfc *BitStreamFilterContext) OutputCodecParameters() *CodecParameters {
	return newCodecParametersFromC(bsfc.c.par_out)
}
