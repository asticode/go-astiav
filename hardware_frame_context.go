package astiav

//#include <libavcodec/avcodec.h>
import "C"
import (
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n7.0/libavutil/hwcontext.h#L115
type HardwareFrameContext struct {
	c *C.struct_AVBufferRef
}

func newHardwareFrameContextFromC(c *C.struct_AVBufferRef) *HardwareFrameContext {
	if c == nil {
		return nil
	}
	return &HardwareFrameContext{c: c}
}

func AllocHardwareFrameContext(hdc *HardwareDeviceContext) *HardwareFrameContext {
	return newHardwareFrameContextFromC(C.av_hwframe_ctx_alloc(hdc.c))
}

func (hfc *HardwareFrameContext) data() *C.AVHWFramesContext {
	return (*C.AVHWFramesContext)(unsafe.Pointer((hfc.c.data)))
}

func (hfc *HardwareFrameContext) SetWidth(width int) {
	hfc.data().width = C.int(width)
}

func (hfc *HardwareFrameContext) SetHeight(height int) {
	hfc.data().height = C.int(height)
}

func (hfc *HardwareFrameContext) SetPixelFormat(format PixelFormat) {
	hfc.data().format = C.enum_AVPixelFormat(format)
}

func (hfc *HardwareFrameContext) SetSoftwarePixelFormat(swFormat PixelFormat) {
	hfc.data().sw_format = C.enum_AVPixelFormat(swFormat)
}

func (hfc *HardwareFrameContext) SetInitialPoolSize(initialPoolSize int) {
	hfc.data().initial_pool_size = C.int(initialPoolSize)
}

func (hfc *HardwareFrameContext) Initialize() error {
	return newError(C.av_hwframe_ctx_init(hfc.c))
}
