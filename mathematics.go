package astiav

//#include <libavutil/mathematics.h>
import "C"
import "unsafe"

// https://ffmpeg.org/doxygen/8.0/group__lavu__math.html#gaf02994a8bbeaa91d4757df179cbe567f
func RescaleQ(a int64, b Rational, c Rational) int64 {
	return int64(C.av_rescale_q(C.int64_t(a), b.c, c.c))
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__math.html#gab706bfec9bf56534e02ca9564cb968f6
func RescaleQRnd(a int64, b Rational, c Rational, r Rounding) int64 {
	return int64(C.av_rescale_q_rnd(C.int64_t(a), b.c, c.c, C.enum_AVRounding(r)))
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__math.html#ga1e2c8b2b71c2cd5db1b5b3b7b7b7b7b7
func Rescale(a, b, c int64) int64 {
	return int64(C.av_rescale(C.int64_t(a), C.int64_t(b), C.int64_t(c)))
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__math.html#ga1e2c8b2b71c2cd5db1b5b3b7b7b7b7b7
func RescaleRnd(a, b, c int64, rnd Rounding) int64 {
	return int64(C.av_rescale_rnd(C.int64_t(a), C.int64_t(b), C.int64_t(c), C.enum_AVRounding(rnd)))
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__math.html#ga1e2c8b2b71c2cd5db1b5b3b7b7b7b7b7
func Gcd(a, b int64) int64 {
	return int64(C.av_gcd(C.int64_t(a), C.int64_t(b)))
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__math.html#ga1e2c8b2b71c2cd5db1b5b3b7b7b7b7b7
func CompareTs(tsA int64, tbA Rational, tsB int64, tbB Rational) int {
	return int(C.av_compare_ts(C.int64_t(tsA), tbA.c, C.int64_t(tsB), tbB.c))
}

// CompareTimestamps is implemented in time.go

// RescaleDelta rescales a timestamp while preserving known durations
// https://ffmpeg.org/doxygen/8.0/group__lavu__math.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func RescaleDelta(inTb Rational, inTs int64, fsTb Rational, duration int, last *int64, outTb Rational) int64 {
	var lastPtr *C.int64_t
	if last != nil {
		lastPtr = (*C.int64_t)(unsafe.Pointer(last))
	}
	return int64(C.av_rescale_delta(inTb.c, C.int64_t(inTs), fsTb.c, C.int(duration), lastPtr, outTb.c))
}

// AddStable adds two rationals (implemented in rational.go)

// AddStable adds two rationals (implemented in rational.go)
// SubStable subtracts one rational from another (implemented in rational.go)

// Reduce reduces a fraction to its simplest form
// https://ffmpeg.org/doxygen/8.0/group__lavu__math.html#ga1e2c8b2b71c2cd5db1b5b3b7b7b7b7b7
func Reduce(num, den int64, max int64) (int, int) {
	var cNum, cDen C.int
	C.av_reduce(&cNum, &cDen, C.int64_t(num), C.int64_t(den), C.int64_t(max))
	return int(cNum), int(cDen)
}

// Q2d converts rational to double
func Q2d(r Rational) float64 {
	if r.Den() == 0 {
		return 0.0
	}
	return float64(r.Num()) / float64(r.Den())
}
