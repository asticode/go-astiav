package astiav

//#include <libavutil/pixdesc.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/structAVPixFmtDescriptor.html
type PixelFormatDescriptor struct {
	c *C.AVPixFmtDescriptor
}

func newPixelFormatDescriptorFromC(c *C.AVPixFmtDescriptor) *PixelFormatDescriptor {
	if c == nil {
		return nil
	}
	return &PixelFormatDescriptor{c: c}
}

// https://ffmpeg.org/doxygen/8.0/structAVPixFmtDescriptor.html#a10736c3f1288eb87b23ede3ffdefb435
func (pfd *PixelFormatDescriptor) Name() string {
	return C.GoString(pfd.c.name)
}

// https://ffmpeg.org/doxygen/8.0/structAVPixFmtDescriptor.html#a5047d1e6b045f637345dbc305bf4357d
func (pfd *PixelFormatDescriptor) Flags() PixelFormatDescriptorFlags {
	return PixelFormatDescriptorFlags(pfd.c.flags)
}
