package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/avio_8h.html#a21e61cb486bd1588eb7f775998cf8c77
type IOContextFlag int64

const (
	IOContextFlagRead      = IOContextFlag(C.AVIO_FLAG_READ)
	IOContextFlagWrite     = IOContextFlag(C.AVIO_FLAG_WRITE)
	IOContextFlagReadWrite = IOContextFlag(C.AVIO_FLAG_READ_WRITE)
	IOContextFlagNonBlock  = IOContextFlag(C.AVIO_FLAG_NONBLOCK)
	IOContextFlagDirect    = IOContextFlag(C.AVIO_FLAG_DIRECT)
)
