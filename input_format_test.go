package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestInputFormat(t *testing.T) {
	formatName := "rawvideo"
	inputFormat := astiav.FindInputFormat(formatName)
	require.NotNil(t, inputFormat)
	require.True(t, inputFormat.Name() == formatName)
}
