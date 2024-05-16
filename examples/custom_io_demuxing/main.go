package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/asticode/go-astiav"
)

var (
	input = flag.String("i", "", "the input path")
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
	if *input == "" {
		log.Println("Usage: <binary path> -i <input path>")
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

	// Open file
	f, err := os.Open(*input)
	if err != nil {
		log.Fatal(fmt.Errorf("main: opening %s failed: %w", *input, err))
	}
	defer f.Close()

	// Alloc io context
	ioContext, err := astiav.AllocIOContext(
		4096,
		false,
		func(b []byte) (n int, err error) {
			return f.Read(b)
		},
		func(offset int64, whence int) (n int64, err error) {
			return f.Seek(offset, whence)
		},
		nil,
	)
	if err != nil {
		log.Fatal(fmt.Errorf("main: allocating io context failed: %w", err))
	}
	defer ioContext.Free()

	// Store io context
	inputFormatContext.SetPb(ioContext)

	// Open input
	if err := inputFormatContext.OpenInput("", nil, nil); err != nil {
		log.Fatal(fmt.Errorf("main: opening input failed: %w", err))
	}
	defer inputFormatContext.CloseInput()

	// Find stream info
	if err := inputFormatContext.FindStreamInfo(nil); err != nil {
		log.Fatal(fmt.Errorf("main: finding stream info failed: %w", err))
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

		// Do something with the packet
		log.Printf("new packet: stream %d - pts: %d", pkt.StreamIndex(), pkt.Pts())
	}

	// Success
	log.Println("success")
}
