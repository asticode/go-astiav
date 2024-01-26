package astiav_test

import (
	"image"
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

type frameDataFrame struct {
	height      int
	imageBytes  []byte
	linesizes   []int
	pixelFormat astiav.PixelFormat
	planesBytes [][]byte
	width       int
}

var _ astiav.FrameDataFrame = (*frameDataFrame)(nil)

func (f *frameDataFrame) Height() int {
	return f.height
}

func (f *frameDataFrame) ImageBufferSize(align int) (int, error) {
	return len(f.imageBytes), nil
}

func (f *frameDataFrame) ImageCopyToBuffer(b []byte, align int) (int, error) {
	copy(b, f.imageBytes)
	return len(f.imageBytes), nil
}

func (f *frameDataFrame) Linesize(i int) int {
	return f.linesizes[i]
}

func (f *frameDataFrame) PixelFormat() astiav.PixelFormat {
	return f.pixelFormat
}

func (f *frameDataFrame) PlaneBytes(i int) []byte {
	return f.planesBytes[i]
}

func (f *frameDataFrame) Width() int {
	return f.width
}

func TestFrameData(t *testing.T) {
	f := &frameDataFrame{}
	fd := astiav.NewFrameData(f)

	for _, v := range []struct {
		err bool
		i   image.Image
		pfs []astiav.PixelFormat
	}{
		{
			i:   &image.Gray{},
			pfs: []astiav.PixelFormat{astiav.PixelFormatGray8},
		},
		{
			i:   &image.Gray16{},
			pfs: []astiav.PixelFormat{astiav.PixelFormatGray16Be},
		},
		{
			i: &image.RGBA{},
			pfs: []astiav.PixelFormat{
				astiav.PixelFormatRgb0,
				astiav.PixelFormat0Rgb,
				astiav.PixelFormatRgb4,
				astiav.PixelFormatRgb8,
			},
		},
		{
			i:   &image.NRGBA{},
			pfs: []astiav.PixelFormat{astiav.PixelFormatRgba},
		},
		{
			i:   &image.NRGBA64{},
			pfs: []astiav.PixelFormat{astiav.PixelFormatRgba64Be},
		},
		{
			i: &image.NYCbCrA{},
			pfs: []astiav.PixelFormat{
				astiav.PixelFormatYuva420P,
				astiav.PixelFormatYuva422P,
				astiav.PixelFormatYuva444P,
			},
		},
		{
			i: &image.YCbCr{},
			pfs: []astiav.PixelFormat{
				astiav.PixelFormatYuv410P,
				astiav.PixelFormatYuv411P,
				astiav.PixelFormatYuvj411P,
				astiav.PixelFormatYuv420P,
				astiav.PixelFormatYuvj420P,
				astiav.PixelFormatYuv422P,
				astiav.PixelFormatYuvj422P,
				astiav.PixelFormatYuv440P,
				astiav.PixelFormatYuvj440P,
				astiav.PixelFormatYuv444P,
				astiav.PixelFormatYuvj444P,
			},
		},
		{
			err: true,
			pfs: []astiav.PixelFormat{astiav.PixelFormatAbgr},
		},
	} {
		for _, pf := range v.pfs {
			f.pixelFormat = pf
			i, err := fd.GuessImageFormat()
			if v.err {
				require.Error(t, err)
			} else {
				require.IsType(t, v.i, i)
			}
		}
	}

	f.imageBytes = []byte{0, 1, 2, 3}
	_, err := fd.Bytes(0)
	require.Error(t, err)
	f.height = 1
	f.width = 2
	b, err := fd.Bytes(0)
	require.NoError(t, err)
	require.Equal(t, f.imageBytes, b)

	for _, v := range []struct {
		e           image.Image
		err         bool
		i           image.Image
		linesizes   []int
		pixelFormat astiav.PixelFormat
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
			pixelFormat: astiav.PixelFormatRgba,
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
			pixelFormat: astiav.PixelFormatRgba,
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
			pixelFormat: astiav.PixelFormatRgba,
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
			pixelFormat: astiav.PixelFormatRgba,
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
			pixelFormat: astiav.PixelFormatRgba,
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
			pixelFormat: astiav.PixelFormatRgba,
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
			pixelFormat: astiav.PixelFormatRgba,
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
			pixelFormat: astiav.PixelFormatYuv444P,
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
			pixelFormat: astiav.PixelFormatRgba,
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
			pixelFormat: astiav.PixelFormatRgba,
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
			pixelFormat: astiav.PixelFormatYuv420P,
			planesBytes: [][]byte{{0, 1}, {2, 3}, {4, 5}},
		},
	} {
		f.linesizes = v.linesizes
		f.pixelFormat = v.pixelFormat
		f.planesBytes = v.planesBytes
		err = fd.ToImage(v.i)
		if v.err {
			require.Error(t, err)
		} else {
			require.Equal(t, v.e, v.i)
		}
	}
}
