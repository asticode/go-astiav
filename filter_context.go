package astiav

//#include <libavfilter/avfilter.h>
//#include <libavfilter/buffersink.h>
//#include <libavfilter/buffersrc.h>
//#include <libavutil/frame.h>
import "C"
import (
	"math"
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
	classers.del(fc)
	C.avfilter_free(fc.c)
}

func (fc *FilterContext) BuffersrcAddFrame(f *Frame, fs BuffersrcFlags) error {
	var cf *C.AVFrame
	if f != nil {
		cf = f.c
	}
	return newError(C.av_buffersrc_add_frame_flags(fc.c, cf, C.int(fs)))
}

func (fc *FilterContext) BuffersinkGetFrame(f *Frame, fs BuffersinkFlags) error {
	var cf *C.AVFrame
	if f != nil {
		cf = f.c
	}
	return newError(C.av_buffersink_get_frame_flags(fc.c, cf, C.int(fs)))
}

func (fc *FilterContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(fc.c))
}

func (fc *FilterContext) NbInputs() int {
	return int(fc.c.nb_inputs)
}

func (fc *FilterContext) NbOutputs() int {
	return int(fc.c.nb_outputs)
}

func (fc *FilterContext) Inputs() (ls []*FilterLink) {
	lcs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVFilterLink)(nil))](*C.AVFilterLink))(unsafe.Pointer(fc.c.inputs))
	for i := 0; i < fc.NbInputs(); i++ {
		ls = append(ls, newFilterLinkFromC(lcs[i]))
	}
	return
}

func (fc *FilterContext) Outputs() (ls []*FilterLink) {
	lcs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVFilterLink)(nil))](*C.AVFilterLink))(unsafe.Pointer(fc.c.outputs))
	for i := 0; i < fc.NbOutputs(); i++ {
		ls = append(ls, newFilterLinkFromC(lcs[i]))
	}
	return
}
