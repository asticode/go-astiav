package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestMediaType(t *testing.T) {
	require.Equal(t, "video", astiav.MediaTypeVideo.String())
}
