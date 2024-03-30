package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPacketSideDataType(t *testing.T) {
	require.Equal(t, "Display Matrix", PacketSideDataTypeDisplaymatrix.Name())
	require.Equal(t, "Display Matrix", PacketSideDataTypeDisplaymatrix.String())
}
