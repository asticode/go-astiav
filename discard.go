package astiav

//#include <libavcodec/defs.h>
import "C"

// https://ffmpeg.org/doxygen/7.1/group__lavc__decoding.html#ga352363bce7d3ed82c101b3bc001d1c16
type Discard C.enum_AVDiscard

const (
	DiscardNone          = Discard(C.AVDISCARD_NONE)
	DiscardDefault       = Discard(C.AVDISCARD_DEFAULT)
	DiscardNonRef        = Discard(C.AVDISCARD_NONREF)
	DiscardBidirectional = Discard(C.AVDISCARD_BIDIR)
	DiscardNonIntra      = Discard(C.AVDISCARD_NONINTRA)
	DiscardNonKey        = Discard(C.AVDISCARD_NONKEY)
	DiscardAll           = Discard(C.AVDISCARD_ALL)
)
