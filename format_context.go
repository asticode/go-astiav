package astiav

//#include <libavcodec/avcodec.h>
//#include <libavformat/avformat.h>
import "C"
import (
	"math"
	"unsafe"
)

// https://github.com/FFmpeg/FFmpeg/blob/n5.0/libavformat/avformat.h#L1202
type FormatContext struct {
	c *C.AVFormatContext
}

func newFormatContextFromC(c *C.AVFormatContext) *FormatContext {
	if c == nil {
		return nil
	}
	fc := &FormatContext{c: c}
	classers.set(fc)
	return fc
}

var _ Classer = (*FormatContext)(nil)

func AllocFormatContext() *FormatContext {
	return newFormatContextFromC(C.avformat_alloc_context())
}

func AllocOutputFormatContext(of *OutputFormat, formatName, filename string) (*FormatContext, error) {
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
	var ofc *C.AVOutputFormat
	if of != nil {
		ofc = of.c
	}
	var fcc *C.AVFormatContext
	if err := newError(C.avformat_alloc_output_context2(&fcc, ofc, fonc, finc)); err != nil {
		return nil, err
	}
	return newFormatContextFromC(fcc), nil
}

func (fc *FormatContext) Free() {
	classers.del(fc)
	C.avformat_free_context(fc.c)
}

func (fc *FormatContext) BitRate() int64 {
	return int64(fc.c.bit_rate)
}

func (fc *FormatContext) Class() *Class {
	return newClassFromC(unsafe.Pointer(fc.c))
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

func (fc *FormatContext) Flags() FormatContextFlags {
	return FormatContextFlags(fc.c.flags)
}

func (fc *FormatContext) SetFlags(f FormatContextFlags) {
	fc.c.flags = C.int(f)
}

func (fc *FormatContext) SetInterruptCallback() IOInterrupter {
	i := newDefaultIOInterrupter()
	fc.c.interrupt_callback = i.c
	return i
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

func (fc *FormatContext) SetMetadata(d *Dictionary) {
	if d == nil {
		fc.c.metadata = nil
	} else {
		fc.c.metadata = d.c
	}
}

func (fc *FormatContext) NbPrograms() int {
	return int(fc.c.nb_programs)
}

func (fc *FormatContext) NbStreams() int {
	return int(fc.c.nb_streams)
}

func (fc *FormatContext) OutputFormat() *OutputFormat {
	return newOutputFormatFromC(fc.c.oformat)
}

func (fc *FormatContext) Pb() *IOContext {
	// If the io context has been created using the format context's OpenInput() method, we need to
	// make sure to return the same go struct as the one stored in classers
	if c, ok := classers.get(unsafe.Pointer(fc.c.pb)); ok {
		if v, ok := c.(*IOContext); ok {
			return v
		}
	}
	return newIOContextFromC(fc.c.pb)
}

func (fc *FormatContext) Programs() (ps []*Program) {
	pcs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVProgram)(nil))](*C.AVProgram))(unsafe.Pointer(fc.c.programs))
	for i := 0; i < fc.NbPrograms(); i++ {
		ps = append(ps, newProgramFromC(pcs[i], fc))
	}
	return
}

func (fc *FormatContext) SetPb(i *IOContext) {
	fc.c.pb = i.c
}

func (fc *FormatContext) StartTime() int64 {
	return int64(fc.c.start_time)
}

func (fc *FormatContext) Streams() (ss []*Stream) {
	scs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVStream)(nil))](*C.AVStream))(unsafe.Pointer(fc.c.streams))
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
	var urlc *C.char
	if url != "" {
		urlc = C.CString(url)
		defer C.free(unsafe.Pointer(urlc))
	}
	var dc **C.AVDictionary
	if d != nil {
		dc = &d.c
	}
	var fmtc *C.AVInputFormat
	if fmt != nil {
		fmtc = fmt.c
	}
	if err := newError(C.avformat_open_input(&fc.c, urlc, fmtc, dc)); err != nil {
		return err
	}
	if pb := fc.Pb(); pb != nil {
		classers.set(pb)
	}
	return nil
}

func (fc *FormatContext) CloseInput() {
	if pb := fc.Pb(); pb != nil {
		classers.del(pb)
	}
	classers.del(fc)
	if fc.c != nil {
		C.avformat_close_input(&fc.c)
	}
}

func (fc *FormatContext) NewProgram(id int) *Program {
	return newProgramFromC(C.av_new_program(fc.c, C.int(id)), fc)
}

func (fc *FormatContext) NewStream(c *Codec) *Stream {
	var cc *C.AVCodec
	if c != nil {
		cc = c.c
	}
	return newStreamFromC(C.avformat_new_stream(fc.c, cc))
}

func (fc *FormatContext) FindStreamInfo(d *Dictionary) error {
	var dc **C.AVDictionary
	if d != nil {
		dc = &d.c
	}
	return newError(C.avformat_find_stream_info(fc.c, dc))
}

func (fc *FormatContext) ReadFrame(p *Packet) error {
	var pc *C.AVPacket
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
	var dc **C.AVDictionary
	if d != nil {
		dc = &d.c
	}
	return newError(C.avformat_write_header(fc.c, dc))
}

func (fc *FormatContext) WriteFrame(p *Packet) error {
	var pc *C.AVPacket
	if p != nil {
		pc = p.c
	}
	return newError(C.av_write_frame(fc.c, pc))
}

func (fc *FormatContext) WriteInterleavedFrame(p *Packet) error {
	var pc *C.AVPacket
	if p != nil {
		pc = p.c
	}
	return newError(C.av_interleaved_write_frame(fc.c, pc))
}

func (fc *FormatContext) WriteTrailer() error {
	return newError(C.av_write_trailer(fc.c))
}

func (fc *FormatContext) GuessSampleAspectRatio(s *Stream, f *Frame) Rational {
	var cf *C.AVFrame
	if f != nil {
		cf = f.c
	}
	return newRationalFromC(C.av_guess_sample_aspect_ratio(fc.c, s.c, cf))
}

func (fc *FormatContext) GuessFrameRate(s *Stream, f *Frame) Rational {
	var cf *C.AVFrame
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
		fccs := []*C.AVFormatContext{}
		for _, fc := range fcs {
			fccs = append(fccs, fc.c)
		}
		return newError(C.av_sdp_create(&fccs[0], C.int(len(fcs)), buf, C.int(size)))
	})
}
