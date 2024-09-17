package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"math"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n7.0/libavfilter/avfilter.h#L1142
type FilterChain struct {
	c *C.AVFilterChain
}

func newFilterChainFromC(c *C.AVFilterChain) *FilterChain {
	if c == nil {
		return nil
	}
	return &FilterChain{c: c}
}

func (fc *FilterChain) Filters() (fs []*FilterParams) {
	fcs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVFilterParams)(nil))](*C.AVFilterParams))(unsafe.Pointer(fc.c.filters))
	for i := 0; i < fc.NbFilters(); i++ {
		fs = append(fs, newFilterParamsFromC(fcs[i]))
	}
	return
}

func (fc *FilterChain) NbFilters() int {
	return int(fc.c.nb_filters)
}
