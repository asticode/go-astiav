package astiav

//#include <libswscale/swscale.h>
import "C"

type SoftwareScaleContextFlag int64

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libswscale/swscale.h#L59
const (
	SoftwareScaleContextFlagArea         = SoftwareScaleContextFlag(C.SWS_AREA)
	SoftwareScaleContextFlagBicubic      = SoftwareScaleContextFlag(C.SWS_BICUBIC)
	SoftwareScaleContextFlagBicublin     = SoftwareScaleContextFlag(C.SWS_BICUBLIN)
	SoftwareScaleContextFlagBilinear     = SoftwareScaleContextFlag(C.SWS_BILINEAR)
	SoftwareScaleContextFlagFastBilinear = SoftwareScaleContextFlag(C.SWS_FAST_BILINEAR)
	SoftwareScaleContextFlagGauss        = SoftwareScaleContextFlag(C.SWS_GAUSS)
	SoftwareScaleContextFlagLanczos      = SoftwareScaleContextFlag(C.SWS_LANCZOS)
	SoftwareScaleContextFlagPoint        = SoftwareScaleContextFlag(C.SWS_POINT)
	SoftwareScaleContextFlagSinc         = SoftwareScaleContextFlag(C.SWS_SINC)
	SoftwareScaleContextFlagSpline       = SoftwareScaleContextFlag(C.SWS_SPLINE)
	SoftwareScaleContextFlagX            = SoftwareScaleContextFlag(C.SWS_X)
)
