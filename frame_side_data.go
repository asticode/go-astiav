package astiav

//#include <libavutil/frame.h>
import "C"
import (
	"math"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/frame.h#L223
type FrameSideData struct {
	c *C.AVFrameSideData
}

func newFrameSideDataFromC(c *C.AVFrameSideData) *FrameSideData {
	if c == nil {
		return nil
	}
	return &FrameSideData{c: c}
}

func (d *FrameSideData) Data() []byte {
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		*size = d.c.size
		return d.c.data
	})
}

func (d *FrameSideData) SetData(b []byte) {
	C.memcpy(unsafe.Pointer(d.c.data), unsafe.Pointer(&b[0]), C.size_t(math.Min(float64(len(b)), float64(d.c.size))))
}

func (d *FrameSideData) Type() FrameSideDataType {
	return FrameSideDataType(d.c._type)
}
