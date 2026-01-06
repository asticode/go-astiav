package astiav

//#include <libswscale/swscale.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#gga71e02a84731a9d51c3b27c1ecb15f4f4a0e2c41ecbb835ddcafb698b324842cb2
type SoftwareScaleContextIntent C.enum_SwsIntent

const (
	SoftwareScaleContextIntentPerceptual           = SoftwareScaleContextIntent(C.SWS_INTENT_PERCEPTUAL)
	SoftwareScaleContextIntentRelativeColorimetric = SoftwareScaleContextIntent(C.SWS_INTENT_RELATIVE_COLORIMETRIC)
	SoftwareScaleContextIntentSaturation           = SoftwareScaleContextIntent(C.SWS_INTENT_SATURATION)
	SoftwareScaleContextIntentAbsoluteColorimetric = SoftwareScaleContextIntent(C.SWS_INTENT_ABSOLUTE_COLORIMETRIC)
)
