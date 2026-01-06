package astiav

//#include <libavcodec/avcodec.h>
import "C"

// https://ffmpeg.org/doxygen/7.1/defs_8h.html#a96808e3862c53c7edb4ace1b2f3e544f
type StrictStdCompliance int

const (
	StrictStdComplianceVeryStrict   = StrictStdCompliance(C.FF_COMPLIANCE_VERY_STRICT)
	StrictStdComplianceStrict       = StrictStdCompliance(C.FF_COMPLIANCE_STRICT)
	StrictStdComplianceNormal       = StrictStdCompliance(C.FF_COMPLIANCE_NORMAL)
	StrictStdComplianceUnofficial   = StrictStdCompliance(C.FF_COMPLIANCE_UNOFFICIAL)
	StrictStdComplianceExperimental = StrictStdCompliance(C.FF_COMPLIANCE_EXPERIMENTAL)
)
