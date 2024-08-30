package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRational(t *testing.T) {
	r := NewRational(2, 1)
	require.Equal(t, 2, r.Num())
	require.Equal(t, 1, r.Den())
	require.Equal(t, 0.5, r.Invert().Float64())
	r.SetNum(1)
	r.SetDen(2)
	require.Equal(t, 1, r.Num())
	require.Equal(t, 2, r.Den())
	require.Equal(t, "1/2", r.String())
	require.Equal(t, 0.5, r.Float64())
	r.SetDen(0)
	require.Equal(t, float64(0), r.Float64())
	require.Equal(t, "0", r.String())
}
