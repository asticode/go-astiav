package astiav

//#include <libavutil/mathematics.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/group__lavu__math.html#gaf02994a8bbeaa91d4757df179cbe567f
func RescaleQ(a int64, b Rational, c Rational) int64 {
	return int64(C.av_rescale_q(C.int64_t(a), b.c, c.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__math__rational.html#ga3f9c69432582e2857147bcba3c75dc32
func MulQ(a Rational, b Rational) Rational {
	return newRationalFromC(C.av_mul_q(a.c, b.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__math.html#gab706bfec9bf56534e02ca9564cb968f6
func RescaleQRnd(a int64, b Rational, c Rational, r Rounding) int64 {
	return int64(C.av_rescale_q_rnd(C.int64_t(a), b.c, c.c, C.enum_AVRounding(r)))
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__math.html#ga29b7c3d60d68ef678ee1f4adc61a25dc
func RescaleDelta(inTB Rational, inTS int64, fsTB Rational, duration, last int64, outTB Rational) (out, lastTS int64) {
	clast := C.int64_t(last)
	outTS := C.av_rescale_delta(inTB.c, C.int64_t(inTS), fsTB.c, C.int(duration), &clast, outTB.c)
	return int64(outTS), int64(clast)
}
