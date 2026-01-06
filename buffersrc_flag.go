package astiav

//#include <libavfilter/buffersrc.h>
import "C"

// https://ffmpeg.org/doxygen/7.1/group__lavfi__buffersrc.html#ggac998c8f2ed6a3c0d0ae07b3822c16f9da6efcf61145ec6d60d3a773fcd0797872
type BuffersrcFlag int64

const (
	BuffersrcFlagNoCheckFormat = BuffersrcFlag(C.AV_BUFFERSRC_FLAG_NO_CHECK_FORMAT)
	BuffersrcFlagPush          = BuffersrcFlag(C.AV_BUFFERSRC_FLAG_PUSH)
	BuffersrcFlagKeepRef       = BuffersrcFlag(C.AV_BUFFERSRC_FLAG_KEEP_REF)
)
