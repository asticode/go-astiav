package astiav

//#include <libavfilter/avfilter.h>
import "C"
import "unsafe"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavfilter/avfilter.h#L165
type Filter struct {
	c *C.AVFilter
}

func newFilterFromC(c *C.AVFilter) *Filter {
	if c == nil {
		return nil
	}
	return &Filter{c: c}
}

func FindFilterByName(n string) *Filter {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newFilterFromC(C.avfilter_get_by_name(cn))
}

func (f *Filter) Name() string {
	return C.GoString(f.c.name)
}

func (f *Filter) String() string {
	return f.Name()
}
