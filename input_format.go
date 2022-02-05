package astiav

//#cgo pkg-config: libavformat
//#include <libavformat/avformat.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L650
type InputFormat struct {
	c *C.struct_AVInputFormat
}

func newInputFormatFromC(c *C.struct_AVInputFormat) *InputFormat {
	if c == nil {
		return nil
	}
	return &InputFormat{c: c}
}

func (f *InputFormat) Flags() IOFormatFlags {
	return IOFormatFlags(f.c.flags)
}
