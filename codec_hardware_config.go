package astiav

//#cgo pkg-config: libavcodec
//#include <libavcodec/avcodec.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/codec.h#L460
type CodecHardwareConfig struct {
	c *C.struct_AVCodecHWConfig
}

func newCodecHardwareConfigFromC(c *C.struct_AVCodecHWConfig) CodecHardwareConfig {
	return CodecHardwareConfig{c: c}
}

func (chc CodecHardwareConfig) HardwareDeviceType() HardwareDeviceType {
	return HardwareDeviceType(chc.c.device_type)
}

func (chc CodecHardwareConfig) MethodFlags() CodecHardwareConfigMethodFlags {
	return CodecHardwareConfigMethodFlags(chc.c.methods)
}

func (chc CodecHardwareConfig) PixelFormat() PixelFormat {
	return PixelFormat(chc.c.pix_fmt)
}
