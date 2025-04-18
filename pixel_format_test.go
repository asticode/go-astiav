package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPixelFormat(t *testing.T) {
	d := p.Descriptor()
	require.NotNil(t, d)
	require.Equal(t, d.Name(), p.String())
	p := FindPixelFormatByName("yuv420p")
	require.Equal(t, PixelFormatYuv420P, p)
	require.Equal(t, "yuv420p", p.String())
	t.Run("FindByName", func(t *testing.T) {
		p := FindPixelFormatByName("yuv420p")
		require.Equal(t, PixelFormatYuv420P, p)
		require.Equal(t, "yuv420p", p.String())
	})
	t.Run("Encode", func(t *testing.T) {
		x1 := PixelFormatAbgr
		s, err := x1.MarshalText()
		require.NoError(t, err)
		var x2 PixelFormat
		require.NoError(t, x2.UnmarshalText(s))
		require.Equal(t, x1, x2)
	})
}
