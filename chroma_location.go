package astiav

//#include <libavutil/pixfmt.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/pixfmt_8h.html#a1f86ed1b6a420faccacf77c98db6c1ff
type ChromaLocation C.enum_AVChromaLocation

const (
	ChromaLocationUnspecified = ChromaLocation(C.AVCHROMA_LOC_UNSPECIFIED)
	ChromaLocationLeft        = ChromaLocation(C.AVCHROMA_LOC_LEFT)
	ChromaLocationCenter      = ChromaLocation(C.AVCHROMA_LOC_CENTER)
	ChromaLocationTopleft     = ChromaLocation(C.AVCHROMA_LOC_TOPLEFT)
	ChromaLocationTop         = ChromaLocation(C.AVCHROMA_LOC_TOP)
	ChromaLocationBottomleft  = ChromaLocation(C.AVCHROMA_LOC_BOTTOMLEFT)
	ChromaLocationBottom      = ChromaLocation(C.AVCHROMA_LOC_BOTTOM)
	ChromaLocationNb          = ChromaLocation(C.AVCHROMA_LOC_NB)
)
