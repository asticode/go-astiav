package astiav

//#include <libavutil/frame.h>
import "C"

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame__flags.html#gafe155269fc8dc3a484490bd19b86cc40
type FrameFlag int64

const (
	FrameFlagCorrupt       = FrameFlag(C.AV_FRAME_FLAG_CORRUPT)
	FrameFlagDiscard       = FrameFlag(C.AV_FRAME_FLAG_DISCARD)
	FrameFlagInterlaced    = FrameFlag(C.AV_FRAME_FLAG_INTERLACED)
	FrameFlagKey           = FrameFlag(C.AV_FRAME_FLAG_KEY)
	FrameFlagTopFieldFirst = FrameFlag(C.AV_FRAME_FLAG_TOP_FIELD_FIRST)
)
