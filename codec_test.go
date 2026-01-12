package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodec(t *testing.T) {
	c := FindDecoder(CodecIDMp3)
	require.NotNil(t, c)
	require.Equal(t, c.ID(), CodecIDMp3)
	require.Nil(t, c.SupportedChannelLayouts())
	require.True(t, c.IsDecoder())
	require.False(t, c.IsEncoder())
	require.Nil(t, c.SupportedPixelFormats())
	require.Equal(t, []SampleFormat{SampleFormatFltp, SampleFormatFlt}, c.SupportedSampleFormats())
	require.Equal(t, "mp3float", c.Name())
	require.Equal(t, "mp3float", c.String())

	c = FindDecoderByName("aac")
	require.NotNil(t, c)
	els := []ChannelLayout{
		ChannelLayoutMono,
		ChannelLayoutStereo,
		ChannelLayoutSurround,
		ChannelLayout4Point0,
		ChannelLayout5Point0Back,
		ChannelLayout5Point1Back,
		ChannelLayout7Point1WideBack,
		ChannelLayout6Point1Back,
		ChannelLayout7Point1,
		ChannelLayout22Point2,
		ChannelLayout5Point1Point2Back,
	}
	gls := c.SupportedChannelLayouts()
	require.Len(t, gls, len(els))
	for idx := range els {
		require.True(t, els[idx].Equal(gls[idx]))
	}
	require.True(t, c.IsDecoder())
	require.False(t, c.IsEncoder())
	require.Equal(t, []SampleFormat{SampleFormatFltp}, c.SupportedSampleFormats())
	require.Equal(t, "aac", c.Name())
	require.Equal(t, "aac", c.String())

	c = FindEncoderByName("aac")
	require.NotNil(t, c)
	require.Equal(t, []int{96000, 88200, 64000, 48000, 44100, 32000, 24000, 22050, 16000, 12000, 11025, 8000, 7350}, c.SupportedSampleRates())

	c = FindEncoder(CodecIDMjpeg)
	require.NotNil(t, c)
	require.False(t, c.IsDecoder())
	require.True(t, c.IsEncoder())
	require.Contains(t, c.SupportedPixelFormats(), PixelFormatYuvj420P)
	require.Nil(t, c.SupportedSampleFormats())
	require.Contains(t, c.Name(), "mjpeg")
	require.Contains(t, c.String(), "mjpeg")

	c = FindEncoderByName("mjpeg")
	require.NotNil(t, c)
	require.False(t, c.IsDecoder())
	require.True(t, c.IsEncoder())
	require.Equal(t, []PixelFormat{
		PixelFormatYuvj420P,
		PixelFormatYuvj422P,
		PixelFormatYuvj444P,
		PixelFormatYuv420P,
		PixelFormatYuv422P,
		PixelFormatYuv444P,
	}, c.SupportedPixelFormats())
	require.Equal(t, "mjpeg", c.Name())
	require.Equal(t, "mjpeg", c.String())
	require.Equal(t, []ColorRange{ColorRangeJpeg}, c.SupportedColorRanges())
	require.Nil(t, c.SupportedColorSpaces())

	c = FindEncoderByName("mpeg1video")
	require.NotNil(t, c)
	require.Equal(t, []Rational{
		NewRational(24000, 1001),
		NewRational(24, 1),
		NewRational(25, 1),
		NewRational(30000, 1001),
		NewRational(30, 1),
		NewRational(50, 1),
		NewRational(60000, 1001),
		NewRational(60, 1),
		NewRational(15, 1),
		NewRational(5, 1),
		NewRational(10, 1),
		NewRational(12, 1),
		NewRational(15, 1),
	}, c.SupportedFrameRates())

	c = FindDecoderByName("invalid")
	require.Nil(t, c)

	var found bool
	for _, c := range Codecs() {
		if c.ID() == CodecIDMjpeg {
			found = true
		}
	}
	require.True(t, found)
}
