package astiav

//#include <stdint.h>
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
	height() int
	imageBufferSize(align int) (int, error)
	imageCopyToBuffer(b []byte, align int) (int, error)
	linesize(i int) int
	pixelFormat() PixelFormat
	planeBytes(i int) []byte
	width() int
}

func newFrameData(f frameDataFramer) *FrameData {
	return &FrameData{f: f}
}

func (d *FrameData) Bytes(align int) ([]byte, error) {
	switch {
	// Video
	case d.f.height() > 0 && d.f.width() > 0:
		// Get buffer size
		s, err := d.f.imageBufferSize(align)
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
		if _, err = d.f.imageCopyToBuffer(b, align); err != nil {
			return nil, fmt.Errorf("astiav: copying image to buffer failed: %w", err)
		}
		return b, nil
	}
	return nil, errors.New("astiav: frame type not implemented")
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

func (d *FrameData) toImagePix(pix *[]uint8, stride *int, rect *image.Rectangle) {
	*pix = d.f.planeBytes(0)
	if v := d.f.linesize(0); *stride != v {
		*stride = v
	}
	if w, h := d.f.width(), d.f.height(); rect.Dy() != w || rect.Dx() != h {
		*rect = image.Rect(0, 0, w, h)
	}
}

func (d *FrameData) toImageYCbCr(y, cb, cr *[]uint8, yStride, cStride *int, subsampleRatio *image.YCbCrSubsampleRatio, rect *image.Rectangle) {
	*y = d.f.planeBytes(0)
	*cb = d.f.planeBytes(1)
	*cr = d.f.planeBytes(2)
	if v := d.f.linesize(0); *yStride != v {
		*yStride = v
	}
	if v := d.f.linesize(1); *cStride != v {
		*cStride = v
	}
	if v := d.imageYCbCrSubsampleRatio(); *subsampleRatio != v {
		*subsampleRatio = v
	}
	if w, h := d.f.width(), d.f.height(); rect.Dy() != w || rect.Dx() != h {
		*rect = image.Rect(0, 0, w, h)
	}
}

func (d *FrameData) toImageYCbCrA(y, cb, cr, a *[]uint8, yStride, cStride, aStride *int, subsampleRatio *image.YCbCrSubsampleRatio, rect *image.Rectangle) {
	d.toImageYCbCr(y, cb, cr, yStride, cStride, subsampleRatio, rect)
	*a = d.f.planeBytes(3)
	if v := d.f.linesize(3); *aStride != v {
		*aStride = v
	}
}

func (d *FrameData) ToImage(dst image.Image) error {
	if v, ok := dst.(*image.Alpha); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect)
	} else if v, ok := dst.(*image.Alpha16); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect)
	} else if v, ok := dst.(*image.CMYK); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect)
	} else if v, ok := dst.(*image.Gray); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect)
	} else if v, ok := dst.(*image.Gray16); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect)
	} else if v, ok := dst.(*image.NRGBA); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect)
	} else if v, ok := dst.(*image.NRGBA64); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect)
	} else if v, ok := dst.(*image.NYCbCrA); ok {
		d.toImageYCbCrA(&v.Y, &v.Cb, &v.Cr, &v.A, &v.YStride, &v.CStride, &v.AStride, &v.SubsampleRatio, &v.Rect)
	} else if v, ok := dst.(*image.RGBA); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect)
	} else if v, ok := dst.(*image.RGBA64); ok {
		d.toImagePix(&v.Pix, &v.Stride, &v.Rect)
	} else if v, ok := dst.(*image.YCbCr); ok {
		d.toImageYCbCr(&v.Y, &v.Cb, &v.Cr, &v.YStride, &v.CStride, &v.SubsampleRatio, &v.Rect)
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

func (f *frameDataFrame) height() int {
	return f.f.Height()
}

func (f *frameDataFrame) imageBufferSize(align int) (int, error) {
	return f.f.ImageBufferSize(align)
}

func (f *frameDataFrame) imageCopyToBuffer(b []byte, align int) (int, error) {
	return f.f.ImageCopyToBuffer(b, align)
}

func (f *frameDataFrame) linesize(i int) int {
	return f.f.Linesize()[i]
}

func (f *frameDataFrame) pixelFormat() PixelFormat {
	return f.f.PixelFormat()
}

func (f *frameDataFrame) planeBytes(i int) []byte {
	return bytesFromC(func(size *C.size_t) *C.uint8_t {
		*size = C.size_t(int(f.f.c.linesize[i]) * f.f.Height())
		return f.f.c.data[i]
	})
}

func (f *frameDataFrame) width() int {
	return f.f.Width()
}
