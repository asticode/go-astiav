package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChannelLayout(t *testing.T) {
	cl := ChannelLayoutStereo
	require.Equal(t, 2, cl.Channels())
	require.Equal(t, "stereo", cl.String())
	require.True(t, cl.Valid())
	require.True(t, cl.Equal(ChannelLayoutStereo))
	require.False(t, cl.Equal(ChannelLayoutMono))
}
