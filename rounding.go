package astiav

//#include <libavutil/mathematics.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/mathematics.h#L79
type Rounding C.enum_AVRounding

const (
	RoundingZero       = Rounding(C.AV_ROUND_ZERO)
	RoundingInf        = Rounding(C.AV_ROUND_INF)
	RoundingDown       = Rounding(C.AV_ROUND_DOWN)
	RoundingUp         = Rounding(C.AV_ROUND_UP)
	RoundingNearInf    = Rounding(C.AV_ROUND_NEAR_INF)
	RoundingPassMinmax = Rounding(C.AV_ROUND_PASS_MINMAX)
)
