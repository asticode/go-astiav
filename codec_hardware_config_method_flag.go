package astiav

//#include <libavcodec/avcodec.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/group__lavc__core.html#gga26e1124d33b4acdb532c49f6498df549a680870b80f0ed65e9ba97ea0905eb2fa
type CodecHardwareConfigMethodFlag int64

const (
	CodecHardwareConfigMethodFlagAdHoc       = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_AD_HOC)
	CodecHardwareConfigMethodFlagHwDeviceCtx = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_HW_DEVICE_CTX)
	CodecHardwareConfigMethodFlagHwFramesCtx = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_HW_FRAMES_CTX)
	CodecHardwareConfigMethodFlagInternal    = CodecHardwareConfigMethodFlag(C.AV_CODEC_HW_CONFIG_METHOD_INTERNAL)
)
