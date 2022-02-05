package astiav

//#cgo pkg-config: libavformat
//#include <libavformat/avformat.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L503
type OutputFormat struct {
	c *C.struct_AVOutputFormat
}

func newOutputFormatFromC(c *C.struct_AVOutputFormat) *OutputFormat {
	if c == nil {
		return nil
	}
	return &OutputFormat{c: c}
}

func (f *OutputFormat) Flags() IOFormatFlags {
	return IOFormatFlags(f.c.flags)
}
