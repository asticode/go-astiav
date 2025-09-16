package astiav

//#include <libavformat/avformat.h>
import "C"
import "unsafe"

// https://ffmpeg.org/doxygen/8.1/structAVInputFormat.html
type InputFormat struct {
	c *C.AVInputFormat
}

func newInputFormatFromC(c *C.AVInputFormat) *InputFormat {
	if c == nil {
		return nil
	}
	return &InputFormat{c: c}
}

// https://ffmpeg.org/doxygen/8.1/group__lavf__decoding.html#ga40034b6d64d372e1c989e16dde4b459a
func FindInputFormat(name string) *InputFormat {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return newInputFormatFromC(C.av_find_input_format(cname))
}

// https://ffmpeg.org/doxygen/8.1/structAVInputFormat.html#a1b30f6647d0c2faf38ba8786d7c3a838
func (f *InputFormat) Flags() IOFormatFlags {
	return IOFormatFlags(f.c.flags)
}

// https://ffmpeg.org/doxygen/8.1/structAVInputFormat.html#a850db3eb225e22b64f3304d72134ca0c
func (f *InputFormat) Name() string {
	return C.GoString(f.c.name)
}

// https://ffmpeg.org/doxygen/8.1/structAVInputFormat.html#a1f67064a527941944017f1dfe65d3aa9
func (f *InputFormat) LongName() string {
	return C.GoString(f.c.long_name)
}

func (f *InputFormat) String() string {
	return f.Name()
}
