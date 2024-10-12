package astiav

import (
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockedFrameDataFrame struct {
	copiedPlanes []frameDataPlane
	h            int
	onBytes      func(align int) ([]byte, error)
	onPlanes     func(b []byte, align int) ([]frameDataPlane, error)
	pf           PixelFormat
	w            int
}

var _ frameDataFramer = (*mockedFrameDataFrame)(nil)

func (f *mockedFrameDataFrame) bytes(align int) ([]byte, error) {
	return f.onBytes(align)
}

func (f *mockedFrameDataFrame) copyPlanes(ps []frameDataPlane) error {
	f.copiedPlanes = ps
	return nil
}

func (f *mockedFrameDataFrame) height() int {
	return f.h
}

func (f *mockedFrameDataFrame) pixelFormat() PixelFormat {
	return f.pf
}

func (f *mockedFrameDataFrame) planes(b []byte, align int) ([]frameDataPlane, error) {
	return f.onPlanes(b, align)
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
	b1 := []byte{0, 1, 2, 3}
	fdf.onBytes = func(align int) ([]byte, error) { return b1, nil }
	fdf.w = 2
	b2, err := fd.Bytes(0)
	require.NoError(t, err)
	require.Equal(t, b1, b2)

	for _, v := range []struct {
		e           image.Image
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
		fdf.onPlanes = func(b []byte, align int) ([]frameDataPlane, error) { return v.planes, nil }
		require.NoError(t, fd.ToImage(v.i))
		require.Equal(t, v.e, v.i)
	}

	b1 = []byte{1, 2, 3, 4}
	fdf.onPlanes = func(b []byte, align int) ([]frameDataPlane, error) {
		return []frameDataPlane{
			{
				bytes:    b1[:2],
				linesize: 1,
			},
			{
				bytes:    b1[2:],
				linesize: 2,
			},
		}, nil
	}
	require.NoError(t, fd.SetBytes(b1, 0))
	require.Equal(t, []frameDataPlane{
		{bytes: b1[:2], linesize: 1},
		{bytes: b1[2:], linesize: 2},
	}, fdf.copiedPlanes)

	for _, v := range []struct {
		expectedCopiedPlanes []frameDataPlane
		i                    image.Image
	}{
		{
			expectedCopiedPlanes: []frameDataPlane{{bytes: []byte{0, 1, 2, 3}, linesize: 1}},
			i: &image.Alpha{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
			},
		},
		{
			expectedCopiedPlanes: []frameDataPlane{{bytes: []byte{0, 1, 2, 3}, linesize: 1}},
			i: &image.Alpha16{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
			},
		},
		{
			expectedCopiedPlanes: []frameDataPlane{{bytes: []byte{0, 1, 2, 3}, linesize: 1}},
			i: &image.CMYK{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
			},
		},
		{
			expectedCopiedPlanes: []frameDataPlane{{bytes: []byte{0, 1, 2, 3}, linesize: 1}},
			i: &image.Gray{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
			},
		},
		{
			expectedCopiedPlanes: []frameDataPlane{{bytes: []byte{0, 1, 2, 3}, linesize: 1}},
			i: &image.Gray16{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
			},
		},
		{
			expectedCopiedPlanes: []frameDataPlane{{bytes: []byte{0, 1, 2, 3}, linesize: 1}},
			i: &image.NRGBA{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
			},
		},
		{
			expectedCopiedPlanes: []frameDataPlane{{bytes: []byte{0, 1, 2, 3}, linesize: 1}},
			i: &image.NRGBA64{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
			},
		},
		{
			expectedCopiedPlanes: []frameDataPlane{
				{bytes: []byte{0, 1}, linesize: 1},
				{bytes: []byte{2, 3}, linesize: 2},
				{bytes: []byte{4, 5}, linesize: 2},
				{bytes: []byte{6, 7}, linesize: 4},
			},
			i: &image.NYCbCrA{
				A:       []byte{6, 7},
				AStride: 4,
				YCbCr: image.YCbCr{
					Y:       []byte{0, 1},
					Cb:      []byte{2, 3},
					Cr:      []byte{4, 5},
					YStride: 1,
					CStride: 2,
				},
			},
		},
		{
			expectedCopiedPlanes: []frameDataPlane{{bytes: []byte{0, 1, 2, 3}, linesize: 1}},
			i: &image.RGBA{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
			},
		},
		{
			expectedCopiedPlanes: []frameDataPlane{{bytes: []byte{0, 1, 2, 3}, linesize: 1}},
			i: &image.RGBA64{
				Pix:    []byte{0, 1, 2, 3},
				Stride: 1,
			},
		},
		{
			expectedCopiedPlanes: []frameDataPlane{
				{bytes: []byte{0, 1}, linesize: 1},
				{bytes: []byte{2, 3}, linesize: 2},
				{bytes: []byte{4, 5}, linesize: 2},
			},
			i: &image.YCbCr{
				Y:       []byte{0, 1},
				Cb:      []byte{2, 3},
				Cr:      []byte{4, 5},
				YStride: 1,
				CStride: 2,
			},
		},
	} {
		require.NoError(t, fd.FromImage(v.i))
		require.Equal(t, v.expectedCopiedPlanes, fdf.copiedPlanes)
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
		f1, err := globalHelper.inputLastFrame(v.name+"."+v.ext, MediaTypeVideo)
		require.NoError(t, err)
		fd1 := f1.Data()

		b1, err := fd1.Bytes(1)
		require.NoError(t, err)
		b2 := []byte(fmt.Sprintf("%+v", b1))
		b3, err := os.ReadFile("testdata/" + v.name + "-bytes")
		require.NoError(t, err)
		require.Equal(t, b3, b2)

		i1, err := fd1.GuessImageFormat()
		require.NoError(t, err)
		require.NoError(t, fd1.ToImage(i1))
		b4 := []byte(fmt.Sprintf("%+v", i1))
		b5, err := os.ReadFile("testdata/" + v.name + "-struct")
		require.NoError(t, err)
		require.Equal(t, b5, b4)

		f2 := AllocFrame()
		defer f2.Free()
		f2.SetHeight(f1.Height())
		f2.SetPixelFormat(f1.PixelFormat())
		f2.SetWidth(f1.Width())
		const align = 1
		require.NoError(t, f2.AllocBuffer(align))
		require.NoError(t, f2.AllocImage(align))
		fd2 := f2.Data()

		require.NoError(t, fd2.FromImage(i1))
		b6, err := fd2.Bytes(align)
		require.NoError(t, err)
		b7 := []byte(fmt.Sprintf("%+v", b6))
		require.Equal(t, b3, b7)

		require.NoError(t, f2.ImageFillBlack())
		require.NoError(t, fd2.SetBytes(b1, align))
		b1[0] -= 1
		b8, err := fd2.Bytes(align)
		require.NoError(t, err)
		b9 := []byte(fmt.Sprintf("%+v", b8))
		require.Equal(t, b3, b9)
	}
}
