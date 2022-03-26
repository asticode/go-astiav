package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDisplayMatrix(t *testing.T) {
	_, err := NewDisplayMatrixFromBytes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	require.Error(t, err)
	dm, err := NewDisplayMatrixFromBytes([]byte{0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64})
	require.NoError(t, err)
	require.Equal(t, DisplayMatrix{0x0, 0xffff0000, 0x0, 0x10000, 0x0, 0x0, 0x0, 0x0, 0x40000000}, *dm)
	require.Equal(t, float64(90), dm.Rotation())
	dm, err = NewDisplayMatrixFromBytes([]byte{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64})
	require.NoError(t, err)
	require.Equal(t, DisplayMatrix{0x0, 0x10000, 0x0, 0xffff0000, 0x0, 0x0, 0x0, 0x0, 0x40000000}, *dm)
	require.Equal(t, float64(-90), dm.Rotation())
}
