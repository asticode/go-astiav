package astiav

//#include <libavformat/avformat.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L937
type Stream struct {
	c *C.AVStream
}

func newStreamFromC(c *C.AVStream) *Stream {
	if c == nil {
		return nil
	}
	return &Stream{c: c}
}

func (s *Stream) AvgFrameRate() Rational {
	return newRationalFromC(s.c.avg_frame_rate)
}

func (s *Stream) SetAvgFrameRate(r Rational) {
	s.c.avg_frame_rate = r.c
}

func (s *Stream) CodecParameters() *CodecParameters {
	return newCodecParametersFromC(s.c.codecpar)
}

func (s *Stream) Duration() int64 {
	return int64(s.c.duration)
}

func (s *Stream) EventFlags() StreamEventFlags {
	return StreamEventFlags(s.c.event_flags)
}

func (s *Stream) ID() int {
	return int(s.c.id)
}

func (s *Stream) SetID(i int) {
	s.c.id = C.int(i)
}

func (s *Stream) Index() int {
	return int(s.c.index)
}

func (s *Stream) SetIndex(i int) {
	s.c.index = C.int(i)
}

func (s *Stream) Metadata() *Dictionary {
	return newDictionaryFromC(s.c.metadata)
}

func (s *Stream) NbFrames() int64 {
	return int64(s.c.nb_frames)
}

func (s *Stream) RFrameRate() Rational {
	return newRationalFromC(s.c.r_frame_rate)
}

func (s *Stream) SetRFrameRate(r Rational) {
	s.c.r_frame_rate = r.c
}

func (s *Stream) SampleAspectRatio() Rational {
	return newRationalFromC(s.c.sample_aspect_ratio)
}

func (s *Stream) SetSampleAspectRatio(r Rational) {
	s.c.sample_aspect_ratio = r.c
}

func (s *Stream) StartTime() int64 {
	return int64(s.c.start_time)
}

func (s *Stream) TimeBase() Rational {
	return newRationalFromC(s.c.time_base)
}

func (s *Stream) SetTimeBase(r Rational) {
	s.c.time_base = r.c
}
