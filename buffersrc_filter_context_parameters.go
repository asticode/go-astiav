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
