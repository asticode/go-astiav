package astiav

//#include <libavcodec/avcodec.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/packet.h#L350
type Packet struct {
	c *C.AVPacket
}

func newPacketFromC(c *C.AVPacket) *Packet {
	if c == nil {
		return nil
	}
	return &Packet{c: c}
}

func AllocPacket() *Packet {
	return newPacketFromC(C.av_packet_alloc())
}

func (p *Packet) Data() []byte {
	if p.c.data == nil {
		return nil
	}
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		*size = C.size_t(p.c.size)
		return p.c.data
	})
}

func (p *Packet) Dts() int64 {
	return int64(p.c.dts)
}

func (p *Packet) SetDts(v int64) {
	p.c.dts = C.int64_t(v)
}

func (p *Packet) Duration() int64 {
	return int64(p.c.duration)
}

func (p *Packet) SetDuration(d int64) {
	p.c.duration = C.int64_t(d)
}

func (p *Packet) Flags() PacketFlags {
	return PacketFlags(p.c.flags)
}

func (p *Packet) SetFlags(f PacketFlags) {
	p.c.flags = C.int(f)
}

func (p *Packet) Pos() int64 {
	return int64(p.c.pos)
}

func (p *Packet) SetPos(v int64) {
	p.c.pos = C.int64_t(v)
}

func (p *Packet) Pts() int64 {
	return int64(p.c.pts)
}

func (p *Packet) SetPts(v int64) {
	p.c.pts = C.int64_t(v)
}

func (p *Packet) SideData() *PacketSideData {
	return newPacketSideDataFromC(&p.c.side_data, &p.c.side_data_elems)
}

func (p *Packet) Size() int {
	return int(p.c.size)
}

func (p *Packet) SetSize(s int) {
	p.c.size = C.int(s)
}

func (p *Packet) StreamIndex() int {
	return int(p.c.stream_index)
}

func (p *Packet) SetStreamIndex(i int) {
	p.c.stream_index = C.int(i)
}

func (p *Packet) Free() {
	C.av_packet_free(&p.c)
}

func (p *Packet) Clone() *Packet {
	return newPacketFromC(C.av_packet_clone(p.c))
}

func (p *Packet) AllocPayload(s int) error {
	return newError(C.av_new_packet(p.c, C.int(s)))
}

func (p *Packet) Ref(src *Packet) error {
	return newError(C.av_packet_ref(p.c, src.c))
}

func (p *Packet) Unref() {
	C.av_packet_unref(p.c)
}

func (p *Packet) MoveRef(src *Packet) {
	C.av_packet_move_ref(p.c, src.c)
}

func (p *Packet) RescaleTs(src, dst Rational) {
	C.av_packet_rescale_ts(p.c, src.c, dst.c)
}

func (p *Packet) FromData(data []byte) (err error) {
	// Create buf
	buf := (*C.uint8_t)(C.av_malloc(C.size_t(len(data))))
	if buf == nil {
		err = errors.New("astiav: allocating buffer failed")
		return
	}

	// Make sure to free buf in case of error
	defer func() {
		if err != nil {
			C.av_freep(unsafe.Pointer(&buf))
		}
	}()

	// Copy
	if len(data) > 0 {
		C.memcpy(unsafe.Pointer(buf), unsafe.Pointer(&data[0]), C.size_t(len(data)))
	}

	// From data
	err = newError(C.av_packet_from_data(p.c, buf, C.int(len(data))))
	return
}
