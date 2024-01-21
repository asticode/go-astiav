package astiav

//#cgo pkg-config: libswscale
//#include <libswscale/swscale.h>
import "C"
import (
	"fmt"
)

// https://github.com/FFmpeg/FFmpeg/blob/n4.2.7/libswscale/swscale_internal.h#L280
type SWSContext struct {
	c         *C.struct_SwsContext
	dstFormat PixelFormat
	srcFormat PixelFormat
	srcW      int
	srcH      int
	dstW      int
	dstH      int
	flags     int
	dstFrame  *Frame
}

const (
	SWS_FAST_BILINEAR = C.SWS_FAST_BILINEAR
	SWS_BILINEAR      = C.SWS_BILINEAR
	SWS_BICUBIC       = C.SWS_BICUBIC
	SWS_X             = C.SWS_X
	SWS_POINT         = C.SWS_POINT
	SWS_AREA          = C.SWS_AREA
	SWS_BICUBLIN      = C.SWS_BICUBLIN
	SWS_GAUSS         = C.SWS_GAUSS
	SWS_SINC          = C.SWS_SINC
	SWS_LANCZOS       = C.SWS_LANCZOS
	SWS_SPLINE        = C.SWS_SPLINE
)

func CreateSwsContext(srcW, srcH int, srcFormat PixelFormat, dstW, dstH int, dstFormat PixelFormat, flags int, dstFrame *Frame) *SWSContext {
	dstFrame.SetPixelFormat(dstFormat)
	dstFrame.SetWidth(dstW)
	dstFrame.SetHeight(dstH)
	dstFrame.AllocBuffer(1)

	swsCtx := C.sws_getContext(
		C.int(srcW),
		C.int(srcH),
		C.enum_AVPixelFormat(srcFormat),
		C.int(dstW),
		C.int(dstH),
		C.enum_AVPixelFormat(dstFormat),
		C.int(flags),
		nil, nil, nil,
	)
	if swsCtx == nil {
		return nil
	}
	return &SWSContext{c: swsCtx, dstFormat: dstFormat, srcFormat: srcFormat, srcW: srcW, srcH: srcH, dstW: dstW, dstH: dstH, flags: flags, dstFrame: dstFrame}
}

func (sc *SWSContext) ChangeResolution(dstW, dstH int) *SWSContext {
	sc.Free()
	return CreateSwsContext(sc.srcW, sc.srcH, sc.srcFormat, dstW, dstH, sc.dstFormat, sc.flags, sc.dstFrame)
}

func (sc *SWSContext) Scale(srcFrame, dstFrame *Frame) error {
	height := int(
		C.sws_scale(
			sc.c,
			&srcFrame.c.data[0],
			&srcFrame.c.linesize[0],
			0,
			C.int(srcFrame.Height()),
			&dstFrame.c.data[0], &dstFrame.c.linesize[0]))

	if height != dstFrame.Height() {
		return fmt.Errorf("sws_scale did not process all lines, expected: %d, got: %d", dstFrame.Height(), height)
	}
	return nil
}

func (sc *SWSContext) Free() {
	C.sws_freeContext(sc.c)
}
