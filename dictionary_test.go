package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDictionary(t *testing.T) {
	t.Run("ParseString", func(t *testing.T) {
		d := NewDictionary()
		defer d.Free()
		err := d.ParseString("invalid,,", ":", ",", 0)
		require.Error(t, err)
		err = d.ParseString("k1=v1,k2=v2", "=", ",", 0)
		require.NoError(t, err)
		e := d.Get("k1", nil, 0)
		require.NotNil(t, e)
		require.Equal(t, "k1", e.Key())
		require.Equal(t, "v1", e.Value())
		e = d.Get("k2", nil, 0)
		require.NotNil(t, e)
		require.Equal(t, "k2", e.Key())
		require.Equal(t, "v2", e.Value())
	})
	t.Run("Set", func(t *testing.T) {
		d := NewDictionary()
		defer d.Free()
		err := d.ParseString("k1=v1,k2=v2", "=", ",", 0)
		require.NoError(t, err)
		err = d.Set("k3", "v3", 0)
		require.NoError(t, err)
		e := d.Get("k1", nil, 0)
		require.NotNil(t, e)
		require.Equal(t, "k1", e.Key())
		require.Equal(t, "v1", e.Value())
		e = d.Get("k2", nil, 0)
		require.NotNil(t, e)
		require.Equal(t, "k2", e.Key())
		require.Equal(t, "v2", e.Value())
		e = d.Get("k3", nil, 0)
		require.NotNil(t, e)
		require.Equal(t, "k3", e.Key())
		require.Equal(t, "v3", e.Value())
		e = d.Get("k4", nil, 0)
		require.Nil(t, e)
	})
	t.Run("Copy", func(t *testing.T) {
		d1 := NewDictionary()
		defer d1.Free()
		err := d1.ParseString("k1=v1,k2=v2", "=", ",", 0)
		require.NoError(t, err)

		d2 := NewDictionary()
		defer d2.Free()
		require.NoError(t, d1.Copy(d2, NewDictionaryFlags()))
		e := d2.Get("k2", nil, 0)
		require.NotNil(t, e)
		require.Equal(t, "k2", e.Key())
		require.Equal(t, "v2", e.Value())
	})
	t.Run("Pack", func(t *testing.T) {
		d := NewDictionary()
		defer d.Free()
		err := d.ParseString("k1=v1,k2=v2", "=", ",", 0)
		require.NoError(t, err)
		require.Equal(t, "k1\x00v1\x00k2\x00v2\x00", string(d.Pack()))
	})
	t.Run("Unpack", func(t *testing.T) {
		d := NewDictionary()
		defer d.Free()
		err := d.Unpack([]byte("k4\x00v4\x00"))
		require.NoError(t, err)
		e := d.Get("k4", nil, 0)
		require.NotNil(t, e)
		require.Equal(t, "k4", e.Key())
		require.Equal(t, "v4", e.Value())
	})
	t.Run("ToMap", func(t *testing.T) {
		d := NewDictionary()
		defer d.Free()
		err := d.ParseString("k1=v1,k2=v2", "=", ",", 0)
		require.NoError(t, err)
		want := map[string]string{
			"k1": "v1",
			"k2": "v2",
		}
		require.Equal(t, want, d.ToMap())
	})
	t.Run("String", func(t *testing.T) {
		d := NewDictionary()
		defer d.Free()
		err := d.ParseString("k1=v1,k2=v2", "=", ",", 0)
		require.NoError(t, err)
		require.Equal(t, "[ k1:v1 k2:v2 ]", d.String())
	})
}
