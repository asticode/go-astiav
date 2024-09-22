package astiav

//#include <libavcodec/avcodec.h>
import "C"

type Level int

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/avcodec.h#L1652
const (
	LevelUnknown = Level(C.FF_LEVEL_UNKNOWN)
)
