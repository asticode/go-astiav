package astiav

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockedFrameDataFrame struct {
	h           int
	imageBytes  []byte
	linesizes   []int
	pf          PixelFormat
	planesBytes [][]byte
	w           int
}

var _ frameDataFramer = (*mockedFrameDataFrame)(nil)

func (f *mockedFrameDataFrame) height() int {
	return f.h
}

func (f *mockedFrameDataFrame) imageBufferSize(align int) (int, error) {
	return len(f.imageBytes), nil
}

func (f *mockedFrameDataFrame) imageCopyToBuffer(b []byte, align int) (int, error) {
	copy(b, f.imageBytes)
	return len(f.imageBytes), nil
}

func (f *mockedFrameDataFrame) linesize(i int) int {
	return f.linesizes[i]
}

func (f *mockedFrameDataFrame) pixelFormat() PixelFormat {
	return f.pf
}

func (f *mockedFrameDataFrame) planeBytes(i int) []byte {
	return f.planesBytes[i]
}

func (f *mockedFrameDataFrame) width() int {
	return f.w
}

func TestFrameDataInternal(t *testing.T) {
	fdf := &mockedFrameDataFrame{}
	fd := newFrameData(fdf)

	for _, v := range []struct {
		err bool
		i   image.Image
		pfs []PixelFormat
	}{
		{
			i:   &image.Gray{},
			pfs: []PixelFormat{PixelFormatGray8},
		},
		{
			i:   &image.Gray16{},
			pfs: []PixelFormat{PixelFormatGray16Be},
		},
		{
			i: &image.RGBA{},
			pfs: []PixelFormat{
				PixelFormatRgb0,
				PixelFormat0Rgb,
				PixelFormatRgb4,
				PixelFormatRgb8,
			},
		},
		{
			i:   &image.NRGBA{},
			pfs: []PixelFormat{PixelFormatRgba},
		},
		{
			i:   &image.NRGBA64{},
			pfs: []PixelFormat{PixelFormatRgba64Be},
		},
		{
			i: &image.NYCbCrA{},
			pfs: []PixelFormat{
				PixelFormatYuva420P,
				PixelFormatYuva422P,
				PixelFormatYuva444P,
			},
		},
		{
			i: &image.YCbCr{},
			pfs: []PixelFormat{
				PixelFormatYuv410P,
				PixelFormatYuv411P,
				PixelFormatYuvj411P,
				PixelFormatYuv420P,
				PixelFormatYuvj420P,
				PixelFormatYuv422P,
				PixelFormatYuvj422P,
				PixelFormatYuv440P,
				PixelFormatYuvj440P,
				PixelFormatYuv444P,
				PixelFormatYuvj444P,
			},
		},
		{
			err: true,
			pfs: []PixelFormat{PixelFormatAbgr},
		},
	} {
		for _, pf := range v.pfs {
			fdf.pf = pf
			i, err := fd.GuessImageFormat()
			if v.err {
				require.Error(t, err)
			} else {
				require.IsType(t, v.i, i)
			}
		}
	}

	fdf.imageBytes = []byte{0, 1, 2, 3}
	_, err := fd.Bytes(0)
	require.Error(t, err)
	fdf.h = 1
	fdf.w = 2
	b, err := fd.Bytes(0)
	require.NoError(t, err)
	require.Equal(t, fdf.imageBytes, b)

	for _, v := range []struct {
		e           image.Image
		err         bool
		i           image.Image
		linesizes   []int
		pixelFormat PixelFormat
		planesBytes [][]byte
	}{
		{
			e: &image.Alpha{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.Alpha{},
			linesizes:   []int{1},
			pixelFormat: PixelFormatRgba,
			planesBytes: [][]byte{{0, 1, 2, 3}},
		},
		{
			e: &image.Alpha16{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.Alpha16{},
			linesizes:   []int{1},
			pixelFormat: PixelFormatRgba,
			planesBytes: [][]byte{{0, 1, 2, 3}},
		},
		{
			e: &image.CMYK{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.CMYK{},
			linesizes:   []int{1},
			pixelFormat: PixelFormatRgba,
			planesBytes: [][]byte{{0, 1, 2, 3}},
		},
		{
			e: &image.Gray{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.Gray{},
			linesizes:   []int{1},
			pixelFormat: PixelFormatRgba,
			planesBytes: [][]byte{{0, 1, 2, 3}},
		},
		{
			e: &image.Gray16{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.Gray16{},
			linesizes:   []int{1},
			pixelFormat: PixelFormatRgba,
			planesBytes: [][]byte{{0, 1, 2, 3}},
		},
		{
			e: &image.NRGBA{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.NRGBA{},
			linesizes:   []int{1},
			pixelFormat: PixelFormatRgba,
			planesBytes: [][]byte{{0, 1, 2, 3}},
		},
		{
			e: &image.NRGBA64{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.NRGBA64{},
			linesizes:   []int{1},
			pixelFormat: PixelFormatRgba,
			planesBytes: [][]byte{{0, 1, 2, 3}},
		},
		{
			e: &image.NYCbCrA{
				A:       []byte{6, 7},
				AStride: 4,
				YCbCr: image.YCbCr{
					Y:              []byte{0, 1},
					Cb:             []byte{2, 3},
					Cr:             []byte{4, 5},
					YStride:        1,
					CStride:        2,
					SubsampleRatio: image.YCbCrSubsampleRatio444,
					Rect:           image.Rect(0, 0, 2, 1),
				},
			},
			i:           &image.NYCbCrA{},
			linesizes:   []int{1, 2, 3, 4},
			pixelFormat: PixelFormatYuv444P,
			planesBytes: [][]byte{{0, 1}, {2, 3}, {4, 5}, {6, 7}},
		},
		{
			e: &image.RGBA{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.RGBA{},
			linesizes:   []int{1},
			pixelFormat: PixelFormatRgba,
			planesBytes: [][]byte{{0, 1, 2, 3}},
		},
		{
			e: &image.RGBA64{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.RGBA64{},
			linesizes:   []int{1},
			pixelFormat: PixelFormatRgba,
			planesBytes: [][]byte{{0, 1, 2, 3}},
		},
		{
			e: &image.YCbCr{
				Y:              []byte{0, 1},
				Cb:             []byte{2, 3},
				Cr:             []byte{4, 5},
				YStride:        1,
				CStride:        2,
				SubsampleRatio: image.YCbCrSubsampleRatio420,
				Rect:           image.Rect(0, 0, 2, 1),
			},
			i:           &image.YCbCr{},
			linesizes:   []int{1, 2, 3},
			pixelFormat: PixelFormatYuv420P,
			planesBytes: [][]byte{{0, 1}, {2, 3}, {4, 5}},
		},
	} {
		fdf.linesizes = v.linesizes
		fdf.pf = v.pixelFormat
		fdf.planesBytes = v.planesBytes
		err = fd.ToImage(v.i)
		if v.err {
			require.Error(t, err)
		} else {
			require.Equal(t, v.e, v.i)
		}
	}
}

func TestFrameData(t *testing.T) {
	const (
		name = "image-rgba"
		ext  = "png"
	)
	f, err := globalHelper.inputLastFrame(name+"."+ext, MediaTypeVideo)
	require.NoError(t, err)
	fd := f.Data()

	b1, err := fd.Bytes(1)
	require.NoError(t, err)

	b2, err := os.ReadFile("testdata/" + name + "-bytes")
	require.NoError(t, err)
	require.Equal(t, b1, b2)

	f1, err := os.Open("testdata/" + name + "." + ext)
	require.NoError(t, err)
	defer f1.Close()

	i1, err := fd.GuessImageFormat()
	require.NoError(t, err)
	require.NoError(t, err)
	require.NoError(t, fd.ToImage(i1))
	i2, err := png.Decode(f1)
	require.NoError(t, err)
	require.Equal(t, i1, i2)
}
