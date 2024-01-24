package astiav_test

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestFrameData(t *testing.T) {
	for _, v := range []struct {
		ext  string
		i    image.Image
		name string
	}{
		{
			ext:  "png",
			i:    &image.NRGBA{},
			name: "image-rgba",
		},
		// TODO Find a way to test yuv and yuva even though result seems to change randomly
	} {
		// We use a closure to ease closing files
		func() {
			f, err := globalHelper.inputLastFrame(v.name+"."+v.ext, astiav.MediaTypeVideo)
			require.NoError(t, err)
			fd := f.Data()

			b1, err := fd.Bytes(1)
			require.NoError(t, err)

			b2, err := os.ReadFile("testdata/" + v.name + "-bytes")
			require.NoError(t, err)
			require.Equal(t, b1, b2)

			f1, err := os.Open("testdata/" + v.name + "." + v.ext)
			require.NoError(t, err)
			defer f1.Close()

			i1, err := fd.Image()
			require.NoError(t, err)
			require.NoError(t, fd.ToImage(v.i))
			i2, err := png.Decode(f1)
			require.NoError(t, err)
			require.Equal(t, i1, i2)
			require.Equal(t, v.i, i2)
		}()
	}
}
