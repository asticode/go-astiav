package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSampleFormat(t *testing.T) {
	require.Equal(t, "s16", SampleFormatS16.String())
	require.Equal(t, 2, SampleFormatS16.BytesPerSample())
	require.False(t, SampleFormatS16.IsPlanar())
	require.True(t, SampleFormatS16P.IsPlanar())
}
