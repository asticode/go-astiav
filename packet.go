package astiav

//#include <libavcodec/avcodec.h>
import "C"
import (
	"errors"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html
type Packet struct {
	c *C.AVPacket
}

func newPacketFromC(c *C.AVPacket) *Packet {
	if c == nil {
		return nil
	}
	return &Packet{c: c}
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#gaaf85aa950695631e0217a16062289b66
func AllocPacket() *Packet {
	return newPacketFromC(C.av_packet_alloc())
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#aaf4fe58dfcc7c232c1f2268b539d8367
func (p *Packet) Data() []byte {
	if p.c.data == nil {
		return nil
	}
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		*size = C.size_t(p.c.size)
		return p.c.data
	})
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a85dbbd306b44b02390cd91c45e6a0f76
func (p *Packet) Dts() int64 {
	return int64(p.c.dts)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a85dbbd306b44b02390cd91c45e6a0f76
func (p *Packet) SetDts(v int64) {
	p.c.dts = C.int64_t(v)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a622e758be29fd500aed0ffdc069550f7
func (p *Packet) Duration() int64 {
	return int64(p.c.duration)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a622e758be29fd500aed0ffdc069550f7
func (p *Packet) SetDuration(d int64) {
	p.c.duration = C.int64_t(d)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a437be96a9da675f12caa228a9c81bd82
func (p *Packet) Flags() PacketFlags {
	return PacketFlags(p.c.flags)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a437be96a9da675f12caa228a9c81bd82
func (p *Packet) SetFlags(f PacketFlags) {
	p.c.flags = C.int(f)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#ab5793d8195cf4789dfb3913b7a693903
func (p *Packet) Pos() int64 {
	return int64(p.c.pos)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#ab5793d8195cf4789dfb3913b7a693903
func (p *Packet) SetPos(v int64) {
	p.c.pos = C.int64_t(v)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a73bde0a37f3b1efc839f11295bfbf42a
func (p *Packet) Pts() int64 {
	return int64(p.c.pts)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a73bde0a37f3b1efc839f11295bfbf42a
func (p *Packet) SetPts(v int64) {
	p.c.pts = C.int64_t(v)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#ac55bfef91c33f02704ba76518d0f294c
func (p *Packet) SideData() *PacketSideData {
	return newPacketSideDataFromC(&p.c.side_data, &p.c.side_data_elems)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a4d1ea19f63eb107111fd650ca514d1f4
func (p *Packet) Size() int {
	return int(p.c.size)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a4d1ea19f63eb107111fd650ca514d1f4
func (p *Packet) SetSize(s int) {
	p.c.size = C.int(s)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a0d1cb9b5a32b00fb6edc81ea3aae2a49
func (p *Packet) StreamIndex() int {
	return int(p.c.stream_index)
}

// https://ffmpeg.org/doxygen/7.1/structAVPacket.html#a0d1cb9b5a32b00fb6edc81ea3aae2a49
func (p *Packet) SetStreamIndex(i int) {
	p.c.stream_index = C.int(i)
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#ga1066464e7cdd1f215df6940db94e5d8e
func (p *Packet) Free() {
	if p.c != nil {
		C.av_packet_free(&p.c)
	}
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#gacbe3e51cf411a7003d706127dc48cbb1
func (p *Packet) Clone() *Packet {
	return newPacketFromC(C.av_packet_clone(p.c))
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#gade00f67930f4e2a3401b67b701d5b3a2
func (p *Packet) CopyProperties(src *Packet) error {
	return newError(C.av_packet_copy_props(p.c, src.c))
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#ga8a6deff6c1809029037ffd760db3e0d4
func (p *Packet) MakeReferenceCounted() error {
	return newError(C.av_packet_make_refcounted(p.c))
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#gaaa304ffdab83984ac995d134e4298d4b
func (p *Packet) MakeWritable() error {
	return newError(C.av_packet_make_writable(p.c))
}

func (p *Packet) IsWritable() bool {
	return p.c.buf != nil && C.av_buffer_is_writable(p.c.buf) != 0
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#gadfa708660b85a56749c753124de2da7d
func (p *Packet) AllocPayload(s int) error {
	return newError(C.av_new_packet(p.c, C.int(s)))
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#gadb05d71fa2ea7b7fd3e8cfc6d9065a47
func (p *Packet) Ref(src *Packet) error {
	return newError(C.av_packet_ref(p.c, src.c))
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#ga63d5a489b419bd5d45cfd09091cbcbc2
func (p *Packet) Unref() {
	C.av_packet_unref(p.c)
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#ga91dbb1359f99547adb544ee96a406b21
func (p *Packet) MoveRef(src *Packet) {
	C.av_packet_move_ref(p.c, src.c)
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#gae5c86e4d93f6e7aa62ef2c60763ea67e
func (p *Packet) RescaleTs(src, dst Rational) {
	C.av_packet_rescale_ts(p.c, src.c, dst.c)
}

// https://ffmpeg.org/doxygen/7.1/group__lavc__packet.html#ga7ca877e1f0ded89a27199b65e9a077dc
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
