package astiav

//#include <libavcodec/avcodec.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type HardwareFrameContext struct {
	c *C.struct_AVBufferRef
}

func AllocHardwareFrameContext(hdc *HardwareDeviceContext) (*HardwareFrameContext, error) {
	if hfc := C.av_hwframe_ctx_alloc(hdc.c); hfc != nil {
		return &HardwareFrameContext{c: hfc}, nil
	}
	return nil, fmt.Errorf("failed to allocate hardware frame context")
}

func (hfc *HardwareFrameContext) SetWidth(width int) {
	frameCtx := (*C.AVHWFramesContext)(unsafe.Pointer((hfc.c.data)))
	frameCtx.width = C.int(width)
}

func (hfc *HardwareFrameContext) SetHeight(height int) {
	frameCtx := (*C.AVHWFramesContext)(unsafe.Pointer((hfc.c.data)))
	frameCtx.height = C.int(height)
}

func (hfc *HardwareFrameContext) SetFormat(format PixelFormat) {
	frameCtx := (*C.AVHWFramesContext)(unsafe.Pointer((hfc.c.data)))
	frameCtx.format = C.enum_AVPixelFormat(format)
}

func (hfc *HardwareFrameContext) SetSWFormat(swFormat PixelFormat) {
	frameCtx := (*C.AVHWFramesContext)(unsafe.Pointer((hfc.c.data)))
	frameCtx.sw_format = C.enum_AVPixelFormat(swFormat)
}

func (hfc *HardwareFrameContext) SetInitialPoolSize(initialPoolSize int) {
	frameCtx := (*C.AVHWFramesContext)(unsafe.Pointer((hfc.c.data)))
	frameCtx.initial_pool_size = C.int(initialPoolSize)
}

func (hfc *HardwareFrameContext) Init() error {
	return newError(C.av_hwframe_ctx_init(hfc.c))
}
