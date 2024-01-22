package astiav

//#cgo pkg-config: libavcodec libavutil
//#include <libavcodec/avcodec.h>
//#include <libavutil/frame.h>
import "C"
import (
	"fmt"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/avcodec.h#L383
type CodecContext struct {
	c *C.struct_AVCodecContext
}

type HWDeviceType int

const (
	HWDeviceTypeNone         HWDeviceType = C.AV_HWDEVICE_TYPE_NONE
	HWDeviceTypeVDPAU        HWDeviceType = C.AV_HWDEVICE_TYPE_VDPAU
	HWDeviceTypeCUDA         HWDeviceType = C.AV_HWDEVICE_TYPE_CUDA
	HWDeviceTypeVAAPI        HWDeviceType = C.AV_HWDEVICE_TYPE_VAAPI
	HWDeviceTypeDXVA2        HWDeviceType = C.AV_HWDEVICE_TYPE_DXVA2
	HWDeviceTypeQSV          HWDeviceType = C.AV_HWDEVICE_TYPE_QSV
	HWDeviceTypeVideoToolbox HWDeviceType = C.AV_HWDEVICE_TYPE_VIDEOTOOLBOX
	HWDeviceTypeD3D11VA      HWDeviceType = C.AV_HWDEVICE_TYPE_D3D11VA
	HWDeviceTypeDRM          HWDeviceType = C.AV_HWDEVICE_TYPE_DRM
	HWDeviceTypeOpenCL       HWDeviceType = C.AV_HWDEVICE_TYPE_OPENCL
	HWDeviceTypeMediaCodec   HWDeviceType = C.AV_HWDEVICE_TYPE_MEDIACODEC
	HWDeviceTypeVulkan       HWDeviceType = C.AV_HWDEVICE_TYPE_VULKAN
)

func AllocCodecContext(c *Codec) *CodecContext {
	var cc *C.struct_AVCodec
	if c != nil {
		cc = c.c
	}
	return newCodecContextFromC(C.avcodec_alloc_context3(cc))
}

func AllocHWDeviceContext(c *Codec, hwType HWDeviceType) *CodecContext {
	ctx := AllocCodecContext(c)
	if ctx == nil {
		return nil
	}

	var hwDeviceCtx *C.AVBufferRef
	errorCode := C.av_hwdevice_ctx_create(&hwDeviceCtx, C.enum_AVHWDeviceType(hwType), nil, nil, 0)
	if errorCode < 0 {
		ctx.Free()
		return nil
	}

	ctx.c.hw_device_ctx = hwDeviceCtx

	return ctx
}

func AllocHWDeviceContextWithDevice(c *Codec, hwType HWDeviceType, device string) *CodecContext {
	ctx := AllocCodecContext(c)
	if ctx == nil {
		return nil
	}

	var hwDeviceCtx *C.AVBufferRef
	deviceC := C.CString(device)
	defer C.free(unsafe.Pointer(deviceC))

	errorCode := C.av_hwdevice_ctx_create(&hwDeviceCtx, C.enum_AVHWDeviceType(hwType), deviceC, nil, 0)
	if errorCode < 0 {
		ctx.Free()
		return nil
	}

	ctx.c.hw_device_ctx = hwDeviceCtx

	return ctx
}

func newCodecContextFromC(c *C.struct_AVCodecContext) *CodecContext {
	if c == nil {
		return nil
	}
	return &CodecContext{c: c}
}

func (cc *CodecContext) Free() {
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
	err := newError(C.avcodec_receive_frame(cc.c, fc))
	if err != nil {
		return err
	}
	if isHardwarePixelFormat(f.PixelFormat()) {
		temp_frame := AllocFrame()
		ret := C.av_hwframe_transfer_data(temp_frame.c, f.c, 0)
		if int(ret) < 0 {
			return fmt.Errorf("Unable to transfer data from GPU: %d", int(ret))
		}
		f.Free()
		f.c = temp_frame.c
	}
	return nil
}

func isHardwarePixelFormat(pf PixelFormat) bool {
	hwPixelFormats := []PixelFormat{
		PixelFormatCuda,
		PixelFormatD3D11,
		PixelFormatQsv,
		PixelFormatD3D11VaVld,
		PixelFormatDxva2Vld,
		PixelFormatVaapi,
		PixelFormatVdpau,
		PixelFormatVideotoolbox,
		PixelFormatMmal,
		PixelFormatXvmc,
	}

	for _, hwPf := range hwPixelFormats {
		if pf == hwPf {
			return true
		}
	}
	return false
}

func (cc *CodecContext) SendFrame(f *Frame) error {
	var fc *C.struct_AVFrame
	if f != nil {
		fc = f.c
	}
	return newError(C.avcodec_send_frame(cc.c, fc))
}
