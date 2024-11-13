package astiav

//#include <libavutil/avutil.h>
//#include <libavutil/time.h>
import "C"

const (
	// https://ffmpeg.org/doxygen/7.0/group__lavu__time.html#ga2eaefe702f95f619ea6f2d08afa01be1
	NoPtsValue = int64(C.AV_NOPTS_VALUE)
	// https://ffmpeg.org/doxygen/7.0/group__lavu__time.html#gaa11ed202b70e1f52bac809811a910e2a
	TimeBase = int(C.AV_TIME_BASE)
)

var (
	// https://ffmpeg.org/doxygen/7.0/group__lavu__time.html#gafd07a13a4ddaa6015275cad6022d9ee3
	TimeBaseQ = newRationalFromC(C.AV_TIME_BASE_Q)
)

// https://ffmpeg.org/doxygen/7.0/time_8c.html#adf0e36df54426fa167e3cc5a3406f3b7
func RelativeTime() int64 {
	return int64(C.av_gettime_relative())
}
