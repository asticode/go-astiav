package astiav

//#cgo pkg-config: libavutil
//#include <libavutil/channel_layout.h>
import "C"

type ChannelLayout uint64

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/channel_layout.h#L90
const (
	ChannelLayoutMono            = ChannelLayout(C.AV_CH_LAYOUT_MONO)
	ChannelLayoutStereo          = ChannelLayout(C.AV_CH_LAYOUT_STEREO)
	ChannelLayout2Point1         = ChannelLayout(C.AV_CH_LAYOUT_2POINT1)
	ChannelLayout21              = ChannelLayout(C.AV_CH_LAYOUT_2_1)
	ChannelLayoutSurround        = ChannelLayout(C.AV_CH_LAYOUT_SURROUND)
	ChannelLayout3Point1         = ChannelLayout(C.AV_CH_LAYOUT_3POINT1)
	ChannelLayout4Point0         = ChannelLayout(C.AV_CH_LAYOUT_4POINT0)
	ChannelLayout4Point1         = ChannelLayout(C.AV_CH_LAYOUT_4POINT1)
	ChannelLayout22              = ChannelLayout(C.AV_CH_LAYOUT_2_2)
	ChannelLayoutQuad            = ChannelLayout(C.AV_CH_LAYOUT_QUAD)
	ChannelLayout5Point0         = ChannelLayout(C.AV_CH_LAYOUT_5POINT0)
	ChannelLayout5Point1         = ChannelLayout(C.AV_CH_LAYOUT_5POINT1)
	ChannelLayout5Point0Back     = ChannelLayout(C.AV_CH_LAYOUT_5POINT0_BACK)
	ChannelLayout5Point1Back     = ChannelLayout(C.AV_CH_LAYOUT_5POINT1_BACK)
	ChannelLayout6Point0         = ChannelLayout(C.AV_CH_LAYOUT_6POINT0)
	ChannelLayout6Point0Front    = ChannelLayout(C.AV_CH_LAYOUT_6POINT0_FRONT)
	ChannelLayoutHexagonal       = ChannelLayout(C.AV_CH_LAYOUT_HEXAGONAL)
	ChannelLayout6Point1         = ChannelLayout(C.AV_CH_LAYOUT_6POINT1)
	ChannelLayout6Point1Back     = ChannelLayout(C.AV_CH_LAYOUT_6POINT1_BACK)
	ChannelLayout6Point1Front    = ChannelLayout(C.AV_CH_LAYOUT_6POINT1_FRONT)
	ChannelLayout7Point0         = ChannelLayout(C.AV_CH_LAYOUT_7POINT0)
	ChannelLayout7Point0Front    = ChannelLayout(C.AV_CH_LAYOUT_7POINT0_FRONT)
	ChannelLayout7Point1         = ChannelLayout(C.AV_CH_LAYOUT_7POINT1)
	ChannelLayout7Point1Wide     = ChannelLayout(C.AV_CH_LAYOUT_7POINT1_WIDE)
	ChannelLayout7Point1WideBack = ChannelLayout(C.AV_CH_LAYOUT_7POINT1_WIDE_BACK)
	ChannelLayoutOctagonal       = ChannelLayout(C.AV_CH_LAYOUT_OCTAGONAL)
	ChannelLayoutHexadecagonal   = ChannelLayout(C.AV_CH_LAYOUT_HEXADECAGONAL)
	ChannelLayoutStereoDownmix   = ChannelLayout(C.AV_CH_LAYOUT_STEREO_DOWNMIX)
)

func (l ChannelLayout) NbChannels() int {
	return int(C.av_get_channel_layout_nb_channels(C.uint64_t(l)))
}

func (l ChannelLayout) String() string {
	return l.StringWithNbChannels(l.NbChannels())
}

func (l ChannelLayout) StringWithNbChannels(nbChannels int) string {
	s, _ := stringFromC(255, func(buf *C.char, size C.size_t) error {
		C.av_get_channel_layout_string(buf, C.int(size), C.int(nbChannels), C.uint64_t(l))
		return nil
	})
	return s
}
