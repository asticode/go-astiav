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
