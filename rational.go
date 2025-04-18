package astiav

//#include <libavutil/rational.h>
import "C"
import (
	"encoding"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// https://ffmpeg.org/doxygen/7.0/structAVRational.html
type Rational struct {
	c C.AVRational
}

var _ json.Marshaler = (*Rational)(nil)
var _ json.Unmarshaler = (*Rational)(nil)
var _ encoding.TextMarshaler = (*Rational)(nil)
var _ encoding.TextUnmarshaler = (*Rational)(nil)

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

type rational struct {
	Num int
	Den int
}

func (r *Rational) UnmarshalJSON(d []byte) error {
	var rs rational
	if err := json.Unmarshal(d, &rs); err != nil {
		return err
	}
	r.SetNum(rs.Num)
	r.SetDen(rs.Den)
	return nil
}

func (r Rational) MarshalJSON() ([]byte, error) {
	return json.Marshal(rational{Num: r.Num(), Den: r.Den()})
}

func (r Rational) MarshalText() ([]byte, error) { return ([]byte)(r.String()), nil }
func (r *Rational) UnmarshalText(d []byte) error {
	s := string(d)
	if s == "" {
		r.SetNum(0)
		r.SetDen(0)
	}
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return fmt.Errorf("value %s is not Rational ('d/d')", s)
	}
	num, err := strconv.ParseInt(parts[0], 10, 0)
	if err != nil {
		return fmt.Errorf("value %s is not Rational ('d/d'): %w", s, err)
	}
	den, err := strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		return fmt.Errorf("value %s is not Rational ('d/d'): %w", s, err)
	}
	r.SetNum(int(num))
	r.SetDen(int(den))
	return nil
}
