package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestCodec(t *testing.T) {
	c := astiav.FindDecoder(astiav.CodecIDMp3)
	require.NotNil(t, c)
	require.Equal(t, c.ID(), astiav.CodecIDMp3)
	require.Nil(t, c.ChannelLayouts())
	require.True(t, c.IsDecoder())
	require.False(t, c.IsEncoder())
	require.Nil(t, c.PixelFormats())
	require.Equal(t, []astiav.SampleFormat{astiav.SampleFormatFltp, astiav.SampleFormatFlt}, c.SampleFormats())
	require.Equal(t, "mp3float", c.Name())
	require.Equal(t, "mp3float", c.String())

	c = astiav.FindDecoderByName("aac")
	require.NotNil(t, c)
	els := []astiav.ChannelLayout{
		astiav.ChannelLayoutMono,
		astiav.ChannelLayoutStereo,
		astiav.ChannelLayoutSurround,
		astiav.ChannelLayout4Point0,
		astiav.ChannelLayout5Point0Back,
		astiav.ChannelLayout5Point1Back,
		astiav.ChannelLayout7Point1WideBack,
	}
	gls := c.ChannelLayouts()
	require.Len(t, gls, len(els))
	for idx := range els {
		require.True(t, els[idx].Equal(gls[idx]))
	}
	require.True(t, c.IsDecoder())
	require.False(t, c.IsEncoder())
	require.Equal(t, []astiav.SampleFormat{astiav.SampleFormatFltp}, c.SampleFormats())
	require.Equal(t, "aac", c.Name())
	require.Equal(t, "aac", c.String())

	c = astiav.FindEncoder(astiav.CodecIDMjpeg)
	require.NotNil(t, c)
	require.False(t, c.IsDecoder())
	require.True(t, c.IsEncoder())
	require.Contains(t, c.PixelFormats(), astiav.PixelFormatYuvj420P)
	require.Nil(t, c.SampleFormats())
	require.Contains(t, c.Name(), "mjpeg")
	require.Contains(t, c.String(), "mjpeg")

	c = astiav.FindEncoderByName("mjpeg")
	require.NotNil(t, c)
	require.False(t, c.IsDecoder())
	require.True(t, c.IsEncoder())
	require.Equal(t, []astiav.PixelFormat{
		astiav.PixelFormatYuvj420P,
		astiav.PixelFormatYuvj422P,
		astiav.PixelFormatYuvj444P,
		astiav.PixelFormatYuv420P,
		astiav.PixelFormatYuv422P,
		astiav.PixelFormatYuv444P,
	}, c.PixelFormats())
	require.Equal(t, "mjpeg", c.Name())
	require.Equal(t, "mjpeg", c.String())

	c = astiav.FindDecoderByName("invalid")
	require.Nil(t, c)

	var found bool
	for _, c := range astiav.Codecs() {
		if c.ID() == astiav.CodecIDMjpeg {
			found = true
		}
	}
	require.True(t, found)
}
