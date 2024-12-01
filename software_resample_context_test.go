package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSoftwareResampleContext(t *testing.T) {
	src := AllocSoftwareResampleContext()
	defer src.Free()

	f1, err := globalHelper.inputLastFrame("video.mp4", MediaTypeAudio, nil)
	require.NoError(t, err)

	f2 := AllocFrame()
	defer f2.Free()
	f2.SetChannelLayout(ChannelLayoutMono)
	f2.SetNbSamples(300)
	f2.SetSampleFormat(SampleFormatS16)
	f2.SetSampleRate(24000)
	require.NoError(t, f2.AllocBuffer(0))
	require.NoError(t, f2.AllocSamples(0))

	for _, v := range []struct {
		expectedDelay     int64
		expectedNbSamples int
		f                 *Frame
	}{
		{
			expectedDelay:     212,
			expectedNbSamples: 300,
			f:                 f1,
		},
		{
			expectedDelay:     17,
			expectedNbSamples: 212,
		},
		{expectedDelay: 17},
	} {
		require.NoError(t, src.ConvertFrame(v.f, f2))
		require.Equal(t, v.expectedNbSamples, f2.NbSamples())
		require.Equal(t, v.expectedDelay, src.Delay(int64(f2.SampleRate())))
	}
}
