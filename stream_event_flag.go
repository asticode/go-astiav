package astiav

//#cgo pkg-config: libavformat
//#include <libavformat/avformat.h>
import "C"

type StreamEventFlag int64

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L1070
const (
	StreamEventFlagMetadataUpdated = StreamEventFlag(C.AVSTREAM_EVENT_FLAG_METADATA_UPDATED)
)
