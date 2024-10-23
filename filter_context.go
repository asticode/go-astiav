package astiav

//#include <libavfilter/avfilter.h>
//#include <libavfilter/buffersink.h>
//#include <libavfilter/buffersrc.h>
//#include <libavutil/frame.h>
import "C"
import (
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavfilter/avfilter.h#L67
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

func (fc *FilterContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(fc.c))
}

type BuffersinkFilterContext struct {
	fc *FilterContext
}

func newBuffersinkFilterContext(fc *FilterContext) *BuffersinkFilterContext {
	return &BuffersinkFilterContext{fc: fc}
}

func (bfc *BuffersinkFilterContext) ChannelLayout() ChannelLayout {
	var cl C.AVChannelLayout
	C.av_buffersink_get_ch_layout(bfc.fc.c, &cl)
	return newChannelLayoutFromC(&cl)
}

func (bfc *BuffersinkFilterContext) ColorRange() ColorRange {
	return ColorRange(C.av_buffersink_get_color_range(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) ColorSpace() ColorSpace {
	return ColorSpace(C.av_buffersink_get_colorspace(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) FilterContext() *FilterContext {
	return bfc.fc
}

func (bfc *BuffersinkFilterContext) FrameRate() Rational {
	return newRationalFromC(C.av_buffersink_get_frame_rate(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) GetFrame(f *Frame, fs BuffersinkFlags) error {
	var cf *C.AVFrame
	if f != nil {
		cf = f.c
	}
	return newError(C.av_buffersink_get_frame_flags(bfc.fc.c, cf, C.int(fs)))
}

func (bfc *BuffersinkFilterContext) Height() int {
	return int(C.av_buffersink_get_h(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) MediaType() MediaType {
	return MediaType(C.av_buffersink_get_type(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) PixelFormat() PixelFormat {
	return PixelFormat(C.av_buffersink_get_format(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) SampleAspectRatio() Rational {
	return newRationalFromC(C.av_buffersink_get_sample_aspect_ratio(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) SampleFormat() SampleFormat {
	return SampleFormat(C.av_buffersink_get_format(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) SampleRate() int {
	return int(C.av_buffersink_get_sample_rate(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) TimeBase() Rational {
	return newRationalFromC(C.av_buffersink_get_time_base(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) Width() int {
	return int(C.av_buffersink_get_w(bfc.fc.c))
}

type BuffersrcFilterContext struct {
	fc *FilterContext
}

func newBuffersrcFilterContext(fc *FilterContext) *BuffersrcFilterContext {
	return &BuffersrcFilterContext{fc: fc}
}

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
