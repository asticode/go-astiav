package astiav

//#include <libavcodec/avcodec.h>
//#include <libavformat/avformat.h>
import "C"
import (
	"fmt"
	"math"
	"unsafe"
)

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html
type FormatContext struct {
	classerHandler
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

// https://ffmpeg.org/doxygen/7.0/group__lavf__core.html#gac7a91abf2f59648d995894711f070f62
func AllocFormatContext() *FormatContext {
	return newFormatContextFromC(C.avformat_alloc_context())
}

// https://ffmpeg.org/doxygen/7.0/avformat_8h.html#af5930942120e38a4766dc0bb9e4cae74
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

// https://ffmpeg.org/doxygen/7.0/group__lavf__core.html#gac2990b13b68e831a408fce8e1d0d6445
func (fc *FormatContext) Free() {
	if fc.c != nil {
		// Make sure to clone the classer before freeing the object since
		// the C free method may reset the pointer
		c := newClonedClasser(fc)
		C.avformat_free_context(fc.c)
		fc.c = nil
		// Make sure to remove from classers after freeing the object since
		// the C free method may use methods needing the classer
		if c != nil {
			classers.del(c)
		}
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a972a02b9e3b542a426e323a8f8e3ea41
func (fc *FormatContext) BitRate() int64 {
	return int64(fc.c.bit_rate)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a0c396740b9a2487aa57d4352d2dc1687
func (fc *FormatContext) Class() *Class {
	if fc.c == nil {
		return nil
	}
	return newClassFromC(unsafe.Pointer(fc.c))
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a4e6076343df1ffc2e16cedbba3f3f397
func (fc *FormatContext) CtxFlags() FormatContextCtxFlags {
	return FormatContextCtxFlags(fc.c.ctx_flags)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#ad0ea78ac48f5bb0a15a15c1c472744d9
func (fc *FormatContext) Duration() int64 {
	return int64(fc.c.duration)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a0302506d4b3434da77b8b3db43821aa0
func (fc *FormatContext) EventFlags() FormatEventFlags {
	return FormatEventFlags(fc.c.event_flags)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a32379cc371463b235d54235d4af06a15
func (fc *FormatContext) Flags() FormatContextFlags {
	return FormatContextFlags(fc.c.flags)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a32379cc371463b235d54235d4af06a15
func (fc *FormatContext) SetFlags(f FormatContextFlags) {
	fc.c.flags = C.int(f)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a5b37acfe4024d92ee510064e80920b40
func (fc *FormatContext) SetIOInterrupter(i *IOInterrupter) {
	if i == nil {
		fc.c.interrupt_callback = C.AVIOInterruptCB{}
	} else {
		fc.c.interrupt_callback = *i.c
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a6c01f25ef062e0398b0b55dd337246ed
func (fc *FormatContext) InputFormat() *InputFormat {
	return newInputFormatFromC(fc.c.iformat)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a5e6814c9de3c272396f07e2ff18c7b27
func (fc *FormatContext) IOFlags() IOContextFlags {
	return IOContextFlags(fc.c.avio_flags)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a4d860662c014f88277c8f20e238fa694
func (fc *FormatContext) MaxAnalyzeDuration() int64 {
	return int64(fc.c.max_analyze_duration)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a3019a56080ed2e3297ff25bc2ff88adf
func (fc *FormatContext) Metadata() *Dictionary {
	return newDictionaryFromC(fc.c.metadata)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a3019a56080ed2e3297ff25bc2ff88adf
func (fc *FormatContext) SetMetadata(d *Dictionary) {
	if d == nil {
		fc.c.metadata = nil
	} else {
		fc.c.metadata = d.c
	}
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a58c8c4d0ea974e0fcb0ce06fb1174f9f
func (fc *FormatContext) NbPrograms() int {
	return int(fc.c.nb_programs)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a4c2c5a4c758966349ff513e95154d062
func (fc *FormatContext) Programs() (ps []*Program) {
	pcs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVProgram)(nil))](*C.AVProgram))(unsafe.Pointer(fc.c.programs))
	for i := 0; i < fc.NbPrograms(); i++ {
		ps = append(ps, newProgramFromC(pcs[i], fc))
	}
	return
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a0b748d924898b08b89ff4974afd17285
func (fc *FormatContext) NbStreams() int {
	return int(fc.c.nb_streams)
}

func (fc *FormatContext) Streams() (ss []*Stream) {
	scs := (*[(math.MaxInt32 - 1) / unsafe.Sizeof((*C.AVStream)(nil))](*C.AVStream))(unsafe.Pointer(fc.c.streams))
	for i := 0; i < fc.NbStreams(); i++ {
		ss = append(ss, newStreamFromC(scs[i]))
	}
	return
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a37ba86cd5630097cdae01afbc2b40743
func (fc *FormatContext) OutputFormat() *OutputFormat {
	return newOutputFormatFromC(fc.c.oformat)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a1e7324262b6b78522e52064daaa7bc87
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

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a1e7324262b6b78522e52064daaa7bc87
func (fc *FormatContext) SetPb(i *IOContext) {
	fc.c.pb = i.c
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#ac4c0777e54085af2f3f1b27130e2b21b
func (fc *FormatContext) PrivateData() *PrivateData {
	return newPrivateDataFromC(fc.c.priv_data)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a2590129e00adfa726ab2033a10e905e9
func (fc *FormatContext) StartTime() int64 {
	return int64(fc.c.start_time)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a5017684cf0a84c990f60c8d50adec144
func (fc *FormatContext) StrictStdCompliance() StrictStdCompliance {
	return StrictStdCompliance(fc.c.strict_std_compliance)
}

// https://ffmpeg.org/doxygen/7.0/structAVFormatContext.html#a5017684cf0a84c990f60c8d50adec144
func (fc *FormatContext) SetStrictStdCompliance(strictStdCompliance StrictStdCompliance) {
	fc.c.strict_std_compliance = C.int(strictStdCompliance)
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__decoding.html#gac05d61a2b492ae3985c658f34622c19d
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
	if err := fc.newError(C.avformat_open_input(&fc.c, urlc, fmtc, dc)); err != nil {
		return err
	}
	if pb := fc.Pb(); pb != nil {
		classers.set(pb)
	}
	return nil
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__decoding.html#gae804b99aec044690162b8b9b110236a4
func (fc *FormatContext) CloseInput() {
	if fc.c != nil {
		// Make sure to clone the classer before freeing the object since
		// the C free method may reset the pointer
		c := newClonedClasser(fc)
		var cpb *ClonedClasser
		if pb := fc.Pb(); pb != nil {
			cpb = newClonedClasser(pb)
		}
		C.avformat_close_input(&fc.c)
		// Make sure to remove from classers after freeing the object since
		// the C free method may use methods needing the classer
		if cpb != nil {
			classers.del(cpb)
		}
		if c != nil {
			classers.del(c)
		}
	}
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__core.html#gab31f7c7c99dcadead38e8e83e0fdb828
func (fc *FormatContext) NewProgram(id int) *Program {
	return newProgramFromC(C.av_new_program(fc.c, C.int(id)), fc)
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__core.html#gaf2c94216a6a19144e86cac843a0a4409
func (fc *FormatContext) NewStream(c *Codec) *Stream {
	var cc *C.AVCodec
	if c != nil {
		cc = c.c
	}
	return newStreamFromC(C.avformat_new_stream(fc.c, cc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__decoding.html#gad42172e27cddafb81096939783b157bb
func (fc *FormatContext) FindStreamInfo(d *Dictionary) error {
	var dc **C.AVDictionary
	if d != nil {
		dc = &d.c
	}
	return fc.newError(C.avformat_find_stream_info(fc.c, dc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__decoding.html#ga4fdb3084415a82e3810de6ee60e46a61
func (fc *FormatContext) ReadFrame(p *Packet) error {
	var pc *C.AVPacket
	if p != nil {
		pc = p.c
	}
	return fc.newError(C.av_read_frame(fc.c, pc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__decoding.html#gaa23f7619d8d4ea0857065d9979c75ac8
func (fc *FormatContext) SeekFrame(streamIndex int, timestamp int64, f SeekFlags) error {
	return fc.newError(C.av_seek_frame(fc.c, C.int(streamIndex), C.int64_t(timestamp), C.int(f)))
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__decoding.html#gaa03a82c5fd4fe3af312d229ca94cd6f3
func (fc *FormatContext) Flush() error {
	return fc.newError(C.avformat_flush(fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__encoding.html#ga18b7b10bb5b94c4842de18166bc677cb
func (fc *FormatContext) WriteHeader(d *Dictionary) error {
	var dc **C.AVDictionary
	if d != nil {
		dc = &d.c
	}
	return fc.newError(C.avformat_write_header(fc.c, dc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__encoding.html#gaa85cc1774f18f306cd20a40fc50d0b36
func (fc *FormatContext) WriteFrame(p *Packet) error {
	var pc *C.AVPacket
	if p != nil {
		pc = p.c
	}
	return fc.newError(C.av_write_frame(fc.c, pc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__encoding.html#ga37352ed2c63493c38219d935e71db6c1
func (fc *FormatContext) WriteInterleavedFrame(p *Packet) error {
	var pc *C.AVPacket
	if p != nil {
		pc = p.c
	}
	return fc.newError(C.av_interleaved_write_frame(fc.c, pc))
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__encoding.html#ga7f14007e7dc8f481f054b21614dfec13
func (fc *FormatContext) WriteTrailer() error {
	return fc.newError(C.av_write_trailer(fc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__misc.html#gafa6fbfe5c1bf6792fd6e33475b6056bd
func (fc *FormatContext) GuessSampleAspectRatio(s *Stream, f *Frame) Rational {
	var cf *C.AVFrame
	if f != nil {
		cf = f.c
	}
	return newRationalFromC(C.av_guess_sample_aspect_ratio(fc.c, s.c, cf))
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__misc.html#ga698e6aa73caa9616851092e2be15875d
func (fc *FormatContext) GuessFrameRate(s *Stream, f *Frame) Rational {
	var cf *C.AVFrame
	if f != nil {
		cf = f.c
	}
	return newRationalFromC(C.av_guess_frame_rate(fc.c, s.c, cf))
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__misc.html#gaa2a7353a6bb0c8726797abd56b176af0
func (fc *FormatContext) SDPCreate() (string, error) {
	return stringFromC(1024, func(buf *C.char, size C.size_t) error {
		fccs := []*C.AVFormatContext{fc.c}
		return fc.newError(C.av_sdp_create(&fccs[0], C.int(len(fccs)), buf, C.int(size)))
	})
}

// https://ffmpeg.org/doxygen/7.0/avformat_8c.html#a8d4609a8f685ad894c1503ffd1b610b4
func (fc *FormatContext) FindBestStream(mt MediaType, wantedStreamIndex, relatedStreamIndex int) (*Stream, *Codec, error) {
	// Find best stream
	var cCodec *C.AVCodec
	ret := C.av_find_best_stream(fc.c, C.enum_AVMediaType(mt), C.int(wantedStreamIndex), C.int(relatedStreamIndex), &cCodec, 0)
	if err := fc.newError(ret); err != nil {
		return nil, nil, err
	}

	// Loop through streams
	for _, s := range fc.Streams() {
		if s.Index() == int(ret) {
			return s, newCodecFromC(cCodec), nil
		}
	}
	return nil, nil, fmt.Errorf("astiav: no stream with index %d", ret)
}

// https://ffmpeg.org/doxygen/7.0/group__lavf__misc.html#gae2645941f2dc779c307eb6314fd39f10
func (fc *FormatContext) Dump(streamIndex int, url string, isOutput bool) {
	curl := (*C.char)(nil)
	if len(url) > 0 {
		curl = C.CString(url)
		defer C.free(unsafe.Pointer(curl))
	}
	cisOutput := 0
	if isOutput {
		cisOutput = 1
	}
	C.av_dump_format(fc.c, C.int(streamIndex), curl, C.int(cisOutput))
}
