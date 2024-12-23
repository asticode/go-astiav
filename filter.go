package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVFilter.html
type Filter struct {
	c *C.AVFilter
}

func newFilterFromC(c *C.AVFilter) *Filter {
	if c == nil {
		return nil
	}
	return &Filter{c: c}
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#gadd774ec49e50edf00158248e1bfe4ae6
func FindFilterByName(n string) *Filter {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newFilterFromC(C.avfilter_get_by_name(cn))
}

// https://ffmpeg.org/doxygen/7.0/structAVFilter.html#a632c76418742ad4f4dccbd4db40badd0
func (f *Filter) Flags() FilterFlags {
	return FilterFlags(f.c.flags)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilter.html#a28a4776f344f91055f42a4c2a1b15c0c
func (f *Filter) Name() string {
	return C.GoString(f.c.name)
}

func (f *Filter) String() string {
	return f.Name()
}

// https://ffmpeg.org/doxygen/7.0/structAVFilter.html#a04e408702054370fbe35c8d5b49f68cb
func (f *Filter) NbInputs() int {
	return int(f.c.nb_inputs)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilter.html#abb166cb2c9349de54d24aefb879608f4
func (f *Filter) NbOutputs() int {
	return int(f.c.nb_outputs)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilter.html#ad311151fe6e8c87a89f895bef7c8b98b
func (f *Filter) Inputs() (ps []*FilterPad) {
	for idx := 0; idx < f.NbInputs(); idx++ {
		ps = append(ps, newFilterPad(MediaType(C.avfilter_pad_get_type(f.c.inputs, C.int(idx)))))
	}
	return
}

// https://ffmpeg.org/doxygen/7.0/structAVFilter.html#ad0608786fa3e1ca6e4cc4b67039f77d7
func (f *Filter) Outputs() (ps []*FilterPad) {
	for idx := 0; idx < f.NbOutputs(); idx++ {
		ps = append(ps, newFilterPad(MediaType(C.avfilter_pad_get_type(f.c.outputs, C.int(idx)))))
	}
	return
}
