package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	require.NotEqual(t, 0, RelativeTime())
}

func TestCompareTimestamps(t *testing.T) {
	a := int64(0)
	timeBaseA := NewRational(1, 1)
	b := int64(2)
	timeBaseB := NewRational(1, 2)
	require.Equal(t, CompareTimestampsResultABeforeB, CompareTimestamps(a, b, timeBaseA, timeBaseB))
	a = 1
	require.Equal(t, CompareTimestampsResultAEqualB, CompareTimestamps(a, b, timeBaseA, timeBaseB))
	a = 2
	require.Equal(t, CompareTimestampsResultAAfterB, CompareTimestamps(a, b, timeBaseA, timeBaseB))
}
