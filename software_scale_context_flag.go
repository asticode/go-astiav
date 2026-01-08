package astiav

//#include <libswscale/swscale.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#ggade664a46fb2652e6050985ebcd316798ae366a8f172c83a868c4a149ef844f2a7
type SoftwareScaleContextFlag int64

const (
	SoftwareScaleContextFlagArea           = SoftwareScaleContextFlag(C.SWS_AREA)
	SoftwareScaleContextFlagBicubic        = SoftwareScaleContextFlag(C.SWS_BICUBIC)
	SoftwareScaleContextFlagBicublin       = SoftwareScaleContextFlag(C.SWS_BICUBLIN)
	SoftwareScaleContextFlagBilinear       = SoftwareScaleContextFlag(C.SWS_BILINEAR)
	SoftwareScaleContextFlagFastBilinear   = SoftwareScaleContextFlag(C.SWS_FAST_BILINEAR)
	SoftwareScaleContextFlagGauss          = SoftwareScaleContextFlag(C.SWS_GAUSS)
	SoftwareScaleContextFlagLanczos        = SoftwareScaleContextFlag(C.SWS_LANCZOS)
	SoftwareScaleContextFlagPoint          = SoftwareScaleContextFlag(C.SWS_POINT)
	SoftwareScaleContextFlagSinc           = SoftwareScaleContextFlag(C.SWS_SINC)
	SoftwareScaleContextFlagSpline         = SoftwareScaleContextFlag(C.SWS_SPLINE)
	SoftwareScaleContextFlagX              = SoftwareScaleContextFlag(C.SWS_X)
	SoftwareScaleContextFlagStrict         = SoftwareScaleContextFlag(C.SWS_STRICT)
	SoftwareScaleContextFlagPrintInfo      = SoftwareScaleContextFlag(C.SWS_PRINT_INFO)
	SoftwareScaleContextFlagFullChrHInt    = SoftwareScaleContextFlag(C.SWS_FULL_CHR_H_INT)
	SoftwareScaleContextFlagFullChrHInp    = SoftwareScaleContextFlag(C.SWS_FULL_CHR_H_INP)
	SoftwareScaleContextFlagAccurateRnd    = SoftwareScaleContextFlag(C.SWS_ACCURATE_RND)
	SoftwareScaleContextFlagBitexact       = SoftwareScaleContextFlag(C.SWS_BITEXACT)
	SoftwareScaleContextFlagDirectBgr      = SoftwareScaleContextFlag(C.SWS_DIRECT_BGR)
	SoftwareScaleContextFlagErrorDiffusion = SoftwareScaleContextFlag(C.SWS_ERROR_DIFFUSION)
)
