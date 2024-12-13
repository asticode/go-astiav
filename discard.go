package astiav

//#include <libavcodec/defs.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/group__lavc__decoding.html#ga352363bce7d3ed82c101b3bc001d1c16
type Discard C.enum_AVDiscard

const (
	DiscardNone          = Discard(C.AVDISCARD_NONE)     // discard nothing
	DiscardDefault       = Discard(C.AVDISCARD_DEFAULT)  // discard useless packets like 0 size packets in avi
	DiscardNonRef        = Discard(C.AVDISCARD_NONREF)   // discard all non reference
	DiscardBidirectional = Discard(C.AVDISCARD_BIDIR)    // discard all bidirectional frames
	DiscardNonIntra      = Discard(C.AVDISCARD_NONINTRA) // discard all non intra frames
	DiscardNonKey        = Discard(C.AVDISCARD_NONKEY)   // discard all frames except keyframes
	DiscardAll           = Discard(C.AVDISCARD_ALL)      // discard all
)
