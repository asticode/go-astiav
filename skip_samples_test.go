package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSkipSamples(t *testing.T) {
	_, err := newSkipSamplesFromBytes([]byte("123456789"))
	require.Error(t, err)
	ss1 := &SkipSamples{
		ReasonEnd:   1,
		ReasonStart: 2,
		SkipEnd:     3,
		SkipStart:   4,
	}
	b1 := ss1.bytes()
	require.Equal(t, []byte{0x4, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x2, 0x1}, b1)
	ss2, err := newSkipSamplesFromBytes(b1)
	require.NoError(t, err)
	require.Equal(t, ss1, ss2)
}
