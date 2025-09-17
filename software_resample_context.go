package astiav

//#include <libswresample/swresample.h>
import "C"
import "unsafe"

// https://ffmpeg.org/doxygen/8.0/structSwrContext.html
type SoftwareResampleContext struct {
	c *C.SwrContext
}

func newSoftwareResampleContextFromC(c *C.SwrContext) *SoftwareResampleContext {
	if c == nil {
		return nil
	}
	src := &SoftwareResampleContext{c: c}
	classers.set(src)
	return src
}

// https://ffmpeg.org/doxygen/8.0/group__lswr.html#gaf58c4ff10f73d74bdab8e5aa7193147c
func AllocSoftwareResampleContext() *SoftwareResampleContext {
	return newSoftwareResampleContextFromC(C.swr_alloc())
}

// https://ffmpeg.org/doxygen/8.0/group__lswr.html#ga818f7d78b1ad7d8d5b70de374b668c34
func (src *SoftwareResampleContext) Free() {
	if src.c != nil {
		// Make sure to clone the classer before freeing the object since
		// the C free method may reset the pointer
		c := newClonedClasser(src)
		C.swr_free(&src.c)
		// Make sure to remove from classers after freeing the object since
		// the C free method may use methods needing the classer
		if c != nil {
			classers.del(c)
		}
	}
}

var _ Classer = (*SoftwareResampleContext)(nil)

// https://ffmpeg.org/doxygen/8.0/structSwrContext.html#a7e13adcdcbc11bcc933cb7d0b9f839a0
func (src *SoftwareResampleContext) Class() *Class {
	if src.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(src.c))
}

// https://ffmpeg.org/doxygen/8.0/group__lswr.html#gac482028c01d95580106183aa84b0930c
func (src_ *SoftwareResampleContext) ConvertFrame(src, dst *Frame) error {
	var csrc *C.AVFrame
	if src != nil {
		csrc = src.c
	}
	return newError(C.swr_convert_frame(src_.c, dst.c, csrc))
}

// https://ffmpeg.org/doxygen/8.0/group__lswr.html#ga5121a5a7890a2d23b72dc871dd0ebb06
func (src_ *SoftwareResampleContext) Delay(base int64) int64 {
	return int64(C.swr_get_delay(src_.c, C.int64_t(base)))
}

// Init initializes the resampling context
// https://ffmpeg.org/doxygen/8.0/group__lswr.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (src *SoftwareResampleContext) Init() error {
	return newError(C.swr_init(src.c))
}

// AllocSetOpts2 allocates and sets options for a resampling context
// https://ffmpeg.org/doxygen/8.0/group__lswr.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func AllocSoftwareResampleContextSetOpts2(outChannelLayout ChannelLayout, outSampleFormat SampleFormat, outSampleRate int,
	inChannelLayout ChannelLayout, inSampleFormat SampleFormat, inSampleRate int) (*SoftwareResampleContext, error) {
	var c *C.SwrContext
	ret := C.swr_alloc_set_opts2(&c,
		outChannelLayout.c, C.enum_AVSampleFormat(outSampleFormat), C.int(outSampleRate),
		inChannelLayout.c, C.enum_AVSampleFormat(inSampleFormat), C.int(inSampleRate),
		0, nil)
	if err := newError(ret); err != nil {
		return nil, err
	}
	return newSoftwareResampleContextFromC(c), nil
}

// Convert converts audio samples using the resampler
// https://ffmpeg.org/doxygen/8.0/group__lswr.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (src *SoftwareResampleContext) Convert(out [][]byte, outCount int, in [][]byte, inCount int) (int, error) {
	var outPtr **C.uint8_t
	var inPtr **C.uint8_t
	
	if len(out) > 0 {
		// 创建C指针数组
		outPtrs := make([]*C.uint8_t, len(out))
		for i, data := range out {
			if len(data) > 0 {
				outPtrs[i] = (*C.uint8_t)(unsafe.Pointer(&data[0]))
			}
		}
		outPtr = (**C.uint8_t)(unsafe.Pointer(&outPtrs[0]))
	}
	
	if len(in) > 0 {
		// 创建C指针数组
		inPtrs := make([]*C.uint8_t, len(in))
		for i, data := range in {
			if len(data) > 0 {
				inPtrs[i] = (*C.uint8_t)(unsafe.Pointer(&data[0]))
			}
		}
		inPtr = (**C.uint8_t)(unsafe.Pointer(&inPtrs[0]))
	}
	
	ret := C.swr_convert(src.c, outPtr, C.int(outCount), inPtr, C.int(inCount))
	if ret < 0 {
		return 0, newError(ret)
	}
	return int(ret), nil
}
