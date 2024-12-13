package astiav

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIOContext(t *testing.T) {
	t.Run("read write seek", func(t *testing.T) {
		var seeked bool
		rb := []byte("read")
		wb := []byte("write")
		var written []byte
		c, err := AllocIOContext(8, true, func(b []byte) (int, error) {
			copy(b, rb)
			return len(rb), nil
		}, func(offset int64, whence int) (n int64, err error) {
			seeked = true
			return offset, nil
		}, func(b []byte) (int, error) {
			written = make([]byte, len(b))
			copy(written, b)
			return len(b), nil
		})
		require.NoError(t, err)
		defer c.Free()
		b := make([]byte, 6)
		n, err := c.Read(b)
		require.NoError(t, err)
		require.Equal(t, 4, n)
		require.Equal(t, rb, b[:n])
		_, err = c.Seek(2, 0)
		require.NoError(t, err)
		require.True(t, seeked)
		c.Write(wb)
		c.Flush()
		require.Equal(t, wb, written)
	})

	t.Run("io.EOF is mapped to AVERROR_EOF when reading", func(t *testing.T) {
		c, err := AllocIOContext(8, false, func(b []byte) (int, error) {
			return 0, io.EOF
		}, nil, nil)
		require.NoError(t, err)
		defer c.Free()
		b := make([]byte, 100)
		n, err := c.Read(b)
		require.ErrorIs(t, err, ErrEof)
		require.Equal(t, 0, n)
	})
}

func TestOpenIOContext(t *testing.T) {
	path := filepath.Join(t.TempDir(), "iocontext.txt")
	c1, err := OpenIOContext(path, NewIOContextFlags(IOContextFlagWrite), nil, nil)
	require.NoError(t, err)
	defer os.RemoveAll(path)
	cl := c1.Class()
	require.NotNil(t, cl)
	require.Equal(t, "AVIOContext", cl.Name())
	c1.Write(nil)
	c1.Write([]byte("test"))
	require.NoError(t, c1.Close())
	b1, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "test", string(b1))

	ii := NewIOInterrupter()
	defer ii.Free()
	c2, err := OpenIOContext(path, NewIOContextFlags(IOContextFlagRead), ii, nil)
	require.NoError(t, err)
	b2 := make([]byte, 10)
	_, err = c2.Read(b2)
	require.NoError(t, err)
	ii.Interrupt()
	_, err = c2.Read(b2)
	require.ErrorIs(t, err, ErrExit)
	require.ErrorIs(t, c2.Close(), ErrExit)

	d := NewDictionary()
	require.NoError(t, d.Set("protocol_whitelist", "rtp", NewDictionaryFlags()))
	_, err = OpenIOContext(path, NewIOContextFlags(IOContextFlagWrite), nil, d)
	require.Error(t, err)
}
