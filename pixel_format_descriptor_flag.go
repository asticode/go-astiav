package astiav

//#include <libavutil/pixdesc.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/pixdesc_8h.html#ac7c7d0be16fb9b6f05b3e0d463cd037b
type PixelFormatDescriptorFlag int64

const (
	PixelFormatDescriptorFlagBe        = PixelFormatDescriptorFlag(C.AV_PIX_FMT_FLAG_BE)
	PixelFormatDescriptorFlagPal       = PixelFormatDescriptorFlag(C.AV_PIX_FMT_FLAG_PAL)
	PixelFormatDescriptorFlagBitStream = PixelFormatDescriptorFlag(C.AV_PIX_FMT_FLAG_BITSTREAM)
	PixelFormatDescriptorFlagHwAccel   = PixelFormatDescriptorFlag(C.AV_PIX_FMT_FLAG_HWACCEL)
	PixelFormatDescriptorFlagPlanar    = PixelFormatDescriptorFlag(C.AV_PIX_FMT_FLAG_PLANAR)
	PixelFormatDescriptorFlagRgb       = PixelFormatDescriptorFlag(C.AV_PIX_FMT_FLAG_RGB)
	PixelFormatDescriptorFlagAlpha     = PixelFormatDescriptorFlag(C.AV_PIX_FMT_FLAG_ALPHA)
	PixelFormatDescriptorFlagBayer     = PixelFormatDescriptorFlag(C.AV_PIX_FMT_FLAG_BAYER)
	PixelFormatDescriptorFlagFloat     = PixelFormatDescriptorFlag(C.AV_PIX_FMT_FLAG_FLOAT)
	PixelFormatDescriptorFlagXyz       = PixelFormatDescriptorFlag(C.AV_PIX_FMT_FLAG_XYZ)
)
