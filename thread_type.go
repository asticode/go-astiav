package astiav

//#include <libavcodec/avcodec.h>
import "C"

type ThreadType int

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/avcodec.h#L1451
const (
	ThreadTypeFrame     = ThreadType(C.FF_THREAD_FRAME)
	ThreadTypeSlice     = ThreadType(C.FF_THREAD_SLICE)
	ThreadTypeUndefined = ThreadType(0)
)
