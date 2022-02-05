package astiav

//#cgo pkg-config: libavutil
//#include <libavutil/dict.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavutil/dict.h#L79
type DictionaryEntry struct {
	c *C.struct_AVDictionaryEntry
}

func newDictionaryEntryFromC(c *C.struct_AVDictionaryEntry) *DictionaryEntry {
	return &DictionaryEntry{c: c}
}

func (e DictionaryEntry) Key() string {
	return C.GoString(e.c.key)
}

func (e DictionaryEntry) Value() string {
	return C.GoString(e.c.value)
}
