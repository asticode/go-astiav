package astiav

//#include <libavcodec/avcodec.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/avcodec_8h.html#a84a993ea19afa2cbda45b3283a598fe6
type Level int

const (
	LevelUnknown = Level(C.FF_LEVEL_UNKNOWN)
)
