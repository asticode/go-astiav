package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodecID(t *testing.T) {
	require.Equal(t, MediaTypeVideo, CodecIDH264.MediaType())
	require.Equal(t, "h264", CodecIDH264.Name())
	require.Equal(t, "h264", CodecIDH264.String())
}
