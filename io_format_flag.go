package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/8.1/avformat_8h.html#a752cce390d480521919aa5d8be24ac0b
type IOFormatFlag int64

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
	// IOFormatFlagAllowFlush removed in FFmpeg 8.0 - AVFMT_ALLOW_FLUSH deprecated and removed
	IOFormatFlagTsNonstrict = IOFormatFlag(C.AVFMT_TS_NONSTRICT)
	IOFormatFlagTsNegative  = IOFormatFlag(C.AVFMT_TS_NEGATIVE)
	IOFormatFlagSeekToPts   = IOFormatFlag(C.AVFMT_SEEK_TO_PTS)
)
