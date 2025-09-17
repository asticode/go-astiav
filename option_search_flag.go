package astiav

//#include <libavutil/opt.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/group__opt__mng.html#ga25801ba4fc9b5313eb33ec84e082dd72
type OptionSearchFlag int64

const (
	OptionSearchFlagChildren   = OptionSearchFlag(C.AV_OPT_SEARCH_CHILDREN)
	OptionSearchFlagFakeObject = OptionSearchFlag(C.AV_OPT_SEARCH_FAKE_OBJ)
)
