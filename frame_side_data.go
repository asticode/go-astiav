package astiav

//#include <libavutil/frame.h>
//#include "frame_side_data.h"
import "C"
import (
	"errors"
	"math"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.1/structAVFrameSideData.html
// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#gae01fa7e427274293aacdf2adc17076bc
type FrameSideData struct {
	sd   ***C.AVFrameSideData
	size *C.int
}

func newFrameSideDataFromC(sd ***C.AVFrameSideData, size *C.int) *FrameSideData {
	return &FrameSideData{
		sd:   sd,
		size: size,
	}
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__frame.html#ggae01fa7e427274293aacdf2adc17076bcaf525ec92d2c5a78d44950bc3f29972aa
func (d *FrameSideData) RegionsOfInterest() *frameSideDataRegionsOfInterest {
	return newFrameSideDataRegionsOfInterest(d)
}

type frameSideDataRegionsOfInterest struct {
	d *FrameSideData
}

func newFrameSideDataRegionsOfInterest(d *FrameSideData) *frameSideDataRegionsOfInterest {
	return &frameSideDataRegionsOfInterest{d: d}
}

func (d *frameSideDataRegionsOfInterest) data(sd *C.AVFrameSideData) *[(math.MaxInt32 - 1) / C.sizeof_AVRegionOfInterest]C.AVRegionOfInterest {
	return (*[(math.MaxInt32 - 1) / C.sizeof_AVRegionOfInterest](C.AVRegionOfInterest))(unsafe.Pointer(C.astiavConvertRegionsOfInterestFrameSideData(sd)))
}

func (d *frameSideDataRegionsOfInterest) Add(rois []RegionOfInterest) error {
	sd := C.av_frame_side_data_new(d.d.sd, d.d.size, C.AV_FRAME_DATA_REGIONS_OF_INTEREST, C.size_t(C.sizeof_AVRegionOfInterest*len(rois)), 0)
	if sd == nil {
		return errors.New("astiav: nil pointer")
	}

	crois := d.data(sd)
	for i, roi := range rois {
		crois[i].bottom = C.int(roi.Bottom)
		crois[i].left = C.int(roi.Left)
		crois[i].qoffset = roi.QuantisationOffset.c
		crois[i].right = C.int(roi.Right)
		crois[i].self_size = C.sizeof_AVRegionOfInterest
		crois[i].top = C.int(roi.Top)
	}
	return nil
}

func (d *frameSideDataRegionsOfInterest) Get() ([]RegionOfInterest, bool) {
	sd := C.av_frame_side_data_get(*d.d.sd, *d.d.size, C.AV_FRAME_DATA_REGIONS_OF_INTEREST)
	if sd == nil {
		return nil, false
	}

	crois := d.data(sd)
	rois := make([]RegionOfInterest, int(sd.size/C.sizeof_AVRegionOfInterest))
	for i := range rois {
		rois[i] = RegionOfInterest{
			Bottom:             int(crois[i].bottom),
			Left:               int(crois[i].left),
			QuantisationOffset: newRationalFromC(crois[i].qoffset),
			Right:              int(crois[i].right),
			Top:                int(crois[i].top),
		}
	}
	return rois, true
}
