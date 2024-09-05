package astiav

//#include <libavutil/dict.h>
import "C"

type DictionaryFlag int64

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/dict.h#L67
const (
	DictionaryFlagMatchCase     = DictionaryFlag(C.AV_DICT_MATCH_CASE)
	DictionaryFlagIgnoreSuffix  = DictionaryFlag(C.AV_DICT_IGNORE_SUFFIX)
	DictionaryFlagDontStrdupKey = DictionaryFlag(C.AV_DICT_DONT_STRDUP_KEY)
	DictionaryFlagDontStrdupVal = DictionaryFlag(C.AV_DICT_DONT_STRDUP_VAL)
	DictionaryFlagDontOverwrite = DictionaryFlag(C.AV_DICT_DONT_OVERWRITE)
	DictionaryFlagAppend        = DictionaryFlag(C.AV_DICT_APPEND)
	DictionaryFlagMultikey      = DictionaryFlag(C.AV_DICT_MULTIKEY)
)
