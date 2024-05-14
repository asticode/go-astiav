package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitStreamFilter(t *testing.T) {
	fn := "null"
	f := FindBitStreamFilterByName(fn)
	require.NotNil(t, f)
	require.Equal(t, f.Name(), fn)
	require.Equal(t, f.String(), fn)

	f = FindBitStreamFilterByName("foobar_non_existing_bsf")
	require.Nil(t, f)
}
