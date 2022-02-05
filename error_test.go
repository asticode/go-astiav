package astiav_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

type testError struct{}

func (err testError) Error() string { return "" }

func TestError(t *testing.T) {
	require.Equal(t, "Decoder not found", astiav.ErrDecoderNotFound.Error())
	err1 := fmt.Errorf("test 1: %w", astiav.ErrDecoderNotFound)
	require.True(t, errors.Is(err1, astiav.ErrDecoderNotFound))
	require.False(t, errors.Is(err1, testError{}))
	err2 := fmt.Errorf("test 2: %w", astiav.ErrDemuxerNotFound)
	require.False(t, errors.Is(err2, astiav.ErrDecoderNotFound))
}
