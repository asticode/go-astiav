package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuffersrcFilterContextParameters(t *testing.T) {
	p := AllocBuffersrcFilterContextParameters()
	defer p.Free()
	p.SetChannelLayout(ChannelLayoutStereo)
	require.Equal(t, ChannelLayoutStereo, p.ChannelLayout())
	p.SetColorRange(ColorRangeMpeg)
	require.Equal(t, ColorRangeMpeg, p.ColorRange())
	p.SetColorSpace(ColorSpaceBt470Bg)
	require.Equal(t, ColorSpaceBt470Bg, p.ColorSpace())
	p.SetFramerate(NewRational(1, 2))
	require.Equal(t, NewRational(1, 2), p.Framerate())
	p.SetHeight(1)
	require.Equal(t, 1, p.Height())
	p.SetPixelFormat(PixelFormatRgba)
	require.Equal(t, PixelFormatRgba, p.PixelFormat())
	p.SetSampleAspectRatio(NewRational(3, 4))
	require.Equal(t, NewRational(3, 4), p.SampleAspectRatio())
	p.SetSampleFormat(SampleFormatDblp)
	require.Equal(t, SampleFormatDblp, p.SampleFormat())
	p.SetSampleRate(2)
	require.Equal(t, 2, p.SampleRate())
	p.SetTimeBase(NewRational(5, 6))
	require.Equal(t, NewRational(5, 6), p.TimeBase())
	p.SetWidth(3)
	require.Equal(t, 3, p.Width())
}
