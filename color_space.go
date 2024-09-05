package astiav

//#include <libavutil/pixfmt.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/pixfmt.h#L523
type ColorSpace C.enum_AVColorSpace

const (
	ColorSpaceRgb              = ColorSpace(C.AVCOL_SPC_RGB)
	ColorSpaceBt709            = ColorSpace(C.AVCOL_SPC_BT709)
	ColorSpaceUnspecified      = ColorSpace(C.AVCOL_SPC_UNSPECIFIED)
	ColorSpaceReserved         = ColorSpace(C.AVCOL_SPC_RESERVED)
	ColorSpaceFcc              = ColorSpace(C.AVCOL_SPC_FCC)
	ColorSpaceBt470Bg          = ColorSpace(C.AVCOL_SPC_BT470BG)
	ColorSpaceSmpte170M        = ColorSpace(C.AVCOL_SPC_SMPTE170M)
	ColorSpaceSmpte240M        = ColorSpace(C.AVCOL_SPC_SMPTE240M)
	ColorSpaceYcgco            = ColorSpace(C.AVCOL_SPC_YCGCO)
	ColorSpaceYcocg            = ColorSpace(C.AVCOL_SPC_YCOCG)
	ColorSpaceBt2020Ncl        = ColorSpace(C.AVCOL_SPC_BT2020_NCL)
	ColorSpaceBt2020Cl         = ColorSpace(C.AVCOL_SPC_BT2020_CL)
	ColorSpaceSmpte2085        = ColorSpace(C.AVCOL_SPC_SMPTE2085)
	ColorSpaceChromaDerivedNcl = ColorSpace(C.AVCOL_SPC_CHROMA_DERIVED_NCL)
	ColorSpaceChromaDerivedCl  = ColorSpace(C.AVCOL_SPC_CHROMA_DERIVED_CL)
	ColorSpaceIctcp            = ColorSpace(C.AVCOL_SPC_ICTCP)
	ColorSpaceNb               = ColorSpace(C.AVCOL_SPC_NB)
)
