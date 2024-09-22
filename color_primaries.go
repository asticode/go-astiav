package astiav

//#include <libavutil/pixfmt.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/pixfmt.h#L469
type ColorPrimaries C.enum_AVColorPrimaries

const (
	ColorPrimariesReserved0   = ColorPrimaries(C.AVCOL_PRI_RESERVED0)
	ColorPrimariesBt709       = ColorPrimaries(C.AVCOL_PRI_BT709)
	ColorPrimariesUnspecified = ColorPrimaries(C.AVCOL_PRI_UNSPECIFIED)
	ColorPrimariesReserved    = ColorPrimaries(C.AVCOL_PRI_RESERVED)
	ColorPrimariesBt470M      = ColorPrimaries(C.AVCOL_PRI_BT470M)
	ColorPrimariesBt470Bg     = ColorPrimaries(C.AVCOL_PRI_BT470BG)
	ColorPrimariesSmpte170M   = ColorPrimaries(C.AVCOL_PRI_SMPTE170M)
	ColorPrimariesSmpte240M   = ColorPrimaries(C.AVCOL_PRI_SMPTE240M)
	ColorPrimariesFilm        = ColorPrimaries(C.AVCOL_PRI_FILM)
	ColorPrimariesBt2020      = ColorPrimaries(C.AVCOL_PRI_BT2020)
	ColorPrimariesSmpte428    = ColorPrimaries(C.AVCOL_PRI_SMPTE428)
	ColorPrimariesSmptest4281 = ColorPrimaries(C.AVCOL_PRI_SMPTEST428_1)
	ColorPrimariesSmpte431    = ColorPrimaries(C.AVCOL_PRI_SMPTE431)
	ColorPrimariesSmpte432    = ColorPrimaries(C.AVCOL_PRI_SMPTE432)
	ColorPrimariesJedecP22    = ColorPrimaries(C.AVCOL_PRI_JEDEC_P22)
	ColorPrimariesNb          = ColorPrimaries(C.AVCOL_PRI_NB)
)
