package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestSampleFormat(t *testing.T) {
	require.Equal(t, "s16", astiav.SampleFormatS16.String())
}
