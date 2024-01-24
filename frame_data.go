package astiav

//#include <stdint.h>
import "C"
import (
	"errors"
	"fmt"
	"image"
	"strings"
)

type frameDataImageFormat int

const (
	frameDataImageFormatNone frameDataImageFormat = iota
	frameDataImageFormatNRGBA
	frameDataImageFormatNYCbCrA
	frameDataImageFormatYCbCr
)

func frameDataImageFormatFromPixelFormat(pf PixelFormat) frameDataImageFormat {
	// Switch on pixel format
	switch pf {
	// NRGBA
	case PixelFormatRgba:
		return frameDataImageFormatNRGBA
	// NYCbCrA
	case PixelFormatYuva420P,
		PixelFormatYuva422P,
		PixelFormatYuva444P:
		return frameDataImageFormatNYCbCrA
	// YCbCr
	case PixelFormatYuv410P,
		PixelFormatYuv411P, PixelFormatYuvj411P,
		PixelFormatYuv420P, PixelFormatYuvj420P,
		PixelFormatYuv422P, PixelFormatYuvj422P,
		PixelFormatYuv440P, PixelFormatYuvj440P,
		PixelFormatYuv444P, PixelFormatYuvj444P:
		return frameDataImageFormatYCbCr
	}
	return frameDataImageFormatNone
}

type FrameData struct {
	f *Frame
}

func newFrameData(f *Frame) *FrameData {
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

func (d *FrameData) planeBytes(i int) []byte {
	return bytesFromC(func(size *cUlong) *C.uint8_t {
		*size = cUlong(int(d.f.c.linesize[i]) * d.f.Height())
		return d.f.c.data[i]
	})
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
	b := d.planeBytes(0)
	if len(b) > cap(*s) {
		*s = make([]uint8, len(b))
	}
	copy(*s, b)
}

func (d *FrameData) toImageNRGBA(i *image.NRGBA) {
	d.copyPlaneBytes(0, &i.Pix)
	if v := d.f.Linesize()[0]; i.Stride != v {
		i.Stride = v
	}
	if w, h := d.f.Width(), d.f.Height(); i.Rect.Dy() != w || i.Rect.Dx() != h {
		i.Rect = image.Rect(0, 0, w, h)
	}
}

func (d *FrameData) toImageYCbCr(i *image.YCbCr) {
	d.copyPlaneBytes(0, &i.Y)
	d.copyPlaneBytes(1, &i.Cb)
	d.copyPlaneBytes(2, &i.Cr)
	if v := d.f.Linesize()[0]; i.YStride != v {
		i.YStride = v
	}
	if v := d.f.Linesize()[1]; i.CStride != v {
		i.CStride = v
	}
	if v := d.imageYCbCrSubsampleRatio(); i.SubsampleRatio != v {
		i.SubsampleRatio = v
	}
	if w, h := d.f.Width(), d.f.Height(); i.Rect.Dy() != w || i.Rect.Dx() != h {
		i.Rect = image.Rect(0, 0, w, h)
	}
}

func (d *FrameData) toImageNYCbCrA(i *image.NYCbCrA) {
	d.toImageYCbCr(&i.YCbCr)
	d.copyPlaneBytes(3, &i.A)
	if v := d.f.Linesize()[3]; i.AStride != v {
		i.AStride = v
	}
}

func (d *FrameData) Image() (image.Image, error) {
	// Switch on image format
	switch frameDataImageFormatFromPixelFormat(d.f.PixelFormat()) {
	// NRGBA
	case frameDataImageFormatNRGBA:
		i := &image.NRGBA{}
		d.toImageNRGBA(i)
		return i, nil
	// NYCbCrA
	case frameDataImageFormatNYCbCrA:
		i := &image.NYCbCrA{}
		d.toImageNYCbCrA(i)
		return i, nil
	// YCbCr
	case frameDataImageFormatYCbCr:
		i := &image.YCbCr{}
		d.toImageYCbCr(i)
		return i, nil
	}
	return nil, fmt.Errorf("astiav: %s pixel format not handled by the Go standard image package", d.f.PixelFormat())
}

func (d *FrameData) ToImage(dst image.Image) error {
	// Switch on image format
	switch frameDataImageFormatFromPixelFormat(d.f.PixelFormat()) {
	// NRGBA
	case frameDataImageFormatNRGBA:
		i, ok := dst.(*image.NRGBA)
		if !ok {
			return errors.New("astiav: image should be *image.NRGBA")
		}
		d.toImageNRGBA(i)
		return nil
	// NYCbCrA
	case frameDataImageFormatNYCbCrA:
		i, ok := dst.(*image.NYCbCrA)
		if !ok {
			return errors.New("astiav: image should be *image.NYCbCrA")
		}
		d.toImageNYCbCrA(i)
		return nil
	// YCbCr
	case frameDataImageFormatYCbCr:
		i, ok := dst.(*image.YCbCr)
		if !ok {
			return errors.New("astiav: image should be *image.YCbCr")
		}
		d.toImageYCbCr(i)
		return nil
	}
	return fmt.Errorf("astiav: %s pixel format not handled by the Go standard image package", d.f.PixelFormat())
}
