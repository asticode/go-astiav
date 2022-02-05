package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestPixelFormat(t *testing.T) {
	require.Equal(t, "yuv420p", astiav.PixelFormatYuv420P.String())
}
