package astiav

//#include <libavcodec/avcodec.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n6.1.1/libavcodec/packet.h#L342
type PacketSideData struct {
	sd   **C.AVPacketSideData
	size *C.int
}

func newPacketSideDataFromC(sd **C.AVPacketSideData, size *C.int) *PacketSideData {
	return &PacketSideData{
		sd:   sd,
		size: size,
	}
}

func (d *PacketSideData) Add(t PacketSideDataType, b []byte) error {
	if len(b) == 0 {
		return nil
	}

	sd := C.av_packet_side_data_new(d.sd, d.size, (C.enum_AVPacketSideDataType)(t), C.size_t(len(b)), 0)
	if sd == nil {
		return errors.New("astiav: nil pointer")
	}

	C.memcpy(unsafe.Pointer(sd.data), unsafe.Pointer(&b[0]), C.size_t(len(b)))
	return nil
}

func (d *PacketSideData) Get(t PacketSideDataType) []byte {
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		if d.sd == nil || d.size == nil {
			return nil
		}
		sd := C.av_packet_side_data_get(*d.sd, *d.size, (C.enum_AVPacketSideDataType)(t))
		if sd == nil {
			return nil
		}
		*size = sd.size
		return sd.data
	})
}
