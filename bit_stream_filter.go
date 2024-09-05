package astiav

//#include <libavcodec/bsf.h>
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/release/5.1/libavcodec/bsf.h#L111
type BitStreamFilter struct {
	c *C.AVBitStreamFilter
}

func newBitStreamFilterFromC(c *C.AVBitStreamFilter) *BitStreamFilter {
	if c == nil {
		return nil
	}
	return &BitStreamFilter{c: c}
}

func FindBitStreamFilterByName(n string) *BitStreamFilter {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newBitStreamFilterFromC(C.av_bsf_get_by_name(cn))
}

func (bsf *BitStreamFilter) Name() string {
	return C.GoString(bsf.c.name)
}

func (bsf *BitStreamFilter) String() string {
	return bsf.Name()
}
