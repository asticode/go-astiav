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
	filter = flag.String("f", "null", "the bit stream filter name")
	input  = flag.String("i", "", "the input path")
)

type stream struct {
	bitStreamFilterContext *astiav.BitStreamFilterContext
	bitStreamPkt           *astiav.Packet
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
		log.Println("Usage: <binary path> -i <input path> -f <bit stream filter name>")
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

	// Loop through streams
	streams := make(map[int]*stream) // Indexed by input stream index
	for _, is := range inputFormatContext.Streams() {
		// Only process audio or video
		if is.CodecParameters().MediaType() != astiav.MediaTypeAudio &&
			is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
			continue
		}

		// Create stream
		s := &stream{}

		// Alloc packet
		s.bitStreamPkt = astiav.AllocPacket()
		defer s.bitStreamPkt.Free()

		// Find bit stream filter
		bsf := astiav.FindBitStreamFilterByName(*filter)
		if bsf == nil {
			log.Fatal(errors.New("main: bit stream filter is nil"))
		}

		// Alloc bit stream filter context
		var err error
		if s.bitStreamFilterContext, err = astiav.AllocBitStreamFilterContext(bsf); err != nil {
			log.Fatal(fmt.Errorf("main: allocating bit stream filter context failed: %w", err))
		}
		defer s.bitStreamFilterContext.Free()

		// Copy codec parameters
		if err := is.CodecParameters().Copy(s.bitStreamFilterContext.InputCodecParameters()); err != nil {
			log.Fatal(fmt.Errorf("main: copying codec parameters failed: %w", err))
		}

		// Update time base
		s.bitStreamFilterContext.SetInputTimeBase(is.TimeBase())

		// Initialize bit stream filter context
		if err := s.bitStreamFilterContext.Initialize(); err != nil {
			log.Fatal(fmt.Errorf("main: initializing bit stream filter context failed: %w", err))
		}

		// Add stream
		streams[is.Index()] = s
	}

	// Loop through packets
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

		// Filter bit stream
		if err := filterBitStream(pkt, s); err != nil {
			log.Fatal(fmt.Errorf("main: filtering bit stream failed: %w", err))
		}
	}

	// Loop through streams
	for _, s := range streams {
		// Flush bit stream filter
		if err := filterBitStream(nil, s); err != nil {
			log.Fatal(fmt.Errorf("main: filtering bit stream failed: %w", err))
		}
	}

	// Success
	log.Println("success")
}

func filterBitStream(pkt *astiav.Packet, s *stream) error {
	// Send packet
	if err := s.bitStreamFilterContext.SendPacket(pkt); err != nil && !errors.Is(err, astiav.ErrEagain) {
		return fmt.Errorf("main: sending packet failed: %w", err)
	}

	// Loop
	for {
		// Receive packet
		if err := s.bitStreamFilterContext.ReceivePacket(s.bitStreamPkt); err != nil {
			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
				break
			}
			return fmt.Errorf("main: receiving packet failed: %w", err)
		}

		// Do something with packet
		log.Printf("new filtered packet: stream %d - pts: %d", s.bitStreamPkt.StreamIndex(), s.bitStreamPkt.Pts())

		// Unref packet
		s.bitStreamPkt.Unref()
	}
	return nil
}
