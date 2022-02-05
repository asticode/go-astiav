package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	f := astiav.FindFilterByName("format")
	require.NotNil(t, f)
	require.Equal(t, "format", f.Name())
	require.Equal(t, "format", f.String())
}
