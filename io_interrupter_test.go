package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultIOInterrupter(t *testing.T) {
	ii := newDefaultIOInterrupter()
	require.Equal(t, 0, int(ii.i))
	ii.Interrupt()
	require.Equal(t, 1, int(ii.i))
	ii.Resume()
	require.Equal(t, 0, int(ii.i))
}
