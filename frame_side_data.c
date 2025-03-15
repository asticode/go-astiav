#include <libavutil/frame.h>

AVRegionOfInterest* astiavConvertRegionsOfInterestFrameSideData(AVFrameSideData *sd) {
    return (AVRegionOfInterest*)sd->data;
}