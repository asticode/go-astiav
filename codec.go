package astiav

//#include <libavcodec/avcodec.h>
//#include <libavutil/channel_layout.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.1/structAVCodec.html
type Codec struct {
	c *C.AVCodec
}

func newCodecFromC(c *C.AVCodec) *Codec {
	if c == nil {
		return nil
	}
	return &Codec{c: c}
}

// https://ffmpeg.org/doxygen/7.1/structAVCodec.html#ad3daa3e729850b573c139a83be8938ca
func (c *Codec) Name() string {
	return C.GoString(c.c.name)
}

func (c *Codec) String() string {
	return c.Name()
}

// https://ffmpeg.org/doxygen/7.1/structAVCodec.html#a01a53d07936f4c7ee280444793b6967b
func (c *Codec) ID() CodecID {
	return CodecID(c.c.id)
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__decoding.html#gadd58e6b0bbca99fdbc547efbaa6b0ef1
func (c *Codec) ChannelLayouts() (o []ChannelLayout) {
	const codecConfig = C.enum_AVCodecConfig(CodecConfigChannelLayout)
	var outConfigs unsafe.Pointer
	var outNumConfigs C.int
	ret := C.avcodec_get_supported_config(nil, c.c, codecConfig, 0, &outConfigs, &outNumConfigs)
	if ret >= 0 && outConfigs != nil {
		numConfigs := uintptr(outNumConfigs)
		for i := uintptr(0); i < numConfigs; i++ {
			size := unsafe.Sizeof(C.AVChannelLayout{})
			v, _ := newChannelLayoutFromC((*C.AVChannelLayout)(unsafe.Pointer(uintptr(outConfigs) + i*size))).clone()
			o = append(o, v)
		}
	}
	return
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__core.html#ga6dc18eef1afca3610644a52565cf8a31
func (c *Codec) IsDecoder() bool {
	return int(C.av_codec_is_decoder(c.c)) != 0
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__core.html#ga2b665824e4d9144f8d4f6c01e3e85aa3
func (c *Codec) IsEncoder() bool {
	return int(C.av_codec_is_encoder(c.c)) != 0
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__decoding.html#gadd58e6b0bbca99fdbc547efbaa6b0ef1
func (c *Codec) ColorRanges() (o []ColorRange) {
	const codecConfig = C.enum_AVCodecConfig(CodecConfigColorRange)
	var outConfigs unsafe.Pointer
	var outNumConfigs C.int
	ret := C.avcodec_get_supported_config(nil, c.c, codecConfig, 0, &outConfigs, &outNumConfigs)
	if ret >= 0 && outConfigs != nil {
		numConfigs := uintptr(outNumConfigs)
		for i := uintptr(0); i < numConfigs; i++ {
			size := unsafe.Sizeof(C.enum_AVColorRange(0))
			o = append(o, ColorRange(*(*C.enum_AVColorRange)(unsafe.Pointer(uintptr(outConfigs) + i*size))))
		}
	}
	return
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__decoding.html#gadd58e6b0bbca99fdbc547efbaa6b0ef1
func (c *Codec) ColorSpaces() (o []ColorSpace) {
	const codecConfig = C.enum_AVCodecConfig(CodecConfigColorSpace)
	var outConfigs unsafe.Pointer
	var outNumConfigs C.int
	ret := C.avcodec_get_supported_config(nil, c.c, codecConfig, 0, &outConfigs, &outNumConfigs)
	if ret >= 0 && outConfigs != nil {
		numConfigs := uintptr(outNumConfigs)
		for i := uintptr(0); i < numConfigs; i++ {
			size := unsafe.Sizeof(C.enum_AVColorSpace(0))
			o = append(o, ColorSpace(*(*C.enum_AVColorSpace)(unsafe.Pointer(uintptr(outConfigs) + i*size))))
		}
	}
	return
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__decoding.html#gadd58e6b0bbca99fdbc547efbaa6b0ef1
func (c *Codec) PixelFormats() (o []PixelFormat) {
	const codecConfig = C.enum_AVCodecConfig(CodecConfigPixFormat)
	var outConfigs unsafe.Pointer
	var outNumConfigs C.int
	ret := C.avcodec_get_supported_config(nil, c.c, codecConfig, 0, &outConfigs, &outNumConfigs)
	if ret >= 0 && outConfigs != nil {
		numConfigs := uintptr(outNumConfigs)
		for i := uintptr(0); i < numConfigs; i++ {
			size := unsafe.Sizeof(C.enum_AVPixelFormat(0))
			o = append(o, PixelFormat(*(*C.enum_AVPixelFormat)(unsafe.Pointer(uintptr(outConfigs) + i*size))))
		}
	}
	return
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__decoding.html#gadd58e6b0bbca99fdbc547efbaa6b0ef1
func (c *Codec) SampleFormats() (o []SampleFormat) {
	const codecConfig = C.enum_AVCodecConfig(CodecConfigSampleFormat)
	var outConfigs unsafe.Pointer
	var outNumConfigs C.int
	ret := C.avcodec_get_supported_config(nil, c.c, codecConfig, 0, &outConfigs, &outNumConfigs)
	if ret >= 0 && outConfigs != nil {
		numConfigs := uintptr(outNumConfigs)
		for i := uintptr(0); i < numConfigs; i++ {
			size := unsafe.Sizeof(C.enum_AVSampleFormat(0))
			o = append(o, SampleFormat(*(*C.enum_AVSampleFormat)(unsafe.Pointer(uintptr(outConfigs) + i*size))))
		}
	}
	return
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__decoding.html#gadd58e6b0bbca99fdbc547efbaa6b0ef1
func (c *Codec) SupportedFramerates() (o []Rational) {
	const codecConfig = C.enum_AVCodecConfig(CodecConfigFrameRate)
	var outConfigs unsafe.Pointer
	var outNumConfigs C.int
	ret := C.avcodec_get_supported_config(nil, c.c, codecConfig, 0, &outConfigs, &outNumConfigs)
	if ret >= 0 && outConfigs != nil {
		numConfigs := uintptr(outNumConfigs)
		for i := uintptr(0); i < numConfigs; i++ {
			size := unsafe.Sizeof(C.AVRational{})
			o = append(o, newRationalFromC(*(*C.AVRational)(unsafe.Pointer(uintptr(outConfigs) + i*size))))
		}
	}
	return
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__decoding.html#gadd58e6b0bbca99fdbc547efbaa6b0ef1
func (c *Codec) SupportedSamplerates() (o []int) {
	const codecConfig = C.enum_AVCodecConfig(CodecConfigSampleRate)
	var outConfigs unsafe.Pointer
	var outNumConfigs C.int
	ret := C.avcodec_get_supported_config(nil, c.c, codecConfig, 0, &outConfigs, &outNumConfigs)
	if ret >= 0 && outConfigs != nil {
		numConfigs := uintptr(outNumConfigs)
		for i := uintptr(0); i < numConfigs; i++ {
			size := unsafe.Sizeof(C.int(0))
			o = append(o, int(*(*C.int)(unsafe.Pointer(uintptr(outConfigs) + i*size))))
		}
	}
	return
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__core.html#ga51e35d01da2b3833b3afa839212c58fa
func FindDecoder(id CodecID) *Codec {
	return newCodecFromC(C.avcodec_find_decoder((C.enum_AVCodecID)(id)))
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__core.html#gad4e08a758f3560006145db074d16cb47
func FindDecoderByName(n string) *Codec {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newCodecFromC(C.avcodec_find_decoder_by_name(cn))
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__core.html#ga68e4b5f31de5e5fc25d5781a1be22331
func FindEncoder(id CodecID) *Codec {
	return newCodecFromC(C.avcodec_find_encoder((C.enum_AVCodecID)(id)))
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__core.html#ga9fa02c1eae54a2ec67beb789c2688d6e
func FindEncoderByName(n string) *Codec {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newCodecFromC(C.avcodec_find_encoder_by_name(cn))
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__core.html#ga4f80582a2ea9c0e141de5d6f6152008f
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

// https://ffmpeg.org/doxygen/7.1/group__lavc__core.html#ga7cd040fcc147340186deb0c54dc996b0
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
