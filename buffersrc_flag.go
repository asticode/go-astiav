package astiav

//#include <libavfilter/buffersrc.h>
import "C"

type BuffersrcFlag int64

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavfilter/buffersrc.h#L36
const (
	BuffersrcFlagNoCheckFormat = BuffersrcFlag(C.AV_BUFFERSRC_FLAG_NO_CHECK_FORMAT)
	BuffersrcFlagPush          = BuffersrcFlag(C.AV_BUFFERSRC_FLAG_PUSH)
	BuffersrcFlagKeepRef       = BuffersrcFlag(C.AV_BUFFERSRC_FLAG_KEEP_REF)
)
