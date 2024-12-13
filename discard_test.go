package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscard(t *testing.T) {
	require.Equal(t, int(DiscardNone), -16)
	require.Equal(t, int(DiscardDefault), 0)
	require.Equal(t, int(DiscardNonRef), 8)
	require.Equal(t, int(DiscardBidirectional), 16)
	require.Equal(t, int(DiscardNonIntra), 24)
	require.Equal(t, int(DiscardNonKey), 32)
	require.Equal(t, int(DiscardAll), 48)
}
