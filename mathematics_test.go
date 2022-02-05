package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestMathematics(t *testing.T) {
	require.Equal(t, int64(1000), astiav.RescaleQ(100, astiav.NewRational(1, 100), astiav.NewRational(1, 1000)))
	require.Equal(t, int64(0), astiav.RescaleQRnd(1, astiav.NewRational(1, 100), astiav.NewRational(1, 10), astiav.RoundingDown))
	require.Equal(t, int64(1), astiav.RescaleQRnd(1, astiav.NewRational(1, 100), astiav.NewRational(1, 10), astiav.RoundingUp))
}
