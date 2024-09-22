package astiav

//#include <libavutil/audio_fifo.h>
import "C"
import "unsafe"

// https://github.com/FFmpeg/FFmpeg/blob/n7.0/libavutil/audio_fifo.c#L37
type AudioFifo struct {
	c *C.AVAudioFifo
}

func newAudioFifoFromC(c *C.AVAudioFifo) *AudioFifo {
	if c == nil {
		return nil
	}
	return &AudioFifo{c: c}
}

func AllocAudioFifo(sampleFmt SampleFormat, channels int, nbSamples int) *AudioFifo {
	return newAudioFifoFromC(C.av_audio_fifo_alloc(C.enum_AVSampleFormat(sampleFmt), C.int(channels), C.int(nbSamples)))
}

func (a *AudioFifo) Realloc(nbSamples int) error {
	return newError(C.av_audio_fifo_realloc(a.c, C.int(nbSamples)))
}

func (a *AudioFifo) Size() int {
	return int(C.av_audio_fifo_size(a.c))
}

func (a *AudioFifo) Space() int {
	return int(C.av_audio_fifo_space(a.c))
}

func (a *AudioFifo) Write(f *Frame) (int, error) {
	ret := C.av_audio_fifo_write(a.c, (*unsafe.Pointer)(unsafe.Pointer(&f.c.data[0])), C.int(f.NbSamples()))
	if err := newError(ret); err != nil {
		return 0, err
	}
	return int(ret), nil
}

func (a *AudioFifo) Read(f *Frame) (int, error) {
	ret := C.av_audio_fifo_read(a.c, (*unsafe.Pointer)(unsafe.Pointer(&f.c.data[0])), C.int(f.NbSamples()))
	if err := newError(ret); err != nil {
		return 0, err
	}
	return int(ret), nil
}

func (a *AudioFifo) Free() {
	C.av_audio_fifo_free(a.c)
}
