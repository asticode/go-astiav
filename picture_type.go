package astiav

//#include <libavutil/avutil.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/group__lavu__picture.html#gae6cbcab1f70d8e476757f1c1f5a0a78e
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

// https://ffmpeg.org/doxygen/8.1/group__lavu__picture.html#gacbf2ea8b2b89924c890ef8ec10a3d922
func (t PictureType) String() string {
	return string(rune(C.av_get_picture_type_char((C.enum_AVPictureType)(t))))
}
