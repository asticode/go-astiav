package astiav_test

import (
	"image"
	"reflect"
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

// Test constants for source and destination dimensions and formats
const (
	srcW = 100
	srcH = 100
	dstW = 200
	dstH = 200
)

// assertImageType is a helper function to check the type of an image.
func assertImageType(t *testing.T, img image.Image, expectedType reflect.Type) {
	actualType := reflect.TypeOf(img)
	require.Equal(t, expectedType, actualType, "Image type does not match")
}

// TestSWS tests the scaling functionality provided by the SWSContext.
func TestSWS(t *testing.T) {
	// Allocate and initialize source and destination frames
	srcFrame := astiav.AllocFrame()
	defer srcFrame.Free()
	dstFrame := astiav.AllocFrame()
	defer dstFrame.Free()

	srcFrame.SetHeight(srcH)
	srcFrame.SetWidth(srcW)
	srcFrame.SetPixelFormat(astiav.PixelFormatYuv420P)
	srcFrame.AllocBuffer(1)

	swsc := astiav.SwsGetContext(srcW, srcH, astiav.PixelFormatYuv420P, dstW, dstH, astiav.PixelFormatRgba, astiav.SWS_BILINEAR, dstFrame)
	require.NotNil(t, swsc)

	err := swsc.Scale(srcFrame, dstFrame)
	require.NoError(t, err)

	require.Equal(t, dstW, dstFrame.Height())
	require.Equal(t, dstH, dstFrame.Width())
	require.Equal(t, astiav.PixelFormatRgba, dstFrame.PixelFormat())

	// Convert frame data to image and perform additional verifications
	i1, err := dstFrame.Data().Image()
	require.NoError(t, err)
	require.Equal(t, dstW, i1.Bounds().Dx())
	require.Equal(t, dstH, i1.Bounds().Dy())
	assertImageType(t, i1, reflect.TypeOf((*image.NRGBA)(nil)))

	// Update sws ctx tests
	err = swsc.UpdateScalingParameters(50, 50, astiav.PixelFormatRgb24)
	require.NoError(t, err)
	require.Equal(t, astiav.PixelFormatRgb24, dstFrame.PixelFormat())
	err = swsc.Scale(srcFrame, dstFrame)
	require.NoError(t, err)
	require.Equal(t, dstFrame.Width(), 50)
	require.Equal(t, dstFrame.Height(), 50)

	swsc.Free()
}
