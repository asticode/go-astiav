package astiav

//#include <libavutil/pixdesc.h>
//#include <libavutil/pixfmt.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/pixfmt_8h.html#a3da0bf691418bc22c4bcbe6583ad589a
type ColorRange C.enum_AVColorRange

const (
	ColorRangeUnspecified = ColorRange(C.AVCOL_RANGE_UNSPECIFIED)
	ColorRangeMpeg        = ColorRange(C.AVCOL_RANGE_MPEG)
	ColorRangeJpeg        = ColorRange(C.AVCOL_RANGE_JPEG)
	ColorRangeNb          = ColorRange(C.AVCOL_RANGE_NB)
)

// https://ffmpeg.org/doxygen/8.1/pixdesc_8c.html#a590decf389632dd3af095f3096a92caf
func (r ColorRange) Name() string {
	return C.GoString(C.av_color_range_name(C.enum_AVColorRange(r)))
}

func (r ColorRange) String() string {
	return r.Name()
}
