package astiav

import (
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockedFrameDataFrame struct {
	h          int
	imageBytes []byte
	pf         PixelFormat
	planes_    []frameDataPlane
	w          int
}

var _ frameDataFramer = (*mockedFrameDataFrame)(nil)

func (f *mockedFrameDataFrame) bytes(align int) ([]byte, error) {
	return f.imageBytes, nil
}

func (f *mockedFrameDataFrame) height() int {
	return f.h
}

func (f *mockedFrameDataFrame) pixelFormat() PixelFormat {
	return f.pf
}

func (f *mockedFrameDataFrame) planes() ([]frameDataPlane, error) {
	return f.planes_, nil
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

	fdf.h = 1
	fdf.imageBytes = []byte{0, 1, 2, 3}
	fdf.w = 2
	b, err := fd.Bytes(0)
	require.NoError(t, err)
	require.Equal(t, fdf.imageBytes, b)

	for _, v := range []struct {
		e           image.Image
		err         bool
		i           image.Image
		pixelFormat PixelFormat
		planes      []frameDataPlane
	}{
		{
			e: &image.Alpha{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.Alpha{},
			pixelFormat: PixelFormatRgba,
			planes: []frameDataPlane{
				{
					bytes:    []byte{0, 1, 2, 3},
					linesize: 1,
				},
			},
		},
		{
			e: &image.Alpha16{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.Alpha16{},
			pixelFormat: PixelFormatRgba,
			planes: []frameDataPlane{
				{
					bytes:    []byte{0, 1, 2, 3},
					linesize: 1,
				},
			},
		},
		{
			e: &image.CMYK{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.CMYK{},
			pixelFormat: PixelFormatRgba,
			planes: []frameDataPlane{
				{
					bytes:    []byte{0, 1, 2, 3},
					linesize: 1,
				},
			},
		},
		{
			e: &image.Gray{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.Gray{},
			pixelFormat: PixelFormatRgba,
			planes: []frameDataPlane{
				{
					bytes:    []byte{0, 1, 2, 3},
					linesize: 1,
				},
			},
		},
		{
			e: &image.Gray16{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.Gray16{},
			pixelFormat: PixelFormatRgba,
			planes: []frameDataPlane{
				{
					bytes:    []byte{0, 1, 2, 3},
					linesize: 1,
				},
			},
		},
		{
			e: &image.NRGBA{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.NRGBA{},
			pixelFormat: PixelFormatRgba,
			planes: []frameDataPlane{
				{
					bytes:    []byte{0, 1, 2, 3},
					linesize: 1,
				},
			},
		},
		{
			e: &image.NRGBA64{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.NRGBA64{},
			pixelFormat: PixelFormatRgba,
			planes: []frameDataPlane{
				{
					bytes:    []byte{0, 1, 2, 3},
					linesize: 1,
				},
			},
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
			pixelFormat: PixelFormatYuv444P,
			planes: []frameDataPlane{
				{
					bytes:    []byte{0, 1},
					linesize: 1,
				},
				{
					bytes:    []byte{2, 3},
					linesize: 2,
				},
				{
					bytes:    []byte{4, 5},
					linesize: 3,
				},
				{
					bytes:    []byte{6, 7},
					linesize: 4,
				},
			},
		},
		{
			e: &image.RGBA{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.RGBA{},
			pixelFormat: PixelFormatRgba,
			planes: []frameDataPlane{
				{
					bytes:    []byte{0, 1, 2, 3},
					linesize: 1,
				},
			},
		},
		{
			e: &image.RGBA64{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
				Rect:   image.Rect(0, 0, 2, 1),
			},
			i:           &image.RGBA64{},
			pixelFormat: PixelFormatRgba,
			planes: []frameDataPlane{
				{
					bytes:    []byte{0, 1, 2, 3},
					linesize: 1,
				},
			},
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
			pixelFormat: PixelFormatYuv420P,
			planes: []frameDataPlane{
				{
					bytes:    []byte{0, 1},
					linesize: 1,
				},
				{
					bytes:    []byte{2, 3},
					linesize: 2,
				},
				{
					bytes:    []byte{4, 5},
					linesize: 3,
				},
			},
		},
	} {
		fdf.pf = v.pixelFormat
		fdf.planes_ = v.planes
		err = fd.ToImage(v.i)
		if v.err {
			require.Error(t, err)
		} else {
			require.Equal(t, v.e, v.i)
		}
	}
}

func TestFrameData(t *testing.T) {
	for _, v := range []struct {
		ext  string
		name string
	}{
		{
			ext:  "png",
			name: "image-rgba",
		},
		{
			ext:  "h264",
			name: "video-yuv420p",
		},
	} {
		f, err := globalHelper.inputLastFrame(v.name+"."+v.ext, MediaTypeVideo)
		require.NoError(t, err)
		fd := f.Data()

		b1, err := fd.Bytes(1)
		require.NoError(t, err)
		b2 := []byte(fmt.Sprintf("%+v", b1))
		b3, err := os.ReadFile("testdata/" + v.name + "-bytes")
		require.NoError(t, err)
		require.Equal(t, b2, b3)

		i1, err := fd.GuessImageFormat()
		require.NoError(t, err)
		require.NoError(t, fd.ToImage(i1))
		b4 := []byte(fmt.Sprintf("%+v", i1))
		b5, err := os.ReadFile("testdata/" + v.name + "-struct")
		require.NoError(t, err)
		require.Equal(t, b4, b5)

	}
}
