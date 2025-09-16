package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/avformat_8h.html#a19485b8b52e579db560875e9a1e44e7a
type FormatEventFlag int64

const (
	FormatEventFlagMetadataUpdated = FormatEventFlag(C.AVFMT_EVENT_FLAG_METADATA_UPDATED)
)
