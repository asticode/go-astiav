package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"log"
	"strings"

	"github.com/asticode/go-astiav"
)

var (
	input = flag.String("i", "", "the input path")
)

type stream struct {
	decCodec        *astiav.Codec
	decCodecContext *astiav.CodecContext
	inputStream     *astiav.Stream
}

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
	if *input == "" {
		log.Println("Usage: <binary path> -i <input path>")
		return
	}

	// Alloc packet
	pkt := astiav.AllocPacket()
	defer pkt.Free()

	// Alloc frame
	f := astiav.AllocFrame()
	defer f.Free()

	// Alloc input format context
	inputFormatContext := astiav.AllocFormatContext()
	if inputFormatContext == nil {
		log.Fatal(errors.New("main: input format context is nil"))
	}
	defer inputFormatContext.Free()

	// Open input
	if err := inputFormatContext.OpenInput(*input, nil, nil); err != nil {
		log.Fatal(fmt.Errorf("main: opening input failed: %w", err))
	}
	defer inputFormatContext.CloseInput()

	// Find stream info
	if err := inputFormatContext.FindStreamInfo(nil); err != nil {
		log.Fatal(fmt.Errorf("main: finding stream info failed: %w", err))
	}

	// Loop through streams
	streams := make(map[int]*stream) // Indexed by input stream index
	for _, is := range inputFormatContext.Streams() {
		// Only process audio or video
		if is.CodecParameters().MediaType() != astiav.MediaTypeAudio &&
			is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
			continue
		}

		// Create stream
		s := &stream{inputStream: is}

		// Find decoder
		if s.decCodec = astiav.FindDecoder(is.CodecParameters().CodecID()); s.decCodec == nil {
			log.Fatal(errors.New("main: codec is nil"))
		}

		// Alloc codec context
		if s.decCodecContext = astiav.AllocCodecContext(s.decCodec); s.decCodecContext == nil {
			log.Fatal(errors.New("main: codec context is nil"))
		}
		defer s.decCodecContext.Free()

		// Update codec context
		if err := is.CodecParameters().ToCodecContext(s.decCodecContext); err != nil {
			log.Fatal(fmt.Errorf("main: updating codec context failed: %w", err))
		}

		// Open codec context
		if err := s.decCodecContext.Open(s.decCodec, nil); err != nil {
			log.Fatal(fmt.Errorf("main: opening codec context failed: %w", err))
		}

		// Add stream
		streams[is.Index()] = s
	}

	// Loop through packets
	var i image.Image
	for {
		// Read frame
		if err := inputFormatContext.ReadFrame(pkt); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			log.Fatal(fmt.Errorf("main: reading frame failed: %w", err))
		}

		// Get stream
		s, ok := streams[pkt.StreamIndex()]
		if !ok {
			continue
		}

		// Send packet
		if err := s.decCodecContext.SendPacket(pkt); err != nil {
			log.Fatal(fmt.Errorf("main: sending packet failed: %w", err))
		}

		// Loop
		for {
			// Receive frame
			if err := s.decCodecContext.ReceiveFrame(f); err != nil {
				if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
					break
				}
				log.Fatal(fmt.Errorf("main: receiving frame failed: %w", err))
			}

			// Do something with decoded frame
			if s.inputStream.CodecParameters().MediaType() == astiav.MediaTypeVideo {
				// In this example, we'll process the frame data but you can do whatever you feel like
				// with the decoded frame
				fd := f.Data()

				// Image has not yet been initialized
				// If the image format can change in the stream, you'll need to guess image format for every frame
				if i == nil {
					// Guess image format
					// It might not return an image.Image in the proper format for your use case, in that case
					// you can skip this step and provide .ToImage() with your own image.Image
					var err error
					if i, err = fd.GuessImageFormat(); err != nil {
						log.Fatal(fmt.Errorf("main: guessing image format failed: %w", err))
					}
				}

				// Copy frame data to the image
				if err := fd.ToImage(i); err != nil {
					log.Fatal(fmt.Errorf("main: copying frame data to the image failed: %w", err))
				}

				// Log
				log.Printf("new video frame: stream %d - pts: %d - size: %dx%d - color at (0,0): %+v", pkt.StreamIndex(), f.Pts(), i.Bounds().Dx(), i.Bounds().Dy(), i.At(0, 0))
			} else {
				// Log
				log.Printf("new audio frame: stream %d - pts: %d", pkt.StreamIndex(), f.Pts())
			}
		}
	}

	// Success
	log.Println("success")
}
