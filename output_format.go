package astiav

//#include <libavformat/avformat.h>
import "C"
import "unsafe"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L503
type OutputFormat struct {
	c *C.AVOutputFormat
}

func newOutputFormatFromC(c *C.AVOutputFormat) *OutputFormat {
	if c == nil {
		return nil
	}
	return &OutputFormat{c: c}
}

func FindOutputFormat(name string) *OutputFormat {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return newOutputFormatFromC(C.av_guess_format(cname, nil, nil))
}

func (f *OutputFormat) Flags() IOFormatFlags {
	return IOFormatFlags(f.c.flags)
}

func (f *OutputFormat) Name() string {
	return C.GoString(f.c.name)
}

// LongName Description of the format, meant to be more human-readable than Name.
func (f *OutputFormat) LongName() string {
	return C.GoString(f.c.long_name)
}

func (f *OutputFormat) String() string {
	return f.Name()
}
