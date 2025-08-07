package astiav

//#include <libavformat/avformat.h>
import "C"
import (
	"strings"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVOutputFormat.html
type OutputFormat struct {
	c *C.AVOutputFormat
}

func newOutputFormatFromC(c *C.AVOutputFormat) *OutputFormat {
	if c == nil {
		return nil
	}
	return &OutputFormat{c: c}
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__encoding.html#ga00bceb049f2b20716e2f36ebc990a350
func FindOutputFormat(name string) *OutputFormat {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return newOutputFormatFromC(C.av_guess_format(cname, nil, nil))
}

// https://ffmpeg.org/doxygen/7.0/structAVOutputFormat.html#aad55a00e728a020c1dcfaaf695320445
func (f *OutputFormat) Flags() IOFormatFlags {
	return IOFormatFlags(f.c.flags)
}

// https://ffmpeg.org/doxygen/7.0/structAVOutputFormat.html#ac3abc5f47f3465b6b7eec89c9476351c
func (f *OutputFormat) Name() string {
	return C.GoString(f.c.name)
}

// https://ffmpeg.org/doxygen/7.0/structAVOutputFormat.html#a4ff98d90aac0047a204a35a758a363fc
func (f *OutputFormat) LongName() string {
	return C.GoString(f.c.long_name)
}

func (f *OutputFormat) String() string {
	return f.Name()
}

// https://ffmpeg.org/doxygen/7.0/structAVOutputFormat.html#a10f19abe463890063659723c90c15335
func (f *OutputFormat) Extensions() []string {
	s := C.GoString(f.c.extensions)
	if s != "" {
		return strings.Split(s, ",")
	}
	return nil
}
