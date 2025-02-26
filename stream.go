package astiav

//#include <libavformat/avformat.h>
import "C"
import "unsafe"

// https://ffmpeg.org/doxygen/7.0/structAVStream.html
type Stream struct {
	c *C.AVStream
}

func newStreamFromC(c *C.AVStream) *Stream {
	if c == nil {
		return nil
	}
	return &Stream{c: c}
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a946e1e9b89eeeae4cab8a833b482c1ad
func (s *Stream) AvgFrameRate() Rational {
	return newRationalFromC(s.c.avg_frame_rate)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a946e1e9b89eeeae4cab8a833b482c1ad
func (s *Stream) SetAvgFrameRate(r Rational) {
	s.c.avg_frame_rate = r.c
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a4737d8b012827558f55a6f559b253496
func (s *Stream) Class() *Class {
	if s.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(s.c))
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a12826d21779289356722971d362c583c
func (s *Stream) CodecParameters() *CodecParameters {
	return newCodecParametersFromC(s.c.codecpar)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a492fcecc45dbbd8da51edd0124e9dd30
func (s *Stream) Discard() Discard {
	return Discard(s.c.discard)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a492fcecc45dbbd8da51edd0124e9dd30
func (s *Stream) SetDiscard(d Discard) {
	s.c.discard = C.enum_AVDiscard(d)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a4e04af7a5a4d8298649850df798dd0bc
func (s *Stream) Duration() int64 {
	return int64(s.c.duration)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#ab76e176c2a1d1ff09ec9c0bb88dc25e9
func (s *Stream) EventFlags() StreamEventFlags {
	return StreamEventFlags(s.c.event_flags)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#ab76e176c2a1d1ff09ec9c0bb88dc25e9
func (s *Stream) SetEventFlags(eventFlags StreamEventFlags) {
	s.c.event_flags = C.int(eventFlags)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a6873ed62f196c24e8bf282609231786f
func (s *Stream) ID() int {
	return int(s.c.id)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a6873ed62f196c24e8bf282609231786f
func (s *Stream) SetID(i int) {
	s.c.id = C.int(i)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a6ca823054632821e085377f7d371a2d1
func (s *Stream) Index() int {
	return int(s.c.index)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a6ca823054632821e085377f7d371a2d1
func (s *Stream) SetIndex(i int) {
	s.c.index = C.int(i)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html
func (s *Stream) Metadata() *Dictionary {
	return newDictionaryFromC(s.c.metadata)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html
func (s *Stream) SetMetadata(d *Dictionary) {
	if d == nil {
		s.c.metadata = nil
	} else {
		s.c.metadata = d.c
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a4382c3064df1c9eb232ac198dec067f9
func (s *Stream) NbFrames() int64 {
	return int64(s.c.nb_frames)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a6cdb0c90a69899f4e1e54704bb654936
func (s *Stream) PTSWrapBits() int {
	return int(s.c.pts_wrap_bits)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a6cdb0c90a69899f4e1e54704bb654936
func (s *Stream) SetPTSWrapBits(bits int) {
	s.c.pts_wrap_bits = C.int(bits)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#ad63fb11cc1415e278e09ddc676e8a1ad
func (s *Stream) RFrameRate() Rational {
	return newRationalFromC(s.c.r_frame_rate)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#ad63fb11cc1415e278e09ddc676e8a1ad
func (s *Stream) SetRFrameRate(r Rational) {
	s.c.r_frame_rate = r.c
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a3f19c60ac6da237cd10e4d97150c118e
func (s *Stream) SampleAspectRatio() Rational {
	return newRationalFromC(s.c.sample_aspect_ratio)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a3f19c60ac6da237cd10e4d97150c118e
func (s *Stream) SetSampleAspectRatio(r Rational) {
	s.c.sample_aspect_ratio = r.c
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a7c67ae70632c91df8b0f721658ec5377
func (s *Stream) StartTime() int64 {
	return int64(s.c.start_time)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a7c67ae70632c91df8b0f721658ec5377
func (s *Stream) SetStartTime(startTime int64) {
	s.c.start_time = C.int64_t(startTime)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a9db755451f14e2bf590d4b85d82b32e6
func (s *Stream) TimeBase() Rational {
	return newRationalFromC(s.c.time_base)
}

// https://ffmpeg.org/doxygen/7.0/structAVStream.html#a9db755451f14e2bf590d4b85d82b32e6
func (s *Stream) SetTimeBase(r Rational) {
	s.c.time_base = r.c
}
