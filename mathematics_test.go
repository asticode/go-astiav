package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMathematics(t *testing.T) {
	require.Equal(t, int64(1000), RescaleQ(100, NewRational(1, 100), NewRational(1, 1000)))
	require.Equal(t, int64(0), RescaleQRnd(1, NewRational(1, 100), NewRational(1, 10), RoundingDown))
	require.Equal(t, int64(1), RescaleQRnd(1, NewRational(1, 100), NewRational(1, 10), RoundingUp))
}
