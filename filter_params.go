package astiav

//#include <libavfilter/avfilter.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/structAVFilterParams.html
type FilterParams struct {
	c *C.AVFilterParams
}

func newFilterParamsFromC(c *C.AVFilterParams) *FilterParams {
	if c == nil {
		return nil
	}
	return &FilterParams{c: c}
}

// https://ffmpeg.org/doxygen/8.0/structAVFilterParams.html#a90edb3817b62f2ca70ea70001b84d001
func (fp *FilterParams) FilterName() string {
	return C.GoString(fp.c.filter_name)
}
