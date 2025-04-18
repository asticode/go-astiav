package astiav

//#include <libavutil/samplefmt.h>
//#include <libavutil/avutil.h>
import "C"
import (
	"encoding"
	"fmt"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/group__lavu__sampfmts.html#gaf9a51ca15301871723577c730b5865c5
type SampleFormat C.enum_AVSampleFormat

var (
	_ encoding.TextMarshaler   = SampleFormatNone
	_ encoding.TextUnmarshaler = (*SampleFormat)(nil)
)

const (
	SampleFormatDbl  = SampleFormat(C.AV_SAMPLE_FMT_DBL)
	SampleFormatDblp = SampleFormat(C.AV_SAMPLE_FMT_DBLP)
	SampleFormatFlt  = SampleFormat(C.AV_SAMPLE_FMT_FLT)
	SampleFormatFltp = SampleFormat(C.AV_SAMPLE_FMT_FLTP)
	SampleFormatNb   = SampleFormat(C.AV_SAMPLE_FMT_NB)
	SampleFormatNone = SampleFormat(C.AV_SAMPLE_FMT_NONE)
	SampleFormatS16  = SampleFormat(C.AV_SAMPLE_FMT_S16)
	SampleFormatS16P = SampleFormat(C.AV_SAMPLE_FMT_S16P)
	SampleFormatS32  = SampleFormat(C.AV_SAMPLE_FMT_S32)
	SampleFormatS32P = SampleFormat(C.AV_SAMPLE_FMT_S32P)
	SampleFormatS64  = SampleFormat(C.AV_SAMPLE_FMT_S64)
	SampleFormatS64P = SampleFormat(C.AV_SAMPLE_FMT_S64P)
	SampleFormatU8   = SampleFormat(C.AV_SAMPLE_FMT_U8)
	SampleFormatU8P  = SampleFormat(C.AV_SAMPLE_FMT_U8P)
)

// https://ffmpeg.org/doxygen/7.0/group__lavu__sampfmts.html#ga31b9d149b2de9821a65f4f5612970838
func (f SampleFormat) Name() string {
	return C.GoString(C.av_get_sample_fmt_name((C.enum_AVSampleFormat)(f)))
}

func (f SampleFormat) String() string {
	return f.Name()
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__sampfmts.html#ga0c3c218e1dd570ad4917c69a35a6c77d
func (f SampleFormat) BytesPerSample() int {
	return int(C.av_get_bytes_per_sample((C.enum_AVSampleFormat)(f)))
}

// https://ffmpeg.org/doxygen/7.0/group__lavu__sampfmts.html#ga06ba8a64dc4382c422789a5d0b6bf592
func (f SampleFormat) IsPlanar() bool {
	return C.av_sample_fmt_is_planar((C.enum_AVSampleFormat)(f)) > 0
}

func FindSampleFormat(name string) SampleFormat {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return SampleFormat(C.av_get_sample_fmt(cname))
}

func (f SampleFormat) MarshalText() ([]byte, error) { return ([]byte)(f.Name()), nil }
func (f *SampleFormat) UnmarshalText(d []byte) error {
	s := string(d)
	if s == "" {
		*f = SampleFormatNone
	}
	pf := FindSampleFormat(s)
	if pf == SampleFormatNone {
		return fmt.Errorf("astiav: invalid SampleFormat: %q", s)
	}
	*f = pf
	return nil
}
