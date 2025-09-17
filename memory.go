package astiav

//#include <libavutil/mem.h>
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

// Malloc allocates a block of size bytes with alignment suitable for all memory accesses
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func Malloc(size int) unsafe.Pointer {
	return unsafe.Pointer(C.av_malloc(C.size_t(size)))
}

// Mallocz allocates a block of size bytes with alignment suitable for all memory accesses and zero it
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func Mallocz(size int) unsafe.Pointer {
	return unsafe.Pointer(C.av_mallocz(C.size_t(size)))
}

// Realloc reallocates the given block if it is not large enough, otherwise does nothing
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func Realloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return unsafe.Pointer(C.av_realloc(ptr, C.size_t(size)))
}

// Free frees a memory block which has been allocated with a function of av_malloc() or av_realloc() family
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func Free(ptr unsafe.Pointer) {
	C.av_free(ptr)
}

// Freep frees a memory block which has been allocated with a function of av_malloc() or av_realloc() family, and set the pointer pointing to it to NULL
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func Freep(ptr *unsafe.Pointer) {
	C.av_freep(unsafe.Pointer(ptr))
}

// MallocArray allocates an array
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func MallocArray(nmemb, size int) unsafe.Pointer {
	return unsafe.Pointer(C.av_malloc_array(C.size_t(nmemb), C.size_t(size)))
}

// MalloczArray allocates an array and zero it (equivalent to av_calloc)
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func MalloczArray(nmemb, size int) unsafe.Pointer {
	return unsafe.Pointer(C.av_calloc(C.size_t(nmemb), C.size_t(size)))
}

// ReallocArray reallocates an array
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func ReallocArray(ptr unsafe.Pointer, nmemb, size int) unsafe.Pointer {
	return unsafe.Pointer(C.av_realloc_array(ptr, C.size_t(nmemb), C.size_t(size)))
}

// ReallocF reallocates the given block if it is not large enough, otherwise does nothing
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func ReallocF(ptr unsafe.Pointer, size int) unsafe.Pointer {
	return unsafe.Pointer(C.av_realloc_f(ptr, C.size_t(1), C.size_t(size)))
}

// Strdup duplicates a string
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga5a5b0b8b8b8b8b8b8b8b8b8b8b8b8b8b
func Strdup(s string) *C.char {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return C.av_strdup(cs)
}

// Strndup duplicates a substring of a string
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func Strndup(s string, len int) *C.char {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return C.av_strndup(cs, C.size_t(len))
}

// Memdup duplicates a buffer
// https://ffmpeg.org/doxygen/8.0/group__lavu__mem.html#ga7b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func Memdup(p unsafe.Pointer, size int) unsafe.Pointer {
	return unsafe.Pointer(C.av_memdup(p, C.size_t(size)))
}

// Note: av_fast_padded_malloc and av_fast_padded_mallocz are in libavcodec, not libavutil
// They should be implemented in codec_context.go or a separate codec utilities file