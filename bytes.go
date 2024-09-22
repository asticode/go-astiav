package astiav

//#include <libavcodec/avcodec.h>
//#include <stdlib.h>
//#include <stdint.h>
import "C"
import (
	"errors"
	"fmt"
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

func bytesFromC(fn func(size *C.size_t) *C.uint8_t) []byte {
	var size uint64
	r := fn((*C.size_t)(unsafe.Pointer(&size)))
	return C.GoBytes(unsafe.Pointer(r), C.int(size))
}

func bytesToC(b []byte, fn func(b *C.uint8_t, size C.size_t) error) error {
	var ptr *C.uint8_t
	if b != nil {
		c := make([]byte, len(b))
		copy(c, b)
		ptr = (*C.uint8_t)(unsafe.Pointer(&c[0]))
	}
	return fn(ptr, C.size_t(len(b)))
}

// TODO Rename?
func setBytesWithIntSizeInC(b []byte, d **C.uint8_t, size *C.int) error {
	if len(b) == 0 {
		return nil
	}

	if *d != nil {
		C.av_freep(unsafe.Pointer(d))
		*size = 0
	}

	if *d = (*C.uint8_t)(C.av_mallocz(C.size_t(len(b) + C.AV_INPUT_BUFFER_PADDING_SIZE))); *d == nil {
		return fmt.Errorf("astiav: allocation is nil")
	}

	C.memcpy(unsafe.Pointer(*d), unsafe.Pointer(&b[0]), C.size_t(len(b)))
	*size = C.int(len(b))
	return nil
}
