package astiav

//#cgo pkg-config: libavformat
//#include <libavformat/avformat.h>
import "C"

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L937
type Stream struct {
	c *C.struct_AVStream
}

func newStreamFromC(c *C.struct_AVStream) *Stream {
	if c == nil {
		return nil
	}
	return &Stream{c: c}
}

func (s *Stream) AvgFrameRate() Rational {
	return newRationalFromC(s.c.avg_frame_rate)
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

func (s *Stream) Index() int {
	return int(s.c.index)
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

func (s *Stream) SampleAspectRatio() Rational {
	return newRationalFromC(s.c.sample_aspect_ratio)
}

func (s *Stream) SideData(t PacketSideDataType) []byte {
	return bytesFromC(func(size *cUlong) *C.uint8_t {
		return C.av_stream_get_side_data(s.c, (C.enum_AVPacketSideDataType)(t), size)
	})
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
