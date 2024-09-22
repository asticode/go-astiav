package astiav

//#include "int_read_write.h"
import "C"
import "unsafe"

func RL32(i []byte) uint32 {
	if len(i) == 0 {
		return 0
	}
	return uint32(C.astiavRL32((*C.uint8_t)(unsafe.Pointer(&i[0]))))
}

func RL32WithOffset(i []byte, offset uint) uint32 {
	if len(i) == 0 {
		return 0
	}
	return uint32(C.astiavRL32WithOffset((*C.uint8_t)(unsafe.Pointer(&i[0])), C.int(offset)))
}
