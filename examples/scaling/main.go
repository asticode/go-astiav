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
		dstFilename string
		dstWidth    int
		dstHeight   int
	)

	flag.StringVar(&dstFilename, "output", "", "Output file name")
	flag.IntVar(&dstWidth, "w", 0, "Destination width")
	flag.IntVar(&dstHeight, "h", 0, "Destination height")
	flag.Parse()

	if dstFilename == "" || dstWidth <= 0 || dstHeight <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s -output output_file -w W -h H\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	dstFile, err := os.Create(dstFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open destination file %s\n", dstFilename)
		os.Exit(1)
	}
	defer dstFile.Close()

	srcW, srcH := 320, 240
	srcPixFmt, dstPixFmt := astiav.PixelFormatYuv420P, astiav.PixelFormatRgba
	srcFrame := astiav.AllocFrame()
	srcFrame.SetHeight(srcH)
	srcFrame.SetWidth(srcW)
	srcFrame.SetPixelFormat(srcPixFmt)
	srcFrame.AllocBuffer(1)
	srcFrame.ImageFillBlack()
	defer srcFrame.Free()

	dstFrame := astiav.AllocFrame()
	defer dstFrame.Free()

	swsCtx := astiav.SwsGetContext(srcW, srcH, srcPixFmt, dstWidth, dstHeight, dstPixFmt, astiav.SWS_POINT, dstFrame)
	if swsCtx == nil {
		fmt.Fprintln(os.Stderr, "Unable to create scale context")
		os.Exit(1)
	}
	defer swsCtx.Free()

	err = swsCtx.Scale(srcFrame, dstFrame)
	if err != nil {
                log.Fatalf("Unable to scale: %s", err.Error())
        }

        img, err := dstFrame.Data().Image()
        if err != nil {
                log.Fatalf("Unable to get image: %s", err.Error())
        }

        err = png.Encode(dstFile, img)
        if err != nil {
                log.Fatalf("Unable to encode image to png: %s", err.Error())
        }

	log.Printf("Successfully scale to %dx%d and write image to: %s", dstWidth, dstHeight, dstFilename)
}
