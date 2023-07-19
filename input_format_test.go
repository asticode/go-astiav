package astiav_test

import (
	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFindInputFormat(t *testing.T) {
	inputFormat := astiav.FindInputFormat("video4linux2")
	require.True(t, inputFormat.Flags().Has(astiav.IOFormatFlagNoByteSeek))
}
