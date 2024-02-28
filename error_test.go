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
	require.Equal(t, "Decoder not found", ErrDecoderNotFound.Error())
	err1 := fmt.Errorf("test 1: %w", ErrDecoderNotFound)
	require.True(t, errors.Is(err1, ErrDecoderNotFound))
	require.False(t, errors.Is(err1, testError{}))
	err2 := fmt.Errorf("test 2: %w", ErrDemuxerNotFound)
	require.False(t, errors.Is(err2, ErrDecoderNotFound))
}
