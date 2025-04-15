package astiav

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testError struct{}

func (err testError) Error() string { return "" }

func TestError(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		require.Equal(t, "Decoder not found", ErrDecoderNotFound.Error())
	})
	t.Run("Is/wrapped/match", func(t *testing.T) {
		err1 := fmt.Errorf("test 1: %w", ErrDecoderNotFound)
		require.True(t, errors.Is(err1, ErrDecoderNotFound))
	})
	t.Run("Is/wrap/wrong-type", func(t *testing.T) {
		err1 := fmt.Errorf("test 1: %w", ErrDecoderNotFound)
		require.False(t, errors.Is(err1, testError{}))
	})
	t.Run("Is/wrap/wrong-code", func(t *testing.T) {
		err2 := fmt.Errorf("test 2: %w", ErrDemuxerNotFound)
		require.False(t, errors.Is(err2, ErrDecoderNotFound))
	})
}

func TestLoggedError(t *testing.T) {
	t.Run("Is/wrapped/match", func(t *testing.T) {
		err1 := fmt.Errorf("test 1: %w", ErrDecoderNotFound)
		require.True(t, errors.Is(err1, ErrDecoderNotFound))
	})
	t.Run("Is/wrap/wrong-type", func(t *testing.T) {
		err1 := fmt.Errorf("test 1: %w", ErrDecoderNotFound)
		require.False(t, errors.Is(err1, testError{}))
	})
	t.Run("Is/wrap/wrong-code", func(t *testing.T) {
		err2 := fmt.Errorf("test 2: %w", ErrDemuxerNotFound)
		require.False(t, errors.Is(err2, ErrDecoderNotFound))
	})
}
