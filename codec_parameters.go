package astiav

//#include <libavcodec/avcodec.h>
import "C"

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html
type CodecParameters struct {
	c *C.AVCodecParameters
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga647755ab2252e93221bb345f3d5e414f
func AllocCodecParameters() *CodecParameters {
	return newCodecParametersFromC(C.avcodec_parameters_alloc())
}

func newCodecParametersFromC(c *C.AVCodecParameters) *CodecParameters {
	if c == nil {
		return nil
	}
	return &CodecParameters{c: c}
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga950c8da55b8112077e640b6a0cb8cf36
func (cp *CodecParameters) Free() {
	C.avcodec_parameters_free(&cp.c)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a5268fcf4ae8ed27edef54f836b926d93
func (cp *CodecParameters) BitRate() int64 {
	return int64(cp.c.bit_rate)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a4d581c205b18108a4f00e8fb3a2b26f9
func (cp *CodecParameters) ChannelLayout() ChannelLayout {
	l, _ := newChannelLayoutFromC(&cp.c.ch_layout).clone()
	return l
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a4d581c205b18108a4f00e8fb3a2b26f9
func (cp *CodecParameters) SetChannelLayout(l ChannelLayout) {
	l.copy(&cp.c.ch_layout) //nolint: errcheck
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a9f76f2475ef24ff4c9771dd53072d040
func (cp *CodecParameters) CodecID() CodecID {
	return CodecID(cp.c.codec_id)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a9f76f2475ef24ff4c9771dd53072d040
func (cp *CodecParameters) SetCodecID(i CodecID) {
	cp.c.codec_id = uint32(i)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a9b6f7d220d100ba73defab295623356b
func (cp *CodecParameters) CodecTag() CodecTag {
	return CodecTag(cp.c.codec_tag)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a9b6f7d220d100ba73defab295623356b
func (cp *CodecParameters) SetCodecTag(t CodecTag) {
	cp.c.codec_tag = C.uint(t)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#ae4c7ac718a75adb31b5f2076a02fdedf
func (cp *CodecParameters) ChromaLocation() ChromaLocation {
	return ChromaLocation(cp.c.chroma_location)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#aa884cae3fd16b30c61201a686664f96b
func (cp *CodecParameters) ColorPrimaries() ColorPrimaries {
	return ColorPrimaries(cp.c.color_primaries)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#afa6744d9b8766db47a5ff7bddf0f2404
func (cp *CodecParameters) ColorRange() ColorRange {
	return ColorRange(cp.c.color_range)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#afa6744d9b8766db47a5ff7bddf0f2404
func (cp *CodecParameters) SetColorRange(r ColorRange) {
	cp.c.color_range = C.enum_AVColorRange(r)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a020398a4963e932853cefc169d90456d
func (cp *CodecParameters) ColorSpace() ColorSpace {
	return ColorSpace(cp.c.color_space)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a020398a4963e932853cefc169d90456d
func (cp *CodecParameters) SetColorSpace(s ColorSpace) {
	cp.c.color_space = C.enum_AVColorSpace(s)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#ac25ed8394e1efdbbcf28932ff0020893
func (cp *CodecParameters) ColorTransferCharacteristic() ColorTransferCharacteristic {
	return ColorTransferCharacteristic(cp.c.color_trc)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a9befe0b86412646017afb0051d144d13
func (cp *CodecParameters) ExtraData() []byte {
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		*size = C.size_t(cp.c.extradata_size)
		return cp.c.extradata
	})
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a9befe0b86412646017afb0051d144d13
func (cp *CodecParameters) SetExtraData(b []byte) error {
	return setBytesWithIntSizeInC(b, &cp.c.extradata, &cp.c.extradata_size)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a0ce9631123719789e4c7b0c23c66d534
func (cp *CodecParameters) FrameSize() int {
	return int(cp.c.frame_size)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a0ce9631123719789e4c7b0c23c66d534
func (cp *CodecParameters) SetFrameSize(i int) {
	cp.c.frame_size = C.int(i)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a1ec57ee84f19cf65d00eaa4d2a2253ce
func (cp *CodecParameters) Height() int {
	return int(cp.c.height)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a1ec57ee84f19cf65d00eaa4d2a2253ce
func (cp *CodecParameters) SetHeight(h int) {
	cp.c.height = C.int(h)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a26ae48eeaf8b315eca03b207e11edc7c
func (cp *CodecParameters) Level() Level {
	return Level(cp.c.level)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a26ae48eeaf8b315eca03b207e11edc7c
func (cp *CodecParameters) SetLevel(l Level) {
	cp.c.level = C.int(l)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a58369c3a8a986935b572df5aa6361ce2
func (cp *CodecParameters) MediaType() MediaType {
	return MediaType(cp.c.codec_type)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a58369c3a8a986935b572df5aa6361ce2
func (cp *CodecParameters) SetMediaType(t MediaType) {
	cp.c.codec_type = C.enum_AVMediaType(t)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#abee943e65d98f9763fa6602a356e774f
func (cp *CodecParameters) PixelFormat() PixelFormat {
	return PixelFormat(cp.c.format)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#abee943e65d98f9763fa6602a356e774f
func (cp *CodecParameters) SetPixelFormat(f PixelFormat) {
	cp.c.format = C.int(f)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a6b13b8a226ed923085718cd1323bfcb5
func (cp *CodecParameters) Profile() Profile {
	return Profile(cp.c.profile)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a6b13b8a226ed923085718cd1323bfcb5
func (cp *CodecParameters) SetProfile(p Profile) {
	cp.c.profile = C.int(p)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a7d6ef91120ffe80040c699e747a1ad68
func (cp *CodecParameters) SampleAspectRatio() Rational {
	return newRationalFromC(cp.c.sample_aspect_ratio)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a7d6ef91120ffe80040c699e747a1ad68
func (cp *CodecParameters) SetSampleAspectRatio(r Rational) {
	cp.c.sample_aspect_ratio = r.c
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#ad54da9241deabb3601e6e0e8fa832c19
func (cp *CodecParameters) SideData() *PacketSideData {
	return newPacketSideDataFromC(&cp.c.coded_side_data, &cp.c.nb_coded_side_data)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#abee943e65d98f9763fa6602a356e774f
func (cp *CodecParameters) SampleFormat() SampleFormat {
	return SampleFormat(cp.c.format)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#abee943e65d98f9763fa6602a356e774f
func (cp *CodecParameters) SetSampleFormat(f SampleFormat) {
	cp.c.format = C.int(f)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#abfc9b0aa975206f7e77a125e6b78536e
func (cp *CodecParameters) SampleRate() int {
	return int(cp.c.sample_rate)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#abfc9b0aa975206f7e77a125e6b78536e
func (cp *CodecParameters) SetSampleRate(r int) {
	cp.c.sample_rate = C.int(r)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a51639f88aef9f4f283f538a0c033fbb8
func (cp *CodecParameters) Width() int {
	return int(cp.c.width)
}

// https://ffmpeg.org/doxygen/7.0/structAVCodecParameters.html#a51639f88aef9f4f283f538a0c033fbb8
func (cp *CodecParameters) SetWidth(w int) {
	cp.c.width = C.int(w)
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga506c1c185ac48bb0086c61e267fc085c
func (cp *CodecParameters) FromCodecContext(cc *CodecContext) error {
	return newError(C.avcodec_parameters_from_context(cp.c, cc.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga8a4998c9d1695abb01d379539d313227
func (cp *CodecParameters) ToCodecContext(cc *CodecContext) error {
	return newError(C.avcodec_parameters_to_context(cc.c, cp.c))
}

// https://ffmpeg.org/doxygen/7.0/group__lavc__core.html#ga6d02e640ccc12c783841ce51d09b9fa7
func (cp *CodecParameters) Copy(dst *CodecParameters) error {
	return newError(C.avcodec_parameters_copy(dst.c, cp.c))
}
