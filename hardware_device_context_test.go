package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestHardwareDeviceContext(t *testing.T) {
	hcd, err := astiav.CreateHardwareDeviceContext(astiav.HardwareDeviceTypeCUDA, "", astiav.NewDictionary())
	require.NoError(t, err)
	defer hcd.Free()
}
