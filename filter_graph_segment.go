package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"math"
	"unsafe"
)

// https://ffmpeg.org/doxygen/8.0/structAVFilterGraphSegment.html
type FilterGraphSegment struct {
	c *C.AVFilterGraphSegment
}

func newFilterGraphSegmentFromC(c *C.AVFilterGraphSegment) *FilterGraphSegment {
	if c == nil {
		return nil
	}
	return &FilterGraphSegment{c: c}
}

// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga51283edd8f3685e1f33239f360e14ae8
func (fgs *FilterGraphSegment) Free() {
	if fgs.c != nil {
		C.avfilter_graph_segment_free(&fgs.c)
	}
}

// https://ffmpeg.org/doxygen/8.0/structAVFilterGraphSegment.html#ad5a2779af221d1520490fe2719f9e39a
func (fgs *FilterGraphSegment) Chains() (cs []*FilterChain) {
	ccs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVFilterChain)(nil))](*C.AVFilterChain))(unsafe.Pointer(fgs.c.chains))
	for i := 0; i < fgs.NbChains(); i++ {
		cs = append(cs, newFilterChainFromC(ccs[i]))
	}
	return
}

// https://ffmpeg.org/doxygen/8.0/structAVFilterGraphSegment.html#ab7563eca151d89e693f6258de5ce0214
func (fgs *FilterGraphSegment) NbChains() int {
	return int(fgs.c.nb_chains)
}

// ParseFilterGraphSegment parses a filter graph segment from a string
// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (fg *FilterGraph) ParseFilterGraphSegment(graph string, flags int) (*FilterGraphSegment, error) {
	cGraph := C.CString(graph)
	defer C.free(unsafe.Pointer(cGraph))
	
	var seg *C.AVFilterGraphSegment
	if err := newError(C.avfilter_graph_segment_parse(fg.c, cGraph, C.int(flags), &seg)); err != nil {
		return nil, err
	}
	
	return newFilterGraphSegmentFromC(seg), nil
}

// CreateFilters creates filters in the segment
// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (fgs *FilterGraphSegment) CreateFilters(flags int) error {
	return newError(C.avfilter_graph_segment_create_filters(fgs.c, C.int(flags)))
}

// ApplyOpts applies options to the segment
// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (fgs *FilterGraphSegment) ApplyOpts(flags int) error {
	return newError(C.avfilter_graph_segment_apply_opts(fgs.c, C.int(flags)))
}

// Init initializes the segment
// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (fgs *FilterGraphSegment) Init(flags int) error {
	return newError(C.avfilter_graph_segment_init(fgs.c, C.int(flags)))
}

// Link links the segment
// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (fgs *FilterGraphSegment) Link(flags int, inputs, outputs **FilterInOut) error {
	var cInputs, cOutputs **C.AVFilterInOut
	if inputs != nil {
		cInputs = (**C.AVFilterInOut)(unsafe.Pointer(inputs))
	}
	if outputs != nil {
		cOutputs = (**C.AVFilterInOut)(unsafe.Pointer(outputs))
	}
	return newError(C.avfilter_graph_segment_link(fgs.c, C.int(flags), cInputs, cOutputs))
}

// Apply applies the segment to the filter graph
// https://ffmpeg.org/doxygen/8.0/group__lavfi.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (fgs *FilterGraphSegment) Apply(flags int, inputs, outputs **FilterInOut) error {
	var cInputs, cOutputs **C.AVFilterInOut
	if inputs != nil {
		cInputs = (**C.AVFilterInOut)(unsafe.Pointer(inputs))
	}
	if outputs != nil {
		cOutputs = (**C.AVFilterInOut)(unsafe.Pointer(outputs))
	}
	return newError(C.avfilter_graph_segment_apply(fgs.c, C.int(flags), cInputs, cOutputs))
}
