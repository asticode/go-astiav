package astiav

//#include <libavcodec/avcodec.h>
//#include <libavutil/dict.h>
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

// https://ffmpeg.org/doxygen/8.0/structAVDictionary.html
type Dictionary struct {
	c *C.AVDictionary
}

func NewDictionary() *Dictionary {
	return &Dictionary{}
}

func newDictionaryFromC(c *C.AVDictionary) *Dictionary {
	if c == nil {
		return nil
	}
	return &Dictionary{c: c}
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__dict.html#ga8d9c2de72b310cef8e6a28c9cd3acbbe
func (d *Dictionary) Set(key, value string, flags DictionaryFlags) error {
	ck := C.CString(key)
	defer C.free(unsafe.Pointer(ck))
	cv := C.CString(value)
	defer C.free(unsafe.Pointer(cv))
	return newError(C.av_dict_set(&d.c, ck, cv, C.int(flags)))
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__dict.html#gaca5ff7c251e60bd13164d13c82f21b79
func (d *Dictionary) ParseString(i, keyValSep, pairsSep string, flags DictionaryFlags) error {
	ci := C.CString(i)
	defer C.free(unsafe.Pointer(ci))
	ck := C.CString(keyValSep)
	defer C.free(unsafe.Pointer(ck))
	cp := C.CString(pairsSep)
	defer C.free(unsafe.Pointer(cp))
	return newError(C.av_dict_parse_string(&d.c, ci, ck, cp, C.int(flags)))
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__dict.html#gae67f143237b2cb2936c9b147aa6dfde3
func (d *Dictionary) Get(key string, prev *DictionaryEntry, flags DictionaryFlags) *DictionaryEntry {
	ck := C.CString(key)
	defer C.free(unsafe.Pointer(ck))
	var cp *C.AVDictionaryEntry
	if prev != nil {
		cp = prev.c
	}
	if e := C.av_dict_get(d.c, ck, cp, C.int(flags)); e != nil {
		return newDictionaryEntryFromC(e)
	}
	return nil
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__dict.html#ga1bafd682b1fbb90e48a4cc3814b820f7
func (d *Dictionary) Free() {
	if d.c != nil {
		C.av_dict_free(&d.c)
	}
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__packet.html#ga2d2c8e143a2c98cf0aa31b072c286186
func (d *Dictionary) Pack() []byte {
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		return C.av_packet_pack_dictionary(d.c, size)
	})
}

// https://ffmpeg.org/doxygen/8.0/group__lavc__packet.html#gaae45c29cb3a29dc80b0b8f4ee9724492
func (d *Dictionary) Unpack(b []byte) error {
	return bytesToC(b, func(b *C.uint8_t, size C.size_t) error {
		return newError(C.av_packet_unpack_dictionary(b, size, &d.c))
	})
}

// https://ffmpeg.org/doxygen/8.0/group__lavu__dict.html#ga59a6372b124b306e3a2233723c5cdc78
func (d *Dictionary) Copy(dst *Dictionary, flags DictionaryFlags) error {
	return newError(C.av_dict_copy(&dst.c, d.c, C.int(flags)))
}
