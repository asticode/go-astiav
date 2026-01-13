package astiav

//#include <libavcodec/avcodec.h>
//#include <libavutil/channel_layout.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/8.0/structAVCodec.html
type Codec struct {
	c *C.AVCodec
}

func newCodecFromC(c *C.AVCodec) *Codec {
	if c == nil {
		return nil
	}
	return &Codec{c: c}
}

// https://ffmpeg.org/doxygen/8.0/structAVCodec.html#ad3daa3e729850b573c139a83be8938ca
func (c *Codec) Name() string {
	return C.GoString(c.c.name)
}

func (c *Codec) String() string {
	return c.Name()
}

// https://ffmpeg.org/doxygen/8.0/structAVCodec.html#a01a53d07936f4c7ee280444793b6967b
func (c *Codec) ID() CodecID {
	return CodecID(c.c.id)
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__core.html#ga6dc18eef1afca3610644a52565cf8a31
func (c *Codec) IsDecoder() bool {
	return int(C.av_codec_is_decoder(c.c)) != 0
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__core.html#ga2b665824e4d9144f8d4f6c01e3e85aa3
func (c *Codec) IsEncoder() bool {
	return int(C.av_codec_is_encoder(c.c)) != 0
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__decoding.html#gadd58e6b0bbca99fdbc547efbaa6b0ef1
func (c *Codec) supportedConfig(config C.enum_AVCodecConfig, fn func(ptr unsafe.Pointer), size C.size_t) {
	var outConfigs unsafe.Pointer
	var outNumConfigs C.int
	ret := C.avcodec_get_supported_config(nil, c.c, config, 0, &outConfigs, &outNumConfigs)
	if ret >= 0 && outConfigs != nil {
		numConfigs := uintptr(outNumConfigs)
		for i := uintptr(0); i < numConfigs; i++ {
			fn(unsafe.Pointer(uintptr(outConfigs) + i*uintptr(size)))
		}
	}
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__decoding.html#ga9abe043ed40f3c270dff77235fcfcd0b
func (c *Codec) SupportedChannelLayouts() (o []ChannelLayout) {
	c.supportedConfig(C.AV_CODEC_CONFIG_CHANNEL_LAYOUT, func(ptr unsafe.Pointer) {
		v, _ := newChannelLayoutFromC((*C.AVChannelLayout)(ptr)).clone()
		o = append(o, v)
	}, C.sizeof_AVChannelLayout)
	return
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__decoding.html#ga9abe043ed40f3c270dff77235fcfcd0b
func (c *Codec) SupportedColorRanges() (o []ColorRange) {
	c.supportedConfig(C.AV_CODEC_CONFIG_COLOR_RANGE, func(ptr unsafe.Pointer) {
		o = append(o, ColorRange(*(*C.enum_AVColorRange)(ptr)))
	}, C.sizeof_enum_AVColorRange)
	return
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__decoding.html#ga9abe043ed40f3c270dff77235fcfcd0b
func (c *Codec) SupportedColorSpaces() (o []ColorSpace) {
	c.supportedConfig(C.AV_CODEC_CONFIG_COLOR_SPACE, func(ptr unsafe.Pointer) {
		o = append(o, ColorSpace(*(*C.enum_AVColorSpace)(ptr)))
	}, C.sizeof_enum_AVColorSpace)
	return
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__decoding.html#ga9abe043ed40f3c270dff77235fcfcd0b
func (c *Codec) SupportedPixelFormats() (o []PixelFormat) {
	c.supportedConfig(C.AV_CODEC_CONFIG_PIX_FORMAT, func(ptr unsafe.Pointer) {
		o = append(o, PixelFormat(*(*C.enum_AVPixelFormat)(ptr)))
	}, C.sizeof_enum_AVPixelFormat)
	return
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__decoding.html#ga9abe043ed40f3c270dff77235fcfcd0b
func (c *Codec) SupportedSampleFormats() (o []SampleFormat) {
	c.supportedConfig(C.AV_CODEC_CONFIG_SAMPLE_FORMAT, func(ptr unsafe.Pointer) {
		o = append(o, SampleFormat(*(*C.enum_AVSampleFormat)(ptr)))
	}, C.sizeof_enum_AVSampleFormat)
	return
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__decoding.html#ga9abe043ed40f3c270dff77235fcfcd0b
func (c *Codec) SupportedFrameRates() (o []Rational) {
	c.supportedConfig(C.AV_CODEC_CONFIG_FRAME_RATE, func(ptr unsafe.Pointer) {
		o = append(o, newRationalFromC(*(*C.AVRational)(ptr)))
	}, C.sizeof_AVRational)
	return
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__decoding.html#ga9abe043ed40f3c270dff77235fcfcd0b
func (c *Codec) SupportedSampleRates() (o []int) {
	c.supportedConfig(C.AV_CODEC_CONFIG_SAMPLE_RATE, func(ptr unsafe.Pointer) {
		o = append(o, int(*(*C.int)(ptr)))
	}, C.sizeof_int)
	return
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__core.html#ga51e35d01da2b3833b3afa839212c58fa
func FindDecoder(id CodecID) *Codec {
	return newCodecFromC(C.avcodec_find_decoder((C.enum_AVCodecID)(id)))
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__core.html#gad4e08a758f3560006145db074d16cb47
func FindDecoderByName(n string) *Codec {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newCodecFromC(C.avcodec_find_decoder_by_name(cn))
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__core.html#ga68e4b5f31de5e5fc25d5781a1be22331
func FindEncoder(id CodecID) *Codec {
	return newCodecFromC(C.avcodec_find_encoder((C.enum_AVCodecID)(id)))
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__core.html#ga9fa02c1eae54a2ec67beb789c2688d6e
func FindEncoderByName(n string) *Codec {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newCodecFromC(C.avcodec_find_encoder_by_name(cn))
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__core.html#ga4f80582a2ea9c0e141de5d6f6152008f
func (c *Codec) HardwareConfigs() (configs []CodecHardwareConfig) {
	var i int
	for {
		config := C.avcodec_get_hw_config(c.c, C.int(i))
		if config == nil {
			break
		}
		configs = append(configs, newCodecHardwareConfigFromC(config))
		i++
	}
	return
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__core.html#ga7cd040fcc147340186deb0c54dc996b0
func Codecs() (cs []*Codec) {
	var opq *C.void = nil
	for {
		c := C.av_codec_iterate((*unsafe.Pointer)(unsafe.Pointer(&opq)))
		if c == nil {
			break
		}
		cs = append(cs, newCodecFromC(c))
	}
	return
}
