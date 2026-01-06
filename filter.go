package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.1/structAVFilter.html
type Filter struct {
	c *C.AVFilter
}

func newFilterFromC(c *C.AVFilter) *Filter {
	if c == nil {
		return nil
	}
	return &Filter{c: c}
}

// https://ffmpeg.org/doxygen/7.1/group__lavfi.html#gadd774ec49e50edf00158248e1bfe4ae6
func FindFilterByName(n string) *Filter {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newFilterFromC(C.avfilter_get_by_name(cn))
}

// https://ffmpeg.org/doxygen/7.1/structAVFilter.html#a632c76418742ad4f4dccbd4db40badd0
func (f *Filter) Flags() FilterFlags {
	return FilterFlags(f.c.flags)
}

// https://ffmpeg.org/doxygen/7.1/structAVFilter.html#a28a4776f344f91055f42a4c2a1b15c0c
func (f *Filter) Name() string {
	return C.GoString(f.c.name)
}

// https://ffmpeg.org/doxygen/8.0/structAVFilter.html#afb208213ea814c722279962fb0228241
func (f *Filter) Description() string {
	if f.c.description == nil {
		return ""
	}
	return C.GoString(f.c.description)
}

func (f *Filter) String() string {
	return f.Name()
}

// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga54dd15771603f3406c124259595e142b
func (f *Filter) CountPads(isOutput bool) int {
	if isOutput {
		return int(C.avfilter_filter_pad_count(f.c, 1))
	}
	return int(C.avfilter_filter_pad_count(f.c, 0))
}

// https://ffmpeg.org/doxygen/8.0/structAVFilter.html#ad311151fe6e8c87a89f895bef7c8b98b
func (f *Filter) Inputs() (ps []*FilterPad) {
	for idx := 0; idx < f.CountPads(false); idx++ {
		ps = append(ps, newFilterPad(MediaType(C.avfilter_pad_get_type(f.c.inputs, C.int(idx)))))
	}
	return
}

// https://ffmpeg.org/doxygen/8.0/structAVFilter.html#ad0608786fa3e1ca6e4cc4b67039f77d7
func (f *Filter) Outputs() (ps []*FilterPad) {
	for idx := 0; idx < f.CountPads(true); idx++ {
		ps = append(ps, newFilterPad(MediaType(C.avfilter_pad_get_type(f.c.outputs, C.int(idx)))))
	}
	return
}
