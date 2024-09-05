package astiav

//#include "channel_layout.h"
import "C"
import "unsafe"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/channel_layout.h#L90
var (
	ChannelLayoutMono              = newChannelLayoutFromC(C.astiavChannelLayoutMono)
	ChannelLayoutStereo            = newChannelLayoutFromC(C.astiavChannelLayoutStereo)
	ChannelLayout2Point1           = newChannelLayoutFromC(C.astiavChannelLayout2Point1)
	ChannelLayout21                = newChannelLayoutFromC(C.astiavChannelLayout21)
	ChannelLayoutSurround          = newChannelLayoutFromC(C.astiavChannelLayoutSurround)
	ChannelLayout3Point1           = newChannelLayoutFromC(C.astiavChannelLayout3Point1)
	ChannelLayout4Point0           = newChannelLayoutFromC(C.astiavChannelLayout4Point0)
	ChannelLayout4Point1           = newChannelLayoutFromC(C.astiavChannelLayout4Point1)
	ChannelLayout22                = newChannelLayoutFromC(C.astiavChannelLayout22)
	ChannelLayoutQuad              = newChannelLayoutFromC(C.astiavChannelLayoutQuad)
	ChannelLayout5Point0           = newChannelLayoutFromC(C.astiavChannelLayout5Point0)
	ChannelLayout5Point1           = newChannelLayoutFromC(C.astiavChannelLayout5Point1)
	ChannelLayout5Point0Back       = newChannelLayoutFromC(C.astiavChannelLayout5Point0Back)
	ChannelLayout5Point1Back       = newChannelLayoutFromC(C.astiavChannelLayout5Point1Back)
	ChannelLayout6Point0           = newChannelLayoutFromC(C.astiavChannelLayout6Point0)
	ChannelLayout6Point0Front      = newChannelLayoutFromC(C.astiavChannelLayout6Point0Front)
	ChannelLayoutHexagonal         = newChannelLayoutFromC(C.astiavChannelLayoutHexagonal)
	ChannelLayout3Point1Point2     = newChannelLayoutFromC(C.astiavChannelLayout3Point1Point2)
	ChannelLayout6Point1           = newChannelLayoutFromC(C.astiavChannelLayout6Point1)
	ChannelLayout6Point1Back       = newChannelLayoutFromC(C.astiavChannelLayout6Point1Back)
	ChannelLayout6Point1Front      = newChannelLayoutFromC(C.astiavChannelLayout6Point1Front)
	ChannelLayout7Point0           = newChannelLayoutFromC(C.astiavChannelLayout7Point0)
	ChannelLayout7Point0Front      = newChannelLayoutFromC(C.astiavChannelLayout7Point0Front)
	ChannelLayout7Point1           = newChannelLayoutFromC(C.astiavChannelLayout7Point1)
	ChannelLayout7Point1Wide       = newChannelLayoutFromC(C.astiavChannelLayout7Point1Wide)
	ChannelLayout7Point1WideBack   = newChannelLayoutFromC(C.astiavChannelLayout7Point1WideBack)
	ChannelLayout5Point1Point2Back = newChannelLayoutFromC(C.astiavChannelLayout5Point1Point2Back)
	ChannelLayoutOctagonal         = newChannelLayoutFromC(C.astiavChannelLayoutOctagonal)
	ChannelLayoutCube              = newChannelLayoutFromC(C.astiavChannelLayoutCube)
	ChannelLayout5Point1Point4Back = newChannelLayoutFromC(C.astiavChannelLayout5Point1Point4Back)
	ChannelLayout7Point1Point2     = newChannelLayoutFromC(C.astiavChannelLayout7Point1Point2)
	ChannelLayout7Point1Point4Back = newChannelLayoutFromC(C.astiavChannelLayout7Point1Point4Back)
	ChannelLayoutHexadecagonal     = newChannelLayoutFromC(C.astiavChannelLayoutHexadecagonal)
	ChannelLayoutStereoDownmix     = newChannelLayoutFromC(C.astiavChannelLayoutStereoDownmix)
	ChannelLayout22Point2          = newChannelLayoutFromC(C.astiavChannelLayout22Point2)
	ChannelLayout7Point1TopBack    = newChannelLayoutFromC(C.astiavChannelLayout7Point1TopBack)
)

type ChannelLayout struct {
	c *C.AVChannelLayout
}

func newChannelLayoutFromC(c *C.AVChannelLayout) ChannelLayout {
	return ChannelLayout{c: c}
}

func (l ChannelLayout) Channels() int {
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
	if err := newError(ret); err != nil {
		return 0, err
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
	if err := newError(ret); err != nil {
		return false, err
	}
	return ret == 0, nil
}

func (l ChannelLayout) Equal(l2 ChannelLayout) bool {
	v, _ := l.Compare(l2)
	return v
}

func (l ChannelLayout) copy(dst *C.AVChannelLayout) error {
	return newError(C.av_channel_layout_copy(dst, l.c))
}

func (l ChannelLayout) clone() (ChannelLayout, error) {
	var cl C.AVChannelLayout
	err := l.copy(&cl)
	dst := newChannelLayoutFromC(&cl)
	return dst, err
}
