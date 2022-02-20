package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestPixelFormat(t *testing.T) {
	p := astiav.FindPixelFormatByName("yuv420p")
	require.Equal(t, astiav.PixelFormatYuv420P, p)
	require.Equal(t, "yuv420p", p.String())
}
