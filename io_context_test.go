package astiav

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIOContext(t *testing.T) {
	path := filepath.Join(t.TempDir(), "iocontext.txt")
	c, err := OpenIOContext(path, NewIOContextFlags(IOContextFlagWrite))
	require.NoError(t, err)
	cl := c.Class()
	require.NotNil(t, cl)
	require.Equal(t, "AVIOContext", cl.Name())
	c.Write(nil)
	c.Write([]byte("test"))
	err = c.Closep()
	require.NoError(t, err)
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "test", string(b))
	err = os.Remove(path)
	require.NoError(t, err)
}
