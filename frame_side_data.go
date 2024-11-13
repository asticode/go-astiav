package astiav

//#include <libavutil/frame.h>
import "C"
import (
	"math"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVFrameSideData.html
type FrameSideData struct {
	c *C.AVFrameSideData
}

func newFrameSideDataFromC(c *C.AVFrameSideData) *FrameSideData {
	if c == nil {
		return nil
	}
	return &FrameSideData{c: c}
}

// https://ffmpeg.org/doxygen/7.0/structAVFrameSideData.html#a76937ad48652a5a0cc4bff65fc6c886e
func (d *FrameSideData) Data() []byte {
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		*size = d.c.size
		return d.c.data
	})
}

// https://ffmpeg.org/doxygen/7.0/structAVFrameSideData.html#a76937ad48652a5a0cc4bff65fc6c886e
func (d *FrameSideData) SetData(b []byte) {
	C.memcpy(unsafe.Pointer(d.c.data), unsafe.Pointer(&b[0]), C.size_t(math.Min(float64(len(b)), float64(d.c.size))))
}

// https://ffmpeg.org/doxygen/7.0/structAVFrameSideData.html#a07ff3499827c124591ff4bae6f68eec0
func (d *FrameSideData) Type() FrameSideDataType {
	return FrameSideDataType(d.c._type)
}
