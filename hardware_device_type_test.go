package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestHardwareDeviceType(t *testing.T) {
	require.Equal(t, "cuda", astiav.HardwareDeviceTypeCUDA.String())
	require.Equal(t, astiav.FindHardwareDeviceTypeByName("cuda"), astiav.HardwareDeviceTypeCUDA)
}
