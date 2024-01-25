package astiav_test

import (
	"image"
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestSoftwareScaleContext(t *testing.T) {
	f1 := astiav.AllocFrame()
	require.NotNil(t, f1)
	defer f1.Free()

	f2 := astiav.AllocFrame()
	require.NotNil(t, f2)
	defer f2.Free()

	f3 := astiav.AllocFrame()
	require.NotNil(t, f3)
	defer f3.Free()

	srcW := 100
	srcH := 100
	srcPixelFormat := astiav.PixelFormatYuv420P
	dstW := 200
	dstH := 200
	dstPixelFormat := astiav.PixelFormatRgba

	f1.SetHeight(srcH)
	f1.SetWidth(srcW)
	f1.SetPixelFormat(srcPixelFormat)
	require.NoError(t, f1.AllocBuffer(1))

	swscf_1 := astiav.NewSoftwareScaleContextFlags(astiav.SoftwareScaleContextBilinear)
	swsc := astiav.NewSoftwareScaleContext(srcW, srcH, srcPixelFormat, dstW, dstH, dstPixelFormat, swscf_1)
	require.NotNil(t, swsc)
	require.Equal(t, swsc.Flags(), swscf_1)

	swscf_2 := astiav.NewSoftwareScaleContextFlags(astiav.SoftwareScaleContextPoint)
	swsc.SetFlags(swscf_2)
	require.Equal(t, swsc.Flags(), swscf_2)

	require.NoError(t, swsc.PrepareDestinationFrameForScaling(f2))
	require.Equal(t, dstH, swsc.ScaleFrame(f1, f2))

	require.Equal(t, dstW, f2.Height())
	require.Equal(t, dstH, f2.Width())
	require.Equal(t, dstPixelFormat, f2.PixelFormat())

	i1, err := f2.Data().Image()
	require.NoError(t, err)
	require.Equal(t, dstW, i1.Bounds().Dx())
	require.Equal(t, dstH, i1.Bounds().Dy())
	_, nrgbaOk := i1.(*image.NRGBA)
	require.True(t, nrgbaOk)

	dstW = 50
	dstH = 50
	dstPixelFormat = astiav.PixelFormatYuv420P

	require.NoError(t, swsc.SetSourceWidth(f2.Width()))
	require.Equal(t, swsc.SourceWidth(), f2.Width())
	require.NoError(t, swsc.SetSourceHeight(f2.Height()))
	require.Equal(t, swsc.SourceHeight(), f2.Height())
	require.NoError(t, swsc.SetSourcePixelFormat(f2.PixelFormat()))
	require.Equal(t, swsc.SourcePixelFormat(), f2.PixelFormat())

	require.NoError(t, swsc.SetDestinationWidth(dstW))
	require.Equal(t, swsc.DestinationWidth(), dstW)
	require.NoError(t, swsc.SetDestinationHeight(dstH))
	require.Equal(t, swsc.DestinationHeight(), dstH)
	require.NoError(t, swsc.SetDestinationPixelFormat(dstPixelFormat))
	require.Equal(t, swsc.DestinationPixelFormat(), dstPixelFormat)

	require.NoError(t, swsc.PrepareDestinationFrameForScaling(f3))
	require.Equal(t, f3.Height(), dstH)
	require.Equal(t, f3.Width(), dstW)
	require.Equal(t, f3.PixelFormat(), dstPixelFormat)
	require.Equal(t, dstH, swsc.ScaleFrame(f2, f3))
	require.Equal(t, dstW, f3.Height())
	require.Equal(t, dstH, f3.Width())
	require.Equal(t, dstPixelFormat, f3.PixelFormat())

	i2, err := f3.Data().Image()
	require.NoError(t, err)
	require.Equal(t, dstW, i2.Bounds().Dx())
	require.Equal(t, dstH, i2.Bounds().Dy())
	_, ycbcrOk := i2.(*image.YCbCr)
	require.True(t, ycbcrOk)

	defer swsc.Free()
}
