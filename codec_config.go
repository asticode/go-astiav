package astiav

//#include <libavcodec/avcodec.h>
import "C"
import (
	"unsafe"
)

// CodecConfig represents AVCodecConfig enum
type CodecConfig int

const (
	CodecConfigPixelFormat    CodecConfig = C.AV_CODEC_CONFIG_PIX_FORMAT
	CodecConfigFrameRate      CodecConfig = C.AV_CODEC_CONFIG_FRAME_RATE
	CodecConfigSampleRate     CodecConfig = C.AV_CODEC_CONFIG_SAMPLE_RATE
	CodecConfigSampleFormat   CodecConfig = C.AV_CODEC_CONFIG_SAMPLE_FORMAT
	CodecConfigChannelLayout  CodecConfig = C.AV_CODEC_CONFIG_CHANNEL_LAYOUT
	CodecConfigColorRange     CodecConfig = C.AV_CODEC_CONFIG_COLOR_RANGE
	CodecConfigColorSpace     CodecConfig = C.AV_CODEC_CONFIG_COLOR_SPACE
)

// 解析像素格式数组
func parsePixelFormats(ptr unsafe.Pointer, count int) []PixelFormat {
	if ptr == nil || count == 0 {
		return nil
	}
	
	formats := (*[1000]C.enum_AVPixelFormat)(ptr)[:count:count]
	result := make([]PixelFormat, count)
	for i, format := range formats {
		result[i] = PixelFormat(format)
	}
	return result
}

// 解析采样格式数组
func parseSampleFormats(ptr unsafe.Pointer, count int) []SampleFormat {
	if ptr == nil || count == 0 {
		return nil
	}
	
	formats := (*[1000]C.enum_AVSampleFormat)(ptr)[:count:count]
	result := make([]SampleFormat, count)
	for i, format := range formats {
		result[i] = SampleFormat(format)
	}
	return result
}

// 解析采样率数组
func parseSampleRates(ptr unsafe.Pointer, count int) []int {
	if ptr == nil || count == 0 {
		return nil
	}
	
	rates := (*[1000]C.int)(ptr)[:count:count]
	result := make([]int, count)
	for i, rate := range rates {
		result[i] = int(rate)
	}
	return result
}

// 解析通道布局数组
func parseChannelLayouts(ptr unsafe.Pointer, count int) []ChannelLayout {
	if ptr == nil || count == 0 {
		return nil
	}
	
	layouts := (*[1000]C.AVChannelLayout)(ptr)[:count:count]
	result := make([]ChannelLayout, count)
	for i, layout := range layouts {
		result[i] = newChannelLayoutFromC(&layout)
	}
	return result
}

// 解析帧率数组
func parseFrameRates(ptr unsafe.Pointer, count int) []Rational {
	if ptr == nil || count == 0 {
		return nil
	}
	
	rates := (*[1000]C.AVRational)(ptr)[:count:count]
	result := make([]Rational, count)
	for i, rate := range rates {
		result[i] = newRationalFromC(rate)
	}
	return result
}

// 解析颜色范围数组
func parseColorRanges(ptr unsafe.Pointer, count int) []ColorRange {
	if ptr == nil || count == 0 {
		return nil
	}
	
	ranges := (*[1000]C.enum_AVColorRange)(ptr)[:count:count]
	result := make([]ColorRange, count)
	for i, colorRange := range ranges {
		result[i] = ColorRange(colorRange)
	}
	return result
}

// 解析颜色空间数组
func parseColorSpaces(ptr unsafe.Pointer, count int) []ColorSpace {
	if ptr == nil || count == 0 {
		return nil
	}
	
	spaces := (*[1000]C.enum_AVColorSpace)(ptr)[:count:count]
	result := make([]ColorSpace, count)
	for i, space := range spaces {
		result[i] = ColorSpace(space)
	}
	return result
}