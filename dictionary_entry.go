package astiav

//#include <libavutil/dict.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/structAVDictionaryEntry.html
type DictionaryEntry struct {
	c *C.AVDictionaryEntry
}

func newDictionaryEntryFromC(c *C.AVDictionaryEntry) *DictionaryEntry {
	return &DictionaryEntry{c: c}
}

// https://ffmpeg.org/doxygen/7.0/structAVDictionaryEntry.html#a38fc80176f8f839282bb61c03392e194
func (e DictionaryEntry) Key() string {
	return C.GoString(e.c.key)
}

// https://ffmpeg.org/doxygen/7.0/structAVDictionaryEntry.html#aa38678f2cad36f120d42e56449c5edb4
func (e DictionaryEntry) Value() string {
	return C.GoString(e.c.value)
}
