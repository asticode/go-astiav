package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"math"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n7.0/libavfilter/avfilter.h#L1156
type FilterGraphSegment struct {
	c *C.AVFilterGraphSegment
}

func newFilterGraphSegmentFromC(c *C.AVFilterGraphSegment) *FilterGraphSegment {
	if c == nil {
		return nil
	}
	return &FilterGraphSegment{c: c}
}

func (fgs *FilterGraphSegment) Free() {
	C.avfilter_graph_segment_free(&fgs.c)
}

func (fgs *FilterGraphSegment) Chains() (cs []*FilterChain) {
	ccs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVFilterChain)(nil))](*C.AVFilterChain))(unsafe.Pointer(fgs.c.chains))
	for i := 0; i < fgs.NbChains(); i++ {
		cs = append(cs, newFilterChainFromC(ccs[i]))
	}
	return
}

func (fgs *FilterGraphSegment) NbChains() int {
	return int(fgs.c.nb_chains)
}
