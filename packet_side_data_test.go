package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPacketSideData(t *testing.T) {
	cp := AllocCodecParameters()
	defer cp.Free()
	b := []byte("test")
	sd := cp.SideData()
	require.NoError(t, sd.Add(PacketSideDataTypeDisplaymatrix, b))
	require.Equal(t, b, sd.Get(PacketSideDataTypeDisplaymatrix))
}
