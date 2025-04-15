package astiav

//#include <libswscale/swscale.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html
type SoftwareScaleContext struct {
	classerHandler
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

// https://ffmpeg.org/doxygen/7.0/group__libsws.html#gaf360d1a9e0e60f906f74d7d44f9abfdd
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

// https://ffmpeg.org/doxygen/7.0/group__libsws.html#gad3af0ca76f071dbe0173444db9882932
func (ssc *SoftwareScaleContext) Free() {
	if ssc.c != nil {
		// Make sure to clone the classer before freeing the object since
		// the C free method may reset the pointer
		c := newClonedClasser(ssc)
		C.sws_freeContext(ssc.c)
		ssc.c = nil
		// Make sure to remove from classers after freeing the object since
		// the C free method may use methods needing the classer
		if c != nil {
			classers.del(c)
		}

	}
}

var _ Classer = (*SoftwareScaleContext)(nil)

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html#a6866f52574bc730833d2580abc806261
func (ssc *SoftwareScaleContext) Class() *Class {
	if ssc.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(ssc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__libsws.html#ga1c72fcf83bd57aea72cf3dadfcf02541
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

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html#a4bad42bdd38e916f045956efe81039bf
func (ssc *SoftwareScaleContext) Flags() SoftwareScaleContextFlags {
	return SoftwareScaleContextFlags(ssc.flags)
}

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html#a4bad42bdd38e916f045956efe81039bf
func (ssc *SoftwareScaleContext) SetFlags(swscf SoftwareScaleContextFlags) error {
	return ssc.update(softwareScaleContextUpdate{flags: &swscf})
}

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html#a0bf831c04c58c12ea7aef32e0ffb2f6d
func (ssc *SoftwareScaleContext) DestinationWidth() int {
	return int(ssc.dstW)
}

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html#a0bf831c04c58c12ea7aef32e0ffb2f6d
func (ssc *SoftwareScaleContext) SetDestinationWidth(i int) error {
	return ssc.update(softwareScaleContextUpdate{dstW: &i})
}

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html#a195564564eff11e1ee181999c13b9a22
func (ssc *SoftwareScaleContext) DestinationHeight() int {
	return int(ssc.dstH)
}

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html#a195564564eff11e1ee181999c13b9a22
func (ssc *SoftwareScaleContext) SetDestinationHeight(i int) error {
	return ssc.update(softwareScaleContextUpdate{dstH: &i})
}

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html#a66adb1238b56e3539ad1145c146348e2
func (ssc *SoftwareScaleContext) DestinationPixelFormat() PixelFormat {
	return PixelFormat(ssc.dstFormat)
}

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html#a66adb1238b56e3539ad1145c146348e2
func (ssc *SoftwareScaleContext) SetDestinationPixelFormat(p PixelFormat) error {
	return ssc.update(softwareScaleContextUpdate{dstFormat: &p})
}

func (ssc *SoftwareScaleContext) DestinationResolution() (width int, height int) {
	return int(ssc.dstW), int(ssc.dstH)
}

func (ssc *SoftwareScaleContext) SetDestinationResolution(w int, h int) error {
	return ssc.update(softwareScaleContextUpdate{dstW: &w, dstH: &h})
}

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html#a1e1455f5d751e9ca639bf8afbda25646
func (ssc *SoftwareScaleContext) SourceWidth() int {
	return int(ssc.srcW)
}

// https://ffmpeg.org/doxygen/7.0/structSwsContext.html#a1e1455f5d751e9ca639bf8afbda25646
func (ssc *SoftwareScaleContext) SetSourceWidth(i int) error {
	return ssc.update(softwareScaleContextUpdate{srcW: &i})
}

// https://ffmpeg.org/doxygensrcH/7.0/structSwsContext.html#a195564564eff11e1ee181999c13b9a22
func (ssc *SoftwareScaleContext) SourceHeight() int {
	return int(ssc.srcH)
}

// https://ffmpeg.org/doxygensrcH/7.0/structSwsContext.html#a195564564eff11e1ee181999c13b9a22
func (ssc *SoftwareScaleContext) SetSourceHeight(i int) error {
	return ssc.update(softwareScaleContextUpdate{srcH: &i})
}

// https://ffmpeg.org/doxygensrcH/7.0/structSwsContext.html#a195564564eff11e1ee181999c13b9a22
func (ssc *SoftwareScaleContext) SourcePixelFormat() PixelFormat {
	return PixelFormat(ssc.srcFormat)
}

// https://ffmpeg.org/doxygensrcH/7.0/structSwsContext.html#a195564564eff11e1ee181999c13b9a22
func (ssc *SoftwareScaleContext) SetSourcePixelFormat(p PixelFormat) error {
	return ssc.update(softwareScaleContextUpdate{srcFormat: &p})
}

func (ssc *SoftwareScaleContext) SourceResolution() (int, int) {
	return int(ssc.srcW), int(ssc.srcH)
}

func (ssc *SoftwareScaleContext) SetSourceResolution(w int, h int) error {
	return ssc.update(softwareScaleContextUpdate{srcW: &w, srcH: &h})
}
