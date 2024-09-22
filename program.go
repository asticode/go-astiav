package astiav

//#include <libavformat/avformat.h>
import "C"
import (
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n7.0/libavformat/avformat.h#L1181
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

func (p *Program) AddStream(s *Stream) {
	C.av_program_add_stream_index(p.fc.c, p.c.id, C.uint(s.c.index))
}

func (p *Program) ID() int {
	return int(p.c.id)
}

func (p *Program) NbStreams() int {
	return int(p.c.nb_stream_indexes)
}

func (p *Program) SetID(i int) {
	p.c.id = C.int(i)
}

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
