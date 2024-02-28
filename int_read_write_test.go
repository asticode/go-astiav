package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntReadWrite(t *testing.T) {
	is := []uint8{1, 2, 3, 4, 5, 6, 7, 8}
	require.Equal(t, uint32(0), RL32([]byte{}))
	require.Equal(t, uint32(0x4030201), RL32(is))
	require.Equal(t, uint32(0), RL32WithOffset([]byte{}, 4))
	require.Equal(t, uint32(0x8070605), RL32WithOffset(is, 4))
}
