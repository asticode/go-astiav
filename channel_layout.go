package astiav

//#cgo pkg-config: libavutil
//#include <libavutil/channel_layout.h>
/*

// Calling C.AV_CHANNEL_LAYOUT_* in Go gives a "could not determine kind of name for X" error
// therefore we need to bridge the channel layout values
AVChannelLayout *c2goChannelLayoutMono              = &(AVChannelLayout)AV_CHANNEL_LAYOUT_MONO;
AVChannelLayout *c2goChannelLayoutStereo            = &(AVChannelLayout)AV_CHANNEL_LAYOUT_STEREO;
AVChannelLayout *c2goChannelLayout2Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_2POINT1;
AVChannelLayout *c2goChannelLayout21                = &(AVChannelLayout)AV_CHANNEL_LAYOUT_2_1;
AVChannelLayout *c2goChannelLayoutSurround          = &(AVChannelLayout)AV_CHANNEL_LAYOUT_SURROUND;
AVChannelLayout *c2goChannelLayout3Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_3POINT1;
AVChannelLayout *c2goChannelLayout4Point0           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_4POINT0;
AVChannelLayout *c2goChannelLayout4Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_4POINT1;
AVChannelLayout *c2goChannelLayout22                = &(AVChannelLayout)AV_CHANNEL_LAYOUT_2_2;
AVChannelLayout *c2goChannelLayoutQuad              = &(AVChannelLayout)AV_CHANNEL_LAYOUT_QUAD;
AVChannelLayout *c2goChannelLayout5Point0           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT0;
AVChannelLayout *c2goChannelLayout5Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT1;
AVChannelLayout *c2goChannelLayout5Point0Back       = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT0_BACK;
AVChannelLayout *c2goChannelLayout5Point1Back       = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT1_BACK;
AVChannelLayout *c2goChannelLayout6Point0           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_6POINT0;
AVChannelLayout *c2goChannelLayout6Point0Front      = &(AVChannelLayout)AV_CHANNEL_LAYOUT_6POINT0_FRONT;
AVChannelLayout *c2goChannelLayoutHexagonal         = &(AVChannelLayout)AV_CHANNEL_LAYOUT_HEXAGONAL;
AVChannelLayout *c2goChannelLayout3Point1Point2     = &(AVChannelLayout)AV_CHANNEL_LAYOUT_3POINT1POINT2;
AVChannelLayout *c2goChannelLayout6Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_6POINT1;
AVChannelLayout *c2goChannelLayout6Point1Back       = &(AVChannelLayout)AV_CHANNEL_LAYOUT_6POINT1_BACK;
AVChannelLayout *c2goChannelLayout6Point1Front      = &(AVChannelLayout)AV_CHANNEL_LAYOUT_6POINT1_FRONT;
AVChannelLayout *c2goChannelLayout7Point0           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT0;
AVChannelLayout *c2goChannelLayout7Point0Front      = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT0_FRONT;
AVChannelLayout *c2goChannelLayout7Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1;
AVChannelLayout *c2goChannelLayout7Point1Wide       = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1_WIDE;
AVChannelLayout *c2goChannelLayout7Point1WideBack   = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1_WIDE_BACK;
AVChannelLayout *c2goChannelLayout5Point1Point2Back = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT1POINT2_BACK;
AVChannelLayout *c2goChannelLayoutOctagonal         = &(AVChannelLayout)AV_CHANNEL_LAYOUT_OCTAGONAL;
AVChannelLayout *c2goChannelLayoutCube              = &(AVChannelLayout)AV_CHANNEL_LAYOUT_CUBE;
AVChannelLayout *c2goChannelLayout5Point1Point4Back = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT1POINT4_BACK;
AVChannelLayout *c2goChannelLayout7Point1Point2     = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1POINT2;
AVChannelLayout *c2goChannelLayout7Point1Point4Back = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1POINT4_BACK;
AVChannelLayout *c2goChannelLayoutHexadecagonal     = &(AVChannelLayout)AV_CHANNEL_LAYOUT_HEXADECAGONAL;
AVChannelLayout *c2goChannelLayoutStereoDownmix     = &(AVChannelLayout)AV_CHANNEL_LAYOUT_STEREO_DOWNMIX;
AVChannelLayout *c2goChannelLayout22Point2          = &(AVChannelLayout)AV_CHANNEL_LAYOUT_22POINT2;
AVChannelLayout *c2goChannelLayout7Point1TopBack    = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1_TOP_BACK;

*/
import "C"
import "unsafe"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/channel_layout.h#L90
var (
	ChannelLayoutMono              = newChannelLayoutFromC(C.c2goChannelLayoutMono)
	ChannelLayoutStereo            = newChannelLayoutFromC(C.c2goChannelLayoutStereo)
	ChannelLayout2Point1           = newChannelLayoutFromC(C.c2goChannelLayout2Point1)
	ChannelLayout21                = newChannelLayoutFromC(C.c2goChannelLayout21)
	ChannelLayoutSurround          = newChannelLayoutFromC(C.c2goChannelLayoutSurround)
	ChannelLayout3Point1           = newChannelLayoutFromC(C.c2goChannelLayout3Point1)
	ChannelLayout4Point0           = newChannelLayoutFromC(C.c2goChannelLayout4Point0)
	ChannelLayout4Point1           = newChannelLayoutFromC(C.c2goChannelLayout4Point1)
	ChannelLayout22                = newChannelLayoutFromC(C.c2goChannelLayout22)
	ChannelLayoutQuad              = newChannelLayoutFromC(C.c2goChannelLayoutQuad)
	ChannelLayout5Point0           = newChannelLayoutFromC(C.c2goChannelLayout5Point0)
	ChannelLayout5Point1           = newChannelLayoutFromC(C.c2goChannelLayout5Point1)
	ChannelLayout5Point0Back       = newChannelLayoutFromC(C.c2goChannelLayout5Point0Back)
	ChannelLayout5Point1Back       = newChannelLayoutFromC(C.c2goChannelLayout5Point1Back)
	ChannelLayout6Point0           = newChannelLayoutFromC(C.c2goChannelLayout6Point0)
	ChannelLayout6Point0Front      = newChannelLayoutFromC(C.c2goChannelLayout6Point0Front)
	ChannelLayoutHexagonal         = newChannelLayoutFromC(C.c2goChannelLayoutHexagonal)
	ChannelLayout3Point1Point2     = newChannelLayoutFromC(C.c2goChannelLayout3Point1Point2)
	ChannelLayout6Point1           = newChannelLayoutFromC(C.c2goChannelLayout6Point1)
	ChannelLayout6Point1Back       = newChannelLayoutFromC(C.c2goChannelLayout6Point1Back)
	ChannelLayout6Point1Front      = newChannelLayoutFromC(C.c2goChannelLayout6Point1Front)
	ChannelLayout7Point0           = newChannelLayoutFromC(C.c2goChannelLayout7Point0)
	ChannelLayout7Point0Front      = newChannelLayoutFromC(C.c2goChannelLayout7Point0Front)
	ChannelLayout7Point1           = newChannelLayoutFromC(C.c2goChannelLayout7Point1)
	ChannelLayout7Point1Wide       = newChannelLayoutFromC(C.c2goChannelLayout7Point1Wide)
	ChannelLayout7Point1WideBack   = newChannelLayoutFromC(C.c2goChannelLayout7Point1WideBack)
	ChannelLayout5Point1Point2Back = newChannelLayoutFromC(C.c2goChannelLayout5Point1Point2Back)
	ChannelLayoutOctagonal         = newChannelLayoutFromC(C.c2goChannelLayoutOctagonal)
	ChannelLayoutCube              = newChannelLayoutFromC(C.c2goChannelLayoutCube)
	ChannelLayout5Point1Point4Back = newChannelLayoutFromC(C.c2goChannelLayout5Point1Point4Back)
	ChannelLayout7Point1Point2     = newChannelLayoutFromC(C.c2goChannelLayout7Point1Point2)
	ChannelLayout7Point1Point4Back = newChannelLayoutFromC(C.c2goChannelLayout7Point1Point4Back)
	ChannelLayoutHexadecagonal     = newChannelLayoutFromC(C.c2goChannelLayoutHexadecagonal)
	ChannelLayoutStereoDownmix     = newChannelLayoutFromC(C.c2goChannelLayoutStereoDownmix)
	ChannelLayout22Point2          = newChannelLayoutFromC(C.c2goChannelLayout22Point2)
	ChannelLayout7Point1TopBack    = newChannelLayoutFromC(C.c2goChannelLayout7Point1TopBack)
)

type ChannelLayout struct {
	c *C.struct_AVChannelLayout
}

func newChannelLayoutFromC(c *C.struct_AVChannelLayout) ChannelLayout {
	return ChannelLayout{c: c}
}

func (l ChannelLayout) NbChannels() int {
	return int(l.c.nb_channels)
}

func (l ChannelLayout) String() string {
	b := make([]byte, 1024)
	n, err := l.Describe(b)
	if err != nil {
		return ""
	}
	return string(b[:n])
}

func (l ChannelLayout) Describe(b []byte) (int, error) {
	ret := C.av_channel_layout_describe(l.c, (*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
	if ret < 0 {
		return 0, newError(ret)
	}
	if ret > 0 && b[ret-1] == '\x00' {
		ret -= 1
	}
	return int(ret), nil
}

func (l ChannelLayout) Valid() bool {
	return C.av_channel_layout_check(l.c) > 0
}

func (l ChannelLayout) Compare(l2 ChannelLayout) (equal bool, err error) {
	ret := C.av_channel_layout_compare(l.c, l2.c)
	if ret < 0 {
		return false, newError(ret)
	}
	return ret == 0, nil
}

func (l ChannelLayout) Equal(l2 ChannelLayout) bool {
	v, _ := l.Compare(l2)
	return v
}

func (l ChannelLayout) copy(dst *C.struct_AVChannelLayout) error {
	return newError(C.av_channel_layout_copy(dst, l.c))
}

func (l ChannelLayout) clone() (ChannelLayout, error) {
	var cl C.struct_AVChannelLayout
	err := l.copy(&cl)
	dst := newChannelLayoutFromC(&cl)
	return dst, err
}
