package astiav

//#cgo pkg-config: libavcodec libavformat
//#include <libavcodec/avcodec.h>
//#include <libavformat/avformat.h>
/*
int astiavInterruptCallback(void *ret)
{
    return *((int*)ret);
}
AVIOInterruptCB astiavNewInterruptCallback(int *ret)
{
	AVIOInterruptCB c = { astiavInterruptCallback, ret };
	return c;
}
*/
import "C"
import (
	"math"
	"unsafe"
)

const (
	maxArraySize = math.MaxInt32 - 1
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L1202
type FormatContext struct {
	c *C.struct_AVFormatContext
}

func newFormatContext() *FormatContext {
	return &FormatContext{}
}

func newFormatContextFromC(c *C.struct_AVFormatContext) *FormatContext {
	if c == nil {
		return nil
	}
	return &FormatContext{c: c}
}

func AllocFormatContext() *FormatContext {
	return newFormatContextFromC(C.avformat_alloc_context())
}

func AllocOutputFormatContext(o *OutputFormat, formatName, filename string) (*FormatContext, error) {
	fonc := (*C.char)(nil)
	if len(formatName) > 0 {
		fonc = C.CString(formatName)
		defer C.free(unsafe.Pointer(fonc))
	}
	finc := (*C.char)(nil)
	if len(filename) > 0 {
		finc = C.CString(filename)
		defer C.free(unsafe.Pointer(finc))
	}
	fc := newFormatContext()
	var oc *C.struct_AVOutputFormat
	if o != nil {
		oc = o.c
	}
	err := newError(C.avformat_alloc_output_context2(&fc.c, oc, fonc, finc))
	return fc, err
}

func (fc *FormatContext) DumpFormat(index int, url string, isOutput bool) {
	curl := (*C.char)(nil)
	if len(url) > 0 {
		curl = C.CString(url)
		defer C.free(unsafe.Pointer(curl))
	}
	outputC := C.int(0)
	if isOutput {
		outputC = C.int(1)
	}
	C.av_dump_format(fc.c, C.int(index), curl, outputC)
}

func (fc *FormatContext) BitRate() int64 {
	return int64(fc.c.bit_rate)
}

func (fc *FormatContext) CtxFlags() FormatContextCtxFlags {
	return FormatContextCtxFlags(fc.c.ctx_flags)
}

func (fc *FormatContext) Duration() int64 {
	return int64(fc.c.duration)
}

func (fc *FormatContext) EventFlags() FormatEventFlags {
	return FormatEventFlags(fc.c.event_flags)
}

func (fc *FormatContext) Filename() string {
	return C.GoString(&fc.c.filename[0])
}

func (fc *FormatContext) Flags() FormatContextFlags {
	return FormatContextFlags(fc.c.flags)
}

func (fc *FormatContext) SetFlags(f FormatContextFlags) {
	fc.c.flags = C.int(f)
}

func (fc *FormatContext) SetInterruptCallback() *int {
	ret := 0
	fc.c.interrupt_callback = C.astiavNewInterruptCallback((*C.int)(unsafe.Pointer(&ret)))
	return &ret
}

func (fc *FormatContext) InputFormat() *InputFormat {
	return newInputFormatFromC(fc.c.iformat)
}

func (fc *FormatContext) IOFlags() IOContextFlags {
	return IOContextFlags(fc.c.avio_flags)
}

func (fc *FormatContext) MaxAnalyzeDuration() int64 {
	return int64(fc.c.max_analyze_duration)
}

func (fc *FormatContext) Metadata() *Dictionary {
	return newDictionaryFromC(fc.c.metadata)
}

func (fc *FormatContext) NbStreams() int {
	return int(fc.c.nb_streams)
}

func (fc *FormatContext) OutputFormat() *OutputFormat {
	return newOutputFormatFromC(fc.c.oformat)
}

func (fc *FormatContext) Pb() *IOContext {
	if fc.c == nil {
		return nil
	}
	return newIOContextFromC(fc.c.pb)
}

func (fc *FormatContext) SetPb(i *IOContext) {
	fc.c.pb = i.c
}

func (fc *FormatContext) StartTime() int64 {
	return int64(fc.c.start_time)
}

func (fc *FormatContext) Streams() (ss []*Stream) {
	scs := (*[maxArraySize](*C.struct_AVStream))(unsafe.Pointer(fc.c.streams))
	for i := 0; i < fc.NbStreams(); i++ {
		ss = append(ss, newStreamFromC(scs[i]))
	}
	return
}

func (fc *FormatContext) StrictStdCompliance() StrictStdCompliance {
	return StrictStdCompliance(fc.c.strict_std_compliance)
}

func (fc *FormatContext) SetStrictStdCompliance(strictStdCompliance StrictStdCompliance) {
	fc.c.strict_std_compliance = C.int(strictStdCompliance)
}

func (fc *FormatContext) OpenInput(url string, fmt *InputFormat, d *Dictionary) error {
	urlc := C.CString(url)
	defer C.free(unsafe.Pointer(urlc))
	var dc **C.struct_AVDictionary
	if d != nil {
		dc = &d.c
	}
	var fmtc *C.struct_AVInputFormat
	if fmt != nil {
		fmtc = fmt.c
	}
	return newError(C.avformat_open_input(&fc.c, urlc, fmtc, dc))
}

func (fc *FormatContext) CloseInput() {
	C.avformat_close_input(&fc.c)
}

func (fc *FormatContext) Free() {
	if fc.c != nil {
		C.avformat_free_context(fc.c)
	}
}

func (fc *FormatContext) NewStream(c *Codec) *Stream {
	var cc *C.struct_AVCodec
	if c != nil {
		cc = c.c
	}
	return newStreamFromC(C.avformat_new_stream(fc.c, cc))
}

func (fc *FormatContext) FindStreamInfo(d *Dictionary) error {
	var dc **C.struct_AVDictionary
	if d != nil {
		dc = &d.c
	}
	return newError(C.avformat_find_stream_info(fc.c, dc))
}

func (fc *FormatContext) ReadFrame(p *Packet) error {
	var pc *C.struct_AVPacket
	if p != nil {
		pc = p.c
	}
	return newError(C.av_read_frame(fc.c, pc))
}

func (fc *FormatContext) SeekFrame(streamIndex int, timestamp int64, f SeekFlags) error {
	return newError(C.av_seek_frame(fc.c, C.int(streamIndex), C.int64_t(timestamp), C.int(f)))
}

func (fc *FormatContext) Flush() error {
	return newError(C.avformat_flush(fc.c))
}

func (fc *FormatContext) WriteHeader(d *Dictionary) error {
	var dc **C.struct_AVDictionary
	if d != nil {
		dc = &d.c
	}
	if fc.c == nil {
		panic("nil format context")
	}
	if fc.c.pb == nil {
		panic("nil format context")
	}
	return newError(C.avformat_write_header(fc.c, dc))
}

func (fc *FormatContext) WriteFrame(p *Packet) error {
	var pc *C.struct_AVPacket
	if p != nil {
		pc = p.c
	}
	return newError(C.av_write_frame(fc.c, pc))
}

func (fc *FormatContext) WriteInterleavedFrame(p *Packet) error {
	var pc *C.struct_AVPacket
	if p != nil {
		pc = p.c
	}
	return newError(C.av_interleaved_write_frame(fc.c, pc))
}

func (fc *FormatContext) WriteTrailer() error {
	return newError(C.av_write_trailer(fc.c))
}

func (fc *FormatContext) GuessSampleAspectRatio(s *Stream, f *Frame) Rational {
	var cf *C.struct_AVFrame
	if f != nil {
		cf = f.c
	}
	return newRationalFromC(C.av_guess_sample_aspect_ratio(fc.c, s.c, cf))
}

func (fc *FormatContext) GuessFrameRate(s *Stream, f *Frame) Rational {
	var cf *C.struct_AVFrame
	if f != nil {
		cf = f.c
	}
	return newRationalFromC(C.av_guess_frame_rate(fc.c, s.c, cf))
}

func (fc *FormatContext) SDPCreate() (string, error) {
	return sdpCreate([]*FormatContext{fc})
}

func sdpCreate(fcs []*FormatContext) (string, error) {
	return stringFromC(1024, func(buf *C.char, size C.size_t) error {
		fccs := []*C.struct_AVFormatContext{}
		for _, fc := range fcs {
			fccs = append(fccs, fc.c)
		}
		return newError(C.av_sdp_create(&fccs[0], C.int(len(fcs)), buf, C.int(size)))
	})
}
