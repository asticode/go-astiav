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

func (l *FilterLink) TimeBase() Rational {
	return newRationalFromC(l.c.time_base)
}
