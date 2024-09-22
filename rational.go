package astiav

//#include <libavutil/rational.h>
import "C"
import "strconv"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/rational.h#L58
type Rational struct {
	c C.AVRational
}

func newRationalFromC(c C.AVRational) Rational {
	return Rational{c: c}
}

func NewRational(num, den int) Rational {
	var r Rational
	r.SetNum(num)
	r.SetDen(den)
	return r
}

func (r Rational) Num() int {
	return int(r.c.num)
}

func (r *Rational) SetNum(num int) {
	r.c.num = C.int(num)
}

func (r Rational) Den() int {
	return int(r.c.den)
}

func (r *Rational) SetDen(den int) {
	r.c.den = C.int(den)
}

func (r Rational) Float64() float64 {
	if r.Num() == 0 || r.Den() == 0 {
		return 0
	}
	return float64(r.Num()) / float64(r.Den())
}

func (r Rational) String() string {
	if r.Num() == 0 || r.Den() == 0 {
		return "0"
	}
	return strconv.Itoa(r.Num()) + "/" + strconv.Itoa(r.Den())
}

func (r Rational) Invert() Rational {
	return NewRational(r.Den(), r.Num())
}
