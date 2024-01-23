package astiav

//#cgo pkg-config: libavutil
//#include <libavutil/hwcontext.h>
import "C"
import (
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/hwcontext.h#L27
type HardwareDeviceType C.enum_AVHWDeviceType

const (
	HardwareDeviceTypeNone         = HardwareDeviceType(C.AV_HWDEVICE_TYPE_NONE)
	HardwareDeviceTypeVDPAU        = HardwareDeviceType(C.AV_HWDEVICE_TYPE_VDPAU)
	HardwareDeviceTypeCUDA         = HardwareDeviceType(C.AV_HWDEVICE_TYPE_CUDA)
	HardwareDeviceTypeVAAPI        = HardwareDeviceType(C.AV_HWDEVICE_TYPE_VAAPI)
	HardwareDeviceTypeDXVA2        = HardwareDeviceType(C.AV_HWDEVICE_TYPE_DXVA2)
	HardwareDeviceTypeQSV          = HardwareDeviceType(C.AV_HWDEVICE_TYPE_QSV)
	HardwareDeviceTypeVideoToolbox = HardwareDeviceType(C.AV_HWDEVICE_TYPE_VIDEOTOOLBOX)
	HardwareDeviceTypeD3D11VA      = HardwareDeviceType(C.AV_HWDEVICE_TYPE_D3D11VA)
	HardwareDeviceTypeDRM          = HardwareDeviceType(C.AV_HWDEVICE_TYPE_DRM)
	HardwareDeviceTypeOpenCL       = HardwareDeviceType(C.AV_HWDEVICE_TYPE_OPENCL)
	HardwareDeviceTypeMediaCodec   = HardwareDeviceType(C.AV_HWDEVICE_TYPE_MEDIACODEC)
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
