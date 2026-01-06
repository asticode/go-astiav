package astiav

//#include <libavutil/dict.h>
import "C"

// https://ffmpeg.org/doxygen/7.1/group__lavu__dict.html#gad9cbc53cec515b72ae7caa2e194c6bc0
type DictionaryFlag int64

const (
	DictionaryFlagMatchCase     = DictionaryFlag(C.AV_DICT_MATCH_CASE)
	DictionaryFlagIgnoreSuffix  = DictionaryFlag(C.AV_DICT_IGNORE_SUFFIX)
	DictionaryFlagDontStrdupKey = DictionaryFlag(C.AV_DICT_DONT_STRDUP_KEY)
	DictionaryFlagDontStrdupVal = DictionaryFlag(C.AV_DICT_DONT_STRDUP_VAL)
	DictionaryFlagDontOverwrite = DictionaryFlag(C.AV_DICT_DONT_OVERWRITE)
	DictionaryFlagAppend        = DictionaryFlag(C.AV_DICT_APPEND)
	DictionaryFlagMultikey      = DictionaryFlag(C.AV_DICT_MULTIKEY)
)
