package astiav

//#include <libswscale/swscale.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#ga6110064d9edfbec77ca5c3279cb75c31
type SoftwareScaleContextFlag int64

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
