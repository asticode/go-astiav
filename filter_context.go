package astiav

//#cgo pkg-config: libavfilter libavutil
//#include <libavfilter/avfilter.h>
//#include <libavfilter/buffersink.h>
//#include <libavfilter/buffersrc.h>
//#include <libavutil/frame.h>
import "C"
import "unsafe"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavfilter/avfilter.h#L67
type FilterContext struct {
	c *C.struct_AVFilterContext
}

func newFilterContext() *FilterContext {
	return &FilterContext{}
}

func (fc *FilterContext) Free() {
	C.avfilter_free(fc.c)
}

func (fc *FilterContext) BuffersrcAddFrame(f *Frame, fs BuffersrcFlags) error {
	var cf *C.struct_AVFrame
	if f != nil {
		cf = f.c
	}
	return newError(C.av_buffersrc_add_frame_flags(fc.c, cf, C.int(fs)))
}

func (fc *FilterContext) BuffersinkGetFrame(f *Frame, fs BuffersinkFlags) error {
	var cf *C.struct_AVFrame
	if f != nil {
		cf = f.c
	}
	return newError(C.av_buffersink_get_frame_flags(fc.c, cf, C.int(fs)))
}

func (fc *FilterContext) NbInputs() int {
	return int(fc.c.nb_inputs)
}

func (fc *FilterContext) NbOutputs() int {
	return int(fc.c.nb_outputs)
}

func (fc *FilterContext) Inputs() (ls []*FilterLink) {
	lcs := (*[maxArraySize](*C.struct_AVFilterLink))(unsafe.Pointer(fc.c.inputs))
	for i := 0; i < fc.NbInputs(); i++ {
		ls = append(ls, newFilterLinkFromC(lcs[i]))
	}
	return
}

func (fc *FilterContext) Outputs() (ls []*FilterLink) {
	lcs := (*[maxArraySize](*C.struct_AVFilterLink))(unsafe.Pointer(fc.c.outputs))
	for i := 0; i < fc.NbOutputs(); i++ {
		ls = append(ls, newFilterLinkFromC(lcs[i]))
	}
	return
}
