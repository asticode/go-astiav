package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHardwareDeviceType(t *testing.T) {
	require.Equal(t, "cuda", HardwareDeviceTypeCUDA.String())
	require.Equal(t, FindHardwareDeviceTypeByName("cuda"), HardwareDeviceTypeCUDA)
}
