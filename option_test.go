package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOption(t *testing.T) {
	fc, err := AllocOutputFormatContext(nil, "mp4", "")
	require.NoError(t, err)
	pd := fc.PrivateData()
	require.NotNil(t, pd)
	os := pd.Options()
	require.NotNil(t, os)
	l := os.List()
	require.Len(t, l, 56)
	const name = "brand"
	o := l[0]
	require.Equal(t, name, o.Name())
	_, err = os.Get("invalid", NewOptionSearchFlags())
	require.Error(t, err)
	v, err := os.Get(name, NewOptionSearchFlags())
	require.NoError(t, err)
	require.Equal(t, "", v)
	require.Error(t, os.Set("invalid", "", NewOptionSearchFlags()))
	const value = "test"
	require.NoError(t, os.Set(name, value, NewOptionSearchFlags()))
	v, err = os.Get(name, NewOptionSearchFlags())
	require.NoError(t, err)
	require.Equal(t, value, v)
}
