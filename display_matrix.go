package astiav

//#cgo pkg-config: libavutil
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

func (m DisplayMatrix) Rotation() float64 {
	return float64(C.av_display_rotation_get((*C.int32_t)(unsafe.Pointer(&m[0]))))
}
