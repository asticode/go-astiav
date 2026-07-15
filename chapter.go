package astiav

//#include <libavformat/avformat.h>
import "C"

// https://ffmpeg.org/doxygen/8.0/structAVChapter.html
type Chapter struct {
	c *C.AVChapter
}

func newChapterFromC(c *C.AVChapter) *Chapter {
	if c == nil {
		return nil
	}
	return &Chapter{c: c}
}

// https://ffmpeg.org/doxygen/8.0/structAVChapter.html#a067d8f6d76affe31403f99877b1c94bb
func (c *Chapter) ID() int64 {
	return int64(c.c.id)
}

// https://ffmpeg.org/doxygen/8.0/structAVChapter.html#a067d8f6d76affe31403f99877b1c94bb
func (c *Chapter) SetID(i int64) {
	c.c.id = C.int64_t(i)
}

// https://ffmpeg.org/doxygen/8.0/structAVChapter.html#acb5cb6ce9bb6b9f4b970a919f4899818
func (c *Chapter) TimeBase() Rational {
	return newRationalFromC(c.c.time_base)
}

// https://ffmpeg.org/doxygen/8.0/structAVChapter.html#acb5cb6ce9bb6b9f4b970a919f4899818
func (c *Chapter) SetTimeBase(r Rational) {
	c.c.time_base = r.c
}

// https://ffmpeg.org/doxygen/8.0/structAVChapter.html#a499a5062224e22249be6f2d16f74c449
func (c *Chapter) Start() int64 {
	return int64(c.c.start)
}

// https://ffmpeg.org/doxygen/8.0/structAVChapter.html#a499a5062224e22249be6f2d16f74c449
func (c *Chapter) SetStart(start int64) {
	c.c.start = C.int64_t(start)
}

// https://ffmpeg.org/doxygen/8.0/structAVChapter.html#ab68f32dc07fac89b4364e86483b00f3e
func (c *Chapter) End() int64 {
	return int64(c.c.end)
}

// https://ffmpeg.org/doxygen/8.0/structAVChapter.html#ab68f32dc07fac89b4364e86483b00f3e
func (c *Chapter) SetEnd(end int64) {
	c.c.end = C.int64_t(end)
}

// https://ffmpeg.org/doxygen/8.0/structAVChapter.html#a998ea9c9f86547970d58f0b405d55332
func (c *Chapter) Metadata() *Dictionary {
	return newDictionaryFromC(c.c.metadata)
}

// https://ffmpeg.org/doxygen/8.0/structAVChapter.html#a998ea9c9f86547970d58f0b405d55332
func (c *Chapter) SetMetadata(d *Dictionary) {
	if d == nil {
		c.c.metadata = nil
	} else {
		c.c.metadata = d.c
	}
}
