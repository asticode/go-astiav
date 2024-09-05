package astiav

//#include <libswscale/swscale.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libswscale/swscale_internal.h#L300
type SoftwareScaleContext struct {
	c *C.struct_SwsContext
	// We need to store attributes in GO since C attributes are internal and therefore not accessible
	dstFormat C.enum_AVPixelFormat
	dstH      C.int
	dstW      C.int
	flags     C.int
	srcFormat C.enum_AVPixelFormat
	srcH      C.int
	srcW      C.int
}

type softwareScaleContextUpdate struct {
	dstFormat *PixelFormat
	dstH      *int
	dstW      *int
	flags     *SoftwareScaleContextFlags
	srcFormat *PixelFormat
	srcH      *int
	srcW      *int
}

func CreateSoftwareScaleContext(srcW, srcH int, srcFormat PixelFormat, dstW, dstH int, dstFormat PixelFormat, flags SoftwareScaleContextFlags) (*SoftwareScaleContext, error) {
	ssc := &SoftwareScaleContext{
		dstFormat: C.enum_AVPixelFormat(dstFormat),
		dstH:      C.int(dstH),
		dstW:      C.int(dstW),
		flags:     C.int(flags),
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
		ssc.flags,
		nil, nil, nil,
	)
	if ssc.c == nil {
		return nil, errors.New("astiav: empty new context")
	}

	classers.set(ssc)
	return ssc, nil
}

func (ssc *SoftwareScaleContext) Free() {
	classers.del(ssc)
	C.sws_freeContext(ssc.c)
}

var _ Classer = (*SoftwareScaleContext)(nil)

func (ssc *SoftwareScaleContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(ssc.c))
}

func (ssc *SoftwareScaleContext) ScaleFrame(src, dst *Frame) error {
	return newError(C.sws_scale_frame(ssc.c, dst.c, src.c))
}

func (ssc *SoftwareScaleContext) update(u softwareScaleContextUpdate) error {
	dstW := ssc.dstW
	if u.dstW != nil {
		dstW = C.int(*u.dstW)
	}

	dstH := ssc.dstH
	if u.dstH != nil {
		dstH = C.int(*u.dstH)
	}

	dstFormat := ssc.dstFormat
	if u.dstFormat != nil {
		dstFormat = C.enum_AVPixelFormat(*u.dstFormat)
	}

	srcW := ssc.srcW
	if u.srcW != nil {
		srcW = C.int(*u.srcW)
	}

	srcH := ssc.srcH
	if u.srcH != nil {
		srcH = C.int(*u.srcH)
	}

	srcFormat := ssc.srcFormat
	if u.srcFormat != nil {
		srcFormat = C.enum_AVPixelFormat(*u.srcFormat)
	}

	flags := ssc.flags
	if u.flags != nil {
		flags = C.int(*u.flags)
	}

	c := C.sws_getCachedContext(
		ssc.c,
		srcW,
		srcH,
		srcFormat,
		dstW,
		dstH,
		dstFormat,
		flags,
		nil, nil, nil,
	)
	if c == nil {
		return errors.New("astiav: empty new context")
	}

	ssc.c = c
	ssc.dstW = dstW
	ssc.dstH = dstH
	ssc.dstFormat = dstFormat
	ssc.srcW = srcW
	ssc.srcH = srcH
	ssc.srcFormat = srcFormat
	ssc.flags = flags

	return nil
}

func (ssc *SoftwareScaleContext) Flags() SoftwareScaleContextFlags {
	return SoftwareScaleContextFlags(ssc.flags)
}

func (ssc *SoftwareScaleContext) SetFlags(swscf SoftwareScaleContextFlags) error {
	return ssc.update(softwareScaleContextUpdate{flags: &swscf})
}

func (ssc *SoftwareScaleContext) DestinationWidth() int {
	return int(ssc.dstW)
}

func (ssc *SoftwareScaleContext) SetDestinationWidth(i int) error {
	return ssc.update(softwareScaleContextUpdate{dstW: &i})
}

func (ssc *SoftwareScaleContext) DestinationHeight() int {
	return int(ssc.dstH)
}

func (ssc *SoftwareScaleContext) SetDestinationHeight(i int) error {
	return ssc.update(softwareScaleContextUpdate{dstH: &i})
}

func (ssc *SoftwareScaleContext) DestinationPixelFormat() PixelFormat {
	return PixelFormat(ssc.dstFormat)
}

func (ssc *SoftwareScaleContext) SetDestinationPixelFormat(p PixelFormat) error {
	return ssc.update(softwareScaleContextUpdate{dstFormat: &p})
}

func (ssc *SoftwareScaleContext) DestinationResolution() (width int, height int) {
	return int(ssc.dstW), int(ssc.dstH)
}

func (ssc *SoftwareScaleContext) SetDestinationResolution(w int, h int) error {
	return ssc.update(softwareScaleContextUpdate{dstW: &w, dstH: &h})
}

func (ssc *SoftwareScaleContext) SourceWidth() int {
	return int(ssc.srcW)
}

func (ssc *SoftwareScaleContext) SetSourceWidth(i int) error {
	return ssc.update(softwareScaleContextUpdate{srcW: &i})
}

func (ssc *SoftwareScaleContext) SourceHeight() int {
	return int(ssc.srcH)
}

func (ssc *SoftwareScaleContext) SetSourceHeight(i int) error {
	return ssc.update(softwareScaleContextUpdate{srcH: &i})
}

func (ssc *SoftwareScaleContext) SourcePixelFormat() PixelFormat {
	return PixelFormat(ssc.srcFormat)
}

func (ssc *SoftwareScaleContext) SetSourcePixelFormat(p PixelFormat) error {
	return ssc.update(softwareScaleContextUpdate{srcFormat: &p})
}

func (ssc *SoftwareScaleContext) SourceResolution() (int, int) {
	return int(ssc.srcW), int(ssc.srcH)
}

func (ssc *SoftwareScaleContext) SetSourceResolution(w int, h int) error {
	return ssc.update(softwareScaleContextUpdate{srcW: &w, srcH: &h})
}
