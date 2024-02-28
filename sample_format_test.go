package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSampleFormat(t *testing.T) {
	require.Equal(t, "s16", SampleFormatS16.String())
}
