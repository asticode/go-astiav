package astiav

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIOContext(t *testing.T) {
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
}

func TestOpenIOContext(t *testing.T) {
	path := filepath.Join(t.TempDir(), "iocontext.txt")
	c, err := OpenIOContext(path, NewIOContextFlags(IOContextFlagWrite))
	require.NoError(t, err)
	cl := c.Class()
	require.NotNil(t, cl)
	require.Equal(t, "AVIOContext", cl.Name())
	c.Write(nil)
	c.Write([]byte("test"))
	require.NoError(t, c.Close())
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "test", string(b))
	err = os.Remove(path)
	require.NoError(t, err)
}
