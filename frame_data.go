package astiav

//#include <stdint.h>
import "C"
import (
	"errors"
	"fmt"
	"image"
	"strings"
)

type FrameDataFrame interface {
	Height() int
	ImageBufferSize(align int) (int, error)
	ImageCopyToBuffer(b []byte, align int) (int, error)
	Linesize(i int) int
	PixelFormat() PixelFormat
	PlaneBytes(i int) []byte
	Width() int
}

var _ FrameDataFrame = (*frameDataFrame)(nil)

type frameDataFrame struct {
	f *Frame
}

func newFrameDataFrame(f *Frame) *frameDataFrame {
	return &frameDataFrame{f: f}
}

func (f *frameDataFrame) Height() int {
	return f.f.Height()
}

func (f *frameDataFrame) ImageBufferSize(align int) (int, error) {
	return f.f.ImageBufferSize(align)
}

func (f *frameDataFrame) ImageCopyToBuffer(b []byte, align int) (int, error) {
	return f.f.ImageCopyToBuffer(b, align)
}

func (f *frameDataFrame) Linesize(i int) int {
	return f.f.Linesize()[i]
}

func (f *frameDataFrame) PixelFormat() PixelFormat {
	return f.f.PixelFormat()
}

func (f *frameDataFrame) PlaneBytes(i int) []byte {
	return bytesFromC(func(size *cUlong) *C.uint8_t {
		*size = cUlong(int(f.f.c.linesize[i]) * f.f.Height())
		return f.f.c.data[i]
	})
}

func (f *frameDataFrame) Width() int {
	return f.f.Width()
}

type FrameData struct {
	f FrameDataFrame
}

func NewFrameData(f FrameDataFrame) *FrameData {
	return &FrameData{f: f}
}

func (d *FrameData) Bytes(align int) ([]byte, error) {
	switch {
	// Video
	case d.f.Height() > 0 && d.f.Width() > 0:
		// Get buffer size
		s, err := d.f.ImageBufferSize(align)
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
		if _, err = d.f.ImageCopyToBuffer(b, align); err != nil {
			return nil, fmt.Errorf("astiav: copying image to buffer failed: %w", err)
		}
		return b, nil
	}
	return nil, errors.New("astiav: frame type not implemented")
}

// Always returns non-premultiplied formats when dealing with alpha channels, however this might not
// always be accurate. In this case, use your own format in .ToImage()
func (d *FrameData) GuessImageFormat() (image.Image, error) {
	switch d.f.PixelFormat() {
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
	return nil, fmt.Errorf("astiav: pixel format %s not handled by Go", d.f.PixelFormat())
}

func (d *FrameData) imageYCbCrSubsampleRatio() image.YCbCrSubsampleRatio {
	name := d.f.PixelFormat().Name()
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

func (d *FrameData) copyPlaneBytes(i int, s *[]uint8) {
	b := d.f.PlaneBytes(i)
	if len(b) > cap(*s) {
		*s = make([]uint8, len(b))
	}
	copy(*s, b)
}

func (d *FrameData) toImagePix(pix *[]uint8, stride *int, rect *image.Rectangle) {
	d.copyPlaneBytes(0, pix)
	if v := d.f.Linesize(0); *stride != v {
		*stride = v
	}
	if w, h := d.f.Width(), d.f.Height(); rect.Dy() != w || rect.Dx() != h {
		*rect = image.Rect(0, 0, w, h)
	}
}

func (d *FrameData) toImageYCbCr(y, cb, cr *[]uint8, yStride, cStride *int, subsampleRatio *image.YCbCrSubsampleRatio, rect *image.Rectangle) {
	d.copyPlaneBytes(0, y)
	d.copyPlaneBytes(1, cb)
	d.copyPlaneBytes(2, cr)
	if v := d.f.Linesize(0); *yStride != v {
		*yStride = v
	}
	if v := d.f.Linesize(1); *cStride != v {
		*cStride = v
	}
	if v := d.imageYCbCrSubsampleRatio(); *subsampleRatio != v {
		*subsampleRatio = v
	}
	if w, h := d.f.Width(), d.f.Height(); rect.Dy() != w || rect.Dx() != h {
		*rect = image.Rect(0, 0, w, h)
	}
}

func (d *FrameData) toImageYCbCrA(y, cb, cr, a *[]uint8, yStride, cStride, aStride *int, subsampleRatio *image.YCbCrSubsampleRatio, rect *image.Rectangle) {
	d.toImageYCbCr(y, cb, cr, yStride, cStride, subsampleRatio, rect)
	d.copyPlaneBytes(3, a)
	if v := d.f.Linesize(3); *aStride != v {
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
