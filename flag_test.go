package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlag(t *testing.T) {
	f := flags(2 | 4)
	r := f.add(1)
	require.Equal(t, 7, r)
	r = f.del(2)
	require.Equal(t, 4, r)
	require.False(t, f.has(1))
	require.True(t, f.has(4))
}
