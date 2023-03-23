package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestChannelLayout(t *testing.T) {
	cl := astiav.ChannelLayoutStereo
	require.Equal(t, 2, cl.NbChannels())
	require.Equal(t, "stereo", cl.String())
	require.True(t, cl.Valid())
	require.True(t, cl.Equal(astiav.ChannelLayoutStereo))
	require.False(t, cl.Equal(astiav.ChannelLayoutMono))
}
