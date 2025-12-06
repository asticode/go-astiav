package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPixelFormat(t *testing.T) {
	t.Run("FindPixelFormatByName", func(t *testing.T) {
		p := FindPixelFormatByName("yuv420p")
		require.Equal(t, PixelFormatYuv420P, p)
		require.Equal(t, "yuv420p", p.String())
	})
	t.Run("Descriptor", func(t *testing.T) {
		d := PixelFormatYuv420P.Descriptor()
		require.NotNil(t, d)
		require.Equal(t, d.Name(), PixelFormatYuv420P.String())
	})
	t.Run("MarshalText", func(t *testing.T) {
		x1 := PixelFormatAbgr
		b, err := x1.MarshalText()
		require.NoError(t, err)
		var x2 PixelFormat
		require.Equal(t, "abgr", string(b))
		require.NoError(t, x2.UnmarshalText(b))
		require.Equal(t, x1, x2)
	})
}
