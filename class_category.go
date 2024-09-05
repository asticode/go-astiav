package astiav

//#include <libavutil/log.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/log.h#L28
type ClassCategory C.AVClassCategory

const (
	ClassCategoryBitstreamFilter   = ClassCategory(C.AV_CLASS_CATEGORY_BITSTREAM_FILTER)
	ClassCategoryDecoder           = ClassCategory(C.AV_CLASS_CATEGORY_DECODER)
	ClassCategoryDemuxer           = ClassCategory(C.AV_CLASS_CATEGORY_DEMUXER)
	ClassCategoryDeviceAudioInput  = ClassCategory(C.AV_CLASS_CATEGORY_DEVICE_AUDIO_INPUT)
	ClassCategoryDeviceAudioOutput = ClassCategory(C.AV_CLASS_CATEGORY_DEVICE_AUDIO_OUTPUT)
	ClassCategoryDeviceInput       = ClassCategory(C.AV_CLASS_CATEGORY_DEVICE_INPUT)
	ClassCategoryDeviceOutput      = ClassCategory(C.AV_CLASS_CATEGORY_DEVICE_OUTPUT)
	ClassCategoryDeviceVideoInput  = ClassCategory(C.AV_CLASS_CATEGORY_DEVICE_VIDEO_INPUT)
	ClassCategoryDeviceVideoOutput = ClassCategory(C.AV_CLASS_CATEGORY_DEVICE_VIDEO_OUTPUT)
	ClassCategoryEncoder           = ClassCategory(C.AV_CLASS_CATEGORY_ENCODER)
	ClassCategoryFilter            = ClassCategory(C.AV_CLASS_CATEGORY_FILTER)
	ClassCategoryInput             = ClassCategory(C.AV_CLASS_CATEGORY_INPUT)
	ClassCategoryMuxer             = ClassCategory(C.AV_CLASS_CATEGORY_MUXER)
	ClassCategoryNa                = ClassCategory(C.AV_CLASS_CATEGORY_NA)
	ClassCategoryNb                = ClassCategory(C.AV_CLASS_CATEGORY_NB)
	ClassCategoryOutput            = ClassCategory(C.AV_CLASS_CATEGORY_OUTPUT)
	ClassCategorySwresampler       = ClassCategory(C.AV_CLASS_CATEGORY_SWRESAMPLER)
	ClassCategorySwscaler          = ClassCategory(C.AV_CLASS_CATEGORY_SWSCALER)
)
