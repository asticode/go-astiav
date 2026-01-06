package astiav

/*
#include <libavcodec/avcodec.h>
*/
import "C"

// https://ffmpeg.org/doxygen/8.0/group__lavc__decoding.html#ga9abe043ed40f3c270dff77235fcfcd0b
type CodecConfig C.enum_AVCodecConfig

const (
	CodecConfigPixFormat     = CodecConfig(C.AV_CODEC_CONFIG_PIX_FORMAT)
	CodecConfigFrameRate     = CodecConfig(C.AV_CODEC_CONFIG_FRAME_RATE)
	CodecConfigSampleRate    = CodecConfig(C.AV_CODEC_CONFIG_SAMPLE_RATE)
	CodecConfigSampleFormat  = CodecConfig(C.AV_CODEC_CONFIG_SAMPLE_FORMAT)
	CodecConfigChannelLayout = CodecConfig(C.AV_CODEC_CONFIG_CHANNEL_LAYOUT)
	CodecConfigColorRange    = CodecConfig(C.AV_CODEC_CONFIG_COLOR_RANGE)
	CodecConfigColorSpace    = CodecConfig(C.AV_CODEC_CONFIG_COLOR_SPACE)
)
