package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAudioFIFO(t *testing.T) {
	af := AllocAudioFifo(
		SampleFormatFltp,
		2,
		960)
	defer af.Free()
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

	written, err := af.Write(writeFrame)
	require.NoError(t, err)
	require.Equal(t, writeSamples, written)
	read, err := af.Read(readFrame)
	require.NoError(t, err)
	require.Equal(t, readSamples, read)
	require.Equal(t, af.Size(), writeSamples-readSamples)
	reallocSamples := 3000
	err = af.Realloc(reallocSamples)
	require.NoError(t, err)
	expectedAfSize := writeSamples - readSamples
	require.Equal(t, af.Space(), reallocSamples-expectedAfSize)
	// It still has the same amount of data
	require.Equal(t, af.Size(), expectedAfSize)
}
