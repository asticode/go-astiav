package astiav

//#include <stdlib.h>
//#include <stdint.h>
import "C"
import (
	"errors"
	"unsafe"
)

func stringFromC(len int, fn func(buf *C.char, size C.size_t) error) (string, error) {
	size := C.size_t(len)
	buf := (*C.char)(C.malloc(size))
	if buf == nil {
		return "", errors.New("astiav: buf is nil")
	}
	defer C.free(unsafe.Pointer(buf))
	if err := fn(buf, size); err != nil {
		return "", err
	}
	return C.GoString(buf), nil
}

func bytesFromC(fn func(size *C.int) *C.uint8_t) []byte {
	var size int
	r := fn((*C.int)(unsafe.Pointer(&size)))
	return C.GoBytes(unsafe.Pointer(r), C.int(size))
}

func bytesToC(b []byte, fn func(b *C.uint8_t, size C.int) error) error {
	var ptr *C.uint8_t
	if b != nil {
		c := make([]byte, len(b))
		copy(c, b)
		ptr = (*C.uint8_t)(unsafe.Pointer(&c[0]))
	}
	return fn(ptr, C.int(len(b)))
}
