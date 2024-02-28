package astiav

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSoftwareScaleContext(t *testing.T) {
	f1 := AllocFrame()
	require.NotNil(t, f1)
	defer f1.Free()

	f2 := AllocFrame()
	require.NotNil(t, f2)
	defer f2.Free()

	f3 := AllocFrame()
	require.NotNil(t, f3)
	defer f3.Free()

	srcW := 4
	srcH := 2
	srcPixelFormat := PixelFormatYuv420P
	dstW := 8
	dstH := 4
	dstPixelFormat := PixelFormatRgba
	swscf1 := SoftwareScaleContextFlags(SoftwareScaleContextFlagBilinear)

	f1.SetHeight(srcH)
	f1.SetWidth(srcW)
	f1.SetPixelFormat(srcPixelFormat)
	require.NoError(t, f1.AllocBuffer(1))

	swsc1, err := CreateSoftwareScaleContext(srcW, srcH, srcPixelFormat, dstW, dstH, dstPixelFormat, swscf1)
	require.NoError(t, err)
	defer swsc1.Free()
	require.Equal(t, dstH, swsc1.DestinationHeight())
	require.Equal(t, dstPixelFormat, swsc1.DestinationPixelFormat())
	w, h := swsc1.DestinationResolution()
	require.Equal(t, w, dstW)
	require.Equal(t, h, dstH)
	require.Equal(t, dstW, swsc1.DestinationWidth())
	require.Equal(t, swscf1, swsc1.Flags())
	require.Equal(t, srcH, swsc1.SourceHeight())
	require.Equal(t, srcPixelFormat, swsc1.SourcePixelFormat())
	w, h = swsc1.SourceResolution()
	require.Equal(t, w, srcW)
	require.Equal(t, h, srcH)
	require.Equal(t, srcW, swsc1.SourceWidth())
	cl := swsc1.Class()
	require.NotNil(t, cl)
	require.Equal(t, "SWScaler", cl.Name())

	require.NoError(t, swsc1.ScaleFrame(f1, f2))
	require.Equal(t, dstH, f2.Height())
	require.Equal(t, dstW, f2.Width())
	require.Equal(t, dstPixelFormat, f2.PixelFormat())

	dstW = 4
	dstH = 3
	dstPixelFormat = PixelFormatYuv420P
	swscf2 := SoftwareScaleContextFlags(SoftwareScaleContextFlagPoint)
	srcW = 2
	srcH = 1
	srcPixelFormat = PixelFormatRgba

	require.NoError(t, swsc1.SetDestinationHeight(dstH))
	require.Equal(t, dstH, swsc1.DestinationHeight())
	require.NoError(t, swsc1.SetDestinationPixelFormat(dstPixelFormat))
	require.Equal(t, dstPixelFormat, swsc1.DestinationPixelFormat())
	require.NoError(t, swsc1.SetDestinationWidth(dstW))
	require.Equal(t, dstW, swsc1.DestinationWidth())
	dstW = 5
	dstH = 4
	require.NoError(t, swsc1.SetDestinationResolution(dstW, dstH))
	w, h = swsc1.DestinationResolution()
	require.Equal(t, w, dstW)
	require.Equal(t, h, dstH)
	require.NoError(t, swsc1.SetFlags(swscf2))
	require.Equal(t, swsc1.Flags(), swscf2)
	require.NoError(t, swsc1.SetSourceHeight(srcH))
	require.Equal(t, srcH, swsc1.SourceHeight())
	require.NoError(t, swsc1.SetSourcePixelFormat(srcPixelFormat))
	require.Equal(t, srcPixelFormat, swsc1.SourcePixelFormat())
	require.NoError(t, swsc1.SetSourceWidth(srcW))
	require.Equal(t, srcW, swsc1.SourceWidth())
	srcW = 3
	srcH = 2
	require.NoError(t, swsc1.SetSourceResolution(srcW, srcH))
	w, h = swsc1.SourceResolution()
	require.Equal(t, w, srcW)
	require.Equal(t, h, srcH)

	f4, err := globalHelper.inputLastFrame("image-rgba.png", MediaTypeVideo)
	require.NoError(t, err)

	f5 := AllocFrame()
	require.NotNil(t, f5)
	defer f5.Free()

	swsc2, err := CreateSoftwareScaleContext(f4.Width(), f4.Height(), f4.PixelFormat(), 512, 512, f4.PixelFormat(), NewSoftwareScaleContextFlags(SoftwareScaleContextFlagBilinear))
	require.NoError(t, err)
	require.NoError(t, swsc2.ScaleFrame(f4, f5))

	b1, err := f5.Data().Bytes(1)
	require.NoError(t, err)

	b2, err := os.ReadFile("testdata/image-rgba-upscaled-bytes")
	require.NoError(t, err)
	require.Equal(t, b2, b1)
}
