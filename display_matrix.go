package astiav

//#include <libavutil/display.h>
import "C"
import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

type DisplayMatrix [9]uint32

func NewDisplayMatrixFromBytes(b []byte) (m *DisplayMatrix, err error) {
	// Check length
	if len(b) < 36 {
		err = fmt.Errorf("astiav: invalid length %d < 36", len(b))
		return
	}

	// Create display matrix
	m = &DisplayMatrix{}

	// Loop
	for idx := 0; idx < 9; idx++ {
		m[idx] = binary.LittleEndian.Uint32(b[idx*4 : (idx+1)*4])
	}
	return
}

func NewDisplayMatrixFromRotation(angle float64) *DisplayMatrix {
	m := &DisplayMatrix{}
	C.av_display_rotation_set((*C.int32_t)(unsafe.Pointer(&m[0])), C.double(angle))
	return m
}

func (m DisplayMatrix) Bytes() []byte {
	b := make([]byte, 0, 36)
	for _, v := range m {
		b = binary.LittleEndian.AppendUint32(b, v)
	}
	return b
}

// Rotation is a clockwise angle
func (m DisplayMatrix) Rotation() float64 {
	return -float64(C.av_display_rotation_get((*C.int32_t)(unsafe.Pointer(&m[0]))))
}
