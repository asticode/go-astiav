package astiav

//#include <libavcodec/avcodec.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/group__lavc__core.html#gga9334a5b9057f32da96db9b5c6a045d67a680870b80f0ed65e9ba97ea0905eb2fa
type CodecHardwareConfigMethodFlag int64

const (
	CodecHardwareConfigMethodFlagAdHoc       = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_AD_HOC)
	CodecHardwareConfigMethodFlagHwDeviceCtx = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_HW_DEVICE_CTX)
	CodecHardwareConfigMethodFlagHwFramesCtx = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_HW_FRAMES_CTX)
	CodecHardwareConfigMethodFlagInternal    = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_INTERNAL)
)
