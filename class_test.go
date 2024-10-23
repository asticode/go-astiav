package astiav

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestClass(t *testing.T) {
	c := FindDecoder(CodecIDMjpeg)
	require.NotNil(t, c)
	cc := AllocCodecContext(c)
	require.NotNil(t, cc)
	defer cc.Free()

	cl := cc.Class()
	require.NotNil(t, cl)
	require.Equal(t, ClassCategoryDecoder, cl.Category())
	require.Equal(t, "mjpeg", cl.ItemName())
	require.Equal(t, "AVCodecContext", cl.Name())
	require.Equal(t, fmt.Sprintf("mjpeg [AVCodecContext] @ %p", cc.c), cl.String())
	// TODO Test parent
}

func TestClassers(t *testing.T) {
	cl := len(classers.p)
	f1 := AllocFilterGraph()
	f2 := AllocFilterGraph()
	c := FindDecoder(CodecIDMjpeg)
	require.NotNil(t, c)
	bf := FindBitStreamFilterByName("null")
	require.NotNil(t, bf)
	bfc, err := AllocBitStreamFilterContext(bf)
	require.NoError(t, err)
	cc := AllocCodecContext(c)
	require.NotNil(t, cc)
	bufferSink := FindFilterByName("buffersink")
	require.NotNil(t, bufferSink)
	fc1, err := f1.NewFilterContext(bufferSink, "filter_out", nil)
	require.NoError(t, err)
	_, err = f2.NewFilterContext(bufferSink, "filter_out", nil)
	require.NoError(t, err)
	fmc1 := AllocFormatContext()
	fmc2 := AllocFormatContext()
	require.NoError(t, fmc2.OpenInput("testdata/video.mp4", nil, nil))
	path := filepath.Join(t.TempDir(), "iocontext.txt")
	ic1, err := OpenIOContext(path, NewIOContextFlags(IOContextFlagWrite))
	require.NoError(t, err)
	defer os.RemoveAll(path)
	ic2, err := AllocIOContext(1, true, nil, nil, nil)
	require.NoError(t, err)
	ssc, err := CreateSoftwareScaleContext(1, 1, PixelFormatRgba, 2, 2, PixelFormatRgba, NewSoftwareScaleContextFlags())
	require.NoError(t, err)

	require.Equal(t, cl+12, len(classers.p))
	v, ok := classers.get(unsafe.Pointer(f1.c))
	require.True(t, ok)
	require.Equal(t, f1, v)

	bfc.Free()
	cc.Free()
	fc1.Free()
	f1.Free()
	f2.Free()
	fmc1.Free()
	fmc2.CloseInput()
	require.NoError(t, ic1.Close())
	ic2.Free()
	ssc.Free()
	require.Equal(t, cl, len(classers.p))
}
