package astiav

//#include <libavcodec/bsf.h>
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVBitStreamFilter.html
type BitStreamFilter struct {
	c *C.AVBitStreamFilter
}

func newBitStreamFilterFromC(c *C.AVBitStreamFilter) *BitStreamFilter {
	if c == nil {
		return nil
	}
	return &BitStreamFilter{c: c}
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__bsf.html#gae491493190b45698ebd43db28c4e8fe9
func FindBitStreamFilterByName(n string) *BitStreamFilter {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newBitStreamFilterFromC(C.av_bsf_get_by_name(cn))
}

// https://ffmpeg.org/doxygen/7.0/structAVBitStreamFilter.html#a33c3cb51bd13060da35481655b41e4e5
func (bsf *BitStreamFilter) Name() string {
	return C.GoString(bsf.c.name)
}

func (bsf *BitStreamFilter) String() string {
	return bsf.Name()
}
