package astiav

//#cgo pkg-config: libavutil
//#include <libavutil/audio_fifo.h>
//#include <stdlib.h>
import "C"
import "unsafe"

type AudioFifo struct {
	c *C.struct_AVAudioFifo
}

func newAudioFifoFromC(c *C.struct_AVAudioFifo) *AudioFifo {
	if c == nil {
		return nil
	}
	return &AudioFifo{c: c}
}

func AllocAudioFifo(sampleFmt SampleFormat, channels int, nbSamples int) *AudioFifo {
	return newAudioFifoFromC(C.av_audio_fifo_alloc(C.enum_AVSampleFormat(sampleFmt), C.int(channels), C.int(nbSamples)))
}

func (a *AudioFifo) AudioFifoRealloc(nbSamples int) int {
	return int(C.av_audio_fifo_realloc((*C.struct_AVAudioFifo)(a.c), C.int(nbSamples)))
}

func (a *AudioFifo) AudioFifoSize() int {
	return int(C.av_audio_fifo_size((*C.struct_AVAudioFifo)(a.c)))
}

func (a *AudioFifo) AudioFifoSpace() int {
	return int(C.av_audio_fifo_space((*C.struct_AVAudioFifo)(a.c)))
}

func (a *AudioFifo) AudioFifoWrite(data **uint8, nbSamples int) int {
	return int(C.av_audio_fifo_write((*C.struct_AVAudioFifo)(a.c), (*unsafe.Pointer)(unsafe.Pointer(data)), C.int(nbSamples)))
}

func (a *AudioFifo) AudioFifoRead(data **uint8, nbSamples int) int {
	return int(C.av_audio_fifo_read((*C.struct_AVAudioFifo)(a.c), (*unsafe.Pointer)(unsafe.Pointer(data)), C.int(nbSamples)))
}

func (a *AudioFifo) AudioFifoFree() {
	C.av_audio_fifo_free((*C.struct_AVAudioFifo)(a.c))
}
