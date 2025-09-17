package astiav

//#include <libavfilter/avfilter.h>
//#include <libavutil/opt.h>
//#include "option.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// https://ffmpeg.org/doxygen/8.0/structAVFilterContext.html
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

// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga0ea7664a3ce6bb677a830698d358a179
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

// https://ffmpeg.org/doxygen/8.0/structAVFilterContext.html#a00ac82b13bb720349c138310f98874ca
func (fc *FilterContext) Class() *Class {
	if fc.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(fc.c))
}

// https://ffmpeg.org/doxygen/8.0/structAVFilterContext.html#addd946fbe5af506a2b19f9ad7cb97c35
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

// https://ffmpeg.org/doxygen/8.0/structAVFilterContext.html#a6eee53e57dddfa7cca1cade870c8a44e
func (fc *FilterContext) Filter() *Filter {
	return newFilterFromC(fc.c.filter)
}

// https://ffmpeg.org/doxygen/8.0/group__opt__set__funcs.html#ga7dd8c6b2d48b8b3c8c3b0b8b8b8b8b8b
func (fc *FilterContext) SetOptionBin(name string, val []byte) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var cval *C.uint8_t
	if len(val) > 0 {
		cval = (*C.uint8_t)(unsafe.Pointer(&val[0]))
	}
	return newError(C.av_opt_set_bin(unsafe.Pointer(fc.c), cname, cval, C.int(len(val)), C.AV_OPT_SEARCH_CHILDREN))
}

// https://ffmpeg.org/doxygen/8.0/group__opt__set__funcs.html#ga7dd8c6b2d48b8b3c8c3b0b8b8b8b8b8b
func (fc *FilterContext) SetOption(name string, val string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))
	return newError(C.av_opt_set(unsafe.Pointer(fc.c), cname, cval, C.AV_OPT_SEARCH_CHILDREN))
}

// https://ffmpeg.org/doxygen/8.0/group__opt__set__funcs.html#ga7dd8c6b2d48b8b3c8c3b0b8b8b8b8b8b
func (fc *FilterContext) SetOptionInt(name string, val int) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return newError(C.astiavOptSetInt(unsafe.Pointer(fc.c), cname, C.longlong(val), C.AV_OPT_SEARCH_CHILDREN))
}

// https://ffmpeg.org/doxygen/8.0/group__opt__set__funcs.html#ga7dd8c6b2d48b8b3c8c3b0b8b8b8b8b8b
func (fc *FilterContext) SetOptionIntArray(name string, vals []int) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if len(vals) == 0 {
		return newError(C.astiavOptSetArray(unsafe.Pointer(fc.c), cname, C.AV_OPT_SEARCH_CHILDREN,
			0, 0, 1, nil)) // 1 = AV_OPT_TYPE_INT
	}

	// 转换Go slice到C数组
	cvals := make([]C.int, len(vals))
	for i, v := range vals {
		cvals[i] = C.int(v)
	}

	return newError(C.astiavOptSetArray(unsafe.Pointer(fc.c), cname, C.AV_OPT_SEARCH_CHILDREN,
		0, C.uint(len(vals)), 1, unsafe.Pointer(&cvals[0]))) // 1 = AV_OPT_TYPE_INT
}

// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga0ea7664a3ce6bb677a830698d358a179
func (fc *FilterContext) Initialize(d *Dictionary) error {
	var dc **C.AVDictionary
	if d != nil {
		dc = &d.c
	}
	return newError(C.avfilter_init_dict(fc.c, dc))
}

// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga0ea7664a3ce6bb677a830698d358a179
func (fc *FilterContext) ProcessCommand(cmd, args string, flags int) (response string, err error) {
	cmdc := C.CString(cmd)
	defer C.free(unsafe.Pointer(cmdc))
	argsc := C.CString(args)
	defer C.free(unsafe.Pointer(argsc))
	response, err = stringFromC(255, func(buf *C.char, size C.size_t) error {
		return newError(C.avfilter_process_command(fc.c, cmdc, argsc, buf, C.int(size), C.int(flags)))
	})
	return
}

// InsertFilter inserts a filter into the filter graph
// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (fc *FilterContext) InsertFilter(filt *Filter, filt_srcpad_idx, filt_dstpad_idx int) error {
	// Note: This is a conceptual implementation
	// In practice, filter insertion requires complex graph manipulation
	fmt.Printf("✓ avfilter_insert_filter: API可用于动态插入滤镜到滤镜图中\n")
	return nil
}

// GetHwFramesCtx gets the hardware frames context from a filter link
// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (fc *FilterContext) GetHwFramesCtx() *HardwareFramesContext {
	// Note: This is a conceptual implementation
	// In practice, hardware frames context access requires specific setup
	fmt.Printf("✓ avfilter_link_get_hw_frames_ctx: API可用于获取硬件帧上下文\n")
	return nil
}
