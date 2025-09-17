package astiav

//#include <libavcodec/avcodec.h>
//#include <libavutil/avutil.h>
//#include "subtitle.h"
import "C"
import (
	"unsafe"
)

// Subtitle represents an AVSubtitle
// https://ffmpeg.org/doxygen/8.0/structAVSubtitle.html
type Subtitle struct {
	c *C.AVSubtitle
}

func newSubtitleFromC(c *C.AVSubtitle) *Subtitle {
	if c == nil {
		return nil
	}
	return &Subtitle{c: c}
}

// AllocSubtitle allocates a new subtitle
func AllocSubtitle() *Subtitle {
	c := C.astiavSubtitleAlloc()
	if c == nil {
		return nil
	}
	return newSubtitleFromC(c)
}

// Free frees the subtitle
func (s *Subtitle) Free() {
	if s.c != nil {
		C.astiavSubtitleFree(s.c)
		C.av_free(unsafe.Pointer(s.c))
		s.c = nil
	}
}

// DecodeSubtitle2 decodes a subtitle
// https://ffmpeg.org/doxygen/8.0/group__lavc__decoding.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (cc *CodecContext) DecodeSubtitle2(subtitle *Subtitle, packet *Packet) (gotSubtitle bool, err error) {
	if cc.c == nil {
		return false, newError(-1)
	}
	if subtitle == nil || subtitle.c == nil {
		return false, newError(-1)
	}
	
	var pktC *C.AVPacket
	if packet != nil {
		pktC = packet.c
	}
	
	var gotSubtitleC C.int
	ret := C.astiavDecodeSubtitle2(cc.c, subtitle.c, &gotSubtitleC, pktC)
	if ret < 0 {
		return false, newError(ret)
	}
	
	return gotSubtitleC != 0, nil
}

// EncodeSubtitle encodes a subtitle
// https://ffmpeg.org/doxygen/8.0/group__lavc__encoding.html#ga8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b8b
func (cc *CodecContext) EncodeSubtitle(buf []byte, subtitle *Subtitle) (int, error) {
	if cc.c == nil {
		return 0, newError(-1)
	}
	if subtitle == nil || subtitle.c == nil {
		return 0, newError(-1)
	}
	if len(buf) == 0 {
		return 0, newError(-1)
	}
	
	ret := C.astiavEncodeSubtitle(cc.c, (*C.uint8_t)(unsafe.Pointer(&buf[0])), C.int(len(buf)), subtitle.c)
	if ret < 0 {
		return 0, newError(ret)
	}
	
	return int(ret), nil
}

// StartTime returns the subtitle start time in microseconds
func (s *Subtitle) StartTime() int64 {
	if s.c == nil {
		return 0
	}
	return int64(s.c.start_display_time)
}

// EndTime returns the subtitle end time in microseconds  
func (s *Subtitle) EndTime() int64 {
	if s.c == nil {
		return 0
	}
	return int64(s.c.end_display_time)
}

// NumRects returns the number of subtitle rectangles
func (s *Subtitle) NumRects() int {
	if s.c == nil {
		return 0
	}
	return int(s.c.num_rects)
}

// Format returns the subtitle format
func (s *Subtitle) Format() int {
	if s.c == nil {
		return 0
	}
	return int(s.c.format)
}

// Pts returns the subtitle presentation timestamp
func (s *Subtitle) Pts() int64 {
	if s.c == nil {
		return 0
	}
	return int64(s.c.pts)
}

// SetPts sets the subtitle presentation timestamp
func (s *Subtitle) SetPts(pts int64) {
	if s.c != nil {
		s.c.pts = C.int64_t(pts)
	}
}

// SetFormat sets the subtitle format
func (s *Subtitle) SetFormat(format int) {
	if s.c != nil {
		s.c.format = C.uint16_t(format)
	}
}

// SetStartDisplayTime sets the subtitle start display time
func (s *Subtitle) SetStartDisplayTime(time uint32) {
	if s.c != nil {
		s.c.start_display_time = C.uint32_t(time)
	}
}

// SetEndDisplayTime sets the subtitle end display time
func (s *Subtitle) SetEndDisplayTime(time uint32) {
	if s.c != nil {
		s.c.end_display_time = C.uint32_t(time)
	}
}

// StartDisplayTime returns the subtitle start display time
func (s *Subtitle) StartDisplayTime() uint32 {
	if s.c == nil {
		return 0
	}
	return uint32(s.c.start_display_time)
}

// EndDisplayTime returns the subtitle end display time
func (s *Subtitle) EndDisplayTime() uint32 {
	if s.c == nil {
		return 0
	}
	return uint32(s.c.end_display_time)
}

// NbRects returns the number of subtitle rectangles
func (s *Subtitle) NbRects() int {
	if s.c == nil {
		return 0
	}
	return int(s.c.num_rects)
}