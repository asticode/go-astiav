package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestCodec(t *testing.T) {
	c := astiav.FindDecoder(astiav.CodecIDMp3)
	require.NotNil(t, c)
	require.Nil(t, c.ChannelLayouts())
	require.True(t, c.IsDecoder())
	require.False(t, c.IsEncoder())
	require.Nil(t, c.PixelFormats())
	require.Equal(t, []astiav.SampleFormat{astiav.SampleFormatFltp, astiav.SampleFormatFlt}, c.SampleFormats())
	require.Equal(t, "mp3float", c.Name())
	require.Equal(t, "mp3float", c.String())

	c = astiav.FindDecoderByName("aac")
	require.NotNil(t, c)
	require.Equal(t, []astiav.ChannelLayout{
		astiav.ChannelLayoutMono,
		astiav.ChannelLayoutStereo,
		astiav.ChannelLayoutSurround,
		astiav.ChannelLayout4Point0,
		astiav.ChannelLayout5Point0Back,
		astiav.ChannelLayout5Point1Back,
		astiav.ChannelLayout7Point1WideBack,
	}, c.ChannelLayouts())
	require.True(t, c.IsDecoder())
	require.False(t, c.IsEncoder())
	require.Equal(t, []astiav.SampleFormat{astiav.SampleFormatFltp}, c.SampleFormats())
	require.Equal(t, "aac", c.Name())
	require.Equal(t, "aac", c.String())

	c = astiav.FindEncoder(astiav.CodecIDH264)
	require.NotNil(t, c)
	require.False(t, c.IsDecoder())
	require.True(t, c.IsEncoder())
	require.Equal(t, []astiav.PixelFormat{
		astiav.PixelFormatVideotoolbox,
		astiav.PixelFormatNv12,
		astiav.PixelFormatYuv420P,
	}, c.PixelFormats())
	require.Nil(t, c.SampleFormats())
	require.Equal(t, "h264_videotoolbox", c.Name())
	require.Equal(t, "h264_videotoolbox", c.String())

	c = astiav.FindEncoderByName("mjpeg")
	require.NotNil(t, c)
	require.False(t, c.IsDecoder())
	require.True(t, c.IsEncoder())
	require.Equal(t, []astiav.PixelFormat{
		astiav.PixelFormatYuvj420P,
		astiav.PixelFormatYuvj422P,
		astiav.PixelFormatYuvj444P,
	}, c.PixelFormats())
	require.Equal(t, "mjpeg", c.Name())
	require.Equal(t, "mjpeg", c.String())

	c = astiav.FindDecoderByName("invalid")
	require.Nil(t, c)
}
