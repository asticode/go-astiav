package astiav

//#cgo pkg-config: libavcodec
//#include <libavcodec/avcodec.h>
import "C"
import (
	"fmt"
	"unsafe"
)

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
	l, _ := newChannelLayoutFromC(&cp.c.ch_layout).clone()
	return l
}

func (cp *CodecParameters) SetChannelLayout(l ChannelLayout) {
	l.copy(&cp.c.ch_layout) //nolint: errcheck
}

func (cp *CodecParameters) Channels() int {
	return int(cp.c.channels)
}

func (cp *CodecParameters) SetChannels(c int) {
	cp.c.channels = C.int(c)
}

func (cp *CodecParameters) CodecID() CodecID {
	return CodecID(cp.c.codec_id)
}

func (cp *CodecParameters) SetCodecID(i CodecID) {
	cp.c.codec_id = uint32(i)
}

func (cp *CodecParameters) CodecTag() CodecTag {
	return CodecTag(cp.c.codec_tag)
}

func (cp *CodecParameters) SetCodecTag(t CodecTag) {
	cp.c.codec_tag = C.uint(t)
}

func (cp *CodecParameters) CodecType() MediaType {
	return MediaType(cp.c.codec_type)
}

func (cp *CodecParameters) SetCodecType(t MediaType) {
	cp.c.codec_type = int32(t)
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

func (cp *CodecParameters) SetColorRange(r ColorRange) {
	cp.c.color_range = C.enum_AVColorRange(r)
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

func (cp *CodecParameters) SetFrameSize(i int) {
	cp.c.frame_size = C.int(i)
}

func (cp *CodecParameters) Height() int {
	return int(cp.c.height)
}

func (cp *CodecParameters) SetHeight(h int) {
	cp.c.height = C.int(h)
}

func (cp *CodecParameters) Level() Level {
	return Level(cp.c.level)
}

func (cp *CodecParameters) MediaType() MediaType {
	return MediaType(cp.c.codec_type)
}

func (cp *CodecParameters) SetMediaType(t MediaType) {
	cp.c.codec_type = C.enum_AVMediaType(t)
}

func (cp *CodecParameters) PixelFormat() PixelFormat {
	return PixelFormat(cp.c.format)
}

func (cp *CodecParameters) SetPixelFormat(f PixelFormat) {
	cp.c.format = C.int(f)
}

func (cp *CodecParameters) Profile() Profile {
	return Profile(cp.c.profile)
}

func (cp *CodecParameters) SampleAspectRatio() Rational {
	return newRationalFromC(cp.c.sample_aspect_ratio)
}

func (cp *CodecParameters) SetSampleAspectRatio(r Rational) {
	cp.c.sample_aspect_ratio = r.c
}

func (cp *CodecParameters) ExtraData() []byte {
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		if cp.c.extradata == nil {
			*size = C.size_t(0)
			return nil
		}
		*size = C.size_t(cp.c.extradata_size)
		return cp.c.extradata
	})
}

func (cp *CodecParameters) SetExtraData(extraData []byte) error {
	if len(extraData) == 0 {
		return nil
	}

	if cp.c.extradata != nil {
		C.av_freep(unsafe.Pointer(&cp.c.extradata))
		cp.c.extradata_size = 0
	}

	extradataSize := len(extraData)
	if cp.c.extradata = (*C.uint8_t)(C.av_mallocz(C.size_t(extradataSize + C.AV_INPUT_BUFFER_PADDING_SIZE))); cp.c.extradata == nil {
		return fmt.Errorf("astiav: allocation is nil")
	}

	C.memcpy(unsafe.Pointer(cp.c.extradata), unsafe.Pointer(&extraData[0]), C.size_t(extradataSize))
	cp.c.extradata_size = C.int(extradataSize)
	return nil
}

func (cp *CodecParameters) SampleFormat() SampleFormat {
	return SampleFormat(cp.c.format)
}

func (cp *CodecParameters) SetSampleFormat(f SampleFormat) {
	cp.c.format = C.int(f)
}

func (cp *CodecParameters) SampleRate() int {
	return int(cp.c.sample_rate)
}

func (cp *CodecParameters) SetSampleRate(r int) {
	cp.c.sample_rate = C.int(r)
}

func (cp *CodecParameters) Width() int {
	return int(cp.c.width)
}

func (cp *CodecParameters) SetWidth(w int) {
	cp.c.width = C.int(w)
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
