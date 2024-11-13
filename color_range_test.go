package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColorRange(t *testing.T) {
	require.Equal(t, "tv", ColorRangeMpeg.Name())
	require.Equal(t, "tv", ColorRangeMpeg.String())
}
