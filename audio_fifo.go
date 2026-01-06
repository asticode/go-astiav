package astiav

//#include <libavutil/audio_fifo.h>
import "C"
import "unsafe"

// https://ffmpeg.org/doxygen/7.1/structAVAudioFifo.html
type AudioFifo struct {
	c *C.AVAudioFifo
}

func newAudioFifoFromC(c *C.AVAudioFifo) *AudioFifo {
	if c == nil {
		return nil
	}
	return &AudioFifo{c: c}
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__audiofifo.html#ga9d792394f0615a329aec47847f8f8784
func AllocAudioFifo(sampleFmt SampleFormat, channels int, nbSamples int) *AudioFifo {
	return newAudioFifoFromC(C.av_audio_fifo_alloc(C.enum_AVSampleFormat(sampleFmt), C.int(channels), C.int(nbSamples)))
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__audiofifo.html#ga27c1e16e5f09940d6016b1971c0b5742
func (a *AudioFifo) Realloc(nbSamples int) error {
	return newError(C.av_audio_fifo_realloc(a.c, C.int(nbSamples)))
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__audiofifo.html#gaa0a4742ecac52a999e8b4478d27f3b9b
func (a *AudioFifo) Size() int {
	return int(C.av_audio_fifo_size(a.c))
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__audiofifo.html#ga2bed2f01fe34228ee8a73617b3177d00
func (a *AudioFifo) Space() int {
	return int(C.av_audio_fifo_space(a.c))
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__audiofifo.html#ga51d81a165872919bbfdee3f00f6d6530
func (a *AudioFifo) Write(f *Frame) (int, error) {
	ret := C.av_audio_fifo_write(a.c, (*unsafe.Pointer)(unsafe.Pointer(&f.c.data[0])), C.int(f.NbSamples()))
	if err := newError(ret); err != nil {
		return 0, err
	}
	return int(ret), nil
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__audiofifo.html#ga5e2c87bbeefba0d229b4109b4b755529
func (a *AudioFifo) Read(f *Frame) (int, error) {
	ret := C.av_audio_fifo_read(a.c, (*unsafe.Pointer)(unsafe.Pointer(&f.c.data[0])), C.int(f.NbSamples()))
	if err := newError(ret); err != nil {
		return 0, err
	}
	return int(ret), nil
}

// https://ffmpeg.org/doxygen/7.1/group__lavu__audiofifo.html#ga74e029e47f7aa99217ad1f315c434875
func (a *AudioFifo) Free() {
	if a.c != nil {
		C.av_audio_fifo_free(a.c)
		a.c = nil
	}
}
