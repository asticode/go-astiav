package astiav

//#include <libavcodec/avcodec.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/structAVCodecHWConfig.html
type CodecHardwareConfig struct {
	c *C.AVCodecHWConfig
}

func newCodecHardwareConfigFromC(c *C.AVCodecHWConfig) CodecHardwareConfig {
	return CodecHardwareConfig{c: c}
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecHWConfig.html#a1474cb73c1f41e377dc5070ae373ac40
func (chc CodecHardwareConfig) HardwareDeviceType() HardwareDeviceType {
	return HardwareDeviceType(chc.c.device_type)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecHWConfig.html#a208c924c3f626b01bf2020eef9eb4905
func (chc CodecHardwareConfig) MethodFlags() CodecHardwareConfigMethodFlags {
	return CodecHardwareConfigMethodFlags(chc.c.methods)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecHWConfig.html#a9352b11d6d6b315fe3c61b65447d5174
func (chc CodecHardwareConfig) PixelFormat() PixelFormat {
	return PixelFormat(chc.c.pix_fmt)
}
