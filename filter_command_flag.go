package astiav

//#include <libavfilter/avfilter.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#gace41bae000b621fc8804a93ce9f2d6e9
type FilterCommandFlag int64

const (
	FilterCommandFlagOne  = FilterCommandFlag(C.AVFILTER_CMD_FLAG_ONE)
	FilterCommandFlagFast = FilterCommandFlag(C.AVFILTER_CMD_FLAG_FAST)
)
