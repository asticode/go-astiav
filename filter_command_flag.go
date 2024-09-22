package astiav

//#include <libavfilter/avfilter.h>
import "C"

type FilterCommandFlag int64

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavfilter/avfilter.h#L739
const (
	FilterCommandFlagOne  = FilterCommandFlag(C.AVFILTER_CMD_FLAG_ONE)
	FilterCommandFlagFast = FilterCommandFlag(C.AVFILTER_CMD_FLAG_FAST)
)
