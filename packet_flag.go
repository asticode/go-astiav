package astiav

//#include <libavcodec/avcodec.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/group__lavc__packet.html#ga75603d7c2b8adf5829f4fd2fb860168f
type PacketFlag int64

const (
	PacketFlagCorrupt = PacketFlag(C.AV_PKT_FLAG_CORRUPT)
	PacketFlagDiscard = PacketFlag(C.AV_PKT_FLAG_DISCARD)
	PacketFlagKey     = PacketFlag(C.AV_PKT_FLAG_KEY)
)
