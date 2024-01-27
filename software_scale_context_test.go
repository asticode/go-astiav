package astiav_test

import (
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

	srcW := 320
	srcH := 280
	srcPixelFormat := astiav.PixelFormatYuv420P
	dstW := 640
	dstH := 480
	dstPixelFormat := astiav.PixelFormatRgba
	swscf1 := astiav.SoftwareScaleContextFlags(astiav.SoftwareScaleContextFlagBilinear)

	f1.SetHeight(srcH)
	f1.SetWidth(srcW)
	f1.SetPixelFormat(srcPixelFormat)
	require.NoError(t, f1.AllocBuffer(1))

	swsc, err := astiav.CreateSoftwareScaleContext(srcW, srcH, srcPixelFormat, dstW, dstH, dstPixelFormat, swscf1)
	require.NoError(t, err)
	defer swsc.Free()

	require.NoError(t, swsc.ScaleFrame(f1, f2))
	require.Equal(t, dstH, f2.Height())
	require.Equal(t, dstW, f2.Width())
	require.Equal(t, dstPixelFormat, f2.PixelFormat())

	dstW = 1024
	dstH = 576
	dstPixelFormat = astiav.PixelFormatYuv420P
	swscf2 := astiav.SoftwareScaleContextFlags(astiav.SoftwareScaleContextFlagPoint)

	require.Equal(t, swsc.Flags(), swscf1)
	swsc.SetFlags(swscf2)
	require.Equal(t, swsc.Flags(), swscf2)

	require.NoError(t, swsc.SetSourceWidth(f2.Width()))
	require.Equal(t, swsc.SourceWidth(), f2.Width())
	require.NoError(t, swsc.SetSourceHeight(f2.Height()))
	require.Equal(t, swsc.SourceHeight(), f2.Height())
	require.NoError(t, swsc.SetSourceResolution(1280, 720))
	w, h := swsc.SourceResolution()
	require.Equal(t, w, 1280)
	require.Equal(t, h, 720)
	require.NoError(t, swsc.SetSourceResolution(f2.Width(), f2.Height()))
	require.NoError(t, swsc.SetSourcePixelFormat(f2.PixelFormat()))
	require.Equal(t, swsc.SourcePixelFormat(), f2.PixelFormat())

	require.NoError(t, swsc.SetDestinationWidth(dstW))
	require.Equal(t, swsc.DestinationWidth(), dstW)
	require.NoError(t, swsc.SetDestinationHeight(dstH))
	require.Equal(t, swsc.DestinationHeight(), dstH)
	require.NoError(t, swsc.SetDestinationResolution(800, 600))
	w, h = swsc.DestinationResolution()
	require.Equal(t, w, 800)
	require.Equal(t, h, 600)
	require.NoError(t, swsc.SetDestinationResolution(dstW, dstH))
	require.NoError(t, swsc.SetDestinationPixelFormat(dstPixelFormat))
	require.Equal(t, swsc.DestinationPixelFormat(), dstPixelFormat)

	require.NoError(t, swsc.ScaleFrame(f2, f3))
	require.Equal(t, dstW, f3.Width())
	require.Equal(t, dstH, f3.Height())
	require.Equal(t, dstPixelFormat, f3.PixelFormat())

}
