package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPictureType(t *testing.T) {
	require.Equal(t, "I", PictureTypeI.String())
}
