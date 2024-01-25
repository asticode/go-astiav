package astiav

//#cgo pkg-config: libswscale
//#include <libswscale/swscale.h>
import "C"

type SoftwareScaleContextFlag int

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libswscale/swscale.h#L59
const (
	SoftwareScaleContextArea         = SoftwareScaleContextFlag(C.SWS_AREA)
	SoftwareScaleContextBicubic      = SoftwareScaleContextFlag(C.SWS_BICUBIC)
	SoftwareScaleContextBicublin     = SoftwareScaleContextFlag(C.SWS_BICUBLIN)
	SoftwareScaleContextBilinear     = SoftwareScaleContextFlag(C.SWS_BILINEAR)
	SoftwareScaleContextFastBilinear = SoftwareScaleContextFlag(C.SWS_FAST_BILINEAR)
	SoftwareScaleContextGauss        = SoftwareScaleContextFlag(C.SWS_GAUSS)
	SoftwareScaleContextLanczos      = SoftwareScaleContextFlag(C.SWS_LANCZOS)
	SoftwareScaleContextPoint        = SoftwareScaleContextFlag(C.SWS_POINT)
	SoftwareScaleContextSinc         = SoftwareScaleContextFlag(C.SWS_SINC)
	SoftwareScaleContextSpline       = SoftwareScaleContextFlag(C.SWS_SPLINE)
	SoftwareScaleContextX            = SoftwareScaleContextFlag(C.SWS_X)
)
