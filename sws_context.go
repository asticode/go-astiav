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
	flags     ScalingAlgorithm
	dstFrame  *Frame
}

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libswscale/swscale.h#L59
type ScalingAlgorithm int

const (
	SWS_FAST_BILINEAR ScalingAlgorithm = ScalingAlgorithm(C.SWS_FAST_BILINEAR)
	SWS_BILINEAR      ScalingAlgorithm = ScalingAlgorithm(C.SWS_BILINEAR)
	SWS_BICUBIC       ScalingAlgorithm = ScalingAlgorithm(C.SWS_BICUBIC)
	SWS_X             ScalingAlgorithm = ScalingAlgorithm(C.SWS_X)
	SWS_POINT         ScalingAlgorithm = ScalingAlgorithm(C.SWS_POINT)
	SWS_AREA          ScalingAlgorithm = ScalingAlgorithm(C.SWS_AREA)
	SWS_BICUBLIN      ScalingAlgorithm = ScalingAlgorithm(C.SWS_BICUBLIN)
	SWS_GAUSS         ScalingAlgorithm = ScalingAlgorithm(C.SWS_GAUSS)
	SWS_SINC          ScalingAlgorithm = ScalingAlgorithm(C.SWS_SINC)
	SWS_LANCZOS       ScalingAlgorithm = ScalingAlgorithm(C.SWS_LANCZOS)
	SWS_SPLINE        ScalingAlgorithm = ScalingAlgorithm(C.SWS_SPLINE)
)

func SwsGetContext(srcW, srcH int, srcFormat PixelFormat, dstW, dstH int, dstFormat PixelFormat, flags ScalingAlgorithm, dstFrame *Frame) *SWSContext {
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

func (sc *SWSContext) UpdateScalingParameters(dstW, dstH int, dstFormat PixelFormat) error {
	if sc.dstW != dstW || sc.dstH != dstH || sc.dstFormat != dstFormat {
		sc.dstW = dstW
		sc.dstH = dstH
		sc.dstFormat = dstFormat

		// Reallocate the destination frame buffer
		sc.dstFrame.SetPixelFormat(dstFormat)
		sc.dstFrame.SetWidth(dstW)
		sc.dstFrame.SetHeight(dstH)
		sc.dstFrame.AllocBuffer(1)

		// Update the sws context
		sc.c = C.sws_getCachedContext(
			sc.c,
			C.int(sc.srcW),
			C.int(sc.srcH),
			C.enum_AVPixelFormat(sc.srcFormat),
			C.int(dstW),
			C.int(dstH),
			C.enum_AVPixelFormat(dstFormat),
			C.int(sc.flags),
			nil, nil, nil,
		)
		if sc.c == nil {
			return fmt.Errorf("failed to update sws context")
		}
	}
	return nil
}

func (sc *SWSContext) Free() {
	C.sws_freeContext(sc.c)
}
