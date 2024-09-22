package astiav

//#include <libavcodec/avcodec.h>
//#include <libavutil/hwcontext.h>
import "C"
import (
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/hwcontext.h#L61
type HardwareDeviceContext struct {
	c *C.AVBufferRef
}

func CreateHardwareDeviceContext(t HardwareDeviceType, device string, options *Dictionary) (*HardwareDeviceContext, error) {
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
	if err := newError(C.av_hwdevice_ctx_create(&hdc.c, (C.enum_AVHWDeviceType)(t), deviceC, optionsC, 0)); err != nil {
		return nil, err
	}
	return &hdc, nil
}
