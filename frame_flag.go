package astiav

//#include <libavutil/frame.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/group__lavu__frame__flags.html#gadddbce4ec0cc2ad4298cf6f266f97f6a
type FrameFlag int64

const (
	FrameFlagCorrupt       = FrameFlag(C.AV_FRAME_FLAG_CORRUPT)
	FrameFlagDiscard       = FrameFlag(C.AV_FRAME_FLAG_DISCARD)
	FrameFlagInterlaced    = FrameFlag(C.AV_FRAME_FLAG_INTERLACED)
	FrameFlagKey           = FrameFlag(C.AV_FRAME_FLAG_KEY)
	FrameFlagLossless      = FrameFlag(C.AV_FRAME_FLAG_LOSSLESS)
	FrameFlagTopFieldFirst = FrameFlag(C.AV_FRAME_FLAG_TOP_FIELD_FIRST)
)
