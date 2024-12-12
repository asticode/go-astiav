package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIOInterrupter(t *testing.T) {
	ii := NewIOInterrupter()
	require.False(t, ii.Interrupted())
	ii.Interrupt()
	require.True(t, ii.Interrupted())
	ii.Resume()
	require.False(t, ii.Interrupted())
}
