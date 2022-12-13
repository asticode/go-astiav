package astiav_test

import (
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestIOContext_Open_ReadWriteSeek(t *testing.T) {
	c := astiav.NewIOContext()
	path := filepath.Join(t.TempDir(), "iocontext.txt")

	// Write Test
	err := c.Open(path, astiav.NewIOContextFlags(astiav.IOContextFlagWrite))
	require.NoError(t, err)

	err = c.Write(nil)
	require.NoError(t, err)
	require.Equal(t, int64(0), c.Size())

	err = c.Write([]byte("testtest"))
	c.Flush()
	require.Equal(t, int64(8), c.Size())
	require.NoError(t, err)

	err = c.Closep()
	require.NoError(t, err)

	// Read Test
	c = astiav.NewIOContext()
	err = c.Open(path, astiav.NewIOContextFlags(astiav.IOContextFlagRead))
	require.NoError(t, err)

	d := make([]byte, 32768)
	j, err := c.Read(d)
	require.NoError(t, err)
	require.Equal(t, 8, j)
	require.Equal(t, "testtest", string(d[:j]))

	// Cleanup
	err = c.Closep()
	require.NoError(t, err)

	b, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "testtest", string(b))

	err = os.Remove(path)
	require.NoError(t, err)
}

func TestIOContext_OpenWith_Write(t *testing.T) {
	c := astiav.NewIOContext()
	path := filepath.Join(t.TempDir(), "iocontext.txt")

	// Write Test
	dict := astiav.NewDictionary()
	defer dict.Free()
	dict.Set("test", "test", 0)
	err := c.OpenWith(path, astiav.NewIOContextFlags(
		astiav.IOContextFlagReadWrite), dict)
	require.NoError(t, err)

	err = c.Write(nil)
	require.NoError(t, err)
	require.Equal(t, int64(0), c.Size())

	err = c.Write([]byte("testtest"))
	c.Flush()
	require.Equal(t, int64(8), c.Size())
	require.NoError(t, err)

	// Cleanup
	err = c.Closep()
	require.NoError(t, err)

	b, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "testtest", string(b))

	err = os.Remove(path)
	require.NoError(t, err)
}

func TestIOContext_BufferReader(t *testing.T) {
	buffer := randomBytes(1024 * 1024)
	c := astiav.AllocIOContextBufferReader(buffer)
	defer c.Free()

	buf := make([]byte, 256)
	n, err := c.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 256, n)

	// Error expected because write is not supported
	err = c.Write(buf)
	require.Error(t, err)
	require.True(t, astiav.ErrEio.Is(err))
}

func TestIOContext_ReadSeeker(t *testing.T) {
	f := createTestFile(t, string(randomBytes(256)))
	c := astiav.AllocIOContextReadSeeker(f)
	defer c.Free()

	buf1 := make([]byte, 256)
	n, err := c.Read(buf1)
	require.NoError(t, err)
	require.Equal(t, 256, n)

	require.True(t, c.Seekable())
	i, err := c.Seek(0, io.SeekStart)
	require.Equal(t, int64(0), i)

	buf2 := make([]byte, 256)
	n, err = c.Read(buf2)
	require.NoError(t, err)
	require.Equal(t, 256, n)
}

func TestIOContext_BufferReadSeeker(t *testing.T) {
	buffer := randomBytes(1024 * 1024)
	c := astiav.AllocIOContextBufferReader(buffer)
	defer c.Free()

	buf1 := make([]byte, 256)
	n, err := c.Read(buf1)
	require.NoError(t, err)
	require.Equal(t, 256, n)

	require.True(t, c.Seekable())
	i, err := c.Seek(0, io.SeekStart)
	require.Equal(t, int64(0), i)

	buf2 := make([]byte, 256)
	n, err = c.Read(buf2)
	require.NoError(t, err)
	require.Equal(t, 256, n)
}

func TestIOContext_WriteSeeker(t *testing.T) {
	randBytes := randomBytes(256)
	f := createTestFile(t, string(randBytes))
	c := astiav.AllocIOContextWriteSeeker(f)
	defer c.Free()

	err := c.Write(randBytes)
	require.NoError(t, err)

	require.True(t, c.Seekable())
	i, err := c.Seek(0, io.SeekStart)
	require.Equal(t, int64(0), i)

	readBytes := make([]byte, 256)
	n, err := f.Read(readBytes)
	require.NoError(t, err)
	require.Equal(t, 256, n)
	require.Equal(t, readBytes, randBytes)

	n, err = c.Read(make([]byte, 256))
	require.Error(t, err)
	require.True(t, astiav.ErrEio.Is(err))
	require.Equal(t, int(astiav.ErrEio), n)
}

func TestIOContext_BufferWriteSeeker(t *testing.T) {
	rbuffer := randomBytes(1024)
	buffer := make([]byte, 1024)

	c := astiav.AllocIOContextBufferWriter(buffer)
	defer c.Free()

	err := c.Write(rbuffer)
	c.Flush()
	require.NoError(t, err)
	require.Equal(t, rbuffer, buffer)

	require.True(t, c.Seekable())
	i, err := c.Seek(0, io.SeekStart)
	require.Equal(t, int64(0), i)

	// Error expected because read is not supported
	n, err := c.Read(make([]byte, 256))
	require.True(t, astiav.ErrEio.Is(err))
	require.Equal(t, int(astiav.ErrEio), n)
}

func TestIOContext_CallbacksWriteRead(t *testing.T) {
	byteArr := make([]byte, 64)
	size := 0
	pos := 0
	c := astiav.AllocIOContextCallback(
		func(buf []byte) int {
			min := len(buf)
			if pos >= size {
				return int(astiav.ErrEof)
			}
			if size < min {
				min = size
			}
			for i := 0; i < min; i++ {
				buf[i] = byteArr[pos+i]
			}
			pos += min
			return min
		}, func(buf []byte) int {
			bufSize := len(buf)

			if pos >= len(byteArr) {
				return 0
			}
			if (pos + bufSize) > len(byteArr) {
				bufSize = (pos + bufSize) - len(byteArr)
			}
			for i := 0; i < bufSize; i++ {
				byteArr[pos+i] = buf[i]
			}
			pos += bufSize
			size += bufSize
			return bufSize
		}, func(offset int64, whence int) int64 {
			pos = int(offset)
			return offset
		})
	defer c.Free()

	original := randomBytes(128)
	err := c.Write(original)
	require.NoError(t, err)

	require.True(t, c.Seekable())
	i, err := c.Seek(0, io.SeekStart)
	require.Equal(t, int64(0), i)

	buf := make([]byte, 64)
	n, err := c.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 64, n)

	buf = make([]byte, 64)
	n, err = c.Read(buf)
	require.Equal(t, astiav.ErrEof, (astiav.Error)(n))
}

func randomBytes(size int) []byte {
	buf := make([]byte, size)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

func BenchmarkIOContext_OpenAndParseAudio(b *testing.B) {
	astiav.SetLogLevel(astiav.LogLevelError)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		openFromReader(b, "testdata/audio.mp3")
	}
}

func BenchmarkIOContext_OpenAndParseVideo(b *testing.B) {
	astiav.SetLogLevel(astiav.LogLevelError)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		openFromReader(b, "testdata/video.mp4")
	}
}

func BenchmarkIOContext_OpenAndParseImage(b *testing.B) {
	astiav.SetLogLevel(astiav.LogLevelError)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		openFromReader(b, "testdata/image.jpeg")
	}
}

func openFromReader(b *testing.B, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()

	b.StartTimer()
	defer b.StopTimer()
	fc := astiav.AllocFormatContext()
	defer fc.Free()
	ioCtx := astiav.AllocIOContextReadSeeker(file)
	if ioCtx == nil {
		b.Fatal("ioCtx is nil")
	}
	defer ioCtx.Free()
	fc.SetPb(ioCtx)
	dict1 := astiav.NewDictionary()
	err = fc.OpenInput("testing", nil, dict1)
	if err != nil {
		b.Fatalf("error: %v %d\n", err, err)
	}
	fc.SetFlags(fc.Flags().Add(astiav.FormatContextFlagCustomIo))
	dict2 := astiav.NewDictionary()
	if dict2 == nil {
		b.Fatal("dict is nil")
	}
	defer dict2.Free()
	err = fc.FindStreamInfo(dict2)
	if err != nil {
		b.Fatal(err)
	}
	for _, is := range fc.Streams() {
		if is.CodecParameters().MediaType() != astiav.MediaTypeAudio &&
			is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
			continue
		}
	}
}

func createTestFile(t *testing.T, data string) *os.File {
	dir := path.Join(t.TempDir(), t.Name())
	f, err := os.Create(dir)
	require.NoError(t, err)
	n, err := f.WriteString(data)
	require.NoError(t, err)
	require.Equal(t, len(data), n)
	require.NoError(t, f.Close())

	f, err = os.OpenFile(dir, os.O_RDWR, 0)
	require.NoError(t, err)
	return f
}
