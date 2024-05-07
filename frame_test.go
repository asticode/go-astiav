package astiav

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestFrame(t *testing.T) {
	f1, err := globalHelper.inputLastFrame("video.mp4", MediaTypeVideo)
	require.NoError(t, err)
	// Should be "{384, 192, 192, 0, 0, 0, 0, 0}" but for some reason it"s "{320, 160, 160, 0, 0, 0, 0, 0}"
	// on darwin when testing using github
	require.Contains(t, [][8]int{
		{384, 192, 192, 0, 0, 0, 0, 0},
		{320, 160, 160, 0, 0, 0, 0, 0},
	}, f1.Linesize())
	require.Equal(t, int64(60928), f1.PktDts())
	require.Equal(t, unsafe.Pointer(f1.c), f1.UnsafePointer())

	f2 := AllocFrame()
	require.NotNil(t, f2)
	defer f2.Free()
	f2.SetChannelLayout(ChannelLayout21)
	f2.SetColorRange(ColorRangeJpeg)
	f2.SetHeight(2)
	f2.SetKeyFrame(true)
	f2.SetNbSamples(4)
	f2.SetPictureType(PictureTypeB)
	f2.SetPixelFormat(PixelFormat0Bgr)
	require.Equal(t, PixelFormat0Bgr, f2.PixelFormat()) // Need to test it right away as sample format actually updates the same field
	f2.SetPts(7)
	f2.SetSampleAspectRatio(NewRational(10, 2))
	f2.SetSampleFormat(SampleFormatDbl)
	require.Equal(t, SampleFormatDbl, f2.SampleFormat())
	f2.SetSampleRate(9)
	f2.SetWidth(10)
	require.True(t, f2.ChannelLayout().Equal(ChannelLayout21))
	require.Equal(t, ColorRangeJpeg, f2.ColorRange())
	require.Equal(t, 2, f2.Height())
	require.True(t, f2.KeyFrame())
	require.Equal(t, 4, f2.NbSamples())
	require.Equal(t, PictureTypeB, f2.PictureType())
	require.Equal(t, int64(7), f2.Pts())
	require.Equal(t, NewRational(10, 2), f2.SampleAspectRatio())
	require.Equal(t, 9, f2.SampleRate())
	require.Equal(t, 10, f2.Width())

	f3 := f1.Clone()
	require.NotNil(t, f3)
	defer f3.Free()
	require.Equal(t, 180, f3.Height())

	err = f2.AllocBuffer(0)
	require.NoError(t, err)
	err = f3.Ref(f2)
	require.NoError(t, err)
	require.Equal(t, 2, f3.Height())

	f3.MoveRef(f1)
	require.Equal(t, 180, f3.Height())
	require.Equal(t, 0, f1.Height())

	f3.Unref()
	require.Equal(t, 0, f3.Height())

	f4 := AllocFrame()
	require.NotNil(t, f4)
	defer f4.Free()
	f4.SetNbSamples(960)
	f4.SetChannelLayout(ChannelLayoutStereo)
	f4.SetSampleFormat(SampleFormatS16)
	f4.SetSampleRate(48000)
	err = f4.AllocBuffer(0)
	require.NoError(t, err)
	err = f4.AllocSamples(0)
	require.NoError(t, err)

	f5 := AllocFrame()
	require.NotNil(t, f5)
	defer f5.Free()
	sd := f5.NewSideData(FrameSideDataTypeAudioServiceType, 4)
	require.NotNil(t, sd)
	sd.SetData([]byte{1, 2, 3})
	sd = f5.SideData(FrameSideDataTypeAudioServiceType)
	require.NotNil(t, sd)
	require.Equal(t, FrameSideDataTypeAudioServiceType, sd.Type())
	require.True(t, bytes.HasPrefix(sd.Data(), []byte{1, 2, 3}))
	require.Len(t, sd.Data(), 4)
	sd.SetData([]byte{1, 2, 3, 4, 5})
	sd = f5.SideData(FrameSideDataTypeAudioServiceType)
	require.NotNil(t, sd)
	require.Equal(t, []byte{1, 2, 3, 4}, sd.Data())

	f6 := AllocFrame()
	require.NotNil(t, f6)
	defer f6.Free()
	f6.SetColorRange(ColorRangeUnspecified)
	f6.SetHeight(2)
	f6.SetPixelFormat(PixelFormatYuv420P)
	f6.SetWidth(4)
	const align = 1
	require.NoError(t, f6.AllocBuffer(align))
	require.NoError(t, f6.AllocImage(align))
	require.NoError(t, f6.ImageFillBlack())
	n, err := f6.ImageBufferSize(align)
	require.NoError(t, err)
	require.Equal(t, 12, n)
	b := make([]byte, n)
	n, err = f6.ImageCopyToBuffer(b, align)
	require.NoError(t, err)
	require.Equal(t, 12, n)
	require.Equal(t, []byte{0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x80, 0x80, 0x80, 0x80}, b)
}
