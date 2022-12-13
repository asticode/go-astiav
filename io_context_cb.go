package astiav

/*
#cgo pkg-config: libavformat
#include <libavformat/avformat.h>
*/
import "C"
import (
	"unsafe"
)

type ioContextCbs struct {
	readCb  IOContextReadFunc
	writeCb IOContextWriteFunc
	seekCb  IOContextSeekFunc
}

//export go_ioctx_proxy_read
func go_ioctx_proxy_read(opaque unsafe.Pointer, buf *C.uint8_t, buf_size C.int) C.int {
	if buf == nil || buf_size == 0 {
		return C.int(ErrEio)
	}
	id := int(*(*C.int)(opaque))
	if ctx, ok := fetchIOCallback(id); ok && ctx.readCb != nil {
		n := ctx.readCb((*[1 << 30]byte)(unsafe.Pointer(buf))[:int(buf_size)])
		cn := C.int(n)
		if n < 0 {
			return C.int(cn)
		} else if n == 0 {
			return C.int(ErrUnknown)
		}
		return cn
	}
	return C.int(ErrEio)
}

//export go_ioctx_proxy_write
func go_ioctx_proxy_write(opaque unsafe.Pointer, buf *C.uint8_t, buf_size C.int) C.int {
	if buf == nil || buf_size == 0 {
		return C.int(ErrEio)
	}
	id := int(*(*C.int)(opaque))
	if ctx, ok := fetchIOCallback(id); ok && ctx.writeCb != nil {
		return C.int(ctx.writeCb((*[1 << 30]byte)(unsafe.Pointer(buf))[:int(buf_size)]))
	}
	return C.int(ErrEio)
}

//export go_ioctx_proxy_seek
func go_ioctx_proxy_seek(opaque unsafe.Pointer, offset C.int64_t, whence C.int) C.int64_t {
	id := int(*(*C.int)(opaque))
	if ctx, ok := fetchIOCallback(id); ok && ctx.seekCb != nil {
		return C.int64_t(ctx.seekCb(int64(offset), int(whence)))
	}
	return C.int64_t(-1)
}
