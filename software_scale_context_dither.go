package astiav

//#include <libswscale/swscale.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/group__libsws.html#gga7305b190c33d84c6575799b4b593b3b9afa87c1d59268142efb82b43f0382f927
type SoftwareScaleContextDither C.enum_SwsDither

const (
	SoftwareScaleContextDitherNone    = SoftwareScaleContextDither(C.SWS_DITHER_NONE)
	SoftwareScaleContextDitherAuto    = SoftwareScaleContextDither(C.SWS_DITHER_AUTO)
	SoftwareScaleContextDitherBayer   = SoftwareScaleContextDither(C.SWS_DITHER_BAYER)
	SoftwareScaleContextDitherEd      = SoftwareScaleContextDither(C.SWS_DITHER_ED)
	SoftwareScaleContextDitherADither = SoftwareScaleContextDither(C.SWS_DITHER_A_DITHER)
	SoftwareScaleContextDitherXDither = SoftwareScaleContextDither(C.SWS_DITHER_X_DITHER)
)
