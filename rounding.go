package astiav

//#include <libavutil/mathematics.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/group__lavu__math.html#ga921d656eaf2c4d6800a734a13af021d0
type Rounding C.enum_AVRounding

const (
	RoundingZero       = Rounding(C.AV_ROUND_ZERO)
	RoundingInf        = Rounding(C.AV_ROUND_INF)
	RoundingDown       = Rounding(C.AV_ROUND_DOWN)
	RoundingUp         = Rounding(C.AV_ROUND_UP)
	RoundingNearInf    = Rounding(C.AV_ROUND_NEAR_INF)
	RoundingPassMinmax = Rounding(C.AV_ROUND_PASS_MINMAX)
)
