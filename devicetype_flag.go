package astiav

// #cgo pkg-config: libavutil
// #include <libavutil/buffer.h>
// #include <libavutil/hwcontext.h>
import "C"

type DeviceTypeFlag uint32

const (
	DeviceTypeVAAPI = DeviceTypeFlag(C.AV_HWDEVICE_TYPE_VAAPI)
)
