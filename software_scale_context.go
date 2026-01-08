package astiav

//#include <libswscale/swscale.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html
type SoftwareScaleContext struct {
	c *C.struct_SwsContext
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

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#ga59cc19eff0434e7ec11676dc5e222ff3
func CreateSoftwareScaleContext(srcW, srcH int, srcFormat PixelFormat, dstW, dstH int, dstFormat PixelFormat, flags SoftwareScaleContextFlags) (*SoftwareScaleContext, error) {
	ssc := &SoftwareScaleContext{}
	ssc.c = C.sws_getContext(
		C.int(srcW),
		C.int(srcH),
		C.enum_AVPixelFormat(srcFormat),
		C.int(dstW),
		C.int(dstH),
		C.enum_AVPixelFormat(dstFormat),
		C.int(flags),
		nil, nil, nil,
	)
	if ssc.c == nil {
		return nil, errors.New("astiav: empty new context")
	}

	classers.set(ssc)
	return ssc, nil
}

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#gad90b463ceeafdfd526994742f9954dbb
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

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a6866f52574bc730833d2580abc806261
func (ssc *SoftwareScaleContext) Class() *Class {
	if ssc.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(ssc.c))
}

// https://ffmpeg.org/doxygen/8.0/swscale-v2_8txt.html#a20ffff3ac1378332422b93ed70264f4c
func (ssc *SoftwareScaleContext) ScaleFrame(src, dst *Frame) error {
	return newError(C.sws_scale_frame(ssc.c, dst.c, src.c))
}

func (ssc *SoftwareScaleContext) update(u softwareScaleContextUpdate) error {
	if ssc.c == nil {
		return errors.New("astiav: empty context")
	}

	dstW := ssc.c.dst_w
	if u.dstW != nil {
		dstW = C.int(*u.dstW)
	}

	dstH := ssc.c.dst_h
	if u.dstH != nil {
		dstH = C.int(*u.dstH)
	}

	dstFormat := C.enum_AVPixelFormat(ssc.c.dst_format)
	if u.dstFormat != nil {
		dstFormat = C.enum_AVPixelFormat(*u.dstFormat)
	}

	srcW := ssc.c.src_w
	if u.srcW != nil {
		srcW = C.int(*u.srcW)
	}

	srcH := ssc.c.src_h
	if u.srcH != nil {
		srcH = C.int(*u.srcH)
	}

	srcFormat := C.enum_AVPixelFormat(ssc.c.src_format)
	if u.srcFormat != nil {
		srcFormat = C.enum_AVPixelFormat(*u.srcFormat)
	}

	flags := C.int(ssc.c.flags)
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

	return nil
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aef45de443b59978fd38ad1531c618574
func (ssc *SoftwareScaleContext) Flags() SoftwareScaleContextFlags {
	if ssc.c == nil {
		return 0
	}
	return SoftwareScaleContextFlags(ssc.c.flags)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aef45de443b59978fd38ad1531c618574
func (ssc *SoftwareScaleContext) SetFlags(swscf SoftwareScaleContextFlags) error {
	return ssc.update(softwareScaleContextUpdate{flags: &swscf})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a883a891c8a2d4ea7a15a3a7055f64386
func (ssc *SoftwareScaleContext) DestinationWidth() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.dst_w)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a883a891c8a2d4ea7a15a3a7055f64386
func (ssc *SoftwareScaleContext) SetDestinationWidth(i int) error {
	return ssc.update(softwareScaleContextUpdate{dstW: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a7facd34608c9258dae8c2942e3dce78f
func (ssc *SoftwareScaleContext) DestinationHeight() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.dst_h)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a7facd34608c9258dae8c2942e3dce78f
func (ssc *SoftwareScaleContext) SetDestinationHeight(i int) error {
	return ssc.update(softwareScaleContextUpdate{dstH: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) DestinationPixelFormat() PixelFormat {
	if ssc.c == nil {
		return 0
	}
	return PixelFormat(ssc.c.dst_format)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) SetDestinationPixelFormat(p PixelFormat) error {
	return ssc.update(softwareScaleContextUpdate{dstFormat: &p})
}

func (ssc *SoftwareScaleContext) DestinationResolution() (width int, height int) {
	if ssc.c == nil {
		return 0, 0
	}
	return int(ssc.c.dst_w), int(ssc.c.dst_h)
}

func (ssc *SoftwareScaleContext) SetDestinationResolution(w int, h int) error {
	return ssc.update(softwareScaleContextUpdate{dstW: &w, dstH: &h})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aa7dc7a4f9ec57a7c37957259a51cd920
func (ssc *SoftwareScaleContext) SourceWidth() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.src_w)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) SetSourceWidth(i int) error {
	return ssc.update(softwareScaleContextUpdate{srcW: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0dbc8c02bd3b4cd472e07008009751ff
func (ssc *SoftwareScaleContext) SourceHeight() int {
	if ssc.c == nil {
		return 0
	}
	return int(ssc.c.src_h)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#a0ff71c9ef5ab6dabf90378fa7bf836ec
func (ssc *SoftwareScaleContext) SetSourceHeight(i int) error {
	return ssc.update(softwareScaleContextUpdate{srcH: &i})
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aab113373f157ee3b255ad97481af0cd9
func (ssc *SoftwareScaleContext) SourcePixelFormat() PixelFormat {
	if ssc.c == nil {
		return 0
	}
	return PixelFormat(ssc.c.src_format)
}

// https://ffmpeg.org/doxygen/8.0/structSwsContext.html#aab113373f157ee3b255ad97481af0cd9
func (ssc *SoftwareScaleContext) SetSourcePixelFormat(p PixelFormat) error {
	return ssc.update(softwareScaleContextUpdate{srcFormat: &p})
}

func (ssc *SoftwareScaleContext) SourceResolution() (int, int) {
	if ssc.c == nil {
		return 0, 0
	}
	return int(ssc.c.src_w), int(ssc.c.src_h)
}

func (ssc *SoftwareScaleContext) SetSourceResolution(w int, h int) error {
	return ssc.update(softwareScaleContextUpdate{srcW: &w, srcH: &h})
}
