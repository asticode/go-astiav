package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChannelLayout(t *testing.T) {
	t.Run("Actual", func(t *testing.T) {
		cl1 := ChannelLayoutStereo
		require.Equal(t, 2, cl1.Channels())
		require.Equal(t, "stereo", cl1.String())
		require.True(t, cl1.Valid())
		require.True(t, cl1.Equal(ChannelLayoutStereo))
		require.False(t, cl1.Equal(ChannelLayoutMono))
	})
	t.Run("Empty", func(t *testing.T) {
		cl2 := ChannelLayout{}
		require.Equal(t, 0, cl2.Channels())
		require.False(t, cl2.Valid())
		require.Equal(t, "", cl2.String())
		require.False(t, cl2.Equal(ChannelLayoutStereo))
		cl3 := ChannelLayout{}
		require.True(t, cl2.Equal(cl3))
	})
	t.Run("Parse", func(t *testing.T) {
		cl1 := ChannelLayoutStereo
		ls := cl1.String()
		cl2, err := ParseChannelLayout(ls)
		require.NoError(t, err)
		require.True(t, cl1.Equal(cl2))
	})
	t.Run("TextMarshal", func(t *testing.T) {
		x1 := ChannelLayoutStereo
		s, err := x1.MarshalText()
		require.NoError(t, err)
		var x2 ChannelLayout
		require.NoError(t, x2.UnmarshalText(s))
		require.True(t, x1.Equal(x2))
	})
}
