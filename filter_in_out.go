package astiav

//#include <libavfilter/avfilter.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavfilter/avfilter.h#L1021
type FilterInOut struct {
	c *C.AVFilterInOut
}

func newFilterInOutFromC(c *C.AVFilterInOut) *FilterInOut {
	if c == nil {
		return nil
	}
	return &FilterInOut{c: c}
}

func AllocFilterInOut() *FilterInOut {
	return newFilterInOutFromC(C.avfilter_inout_alloc())
}

func (i *FilterInOut) Free() {
	C.avfilter_inout_free(&i.c)
}

func (i *FilterInOut) SetName(n string) {
	i.c.name = C.CString(n)
}

func (i *FilterInOut) SetFilterContext(c *FilterContext) {
	i.c.filter_ctx = (*C.AVFilterContext)(c.c)
}

func (i *FilterInOut) SetPadIdx(idx int) {
	i.c.pad_idx = C.int(idx)
}

func (i *FilterInOut) SetNext(n *FilterInOut) {
	var nc *C.AVFilterInOut
	if n != nil {
		nc = n.c
	}
	i.c.next = nc
}
