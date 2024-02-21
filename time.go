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

func RelativeTime()  int64{
	return int64(C.av_gettime_relative())
}
