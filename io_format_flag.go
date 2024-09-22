package astiav

//#include <libavformat/avformat.h>
import "C"

type IOFormatFlag int64

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L464
const (
	IOFormatFlagNofile       = IOFormatFlag(C.AVFMT_NOFILE)
	IOFormatFlagNeednumber   = IOFormatFlag(C.AVFMT_NEEDNUMBER)
	IOFormatFlagShowIds      = IOFormatFlag(C.AVFMT_SHOW_IDS)
	IOFormatFlagGlobalheader = IOFormatFlag(C.AVFMT_GLOBALHEADER)
	IOFormatFlagNotimestamps = IOFormatFlag(C.AVFMT_NOTIMESTAMPS)
	IOFormatFlagGenericIndex = IOFormatFlag(C.AVFMT_GENERIC_INDEX)
	IOFormatFlagTsDiscont    = IOFormatFlag(C.AVFMT_TS_DISCONT)
	IOFormatFlagVariableFps  = IOFormatFlag(C.AVFMT_VARIABLE_FPS)
	IOFormatFlagNodimensions = IOFormatFlag(C.AVFMT_NODIMENSIONS)
	IOFormatFlagNostreams    = IOFormatFlag(C.AVFMT_NOSTREAMS)
	IOFormatFlagNobinsearch  = IOFormatFlag(C.AVFMT_NOBINSEARCH)
	IOFormatFlagNogensearch  = IOFormatFlag(C.AVFMT_NOGENSEARCH)
	IOFormatFlagNoByteSeek   = IOFormatFlag(C.AVFMT_NO_BYTE_SEEK)
	IOFormatFlagAllowFlush   = IOFormatFlag(C.AVFMT_ALLOW_FLUSH)
	IOFormatFlagTsNonstrict  = IOFormatFlag(C.AVFMT_TS_NONSTRICT)
	IOFormatFlagTsNegative   = IOFormatFlag(C.AVFMT_TS_NEGATIVE)
	IOFormatFlagSeekToPts    = IOFormatFlag(C.AVFMT_SEEK_TO_PTS)
)
