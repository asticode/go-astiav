package astiav

//#cgo pkg-config: libavcodec
//#include <libavcodec/avcodec.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/codec_par.h#L52
type CodecParameters struct {
	c *C.struct_AVCodecParameters
}

func AllocCodecParameters() *CodecParameters {
	return newCodecParametersFromC(C.avcodec_parameters_alloc())
}

func newCodecParametersFromC(c *C.struct_AVCodecParameters) *CodecParameters {
	if c == nil {
		return nil
	}
	return &CodecParameters{c: c}
}

func (cp *CodecParameters) Free() {
	C.avcodec_parameters_free(&cp.c)
}

func (cp *CodecParameters) BitRate() int64 {
	return int64(cp.c.bit_rate)
}

func (cp *CodecParameters) ChannelLayout() ChannelLayout {
	return ChannelLayout(cp.c.channel_layout)
}

func (cp *CodecParameters) Channels() int {
	return int(cp.c.channels)
}

func (cp *CodecParameters) CodecID() CodecID {
	return CodecID(cp.c.codec_id)
}

func (cp *CodecParameters) CodecTag() CodecTag {
	return CodecTag(cp.c.codec_tag)
}

func (cp *CodecParameters) SetCodecTag(t CodecTag) {
	cp.c.codec_tag = C.uint(t)
}

func (cp *CodecParameters) ChromaLocation() ChromaLocation {
	return ChromaLocation(cp.c.chroma_location)
}

func (cp *CodecParameters) ColorPrimaries() ColorPrimaries {
	return ColorPrimaries(cp.c.color_primaries)
}

func (cp *CodecParameters) ColorRange() ColorRange {
	return ColorRange(cp.c.color_range)
}

func (cp *CodecParameters) ColorSpace() ColorSpace {
	return ColorSpace(cp.c.color_space)
}

func (cp *CodecParameters) ColorTransferCharacteristic() ColorTransferCharacteristic {
	return ColorTransferCharacteristic(cp.c.color_trc)
}

func (cp *CodecParameters) FrameSize() int {
	return int(cp.c.frame_size)
}

func (cp *CodecParameters) Height() int {
	return int(cp.c.height)
}

func (cp *CodecParameters) Level() Level {
	return Level(cp.c.level)
}

func (cp *CodecParameters) MediaType() MediaType {
	return MediaType(cp.c.codec_type)
}

func (cp *CodecParameters) PixelFormat() PixelFormat {
	return PixelFormat(cp.c.format)
}

func (cp *CodecParameters) Profile() Profile {
	return Profile(cp.c.profile)
}

func (cp *CodecParameters) SampleAspectRatio() Rational {
	return newRationalFromC(cp.c.sample_aspect_ratio)
}

func (cp *CodecParameters) SampleFormat() SampleFormat {
	return SampleFormat(cp.c.format)
}

func (cp *CodecParameters) SampleRate() int {
	return int(cp.c.sample_rate)
}

func (cp *CodecParameters) Width() int {
	return int(cp.c.width)
}

func (cp *CodecParameters) FromCodecContext(cc *CodecContext) error {
	return newError(C.avcodec_parameters_from_context(cp.c, cc.c))
}

func (cp *CodecParameters) ToCodecContext(cc *CodecContext) error {
	return newError(C.avcodec_parameters_to_context(cc.c, cp.c))
}

func (cp *CodecParameters) Copy(dst *CodecParameters) error {
	return newError(C.avcodec_parameters_copy(dst.c, cp.c))
}
