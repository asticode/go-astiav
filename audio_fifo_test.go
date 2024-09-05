package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAudioFIFO(t *testing.T) {
	audioFifo := AllocAudioFifo(
		SampleFormatFltp,
		2,
		960)

	writeSamples := 1024
	readSamples := 120
	writeFrame := AllocFrame()
	writeFrame.SetNbSamples(writeSamples)
	writeFrame.SetChannelLayout(ChannelLayoutStereo)
	writeFrame.SetSampleFormat(SampleFormatFltp)
	writeFrame.SetSampleRate(48000)
	writeFrame.AllocBuffer(0)

	readFrame := AllocFrame()
	readFrame.SetNbSamples(readSamples)
	readFrame.SetChannelLayout(ChannelLayoutStereo)
	readFrame.SetSampleFormat(SampleFormatFltp)
	readFrame.SetSampleRate(48000)
	readFrame.AllocBuffer(0)

	written := audioFifo.AudioFifoWrite(writeFrame.DataPtr(), writeFrame.NbSamples())
	require.Equal(t, writeSamples, written)
	read := audioFifo.AudioFifoRead(readFrame.DataPtr(), readFrame.NbSamples())
	require.Equal(t, readSamples, read)
	require.Equal(t, audioFifo.AudioFifoSize(), writeSamples-readSamples)
}
