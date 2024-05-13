package astiav

//#cgo pkg-config: libavcodec
//#include <libavcodec/bsf.h>
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/release/5.1/libavcodec/bsf.h#L111
type BitStreamFilter struct {
	c *C.struct_AVBitStreamFilter
}

func newBitStreamFilterFromC(c *C.struct_AVBitStreamFilter) *BitStreamFilter {
	if c == nil {
		return nil
	}
	cc := &BitStreamFilter{c: c}
	return cc
}

func FindBitStreamFilterByName(n string) *BitStreamFilter {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newBitStreamFilterFromC(C.av_bsf_get_by_name(cn))
}
