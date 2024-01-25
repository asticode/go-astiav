package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/asticode/go-astiav"
)

func main() {

	var (
		output    = flag.String("o", "", "the png output path")
		dstWidth  = flag.Int("w", 50, "destination width")
		dstHeight = flag.Int("h", 50, "destination height")
	)

	// Parse flags
	flag.Parse()

	// Usage
	if *output == "" || *dstWidth <= 0 || *dstHeight <= 0 {
		log.Println("Usage: <binary path> -o <output path> -w <output width> -h <output height>")
		return
	}

	// Create destination file
	dstFile, err := os.Create(*output)
	if err != nil {
		log.Fatal(fmt.Errorf("main: creating %s failed: %w", *output, err))
	}
	defer dstFile.Close()

	// Create source frame
	srcFrame := astiav.AllocFrame()
	defer srcFrame.Free()
	srcFrame.SetWidth(320)
	srcFrame.SetHeight(240)
	srcFrame.SetPixelFormat(astiav.PixelFormatYuv420P)
	if err = srcFrame.AllocBuffer(1); err != nil {
		log.Fatal(fmt.Errorf("main: allocating source frame buffer failed: %w", err))
	}
	if err = srcFrame.ImageFillBlack(); err != nil {
		log.Fatal(fmt.Errorf("main: filling source frame with black image failed: %w", err))
	}

	// Create destination frame
	dstFrame := astiav.AllocFrame()
	defer dstFrame.Free()

	// Create software scale context flags
	swscf := astiav.NewSoftwareScaleContextFlags(astiav.SoftwareScaleContextBilinear)

	// Create software scale context
	swsCtx := astiav.NewSoftwareScaleContext(srcFrame.Width(), srcFrame.Height(), srcFrame.PixelFormat(), *dstWidth, *dstHeight, astiav.PixelFormatRgba, swscf)
	if swsCtx == nil {
		log.Fatal("main: creating software scale context failed")
	}
	defer swsCtx.Free()

	// Prepare destination frame (Width, Height and Buffer for correct scaling would be set)
	if err = swsCtx.PrepareDestinationFrameForScaling(dstFrame); err != nil {
		log.Fatal(fmt.Errorf("main: prepare destination image failed: %w", err))
	}

	// Scale frame
	if output_slice_height := swsCtx.ScaleFrame(srcFrame, dstFrame); output_slice_height != *dstHeight {
		log.Fatal(fmt.Errorf("main: scale error, expected output slice height %d, but got %d", *dstHeight, output_slice_height))
	}

	// Get image
	img, err := dstFrame.Data().Image()
	if err != nil {
		log.Fatal(fmt.Errorf("main: getting destination image failed: %w", err))
	}

	// Encode to png
	if err = png.Encode(dstFile, img); err != nil {
		log.Fatal(fmt.Errorf("main: encoding to png failed: %w", err))
	}

	log.Println("done")
}
