package astiav

//#include <libavcodec/avcodec.h>
//#include <libavutil/avutil.h>
import "C"

// https://ffmpeg.org/doxygen/7.1/group__lavu__misc.html#ga9a84bba4713dfced21a1a56163be1f48
type MediaType C.enum_AVMediaType

const (
	MediaTypeAttachment = MediaType(C.AVMEDIA_TYPE_ATTACHMENT)
	MediaTypeAudio      = MediaType(C.AVMEDIA_TYPE_AUDIO)
	MediaTypeData       = MediaType(C.AVMEDIA_TYPE_DATA)
	MediaTypeNb         = MediaType(C.AVMEDIA_TYPE_NB)
	MediaTypeSubtitle   = MediaType(C.AVMEDIA_TYPE_SUBTITLE)
	MediaTypeUnknown    = MediaType(C.AVMEDIA_TYPE_UNKNOWN)
	MediaTypeVideo      = MediaType(C.AVMEDIA_TYPE_VIDEO)
)

// https://ffmpeg.org/doxygen/7.1/group__lavu__misc.html#gaf21645cfa855b2caf9699d7dc7b2d08e
func (t MediaType) String() string {
	return C.GoString(C.av_get_media_type_string((C.enum_AVMediaType)(t)))
}
