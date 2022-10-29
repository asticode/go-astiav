package astiav

//#cgo pkg-config: libavformat
//#include <libavformat/avformat.h>
import "C"
import (
	"io"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avio.h#L161
type IOContext struct {
	c *C.struct_AVIOContext
	io.ReadWriteCloser
	io.Seeker
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

func (ic *IOContext) Close() error {
	return newError(C.avio_close(ic.c))
}

func (ic *IOContext) Open(filename string, flags IOContextFlags) error {
	cfi := C.CString(filename)
	defer C.free(unsafe.Pointer(cfi))
	return newError(C.avio_open(&ic.c, cfi, C.int(flags)))
}

func (ic *IOContext) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	C.avio_write(ic.c, (*C.uchar)(unsafe.Pointer(&b[0])), C.int(len(b)))
	return len(b), nil
}

func (ic *IOContext) Read(b []byte) (n int, err error) {
	n = int(C.avio_read(ic.c, (*C.uchar)(unsafe.Pointer(&b[0])), C.int(len(b))))
	if n < 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (ic *IOContext) Seekable() bool {
	return int(ic.c.seekable) != 0
}

func (ic *IOContext) Seek(offset int64, whence int) (int64, error) {
	return int64(C.avio_seek(ic.c, C.int64_t(offset), C.int(whence))), nil
}
