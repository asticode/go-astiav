package astiav

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRational(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		r := NewRational(2, 1)
		require.Equal(t, 2, r.Num())
		require.Equal(t, 1, r.Den())
		require.Equal(t, 0.5, r.Invert().Float64())
	})
	t.Run("Set", func(t *testing.T) {
		r := NewRational(0, 0)
		r.SetNum(1)
		r.SetDen(2)
		require.Equal(t, 1, r.Num())
		require.Equal(t, 2, r.Den())
	})
	t.Run("Float", func(t *testing.T) {
		r := NewRational(0, 0)
		require.Equal(t, 0.0, r.Float64())
		r = NewRational(2, 1)
		require.Equal(t, 2.0, r.Float64())
		r = NewRational(1, 2)
		require.Equal(t, 0.5, r.Float64())
	})
	t.Run("String", func(t *testing.T) {
		r := NewRational(2, 1)
		require.Equal(t, "2/1", r.String())
	})
	t.Run("TextMarshal", func(t *testing.T) {
		r := NewRational(2, 1)
		s, err := r.MarshalText()
		require.NoError(t, err)
		var r2 Rational
		require.NoError(t, r2.UnmarshalText(s))
		require.Equal(t, r.Num(), r2.Num())
		require.Equal(t, r.Den(), r2.Den())
	})
	t.Run("json.Marshal", func(t *testing.T) {
		type test struct {
			Timebase Rational
		}
		x1 := test{Timebase: NewRational(2, 1)}
		data, err := json.Marshal(x1)
		require.NoError(t, err)
		x2 := test{}
		require.NoError(t, json.Unmarshal(data, &x2))
		require.Equal(t, x1.Timebase.Num(), x2.Timebase.Num())
		require.Equal(t, x1.Timebase.Den(), x2.Timebase.Den())
	})
}
