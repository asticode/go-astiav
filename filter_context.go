package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVFilterContext.html
type FilterContext struct {
	classerHandler
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
	if fc.c != nil {
		// Make sure to clone the classer before freeing the object since
		// the C free method may reset the pointer
		c := newClonedClasser(fc)
		C.avfilter_free(fc.c)
		fc.c = nil
		// Make sure to remove from classers after freeing the object since
		// the C free method may use methods needing the classer
		if c != nil {
			classers.del(c)
		}
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterContext.html#a00ac82b13bb720349c138310f98874ca
func (fc *FilterContext) Class() *Class {
	if fc.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(fc.c))
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterContext.html#addd946fbe5af506a2b19f9ad7cb97c35
func (fc *FilterContext) SetHardwareDeviceContext(hdc *HardwareDeviceContext) {
	if fc.c.hw_device_ctx != nil {
		C.av_buffer_unref(&fc.c.hw_device_ctx)
	}
	if hdc != nil {
		fc.c.hw_device_ctx = C.av_buffer_ref(hdc.c)
	} else {
		fc.c.hw_device_ctx = nil
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterContext.html#a6eee53e57dddfa7cca1cade870c8a44e
func (fc *FilterContext) Filter() *Filter {
	return newFilterFromC(fc.c.filter)
}
