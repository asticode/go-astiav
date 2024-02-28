package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	f := FindFilterByName("format")
	require.NotNil(t, f)
	require.Equal(t, "format", f.Name())
	require.Equal(t, "format", f.String())
}
