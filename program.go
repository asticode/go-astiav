package astiav

//#include <libavformat/avformat.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html
type Program struct {
	c  *C.AVProgram
	fc *FormatContext
}

func newProgramFromC(c *C.AVProgram, fc *FormatContext) *Program {
	if c == nil {
		return nil
	}
	return &Program{
		c:  c,
		fc: fc,
	}
}

// https://ffmpeg.org/doxygen/7.0/avformat_8c.html#ae1eb83cf16060217c805e61f0f62fa4e
func (p *Program) AddStream(s *Stream) {
	C.av_program_add_stream_index(p.fc.c, p.c.id, C.uint(s.c.index))
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a10cc799a98b37335e820b0bdb386eb95
func (p *Program) ID() int {
	return int(p.c.id)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a10cc799a98b37335e820b0bdb386eb95
func (p *Program) SetID(i int) {
	p.c.id = C.int(i)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a136cf29d2aa5b0e4c6d743406c5e39d1
func (p *Program) NbStreams() int {
	return int(p.c.nb_stream_indexes)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a7967d41af4812ed61a28762e988c7a02
func (p *Program) Streams() (ss []*Stream) {
	is := make(map[int]bool)
	for _, idx := range unsafe.Slice(p.c.stream_index, p.c.nb_stream_indexes) {
		is[int(idx)] = true
	}
	for _, s := range p.fc.Streams() {
		if _, ok := is[s.Index()]; ok {
			ss = append(ss, s)
		}
	}
	return
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a9c7a07c08a1f960aaa49f3f47633af5c
func (p *Program) Discard() Discard {
	return Discard(p.c.discard)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a9c7a07c08a1f960aaa49f3f47633af5c
func (p *Program) SetDiscard(d Discard) {
	p.c.discard = int32(d)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#ae9dab38d4694e3da9cba0f882f4e43d3
func (p *Program) Metadata() *Dictionary {
	return newDictionaryFromC(p.c.metadata)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#ae9dab38d4694e3da9cba0f882f4e43d3
func (p *Program) SetMetadata(d *Dictionary) {
	p.c.metadata = d.c
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a4c1539ea3c98da979b95a59a3ea163cb
func (p *Program) ProgramNumber() int {
	return int(p.c.program_num)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a4c1539ea3c98da979b95a59a3ea163cb
func (p *Program) SetProgramNumber(n int) {
	p.c.program_num = C.int(n)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a02011963a63c291c6dc6d4eefa56cd69
func (p *Program) PmtPid() int {
	return int(p.c.pmt_pid)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a02011963a63c291c6dc6d4eefa56cd69
func (p *Program) SetPmtPid(n int) {
	p.c.pmt_pid = C.int(n)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a7e026323df87e84a72ec5e5c8ce341a5
func (p *Program) PcrPid() int {
	return int(p.c.pcr_pid)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a7e026323df87e84a72ec5e5c8ce341a5
func (p *Program) SetPcrPid(n int) {
	p.c.pcr_pid = C.int(n)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a2276db4d51695120664d527f20b7c532
func (p *Program) StartTime() int64 {
	return int64(p.c.start_time)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a5a7795c918153d0f64d68a838e172db4
func (p *Program) EndTime() int64 {
	return int64(p.c.end_time)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#a7e539e286876577e158039f6e7678452
func (p *Program) PtsWrapReference() int64 {
	return int64(p.c.pts_wrap_reference)
}

// https://ffmpeg.org/doxygen/7.0/structAVProgram.html#aa3f8af78093a910ff766ac5af381758b
func (p *Program) PtsWrapBehavior() int {
	return int(p.c.pts_wrap_behavior)
}
