package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPixelFormat(t *testing.T) {
	p := FindPixelFormatByName("yuv420p")
	require.Equal(t, PixelFormatYuv420P, p)
	require.Equal(t, "yuv420p", p.String())
	d := p.Descriptor()
	require.NotNil(t, d)
	require.Equal(t, d.Name(), p.String())
}
