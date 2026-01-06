package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/7.1/avformat_8h.html#a69e2c8bc119c0245ff6092f9db4d12ae
type FormatContextFlag int64

const (
	FormatContextFlagGenPts         = FormatContextFlag(C.AVFMT_FLAG_GENPTS)
	FormatContextFlagIgnidx         = FormatContextFlag(C.AVFMT_FLAG_IGNIDX)
	FormatContextFlagNonblock       = FormatContextFlag(C.AVFMT_FLAG_NONBLOCK)
	FormatContextFlagIgnDts         = FormatContextFlag(C.AVFMT_FLAG_IGNDTS)
	FormatContextFlagNofillin       = FormatContextFlag(C.AVFMT_FLAG_NOFILLIN)
	FormatContextFlagNoparse        = FormatContextFlag(C.AVFMT_FLAG_NOPARSE)
	FormatContextFlagNobuffer       = FormatContextFlag(C.AVFMT_FLAG_NOBUFFER)
	FormatContextFlagCustomIo       = FormatContextFlag(C.AVFMT_FLAG_CUSTOM_IO)
	FormatContextFlagDiscardCorrupt = FormatContextFlag(C.AVFMT_FLAG_DISCARD_CORRUPT)
	FormatContextFlagFlushPackets   = FormatContextFlag(C.AVFMT_FLAG_FLUSH_PACKETS)
	FormatContextFlagBitexact       = FormatContextFlag(C.AVFMT_FLAG_BITEXACT)
	FormatContextFlagSortDts        = FormatContextFlag(C.AVFMT_FLAG_SORT_DTS)
	FormatContextFlagFastSeek       = FormatContextFlag(C.AVFMT_FLAG_FAST_SEEK)
	FormatContextFlagShortest       = FormatContextFlag(C.AVFMT_FLAG_SHORTEST)
	FormatContextFlagAutoBsf        = FormatContextFlag(C.AVFMT_FLAG_AUTO_BSF)
)
