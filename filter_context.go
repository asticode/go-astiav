package astiav

//#include <libavfilter/avfilter.h>
//#include <libavfilter/buffersink.h>
//#include <libavfilter/buffersrc.h>
//#include <libavutil/frame.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVFilterContext.html
type FilterContext struct {
	c *C.AVFilterContext
}

func newFilterContext(c *C.AVFilterContext) *FilterContext {
	if c == nil {
		return nil
	}
	fc := &FilterContext{c: c}
	classers.set(fc)
	return fc
}

var _ Classer = (*FilterContext)(nil)

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#ga0ea7664a3ce6bb677a830698d358a179
func (fc *FilterContext) Free() {
	// Make sure to clone the classer before freeing the object since
	// the C free method may reset the pointer
	c := newClonedClasser(fc)
	C.avfilter_free(fc.c)
	// Make sure to remove from classers after freeing the object since
	// the C free method may use methods needing the classer
	if c != nil {
		classers.del(c)
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterContext.html#a00ac82b13bb720349c138310f98874ca
func (fc *FilterContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(fc.c))
}

type BuffersinkFilterContext struct {
	fc *FilterContext
}

func newBuffersinkFilterContext(fc *FilterContext) *BuffersinkFilterContext {
	return &BuffersinkFilterContext{fc: fc}
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gaad918036937648701c09f9612f42706e
func (bfc *BuffersinkFilterContext) ChannelLayout() ChannelLayout {
	var cl C.AVChannelLayout
	C.av_buffersink_get_ch_layout(bfc.fc.c, &cl)
	return newChannelLayoutFromC(&cl)
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gab80976e506ab88d23d94bb6d7a4051bd
func (bfc *BuffersinkFilterContext) ColorRange() ColorRange {
	return ColorRange(C.av_buffersink_get_color_range(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gaad817cdcf5493c385126e8e17c5717f2
func (bfc *BuffersinkFilterContext) ColorSpace() ColorSpace {
	return ColorSpace(C.av_buffersink_get_colorspace(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) FilterContext() *FilterContext {
	return bfc.fc
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga55614fd28de2fa05b04f427390061d5b
func (bfc *BuffersinkFilterContext) FrameRate() Rational {
	return newRationalFromC(C.av_buffersink_get_frame_rate(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink.html#ga71ae9c529c8da51681e12faa37d1a395
func (bfc *BuffersinkFilterContext) GetFrame(f *Frame, fs BuffersinkFlags) error {
	var cf *C.AVFrame
	if f != nil {
		cf = f.c
	}
	return newError(C.av_buffersink_get_frame_flags(bfc.fc.c, cf, C.int(fs)))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga955ecf3680e71e10429d7500343be25c
func (bfc *BuffersinkFilterContext) Height() int {
	return int(C.av_buffersink_get_h(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga1eb8bbf583ffb7cc29aaa1944b1e699c
func (bfc *BuffersinkFilterContext) MediaType() MediaType {
	return MediaType(C.av_buffersink_get_type(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga402ddbef6f7347869725696846ac81eb
func (bfc *BuffersinkFilterContext) PixelFormat() PixelFormat {
	return PixelFormat(C.av_buffersink_get_format(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gaa38ee33e1c7f6f7cb190bd2330e5f848
func (bfc *BuffersinkFilterContext) SampleAspectRatio() Rational {
	return newRationalFromC(C.av_buffersink_get_sample_aspect_ratio(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga402ddbef6f7347869725696846ac81eb
func (bfc *BuffersinkFilterContext) SampleFormat() SampleFormat {
	return SampleFormat(C.av_buffersink_get_format(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga2af714e82f48759551acdbc4488ded4a
func (bfc *BuffersinkFilterContext) SampleRate() int {
	return int(C.av_buffersink_get_sample_rate(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gabc82f65ec7f4fa47c5216260639258a1
func (bfc *BuffersinkFilterContext) TimeBase() Rational {
	return newRationalFromC(C.av_buffersink_get_time_base(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gac8c86515d2ef56090395dfd74854c835
func (bfc *BuffersinkFilterContext) Width() int {
	return int(C.av_buffersink_get_w(bfc.fc.c))
}

type BuffersrcFilterContext struct {
	fc *FilterContext
}

func newBuffersrcFilterContext(fc *FilterContext) *BuffersrcFilterContext {
	return &BuffersrcFilterContext{fc: fc}
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersrc.html#ga73ed90c3c3407f36e54d65f91faaaed9
func (bfc *BuffersrcFilterContext) AddFrame(f *Frame, fs BuffersrcFlags) error {
	var cf *C.AVFrame
	if f != nil {
		cf = f.c
	}
	return newError(C.av_buffersrc_add_frame_flags(bfc.fc.c, cf, C.int(fs)))
}

func (bfc *BuffersrcFilterContext) FilterContext() *FilterContext {
	return bfc.fc
}
