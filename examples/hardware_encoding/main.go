package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/asticode/go-astiav"
)

var (
	encoderCodecName        = flag.String("c", "", "the encoder codec name (e.g. h264_nvenc)")
	hardwareDeviceName      = flag.String("n", "", "the hardware device name (e.g. 0)")
	hardwareDeviceTypeName  = flag.String("t", "", "the hardware device type (e.g. cuda)")
	hardwarePixelFormatName = flag.String("hpf", "", "the hardware pixel format name (e.g. cuda)")
	height                  = flag.Int("h", 1080, "the height")
	softwarePixelFormatName = flag.String("spf", "", "the software pixel format name (e.g. nv12)")
	width                   = flag.Int("w", 1920, "the width")
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
	if *hardwareDeviceTypeName == "" || *encoderCodecName == "" || *hardwarePixelFormatName == "" {
		log.Println("Usage: <binary path> -t <hardware device type> -c <encoder codec> -hpf <hardware pixel format> [-n <hardware device name> -w <width> -h <height>]")
		return
	}

	// Get hardware device type
	hardwareDeviceType := astiav.FindHardwareDeviceTypeByName(*hardwareDeviceTypeName)
	if hardwareDeviceType == astiav.HardwareDeviceTypeNone {
		log.Fatal(errors.New("main: hardware device not found"))
	}

	// Create hardware device context
	hardwareDeviceContext, err := astiav.CreateHardwareDeviceContext(hardwareDeviceType, *hardwareDeviceName, nil, 0)
	if err != nil {
		log.Fatal(fmt.Errorf("main: creating hardware device context failed: %w", err))
	}
	defer hardwareDeviceContext.Free()

	hardwareFramesConstraints := hardwareDeviceContext.HardwareFramesConstraints()
	if hardwareFramesConstraints == nil {
		log.Fatal("main: hardware frames constraints is nil")
		return
	}
	defer hardwareFramesConstraints.Free()

	validHardwarePixelFormats := hardwareFramesConstraints.ValidHardwarePixelFormats()
	if len(validHardwarePixelFormats) == 0 {
		log.Fatal("main: no valid hardware pixel formats")
		return
	}

	// Find encoder codec
	encCodec := astiav.FindEncoderByName(*encoderCodecName)
	if encCodec == nil {
		log.Fatal("main: encoder codec is nil")
	}

	// Alloc codec context
	encCodecContext := astiav.AllocCodecContext(encCodec)
	if encCodecContext == nil {
		log.Fatal("main: codec context is nil")
	}
	defer encCodecContext.Free()

	// Get hardware pixel format
	hardwarePixelFormat := astiav.FindPixelFormatByName(*hardwarePixelFormatName)
	if hardwarePixelFormat == astiav.PixelFormatNone {
		log.Fatal("main: hardware pixel format not found")
	} else if !slices.Contains(validHardwarePixelFormats, hardwarePixelFormat) {
		log.Fatalf("main: hardware pixel format not supported : %s", hardwarePixelFormat)
	}

	// Set codec context
	encCodecContext.SetWidth(*width)
	encCodecContext.SetHeight(*height)
	encCodecContext.SetTimeBase(astiav.NewRational(1, 25))
	encCodecContext.SetFramerate(encCodecContext.TimeBase().Invert())
	encCodecContext.SetPixelFormat(hardwarePixelFormat)

	// Alloc hardware frames context
	hardwareFramesContext := astiav.AllocHardwareFramesContext(hardwareDeviceContext)
	if hardwareFramesContext == nil {
		log.Fatal("main: hardware frames context is nil")
	}
	defer hardwareFramesContext.Free()

	validSoftwarePixelFormats := hardwareFramesConstraints.ValidSoftwarePixelFormats()
	if len(validSoftwarePixelFormats) == 0 {
		log.Fatal("main: no valid software pixel formats")
		return
	}

	// Get software pixel format
	softwarePixelFormat := astiav.FindPixelFormatByName(*softwarePixelFormatName)
	if softwarePixelFormat == astiav.PixelFormatNone {
		log.Fatal("main: software pixel format not found")
	} else if !slices.Contains(validSoftwarePixelFormats, softwarePixelFormat) {
		log.Fatalf("main: software pixel format not supported : %s", softwarePixelFormat)
	}

	// Set hardware frame content
	hardwareFramesContext.SetHardwarePixelFormat(hardwarePixelFormat)
	hardwareFramesContext.SetSoftwarePixelFormat(softwarePixelFormat)
	hardwareFramesContext.SetWidth(*width)
	hardwareFramesContext.SetHeight(*height)
	hardwareFramesContext.SetInitialPoolSize(20)

	// Initialize hardware frame context
	if err := hardwareFramesContext.Initialize(); err != nil {
		log.Fatal(fmt.Errorf("main: initializing hardware frame context failed: %w", err))
	}

	// Update encoder codec context hardware frame context
	encCodecContext.SetHardwareFramesContext(hardwareFramesContext)

	// Open codec context
	if err := encCodecContext.Open(encCodec, nil); err != nil {
		log.Fatal(fmt.Errorf("main: opening codec context failed: %w", err))
	}

	// Alloc software frame
	softwareFrame := astiav.AllocFrame()
	defer softwareFrame.Free()

	// Set software frame
	softwareFrame.SetWidth(*width)
	softwareFrame.SetHeight(*height)
	softwareFrame.SetPixelFormat(softwarePixelFormat)

	// Alloc software frame buffer
	if err := softwareFrame.AllocBuffer(0); err != nil {
		log.Fatal(fmt.Errorf("main: allocating buffer failed: %w", err))
	}

	// Fill software frame with black
	if err = softwareFrame.ImageFillBlack(); err != nil {
		log.Fatal(fmt.Errorf("main: filling software frame with black failed: %w", err))
	}

	// Alloc hardware frame
	hardwareFrame := astiav.AllocFrame()
	defer hardwareFrame.Free()

	// Alloc hardware frame buffer
	if err := hardwareFrame.AllocHardwareBuffer(hardwareFramesContext); err != nil {
		log.Fatal(fmt.Errorf("main: allocating hardware buffer failed: %w", err))
	}

	// Transfer software frame data to hardware frame
	if err := softwareFrame.TransferHardwareData(hardwareFrame); err != nil {
		log.Fatal(fmt.Errorf("main: transferring hardware data failed: %w", err))
	}

	// Encode frame
	if err := encCodecContext.SendFrame(hardwareFrame); err != nil {
		log.Fatal(fmt.Errorf("main: sending frame failed: %w", err))
	}

	// Alloc packet
	pkt := astiav.AllocPacket()
	defer pkt.Free()

	// Loop
	for {
		// We use a closure to ease unreferencing the packet
		if stop := func() bool {
			// Receive packet
			if err = encCodecContext.ReceivePacket(pkt); err != nil {
				if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
					return true
				}
				log.Fatal(fmt.Errorf("main: receiving packet failed: %w", err))
			}

			// Make sure to unreference packet
			defer pkt.Unref()

			// Log
			log.Println("new packet")
			return false
		}(); stop {
			break
		}
	}
}
