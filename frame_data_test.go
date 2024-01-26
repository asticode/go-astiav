package astiav_test

import (
	"image/png"
	"os"
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestFrameData(t *testing.T) {
	const (
		name = "image-rgba"
		ext  = "png"
	)
	f, err := globalHelper.inputLastFrame(name+"."+ext, astiav.MediaTypeVideo)
	require.NoError(t, err)
	fd := f.Data()

	b1, err := fd.Bytes(1)
	require.NoError(t, err)

	b2, err := os.ReadFile("testdata/" + name + "-bytes")
	require.NoError(t, err)
	require.Equal(t, b1, b2)

	f1, err := os.Open("testdata/" + name + "." + ext)
	require.NoError(t, err)
	defer f1.Close()

	i1, err := fd.GuessImageFormat()
	require.NoError(t, err)
	require.NoError(t, err)
	require.NoError(t, fd.ToImage(i1))
	i2, err := png.Decode(f1)
	require.NoError(t, err)
	require.Equal(t, i1, i2)
}
