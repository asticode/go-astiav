package astiav

//#include <libavcodec/avcodec.h>
import "C"

type StrictStdCompliance int

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/avcodec.h#L1281
const (
	StrictStdComplianceVeryStrict   = StrictStdCompliance(C.FF_COMPLIANCE_VERY_STRICT)
	StrictStdComplianceStrict       = StrictStdCompliance(C.FF_COMPLIANCE_STRICT)
	StrictStdComplianceNormal       = StrictStdCompliance(C.FF_COMPLIANCE_NORMAL)
	StrictStdComplianceUnofficial   = StrictStdCompliance(C.FF_COMPLIANCE_UNOFFICIAL)
	StrictStdComplianceExperimental = StrictStdCompliance(C.FF_COMPLIANCE_EXPERIMENTAL)
)
