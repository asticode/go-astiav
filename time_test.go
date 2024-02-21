package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	require.NotEqual(t, 0, astiav.RelativeTime())
}
