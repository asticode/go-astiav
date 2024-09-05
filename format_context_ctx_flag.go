package astiav

//#include <libavformat/avformat.h>
import "C"

type FormatContextCtxFlag int64

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L1153
const (
	FormatContextCtxFlagNoHeader   = FormatContextCtxFlag(C.AVFMTCTX_NOHEADER)
	FormatContextCtxFlagUnseekable = FormatContextCtxFlag(C.AVFMTCTX_UNSEEKABLE)
)
