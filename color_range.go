package astiav

//#include <libavutil/pixfmt.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/pixfmt.h#L562
type ColorRange C.enum_AVColorRange

const (
	ColorRangeUnspecified = ColorRange(C.AVCOL_RANGE_UNSPECIFIED)
	ColorRangeMpeg        = ColorRange(C.AVCOL_RANGE_MPEG)
	ColorRangeJpeg        = ColorRange(C.AVCOL_RANGE_JPEG)
	ColorRangeNb          = ColorRange(C.AVCOL_RANGE_NB)
)
