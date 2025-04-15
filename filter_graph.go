package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"errors"
	"math"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVFilterGraph.html
type FilterGraph struct {
	classerHandler
	c *C.AVFilterGraph
	// We need to store filter contexts to clean classer once filter graph is freed
	fcs []*FilterContext
}

func newFilterGraphFromC(c *C.AVFilterGraph) *FilterGraph {
	if c == nil {
		return nil
	}
	g := &FilterGraph{c: c}
	classers.set(g)
	return g
}

var _ Classer = (*FilterGraph)(nil)

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#ga6c778454b86f845805ffd814b4ce51d4
func AllocFilterGraph() *FilterGraph {
	return newFilterGraphFromC(C.avfilter_graph_alloc())
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#ga871684449dac05050df238a18d0d493b
func (g *FilterGraph) Free() {
	if g.c != nil {
		// Make sure to clone the classer before freeing the object since
		// the C free method may reset the pointer
		c := newClonedClasser(g)
		var cfcs []*ClonedClasser
		for _, fc := range g.fcs {
			cfcs = append(cfcs, newClonedClasser(fc))
		}
		C.avfilter_graph_free(&g.c)
		// Make sure to remove from classers after freeing the object since
		// the C free method may use methods needing the classer
		for _, cfc := range cfcs {
			if cfc != nil {
				classers.del(cfc)
			}
		}
		if c != nil {
			classers.del(c)
		}
	}
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#gadb442aca4e5a8c3ba740f6049f0a288b
func (g *FilterGraph) String() string {
	return C.GoString(C.avfilter_graph_dump(g.c, nil))
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterGraph.html#af00925dd69b474fac48887efc0e1ac94
func (g *FilterGraph) Class() *Class {
	if g.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(g.c))
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterGraph.html#ac28dcbf76e6fdd800295a2738d41660e
func (g *FilterGraph) ThreadCount() int {
	return int(g.c.nb_threads)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterGraph.html#ac28dcbf76e6fdd800295a2738d41660e
func (g *FilterGraph) SetThreadCount(threadCount int) {
	g.c.nb_threads = C.int(threadCount)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterGraph.html#a7fd96bbd6d1a3b730681dc0bf5107a5e
func (g *FilterGraph) ThreadType() ThreadType {
	return ThreadType(g.c.thread_type)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterGraph.html#a7fd96bbd6d1a3b730681dc0bf5107a5e
func (g *FilterGraph) SetThreadType(t ThreadType) {
	g.c.thread_type = C.int(t)
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#gac0788a9ab6966dba9318b5d5c7524fea
func (g *FilterGraph) NewBuffersinkFilterContext(f *Filter, name string) (*BuffersinkFilterContext, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var c *C.AVFilterContext
	if err := newError(C.avfilter_graph_create_filter(&c, f.c, cname, nil, nil, g.c)); err != nil {
		return nil, err
	}
	fc := newFilterContext(c)
	g.fcs = append(g.fcs, fc)
	return newBuffersinkFilterContext(fc), nil
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#gaa9af17ecf4c5c87307b57cf08411088b
func (g *FilterGraph) NewBuffersrcFilterContext(f *Filter, name string) (*BuffersrcFilterContext, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	c := C.avfilter_graph_alloc_filter(g.c, f.c, cname)
	if c == nil {
		return nil, errors.New("astiav: allocating filter context failed")
	}
	fc := newFilterContext(c)
	g.fcs = append(g.fcs, fc)
	return newBuffersrcFilterContext(fc), nil
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#ga34f4ff420bd58da6747a3ff1fbedd001
func (g *FilterGraph) Parse(content string, inputs, outputs *FilterInOut) error {
	cc := C.CString(content)
	defer C.free(unsafe.Pointer(cc))
	var ic **C.AVFilterInOut
	if inputs != nil {
		ic = &inputs.c
	}
	var oc **C.AVFilterInOut
	if outputs != nil {
		oc = &outputs.c
	}
	return newError(C.avfilter_graph_parse_ptr(g.c, cc, ic, oc, nil))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#ga2ecfd3667219b6cd1e37b7047cc0ef2b
func (g *FilterGraph) ParseSegment(content string) (*FilterGraphSegment, error) {
	cc := C.CString(content)
	defer C.free(unsafe.Pointer(cc))
	var cs *C.AVFilterGraphSegment
	if err := newError(C.avfilter_graph_segment_parse(g.c, cc, 0, &cs)); err != nil {
		return nil, err
	}
	return newFilterGraphSegmentFromC(cs), nil
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#ga1896c46b7bc6ff1bdb1a4815faa9ad07
func (g *FilterGraph) Configure() error {
	return newError(C.avfilter_graph_config(g.c, nil))
}

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#gaaad7850fb5fe26d35e5d371ca75b79e1
func (g *FilterGraph) SendCommand(target, cmd, args string, f FilterCommandFlags) (response string, err error) {
	targetc := C.CString(target)
	defer C.free(unsafe.Pointer(targetc))
	cmdc := C.CString(cmd)
	defer C.free(unsafe.Pointer(cmdc))
	argsc := C.CString(args)
	defer C.free(unsafe.Pointer(argsc))
	response, err = stringFromC(255, func(buf *C.char, size C.size_t) error {
		return newError(C.avfilter_graph_send_command(g.c, targetc, cmdc, argsc, buf, C.int(size), C.int(f)))
	})
	return
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterGraph.html#a0ba5c820c760788ea5f8e40c476f9704
func (g *FilterGraph) NbFilters() int {
	return int(g.c.nb_filters)
}

// https://ffmpeg.org/doxygen/7.0/structAVFilterGraph.html#a1dafd3d239f7c2f5e3ac109578ef926d
func (g *FilterGraph) Filters() (fs []*FilterContext) {
	fcs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVFilterContext)(nil))](*C.AVFilterContext))(unsafe.Pointer(g.c.filters))
	for i := 0; i < g.NbFilters(); i++ {
		fs = append(fs, newFilterContext(fcs[i]))
	}
	return
}
