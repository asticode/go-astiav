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
	decoderCodecName       = flag.String("c", "", "the decoder codec name (e.g. h264_cuvid)")
	hardwareDeviceName     = flag.String("n", "", "the hardware device name (e.g. 0)")
	hardwareDeviceTypeName = flag.String("t", "", "the hardware device type (e.g. cuda)")
	input                  = flag.String("i", "", "the input path")
)

type stream struct {
	decCodec              *astiav.Codec
	decCodecContext       *astiav.CodecContext
	hardwareDeviceContext *astiav.HardwareDeviceContext
	hardwarePixelFormat   astiav.PixelFormat
	inputStream           *astiav.Stream
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
	if *input == "" || *hardwareDeviceTypeName == "" {
		log.Println("Usage: <binary path> -t <hardware device type> -i <input path> [-n <hardware device name> -c <decoder codec>]")
		return
	}

	// Get hardware device type
	hardwareDeviceType := astiav.FindHardwareDeviceTypeByName(*hardwareDeviceTypeName)
	if hardwareDeviceType == astiav.HardwareDeviceTypeNone {
		log.Fatal(errors.New("main: hardware device not found"))
	}

	// Alloc packet
	pkt := astiav.AllocPacket()
	defer pkt.Free()

	// Alloc hardware frame
	hardwareFrame := astiav.AllocFrame()
	defer hardwareFrame.Free()

	// Alloc software frame
	softwareFrame := astiav.AllocFrame()
	defer softwareFrame.Free()

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
		// Only process video
		if is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
			continue
		}

		// Create stream
		s := &stream{inputStream: is}

		// Find decoder
		if *decoderCodecName != "" {
			s.decCodec = astiav.FindDecoderByName(*decoderCodecName)
		} else {
			s.decCodec = astiav.FindDecoder(is.CodecParameters().CodecID())
		}

		// No codec
		if s.decCodec == nil {
			log.Fatal(errors.New("main: codec is nil"))
		}

		// Alloc codec context
		if s.decCodecContext = astiav.AllocCodecContext(s.decCodec); s.decCodecContext == nil {
			log.Fatal(errors.New("main: codec context is nil"))
		}
		defer s.decCodecContext.Free()

		// Loop through codec hardware configs
		for _, p := range s.decCodec.HardwareConfigs() {
			// Valid hardware config
			if p.MethodFlags().Has(astiav.CodecHardwareConfigMethodFlagHwDeviceCtx) && p.HardwareDeviceType() == hardwareDeviceType {
				s.hardwarePixelFormat = p.PixelFormat()
				break
			}
		}

		// No valid hardware pixel format
		if s.hardwarePixelFormat == astiav.PixelFormatNone {
			log.Fatal(errors.New("main: hardware device type not supported by decoder"))
		}

		// Update codec context
		if err := is.CodecParameters().ToCodecContext(s.decCodecContext); err != nil {
			log.Fatal(fmt.Errorf("main: updating codec context failed: %w", err))
		}

		// Create hardware device context
		var err error
		if s.hardwareDeviceContext, err = astiav.CreateHardwareDeviceContext(hardwareDeviceType, *hardwareDeviceName, nil); err != nil {
			log.Fatal(fmt.Errorf("main: creating hardware device context failed: %w", err))
		}

		// Update decoder context
		s.decCodecContext.SetHardwareDeviceContext(s.hardwareDeviceContext)
		s.decCodecContext.SetPixelFormatCallback(func(pfs []astiav.PixelFormat) astiav.PixelFormat {
			for _, pf := range pfs {
				if pf == s.hardwarePixelFormat {
					return pf
				}
			}
			log.Fatal(errors.New("main: using hardware pixel format failed"))
			return astiav.PixelFormatNone
		})

		// Open codec context
		if err := s.decCodecContext.Open(s.decCodec, nil); err != nil {
			log.Fatal(fmt.Errorf("main: opening codec context failed: %w", err))
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

		// Send packet
		if err := s.decCodecContext.SendPacket(pkt); err != nil {
			log.Fatal(fmt.Errorf("main: sending packet failed: %w", err))
		}

		// Loop
		for {
			// Receive frame
			if err := s.decCodecContext.ReceiveFrame(hardwareFrame); err != nil {
				if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
					break
				}
				log.Fatal(fmt.Errorf("main: receiving frame failed: %w", err))
			}

			// Get final frame
			var finalFrame *astiav.Frame
			if hardwareFrame.PixelFormat() == s.hardwarePixelFormat {
				// Transfer hardware data
				if err := hardwareFrame.TransferHardwareData(softwareFrame); err != nil {
					log.Fatal(fmt.Errorf("main: transferring hardware data failed: %w", err))
				}

				// Update pts
				softwareFrame.SetPts(hardwareFrame.Pts())

				// Update final frame
				finalFrame = softwareFrame
			} else {
				// Update final frame
				finalFrame = hardwareFrame
			}

			// Do something with decoded frame
			log.Printf("new frame: stream %d - pts: %d - transferred: %v", pkt.StreamIndex(), finalFrame.Pts(), hardwareFrame.PixelFormat() == s.hardwarePixelFormat)
		}
	}

	// Success
	log.Println("success")
}
