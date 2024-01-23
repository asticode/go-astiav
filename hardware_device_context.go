package astiav

//#cgo pkg-config: libavutil libavcodec
//#include <libavcodec/avcodec.h>
//#include <libavutil/hwcontext.h>
import "C"
import (
	"fmt"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/hwcontext.h#L141
type HardwareDeviceContext struct {
	c *C.AVBufferRef
}

func CreateHardwareDeviceContext(t HardwareDeviceType, device string, options *Dictionary) (*HardwareDeviceContext, error) {
	hdc := HardwareDeviceContext{}

	// Check for an emtpy string and pass NULL to av_hwdevice_ctx_create if its emtpy
	var deviceC *C.char
	if device != "" {
		deviceC = C.CString(device)
		defer C.free(unsafe.Pointer(deviceC))
	} else {
		deviceC = (*C.char)(nil)
	}
	errorCode := C.av_hwdevice_ctx_create(&hdc.c, (C.enum_AVHWDeviceType)(t), deviceC, options.c, 0)
	if errorCode < 0 {
		return nil, newError(errorCode)
	}
	return &hdc, nil
}

func FindSuitableHardwareFormat(decoder *Codec, deviceType HardwareDeviceType) (PixelFormat, error) {
	var i int
	for {
		config := C.avcodec_get_hw_config(decoder.c, C.int(i))
		if config == nil {
			return 0, fmt.Errorf("Decoder %s does not support device type %s", decoder.Name(), deviceType.String())
		}
		if config.methods&C.AV_CODEC_HW_CONFIG_METHOD_HW_DEVICE_CTX != 0 && config.device_type == C.enum_AVHWDeviceType(deviceType) {
			return PixelFormat(config.pix_fmt), nil
		}
		i++
	}
}

func (hdc *HardwareDeviceContext) Free() {
	if hdc.c != nil {
		C.av_buffer_unref(&hdc.c)
	}
}
