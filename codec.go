package astiav

//#include <libavcodec/avcodec.h>
//#include <libavutil/channel_layout.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVCodec.html
type Codec struct {
	c *C.AVCodec
}

func newCodecFromC(c *C.AVCodec) *Codec {
	if c == nil {
		return nil
	}
	return &Codec{c: c}
}

// https://ffmpeg.org/doxygen/7.0/structAVCodec.html#ad3daa3e729850b573c139a83be8938ca
func (c *Codec) Name() string {
	return C.GoString(c.c.name)
}

func (c *Codec) String() string {
	return c.Name()
}

// https://ffmpeg.org/doxygen/7.0/structAVCodec.html#a01a53d07936f4c7ee280444793b6967b
func (c *Codec) ID() CodecID {
	return CodecID(c.c.id)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodec.html#a710e3bd3081124ef3364b0c520379dd8
func (c *Codec) ChannelLayouts() (o []ChannelLayout) {
	if c.c.ch_layouts == nil {
		return nil
	}
	size := unsafe.Sizeof(*c.c.ch_layouts)
	for i := 0; ; i++ {
		v, _ := newChannelLayoutFromC((*C.AVChannelLayout)(unsafe.Pointer(uintptr(unsafe.Pointer(c.c.ch_layouts)) + uintptr(i)*size))).clone()
		if !v.Valid() {
			break
		}
		o = append(o, v)
	}
	return
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga6dc18eef1afca3610644a52565cf8a31
func (c *Codec) IsDecoder() bool {
	return int(C.av_codec_is_decoder(c.c)) != 0
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga2b665824e4d9144f8d4f6c01e3e85aa3
func (c *Codec) IsEncoder() bool {
	return int(C.av_codec_is_encoder(c.c)) != 0
}

// https://ffmpeg.org/doxygen/7.0/structAVCodec.html#ac2b97bd3c19686025e1b7d577329c250
func (c *Codec) PixelFormats() (o []PixelFormat) {
	if c.c.pix_fmts == nil {
		return nil
	}
	size := unsafe.Sizeof(*c.c.pix_fmts)
	for i := 0; ; i++ {
		p := *(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(c.c.pix_fmts)) + uintptr(i)*size))
		if p == C.AV_PIX_FMT_NONE {
			break
		}
		o = append(o, PixelFormat(p))
	}
	return
}

// https://ffmpeg.org/doxygen/7.0/structAVCodec.html#aac19f4c45370f715412ad5c7b78daacf
func (c *Codec) SampleFormats() (o []SampleFormat) {
	if c.c.sample_fmts == nil {
		return nil
	}
	size := unsafe.Sizeof(*c.c.sample_fmts)
	for i := 0; ; i++ {
		p := *(*C.int)(unsafe.Pointer(uintptr(unsafe.Pointer(c.c.sample_fmts)) + uintptr(i)*size))
		if p == C.AV_SAMPLE_FMT_NONE {
			break
		}
		o = append(o, SampleFormat(p))
	}
	return
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga51e35d01da2b3833b3afa839212c58fa
func FindDecoder(id CodecID) *Codec {
	return newCodecFromC(C.avcodec_find_decoder((C.enum_AVCodecID)(id)))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#gad4e08a758f3560006145db074d16cb47
func FindDecoderByName(n string) *Codec {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newCodecFromC(C.avcodec_find_decoder_by_name(cn))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga68e4b5f31de5e5fc25d5781a1be22331
func FindEncoder(id CodecID) *Codec {
	return newCodecFromC(C.avcodec_find_encoder((C.enum_AVCodecID)(id)))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga9fa02c1eae54a2ec67beb789c2688d6e
func FindEncoderByName(n string) *Codec {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newCodecFromC(C.avcodec_find_encoder_by_name(cn))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga4f80582a2ea9c0e141de5d6f6152008f
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

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga7cd040fcc147340186deb0c54dc996b0
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
