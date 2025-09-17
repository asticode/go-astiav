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
	input  = flag.String("i", "", "the input video file path")
	output = flag.String("o", "", "the output raw video file path")
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
	if *input == "" || *output == "" {
		log.Println("Usage: <binary path> -i <input video file> -o <output raw video file>")
		return
	}

	// Allocate format context
	formatContext := astiav.AllocFormatContext()
	if formatContext == nil {
		log.Fatal("failed to allocate format context")
	}
	defer formatContext.Free()

	// Open input file
	if err := formatContext.OpenInput(*input, nil, nil); err != nil {
		log.Fatal(fmt.Errorf("failed to open input: %w", err))
	}
	defer formatContext.CloseInput()

	// Find stream info
	if err := formatContext.FindStreamInfo(nil); err != nil {
		log.Fatal(fmt.Errorf("failed to find stream info: %w", err))
	}

	// Find video stream
	videoStreamIndex := -1
	var videoStream *astiav.Stream
	for i, stream := range formatContext.Streams() {
		if stream.CodecParameters().MediaType() == astiav.MediaTypeVideo {
			videoStreamIndex = i
			videoStream = stream
			break
		}
	}

	if videoStreamIndex == -1 {
		log.Fatal("no video stream found")
	}

	// Find decoder
	codec := astiav.FindDecoder(videoStream.CodecParameters().CodecID())
	if codec == nil {
		log.Fatal("failed to find decoder")
	}

	// Allocate codec context
	codecContext := astiav.AllocCodecContext(codec)
	if codecContext == nil {
		log.Fatal("failed to allocate codec context")
	}
	defer codecContext.Free()

	// Copy codec parameters to codec context
	if err := codecContext.FromCodecParameters(videoStream.CodecParameters()); err != nil {
		log.Fatal(fmt.Errorf("failed to copy codec parameters: %w", err))
	}

	// Open codec
	if err := codecContext.Open(codec, nil); err != nil {
		log.Fatal(fmt.Errorf("failed to open codec: %w", err))
	}

	// Create output file
	outputFile, err := os.Create(*output)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create output file: %w", err))
	}
	defer outputFile.Close()

	// Allocate packet and frame
	packet := astiav.AllocPacket()
	defer packet.Free()

	frame := astiav.AllocFrame()
	defer frame.Free()

	// Print video stream information
	log.Printf("Video stream found:")
	log.Printf("  Codec: %s", codec.Name())
	log.Printf("  Resolution: %dx%d", codecContext.Width(), codecContext.Height())
	log.Printf("  Pixel format: %s", codecContext.PixelFormat().Name())
	log.Printf("  Frame rate: %v", codecContext.Framerate())
	log.Printf("  Duration: %v", videoStream.Duration())

	frameCount := 0

	// Read and decode packets
	for {
		// Read packet
		if err := formatContext.ReadFrame(packet); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			log.Fatal(fmt.Errorf("failed to read frame: %w", err))
		}

		// Skip non-video packets
		if packet.StreamIndex() != videoStreamIndex {
			packet.Unref()
			continue
		}

		// Decode packet
		if err := decodePacket(codecContext, packet, frame, outputFile, &frameCount); err != nil {
			log.Fatal(fmt.Errorf("failed to decode packet: %w", err))
		}

		packet.Unref()
	}

	// Flush decoder
	if err := decodePacket(codecContext, nil, frame, outputFile, &frameCount); err != nil {
		log.Fatal(fmt.Errorf("failed to flush decoder: %w", err))
	}

	log.Printf("Video decoding completed successfully. Decoded %d frames", frameCount)
}

func decodePacket(codecContext *astiav.CodecContext, packet *astiav.Packet, frame *astiav.Frame, outputFile *os.File, frameCount *int) error {
	// Send packet to decoder
	if err := codecContext.SendPacket(packet); err != nil {
		return fmt.Errorf("failed to send packet: %w", err)
	}

	// Receive frames from decoder
	for {
		if err := codecContext.ReceiveFrame(frame); err != nil {
			if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("failed to receive frame: %w", err)
		}

		// Write frame data to output file
		if err := writeFrame(frame, outputFile); err != nil {
			return fmt.Errorf("failed to write frame: %w", err)
		}

		*frameCount++
		if *frameCount%100 == 0 {
			log.Printf("Decoded %d frames", *frameCount)
		}

		frame.Unref()
	}

	return nil
}

func writeFrame(frame *astiav.Frame, outputFile *os.File) error {
	// Get frame properties
	width := frame.Width()
	height := frame.Height()
	pixelFormat := frame.PixelFormat()

	// Write frame data based on pixel format
	switch pixelFormat {
	case astiav.PixelFormatYuv420P:
		return writeYUV420P(frame, outputFile, width, height)
	case astiav.PixelFormatYuv422P:
		return writeYUV422P(frame, outputFile, width, height)
	case astiav.PixelFormatYuv444P:
		return writeYUV444P(frame, outputFile, width, height)
	case astiav.PixelFormatRgb24:
		return writeRGB24(frame, outputFile, width, height)
	case astiav.PixelFormatBgr24:
		return writeBGR24(frame, outputFile, width, height)
	default:
		// For other formats, try to write raw data
		return writeRawFrame(frame, outputFile)
	}
}

func writeYUV420P(frame *astiav.Frame, outputFile *os.File, width, height int) error {
	// Y plane
	yData := frame.DataSlice(0, width*height)
	if yData != nil {
		linesize := frame.Linesize()[0]
		for y := 0; y < height; y++ {
			start := y * linesize
			end := start + width
			if end <= len(yData) {
				if _, err := outputFile.Write(yData[start:end]); err != nil {
					return err
				}
			}
		}
	}

	// U plane (half width, half height)
	uData := frame.DataSlice(1, (width/2)*(height/2))
	if uData != nil {
		linesize := frame.Linesize()[1]
		for y := 0; y < height/2; y++ {
			start := y * linesize
			end := start + width/2
			if end <= len(uData) {
				if _, err := outputFile.Write(uData[start:end]); err != nil {
					return err
				}
			}
		}
	}

	// V plane (half width, half height)
	vData := frame.DataSlice(2, (width/2)*(height/2))
	if vData != nil {
		linesize := frame.Linesize()[2]
		for y := 0; y < height/2; y++ {
			start := y * linesize
			end := start + width/2
			if end <= len(vData) {
				if _, err := outputFile.Write(vData[start:end]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func writeYUV422P(frame *astiav.Frame, outputFile *os.File, width, height int) error {
	// Y plane
	yData := frame.DataSlice(0, width*height)
	if yData != nil {
		linesize := frame.Linesize()[0]
		for y := 0; y < height; y++ {
			start := y * linesize
			end := start + width
			if end <= len(yData) {
				if _, err := outputFile.Write(yData[start:end]); err != nil {
					return err
				}
			}
		}
	}

	// U plane (half width, full height)
	uData := frame.DataSlice(1, (width/2)*height)
	if uData != nil {
		linesize := frame.Linesize()[1]
		for y := 0; y < height; y++ {
			start := y * linesize
			end := start + width/2
			if end <= len(uData) {
				if _, err := outputFile.Write(uData[start:end]); err != nil {
					return err
				}
			}
		}
	}

	// V plane (half width, full height)
	vData := frame.DataSlice(2, (width/2)*height)
	if vData != nil {
		linesize := frame.Linesize()[2]
		for y := 0; y < height; y++ {
			start := y * linesize
			end := start + width/2
			if end <= len(vData) {
				if _, err := outputFile.Write(vData[start:end]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func writeYUV444P(frame *astiav.Frame, outputFile *os.File, width, height int) error {
	// Y, U, V planes (all full resolution)
	for plane := 0; plane < 3; plane++ {
		data := frame.DataSlice(plane, width*height)
		if data != nil {
			linesize := frame.Linesize()[plane]
			for y := 0; y < height; y++ {
				start := y * linesize
				end := start + width
				if end <= len(data) {
					if _, err := outputFile.Write(data[start:end]); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func writeRGB24(frame *astiav.Frame, outputFile *os.File, width, height int) error {
	data := frame.DataSlice(0, width*height*3)
	if data != nil {
		linesize := frame.Linesize()[0]
		for y := 0; y < height; y++ {
			start := y * linesize
			end := start + width*3
			if end <= len(data) {
				if _, err := outputFile.Write(data[start:end]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func writeBGR24(frame *astiav.Frame, outputFile *os.File, width, height int) error {
	data := frame.DataSlice(0, width*height*3)
	if data != nil {
		linesize := frame.Linesize()[0]
		for y := 0; y < height; y++ {
			start := y * linesize
			end := start + width*3
			if end <= len(data) {
				if _, err := outputFile.Write(data[start:end]); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func writeRawFrame(frame *astiav.Frame, outputFile *os.File) error {
	// Try to get the frame data as bytes
	data, err := frame.Data().Bytes(1)
	if err != nil {
		return fmt.Errorf("failed to get frame bytes: %w", err)
	}

	if _, err := outputFile.Write(data); err != nil {
		return err
	}

	return nil
}
