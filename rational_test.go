package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestRational(t *testing.T) {
	r := astiav.NewRational(2, 1)
	require.Equal(t, 2, r.Num())
	require.Equal(t, 1, r.Den())
	r.SetNum(1)
	r.SetDen(2)
	require.Equal(t, 1, r.Num())
	require.Equal(t, 2, r.Den())
	require.Equal(t, "1/2", r.String())
	require.Equal(t, 0.5, r.ToDouble())
	r.SetDen(0)
	require.Equal(t, float64(0), r.ToDouble())
	require.Equal(t, "0", r.String())
}
