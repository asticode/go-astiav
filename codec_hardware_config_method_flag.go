package astiav

//#cgo pkg-config: libavcodec
//#include <libavcodec/avcodec.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/codec.h#L420
type CodecHardwareConfigMethodFlag int

const (
	CodecHardwareConfigMethodAdHoc       = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_AD_HOC)
	CodecHardwareConfigMethodHwDeviceCtx = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_HW_DEVICE_CTX)
	CodecHardwareConfigMethodHwFramesCtx = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_HW_FRAMES_CTX)
	CodecHardwareConfigMethodInternal    = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_INTERNAL)
)
