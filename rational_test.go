package astiav

import (
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
		r = NewRational(1, 0)
		require.Equal(t, 0.0, r.Float64())
		r = NewRational(2, 1)
		require.Equal(t, 2.0, r.Float64())
		r = NewRational(1, 2)
		require.Equal(t, 0.5, r.Float64())
	})
	t.Run("String", func(t *testing.T) {
		r := NewRational(1, 2)
		require.Equal(t, "1/2", r.String())
		r = NewRational(2, 1)
		require.Equal(t, "2", r.String())
	})
	t.Run("NewRationalFromString", func(t *testing.T) {
		t.Run("Fraction", func(t *testing.T) {
			got, err := NewRationalFromString("1/2")
			require.NoError(t, err)
			want := NewRational(1, 2)
			require.Equal(t, want.Num(), got.Num())
			require.Equal(t, want.Den(), got.Den())
		})
		t.Run("Whole", func(t *testing.T) {
			got, err := NewRationalFromString("2")
			require.NoError(t, err)
			want := NewRational(2, 1)
			require.Equal(t, want.Num(), got.Num())
			require.Equal(t, want.Den(), got.Den())
		})
		t.Run("Invalid/Denominator", func(t *testing.T) {
			_, err := NewRationalFromString("1/abc")
			require.EqualError(t, err, "astiav: parsing denominator failed: strconv.ParseInt: parsing \"abc\": invalid syntax")
		})
		t.Run("Invalid/Numerator", func(t *testing.T) {
			_, err := NewRationalFromString("abc")
			require.EqualError(t, err, "astiav: parsing numerator failed: strconv.ParseInt: parsing \"abc\": invalid syntax")
		})
	})
	t.Run("TextMarshal", func(t *testing.T) {
		r := NewRational(1, 2)
		b, err := r.MarshalText()
		require.Equal(t, "1/2", string(b))
		require.NoError(t, err)
		var r2 Rational
		require.NoError(t, r2.UnmarshalText(b))
		require.Equal(t, r.Num(), r2.Num())
		require.Equal(t, r.Den(), r2.Den())
	})
}
