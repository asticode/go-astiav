package astiav

// #cgo pkg-config: libavutil
// #include <libavutil/hwcontext.h>
import "C"

func NewDeviceContext(deviceType DeviceTypeFlag) *BufferRef {
	br := newBufferRef()
	ret := int(C.av_hwdevice_ctx_create(&br.c, uint32(deviceType), nil, nil, 0))
	if ret < 0 {
		return nil
	} else {
		return br
	}
}
