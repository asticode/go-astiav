package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOutputFormat(t *testing.T) {
	formatName := "rawvideo"
	outputFormat := FindOutputFormat(formatName)
	require.NotNil(t, outputFormat)
	require.Equal(t, formatName, outputFormat.Name())
	require.Equal(t, formatName, outputFormat.String())
	require.Equal(t, "raw video", outputFormat.LongName())
}
