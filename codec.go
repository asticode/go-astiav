package astiav

//#include <libavcodec/avcodec.h>
//#include <libavutil/channel_layout.h>
import "C"
import (
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/codec.h#L202
type Codec struct {
	c *C.AVCodec
}

func newCodecFromC(c *C.AVCodec) *Codec {
	if c == nil {
		return nil
	}
	return &Codec{c: c}
}

func (c *Codec) Name() string {
	return C.GoString(c.c.name)
}

func (c *Codec) String() string {
	return c.Name()
}

func (c *Codec) ID() CodecID {
	return CodecID(c.c.id)
}

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

func (c *Codec) IsDecoder() bool {
	return int(C.av_codec_is_decoder(c.c)) != 0
}

func (c *Codec) IsEncoder() bool {
	return int(C.av_codec_is_encoder(c.c)) != 0
}

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

func FindDecoder(id CodecID) *Codec {
	return newCodecFromC(C.avcodec_find_decoder((C.enum_AVCodecID)(id)))
}

func FindDecoderByName(n string) *Codec {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newCodecFromC(C.avcodec_find_decoder_by_name(cn))
}

func FindEncoder(id CodecID) *Codec {
	return newCodecFromC(C.avcodec_find_encoder((C.enum_AVCodecID)(id)))
}

func FindEncoderByName(n string) *Codec {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	return newCodecFromC(C.avcodec_find_encoder_by_name(cn))
}

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
