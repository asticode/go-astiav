package astiav

//#include <libavutil/frame.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/structAVRegionOfInterest.html
type RegionOfInterest struct {
	Bottom             int
	Left               int
	QuantisationOffset Rational
	Right              int
	Top                int
}
