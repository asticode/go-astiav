package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMathematics(t *testing.T) {
	require.Equal(t, int64(1000), RescaleQ(100, NewRational(1, 100), NewRational(1, 1000)))
	require.Equal(t, int64(0), RescaleQRnd(1, NewRational(1, 100), NewRational(1, 10), RoundingDown))
	require.Equal(t, int64(1), RescaleQRnd(1, NewRational(1, 100), NewRational(1, 10), RoundingUp))
	require.Equal(t, 0.04, Q2D(NewRational(1, 25)))
	require.Equal(t, 0.3, Q2D(NewRational(3, 10)))
	require.Equal(t, float64(30), Q2D(NewRational(30, 1)))
}
