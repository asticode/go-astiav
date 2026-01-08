package astiav

//#include <libavcodec/avcodec.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/defs_8h.html#a3b66f83dc9a672ac570d1acfb27b1057
type Level int

const (
	LevelUnknown = Level(C.AV_LEVEL_UNKNOWN)
)
