package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/avformat_8h.html#a42e3c3d72e561fdc501613962fccc4aa
type FormatContextCtxFlag int64

const (
	FormatContextCtxFlagNoHeader   = FormatContextCtxFlag(C.AVFMTCTX_NOHEADER)
	FormatContextCtxFlagUnseekable = FormatContextCtxFlag(C.AVFMTCTX_UNSEEKABLE)
)
