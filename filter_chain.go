package astiav

//#include <libavfilter/avfilter.h>
import "C"
import (
	"math"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.1/structAVFilterChain.html
type FilterChain struct {
	c *C.AVFilterChain
}

func newFilterChainFromC(c *C.AVFilterChain) *FilterChain {
	if c == nil {
		return nil
	}
	return &FilterChain{c: c}
}

// https://ffmpeg.org/doxygen/7.1/structAVFilterChain.html#aedebb337fac024e27b499fb3a0321f3e
func (fc *FilterChain) Filters() (fs []*FilterParams) {
	fcs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVFilterParams)(nil))](*C.AVFilterParams))(unsafe.Pointer(fc.c.filters))
	for i := 0; i < fc.NbFilters(); i++ {
		fs = append(fs, newFilterParamsFromC(fcs[i]))
	}
	return
}

// https://ffmpeg.org/doxygen/7.1/structAVFilterChain.html#abacf5280bd6db0d37a304b0dd0b6c54d
func (fc *FilterChain) NbFilters() int {
	return int(fc.c.nb_filters)
}
