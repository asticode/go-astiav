package astiav

//#cgo pkg-config: libavutil
//#include <libavutil/avutil.h>
import "C"

const (
	NoPtsValue = int64(C.AV_NOPTS_VALUE)
	TimeBase   = int(C.AV_TIME_BASE)
)

var (
	TimeBaseQ = newRationalFromC(C.AV_TIME_BASE_Q)
)

func Q2D(a Rational)  float64{
	return float64(C.av_q2d(a.c))
}

func GetTimeRelative()  int64{
	return int64(C.av_gettime_relative())
}
