package astiav

//#cgo pkg-config: libavutil
//#include <libavutil/frame.h>
//#include <libavutil/samplefmt.h>
import "C"
import (
	"errors"
	"unsafe"
)

const NumDataPointers = uint(C.AV_NUM_DATA_POINTERS)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/frame.h#L317
type Frame struct {
	c *C.struct_AVFrame
}

func newFrameFromC(c *C.struct_AVFrame) *Frame {
	if c == nil {
		return nil
	}
	return &Frame{c: c}
}

func AllocFrame() *Frame {
	return newFrameFromC(C.av_frame_alloc())
}

func (f *Frame) AllocBuffer(align int) error {
	return newError(C.av_frame_get_buffer(f.c, C.int(align)))
}

func (f *Frame) AllocSamples(sf SampleFormat, nbChannels, nbSamples, align int) error {
	return newError(C.av_samples_alloc(&f.c.data[0], &f.c.linesize[0], C.int(nbChannels), C.int(nbSamples), (C.enum_AVSampleFormat)(sf), C.int(align)))
}

func (f *Frame) ChannelLayout() ChannelLayout {
	return ChannelLayout(f.c.channel_layout)
}

func (f *Frame) SetChannelLayout(l ChannelLayout) {
	f.c.channel_layout = C.uint64_t(l)
}

func (f *Frame) Data() [NumDataPointers][]byte {
	b := [NumDataPointers][]byte{}
	for i := 0; i < int(NumDataPointers); i++ {
		b[i] = bytesFromC(func(size *C.int) *C.uint8_t {
			*size = f.c.linesize[i]
			if f.c.height > 0 {
				*size = *size * f.c.height
			} else if f.c.channels > 0 {
				*size = *size * f.c.channels
			}
			return f.c.data[i]
		})
	}
	return b
}

// DataBytes returns the frame's data as byte slice.
func (f *Frame) DataBytes() ([]byte, error) {

	isVideoFrame := f.Height() > 0 && f.Width() > 0

	if isVideoFrame {

		// retrieve the image buffer size
		bufferSize := C.av_image_get_buffer_size(C.enum_AVPixelFormat(f.c.format), f.c.width, f.c.height, 1)
		if bufferSize < 0 {
			return nil, newError(bufferSize)
		}

		// create a buffer and copy the raw image data into it
		buffer := make([]byte, int(bufferSize))
		n := C.av_image_copy_to_buffer((*C.uint8_t)(unsafe.Pointer(&buffer[0])), bufferSize, (**C.uint8_t)(unsafe.Pointer(&f.c.data)), (*C.int)(unsafe.Pointer(&f.c.linesize)), C.enum_AVPixelFormat(f.c.format), f.c.width, f.c.height, 1)
		if n < 0 {
			return nil, newError(n)
		}

		return buffer, nil

	} else {
		return nil, errors.New("astiav: not implemented")
	}

}

func (f *Frame) Height() int {
	return int(f.c.height)
}

func (f *Frame) SetHeight(h int) {
	f.c.height = C.int(h)
}

func (f *Frame) KeyFrame() bool {
	return int(f.c.key_frame) > 0
}

func (f *Frame) SetKeyFrame(k bool) {
	i := 0
	if k {
		i = 1
	}
	f.c.key_frame = C.int(i)
}

func (f *Frame) Linesize() [NumDataPointers]int {
	o := [NumDataPointers]int{}
	for i := 0; i < int(NumDataPointers); i++ {
		o[i] = int(f.c.linesize[i])
	}
	return o
}

func (f *Frame) NbSamples() int {
	return int(f.c.nb_samples)
}

func (f *Frame) SetNbSamples(n int) {
	f.c.nb_samples = C.int(n)
}

func (f *Frame) PictureType() PictureType {
	return PictureType(f.c.pict_type)
}

func (f *Frame) SetPictureType(t PictureType) {
	f.c.pict_type = C.enum_AVPictureType(t)
}

func (f *Frame) PixelFormat() PixelFormat {
	return PixelFormat(f.c.format)
}

func (f *Frame) SetPixelFormat(pf PixelFormat) {
	f.c.format = C.int(pf)
}

func (f *Frame) PktPts() int64 {
	return int64(f.c.pkt_pts)
}

func (f *Frame) PktDts() int64 {
	return int64(f.c.pkt_dts)
}

func (f *Frame) Pts() int64 {
	return int64(f.c.pts)
}

func (f *Frame) SetPts(i int64) {
	f.c.pts = C.int64_t(i)
}

func (f *Frame) SampleFormat() SampleFormat {
	return SampleFormat(f.c.format)
}

func (f *Frame) SetSampleFormat(sf SampleFormat) {
	f.c.format = C.int(sf)
}

func (f *Frame) SampleRate() int {
	return int(f.c.sample_rate)
}

func (f *Frame) SetSampleRate(r int) {
	f.c.sample_rate = C.int(r)
}

func (f *Frame) Width() int {
	return int(f.c.width)
}

func (f *Frame) SetWidth(w int) {
	f.c.width = C.int(w)
}

func (f *Frame) Free() {
	C.av_frame_free(&f.c)
}

func (f *Frame) Ref(src *Frame) error {
	return newError(C.av_frame_ref(f.c, src.c))
}

func (f *Frame) Clone() *Frame {
	return newFrameFromC(C.av_frame_clone(f.c))
}

func (f *Frame) Unref() {
	C.av_frame_unref(f.c)
}

func (f *Frame) MoveRef(src *Frame) {
	C.av_frame_move_ref(f.c, src.c)
}
