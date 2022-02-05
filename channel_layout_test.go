package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestChannelLayout(t *testing.T) {
	require.Equal(t, 2, astiav.ChannelLayoutStereo.NbChannels())
	require.Equal(t, "stereo", astiav.ChannelLayoutStereo.String())
	require.Equal(t, "1 channels (FL+FR)", astiav.ChannelLayoutStereo.StringWithNbChannels(1))
}
