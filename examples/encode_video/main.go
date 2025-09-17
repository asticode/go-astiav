package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/asticode/go-astiav"
)

var (
	output = flag.String("o", "", "the output video file path")
	codec  = flag.String("c", "libx264", "the video codec to use (libx264, mpeg4, etc.)")
	width  = flag.Int("w", 320, "video width")
	height = flag.Int("h", 240, "video height")
)

func main() {
	// Handle ffmpeg logs
	astiav.SetLogLevel(astiav.LogLevelInfo)
	astiav.SetLogCallback(func(c astiav.Classer, l astiav.LogLevel, fmtStr, msg string) {
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
	if *output == "" {
		log.Println("Usage: <binary path> -o <output video file> [-c <codec>] [-w <width>] [-h <height>]")
		return
	}

	// Find encoder
	encoder := astiav.FindEncoderByName(*codec)
	if encoder == nil {
		log.Fatal(fmt.Errorf("failed to find encoder for codec: %s", *codec))
	}

	// Allocate codec context
	codecContext := astiav.AllocCodecContext(encoder)
	if codecContext == nil {
		log.Fatal("failed to allocate codec context")
	}
	defer codecContext.Free()

	// Set codec parameters
	codecContext.SetBitRate(400000)
	codecContext.SetWidth(*width)
	codecContext.SetHeight(*height)
	codecContext.SetTimeBase(astiav.NewRational(1, 25)) // 25 FPS
	codecContext.SetFramerate(astiav.NewRational(25, 1))
	codecContext.SetPixelFormat(astiav.PixelFormatYuv420P)

	// Set GOP size (Group of Pictures)
	codecContext.SetGopSize(10)
	codecContext.SetMaxBFrames(1)

	// Check if the encoder supports the pixel format
	if !checkPixelFormat(encoder, codecContext.PixelFormat()) {
		log.Printf("Pixel format %s not supported by encoder, trying alternatives...", codecContext.PixelFormat().Name())

		// Try common formats
		supportedFormats := []astiav.PixelFormat{
			astiav.PixelFormatYuv420P,
			astiav.PixelFormatYuv422P,
			astiav.PixelFormatYuv444P,
			astiav.PixelFormatRgb24,
		}

		found := false
		for _, pf := range supportedFormats {
			if checkPixelFormat(encoder, pf) {
				codecContext.SetPixelFormat(pf)
				log.Printf("Using pixel format: %s", pf.Name())
				found = true
				break
			}
		}

		if !found {
			log.Fatal("No supported pixel format found")
		}
	}

	// Open codec
	if err := codecContext.Open(encoder, nil); err != nil {
		log.Fatal(fmt.Errorf("failed to open codec: %w", err))
	}

	// Create output file
	outputFile, err := os.Create(*output)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create output file: %w", err))
	}
	defer outputFile.Close()

	// Allocate frame and packet
	frame := astiav.AllocFrame()
	defer frame.Free()

	packet := astiav.AllocPacket()
	defer packet.Free()

	// Set frame parameters
	frame.SetWidth(*width)
	frame.SetHeight(*height)
	frame.SetPixelFormat(codecContext.PixelFormat())

	// Allocate frame buffer
	if err := frame.AllocBuffer(32); err != nil {
		log.Fatal(fmt.Errorf("failed to allocate frame buffer: %w", err))
	}

	// Print encoding information
	log.Printf("Encoding video:")
	log.Printf("  Codec: %s", encoder.Name())
	log.Printf("  Resolution: %dx%d", *width, *height)
	log.Printf("  Pixel format: %s", codecContext.PixelFormat().Name())
	log.Printf("  Frame rate: %v", codecContext.Framerate())
	log.Printf("  Bit rate: %d bps", codecContext.BitRate())

	// Generate and encode video frames
	totalFrames := 125 // 5 seconds at 25 FPS

	for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
		// Make frame writable
		if err := frame.MakeWritable(); err != nil {
			log.Fatal(fmt.Errorf("failed to make frame writable: %w", err))
		}

		// Generate video data (animated pattern)
		if err := generateVideoData(frame, frameIndex); err != nil {
			log.Fatal(fmt.Errorf("failed to generate video data: %w", err))
		}

		// Set frame timestamp
		frame.SetPts(int64(frameIndex))

		// Encode frame
		if err := encodeFrame(codecContext, frame, packet, outputFile); err != nil {
			log.Fatal(fmt.Errorf("failed to encode frame: %w", err))
		}

		if frameIndex%25 == 0 {
			log.Printf("Encoded %d/%d frames", frameIndex, totalFrames)
		}
	}

	// Flush encoder
	if err := encodeFrame(codecContext, nil, packet, outputFile); err != nil {
		log.Fatal(fmt.Errorf("failed to flush encoder: %w", err))
	}

	log.Println("Video encoding completed successfully")
}

func checkPixelFormat(codec *astiav.Codec, pixelFormat astiav.PixelFormat) bool {
	formats := codec.PixelFormats()
	for _, pf := range formats {
		if pf == pixelFormat {
			return true
		}
	}
	return false
}

func generateVideoData(frame *astiav.Frame, frameIndex int) error {
	width := frame.Width()
	height := frame.Height()
	pixelFormat := frame.PixelFormat()

	switch pixelFormat {
	case astiav.PixelFormatYuv420P:
		return generateYUV420P(frame, width, height, frameIndex)
	case astiav.PixelFormatYuv422P:
		return generateYUV422P(frame, width, height, frameIndex)
	case astiav.PixelFormatYuv444P:
		return generateYUV444P(frame, width, height, frameIndex)
	case astiav.PixelFormatRgb24:
		return generateRGB24(frame, width, height, frameIndex)
	default:
		return fmt.Errorf("unsupported pixel format for generation: %s", pixelFormat.Name())
	}
}

func generateYUV420P(frame *astiav.Frame, width, height, frameIndex int) error {
	// Generate Y plane (luminance)
	yData := frame.DataSlice(0, width*height)
	if yData != nil {
		linesize := frame.Linesize()[0]
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Create animated pattern
				value := int((math.Sin(float64(x+frameIndex)*0.1)+math.Sin(float64(y+frameIndex)*0.1))*127 + 128)
				if value < 0 {
					value = 0
				}
				if value > 255 {
					value = 255
				}

				offset := y*linesize + x
				if offset < len(yData) {
					yData[offset] = byte(value)
				}
			}
		}
	}

	// Generate U plane (chroma)
	uData := frame.DataSlice(1, (width/2)*(height/2))
	if uData != nil {
		linesize := frame.Linesize()[1]
		for y := 0; y < height/2; y++ {
			for x := 0; x < width/2; x++ {
				// Blue-ish chroma
				value := 128 + int(math.Sin(float64(frameIndex)*0.1)*50)
				if value < 0 {
					value = 0
				}
				if value > 255 {
					value = 255
				}

				offset := y*linesize + x
				if offset < len(uData) {
					uData[offset] = byte(value)
				}
			}
		}
	}

	// Generate V plane (chroma)
	vData := frame.DataSlice(2, (width/2)*(height/2))
	if vData != nil {
		linesize := frame.Linesize()[2]
		for y := 0; y < height/2; y++ {
			for x := 0; x < width/2; x++ {
				// Red-ish chroma
				value := 128 + int(math.Cos(float64(frameIndex)*0.1)*50)
				if value < 0 {
					value = 0
				}
				if value > 255 {
					value = 255
				}

				offset := y*linesize + x
				if offset < len(vData) {
					vData[offset] = byte(value)
				}
			}
		}
	}

	return nil
}

func generateYUV422P(frame *astiav.Frame, width, height, frameIndex int) error {
	// Similar to YUV420P but with different chroma subsampling
	// Y plane (full resolution)
	yData := frame.DataSlice(0, width*height)
	if yData != nil {
		linesize := frame.Linesize()[0]
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				value := int((math.Sin(float64(x+frameIndex)*0.1)+math.Sin(float64(y+frameIndex)*0.1))*127 + 128)
				if value < 0 {
					value = 0
				}
				if value > 255 {
					value = 255
				}

				offset := y*linesize + x
				if offset < len(yData) {
					yData[offset] = byte(value)
				}
			}
		}
	}

	// U and V planes (half width, full height)
	for plane := 1; plane <= 2; plane++ {
		data := frame.DataSlice(plane, (width/2)*height)
		if data != nil {
			linesize := frame.Linesize()[plane]
			for y := 0; y < height; y++ {
				for x := 0; x < width/2; x++ {
					var value int
					if plane == 1 {
						value = 128 + int(math.Sin(float64(frameIndex)*0.1)*50)
					} else {
						value = 128 + int(math.Cos(float64(frameIndex)*0.1)*50)
					}

					if value < 0 {
						value = 0
					}
					if value > 255 {
						value = 255
					}

					offset := y*linesize + x
					if offset < len(data) {
						data[offset] = byte(value)
					}
				}
			}
		}
	}

	return nil
}

func generateYUV444P(frame *astiav.Frame, width, height, frameIndex int) error {
	// All planes at full resolution
	for plane := 0; plane < 3; plane++ {
		data := frame.DataSlice(plane, width*height)
		if data != nil {
			linesize := frame.Linesize()[plane]
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					var value int
					if plane == 0 {
						// Y (luminance)
						value = int((math.Sin(float64(x+frameIndex)*0.1)+math.Sin(float64(y+frameIndex)*0.1))*127 + 128)
					} else if plane == 1 {
						// U (chroma)
						value = 128 + int(math.Sin(float64(frameIndex)*0.1)*50)
					} else {
						// V (chroma)
						value = 128 + int(math.Cos(float64(frameIndex)*0.1)*50)
					}

					if value < 0 {
						value = 0
					}
					if value > 255 {
						value = 255
					}

					offset := y*linesize + x
					if offset < len(data) {
						data[offset] = byte(value)
					}
				}
			}
		}
	}

	return nil
}

func generateRGB24(frame *astiav.Frame, width, height, frameIndex int) error {
	data := frame.DataSlice(0, width*height*3)
	if data != nil {
		linesize := frame.Linesize()[0]
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Generate RGB values
				r := int((math.Sin(float64(x+frameIndex)*0.1) + 1) * 127)
				g := int((math.Sin(float64(y+frameIndex)*0.1) + 1) * 127)
				b := int((math.Sin(float64(frameIndex)*0.1) + 1) * 127)

				// Clamp values
				if r < 0 {
					r = 0
				}
				if r > 255 {
					r = 255
				}
				if g < 0 {
					g = 0
				}
				if g > 255 {
					g = 255
				}
				if b < 0 {
					b = 0
				}
				if b > 255 {
					b = 255
				}

				offset := y*linesize + x*3
				if offset+2 < len(data) {
					data[offset] = byte(r)
					data[offset+1] = byte(g)
					data[offset+2] = byte(b)
				}
			}
		}
	}

	return nil
}

func encodeFrame(codecContext *astiav.CodecContext, frame *astiav.Frame, packet *astiav.Packet, outputFile *os.File) error {
	// Send frame to encoder
	if err := codecContext.SendFrame(frame); err != nil {
		return fmt.Errorf("failed to send frame: %w", err)
	}

	// Receive packets from encoder
	for {
		if err := codecContext.ReceivePacket(packet); err != nil {
			if err == astiav.ErrEagain || err == astiav.ErrEof {
				break
			}
			return fmt.Errorf("failed to receive packet: %w", err)
		}

		// Write packet data to output file
		data := packet.Data()
		if data != nil && len(data) > 0 {
			if _, err := outputFile.Write(data); err != nil {
				return fmt.Errorf("failed to write packet data: %w", err)
			}
		}

		packet.Unref()
	}

	return nil
}
