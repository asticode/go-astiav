package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"strings"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavfilter/avfilter.h#L861
type FilterGraph struct {
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

func AllocFilterGraph() *FilterGraph {
	return newFilterGraphFromC(C.avfilter_graph_alloc())
}

func (g *FilterGraph) Free() {
	if g.c != nil {
		// Make sure to clone the classer before freeing the object since
		// the C free method may reset the pointer
		c := newClonedClasser(g)
		var cfcs []Classer
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

func (g *FilterGraph) String() string {
	return C.GoString(C.avfilter_graph_dump(g.c, nil))
}

func (g *FilterGraph) Class() *Class {
	return newClassFromC(unsafe.Pointer(g.c))
}

func (g *FilterGraph) ThreadCount() int {
	return int(g.c.nb_threads)
}

func (g *FilterGraph) SetThreadCount(threadCount int) {
	g.c.nb_threads = C.int(threadCount)
}

func (g *FilterGraph) ThreadType() ThreadType {
	return ThreadType(g.c.thread_type)
}

func (g *FilterGraph) SetThreadType(t ThreadType) {
	g.c.thread_type = C.int(t)
}

type FilterArgs map[string]string

func (args FilterArgs) String() string {
	var ss []string
	for k, v := range args {
		ss = append(ss, k+"="+v)
	}
	return strings.Join(ss, ":")
}

func (g *FilterGraph) NewFilterContext(f *Filter, name string, args FilterArgs) (*FilterContext, error) {
	ca := (*C.char)(nil)
	if len(args) > 0 {
		ca = C.CString(args.String())
		defer C.free(unsafe.Pointer(ca))
	}
	cn := C.CString(name)
	defer C.free(unsafe.Pointer(cn))
	var c *C.AVFilterContext
	if err := newError(C.avfilter_graph_create_filter(&c, f.c, cn, ca, nil, g.c)); err != nil {
		return nil, err
	}
	fc := newFilterContext(c)
	g.fcs = append(g.fcs, fc)
	return fc, nil
}

func (g *FilterGraph) NewBuffersinkFilterContext(f *Filter, name string, args FilterArgs) (*BuffersinkFilterContext, error) {
	fc, err := g.NewFilterContext(f, name, args)
	if err != nil {
		return nil, err
	}
	return newBuffersinkFilterContext(fc), nil
}

func (g *FilterGraph) NewBuffersrcFilterContext(f *Filter, name string, args FilterArgs) (*BuffersrcFilterContext, error) {
	fc, err := g.NewFilterContext(f, name, args)
	if err != nil {
		return nil, err
	}
	return newBuffersrcFilterContext(fc), nil
}

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

func (g *FilterGraph) ParseSegment(content string) (*FilterGraphSegment, error) {
	cc := C.CString(content)
	defer C.free(unsafe.Pointer(cc))
	var cs *C.AVFilterGraphSegment
	if err := newError(C.avfilter_graph_segment_parse(g.c, cc, 0, &cs)); err != nil {
		return nil, err
	}
	return newFilterGraphSegmentFromC(cs), nil
}

func (g *FilterGraph) Configure() error {
	return newError(C.avfilter_graph_config(g.c, nil))
}

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
