package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVFilterContext.html
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

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#ga0ea7664a3ce6bb677a830698d358a179
func (fc *FilterContext) Free() {
	// Make sure to clone the classer before freeing the object since
	// the C free method may reset the pointer
	c := newClonedClasser(fc)
	C.avfilter_free(fc.c)
	// Make sure to remove from classers after freeing the object since
	// the C free method may use methods needing the classer
	if c != nil {
		classers.del(c)
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterContext.html#a00ac82b13bb720349c138310f98874ca
func (fc *FilterContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(fc.c))
}
