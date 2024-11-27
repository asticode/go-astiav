package astiav

//#include <libavutil/mathematics.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/group__lavu__math.html#gaf02994a8bbeaa91d4757df179cbe567f
func RescaleQ(a int64, b Rational, c Rational) int64 {
	return int64(C.av_rescale_q(C.int64_t(a), b.c, c.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__math.html#gab706bfec9bf56534e02ca9564cb968f6
func RescaleQRnd(a int64, b Rational, c Rational, r Rounding) int64 {
	return int64(C.av_rescale_q_rnd(C.int64_t(a), b.c, c.c, C.enum_AVRounding(r)))
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__math__rational.html#ga935dbbf6bde8dfe5fa7ddb1da582eb07
func Q2D(a Rational) float64 {
	return float64(C.av_q2d(a.c))
}
