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
	encoderCodecName        = flag.String("c", "", "the encoder codec name (e.g. h264_nvenc)")
	hardwareDeviceName      = flag.String("n", "", "the hardware device name (e.g. 0)")
	hardwareDeviceTypeName  = flag.String("t", "", "the hardware device type (e.g. cuda)")
	hardwarePixelFormatName = flag.String("hpf", "", "the hardware pixel format name (e.g. cuda)")

	width           = flag.Int("w", 1920, "the width")
	height          = flag.Int("h", 1080, "the height")
	fps             = flag.Int("f", 25, "the fps")
	initialPoolSize = flag.Int("p", 20, "the initial pool size")
	patternGridSize = flag.Int("g", 128, "the pattern grid size")

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
	if *hardwareDeviceTypeName == "" || *encoderCodecName == "" || *hardwarePixelFormatName == "" || *output == "" {
		log.Println("Usage: <binary path> -t <hardware device type> -c <encoder codec> -hpf <hardware pixel format> -o <output path> [-n <hardware device name> -w <width> -h <height> -f <fps> -p <initial pool size> -g <pattern grid size>]")
		return
	}

	// Open output file
	output, err := os.OpenFile(*output, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal(fmt.Errorf("main: opening output file failed: %w", err))
	}
	defer output.Close()

	// Get hardware device type
	hardwareDeviceType := astiav.FindHardwareDeviceTypeByName(*hardwareDeviceTypeName)
	if hardwareDeviceType == astiav.HardwareDeviceTypeNone {
		log.Fatal(errors.New("main: hardware device not found"))
	}

	// Create hardware device context
	hardwareDeviceContext, err := astiav.CreateHardwareDeviceContext(hardwareDeviceType, *hardwareDeviceName, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("main: creating hardware device context failed: %w", err))
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

	// Set codec context
	encCodecContext.SetWidth(*width)
	encCodecContext.SetHeight(*height)
	encCodecContext.SetTimeBase(astiav.NewRational(1, *fps))
	encCodecContext.SetFramerate(astiav.NewRational(*fps, 1))
	hardwarePixelFormatName := astiav.FindPixelFormatByName(*hardwarePixelFormatName)
	if hardwarePixelFormatName == astiav.PixelFormatNone {
		log.Fatal("main: hardware pixel format not found")
	}
	encCodecContext.SetPixelFormat(hardwarePixelFormatName)

	// Set hardware frame context
	hardwareFrameCtx := astiav.AllocHardwareFrameContext(hardwareDeviceContext)
	if hardwareFrameCtx == nil {
		log.Fatal("main: hardware frame context is nil")
	}
	hardwareFrameCtx.SetPixelFormat(hardwarePixelFormatName)
	hardwareFrameCtx.SetSoftwarePixelFormat(astiav.PixelFormatNv12)
	hardwareFrameCtx.SetWidth(*width)
	hardwareFrameCtx.SetHeight(*height)
	hardwareFrameCtx.SetInitialPoolSize(*initialPoolSize)

	// Initialize hardware frame context
	if err := hardwareFrameCtx.Initialize(); err != nil {
		log.Fatal(fmt.Errorf("main: initializing hardware frame context failed: %w", err))
	}
	encCodecContext.SetHardwareFrameContext(hardwareFrameCtx)

	// Open codec context
	if err := encCodecContext.Open(encCodec, nil); err != nil {
		log.Fatal(fmt.Errorf("main: opening codec context failed: %w", err))
	}

	frameIndex := 0

	// Draw frames, upload them to hardware devices, and encode them
	for {
		// Alloc software frame
		softwareFrame := astiav.AllocFrame()

		// Set software frame
		softwareFrame.SetWidth(*width)
		softwareFrame.SetHeight(*height)
		softwareFrame.SetPixelFormat(astiav.PixelFormatNv12)

		// Alloc software frame buffer
		if err := softwareFrame.AllocBuffer(0); err != nil {
			log.Fatal(fmt.Errorf("main: allocating buffer failed: %w", err))
		}

		// Fill software frame
		yPlane, uvPlane := MakeNV12MovingCheckerboardPattern(*width, *height, *patternGridSize, frameIndex)
		softwareFrame.SetData(0, yPlane)
		softwareFrame.SetData(1, uvPlane)

		// Alloc hardware frame
		hardwareFrame := astiav.AllocFrame()

		// Alloc hardware frame buffer
		if err := hardwareFrame.AllocHardwareBuffer(hardwareFrameCtx); err != nil {
			log.Fatal(fmt.Errorf("main: allocating hardware buffer failed: %w", err))
		}

		// Upload software frame to hardware frame
		if err := softwareFrame.TransferHardwareData(hardwareFrame); err != nil {
			log.Fatal(fmt.Errorf("main: uploading from frame failed: %w", err))
		}
		softwareFrame.Free()

		// Encode frame
		if err := encCodecContext.SendFrame(hardwareFrame); err != nil {
			log.Fatal(fmt.Errorf("main: sending frame failed: %w", err))
		}
		hardwareFrame.Free()

		// Receive packet
		packet := astiav.AllocPacket()
		for {
			if err := encCodecContext.ReceivePacket(packet); err != nil {
				break
			}

			// Write packet
			if _, err := output.Write(packet.Data()); err != nil {
				log.Fatal(fmt.Errorf("main: writing packet failed: %w", err))
			}

			packet.Unref()
		}
		packet.Free()

		frameIndex++
		log.Printf("Finished encoding frame %d\n", frameIndex)
	}
}

func MakeNV12MovingCheckerboardPattern(width, height, blockSize, frame int) ([]byte, []byte) {
	yPlane := make([]byte, width*height)
	uvPlane := make([]byte, width*height/2)

	xOffset := frame % blockSize
	yOffset := frame % blockSize

	// Y plane (checkerboard pattern)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if ((x+xOffset)/blockSize+(y+yOffset)/blockSize)%2 == 0 {
				yPlane[y*width+x] = 255 // White
			} else {
				yPlane[y*width+x] = 0 // Black
			}
		}
	}

	// UV plane (gray)
	for i := 0; i < len(uvPlane); i += 2 {
		uvPlane[i] = 128   // U component (neutral)
		uvPlane[i+1] = 128 // V component (neutral)
	}

	return yPlane, uvPlane
}
