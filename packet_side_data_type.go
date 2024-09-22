package astiav

//#include <libavcodec/avcodec.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavcodec/packet.h#L40
type PacketSideDataType C.enum_AVPacketSideDataType

const (
	PacketSideDataTypeA53Cc                    = PacketSideDataType(C.AV_PKT_DATA_A53_CC)
	PacketSideDataTypeAfd                      = PacketSideDataType(C.AV_PKT_DATA_AFD)
	PacketSideDataTypeAudioServiceType         = PacketSideDataType(C.AV_PKT_DATA_AUDIO_SERVICE_TYPE)
	PacketSideDataTypeContentLightLevel        = PacketSideDataType(C.AV_PKT_DATA_CONTENT_LIGHT_LEVEL)
	PacketSideDataTypeCpbProperties            = PacketSideDataType(C.AV_PKT_DATA_CPB_PROPERTIES)
	PacketSideDataTypeDisplaymatrix            = PacketSideDataType(C.AV_PKT_DATA_DISPLAYMATRIX)
	PacketSideDataTypeEncryptionInfo           = PacketSideDataType(C.AV_PKT_DATA_ENCRYPTION_INFO)
	PacketSideDataTypeEncryptionInitInfo       = PacketSideDataType(C.AV_PKT_DATA_ENCRYPTION_INIT_INFO)
	PacketSideDataTypeFallbackTrack            = PacketSideDataType(C.AV_PKT_DATA_FALLBACK_TRACK)
	PacketSideDataTypeH263MbInfo               = PacketSideDataType(C.AV_PKT_DATA_H263_MB_INFO)
	PacketSideDataTypeJpDualmono               = PacketSideDataType(C.AV_PKT_DATA_JP_DUALMONO)
	PacketSideDataTypeMasteringDisplayMetadata = PacketSideDataType(C.AV_PKT_DATA_MASTERING_DISPLAY_METADATA)
	PacketSideDataTypeMatroskaBlockadditional  = PacketSideDataType(C.AV_PKT_DATA_MATROSKA_BLOCKADDITIONAL)
	PacketSideDataTypeMetadataUpdate           = PacketSideDataType(C.AV_PKT_DATA_METADATA_UPDATE)
	PacketSideDataTypeMpegtsStreamId           = PacketSideDataType(C.AV_PKT_DATA_MPEGTS_STREAM_ID)
	PacketSideDataTypeNb                       = PacketSideDataType(C.AV_PKT_DATA_NB)
	PacketSideDataTypeNewExtradata             = PacketSideDataType(C.AV_PKT_DATA_NEW_EXTRADATA)
	PacketSideDataTypePalette                  = PacketSideDataType(C.AV_PKT_DATA_PALETTE)
	PacketSideDataTypeParamChange              = PacketSideDataType(C.AV_PKT_DATA_PARAM_CHANGE)
	PacketSideDataTypeQualityStats             = PacketSideDataType(C.AV_PKT_DATA_QUALITY_STATS)
	PacketSideDataTypeReplaygain               = PacketSideDataType(C.AV_PKT_DATA_REPLAYGAIN)
	PacketSideDataTypeSkipSamples              = PacketSideDataType(C.AV_PKT_DATA_SKIP_SAMPLES)
	PacketSideDataTypeSpherical                = PacketSideDataType(C.AV_PKT_DATA_SPHERICAL)
	PacketSideDataTypeStereo3D                 = PacketSideDataType(C.AV_PKT_DATA_STEREO3D)
	PacketSideDataTypeStringsMetadata          = PacketSideDataType(C.AV_PKT_DATA_STRINGS_METADATA)
	PacketSideDataTypeSubtitlePosition         = PacketSideDataType(C.AV_PKT_DATA_SUBTITLE_POSITION)
	PacketSideDataTypeWebvttIdentifier         = PacketSideDataType(C.AV_PKT_DATA_WEBVTT_IDENTIFIER)
	PacketSideDataTypeWebvttSettings           = PacketSideDataType(C.AV_PKT_DATA_WEBVTT_SETTINGS)
)

func (t PacketSideDataType) Name() string {
	return C.GoString(C.av_packet_side_data_name((C.enum_AVPacketSideDataType)(t)))
}

func (t PacketSideDataType) String() string {
	return t.Name()
}
