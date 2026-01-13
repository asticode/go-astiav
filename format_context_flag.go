package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/avformat_8h.html#a69e2c8bc119c0245ff6092f9db4d12ae
type FormatContextFlag int64

const (
	FormatContextFlagAutoBsf        = FormatContextFlag(C.AVFMT_FLAG_AUTO_BSF)
	FormatContextFlagBitexact       = FormatContextFlag(C.AVFMT_FLAG_BITEXACT)
	FormatContextFlagCustomIo       = FormatContextFlag(C.AVFMT_FLAG_CUSTOM_IO)
	FormatContextFlagDiscardCorrupt = FormatContextFlag(C.AVFMT_FLAG_DISCARD_CORRUPT)
	FormatContextFlagFastSeek       = FormatContextFlag(C.AVFMT_FLAG_FAST_SEEK)
	FormatContextFlagFlushPackets   = FormatContextFlag(C.AVFMT_FLAG_FLUSH_PACKETS)
	FormatContextFlagGenPts         = FormatContextFlag(C.AVFMT_FLAG_GENPTS)
	FormatContextFlagIgnDts         = FormatContextFlag(C.AVFMT_FLAG_IGNDTS)
	FormatContextFlagIgnidx         = FormatContextFlag(C.AVFMT_FLAG_IGNIDX)
	FormatContextFlagNobuffer       = FormatContextFlag(C.AVFMT_FLAG_NOBUFFER)
	FormatContextFlagNofillin       = FormatContextFlag(C.AVFMT_FLAG_NOFILLIN)
	FormatContextFlagNonblock       = FormatContextFlag(C.AVFMT_FLAG_NONBLOCK)
	FormatContextFlagNoparse        = FormatContextFlag(C.AVFMT_FLAG_NOPARSE)
	FormatContextFlagSortDts        = FormatContextFlag(C.AVFMT_FLAG_SORT_DTS)
)
