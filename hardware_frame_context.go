package astiav

//#include <libavcodec/avcodec.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVHWFramesContext.html
type HardwareFrameContext struct {
	c *C.struct_AVBufferRef
}

func newHardwareFrameContextFromC(c *C.struct_AVBufferRef) *HardwareFrameContext {
	if c == nil {
		return nil
	}
	return &HardwareFrameContext{c: c}
}

// https://ffmpeg.org/doxygen/7.0/hwcontext_8c.html#ac45a7c039eb4e084b692f69ff5f2e217
func AllocHardwareFrameContext(hdc *HardwareDeviceContext) *HardwareFrameContext {
	return newHardwareFrameContextFromC(C.av_hwframe_ctx_alloc(hdc.c))
}

func (hfc *HardwareFrameContext) Free() {
	if hfc.c != nil {
		C.av_buffer_unref(&hfc.c)
	}
}

func (hfc *HardwareFrameContext) data() *C.AVHWFramesContext {
	return (*C.AVHWFramesContext)(unsafe.Pointer((hfc.c.data)))
}

// https://ffmpeg.org/doxygen/7.0/structAVHWFramesContext.html#a9e6f29d0f744930cdd0e8bdff8771520
func (hfc *HardwareFrameContext) SetWidth(width int) {
	hfc.data().width = C.int(width)
}

// https://ffmpeg.org/doxygen/7.0/structAVHWFramesContext.html#ae61bbe1d8645a0c573085e29f1d0a58f
func (hfc *HardwareFrameContext) SetHeight(height int) {
	hfc.data().height = C.int(height)
}

// https://ffmpeg.org/doxygen/7.0/structAVHWFramesContext.html#a045bc1713932804f6ceef170a5578e0e
func (hfc *HardwareFrameContext) SetPixelFormat(format PixelFormat) {
	hfc.data().format = C.enum_AVPixelFormat(format)
}

// https://ffmpeg.org/doxygen/7.0/structAVHWFramesContext.html#a663a9aceca97aa7b2426c9aba6543e4a
func (hfc *HardwareFrameContext) SetSoftwarePixelFormat(swFormat PixelFormat) {
	hfc.data().sw_format = C.enum_AVPixelFormat(swFormat)
}

// https://ffmpeg.org/doxygen/7.0/structAVHWFramesContext.html#a9c3a94dcd9c96e19059b56a6bae9c764
func (hfc *HardwareFrameContext) SetInitialPoolSize(initialPoolSize int) {
	hfc.data().initial_pool_size = C.int(initialPoolSize)
}

// https://ffmpeg.org/doxygen/7.0/hwcontext_8c.html#a66a7e1ebc7e459ce07d3de6639ac7e38
func (hfc *HardwareFrameContext) Initialize() error {
	return newError(C.av_hwframe_ctx_init(hfc.c))
}
