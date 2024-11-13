package astiav

//#include <libavcodec/avcodec.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/avcodec_8h.html#a116c7fb56ac57ccca3e08b80467b2a40
type ThreadType int

const (
	ThreadTypeFrame     = ThreadType(C.FF_THREAD_FRAME)
	ThreadTypeSlice     = ThreadType(C.FF_THREAD_SLICE)
	ThreadTypeUndefined = ThreadType(0)
)
