package astiav

//#include <libavfilter/avfilter.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n7.0/libavfilter/avfilter.h#L1075
type FilterParams struct {
	c *C.AVFilterParams
}

func newFilterParamsFromC(c *C.AVFilterParams) *FilterParams {
	if c == nil {
		return nil
	}
	return &FilterParams{c: c}
}

func (fp *FilterParams) FilterName() string {
	return C.GoString(fp.c.filter_name)
}
