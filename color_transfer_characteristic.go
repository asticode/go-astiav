package astiav

//#include <libavutil/pixfmt.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/pixfmt.h#L494
type ColorTransferCharacteristic C.enum_AVColorTransferCharacteristic

const (
	ColorTransferCharacteristicReserved0   = ColorTransferCharacteristic(C.AVCOL_TRC_RESERVED0)
	ColorTransferCharacteristicBt709       = ColorTransferCharacteristic(C.AVCOL_TRC_BT709)
	ColorTransferCharacteristicUnspecified = ColorTransferCharacteristic(C.AVCOL_TRC_UNSPECIFIED)
	ColorTransferCharacteristicReserved    = ColorTransferCharacteristic(C.AVCOL_TRC_RESERVED)
	ColorTransferCharacteristicGamma22     = ColorTransferCharacteristic(C.AVCOL_TRC_GAMMA22)
	ColorTransferCharacteristicGamma28     = ColorTransferCharacteristic(C.AVCOL_TRC_GAMMA28)
	ColorTransferCharacteristicSmpte170M   = ColorTransferCharacteristic(C.AVCOL_TRC_SMPTE170M)
	ColorTransferCharacteristicSmpte240M   = ColorTransferCharacteristic(C.AVCOL_TRC_SMPTE240M)
	ColorTransferCharacteristicLinear      = ColorTransferCharacteristic(C.AVCOL_TRC_LINEAR)
	ColorTransferCharacteristicLog         = ColorTransferCharacteristic(C.AVCOL_TRC_LOG)
	ColorTransferCharacteristicLogSqrt     = ColorTransferCharacteristic(C.AVCOL_TRC_LOG_SQRT)
	ColorTransferCharacteristicIec6196624  = ColorTransferCharacteristic(C.AVCOL_TRC_IEC61966_2_4)
	ColorTransferCharacteristicBt1361Ecg   = ColorTransferCharacteristic(C.AVCOL_TRC_BT1361_ECG)
	ColorTransferCharacteristicIec6196621  = ColorTransferCharacteristic(C.AVCOL_TRC_IEC61966_2_1)
	ColorTransferCharacteristicBt202010    = ColorTransferCharacteristic(C.AVCOL_TRC_BT2020_10)
	ColorTransferCharacteristicBt202012    = ColorTransferCharacteristic(C.AVCOL_TRC_BT2020_12)
	ColorTransferCharacteristicSmpte2084   = ColorTransferCharacteristic(C.AVCOL_TRC_SMPTE2084)
	ColorTransferCharacteristicSmptest2084 = ColorTransferCharacteristic(C.AVCOL_TRC_SMPTEST2084)
	ColorTransferCharacteristicSmpte428    = ColorTransferCharacteristic(C.AVCOL_TRC_SMPTE428)
	ColorTransferCharacteristicSmptest4281 = ColorTransferCharacteristic(C.AVCOL_TRC_SMPTEST428_1)
	ColorTransferCharacteristicAribStdB67  = ColorTransferCharacteristic(C.AVCOL_TRC_ARIB_STD_B67)
	ColorTransferCharacteristicNb          = ColorTransferCharacteristic(C.AVCOL_TRC_NB)
)
