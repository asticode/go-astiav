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

	// Open file
	f, err := os.Open(*input)
	if err != nil {
		log.Println(fmt.Errorf("main: opening %s failed: %w", *input, err))
		return
	}
	defer f.Close()

	// Allocate io context
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
		log.Println(fmt.Errorf("main: allocating io context failed: %w", err))
		return
	}
	defer ioContext.Free()

	// Store io context
	inputFormatContext.SetPb(ioContext)

	// Open input
	if err := inputFormatContext.OpenInput("", nil, nil); err != nil {
		log.Println(fmt.Errorf("main: opening input failed: %w", err))
		return
	}
	defer inputFormatContext.CloseInput()

	// Find stream info
	if err := inputFormatContext.FindStreamInfo(nil); err != nil {
		log.Println(fmt.Errorf("main: finding stream info failed: %w", err))
		return
	}

	// Loop through packets
	for {
		// We use a closure to ease unreferencing the packet
		if stop := func() bool {
			// Read frame
			if err := inputFormatContext.ReadFrame(pkt); err != nil {
				if !errors.Is(err, astiav.ErrEof) {
					log.Println(fmt.Errorf("main: reading frame failed: %w", err))
				}
				return true
			}

			// Make sure to unreference the packet
			defer pkt.Unref()

			// Do something with the packet
			log.Printf("new packet: stream %d - pts: %d", pkt.StreamIndex(), pkt.Pts())
			return false
		}(); stop {
			break
		}
	}

	// Done
	log.Println("done")
}
