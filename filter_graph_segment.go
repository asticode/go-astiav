package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"math"
	"unsafe"
)

// https://ffmpeg.org/doxygen/8.1/structAVFilterGraphSegment.html
type FilterGraphSegment struct {
	c *C.AVFilterGraphSegment
}

func newFilterGraphSegmentFromC(c *C.AVFilterGraphSegment) *FilterGraphSegment {
	if c == nil {
		return nil
	}
	return &FilterGraphSegment{c: c}
}

// https://ffmpeg.org/doxygen/8.1/group__lavfi.html#ga51283edd8f3685e1f33239f360e14ae8
func (fgs *FilterGraphSegment) Free() {
	if fgs.c != nil {
		C.avfilter_graph_segment_free(&fgs.c)
	}
}

// https://ffmpeg.org/doxygen/8.1/structAVFilterGraphSegment.html#ad5a2779af221d1520490fe2719f9e39a
func (fgs *FilterGraphSegment) Chains() (cs []*FilterChain) {
	ccs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVFilterChain)(nil))](*C.AVFilterChain))(unsafe.Pointer(fgs.c.chains))
	for i := 0; i < fgs.NbChains(); i++ {
		cs = append(cs, newFilterChainFromC(ccs[i]))
	}
	return
}

// https://ffmpeg.org/doxygen/8.1/structAVFilterGraphSegment.html#ab7563eca151d89e693f6258de5ce0214
func (fgs *FilterGraphSegment) NbChains() int {
	return int(fgs.c.nb_chains)
}
