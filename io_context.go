package astiav

//#cgo pkg-config: libavformat
//#include <libavformat/avformat.h>
import "C"
import "unsafe"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avio.h#L161
type IOContext struct {
	c *C.struct_AVIOContext
}

func NewIOContext() *IOContext {
	return &IOContext{}
}

func newIOContextFromC(c *C.struct_AVIOContext) *IOContext {
	if c == nil {
		return nil
	}
	return &IOContext{c: c}
}

func (ic *IOContext) Closep() error {
	return newError(C.avio_closep(&ic.c))
}

func (ic *IOContext) Open(filename string, flags IOContextFlags) error {
	cfi := C.CString(filename)
	defer C.free(unsafe.Pointer(cfi))
	return newError(C.avio_open(&ic.c, cfi, C.int(flags)))
}

func (ic *IOContext) Write(b []byte) {
	if b == nil {
		return
	}
	C.avio_write(ic.c, (*C.uchar)(unsafe.Pointer(&b[0])), C.int(len(b)))
}
