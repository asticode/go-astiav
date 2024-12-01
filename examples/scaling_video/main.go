package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/asticode/go-astiav"
)

var (
	output    = flag.String("o", "", "the png output path")
	dstWidth  = flag.Int("w", 50, "destination width")
	dstHeight = flag.Int("h", 50, "destination height")
)

func main() {
	// Handle ffmpeg logs
	astiav.SetLogLevel(astiav.LogLevelDebug)
	astiav.SetLogCallback(func(c astiav.Classer, l astiav.LogLevel, fmt, msg string) {
		var cs string
		if c != nil {
			if cl := c.Class(); cl != nil {
				cs = " - class: " + cl.String()
			}
		}
		log.Printf("ffmpeg log: %s%s - level: %d\n", strings.TrimSpace(msg), cs, l)
	})

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

	// Create software scale context
	swsCtx, err := astiav.CreateSoftwareScaleContext(
		srcFrame.Width(),
		srcFrame.Height(),
		srcFrame.PixelFormat(),
		*dstWidth,
		*dstHeight,
		astiav.PixelFormatRgba,
		astiav.NewSoftwareScaleContextFlags(astiav.SoftwareScaleContextFlagBilinear),
	)
	if err != nil {
		log.Fatal(fmt.Errorf("main: creating software scale context failed: %w", err))
	}
	defer swsCtx.Free()

	// Scale frame
	if err := swsCtx.ScaleFrame(srcFrame, dstFrame); err != nil {
		log.Fatal(fmt.Errorf("main: scaling frame failed: %w", err))
	}

	// Guess destination image format
	img, err := dstFrame.Data().GuessImageFormat()
	if err != nil {
		log.Fatal(fmt.Errorf("main: guessing destination image format failed: %w", err))
	}

	// Copy frame data to destination image
	if err = dstFrame.Data().ToImage(img); err != nil {
		log.Fatal(fmt.Errorf("main: copying frame data to destination image failed: %w", err))
	}

	// Encode to png
	if err = png.Encode(dstFile, img); err != nil {
		log.Fatal(fmt.Errorf("main: encoding to png failed: %w", err))
	}

	// Success
	log.Println("success")
}
