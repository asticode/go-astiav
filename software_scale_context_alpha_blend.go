package astiav

//#include <libswscale/swscale.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#gga878a03b3ee6c5066ee77c5c02b8f1016a8661b3199ee684757f139875e8541803
type SoftwareScaleContextAlphaBlend C.enum_SwsAlphaBlend

const (
	SoftwareScaleContextAlphaBlendNone         = SoftwareScaleContextAlphaBlend(C.SWS_ALPHA_BLEND_NONE)
	SoftwareScaleContextAlphaBlendUniform      = SoftwareScaleContextAlphaBlend(C.SWS_ALPHA_BLEND_UNIFORM)
	SoftwareScaleContextAlphaBlendCheckerboard = SoftwareScaleContextAlphaBlend(C.SWS_ALPHA_BLEND_CHECKERBOARD)
)
