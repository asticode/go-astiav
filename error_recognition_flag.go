package astiav

//#include <libavcodec/defs.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/defs_8h.html#a16258b3283a3d34c893dd9a39e29e369
type ErrorRecognitionFlag int64

const (
	ErrorRecognitionFlagAggressive = ErrorRecognitionFlag(C.AV_EF_AGGRESSIVE)
	ErrorRecognitionFlagBitstream  = ErrorRecognitionFlag(C.AV_EF_BITSTREAM)
	ErrorRecognitionFlagBuffer     = ErrorRecognitionFlag(C.AV_EF_BUFFER)
	ErrorRecognitionFlagCareful    = ErrorRecognitionFlag(C.AV_EF_CAREFUL)
	ErrorRecognitionFlagCompliant  = ErrorRecognitionFlag(C.AV_EF_COMPLIANT)
	ErrorRecognitionFlagCRCCheck   = ErrorRecognitionFlag(C.AV_EF_CRCCHECK)
	ErrorRecognitionFlagExplode    = ErrorRecognitionFlag(C.AV_EF_EXPLODE)
	ErrorRecognitionFlagIgnoreErr  = ErrorRecognitionFlag(C.AV_EF_IGNORE_ERR)
)
