package astiav

//#include <libavutil/channel_layout.h>
//#include <libavutil/frame.h>
//#include <libavutil/imgutils.h>
//#include <libavutil/samplefmt.h>
//#include <libavutil/hwcontext.h>
//#include "frame.h"
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.1/frame_8h.html#add80189702cf0f5ea82718576fb43201
const NumDataPointers = uint(C.AV_NUM_DATA_POINTERS)

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html
type Frame struct {
	c *C.AVFrame
}

func newFrameFromC(c *C.AVFrame) *Frame {
	if c == nil {
		return nil
	}
	return &Frame{c: c}
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#gac700017c5270c79c1e1befdeeb008b2f
func AllocFrame() *Frame {
	return newFrameFromC(C.av_frame_alloc())
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#ga6b1acbfa82c79bf7fd78d868572f0ceb
func (f *Frame) AllocBuffer(align int) error {
	return newError(C.av_frame_get_buffer(f.c, C.int(align)))
}

// https://ffmpeg.org/doxygen/7.1/hwcontext_8c.html#adfa5aaa3a4f69b163ea30cadc6d663dc
func (f *Frame) AllocHardwareBuffer(hfc *HardwareFramesContext) error {
	return newError(C.av_hwframe_get_buffer(hfc.c, f.c, 0))
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#ae291cdec7758599e765bc9e3edbb3065
func (f *Frame) ChannelLayout() ChannelLayout {
	l, _ := newChannelLayoutFromC(&f.c.ch_layout).clone()
	return l
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#ae291cdec7758599e765bc9e3edbb3065
func (f *Frame) SetChannelLayout(l ChannelLayout) {
	l.copy(&f.c.ch_layout) //nolint: errcheck
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a853afbad220bbc58549b4860732a3aa5
func (f *Frame) ColorRange() ColorRange {
	return ColorRange(f.c.color_range)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a853afbad220bbc58549b4860732a3aa5
func (f *Frame) SetColorRange(r ColorRange) {
	f.c.color_range = C.enum_AVColorRange(r)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a9262c231f1f64869439b4fe587fe1710
func (f *Frame) ColorSpace() ColorSpace {
	return ColorSpace(f.c.colorspace)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a9262c231f1f64869439b4fe587fe1710
func (f *Frame) SetColorSpace(s ColorSpace) {
	f.c.colorspace = C.enum_AVColorSpace(s)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a1d0f65014a8d1bf78cec8cbed2304992
func (f *Frame) Data() *FrameData {
	return newFrameData(newFrameDataFrame(f))
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a29493fbfabaa21432c360a090426aa8e
func (f *Frame) HardwareFramesContext() *HardwareFramesContext {
	return newHardwareFramesContextFromC(f.c.hw_frames_ctx)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a3f89733f429c98ba5bc64373fb0a3f13
func (f *Frame) Height() int {
	return int(f.c.height)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a3f89733f429c98ba5bc64373fb0a3f13
func (f *Frame) SetHeight(h int) {
	f.c.height = C.int(h)
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame__flags.html#gafe155269fc8dc3a484490bd19b86cc40
func (f *Frame) KeyFrame() bool {
	return f.Flags().Has(FrameFlagKey)
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame__flags.html#gafe155269fc8dc3a484490bd19b86cc40
func (f *Frame) SetKeyFrame(k bool) {
	if k {
		f.SetFlags(f.Flags().Add(FrameFlagKey))
	} else {
		f.SetFlags(f.Flags().Del(FrameFlagKey))
	}
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a49020cc320b8fb1f5449167b6c97515b
func (f *Frame) Flags() FrameFlags {
	return FrameFlags(f.c.flags)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a49020cc320b8fb1f5449167b6c97515b
func (f *Frame) SetFlags(fs FrameFlags) {
	f.c.flags = C.int(fs)
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__picture.html#ga24a67963c3ae0054a2a4bab35930e694
func (f *Frame) ImageBufferSize(align int) (int, error) {
	ret := C.av_image_get_buffer_size((C.enum_AVPixelFormat)(f.c.format), f.c.width, f.c.height, C.int(align))
	if err := newError(ret); err != nil {
		return 0, err
	}
	return int(ret), nil
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__picture.html#ga6f8576f1ef0c2d9a9f7c5ac7f9a28c52
func (f *Frame) ImageCopyToBuffer(b []byte, align int) (int, error) {
	ret := C.av_image_copy_to_buffer((*C.uint8_t)(unsafe.Pointer(&b[0])), C.int(len(b)), &f.c.data[0], &f.c.linesize[0], (C.enum_AVPixelFormat)(f.c.format), f.c.width, f.c.height, C.int(align))
	if err := newError(ret); err != nil {
		return 0, err
	}
	return int(ret), nil
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__picture.html#ga3fa8e484cc214e8c7b9026825b5f4078
func (f *Frame) ImageFillBlack() error {
	linesize := [NumDataPointers]C.ptrdiff_t{}
	for i := 0; i < int(NumDataPointers); i++ {
		linesize[i] = C.ptrdiff_t(f.c.linesize[i])
	}
	return newError(C.av_image_fill_black(&f.c.data[0], &linesize[0], (C.enum_AVPixelFormat)(f.c.format), (C.enum_AVColorRange)(f.c.color_range), f.c.width, f.c.height))
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__sampfmts.html#gaa7368bc4e3a366b688e81938ed55eb06
func (f *Frame) SamplesBufferSize(align int) (int, error) {
	ret := C.av_samples_get_buffer_size(nil, f.c.ch_layout.nb_channels, f.c.nb_samples, (C.enum_AVSampleFormat)(f.c.format), C.int(align))
	if err := newError(ret); err != nil {
		return 0, err
	}
	return int(ret), nil
}

func (f *Frame) SamplesCopyToBuffer(b []byte, align int) (int, error) {
	ret := C.astiavSamplesCopyToBuffer((*C.uint8_t)(unsafe.Pointer(&b[0])), C.int(len(b)), &f.c.data[0], f.c.ch_layout.nb_channels, f.c.nb_samples, (C.enum_AVSampleFormat)(f.c.format), C.int(align))
	if err := newError(ret); err != nil {
		return 0, err
	}
	return int(ret), nil
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__sampmanip.html#gabcb166e22938c7d93c2d609529c458bb
func (f *Frame) SamplesFillSilence() error {
	return newError(C.av_samples_set_silence(&f.c.data[0], 0, f.c.nb_samples, f.c.ch_layout.nb_channels, (C.enum_AVSampleFormat)(f.c.format)))
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#aa52bfc6605f6a3059a0c3226cc0f6567
func (f *Frame) Linesize() [NumDataPointers]int {
	o := [NumDataPointers]int{}
	for i := 0; i < int(NumDataPointers); i++ {
		o[i] = int(f.c.linesize[i])
	}
	return o
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a02f45ab8191aea1660159f1e464237ea
func (f *Frame) NbSamples() int {
	return int(f.c.nb_samples)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a02f45ab8191aea1660159f1e464237ea
func (f *Frame) SetNbSamples(n int) {
	f.c.nb_samples = C.int(n)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#af9920fc3fbfa347b8943ae461b50d18b
func (f *Frame) PictureType() PictureType {
	return PictureType(f.c.pict_type)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#af9920fc3fbfa347b8943ae461b50d18b
func (f *Frame) SetPictureType(t PictureType) {
	f.c.pict_type = C.enum_AVPictureType(t)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#aed14fa772ce46881020fd1545c86432c
func (f *Frame) PixelFormat() PixelFormat {
	return PixelFormat(f.c.format)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#aed14fa772ce46881020fd1545c86432c
func (f *Frame) SetPixelFormat(pf PixelFormat) {
	f.c.format = C.int(pf)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#aa52951f35ec9e303d3dfeb4b3e44248a
func (f *Frame) PktDts() int64 {
	return int64(f.c.pkt_dts)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a0452833e3ab6ddd7acbf82817a7818a4
func (f *Frame) Pts() int64 {
	return int64(f.c.pts)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a0452833e3ab6ddd7acbf82817a7818a4
func (f *Frame) SetPts(i int64) {
	f.c.pts = C.int64_t(i)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a62f9c20541a83d37db7072126ff0060d
func (f *Frame) SampleAspectRatio() Rational {
	return newRationalFromC(f.c.sample_aspect_ratio)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a62f9c20541a83d37db7072126ff0060d
func (f *Frame) SetSampleAspectRatio(r Rational) {
	f.c.sample_aspect_ratio = r.c
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#aed14fa772ce46881020fd1545c86432c
func (f *Frame) SampleFormat() SampleFormat {
	return SampleFormat(f.c.format)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#aed14fa772ce46881020fd1545c86432c
func (f *Frame) SetSampleFormat(sf SampleFormat) {
	f.c.format = C.int(sf)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#ac85daa1316e1f47e78da0ca19b7c60e6
func (f *Frame) SampleRate() int {
	return int(f.c.sample_rate)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#ac85daa1316e1f47e78da0ca19b7c60e6
func (f *Frame) SetSampleRate(r int) {
	f.c.sample_rate = C.int(r)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a44d40e03fe22a0511c9157dab22143ee
func (f *Frame) SideData() *FrameSideData {
	return newFrameSideDataFromC(&f.c.side_data, &f.c.nb_side_data)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a5bde87fd101f66d6263bb451056dba13
func (f *Frame) Metadata() *Dictionary {
	return newDictionaryFromC(f.c.metadata)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a5bde87fd101f66d6263bb451056dba13
func (f *Frame) SetMetadata(d *Dictionary) {
	if d == nil {
		f.c.metadata = nil
	} else {
		f.c.metadata = d.c
	}
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a1e71ce60cedd5f3b6811714a9f7f9e0a
func (f *Frame) Width() int {
	return int(f.c.width)
}

// https://ffmpeg.org/doxygen/7.1/structAVFrame.html#a1e71ce60cedd5f3b6811714a9f7f9e0a
func (f *Frame) SetWidth(w int) {
	f.c.width = C.int(w)
}

// https://ffmpeg.org/doxygen/7.1/hwcontext_8c.html#abf1b1664b8239d953ae2cac8b643815a
func (f *Frame) TransferHardwareData(dst *Frame) error {
	return newError(C.av_hwframe_transfer_data(dst.c, f.c, 0))
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#ga979d73f3228814aee56aeca0636e37cc
func (f *Frame) Free() {
	if f.c != nil {
		C.av_frame_free(&f.c)
	}
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#ga88b0ecbc4eb3453eef3fbefa3bddeb7c
func (f *Frame) Ref(src *Frame) error {
	return newError(C.av_frame_ref(f.c, src.c))
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#ga46d6d32f6482a3e9c19203db5877105b
func (f *Frame) Clone() *Frame {
	return newFrameFromC(C.av_frame_clone(f.c))
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#ga0a2b687f9c1c5ed0089b01fd61227108
func (f *Frame) Unref() {
	C.av_frame_unref(f.c)
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#ga709e62bc2917ffd84c5c0f4e1dfc48f7
func (f *Frame) MoveRef(src *Frame) {
	C.av_frame_move_ref(f.c, src.c)
}

func (f *Frame) UnsafePointer() unsafe.Pointer {
	return unsafe.Pointer(f.c)
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#ga3ba755bada5c3c8883361ef43fb5fb7a
func (f *Frame) IsWritable() bool {
	return C.av_frame_is_writable(f.c) > 0
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#gadd5417c06f5a6b419b0dbd8f0ff363fd
func (f *Frame) MakeWritable() error {
	return newError(C.av_frame_make_writable(f.c))
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#gaec4e92f6e1e75ffaf76e07586fb0c9ed
func (f *Frame) Copy(dst *Frame) error {
	return newError(C.av_frame_copy(dst.c, f.c))
}
