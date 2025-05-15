package astiav

//#include <libavcodec/avcodec.h>
//#include <libavutil/hwcontext.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVHWDeviceContext.html
type HardwareDeviceContext struct {
	c *C.AVBufferRef
}

// https://ffmpeg.org/doxygen/7.0/hwcontext_8c.html#a21fbd088225e4e25c4d9a01b3f5e8c51
func CreateHardwareDeviceContext(t HardwareDeviceType, device string, options *Dictionary, flags int) (*HardwareDeviceContext, error) {
	hdc := HardwareDeviceContext{}
	deviceC := (*C.char)(nil)
	if device != "" {
		deviceC = C.CString(device)
		defer C.free(unsafe.Pointer(deviceC))
	}
	optionsC := (*C.AVDictionary)(nil)
	if options != nil {
		optionsC = options.c
	}
	if err := newError(C.av_hwdevice_ctx_create(&hdc.c, (C.enum_AVHWDeviceType)(t), deviceC, optionsC, C.int(flags))); err != nil {
		return nil, err
	}
	return &hdc, nil
}

// https://ffmpeg.org/doxygen/7.0/hwcontext_8c.html#a80f4c1184e1758150b6d9bc0adf2c1df
func (hdc *HardwareDeviceContext) HardwareFramesConstraints() *HardwareFramesConstraints {
	return newHardwareFramesConstraintsFromC(C.av_hwdevice_get_hwframe_constraints(hdc.c, nil))
}

func (hdc *HardwareDeviceContext) Free() {
	if hdc.c != nil {
		C.av_buffer_unref(&hdc.c)
	}
}
