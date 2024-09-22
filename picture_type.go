package astiav

//#include <libavutil/avutil.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/avutil.h#L272
type PictureType C.enum_AVPictureType

const (
	PictureTypeNone = PictureType(C.AV_PICTURE_TYPE_NONE)
	PictureTypeI    = PictureType(C.AV_PICTURE_TYPE_I)
	PictureTypeP    = PictureType(C.AV_PICTURE_TYPE_P)
	PictureTypeB    = PictureType(C.AV_PICTURE_TYPE_B)
	PictureTypeS    = PictureType(C.AV_PICTURE_TYPE_S)
	PictureTypeSi   = PictureType(C.AV_PICTURE_TYPE_SI)
	PictureTypeSp   = PictureType(C.AV_PICTURE_TYPE_SP)
	PictureTypeBi   = PictureType(C.AV_PICTURE_TYPE_BI)
)

func (t PictureType) String() string {
	return string(rune(C.av_get_picture_type_char((C.enum_AVPictureType)(t))))
}
