package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/asticode/go-astiav"
)

var (
	input  = flag.String("i", "", "the input path")
	output = flag.String("o", "", "the output path")
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
	if *input == "" || *output == "" {
		log.Println("Usage: <binary path> -i <input path> -o <output path>")
		return
	}

	// Allocate packet
	pkt := astiav.AllocPacket()
	defer pkt.Free()

	// Allocate input format context
	inputFormatContext := astiav.AllocFormatContext()
	if inputFormatContext == nil {
		log.Println(errors.New("main: input format context is nil"))
		return
	}
	defer inputFormatContext.Free()

	// Open input
	if err := inputFormatContext.OpenInput(*input, nil, nil); err != nil {
		log.Println(fmt.Errorf("main: opening input failed: %w", err))
		return
	}
	defer inputFormatContext.CloseInput()

	// Find stream info
	if err := inputFormatContext.FindStreamInfo(nil); err != nil {
		log.Println(fmt.Errorf("main: finding stream info failed: %w", err))
		return
	}

	// Allocate output format context
	outputFormatContext, err := astiav.AllocOutputFormatContext(nil, "", *output)
	if err != nil {
		log.Println(fmt.Errorf("main: allocating output format context failed: %w", err))
		return
	}
	if outputFormatContext == nil {
		log.Println(errors.New("main: output format context is nil"))
		return
	}
	defer outputFormatContext.Free()

	// Loop through streams
	inputStreams := make(map[int]*astiav.Stream)  // Indexed by input stream index
	outputStreams := make(map[int]*astiav.Stream) // Indexed by input stream index
	for _, is := range inputFormatContext.Streams() {
		// Only process audio or video
		if is.CodecParameters().MediaType() != astiav.MediaTypeAudio &&
			is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
			continue
		}

		// Add input stream
		inputStreams[is.Index()] = is

		// Add stream to output format context
		os := outputFormatContext.NewStream(nil)
		if os == nil {
			log.Println(errors.New("main: output stream is nil"))
			return
		}

		// Copy codec parameters
		if err = is.CodecParameters().Copy(os.CodecParameters()); err != nil {
			log.Println(fmt.Errorf("main: copying codec parameters failed: %w", err))
			return
		}

		// Reset codec tag
		os.CodecParameters().SetCodecTag(0)

		// Add output stream
		outputStreams[is.Index()] = os
	}

	// If this is a file, we need to use an io context
	if !outputFormatContext.OutputFormat().Flags().Has(astiav.IOFormatFlagNofile) {
		// Open io context
		ioContext, err := astiav.OpenIOContext(*output, astiav.NewIOContextFlags(astiav.IOContextFlagWrite), nil, nil)
		if err != nil {
			log.Println(fmt.Errorf("main: opening io context failed: %w", err))
			return
		}
		defer ioContext.Close() //nolint:errcheck

		// Update output format context
		outputFormatContext.SetPb(ioContext)
	}

	// Write header
	if err = outputFormatContext.WriteHeader(nil); err != nil {
		log.Println(fmt.Errorf("main: writing header failed: %w", err))
		return
	}

	// Loop through packets
	for {
		// We use a closure to ease unreferencing packet
		if stop := func() bool {
			// Read frame
			if err = inputFormatContext.ReadFrame(pkt); err != nil {
				if !errors.Is(err, astiav.ErrEof) {
					log.Println(fmt.Errorf("main: reading frame failed: %w", err))
				}
				return true
			}

			// Make sure to unreference packet
			defer pkt.Unref()

			// Get input stream
			inputStream, ok := inputStreams[pkt.StreamIndex()]
			if !ok {
				return false
			}

			// Get output stream
			outputStream, ok := outputStreams[pkt.StreamIndex()]
			if !ok {
				return false
			}

			// Update packet
			pkt.SetStreamIndex(outputStream.Index())
			pkt.RescaleTs(inputStream.TimeBase(), outputStream.TimeBase())
			pkt.SetPos(-1)

			// Write frame
			if err = outputFormatContext.WriteInterleavedFrame(pkt); err != nil {
				log.Println(fmt.Errorf("main: writing interleaved frame failed: %w", err))
				return true
			}
			return false
		}(); stop {
			break
		}
	}

	// Write trailer
	if err = outputFormatContext.WriteTrailer(); err != nil {
		log.Println(fmt.Errorf("main: writing trailer failed: %w", err))
		return
	}

	// Done
	log.Println("done")
}
