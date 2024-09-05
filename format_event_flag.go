package astiav

//#include <libavformat/avformat.h>
import "C"

type FormatEventFlag int64

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L1519
const (
	FormatEventFlagMetadataUpdated = FormatEventFlag(C.AVFMT_EVENT_FLAG_METADATA_UPDATED)
)
