package astiav

//#include <libavutil/pixfmt.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/pixfmt.h#L616
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
