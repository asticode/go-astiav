package astiav

//#include <libswresample/swresample.h>
import "C"
import "unsafe"

// https://ffmpeg.org/doxygen/7.0/structSwrContext.html
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

// https://ffmpeg.org/doxygen/7.0/group__lswr.html#gaf58c4ff10f73d74bdab8e5aa7193147c
func AllocSoftwareResampleContext() *SoftwareResampleContext {
	return newSoftwareResampleContextFromC(C.swr_alloc())
}

// https://ffmpeg.org/doxygen/7.0/group__lswr.html#ga818f7d78b1ad7d8d5b70de374b668c34
func (src *SoftwareResampleContext) Free() {
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

var _ Classer = (*SoftwareResampleContext)(nil)

// https://ffmpeg.org/doxygen/7.0/structSwrContext.html#a7e13adcdcbc11bcc933cb7d0b9f839a0
func (src *SoftwareResampleContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(src.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lswr.html#gac482028c01d95580106183aa84b0930c
func (src_ *SoftwareResampleContext) ConvertFrame(src, dst *Frame) error {
	var csrc *C.AVFrame
	if src != nil {
		csrc = src.c
	}
	return newError(C.swr_convert_frame(src_.c, dst.c, csrc))
}

// https://ffmpeg.org/doxygen/7.0/group__lswr.html#ga5121a5a7890a2d23b72dc871dd0ebb06
func (src_ *SoftwareResampleContext) Delay(base int64) int64 {
	return int64(C.swr_get_delay(src_.c, C.int64_t(base)))
}
