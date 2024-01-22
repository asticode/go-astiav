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
	srcW       = 100
	srcH       = 100
	dstW       = 200
	dstH       = 200
	secondDstW = 300
	secondDstH = 300
	srcFormat  = astiav.PixelFormatYuv420P
	dstFormat  = astiav.PixelFormatRgba
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
	srcFrame.SetPixelFormat(srcFormat)
	srcFrame.AllocBuffer(1)
	srcFrame.ImageFillBlack() // Fill the source frame with black for testing

	// Create SWSContext for scaling and verify it's not nil
	swsc := astiav.AllocSwsContext(srcW, srcH, srcFormat, dstW, dstH, dstFormat, astiav.SWS_BILINEAR, dstFrame)
	require.NotNil(t, swsc)

	// Perform scaling and verify no errors
	err := swsc.Scale(srcFrame, dstFrame)
	require.NoError(t, err)

	// Verify the dimensions and format of the destination frame
	require.Equal(t, dstW, dstFrame.Height())
	require.Equal(t, dstH, dstFrame.Width())
	require.Equal(t, dstFormat, dstFrame.PixelFormat())

	// Convert frame data to image and perform additional verifications
	i1, err := dstFrame.Data().Image()
	require.NoError(t, err)
	require.Equal(t, dstW, i1.Bounds().Dx())
	require.Equal(t, dstH, i1.Bounds().Dy())
	assertImageType(t, i1, reflect.TypeOf((*image.NRGBA)(nil)))
	swsc.Free()
}
