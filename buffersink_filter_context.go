package astiav

//#include <libavfilter/buffersink.h>
import "C"

type BuffersinkFilterContext struct {
	fc *FilterContext
}

func newBuffersinkFilterContext(fc *FilterContext) *BuffersinkFilterContext {
	return &BuffersinkFilterContext{fc: fc}
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gaad918036937648701c09f9612f42706e
func (bfc *BuffersinkFilterContext) ChannelLayout() ChannelLayout {
	var cl C.AVChannelLayout
	C.av_buffersink_get_ch_layout(bfc.fc.c, &cl)
	return newChannelLayoutFromC(&cl)
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gab80976e506ab88d23d94bb6d7a4051bd
func (bfc *BuffersinkFilterContext) ColorRange() ColorRange {
	return ColorRange(C.av_buffersink_get_color_range(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gaad817cdcf5493c385126e8e17c5717f2
func (bfc *BuffersinkFilterContext) ColorSpace() ColorSpace {
	return ColorSpace(C.av_buffersink_get_colorspace(bfc.fc.c))
}

func (bfc *BuffersinkFilterContext) FilterContext() *FilterContext {
	return bfc.fc
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga55614fd28de2fa05b04f427390061d5b
func (bfc *BuffersinkFilterContext) FrameRate() Rational {
	return newRationalFromC(C.av_buffersink_get_frame_rate(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink.html#ga71ae9c529c8da51681e12faa37d1a395
func (bfc *BuffersinkFilterContext) GetFrame(f *Frame, fs BuffersinkFlags) error {
	var cf *C.AVFrame
	if f != nil {
		cf = f.c
	}
	return newError(C.av_buffersink_get_frame_flags(bfc.fc.c, cf, C.int(fs)))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga955ecf3680e71e10429d7500343be25c
func (bfc *BuffersinkFilterContext) Height() int {
	return int(C.av_buffersink_get_h(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga1eb8bbf583ffb7cc29aaa1944b1e699c
func (bfc *BuffersinkFilterContext) MediaType() MediaType {
	return MediaType(C.av_buffersink_get_type(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga402ddbef6f7347869725696846ac81eb
func (bfc *BuffersinkFilterContext) PixelFormat() PixelFormat {
	return PixelFormat(C.av_buffersink_get_format(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gaa38ee33e1c7f6f7cb190bd2330e5f848
func (bfc *BuffersinkFilterContext) SampleAspectRatio() Rational {
	return newRationalFromC(C.av_buffersink_get_sample_aspect_ratio(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga402ddbef6f7347869725696846ac81eb
func (bfc *BuffersinkFilterContext) SampleFormat() SampleFormat {
	return SampleFormat(C.av_buffersink_get_format(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#ga2af714e82f48759551acdbc4488ded4a
func (bfc *BuffersinkFilterContext) SampleRate() int {
	return int(C.av_buffersink_get_sample_rate(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gabc82f65ec7f4fa47c5216260639258a1
func (bfc *BuffersinkFilterContext) TimeBase() Rational {
	return newRationalFromC(C.av_buffersink_get_time_base(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink__accessors.html#gac8c86515d2ef56090395dfd74854c835
func (bfc *BuffersinkFilterContext) Width() int {
	return int(C.av_buffersink_get_w(bfc.fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersink.html#ga359d7d1e42c27ca14c07559d4e9adba7
func (bfc *BuffersinkFilterContext) SetFrameSize(size int) {
	C.av_buffersink_set_frame_size(bfc.fc.c, C.uint(size))
}
