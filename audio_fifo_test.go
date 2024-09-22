package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAudioFIFO(t *testing.T) {
	afn := 2000
	af := AllocAudioFifo(SampleFormatFltp, 2, afn)
	defer af.Free()

	wn := 1024
	wf := AllocFrame()
	wf.SetNbSamples(wn)
	wf.SetChannelLayout(ChannelLayoutStereo)
	wf.SetSampleFormat(SampleFormatFltp)
	wf.SetSampleRate(48000)
	wf.AllocBuffer(0)

	rn := 120
	rf := AllocFrame()
	rf.SetNbSamples(rn)
	rf.SetChannelLayout(ChannelLayoutStereo)
	rf.SetSampleFormat(SampleFormatFltp)
	rf.SetSampleRate(48000)
	rf.AllocBuffer(0)

	w, err := af.Write(wf)
	require.NoError(t, err)
	require.Equal(t, wn, w)
	r, err := af.Read(rf)
	require.NoError(t, err)
	require.Equal(t, rn, r)
	require.Equal(t, wn-rn, af.Size())
	require.Equal(t, afn-af.Size(), af.Space())

	afn = 3000
	require.NoError(t, af.Realloc(afn))
	require.Equal(t, wn-rn, af.Size())
	require.Equal(t, afn-af.Size(), af.Space())
}
