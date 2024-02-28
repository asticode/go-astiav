package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	require.NotEqual(t, 0, RelativeTime())
}
