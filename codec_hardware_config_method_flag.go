package astiav

//#cgo pkg-config: libavcodec
//#include <libavcodec/avcodec.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/codec.h#L420
type CodecHardwareConfigMethodFlag int

const (
	CodecHardwareConfigMethodFlagAdHoc       = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_AD_HOC)
	CodecHardwareConfigMethodFlagHwDeviceCtx = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_HW_DEVICE_CTX)
	CodecHardwareConfigMethodFlagHwFramesCtx = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_HW_FRAMES_CTX)
	CodecHardwareConfigMethodFlagInternal    = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_INTERNAL)
)
