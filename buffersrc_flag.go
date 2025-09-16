package astiav

//#include <libavfilter/buffersrc.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/group__lavfi__buffersrc.html#gga9e9505a91a84992d04ba4e85217fb4e4a6efcf61145ec6d60d3a773fcd0797872
type BuffersrcFlag int64

const (
	BuffersrcFlagNoCheckFormat = BuffersrcFlag(C.AV_BUFFERSRC_FLAG_NO_CHECK_FORMAT)
	BuffersrcFlagPush          = BuffersrcFlag(C.AV_BUFFERSRC_FLAG_PUSH)
	BuffersrcFlagKeepRef       = BuffersrcFlag(C.AV_BUFFERSRC_FLAG_KEEP_REF)
)
