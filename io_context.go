package astiav

//#include <libavformat/avformat.h>
//#include "io_context.h"
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avio.h#L161
type IOContext struct {
	c         *C.AVIOContext
	handlerID unsafe.Pointer
}

func newIOContextFromC(c *C.AVIOContext) *IOContext {
	if c == nil {
		return nil
	}
	ic := &IOContext{c: c}
	classers.set(ic)
	return ic
}

var _ Classer = (*IOContext)(nil)

type IOContextReadFunc func(b []byte) (n int, err error)

type IOContextSeekFunc func(offset int64, whence int) (n int64, err error)

type IOContextWriteFunc func(b []byte) (n int, err error)

func AllocIOContext(bufferSize int, writable bool, readFunc IOContextReadFunc, seekFunc IOContextSeekFunc, writeFunc IOContextWriteFunc) (ic *IOContext, err error) {
	// Invalid buffer size
	if bufferSize <= 0 {
		err = errors.New("astiav: buffer size <= 0")
		return
	}

	// Alloc buffer
	buffer := C.av_malloc(C.size_t(bufferSize))
	if buffer == nil {
		err = errors.New("astiav: allocating buffer failed")
		return
	}

	// Make sure buffer is freed in case of error
	defer func() {
		if err != nil {
			C.av_free(buffer)
		}
	}()

	// Since go doesn't allow c to store pointers to go data, we need to create this C pointer
	handlerID := C.malloc(C.size_t(1))
	if handlerID == nil {
		err = errors.New("astiav: allocating handler id failed")
		return
	}

	// Make sure handler id is freed in case of error
	defer func() {
		if err != nil {
			C.free(handlerID)
		}
	}()

	// Get callbacks
	var cReadFunc, cSeekFunc, cWriteFunc *[0]byte
	if readFunc != nil {
		cReadFunc = (*[0]byte)(C.astiavIOContextReadFunc)
	}
	if seekFunc != nil {
		cSeekFunc = (*[0]byte)(C.astiavIOContextSeekFunc)
	}
	if writeFunc != nil {
		cWriteFunc = (*[0]byte)(C.astiavIOContextWriteFunc)
	}

	// Get write flag
	wf := C.int(0)
	if writable {
		wf = C.int(1)
	}

	// Alloc io context
	cic := C.avio_alloc_context((*C.uchar)(buffer), C.int(bufferSize), wf, handlerID, cReadFunc, cWriteFunc, cSeekFunc)
	if cic == nil {
		err = errors.New("astiav: allocating io context failed: %w")
		return
	}

	// Create io context
	ic = newIOContextFromC(cic)

	// Store handler
	ic.handlerID = handlerID
	ioContextHandlers.set(handlerID, &ioContextHandler{
		r: readFunc,
		s: seekFunc,
		w: writeFunc,
	})
	return
}

func OpenIOContext(filename string, flags IOContextFlags) (*IOContext, error) {
	cfi := C.CString(filename)
	defer C.free(unsafe.Pointer(cfi))
	var c *C.AVIOContext
	if err := newError(C.avio_open(&c, cfi, C.int(flags))); err != nil {
		return nil, err
	}
	return newIOContextFromC(c), nil
}

func (ic *IOContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(ic.c))
}

func (ic *IOContext) Close() error {
	classers.del(ic)
	if ic.c != nil {
		return newError(C.avio_closep(&ic.c))
	}
	return nil
}

func (ic *IOContext) Free() {
	classers.del(ic)
	if ic.c != nil {
		if ic.c.buffer != nil {
			C.av_freep(unsafe.Pointer(&ic.c.buffer))
		}
		if ic.handlerID != nil {
			C.free(ic.handlerID)
			ic.handlerID = nil
		}
		C.avio_context_free(&ic.c)
	}
	return
}

func (ic *IOContext) Read(b []byte) (n int, err error) {
	// Nothing to read
	if b == nil || len(b) <= 0 {
		return
	}

	// Alloc buffer
	buf := C.av_malloc(C.size_t(len(b)))
	if buf == nil {
		err = errors.New("astiav: allocating buffer failed")
		return
	}

	// Make sure buffer is freed
	defer C.av_free(buf)

	// Read
	ret := C.avio_read_partial(ic.c, (*C.uchar)(unsafe.Pointer(buf)), C.int(len(b)))
	if err = newError(ret); err != nil {
		err = fmt.Errorf("astiav: reading failed: %w", err)
		return
	}

	// Copy
	C.memcpy(unsafe.Pointer(&b[0]), unsafe.Pointer(buf), C.size_t(ret))
	n = int(ret)
	return
}

func (ic *IOContext) Seek(offset int64, whence int) (int64, error) {
	ret := C.avio_seek(ic.c, C.int64_t(offset), C.int(whence))
	if err := newError(C.int(ret)); err != nil {
		return 0, err
	}
	return int64(ret), nil
}

func (ic *IOContext) Write(b []byte) {
	// Nothing to write
	if b == nil || len(b) <= 0 {
		return
	}

	// Write
	C.avio_write(ic.c, (*C.uchar)(unsafe.Pointer(&b[0])), C.int(len(b)))
}

func (ic *IOContext) Flush() {
	C.avio_flush(ic.c)
}

type ioContextHandler struct {
	r IOContextReadFunc
	s IOContextSeekFunc
	w IOContextWriteFunc
}

var ioContextHandlers = newIOContextHandlerPool()

type ioContextHandlerPool struct {
	m sync.Mutex
	p map[unsafe.Pointer]*ioContextHandler
}

func newIOContextHandlerPool() *ioContextHandlerPool {
	return &ioContextHandlerPool{p: make(map[unsafe.Pointer]*ioContextHandler)}
}

func (p *ioContextHandlerPool) set(id unsafe.Pointer, h *ioContextHandler) {
	p.m.Lock()
	defer p.m.Unlock()
	p.p[id] = h
}

func (p *ioContextHandlerPool) get(id unsafe.Pointer) (h *ioContextHandler, ok bool) {
	p.m.Lock()
	defer p.m.Unlock()
	h, ok = p.p[id]
	return
}

//export goAstiavIOContextReadFunc
func goAstiavIOContextReadFunc(opaque unsafe.Pointer, buf *C.uint8_t, bufSize C.int) C.int {
	// Get handler
	h, ok := ioContextHandlers.get(opaque)
	if !ok {
		return C.AVERROR_UNKNOWN
	}

	// Create go buffer
	b := make([]byte, int(bufSize), int(bufSize))

	// Read
	n, err := h.r(b)
	if err != nil {
		var e Error
		if errors.As(err, &e) {
			return C.int(e)
		}
		return C.AVERROR_UNKNOWN
	}

	// Copy
	C.memcpy(unsafe.Pointer(buf), unsafe.Pointer(&b[0]), C.size_t(n))
	return C.int(n)
}

//export goAstiavIOContextSeekFunc
func goAstiavIOContextSeekFunc(opaque unsafe.Pointer, offset C.int64_t, whence C.int) C.int64_t {
	// Get handler
	h, ok := ioContextHandlers.get(opaque)
	if !ok {
		return C.AVERROR_UNKNOWN
	}

	// Seek
	n, err := h.s(int64(offset), int(whence))
	if err != nil {
		var e Error
		if errors.As(err, &e) {
			return C.int64_t(e)
		}
		return C.int64_t(C.AVERROR_UNKNOWN)
	}
	return C.int64_t(n)
}

//export goAstiavIOContextWriteFunc
func goAstiavIOContextWriteFunc(opaque unsafe.Pointer, buf *C.uint8_t, bufSize C.int) C.int {
	// Get handler
	h, ok := ioContextHandlers.get(opaque)
	if !ok {
		return C.AVERROR_UNKNOWN
	}

	// Write
	n, err := h.w(C.GoBytes(unsafe.Pointer(buf), bufSize))
	if err != nil {
		var e Error
		if errors.As(err, &e) {
			return C.int(e)
		}
		return C.AVERROR_UNKNOWN
	}
	return C.int(n)
}
