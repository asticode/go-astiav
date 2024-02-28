package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMediaType(t *testing.T) {
	require.Equal(t, "video", MediaTypeVideo.String())
}
