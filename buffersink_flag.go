package astiav

//#include <libavfilter/buffersink.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/group__lavfi__buffersink.html#ga9453fc0e81d30237080b51575da0f0d8
type BuffersinkFlag int64

const (
	BuffersinkFlagPeek      = BuffersinkFlag(C.AV_BUFFERSINK_FLAG_PEEK)
	BuffersinkFlagNoRequest = BuffersinkFlag(C.AV_BUFFERSINK_FLAG_NO_REQUEST)
)
