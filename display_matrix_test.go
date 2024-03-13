package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDisplayMatrix(t *testing.T) {
	_, err := NewDisplayMatrixFromBytes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	require.Error(t, err)
	b := []byte{0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64}
	dm, err := NewDisplayMatrixFromBytes(b)
	require.NoError(t, err)
	require.Equal(t, DisplayMatrix{0x0, 0xffff0000, 0x0, 0x10000, 0x0, 0x0, 0x0, 0x0, 0x40000000}, *dm)
	require.Equal(t, -90.0, dm.Rotation())
	require.Equal(t, b, dm.Bytes())
	dm = NewDisplayMatrixFromRotation(-90)
	require.Equal(t, -90.0, dm.Rotation())
	dm, err = NewDisplayMatrixFromBytes([]byte{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64})
	require.NoError(t, err)
	require.Equal(t, DisplayMatrix{0x0, 0x10000, 0x0, 0xffff0000, 0x0, 0x0, 0x0, 0x0, 0x40000000}, *dm)
	require.Equal(t, 90.0, dm.Rotation())
}
