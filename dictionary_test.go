package astiav

import (
	"maps"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDictionary(t *testing.T) {
	d1 := NewDictionary()
	defer d1.Free()
	err := d1.ParseString("invalid,,", ":", ",", 0)
	require.Error(t, err)
	err = d1.ParseString("k1=v1,k2=v2", "=", ",", 0)
	require.NoError(t, err)
	err = d1.Set("k3", "v3", 0)
	require.NoError(t, err)
	e := d1.Get("k1", nil, 0)
	require.NotNil(t, e)
	require.Equal(t, "k1", e.Key())
	require.Equal(t, "v1", e.Value())
	e = d1.Get("k2", nil, 0)
	require.NotNil(t, e)
	require.Equal(t, "k2", e.Key())
	require.Equal(t, "v2", e.Value())
	e = d1.Get("k3", nil, 0)
	require.NotNil(t, e)
	require.Equal(t, "k3", e.Key())
	require.Equal(t, "v3", e.Value())
	e = d1.Get("k4", nil, 0)
	require.Nil(t, e)

	b := d1.Pack()
	require.Equal(t, "k1\x00v1\x00k2\x00v2\x00k3\x00v3\x00", string(b))

	err = d1.Unpack([]byte("k4\x00v4\x00"))
	require.NoError(t, err)
	e = d1.Get("k4", nil, 0)
	require.NotNil(t, e)
	require.Equal(t, "k4", e.Key())
	require.Equal(t, "v4", e.Value())

	d2 := NewDictionary()
	defer d2.Free()
	require.NoError(t, d1.Copy(d2, NewDictionaryFlags()))
	e = d2.Get("k4", nil, 0)
	require.NotNil(t, e)
	require.Equal(t, "k4", e.Key())
	require.Equal(t, "v4", e.Value())

	d3 := NewDictionary()
	defer d3.Free()
	err = d3.Set("k1", "v1", 0)
	require.NoError(t, err)
	err = d3.Set("k2", "v2", 0)
	require.NoError(t, err)
	err = d3.Set("k3", "v3", 0)
	require.NoError(t, err)

	m1 := maps.Collect(d3.Seq())
	m2 := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}
	require.Equal(t, m1, m2)
}
