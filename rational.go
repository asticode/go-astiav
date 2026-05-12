package astiav

//#include <libavutil/rational.h>
import "C"
import "strconv"

// https://ffmpeg.org/doxygen/8.0/structAVRational.html
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

// https://ffmpeg.org/doxygen/8.0/group__lavu__math__rational.html#ga2eb3a275aabacd8421f140a12bab4a91
func (r Rational) Add(v Rational) Rational {
	return newRationalFromC(C.av_add_q(r.c, v.c))
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__math__rational.html#gac66c6198ce5e8a8caf88dfc20782fa59
func (r Rational) Sub(v Rational) Rational {
	return newRationalFromC(C.av_sub_q(r.c, v.c))
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__math__rational.html#ga3f9c69432582e2857147bcba3c75dc32
func (r Rational) Mul(v Rational) Rational {
	return newRationalFromC(C.av_mul_q(r.c, v.c))
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__math__rational.html#gaffa24e7bd38e12dbac540d8b66461f97
func (r Rational) Div(v Rational) Rational {
	return newRationalFromC(C.av_div_q(r.c, v.c))
}
