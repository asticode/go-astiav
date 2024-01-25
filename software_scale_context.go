package astiav

//#cgo pkg-config: libswscale
//#include <libswscale/swscale.h>
import "C"
import (
	"fmt"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libswscale/swscale_internal.h#L300
type SoftwareScaleContext struct {
	c *C.struct_SwsContext
	// We need to store attributes in GO since C attributes are internal and therefore not accessible
	dstFormat C.enum_AVPixelFormat
	dstH      C.int
	dstW      C.int
	flags     SoftwareScaleContextFlags
	srcFormat C.enum_AVPixelFormat
	srcH      C.int
	srcW      C.int
}

func NewSoftwareScaleContext(srcW, srcH int, srcFormat PixelFormat, dstW, dstH int, dstFormat PixelFormat, flags SoftwareScaleContextFlags) *SoftwareScaleContext {
	ssc := SoftwareScaleContext{
		dstFormat: C.enum_AVPixelFormat(dstFormat),
		dstH:      C.int(dstH),
		dstW:      C.int(dstW),
		flags:     flags,
		srcFormat: C.enum_AVPixelFormat(srcFormat),
		srcH:      C.int(srcH),
		srcW:      C.int(srcW),
	}

	ssc.c = C.sws_getContext(
		ssc.srcW,
		ssc.srcH,
		ssc.srcFormat,
		ssc.dstW,
		ssc.dstH,
		ssc.dstFormat,
		C.int(ssc.flags),
		nil, nil, nil,
	)
	if ssc.c == nil {
		return nil
	}
	return &ssc
}

func (ssc *SoftwareScaleContext) ScaleFrame(src, dst *Frame) (height int) {
	height = int(
		C.sws_scale(
			ssc.c,
			&src.c.data[0],
			&src.c.linesize[0],
			0,
			C.int(src.Height()),
			&dst.c.data[0], &dst.c.linesize[0]))
	return
}

func (ssc *SoftwareScaleContext) updateContext() error {
	ssc.c = C.sws_getContext(
		ssc.srcW,
		ssc.srcH,
		ssc.srcFormat,
		ssc.dstW,
		ssc.dstH,
		ssc.dstFormat,
		C.int(ssc.flags),
		nil, nil, nil,
	)
	if ssc.c == nil {
		return fmt.Errorf("failed to update sws context")
	}
	return nil
}

func (ssc *SoftwareScaleContext) PrepareDestinationFrameForScaling(dstFrame *Frame) error {
	dstFrame.SetPixelFormat(PixelFormat(ssc.dstFormat))
	dstFrame.SetWidth(int(ssc.dstW))
	dstFrame.SetHeight(int(ssc.dstH))
	return dstFrame.AllocBuffer(1)
}

func (ssc *SoftwareScaleContext) DestinationHeight() int {
	return int(ssc.dstH)
}

func (ssc *SoftwareScaleContext) DestinationPixelFormat() PixelFormat {
	return PixelFormat(ssc.dstFormat)
}

func (ssc *SoftwareScaleContext) DestinationResolution() (int, int) {
	return int(ssc.dstW), int(ssc.dstH)
}

func (ssc *SoftwareScaleContext) DestinationWidth() int {
	return int(ssc.dstW)
}

func (ssc *SoftwareScaleContext) Flags() SoftwareScaleContextFlags {
	return ssc.flags
}

func (ssc *SoftwareScaleContext) SetDestinationHeight(i int) error {
	ssc.dstH = C.int(i)
	return ssc.updateContext()
}

func (ssc *SoftwareScaleContext) SetDestinationPixelFormat(p PixelFormat) error {
	ssc.dstFormat = C.enum_AVPixelFormat(p)
	return ssc.updateContext()
}

func (ssc *SoftwareScaleContext) SetDestinationResolution(w int, h int) error {
	ssc.dstW = C.int(w)
	ssc.dstH = C.int(h)
	return ssc.updateContext()
}

func (ssc *SoftwareScaleContext) SetDestinationWidth(i int) error {
	ssc.dstW = C.int(i)
	return ssc.updateContext()
}

func (ssc *SoftwareScaleContext) SetFlags(swscf SoftwareScaleContextFlags) {
	ssc.flags = swscf
}

func (ssc *SoftwareScaleContext) SetSourceHeight(i int) error {
	ssc.srcH = C.int(i)
	return ssc.updateContext()
}

func (ssc *SoftwareScaleContext) SetSourcePixelFormat(p PixelFormat) error {
	ssc.srcFormat = C.enum_AVPixelFormat(p)
	return ssc.updateContext()
}

func (ssc *SoftwareScaleContext) SetSourceResolution(w int, h int) error {
	ssc.srcW = C.int(w)
	ssc.srcH = C.int(h)
	return ssc.updateContext()
}

func (ssc *SoftwareScaleContext) SetSourceWidth(i int) error {
	ssc.srcW = C.int(i)
	return ssc.updateContext()
}

func (ssc *SoftwareScaleContext) SourceHeight() int {
	return int(ssc.srcH)
}

func (ssc *SoftwareScaleContext) SourcePixelFormat() PixelFormat {
	return PixelFormat(ssc.srcFormat)
}

func (ssc *SoftwareScaleContext) SourceResolution() (int, int) {
	return int(ssc.srcW), int(ssc.srcH)
}

func (ssc *SoftwareScaleContext) SourceWidth() int {
	return int(ssc.srcW)
}

func (sc *SoftwareScaleContext) Free() {
	C.sws_freeContext(sc.c)
}
