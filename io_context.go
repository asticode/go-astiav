package astiav

//#cgo pkg-config: libavformat
//#include <libavformat/avformat.h>
import "C"
import "unsafe"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avio.h#L161
type IOContext struct {
	c *C.struct_AVIOContext
}

func newIOContextFromC(c *C.struct_AVIOContext) *IOContext {
	if c == nil {
		return nil
	}
	ic := &IOContext{c: c}
	classers.set(ic)
	return ic
}

var _ Classer = (*IOContext)(nil)

func OpenIOContext(filename string, flags IOContextFlags) (*IOContext, error) {
	cfi := C.CString(filename)
	defer C.free(unsafe.Pointer(cfi))
	var c *C.struct_AVIOContext
	if err := newError(C.avio_open(&c, cfi, C.int(flags))); err != nil {
		return nil, err
	}
	return newIOContextFromC(c), nil
}

func (ic *IOContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(ic.c))
}

func (ic *IOContext) Closep() error {
	classers.del(ic)
	if ic.c != nil {
		return newError(C.avio_closep(&ic.c))
	}
	return nil
}

func (ic *IOContext) Write(b []byte) {
	if b == nil {
		return
	}
	C.avio_write(ic.c, (*C.uchar)(unsafe.Pointer(&b[0])), C.int(len(b)))
}
