package astiav

//#include <libavutil/hwcontext.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.1/hwcontext_8h.html#acf25724be4b066a51ad86aa9214b0d34
type HardwareDeviceType C.enum_AVHWDeviceType

const (
	HardwareDeviceTypeCUDA         = HardwareDeviceType(C.AV_HWDEVICE_TYPE_CUDA)
	HardwareDeviceTypeD3D11VA      = HardwareDeviceType(C.AV_HWDEVICE_TYPE_D3D11VA)
	HardwareDeviceTypeDRM          = HardwareDeviceType(C.AV_HWDEVICE_TYPE_DRM)
	HardwareDeviceTypeDXVA2        = HardwareDeviceType(C.AV_HWDEVICE_TYPE_DXVA2)
	HardwareDeviceTypeMediaCodec   = HardwareDeviceType(C.AV_HWDEVICE_TYPE_MEDIACODEC)
	HardwareDeviceTypeNone         = HardwareDeviceType(C.AV_HWDEVICE_TYPE_NONE)
	HardwareDeviceTypeOpenCL       = HardwareDeviceType(C.AV_HWDEVICE_TYPE_OPENCL)
	HardwareDeviceTypeQSV          = HardwareDeviceType(C.AV_HWDEVICE_TYPE_QSV)
	HardwareDeviceTypeVAAPI        = HardwareDeviceType(C.AV_HWDEVICE_TYPE_VAAPI)
	HardwareDeviceTypeVDPAU        = HardwareDeviceType(C.AV_HWDEVICE_TYPE_VDPAU)
	HardwareDeviceTypeVideoToolbox = HardwareDeviceType(C.AV_HWDEVICE_TYPE_VIDEOTOOLBOX)
	HardwareDeviceTypeVulkan       = HardwareDeviceType(C.AV_HWDEVICE_TYPE_VULKAN)
)

// https://ffmpeg.org/doxygen/7.1/hwcontext_8h.html#afb2b99a15f3fdde25a2fd19353ac5a67
func (t HardwareDeviceType) Name() string {
	return C.GoString(C.av_hwdevice_get_type_name((C.enum_AVHWDeviceType)(t)))
}

func (t HardwareDeviceType) String() string {
	return t.Name()
}

// https://ffmpeg.org/doxygen/7.1/hwcontext_8h.html#a541943ddced791765349645a30adfa4d
func FindHardwareDeviceTypeByName(n string) HardwareDeviceType {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return HardwareDeviceType(C.av_hwdevice_find_type_by_name(cn))
}
