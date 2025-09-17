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
	input  = flag.String("i", "", "the input audio file path")
	output = flag.String("o", "", "the output raw audio file path")
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
		log.Println("Usage: <binary path> -i <input audio file> -o <output raw audio file>")
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

	// Find audio stream
	audioStreamIndex := -1
	var audioStream *astiav.Stream
	for i, stream := range formatContext.Streams() {
		if stream.CodecParameters().MediaType() == astiav.MediaTypeAudio {
			audioStreamIndex = i
			audioStream = stream
			break
		}
	}

	if audioStreamIndex == -1 {
		log.Fatal("no audio stream found")
	}

	// Find decoder
	codec := astiav.FindDecoder(audioStream.CodecParameters().CodecID())
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
	if err := codecContext.FromCodecParameters(audioStream.CodecParameters()); err != nil {
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

	// Print audio stream information
	log.Printf("Audio stream found:")
	log.Printf("  Codec: %s", codec.Name())
	log.Printf("  Sample rate: %d Hz", codecContext.SampleRate())
	log.Printf("  Channels: %d", codecContext.ChannelLayout().Channels())
	log.Printf("  Sample format: %s", codecContext.SampleFormat().Name())
	log.Printf("  Duration: %v", audioStream.Duration())

	// Read and decode packets
	for {
		// Read packet
		if err := formatContext.ReadFrame(packet); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			log.Fatal(fmt.Errorf("failed to read frame: %w", err))
		}

		// Skip non-audio packets
		if packet.StreamIndex() != audioStreamIndex {
			packet.Unref()
			continue
		}

		// Decode packet
		if err := decodePacket(codecContext, packet, frame, outputFile); err != nil {
			log.Fatal(fmt.Errorf("failed to decode packet: %w", err))
		}

		packet.Unref()
	}

	// Flush decoder
	if err := decodePacket(codecContext, nil, frame, outputFile); err != nil {
		log.Fatal(fmt.Errorf("failed to flush decoder: %w", err))
	}

	log.Println("Audio decoding completed successfully")
}

func decodePacket(codecContext *astiav.CodecContext, packet *astiav.Packet, frame *astiav.Frame, outputFile *os.File) error {
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

		frame.Unref()
	}

	return nil
}

func writeFrame(frame *astiav.Frame, outputFile *os.File) error {
	// Get frame properties
	channels := frame.ChannelLayout().Channels()
	sampleFormat := frame.SampleFormat()
	samplesPerChannel := frame.NbSamples()

	// Calculate bytes per sample
	bytesPerSample := sampleFormat.BytesPerSample()
	if bytesPerSample <= 0 {
		return fmt.Errorf("unsupported sample format: %s", sampleFormat.Name())
	}

	if sampleFormat.IsPlanar() {
		// Planar format: interleave channels
		for i := 0; i < samplesPerChannel; i++ {
			for ch := 0; ch < channels; ch++ {
				// Get data slice for this channel
				data := frame.DataSlice(ch, samplesPerChannel*bytesPerSample)
				if data == nil {
					continue
				}

				start := i * bytesPerSample
				end := start + bytesPerSample
				if end <= len(data) {
					if _, err := outputFile.Write(data[start:end]); err != nil {
						return err
					}
				}
			}
		}
	} else {
		// Packed format: write directly
		totalBytes := samplesPerChannel * channels * bytesPerSample
		data := frame.DataSlice(0, totalBytes)
		if data != nil && len(data) >= totalBytes {
			if _, err := outputFile.Write(data[:totalBytes]); err != nil {
				return err
			}
		}
	}

	return nil
}
