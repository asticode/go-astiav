package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChannelLayout(t *testing.T) {
	cl1 := ChannelLayoutStereo
	require.Equal(t, 2, cl1.Channels())
	require.Equal(t, "stereo", cl1.String())
	require.True(t, cl1.Valid())
	require.True(t, cl1.Equal(ChannelLayoutStereo))
	require.False(t, cl1.Equal(ChannelLayoutMono))
	cl2 := ChannelLayout{}
	require.Equal(t, 0, cl2.Channels())
	require.False(t, cl2.Valid())
	require.Equal(t, "", cl2.String())
	require.False(t, cl1.Equal(cl2))
	cl3 := ChannelLayout{}
	require.True(t, cl2.Equal(cl3))
}
