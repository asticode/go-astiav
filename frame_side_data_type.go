package astiav

//#include <libavutil/frame.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/frame.h#L48
type FrameSideDataType C.enum_AVFrameSideDataType

const (
	FrameSideDataTypePanscan                  = FrameSideDataType(C.AV_FRAME_DATA_PANSCAN)
	FrameSideDataTypeA53Cc                    = FrameSideDataType(C.AV_FRAME_DATA_A53_CC)
	FrameSideDataTypeStereo3D                 = FrameSideDataType(C.AV_FRAME_DATA_STEREO3D)
	FrameSideDataTypeMatrixencoding           = FrameSideDataType(C.AV_FRAME_DATA_MATRIXENCODING)
	FrameSideDataTypeDownmixInfo              = FrameSideDataType(C.AV_FRAME_DATA_DOWNMIX_INFO)
	FrameSideDataTypeReplaygain               = FrameSideDataType(C.AV_FRAME_DATA_REPLAYGAIN)
	FrameSideDataTypeDisplaymatrix            = FrameSideDataType(C.AV_FRAME_DATA_DISPLAYMATRIX)
	FrameSideDataTypeAfd                      = FrameSideDataType(C.AV_FRAME_DATA_AFD)
	FrameSideDataTypeMotionVectors            = FrameSideDataType(C.AV_FRAME_DATA_MOTION_VECTORS)
	FrameSideDataTypeSkipSamples              = FrameSideDataType(C.AV_FRAME_DATA_SKIP_SAMPLES)
	FrameSideDataTypeAudioServiceType         = FrameSideDataType(C.AV_FRAME_DATA_AUDIO_SERVICE_TYPE)
	FrameSideDataTypeMasteringDisplayMetadata = FrameSideDataType(C.AV_FRAME_DATA_MASTERING_DISPLAY_METADATA)
	FrameSideDataTypeGopTimecode              = FrameSideDataType(C.AV_FRAME_DATA_GOP_TIMECODE)
	FrameSideDataTypeSpherical                = FrameSideDataType(C.AV_FRAME_DATA_SPHERICAL)
	FrameSideDataTypeContentLightLevel        = FrameSideDataType(C.AV_FRAME_DATA_CONTENT_LIGHT_LEVEL)
	FrameSideDataTypeIccProfile               = FrameSideDataType(C.AV_FRAME_DATA_ICC_PROFILE)
	FrameSideDataTypeS12MTimecode             = FrameSideDataType(C.AV_FRAME_DATA_S12M_TIMECODE)
	FrameSideDataTypeDynamicHdrPlus           = FrameSideDataType(C.AV_FRAME_DATA_DYNAMIC_HDR_PLUS)
	FrameSideDataTypeRegionsOfInterest        = FrameSideDataType(C.AV_FRAME_DATA_REGIONS_OF_INTEREST)
	FrameSideDataTypeVideoEncParams           = FrameSideDataType(C.AV_FRAME_DATA_VIDEO_ENC_PARAMS)
	FrameSideDataTypeSeiUnregistered          = FrameSideDataType(C.AV_FRAME_DATA_SEI_UNREGISTERED)
	FrameSideDataTypeFilmGrainParams          = FrameSideDataType(C.AV_FRAME_DATA_FILM_GRAIN_PARAMS)
)
