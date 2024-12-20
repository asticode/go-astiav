package astiav

//#include <libavfilter/buffersrc.h>
import "C"
import "unsafe"

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html
type BuffersrcFilterContextParameters struct {
	c *C.AVBufferSrcParameters
}

func newBuffersrcFilterContextParametersFromC(c *C.AVBufferSrcParameters) *BuffersrcFilterContextParameters {
	if c == nil {
		return nil
	}
	return &BuffersrcFilterContextParameters{c: c}
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersrc.html#gaae82d4f8a69757ce01421dd3167861a5
func AllocBuffersrcFilterContextParameters() *BuffersrcFilterContextParameters {
	return newBuffersrcFilterContextParametersFromC(C.av_buffersrc_parameters_alloc())
}

func (bfcp *BuffersrcFilterContextParameters) Free() {
	if bfcp.c != nil {
		if bfcp.c.hw_frames_ctx != nil {
			C.av_buffer_unref(&bfcp.c.hw_frames_ctx)
		}
		C.av_freep(unsafe.Pointer(&bfcp.c))
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a5267368bf88b4f2a65a5e06ac3f9ecd4
func (bfcp *BuffersrcFilterContextParameters) ChannelLayout() ChannelLayout {
	l, _ := newChannelLayoutFromC(&bfcp.c.ch_layout).clone()
	return l
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a5267368bf88b4f2a65a5e06ac3f9ecd4
func (bfcp *BuffersrcFilterContextParameters) SetChannelLayout(l ChannelLayout) {
	l.copy(&bfcp.c.ch_layout) //nolint: errcheck
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a442add2b039f416dd7c92ccf1ccd0d3b
func (bfcp *BuffersrcFilterContextParameters) ColorRange() ColorRange {
	return ColorRange(bfcp.c.color_range)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a442add2b039f416dd7c92ccf1ccd0d3b
func (bfcp *BuffersrcFilterContextParameters) SetColorRange(r ColorRange) {
	bfcp.c.color_range = C.enum_AVColorRange(r)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a700226626af70f787c930d7506554757
func (bfcp *BuffersrcFilterContextParameters) ColorSpace() ColorSpace {
	return ColorSpace(bfcp.c.color_space)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a700226626af70f787c930d7506554757
func (bfcp *BuffersrcFilterContextParameters) SetColorSpace(s ColorSpace) {
	bfcp.c.color_space = C.enum_AVColorSpace(s)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a032a202496206e18449c66233058647a
func (bfcp *BuffersrcFilterContextParameters) Framerate() Rational {
	return newRationalFromC(bfcp.c.frame_rate)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a032a202496206e18449c66233058647a
func (bfcp *BuffersrcFilterContextParameters) SetFramerate(f Rational) {
	bfcp.c.frame_rate = f.c
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a86c49b4202433037c9e2b0b6ae541534
func (bfcp *BuffersrcFilterContextParameters) SetHardwareFrameContext(hfc *HardwareFrameContext) {
	if bfcp.c.hw_frames_ctx != nil {
		C.av_buffer_unref(&bfcp.c.hw_frames_ctx)
	}
	if hfc != nil {
		bfcp.c.hw_frames_ctx = C.av_buffer_ref(hfc.c)
	} else {
		bfcp.c.hw_frames_ctx = nil
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a89783d603b84908fb1998bbbea981b70
func (bfcp *BuffersrcFilterContextParameters) Height() int {
	return int(bfcp.c.height)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a89783d603b84908fb1998bbbea981b70
func (bfcp *BuffersrcFilterContextParameters) SetHeight(height int) {
	bfcp.c.height = C.int(height)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a56f28f81f1a86cecc39a8d61674912b8
func (bfcp *BuffersrcFilterContextParameters) PixelFormat() PixelFormat {
	return PixelFormat(bfcp.c.format)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a56f28f81f1a86cecc39a8d61674912b8
func (bfcp *BuffersrcFilterContextParameters) SetPixelFormat(f PixelFormat) {
	bfcp.c.format = C.int(f)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#ae47c141ea7a7770351613242229f951a
func (bfcp *BuffersrcFilterContextParameters) SampleAspectRatio() Rational {
	return newRationalFromC(bfcp.c.sample_aspect_ratio)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#ae47c141ea7a7770351613242229f951a
func (bfcp *BuffersrcFilterContextParameters) SetSampleAspectRatio(r Rational) {
	bfcp.c.sample_aspect_ratio = r.c
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a56f28f81f1a86cecc39a8d61674912b8
func (bfcp *BuffersrcFilterContextParameters) SampleFormat() SampleFormat {
	return SampleFormat(bfcp.c.format)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a56f28f81f1a86cecc39a8d61674912b8
func (bfcp *BuffersrcFilterContextParameters) SetSampleFormat(f SampleFormat) {
	bfcp.c.format = C.int(f)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a34a1613f1e80f8520c159fac59e29834
func (bfcp *BuffersrcFilterContextParameters) SampleRate() int {
	return int(bfcp.c.sample_rate)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a34a1613f1e80f8520c159fac59e29834
func (bfcp *BuffersrcFilterContextParameters) SetSampleRate(sampleRate int) {
	bfcp.c.sample_rate = C.int(sampleRate)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a7767325c1259942a33586f05c90e38b0
func (bfcp *BuffersrcFilterContextParameters) TimeBase() Rational {
	return newRationalFromC(bfcp.c.time_base)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a7767325c1259942a33586f05c90e38b0
func (bfcp *BuffersrcFilterContextParameters) SetTimeBase(r Rational) {
	bfcp.c.time_base = r.c
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a6c6f3d9ed8b427070e9055e7ac61263f
func (bfcp *BuffersrcFilterContextParameters) Width() int {
	return int(bfcp.c.width)
}

// https://ffmpeg.org/doxygen/7.0/structAVBufferSrcParameters.html#a6c6f3d9ed8b427070e9055e7ac61263f
func (bfcp *BuffersrcFilterContextParameters) SetWidth(width int) {
	bfcp.c.width = C.int(width)
}
