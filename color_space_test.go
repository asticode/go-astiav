package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColorSpace(t *testing.T) {
	require.Equal(t, "bt709", ColorSpaceBt709.String())
}
