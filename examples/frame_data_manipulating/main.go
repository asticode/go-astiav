package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/asticode/go-astiav"
)

const (
	align         = 1
	pngPath       = "testdata/image-rgba.png"
	rawBufferPath = "testdata/image-rgba.rgba"
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

	// Alloc frames
	f1 := astiav.AllocFrame()
	defer f1.Free()
	f2 := astiav.AllocFrame()
	defer f2.Free()

	// To write data manually into a frame, proper attributes need to be set and allocated
	for _, f := range []*astiav.Frame{f1, f2} {
		// Set attributes
		f.SetHeight(256)
		f.SetPixelFormat(astiav.PixelFormatRgba)
		f.SetWidth(256)

		// Alloc buffer
		if err := f.AllocBuffer(align); err != nil {
			log.Fatal(fmt.Errorf("main: allocating buffer failed: %w", err))
		}

		// Alloc image
		if err := f.AllocImage(align); err != nil {
			log.Fatal(fmt.Errorf("main: allocating image failed: %w", err))
		}
	}

	// When writing data manually into a frame, you usually need to make sure the frame is writable
	// Don't forget this step above all if the frame's buffer is referenced elsewhere
	for _, f := range []*astiav.Frame{f1, f2} {
		// Make writable
		if err := f.MakeWritable(); err != nil {
			log.Fatal(fmt.Errorf("main: making frame writable failed: %w", err))
		}
	}

	// As an example, we're going to write data manually into the first frame based on a buffer (i.e. raw data)
	b, err := os.ReadFile(rawBufferPath)
	if err != nil {
		log.Fatal(fmt.Errorf("main: reading %s failed: %w", rawBufferPath, err))
	}
	if err := f1.Data().SetBytes(b, align); err != nil {
		log.Fatal(fmt.Errorf("main: setting frame's data based on bytes failed: %w", err))
	}

	// As an example, we're going to write data manually into the second frame based on a Go image
	fl1, err := os.Open(pngPath)
	if err != nil {
		log.Fatal(fmt.Errorf("main: opening %s failed: %w", pngPath, err))
	}
	defer fl1.Close()
	i1, err := png.Decode(fl1)
	if err != nil {
		log.Fatal(fmt.Errorf("main: decoding %s failed: %w", pngPath, err))
	}
	if err := f2.Data().FromImage(i1); err != nil {
		log.Fatal(fmt.Errorf("main: setting frame's data based on Go image failed: %w", err))
	}

	// This is the place where you do stuff with the frames

	// As an example, we're going to read the first frame's data as a buffer (i.e. raw data)
	if _, err = f1.Data().Bytes(align); err != nil {
		log.Fatal(fmt.Errorf("main: getting frame's data as bytes failed: %w", err))
	}

	// As an example, we're going to read the second frame's data as a Go image
	// For that we first need to guess the Go image format based on the frame's attributes before providing
	// it to .ToImage(). You may not need this and can provide your own image.Image to .ToImage()
	i2, err := f2.Data().GuessImageFormat()
	if err != nil {
		log.Fatal(fmt.Errorf("main: guessing image format failed: %w", err))
	}
	if err := f2.Data().ToImage(i2); err != nil {
		log.Fatal(fmt.Errorf("main: getting frame's data as Go image failed: %w", err))
	}

	// Success
	log.Println("success")
}
