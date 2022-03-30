package astiav

//#cgo pkg-config: libavcodec
//#include <libavcodec/avcodec.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/packet.h#L350
type Packet struct {
	c *C.struct_AVPacket
}

func newPacketFromC(c *C.struct_AVPacket) *Packet {
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
	return bytesFromC(func(size *C.int) *C.uint8_t {
		*size = p.c.size
		return p.c.data
	})
}

func (p *Packet) SetData(data []byte) {
	bytesToC(data, func(b *C.uint8_t, size C.int) error {
		p.c.data = b
		p.c.size = size
		return nil
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

func (p *Packet) SideData(t PacketSideDataType) []byte {
	return bytesFromC(func(size *C.int) *C.uint8_t {
		return C.av_packet_get_side_data(p.c, (C.enum_AVPacketSideDataType)(t), size)
	})
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
