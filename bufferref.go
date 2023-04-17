package astiav

// #cgo pkg-config: libavutil
// #include <libavutil/buffer.h>
import "C"

type BufferRef struct {
	c *C.struct_AVBufferRef
}

func newBufferRef() *BufferRef {
	return &BufferRef{}
}

func newBufferFromC(c *C.struct_AVBufferRef) *BufferRef {
	if c == nil {
		return nil
	}
	return &BufferRef{c: c}
}

func (br *BufferRef) Ref() *BufferRef {
	c := C.av_buffer_ref(br.c)
	return newBufferFromC(c)
}

func (br *BufferRef) Unref() {
	C.av_buffer_unref(&br.c)
}
