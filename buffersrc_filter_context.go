package astiav

//#include <libavfilter/buffersrc.h>
import "C"

type BuffersrcFilterContext struct {
	fc *FilterContext
}

func newBuffersrcFilterContext(fc *FilterContext) *BuffersrcFilterContext {
	return &BuffersrcFilterContext{fc: fc}
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersrc.html#ga73ed90c3c3407f36e54d65f91faaaed9
func (bfc *BuffersrcFilterContext) AddFrame(f *Frame, fs BuffersrcFlags) error {
	var cf *C.AVFrame
	if f != nil {
		cf = f.c
	}
	return newError(C.av_buffersrc_add_frame_flags(bfc.fc.c, cf, C.int(fs)))
}

func (bfc *BuffersrcFilterContext) FilterContext() *FilterContext {
	return bfc.fc
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#ga8c15af28902395399fe455f6f8236848
func (bfc *BuffersrcFilterContext) Initialize() error {
	return newError(C.avfilter_init_dict(bfc.fc.c, nil))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi__buffersrc.html#ga398cd2a84f8b4a588197ab9d90135048
func (bfc *BuffersrcFilterContext) SetParameters(bfcp *BuffersrcFilterContextParameters) error {
	return newError(C.av_buffersrc_parameters_set(bfc.fc.c, bfcp.c))
}
