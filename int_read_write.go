package astiav

//#include "int_read_write.h"
import "C"
import "unsafe"

// https://ffmpeg.org/doxygen/7.0/avr32_2intreadwrite_8h.html#ace46e41b9bd6cac88fb7109ffd657f9a
func RL32(i []byte) uint32 {
	if len(i) == 0 {
		return 0
	}
	return uint32(C.astiavRL32((*C.uint8_t)(unsafe.Pointer(&i[0]))))
}

// https://ffmpeg.org/doxygen/7.0/avr32_2intreadwrite_8h.html#ace46e41b9bd6cac88fb7109ffd657f9a
func RL32WithOffset(i []byte, offset uint) uint32 {
	if len(i) == 0 {
		return 0
	}
	return uint32(C.astiavRL32WithOffset((*C.uint8_t)(unsafe.Pointer(&i[0])), C.int(offset)))
}
