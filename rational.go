package astiav

//#include <libavutil/rational.h>
import "C"
import (
	"encoding"
	"fmt"
	"strconv"
	"strings"
)

// https://ffmpeg.org/doxygen/7.0/structAVRational.html
type Rational struct {
	c C.AVRational
}

var (
	_ encoding.TextMarshaler   = (*Rational)(nil)
	_ encoding.TextUnmarshaler = (*Rational)(nil)
)

func newRationalFromC(c C.AVRational) Rational {
	return Rational{c: c}
}

func NewRational(num, den int) Rational {
	var r Rational
	r.SetNum(num)
	r.SetDen(den)
	return r
}

// NewRationalFromString returns a Rational by parsing the specified string in
// the form of 'num/den' matching to the output of String(). If the string
// contains a single numeric value the denominator is set to 1.
func NewRationalFromString(s string) (Rational, error) {
	if s == "" || s == "0" {
		return Rational{}, nil
	}
	parts := strings.Split(s, "/")
	num, err := strconv.ParseInt(parts[0], 10, 0)
	if err != nil {
		return Rational{}, fmt.Errorf("astiav: parsing numerator failed: %w", err)
	}
	var den int64 = 1
	if len(parts) > 1 {
		if den, err = strconv.ParseInt(parts[1], 10, 0); err != nil {
			return Rational{}, fmt.Errorf("astiav: parsing denominator failed: %w", err)
		}
	}
	return NewRational(int(num), int(den)), nil
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
	if r.Den() == 1 {
		return strconv.Itoa(r.Num())
	}
	return strconv.Itoa(r.Num()) + "/" + strconv.Itoa(r.Den())
}

func (r Rational) Invert() Rational {
	return NewRational(r.Den(), r.Num())
}

func (r Rational) MarshalText() ([]byte, error) { return ([]byte)(r.String()), nil }
func (r *Rational) UnmarshalText(d []byte) error {
	rr, err := NewRationalFromString(string(d))
	if err != nil {
		return err
	}
	*r = rr
	return nil
}
