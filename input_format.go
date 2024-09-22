package astiav

//#include <libavformat/avformat.h>
import "C"
import "unsafe"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L650
type InputFormat struct {
	c *C.AVInputFormat
}

func newInputFormatFromC(c *C.AVInputFormat) *InputFormat {
	if c == nil {
		return nil
	}
	return &InputFormat{c: c}
}

func FindInputFormat(name string) *InputFormat {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return newInputFormatFromC(C.av_find_input_format(cname))
}

func (f *InputFormat) Flags() IOFormatFlags {
	return IOFormatFlags(f.c.flags)
}

func (f *InputFormat) Name() string {
	return C.GoString(f.c.name)
}

// LongName Description of the format, meant to be more human-readable than Name.
func (f *InputFormat) LongName() string {
	return C.GoString(f.c.long_name)
}

func (f *InputFormat) String() string {
	return f.Name()
}
