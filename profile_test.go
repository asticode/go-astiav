package astiav

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProfileName(t *testing.T) {
	got := ProfileName(CodecIDH264, ProfileH264High)
	require.Equal(t, "High", got)
}
