package astiav

//#cgo pkg-config: libavformat
//#include <libavformat/avformat.h>
import "C"

type SeekFlag int

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L2277
const (
	SeekFlagAny      = SeekFlag(C.AVSEEK_FLAG_ANY)
	SeekFlagBackward = SeekFlag(C.AVSEEK_FLAG_BACKWARD)
	SeekFlagByte     = SeekFlag(C.AVSEEK_FLAG_BYTE)
	SeekFlagFrame    = SeekFlag(C.AVSEEK_FLAG_FRAME)
)
