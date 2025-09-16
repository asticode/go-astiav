package astiav

//#include <libavcodec/avcodec.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/avcodec_8h.html#a84a993ea19afa2cbda45b3283a598fe6
type Level int

const (
	// LevelUnknown removed in FFmpeg 8.0 - FF_LEVEL_UNKNOWN was removed
	// Using -99 as the equivalent value that was previously defined
	LevelUnknown = Level(-99)
)
