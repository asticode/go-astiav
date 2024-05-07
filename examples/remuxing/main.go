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

	// Alloc packet
	pkt := astiav.AllocPacket()
	defer pkt.Free()

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

	// Alloc output format context
	outputFormatContext, err := astiav.AllocOutputFormatContext(nil, "", *output)
	if err != nil {
		log.Fatal(fmt.Errorf("main: allocating output format context failed: %w", err))
	}
	if outputFormatContext == nil {
		log.Fatal(errors.New("main: output format context is nil"))
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
			log.Fatal(errors.New("main: output stream is nil"))
		}

		// Copy codec parameters
		if err = is.CodecParameters().Copy(os.CodecParameters()); err != nil {
			log.Fatal(fmt.Errorf("main: copying codec parameters failed: %w", err))
		}

		// Reset codec tag
		os.CodecParameters().SetCodecTag(0)

		// Add output stream
		outputStreams[is.Index()] = os
	}

	// If this is a file, we need to use an io context
	if !outputFormatContext.OutputFormat().Flags().Has(astiav.IOFormatFlagNofile) {
		// Open io context
		ioContext, err := astiav.OpenIOContext(*output, astiav.NewIOContextFlags(astiav.IOContextFlagWrite))
		if err != nil {
			log.Fatal(fmt.Errorf("main: opening io context failed: %w", err))
		}
		defer ioContext.Close() //nolint:errcheck

		// Update output format context
		outputFormatContext.SetPb(ioContext)
	}

	// Write header
	if err = outputFormatContext.WriteHeader(nil); err != nil {
		log.Fatal(fmt.Errorf("main: writing header failed: %w", err))
	}

	// Loop through packets
	for {
		// Read frame
		if err = inputFormatContext.ReadFrame(pkt); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			log.Fatal(fmt.Errorf("main: reading frame failed: %w", err))
		}

		// Get input stream
		inputStream, ok := inputStreams[pkt.StreamIndex()]
		if !ok {
			pkt.Unref()
			continue
		}

		// Get output stream
		outputStream, ok := outputStreams[pkt.StreamIndex()]
		if !ok {
			pkt.Unref()
			continue
		}

		// Update packet
		pkt.SetStreamIndex(outputStream.Index())
		pkt.RescaleTs(inputStream.TimeBase(), outputStream.TimeBase())
		pkt.SetPos(-1)

		// Write frame
		if err = outputFormatContext.WriteInterleavedFrame(pkt); err != nil {
			log.Fatal(fmt.Errorf("main: writing interleaved frame failed: %w", err))
		}
	}

	// Write trailer
	if err = outputFormatContext.WriteTrailer(); err != nil {
		log.Fatal(fmt.Errorf("main: writing trailer failed: %w", err))
	}

	// Success
	log.Println("success")
}
