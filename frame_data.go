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

func (d *FrameData) PlaneData(i int, sizeFunc func(linesize int) int) []byte {
	return bytesFromC(func(size *cUlong) *C.uint8_t {
		*size = cUlong(sizeFunc(int(d.f.c.linesize[i])))
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

func (d *FrameData) imageNRGBA() *image.NRGBA {
	return &image.NRGBA{
		Pix:    d.PlaneData(0, func(linesize int) int { return linesize * d.f.Height() }),
		Stride: d.f.Linesize()[0],
		Rect:   image.Rect(0, 0, d.f.Width(), d.f.Height()),
	}
}

func (d *FrameData) imageYCbCr() *image.YCbCr {
	return &image.YCbCr{
		Y:              d.PlaneData(0, func(linesize int) int { return linesize * d.f.Height() }),
		Cb:             d.PlaneData(1, func(linesize int) int { return linesize * d.f.Height() }),
		Cr:             d.PlaneData(2, func(linesize int) int { return linesize * d.f.Height() }),
		YStride:        d.f.Linesize()[0],
		CStride:        d.f.Linesize()[1],
		SubsampleRatio: d.imageYCbCrSubsampleRatio(),
		Rect:           image.Rect(0, 0, d.f.Width(), d.f.Height()),
	}
}

func (d *FrameData) imageNYCbCrA() *image.NYCbCrA {
	return &image.NYCbCrA{
		YCbCr:   *d.imageYCbCr(),
		A:       d.PlaneData(3, func(linesize int) int { return linesize * d.f.Height() }),
		AStride: d.f.Linesize()[3],
	}
}

func (d *FrameData) Image() (image.Image, error) {
	// Switch on pixel format
	switch d.f.PixelFormat() {
	// NRGBA
	case PixelFormatRgba:
		return d.imageNRGBA(), nil
	// NYCbCrA
	case PixelFormatYuva420P,
		PixelFormatYuva422P,
		PixelFormatYuva444P:
		return d.imageNYCbCrA(), nil
	// YCbCr
	case PixelFormatYuv410P,
		PixelFormatYuv411P, PixelFormatYuvj411P,
		PixelFormatYuv420P, PixelFormatYuvj420P,
		PixelFormatYuv422P, PixelFormatYuvj422P,
		PixelFormatYuv440P, PixelFormatYuvj440P,
		PixelFormatYuv444P, PixelFormatYuvj444P:
		return d.imageYCbCr(), nil
	}
	return nil, fmt.Errorf("astiav: %s pixel format not handled by the Go standard image package", d.f.PixelFormat())
}
