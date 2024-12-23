package astiav

//#include <libavfilter/avfilter.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/group__lavfi.html#gae6ed6c10a03508829bdf17560e3e10e5
type FilterFlag int64

const (
	FilterFlagDynamicInputs           = FilterFlag(C.AVFILTER_FLAG_DYNAMIC_INPUTS)
	FilterFlagDynamicOutputs          = FilterFlag(C.AVFILTER_FLAG_DYNAMIC_OUTPUTS)
	FilterFlagSliceThreads            = FilterFlag(C.AVFILTER_FLAG_SLICE_THREADS)
	FilterFlagMetadataOnly            = FilterFlag(C.AVFILTER_FLAG_METADATA_ONLY)
	FilterFlagHardwareDevice          = FilterFlag(C.AVFILTER_FLAG_HWDEVICE)
	FilterFlagSupportTimelineGeneric  = FilterFlag(C.AVFILTER_FLAG_SUPPORT_TIMELINE_GENERIC)
	FilterFlagSupportTimelineInternal = FilterFlag(C.AVFILTER_FLAG_SUPPORT_TIMELINE_INTERNAL)
	FilterFlagSupportTimeline         = FilterFlag(C.AVFILTER_FLAG_SUPPORT_TIMELINE)
)
