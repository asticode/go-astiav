package main

import (
	"fmt"
	"image"
	"log"
	"strings"

	"github.com/asticode/go-astiav"
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

	/*

		In this first part we're going to manipulate an audio frame

	*/

	// Alloc frame
	audioFrame := astiav.AllocFrame()
	defer audioFrame.Free()

	// To write data manually into a frame, proper attributes need to be set and allocated
	audioFrame.SetChannelLayout(astiav.ChannelLayoutStereo)
	audioFrame.SetNbSamples(960)
	audioFrame.SetSampleFormat(astiav.SampleFormatFlt)
	audioFrame.SetSampleRate(48000)

	// Alloc buffer
	align := 0
	if err := audioFrame.AllocBuffer(align); err != nil {
		log.Fatal(fmt.Errorf("main: allocating buffer failed: %w", err))
	}

	// Alloc samples
	if err := audioFrame.AllocSamples(align); err != nil {
		log.Fatal(fmt.Errorf("main: allocating image failed: %w", err))
	}

	// When writing data manually into a frame, you need to make sure the frame is writable
	if err := audioFrame.MakeWritable(); err != nil {
		log.Fatal(fmt.Errorf("main: making frame writable failed: %w", err))
	}

	// Let's say b1 contains an actual audio buffer, we can update the audio frame's data based on the buffer
	var b1 []byte
	if err := audioFrame.Data().SetBytes(b1, align); err != nil {
		log.Fatal(fmt.Errorf("main: setting frame's data based from bytes failed: %w", err))
	}

	// We can also retrieve the audio frame's data as buffer
	if _, err := audioFrame.Data().Bytes(align); err != nil {
		log.Fatal(fmt.Errorf("main: getting frame's data as bytes failed: %w", err))
	}

	/*

		In this second part we're going to manipulate a video frame

	*/

	// Alloc frame
	videoFrame := astiav.AllocFrame()
	defer videoFrame.Free()

	// To write data manually into a frame, proper attributes need to be set and allocated
	videoFrame.SetHeight(256)
	videoFrame.SetPixelFormat(astiav.PixelFormatRgba)
	videoFrame.SetWidth(256)

	// Alloc buffer
	align = 1
	if err := videoFrame.AllocBuffer(align); err != nil {
		log.Fatal(fmt.Errorf("main: allocating buffer failed: %w", err))
	}

	// Alloc image
	if err := videoFrame.AllocImage(align); err != nil {
		log.Fatal(fmt.Errorf("main: allocating image failed: %w", err))
	}

	// When writing data manually into a frame, you need to make sure the frame is writable
	if err := videoFrame.MakeWritable(); err != nil {
		log.Fatal(fmt.Errorf("main: making frame writable failed: %w", err))
	}

	// Let's say b2 contains an actual video buffer, we can update the video frame's data based on the buffer
	var b2 []byte
	if err := videoFrame.Data().SetBytes(b2, align); err != nil {
		log.Fatal(fmt.Errorf("main: setting frame's data based from bytes failed: %w", err))
	}

	// We can also retrieve the video frame's data as buffer
	if _, err := videoFrame.Data().Bytes(align); err != nil {
		log.Fatal(fmt.Errorf("main: getting frame's data as bytes failed: %w", err))
	}

	// Let's say i1 is an actual Go image.Image, we can update the video frame's data based on the image
	var i1 image.Image
	if err := videoFrame.Data().FromImage(i1); err != nil {
		log.Fatal(fmt.Errorf("main: setting frame's data based on Go image failed: %w", err))
	}

	// We can also retrieve the video frame's data as a Go image
	// For that we first need to guess the Go image format based on the frame's attributes before providing
	// it to .ToImage(). You may not need this and can provide your own image.Image to .ToImage()
	i2, err := videoFrame.Data().GuessImageFormat()
	if err != nil {
		log.Fatal(fmt.Errorf("main: guessing image format failed: %w", err))
	}
	if err := videoFrame.Data().ToImage(i2); err != nil {
		log.Fatal(fmt.Errorf("main: getting frame's data as Go image failed: %w", err))
	}

	// Success
	log.Println("success")
}
