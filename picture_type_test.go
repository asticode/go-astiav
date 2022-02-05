package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestPictureType(t *testing.T) {
	require.Equal(t, "I", astiav.PictureTypeI.String())
}
