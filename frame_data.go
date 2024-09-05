package astiav

//#include <libavutil/imgutils.h>
//#include "macros.h"
import "C"
import (
	"errors"
	"fmt"
	"image"
	"strings"
)

type FrameData struct {
	f frameDataFramer
}

type frameDataFramer interface {
	bytes(align int) ([]byte, error)
	height() int
	pixelFormat() PixelFormat
	planes() ([]frameDataPlane, error)
	width() int
}

type frameDataPlane struct {
	bytes    []byte
	linesize int
}

func newFrameData(f frameDataFramer) *FrameData {
	return &FrameData{f: f}
}

func (d *FrameData) Bytes(align int) ([]byte, error) {
	return d.f.bytes(align)
}

// Always returns non-premultiplied formats when dealing with alpha channels, however this might not
// always be accurate. In this case, use your own format in .ToImage()
func (d *FrameData) GuessImageFormat() (image.Image, error) {
	switch d.f.pixelFormat() {
	case PixelFormatGray8:
		return &image.Gray{}, nil
	case PixelFormatGray16Be:
		return &image.Gray16{}, nil
	case PixelFormatRgb0, PixelFormat0Rgb, PixelFormatRgb4, PixelFormatRgb8:
		return &image.RGBA{}, nil
	case PixelFormatRgba:
		return &image.NRGBA{}, nil
	case PixelFormatRgba64Be:
		return &image.NRGBA64{}, nil
	case PixelFormatYuva420P,
		PixelFormatYuva422P,
		PixelFormatYuva444P:
		return &image.NYCbCrA{}, nil
	case PixelFormatYuv410P,
		PixelFormatYuv411P, PixelFormatYuvj411P,
		PixelFormatYuv420P, PixelFormatYuvj420P,
		PixelFormatYuv422P, PixelFormatYuvj422P,
		PixelFormatYuv440P, PixelFormatYuvj440P,
		PixelFormatYuv444P, PixelFormatYuvj444P:
		return &image.YCbCr{}, nil
	}
	return nil, fmt.Errorf("astiav: pixel format %s not handled by Go", d.f.pixelFormat())
}

func (d *FrameData) imageYCbCrSubsampleRatio() image.YCbCrSubsampleRatio {
	name := d.f.pixelFormat().Name()
	for s, r := range map[string]image.YCbCrSubsampleRatio{
		"410": image.YCbCrSubsampleRatio410,
		"411": image.YCbCrSubsampleRatio411,
		"420": image.YCbCrSubsampleRatio420,
		"422": image.YCbCrSubsampleRatio422,
		"440": image.YCbCrSubsampleRatio440,
		"444": image.YCbCrSubsampleRatio444,
	} {
		if strings.Contains(name, s) {
			return r
		}
	}
	return image.YCbCrSubsampleRatio444
}

func (d *FrameData) toImagePix(pix *[]uint8, stride *int, rect *image.Rectangle, planes []frameDataPlane) {
	*pix = planes[0].bytes
	if v := planes[0].linesize; *stride != v {
		*stride = v
	}
	if w, h := d.f.width(), d.f.height(); rect.Dy() != w || rect.Dx() != h {
		*rect = image.Rect(0, 0, w, h)
	}
}

func (d *FrameData) toImageYCbCr(y, cb, cr *[]uint8, yStride, cStride *int, subsampleRatio *image.YCbCrSubsampleRatio, rect *image.Rectangle, planes []frameDataPlane) {
	*y = planes[0].bytes
	*cb = planes[1].bytes
	*cr = planes[2].bytes
	if v := planes[0].linesize; *yStride != v {
		*yStride = v
	}
	if v := planes[1].linesize; *cStride != v {
		*cStride = v
	}
	if v := d.imageYCbCrSubsampleRatio(); *subsampleRatio != v {
		*subsampleRatio = v
	}
	if w, h := d.f.width(), d.f.height(); rect.Dy() != w || rect.Dx() != h {
		*rect = image.Rect(0, 0, w, h)
	}
}

func (d *FrameData) toImageYCbCrA(y, cb, cr, a *[]uint8, yStride, cStride, aStride *int, subsampleRatio *image.YCbCrSubsampleRatio, rect *image.Rectangle, planes []frameDataPlane) {
	d.toImageYCbCr(y, cb, cr, yStride, cStride, subsampleRatio, rect, planes)
	*a = planes[3].bytes
	if v := planes[3].linesize; *aStride != v {
		*aStride = v
	}
}

func (d *FrameData) ToImage(dst image.Image) error {
	// Get planes
	planes, err := d.f.planes()
	if err != nil {
		return fmt.Errorf("astiav: getting planes failed: %w", err)
	}

	// Update image
	if v, ok := dst.(*image.Alpha); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect, planes)
	} else if v, ok := dst.(*image.Alpha16); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect, planes)
	} else if v, ok := dst.(*image.CMYK); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect, planes)
	} else if v, ok := dst.(*image.Gray); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect, planes)
	} else if v, ok := dst.(*image.Gray16); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect, planes)
	} else if v, ok := dst.(*image.NRGBA); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect, planes)
	} else if v, ok := dst.(*image.NRGBA64); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect, planes)
	} else if v, ok := dst.(*image.NYCbCrA); ok {
		d.toImageYCbCrA(&v.Y, &v.Cb, &v.Cr, &v.A, &v.YStride, &v.CStride, &v.AStride, &v.SubsampleRatio, &v.Rect, planes)
	} else if v, ok := dst.(*image.RGBA); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect, planes)
	} else if v, ok := dst.(*image.RGBA64); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect, planes)
	} else if v, ok := dst.(*image.YCbCr); ok {
		d.toImageYCbCr(&v.Y, &v.Cb, &v.Cr, &v.YStride, &v.CStride, &v.SubsampleRatio, &v.Rect, planes)
	} else {
		return errors.New("astiav: image format is not handled")
	}
	return nil
}

var _ frameDataFramer = (*frameDataFrame)(nil)

type frameDataFrame struct {
	f *Frame
}

func newFrameDataFrame(f *Frame) *frameDataFrame {
	return &frameDataFrame{f: f}
}

func (f *frameDataFrame) bytes(align int) ([]byte, error) {
	switch {
	// Video
	case f.height() > 0 && f.width() > 0:
		// Get buffer size
		s, err := f.f.ImageBufferSize(align)
		if err != nil {
			return nil, fmt.Errorf("astiav: getting image buffer size failed: %w", err)
		}

		// Invalid buffer size
		if s == 0 {
			return nil, errors.New("astiav: invalid image buffer size")
		}

		// Create buffer
		b := make([]byte, s)

		// Copy image to buffer
		if _, err = f.f.ImageCopyToBuffer(b, align); err != nil {
			return nil, fmt.Errorf("astiav: copying image to buffer failed: %w", err)
		}
		return b, nil
	}
	return nil, errors.New("astiav: frame type not implemented")
}

func (f *frameDataFrame) height() int {
	return f.f.Height()
}

func (f *frameDataFrame) pixelFormat() PixelFormat {
	return f.f.PixelFormat()
}

// Using bytesFromC on f.c.data caused random segfaults
func (f *frameDataFrame) planes() ([]frameDataPlane, error) {
	// Get bytes
	const align = 1
	b, err := f.bytes(align)
	if err != nil {
		return nil, fmt.Errorf("astiav: getting bytes failed: %w", err)
	}

	switch {
	// Video
	case f.height() > 0 && f.width() > 0:
		// Below is mostly inspired by https://github.com/FFmpeg/FFmpeg/blob/n5.1.2/libavutil/imgutils.c#L466

		// Get linesize
		var linesize [4]C.int
		if err := newError(C.av_image_fill_linesizes(&linesize[0], (C.enum_AVPixelFormat)(f.f.c.format), f.f.c.width)); err != nil {
			return nil, fmt.Errorf("astiav: getting linesize failed: %w", err)
		}

		// Align linesize
		var alignedLinesize [4]C.ptrdiff_t
		for i := 0; i < 4; i++ {
			alignedLinesize[i] = C.astiavFFAlign(linesize[i], C.int(align))
		}

		// Get plane sizes
		var planeSizes [4]C.size_t
		if err := newError(C.av_image_fill_plane_sizes(&planeSizes[0], (C.enum_AVPixelFormat)(f.f.c.format), f.f.c.height, &alignedLinesize[0])); err != nil {
			return nil, fmt.Errorf("astiav: getting plane sizes failed: %w", err)
		}

		// Loop through plane sizes
		var ps []frameDataPlane
		start := 0
		for idx, planeSize := range planeSizes {
			// Get end
			end := start + int(planeSize)
			if len(b) < end {
				return nil, fmt.Errorf("astiav: buffer length %d is invalid for [%d:%d]", len(b), start, end)
			}

			// Append plane
			ps = append(ps, frameDataPlane{
				bytes:    b[start:end],
				linesize: int(alignedLinesize[idx]),
			})

			// Update start
			start += int(planeSize)
		}
		return ps, nil
	}
	return nil, errors.New("astiav: frame type not implemented")
}

func (f *frameDataFrame) width() int {
	return f.f.Width()
}
