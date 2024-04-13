package astiav

//#cgo pkg-config: libavcodec libavutil
//#include <libavcodec/avcodec.h>
//#include <libavutil/frame.h>
/*
extern enum AVPixelFormat goAstiavCodecContextGetFormat(AVCodecContext *ctx, enum AVPixelFormat *pix_fmts, int pix_fmts_size);

static inline enum AVPixelFormat astiavCodecContextGetFormat(AVCodecContext *ctx, const enum AVPixelFormat *pix_fmts)
{
	int pix_fmts_size = 0;
	while (*pix_fmts != AV_PIX_FMT_NONE) {
		pix_fmts_size++;
		pix_fmts++;
	}
	pix_fmts -= pix_fmts_size;
	return goAstiavCodecContextGetFormat(ctx, (enum AVPixelFormat*)(pix_fmts), pix_fmts_size);
}
static inline void astiavSetCodecContextGetFormat(AVCodecContext *ctx)
{
	ctx->get_format = astiavCodecContextGetFormat;
}
static inline void astiavResetCodecContextGetFormat(AVCodecContext *ctx)
{
	ctx->get_format = NULL;
}

*/
import "C"
import (
	"sync"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/avcodec.h#L383
type CodecContext struct {
	c *C.struct_AVCodecContext
	// We need to store this to unref it properly
	hdc *HardwareDeviceContext
}

func newCodecContextFromC(c *C.struct_AVCodecContext) *CodecContext {
	if c == nil {
		return nil
	}
	cc := &CodecContext{c: c}
	classers.set(cc)
	return cc
}

var _ Classer = (*CodecContext)(nil)

func AllocCodecContext(c *Codec) *CodecContext {
	var cc *C.struct_AVCodec
	if c != nil {
		cc = c.c
	}
	return newCodecContextFromC(C.avcodec_alloc_context3(cc))
}

func (cc *CodecContext) Free() {
	if cc.hdc != nil {
		C.av_buffer_unref(&cc.hdc.c)
		cc.hdc = nil
	}
	classers.del(cc)
	C.avcodec_free_context(&cc.c)
}

func (cc *CodecContext) String() string {
	s, _ := stringFromC(255, func(buf *C.char, size C.size_t) error {
		C.avcodec_string(buf, C.int(size), cc.c, C.int(0))
		return nil
	})
	return s
}

func (cc *CodecContext) BitRate() int64 {
	return int64(cc.c.bit_rate)
}

func (cc *CodecContext) SetBitRate(bitRate int64) {
	cc.c.bit_rate = C.int64_t(bitRate)
}

func (cc *CodecContext) Channels() int {
	return int(cc.c.channels)
}

func (cc *CodecContext) SetChannels(channels int) {
	cc.c.channels = C.int(channels)
}

func (cc *CodecContext) ChannelLayout() ChannelLayout {
	l, _ := newChannelLayoutFromC(&cc.c.ch_layout).clone()
	return l
}

func (cc *CodecContext) SetChannelLayout(l ChannelLayout) {
	l.copy(&cc.c.ch_layout) //nolint: errcheck
}

func (cc *CodecContext) ChromaLocation() ChromaLocation {
	return ChromaLocation(cc.c.chroma_sample_location)
}

func (cc *CodecContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(cc.c))
}

func (cc *CodecContext) CodecID() CodecID {
	return CodecID(cc.c.codec_id)
}

func (cc *CodecContext) ColorPrimaries() ColorPrimaries {
	return ColorPrimaries(cc.c.color_primaries)
}

func (cc *CodecContext) ColorRange() ColorRange {
	return ColorRange(cc.c.color_range)
}

func (cc *CodecContext) ColorSpace() ColorSpace {
	return ColorSpace(cc.c.colorspace)
}

func (cc *CodecContext) ColorTransferCharacteristic() ColorTransferCharacteristic {
	return ColorTransferCharacteristic(cc.c.color_trc)
}

func (cc *CodecContext) ExtraData() []byte {
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		*size = C.size_t(cc.c.extradata_size)
		return cc.c.extradata
	})
}

func (cc *CodecContext) SetExtraData(b []byte) error {
	return setBytesWithIntSizeInC(b, &cc.c.extradata, &cc.c.extradata_size)
}

func (cc *CodecContext) Flags() CodecContextFlags {
	return CodecContextFlags(cc.c.flags)
}

func (cc *CodecContext) SetFlags(fs CodecContextFlags) {
	cc.c.flags = C.int(fs)
}

func (cc *CodecContext) Flags2() CodecContextFlags2 {
	return CodecContextFlags2(cc.c.flags2)
}

func (cc *CodecContext) SetFlags2(fs CodecContextFlags2) {
	cc.c.flags2 = C.int(fs)
}

func (cc *CodecContext) Framerate() Rational {
	return newRationalFromC(cc.c.framerate)
}

func (cc *CodecContext) SetFramerate(f Rational) {
	cc.c.framerate = f.c
}

func (cc *CodecContext) FrameSize() int {
	return int(cc.c.frame_size)
}

func (cc *CodecContext) GopSize() int {
	return int(cc.c.gop_size)
}

func (cc *CodecContext) SetGopSize(gopSize int) {
	cc.c.gop_size = C.int(gopSize)
}

func (cc *CodecContext) Height() int {
	return int(cc.c.height)
}

func (cc *CodecContext) SetHeight(height int) {
	cc.c.height = C.int(height)
}

func (cc *CodecContext) Level() Level {
	return Level(cc.c.level)
}

func (cc *CodecContext) SetLevel(l Level) {
	cc.c.level = C.int(l)
}

func (cc *CodecContext) MediaType() MediaType {
	return MediaType(cc.c.codec_type)
}

func (cc *CodecContext) PixelFormat() PixelFormat {
	return PixelFormat(cc.c.pix_fmt)
}

func (cc *CodecContext) SetPixelFormat(pixFmt PixelFormat) {
	cc.c.pix_fmt = C.enum_AVPixelFormat(pixFmt)
}

func (cc *CodecContext) Profile() Profile {
	return Profile(cc.c.profile)
}

func (cc *CodecContext) SetProfile(p Profile) {
	cc.c.profile = C.int(p)
}

func (cc *CodecContext) Qmin() int {
	return int(cc.c.qmin)
}

func (cc *CodecContext) SetQmin(qmin int) {
	cc.c.qmin = C.int(qmin)
}

func (cc *CodecContext) SampleAspectRatio() Rational {
	return newRationalFromC(cc.c.sample_aspect_ratio)
}

func (cc *CodecContext) SetSampleAspectRatio(r Rational) {
	cc.c.sample_aspect_ratio = r.c
}

func (cc *CodecContext) SampleFormat() SampleFormat {
	return SampleFormat(cc.c.sample_fmt)
}

func (cc *CodecContext) SetSampleFormat(f SampleFormat) {
	cc.c.sample_fmt = C.enum_AVSampleFormat(f)
}

func (cc *CodecContext) SampleRate() int {
	return int(cc.c.sample_rate)
}

func (cc *CodecContext) SetSampleRate(sampleRate int) {
	cc.c.sample_rate = C.int(sampleRate)
}

func (cc *CodecContext) StrictStdCompliance() StrictStdCompliance {
	return StrictStdCompliance(cc.c.strict_std_compliance)
}

func (cc *CodecContext) SetStrictStdCompliance(c StrictStdCompliance) {
	cc.c.strict_std_compliance = C.int(c)
}

func (cc *CodecContext) TimeBase() Rational {
	return newRationalFromC(cc.c.time_base)
}

func (cc *CodecContext) SetTimeBase(r Rational) {
	cc.c.time_base = r.c
}

func (cc *CodecContext) ThreadCount() int {
	return int(cc.c.thread_count)
}

func (cc *CodecContext) SetThreadCount(threadCount int) {
	cc.c.thread_count = C.int(threadCount)
}

func (cc *CodecContext) ThreadType() ThreadType {
	return ThreadType(cc.c.thread_type)
}

func (cc *CodecContext) SetThreadType(t ThreadType) {
	cc.c.thread_type = C.int(t)
}

func (cc *CodecContext) Width() int {
	return int(cc.c.width)
}

func (cc *CodecContext) SetWidth(width int) {
	cc.c.width = C.int(width)
}

func (cc *CodecContext) Open(c *Codec, d *Dictionary) error {
	var dc **C.struct_AVDictionary
	if d != nil {
		dc = &d.c
	}
	return newError(C.avcodec_open2(cc.c, c.c, dc))
}

func (cc *CodecContext) ReceivePacket(p *Packet) error {
	var pc *C.struct_AVPacket
	if p != nil {
		pc = p.c
	}
	return newError(C.avcodec_receive_packet(cc.c, pc))
}

func (cc *CodecContext) SendPacket(p *Packet) error {
	var pc *C.struct_AVPacket
	if p != nil {
		pc = p.c
	}
	return newError(C.avcodec_send_packet(cc.c, pc))
}

func (cc *CodecContext) ReceiveFrame(f *Frame) error {
	var fc *C.struct_AVFrame
	if f != nil {
		fc = f.c
	}
	return newError(C.avcodec_receive_frame(cc.c, fc))
}

func (cc *CodecContext) SendFrame(f *Frame) error {
	var fc *C.struct_AVFrame
	if f != nil {
		fc = f.c
	}
	return newError(C.avcodec_send_frame(cc.c, fc))
}

func (cc *CodecContext) ToCodecParameters(cp *CodecParameters) error {
	return cp.FromCodecContext(cc)
}

func (cc *CodecContext) FromCodecParameters(cp *CodecParameters) error {
	return cp.ToCodecContext(cc)
}

func (cc *CodecContext) SetHardwareDeviceContext(hdc *HardwareDeviceContext) {
	if cc.hdc != nil {
		C.av_buffer_unref(&cc.hdc.c)
	}
	cc.hdc = hdc
	if cc.hdc != nil {
		cc.c.hw_device_ctx = C.av_buffer_ref(cc.hdc.c)
	}
}

func (cc *CodecContext) ExtraHardwareFrames() int {
	return int(cc.c.extra_hw_frames)
}

func (cc *CodecContext) SetExtraHardwareFrames(n int) {
	cc.c.extra_hw_frames = C.int(n)
}

type CodecContextPixelFormatCallback func(pfs []PixelFormat) PixelFormat

var (
	codecContextPixelFormatCallbacks      = make(map[*C.struct_AVCodecContext]CodecContextPixelFormatCallback)
	codecContextPixelFormatCallbacksMutex = &sync.Mutex{}
)

func (cc *CodecContext) SetPixelFormatCallback(c CodecContextPixelFormatCallback) {
	// Lock
	codecContextPixelFormatCallbacksMutex.Lock()
	defer codecContextPixelFormatCallbacksMutex.Unlock()

	// Update callback
	if c == nil {
		C.astiavResetCodecContextGetFormat(cc.c)
		delete(codecContextPixelFormatCallbacks, cc.c)
	} else {
		C.astiavSetCodecContextGetFormat(cc.c)
		codecContextPixelFormatCallbacks[cc.c] = c
	}
}

//export goAstiavCodecContextGetFormat
func goAstiavCodecContextGetFormat(cc *C.struct_AVCodecContext, pfsCPtr *C.enum_AVPixelFormat, pfsCSize C.int) C.enum_AVPixelFormat {
	// Lock
	codecContextPixelFormatCallbacksMutex.Lock()
	defer codecContextPixelFormatCallbacksMutex.Unlock()

	// Get callback
	c, ok := codecContextPixelFormatCallbacks[cc]
	if !ok {
		return C.enum_AVPixelFormat(PixelFormatNone)
	}

	// Get pixel formats
	var pfs []PixelFormat
	for _, v := range unsafe.Slice(pfsCPtr, pfsCSize) {
		pfs = append(pfs, PixelFormat(v))
	}

	// Callback
	return C.enum_AVPixelFormat(c(pfs))

}
