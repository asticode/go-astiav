package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/7.1/avformat_8h.html#ac736f8f4afc930ca1cda0b43638cc678
type SeekFlag int64

const (
	SeekFlagAny      = SeekFlag(C.AVSEEK_FLAG_ANY)
	SeekFlagBackward = SeekFlag(C.AVSEEK_FLAG_BACKWARD)
	SeekFlagByte     = SeekFlag(C.AVSEEK_FLAG_BYTE)
	SeekFlagFrame    = SeekFlag(C.AVSEEK_FLAG_FRAME)
)
