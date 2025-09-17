package astiav

//#include <libavutil/hwcontext.h>
import "C"
import "unsafe"

// https://ffmpeg.org/doxygen/8.0/structAVHWFramesConstraints.html
type HardwareFramesConstraints struct {
	c *C.AVHWFramesConstraints
}

func newHardwareFramesConstraintsFromC(c *C.AVHWFramesConstraints) *HardwareFramesConstraints {
	if c == nil {
		return nil
	}
	return &HardwareFramesConstraints{c: c}
}

func (hfc *HardwareFramesConstraints) pixelFormats(formats *C.enum_AVPixelFormat) (o []PixelFormat) {
	if formats == nil {
		return nil
	}
	size := unsafe.Sizeof(*formats)
	for i := 0; ; i++ {
		p := *(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(formats)) + uintptr(i)*size))
		if p == C.AV_PIX_FMT_NONE {
			break
		}
		o = append(o, PixelFormat(p))
	}
	return
}

// https://ffmpeg.org/doxygen/8.0/structAVHWFramesConstraints.html#a4258bbe81f927b76b7ca5af44ba7ef6b
func (hfc *HardwareFramesConstraints) ValidHardwarePixelFormats() (o []PixelFormat) {
	return hfc.pixelFormats(hfc.c.valid_hw_formats)
}

// https://ffmpeg.org/doxygen/8.0/structAVHWFramesConstraints.html#aabea88093c6f85d6185ffb0852a2217f
func (hfc *HardwareFramesConstraints) ValidSoftwarePixelFormats() (o []PixelFormat) {
	return hfc.pixelFormats(hfc.c.valid_sw_formats)
}

// https://ffmpeg.org/doxygen/8.0/structAVHWFramesConstraints.html#af220776925452091085139081d5d7251
func (hfc *HardwareFramesConstraints) MinWidth() int {
	return int(hfc.c.min_width)
}

// https://ffmpeg.org/doxygen/8.0/structAVHWFramesConstraints.html#a3f1aec6d1c90f77837875c2a3598be46
func (hfc *HardwareFramesConstraints) MinHeight() int {
	return int(hfc.c.min_height)
}

// https://ffmpeg.org/doxygen/8.0/structAVHWFramesConstraints.html#a34e06e3397af2b83de9d78f893bf4168
func (hfc *HardwareFramesConstraints) MaxWidth() int {
	return int(hfc.c.max_width)
}

// https://ffmpeg.org/doxygen/8.0/structAVHWFramesConstraints.html#af5d3a683727f7b92abca7b7114d4e15c
func (hfc *HardwareFramesConstraints) MaxHeight() int {
	return int(hfc.c.max_height)
}

// https://ffmpeg.org/doxygen/8.0/hwcontext_8c.html#a29da7fa7ffa73266d1cbfccb116ed634
func (hfc *HardwareFramesConstraints) Free() {
	if hfc.c != nil {
		C.av_hwframe_constraints_free(&hfc.c)
		hfc.c = nil
	}
}
