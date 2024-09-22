package astiav

//#include <libavcodec/avcodec.h>
//#include <libavutil/avutil.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/avutil.h#L199
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

func (t MediaType) String() string {
	return C.GoString(C.av_get_media_type_string((C.enum_AVMediaType)(t)))
}
