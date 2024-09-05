package astiav

//#include <libavutil/hwcontext.h>
import "C"
import (
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/hwcontext.h#L27
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

func (t HardwareDeviceType) String() string {
	return C.GoString(C.av_hwdevice_get_type_name((C.enum_AVHWDeviceType)(t)))
}

func FindHardwareDeviceTypeByName(n string) HardwareDeviceType {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return HardwareDeviceType(C.av_hwdevice_find_type_by_name(cn))
}
