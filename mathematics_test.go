package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMathematics(t *testing.T) {
	t.Run("RescaleQ", func(t *testing.T) {
		require.Equal(t, int64(1000), RescaleQ(100, NewRational(1, 100), NewRational(1, 1000)))
	})
	t.Run("RescaleQRnd", func(t *testing.T) {
		require.Equal(t, int64(0), RescaleQRnd(1, NewRational(1, 100), NewRational(1, 10), RoundingDown))
		require.Equal(t, int64(1), RescaleQRnd(1, NewRational(1, 100), NewRational(1, 10), RoundingUp))
	})
	t.Run("MulQ", func(t *testing.T) {
		require.Equal(t, NewRational(3, 5000), MulQ(NewRational(2, 100), NewRational(3, 100)))
	})
	t.Run("RescaleDelta", func(t *testing.T) {
		var lastDuration = NoPtsValue
		inputTS, inputTb := int64(9999), NewRational(1, 1000)
		duration, durationTb := int64(100), NewRational(1, 48000)
		outTb := NewRational(1, 2000)
		outTS, newLast := RescaleDelta(inputTb, inputTS, durationTb, duration, lastDuration, outTb)
		require.Equal(t, int64(19998), outTS)
		require.Equal(t, int64(480052), newLast)
	})
}
