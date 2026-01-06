package astiav

//#include <libavutil/pixdesc.h>
//#include <libavutil/pixfmt.h>
import "C"

// https://ffmpeg.org/doxygen/7.1/pixfmt_8h.html#aff71a069509a1ad3ff54d53a1c894c85
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
	ColorSpaceYcgcoRe          = ColorSpace(C.AVCOL_SPC_YCGCO_RE)
	ColorSpaceYcgcoRo          = ColorSpace(C.AVCOL_SPC_YCGCO_RO)
	ColorSpaceYcocg            = ColorSpace(C.AVCOL_SPC_YCOCG)
	ColorSpaceBt2020Ncl        = ColorSpace(C.AVCOL_SPC_BT2020_NCL)
	ColorSpaceBt2020Cl         = ColorSpace(C.AVCOL_SPC_BT2020_CL)
	ColorSpaceSmpte2085        = ColorSpace(C.AVCOL_SPC_SMPTE2085)
	ColorSpaceChromaDerivedNcl = ColorSpace(C.AVCOL_SPC_CHROMA_DERIVED_NCL)
	ColorSpaceChromaDerivedCl  = ColorSpace(C.AVCOL_SPC_CHROMA_DERIVED_CL)
	ColorSpaceIctcp            = ColorSpace(C.AVCOL_SPC_ICTCP)
	ColorSpaceIptC2            = ColorSpace(C.AVCOL_SPC_IPT_C2)
	ColorSpaceNb               = ColorSpace(C.AVCOL_SPC_NB)
)

// https://ffmpeg.org/doxygen/7.1/pixdesc_8c.html#a7a5b3f4d128f0a0112b4a91f75055339
func (s ColorSpace) Name() string {
	return C.GoString(C.av_color_space_name(C.enum_AVColorSpace(s)))
}

func (s ColorSpace) String() string {
	return s.Name()
}
