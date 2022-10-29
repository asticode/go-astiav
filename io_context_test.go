package astiav_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestIOContext_Write(t *testing.T) {
	ioctxWrapper(t, func(c *astiav.IOContext, path string) {
		err := c.Open(path, astiav.NewIOContextFlags(
			astiav.IOContextFlagWrite))
		require.NoError(t, err)

		n, err := c.Write(nil)
		require.NoError(t, err)
		require.Equal(t, 0, n)

		l, err := c.Write([]byte("test"))
		require.Equal(t, 4, l)
		require.NoError(t, err)

		err = c.Closep()
		require.NoError(t, err)

		b, err := os.ReadFile(path)
		require.NoError(t, err)
		require.Equal(t, "test", string(b))
	})
}

func TestIOContext_WriteCopy(t *testing.T) {
	ioctxWrapper(t, func(c *astiav.IOContext, path string) {
		err := c.Open(path, astiav.NewIOContextFlags(
			astiav.IOContextFlagWrite))
		require.NoError(t, err)

		n, err := c.Write(nil)
		require.NoError(t, err)
		require.Equal(t, 0, n)

		l, err := io.Copy(c, strings.NewReader("test"))
		require.Equal(t, int64(4), l)
		require.NoError(t, err)

		err = c.Closep()
		require.NoError(t, err)

		b, err := os.ReadFile(path)
		require.NoError(t, err)
		require.Equal(t, "test", string(b))
	})
}

func TestIOContext_Read(t *testing.T) {
	ioctxWrapper(t, func(c *astiav.IOContext, path string) {
		err := os.WriteFile(path, []byte("testing"), 0644)
		require.NoError(t, err)

		err = c.Open(path, astiav.NewIOContextFlags(
			astiav.IOContextFlagRead))
		require.NoError(t, err)

		// Read from file (using astiav.IOContext.Read())
		d := make([]byte, 16)
		n, err := c.Read(d)
		require.NoError(t, err)
		require.Equal(t, 7, n)
		require.Equal(t, "testing", string(d[:n]))

		// Close context
		err = c.Closep()
		require.NoError(t, err)
	})
}

func TestIOContext_ReadCopy(t *testing.T) {
	ioctxWrapper(t, func(c *astiav.IOContext, path string) {
		err := os.WriteFile(path, []byte("testing"), 0644)
		require.NoError(t, err)

		err = c.Open(path, astiav.NewIOContextFlags(
			astiav.IOContextFlagRead))
		require.NoError(t, err)

		// Read from file (using astiav.IOContext.Read())
		buf := bytes.NewBufferString("")
		n, err := io.Copy(buf, c)
		require.Equal(t, int64(7), n)
		require.NoError(t, err)
		require.Equal(t, "testing", buf.String())

		err = c.Closep()
		require.NoError(t, err)
	})
}

func TestIOSeekerContext_Seek(t *testing.T) {
	ioctxWrapper(t, func(c *astiav.IOContext, path string) {
		err := os.WriteFile(path, []byte("testing"), 0644)
		require.NoError(t, err)

		err = c.Open(path, astiav.NewIOContextFlags(
			astiav.IOContextFlagRead))
		require.NoError(t, err)

		l, err := c.Seek(4, io.SeekStart)
		require.Equal(t, int64(4), l)
		require.NoError(t, err)

		d := make([]byte, 16)
		n, err := c.Read(d)
		require.Equal(t, 3, n)
		require.NoError(t, err)
		require.Equal(t, "ing", string(d[:n]))

		err = c.Closep()
		require.NoError(t, err)
	})
}

func ioctxWrapper(t *testing.T, cb func(c *astiav.IOContext, path string)) {
	c := astiav.NewIOContext()
	path := filepath.Join(t.TempDir(), "iocontext.txt")

	cb(c, path)

	err := os.Remove(path)
	require.NoError(t, err)
}
