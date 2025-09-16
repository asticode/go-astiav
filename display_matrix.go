package astiav

//#include <libavutil/display.h>
import "C"
import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

// https://ffmpeg.org/doxygen/8.1/group__lavu__video__display.html
type DisplayMatrix [9]uint32

func newDisplayMatrixFromBytes(b []byte) (m *DisplayMatrix, err error) {
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

// https://ffmpeg.org/doxygen/8.1/group__lavu__video__display.html#ga5964303bfe085ad33683bc2454768d4a
func NewDisplayMatrixFromRotation(angle float64) *DisplayMatrix {
	m := &DisplayMatrix{}
	C.av_display_rotation_set((*C.int32_t)(unsafe.Pointer(&m[0])), C.double(angle))
	return m
}

func (m DisplayMatrix) bytes() []byte {
	b := make([]byte, 0, 36)
	for _, v := range m {
		b = binary.LittleEndian.AppendUint32(b, v)
	}
	return b
}

// Rotation is a clockwise angle
// https://ffmpeg.org/doxygen/8.1/group__lavu__video__display.html#gaac2ea94d3f66496c758349450b5b0217
func (m DisplayMatrix) Rotation() float64 {
	return -float64(C.av_display_rotation_get((*C.int32_t)(unsafe.Pointer(&m[0]))))
}
