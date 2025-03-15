package astiav

//#include <libavcodec/avcodec.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVPacketSideData.html
// https://ffmpeg.org/doxygen/7.0/group__lavc__packet__side__data.html#ga9a80bfcacc586b483a973272800edb97
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

// https://ffmpeg.org/doxygen/7.0/group__lavc__packet__side__data.html#gga9a80bfcacc586b483a973272800edb97aab8c149a1e6c67aad340733becec87e1
func (d *PacketSideData) DisplayMatrix() *packetSideDataDisplayMatrix {
	return newPacketSideDataDisplayMatrix(d)
}

type packetSideDataDisplayMatrix struct {
	d *PacketSideData
}

func newPacketSideDataDisplayMatrix(d *PacketSideData) *packetSideDataDisplayMatrix {
	return &packetSideDataDisplayMatrix{d: d}
}

func (d *packetSideDataDisplayMatrix) Add(m *DisplayMatrix) error {
	return d.d.addBytes(C.AV_PKT_DATA_DISPLAYMATRIX, m.bytes())
}

func (d *packetSideDataDisplayMatrix) Get() (*DisplayMatrix, bool) {
	b := d.d.getBytes(C.AV_PKT_DATA_DISPLAYMATRIX)
	if len(b) == 0 {
		return nil, false
	}
	m, err := newDisplayMatrixFromBytes(b)
	if err != nil {
		return nil, false
	}
	return m, true
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__packet__side__data.html#gga9a80bfcacc586b483a973272800edb97a2093332d8086d25a04942ede61007f6a
func (d *PacketSideData) SkipSamples() *packetSideDataSkipSamples {
	return newPacketSideDataSkipSamples(d)
}

type packetSideDataSkipSamples struct {
	d *PacketSideData
}

func newPacketSideDataSkipSamples(d *PacketSideData) *packetSideDataSkipSamples {
	return &packetSideDataSkipSamples{d: d}
}

func (d *packetSideDataSkipSamples) Add(ss *SkipSamples) error {
	return d.d.addBytes(C.AV_PKT_DATA_SKIP_SAMPLES, ss.bytes())
}

func (d *packetSideDataSkipSamples) Get() (*SkipSamples, bool) {
	b := d.d.getBytes(C.AV_PKT_DATA_SKIP_SAMPLES)
	if len(b) == 0 {
		return nil, false
	}
	ss, err := newSkipSamplesFromBytes(b)
	if err != nil {
		return nil, false
	}
	return ss, true
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__packet__side__data.html#gad208a666db035802403ea994912a83db
func (d *PacketSideData) addBytes(t C.enum_AVPacketSideDataType, b []byte) error {
	if len(b) == 0 {
		return nil
	}

	sd := C.av_packet_side_data_new(d.sd, d.size, t, C.size_t(len(b)), 0)
	if sd == nil {
		return errors.New("astiav: nil pointer")
	}

	C.memcpy(unsafe.Pointer(sd.data), unsafe.Pointer(&b[0]), C.size_t(len(b)))
	return nil
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__packet__side__data.html#ga61a3a0fba92a308208c8ab957472d23c
func (d *PacketSideData) getBytes(t C.enum_AVPacketSideDataType) []byte {
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		if d.sd == nil || d.size == nil {
			return nil
		}
		sd := C.av_packet_side_data_get(*d.sd, *d.size, t)
		if sd == nil {
			return nil
		}
		*size = sd.size
		return sd.data
	})
}
