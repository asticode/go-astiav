package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInputFormat(t *testing.T) {
	formatName := "rawvideo"
	inputFormat := FindInputFormat(formatName)
	require.NotNil(t, inputFormat)
	require.Equal(t, formatName, inputFormat.Name())
	require.Equal(t, formatName, inputFormat.String())
	require.Equal(t, "raw video", inputFormat.LongName())
}
