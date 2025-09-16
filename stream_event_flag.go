package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/avformat_8h.html#ab3a5958310f614671f5030ed10753ba9
type StreamEventFlag int64

const (
	StreamEventFlagMetadataUpdated = StreamEventFlag(C.AVSTREAM_EVENT_FLAG_METADATA_UPDATED)
)
