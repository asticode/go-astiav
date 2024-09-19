package astiav

//#include <libavfilter/avfilter.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavfilter/avfilter.h#L471
type FilterLink struct {
	c *C.AVFilterLink
}

func newFilterLinkFromC(c *C.AVFilterLink) *FilterLink {
	if c == nil {
		return nil
	}
	return &FilterLink{c: c}
}

func (l *FilterLink) ChannelLayout() ChannelLayout {
	v, _ := newChannelLayoutFromC(&l.c.ch_layout).clone()
	return v
}

func (l *FilterLink) ColorRange() ColorRange {
	return ColorRange(l.c.color_range)
}

func (l *FilterLink) ColorSpace() ColorSpace {
	return ColorSpace(l.c.colorspace)
}

func (l *FilterLink) FrameRate() Rational {
	return newRationalFromC(l.c.frame_rate)
}

func (l *FilterLink) Height() int {
	return int(l.c.h)
}

func (l *FilterLink) MediaType() MediaType {
	return MediaType(l.c._type)
}

func (l *FilterLink) PixelFormat() PixelFormat {
	return PixelFormat(l.c.format)
}

func (l *FilterLink) SampleAspectRatio() Rational {
	return newRationalFromC(l.c.sample_aspect_ratio)
}

func (l *FilterLink) SampleFormat() SampleFormat {
	return SampleFormat(l.c.format)
}

func (l *FilterLink) SampleRate() int {
	return int(l.c.sample_rate)
}

func (l *FilterLink) TimeBase() Rational {
	return newRationalFromC(l.c.time_base)
}

func (l *FilterLink) Width() int {
	return int(l.c.w)
}
