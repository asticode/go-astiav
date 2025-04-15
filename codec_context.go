package astiav

//#include "codec_context.h"
import "C"
import (
	"sync"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html
type CodecContext struct {
	classerHandler
	c *C.AVCodecContext
}

func newCodecContextFromC(c *C.AVCodecContext) *CodecContext {
	if c == nil {
		return nil
	}
	cc := &CodecContext{c: c}
	classers.set(cc)
	return cc
}

var _ Classer = (*CodecContext)(nil)

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#gae80afec6f26df6607eaacf39b561c315
func AllocCodecContext(c *Codec) *CodecContext {
	var cc *C.AVCodec
	if c != nil {
		cc = c.c
	}
	return newCodecContextFromC(C.avcodec_alloc_context3(cc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#gaf869d0829ed607cec3a4a02a1c7026b3
func (cc *CodecContext) Free() {
	if cc.c != nil {
		if cc.c.hw_device_ctx != nil {
			C.av_buffer_unref(&cc.c.hw_device_ctx)
		}
		if cc.c.hw_frames_ctx != nil {
			C.av_buffer_unref(&cc.c.hw_frames_ctx)
		}
		// Make sure to clone the classer before freeing the object since
		// the C free method may reset the pointer
		c := newClonedClasser(cc)
		C.avcodec_free_context(&cc.c)
		// Make sure to remove from classers after freeing the object since
		// the C free method may use methods needing the classer
		if c != nil {
			classers.del(c)
		}
	}
}

func (cc *CodecContext) String() string {
	s, _ := stringFromC(255, func(buf *C.char, size C.size_t) error {
		C.avcodec_string(buf, C.int(size), cc.c, C.int(0))
		return nil
	})
	return s
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a6b53fda85ad61baa345edbd96cb8a33c
func (cc *CodecContext) BitRate() int64 {
	return int64(cc.c.bit_rate)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a6b53fda85ad61baa345edbd96cb8a33c
func (cc *CodecContext) SetBitRate(bitRate int64) {
	cc.c.bit_rate = C.int64_t(bitRate)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a167ff73c67960acf2d5ca73d93e13f64
func (cc *CodecContext) ChannelLayout() ChannelLayout {
	l, _ := newChannelLayoutFromC(&cc.c.ch_layout).clone()
	return l
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a167ff73c67960acf2d5ca73d93e13f64
func (cc *CodecContext) SetChannelLayout(l ChannelLayout) {
	l.copy(&cc.c.ch_layout) //nolint: errcheck
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#ac60a0209642b5d74068cab0ac35a78b2
func (cc *CodecContext) ChromaLocation() ChromaLocation {
	return ChromaLocation(cc.c.chroma_sample_location)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a90622d3af2a9abba986a1c9f7ca21b16
func (cc *CodecContext) Class() *Class {
	if cc.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(cc.c))
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#adc5f65d6099fd8339c1580c091777223
func (cc *CodecContext) CodecID() CodecID {
	return CodecID(cc.c.codec_id)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a3a41b3e5bde23b877799f6e72dac8ef3
func (cc *CodecContext) ColorPrimaries() ColorPrimaries {
	return ColorPrimaries(cc.c.color_primaries)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a255bf7100a4ba6dcb6ee5d87740a4f35
func (cc *CodecContext) ColorRange() ColorRange {
	return ColorRange(cc.c.color_range)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a8cd8caa7d40319324ce3d879a2edbd9f
func (cc *CodecContext) ColorSpace() ColorSpace {
	return ColorSpace(cc.c.colorspace)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#ab649e8c599f5a0e2a30448e67a36deb6
func (cc *CodecContext) ColorTransferCharacteristic() ColorTransferCharacteristic {
	return ColorTransferCharacteristic(cc.c.color_trc)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#abe964316aaaa61967b012efdcced79c4
func (cc *CodecContext) ExtraData() []byte {
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		*size = C.size_t(cc.c.extradata_size)
		return cc.c.extradata
	})
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#abe964316aaaa61967b012efdcced79c4
func (cc *CodecContext) SetExtraData(b []byte) error {
	return setBytesWithIntSizeInC(b, &cc.c.extradata, &cc.c.extradata_size)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#abb01e291550fa3fb96188af4d494587e
func (cc *CodecContext) Flags() CodecContextFlags {
	return CodecContextFlags(cc.c.flags)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#abb01e291550fa3fb96188af4d494587e
func (cc *CodecContext) SetFlags(fs CodecContextFlags) {
	cc.c.flags = C.int(fs)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a1944f9a4f8f2e123c087e1fe7613d571
func (cc *CodecContext) Flags2() CodecContextFlags2 {
	return CodecContextFlags2(cc.c.flags2)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a1944f9a4f8f2e123c087e1fe7613d571
func (cc *CodecContext) SetFlags2(fs CodecContextFlags2) {
	cc.c.flags2 = C.int(fs)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a4d08b297e97eefd66c714df4fff493c8
func (cc *CodecContext) Framerate() Rational {
	return newRationalFromC(cc.c.framerate)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a4d08b297e97eefd66c714df4fff493c8
func (cc *CodecContext) SetFramerate(f Rational) {
	cc.c.framerate = f.c
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#aec57f0d859a6df8b479cd93ca3a44a33
func (cc *CodecContext) FrameSize() int {
	return int(cc.c.frame_size)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a9b6b3f1fcbdcc2ad9f4dbb4370496e38
func (cc *CodecContext) GopSize() int {
	return int(cc.c.gop_size)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a9b6b3f1fcbdcc2ad9f4dbb4370496e38
func (cc *CodecContext) SetGopSize(gopSize int) {
	cc.c.gop_size = C.int(gopSize)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a0449afd803eb107bd4dbc8b5ea22e363
func (cc *CodecContext) Height() int {
	return int(cc.c.height)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a0449afd803eb107bd4dbc8b5ea22e363
func (cc *CodecContext) SetHeight(height int) {
	cc.c.height = C.int(height)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a6927dc652ae6241f1dfdbad4e12d3a40
func (cc *CodecContext) Level() Level {
	return Level(cc.c.level)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a6927dc652ae6241f1dfdbad4e12d3a40
func (cc *CodecContext) SetLevel(l Level) {
	cc.c.level = C.int(l)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a3f99ca3115c44e6d7772c9384faf15e6
func (cc *CodecContext) MediaType() MediaType {
	return MediaType(cc.c.codec_type)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a0425c77b3d06d71e5db88b1d7e1b37f2
func (cc *CodecContext) PixelFormat() PixelFormat {
	return PixelFormat(cc.c.pix_fmt)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a0425c77b3d06d71e5db88b1d7e1b37f2
func (cc *CodecContext) SetPixelFormat(pixFmt PixelFormat) {
	cc.c.pix_fmt = C.enum_AVPixelFormat(pixFmt)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#af3379123060ad8cc9c321c29af4f8360
func (cc *CodecContext) PrivateData() *PrivateData {
	return newPrivateDataFromC(cc.c.priv_data)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a7abe7095de73df98df4895bf9e25fc6b
func (cc *CodecContext) Profile() Profile {
	return Profile(cc.c.profile)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a7abe7095de73df98df4895bf9e25fc6b
func (cc *CodecContext) SetProfile(p Profile) {
	cc.c.profile = C.int(p)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a3f63bc9141e25bf7f1cda0cef7cd4a60
func (cc *CodecContext) Qmin() int {
	return int(cc.c.qmin)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a3f63bc9141e25bf7f1cda0cef7cd4a60
func (cc *CodecContext) SetQmin(qmin int) {
	cc.c.qmin = C.int(qmin)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a5252d34fbce300228d4dbda19a8c3293
func (cc *CodecContext) SampleAspectRatio() Rational {
	return newRationalFromC(cc.c.sample_aspect_ratio)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a5252d34fbce300228d4dbda19a8c3293
func (cc *CodecContext) SetSampleAspectRatio(r Rational) {
	cc.c.sample_aspect_ratio = r.c
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a1bdba69ea111e2a9d03fdaa7a46a4c45
func (cc *CodecContext) SampleFormat() SampleFormat {
	return SampleFormat(cc.c.sample_fmt)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a1bdba69ea111e2a9d03fdaa7a46a4c45
func (cc *CodecContext) SetSampleFormat(f SampleFormat) {
	cc.c.sample_fmt = C.enum_AVSampleFormat(f)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a8ff0b000c463361e234af48d03aadfc0
func (cc *CodecContext) SampleRate() int {
	return int(cc.c.sample_rate)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a8ff0b000c463361e234af48d03aadfc0
func (cc *CodecContext) SetSampleRate(sampleRate int) {
	cc.c.sample_rate = C.int(sampleRate)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a3090804569341ca235e3adbdc03318d2
func (cc *CodecContext) StrictStdCompliance() StrictStdCompliance {
	return StrictStdCompliance(cc.c.strict_std_compliance)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a3090804569341ca235e3adbdc03318d2
func (cc *CodecContext) SetStrictStdCompliance(c StrictStdCompliance) {
	cc.c.strict_std_compliance = C.int(c)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#ab7bfeb9fa5840aac090e2b0bd0ef7589
func (cc *CodecContext) TimeBase() Rational {
	return newRationalFromC(cc.c.time_base)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#ab7bfeb9fa5840aac090e2b0bd0ef7589
func (cc *CodecContext) SetTimeBase(r Rational) {
	cc.c.time_base = r.c
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#aa852b6227d0778b62e9cc4034ad3720c
func (cc *CodecContext) ThreadCount() int {
	return int(cc.c.thread_count)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#aa852b6227d0778b62e9cc4034ad3720c
func (cc *CodecContext) SetThreadCount(threadCount int) {
	cc.c.thread_count = C.int(threadCount)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a7651614f4309122981d70e06a4b42fcb
func (cc *CodecContext) ThreadType() ThreadType {
	return ThreadType(cc.c.thread_type)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a7651614f4309122981d70e06a4b42fcb
func (cc *CodecContext) SetThreadType(t ThreadType) {
	cc.c.thread_type = C.int(t)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a0d8f46461754e8abea0847dcbc41b956
func (cc *CodecContext) Width() int {
	return int(cc.c.width)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a0d8f46461754e8abea0847dcbc41b956
func (cc *CodecContext) SetWidth(width int) {
	cc.c.width = C.int(width)
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga11f785a188d7d9df71621001465b0f1d
func (cc *CodecContext) Open(c *Codec, d *Dictionary) error {
	var dc **C.AVDictionary
	if d != nil {
		dc = &d.c
	}
	return cc.newError(C.avcodec_open2(cc.c, c.c, dc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__decoding.html#ga5b8eff59cf259747cf0b31563e38ded6
func (cc *CodecContext) ReceivePacket(p *Packet) error {
	var pc *C.AVPacket
	if p != nil {
		pc = p.c
	}
	return cc.newError(C.avcodec_receive_packet(cc.c, pc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__decoding.html#ga58bc4bf1e0ac59e27362597e467efff3
func (cc *CodecContext) SendPacket(p *Packet) error {
	var pc *C.AVPacket
	if p != nil {
		pc = p.c
	}
	return cc.newError(C.avcodec_send_packet(cc.c, pc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__decoding.html#ga11e6542c4e66d3028668788a1a74217c
func (cc *CodecContext) ReceiveFrame(f *Frame) error {
	var fc *C.AVFrame
	if f != nil {
		fc = f.c
	}
	return cc.newError(C.avcodec_receive_frame(cc.c, fc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__decoding.html#ga9395cb802a5febf1f00df31497779169
func (cc *CodecContext) SendFrame(f *Frame) error {
	var fc *C.AVFrame
	if f != nil {
		fc = f.c
	}
	return cc.newError(C.avcodec_send_frame(cc.c, fc))
}

func (cc *CodecContext) ToCodecParameters(cp *CodecParameters) error {
	return cp.FromCodecContext(cc)
}

func (cc *CodecContext) FromCodecParameters(cp *CodecParameters) error {
	return cp.ToCodecContext(cc)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#acf8113e490f9e7b57465e65af9c0c75c
func (cc *CodecContext) SetHardwareDeviceContext(hdc *HardwareDeviceContext) {
	if cc.c.hw_device_ctx != nil {
		C.av_buffer_unref(&cc.c.hw_device_ctx)
	}
	if hdc != nil {
		cc.c.hw_device_ctx = C.av_buffer_ref(hdc.c)
	} else {
		cc.c.hw_device_ctx = nil
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a3bac44bb0b016ab838780cc19ac277d6
func (cc *CodecContext) HardwareFramesContext() *HardwareFramesContext {
	return newHardwareFramesContextFromC(cc.c.hw_frames_ctx)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a3bac44bb0b016ab838780cc19ac277d6
func (cc *CodecContext) SetHardwareFramesContext(hfc *HardwareFramesContext) {
	if cc.c.hw_frames_ctx != nil {
		C.av_buffer_unref(&cc.c.hw_frames_ctx)
	}
	if hfc != nil {
		cc.c.hw_frames_ctx = C.av_buffer_ref(hfc.c)
	} else {
		cc.c.hw_frames_ctx = nil
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#ad2f772bd948d8f3be4d674a3a52ee00e
func (cc *CodecContext) ExtraHardwareFrames() int {
	return int(cc.c.extra_hw_frames)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#ad2f772bd948d8f3be4d674a3a52ee00e
func (cc *CodecContext) SetExtraHardwareFrames(n int) {
	cc.c.extra_hw_frames = C.int(n)
}

func (cc *CodecContext) UnsafePointer() unsafe.Pointer {
	return unsafe.Pointer(cc.c)
}

type CodecContextPixelFormatCallback func(pfs []PixelFormat) PixelFormat

var (
	codecContextPixelFormatCallbacks      = make(map[*C.AVCodecContext]CodecContextPixelFormatCallback)
	codecContextPixelFormatCallbacksMutex = &sync.Mutex{}
)

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a360a2b8508a67c4234d97f4c13ba1bb5
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
func goAstiavCodecContextGetFormat(cc *C.AVCodecContext, pfsCPtr *C.enum_AVPixelFormat, pfsCSize C.int) C.enum_AVPixelFormat {
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

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a3e5334a611a3e2a6a653805bb9e2d4d4
func (cc *CodecContext) MaxBFrames() int {
	return int(cc.c.max_b_frames)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a3e5334a611a3e2a6a653805bb9e2d4d4
func (cc *CodecContext) SetMaxBFrames(n int) {
	cc.c.max_b_frames = C.int(n)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#aa2b5582f1a360534310b686cc3f7c668
func (cc *CodecContext) RateControlMaxRate() int64 {
	return int64(cc.c.rc_max_rate)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#aa2b5582f1a360534310b686cc3f7c668
func (cc *CodecContext) SetRateControlMaxRate(n int64) {
	cc.c.rc_max_rate = C.int64_t(n)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#ac265c70b89e87455ec05eb2978def81b
func (cc *CodecContext) RateControlMinRate() int64 {
	return int64(cc.c.rc_min_rate)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#ac265c70b89e87455ec05eb2978def81b
func (cc *CodecContext) SetRateControlMinRate(n int64) {
	cc.c.rc_min_rate = C.int64_t(n)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a15000607a7e2371162348bb35b0184c1
func (cc *CodecContext) RateControlBufferSize() int {
	return int(cc.c.rc_buffer_size)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecContext.html#a15000607a7e2371162348bb35b0184c1
func (cc *CodecContext) SetRateControlBufferSize(n int) {
	cc.c.rc_buffer_size = C.int(n)
}
