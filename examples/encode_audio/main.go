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
	output = flag.String("o", "", "the output audio file path")
	codec  = flag.String("c", "mp2", "the audio codec to use (mp2, aac, etc.)")
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
		log.Println("Usage: <binary path> -o <output audio file> [-c <codec>]")
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
	codecContext.SetBitRate(64000)
	codecContext.SetSampleFormat(astiav.SampleFormatS16) // 16-bit signed integer
	codecContext.SetSampleRate(44100)
	codecContext.SetChannelLayout(astiav.ChannelLayoutStereo)

	// Check if the encoder supports the sample format
	if !checkSampleFormat(encoder, codecContext.SampleFormat()) {
		log.Printf("Sample format %s not supported by encoder, trying alternatives...", codecContext.SampleFormat().Name())

		// Try common formats
		supportedFormats := []astiav.SampleFormat{
			astiav.SampleFormatFltp,
			astiav.SampleFormatS16P,
			astiav.SampleFormatS32,
			astiav.SampleFormatFlt,
		}

		found := false
		for _, sf := range supportedFormats {
			if checkSampleFormat(encoder, sf) {
				codecContext.SetSampleFormat(sf)
				log.Printf("Using sample format: %s", sf.Name())
				found = true
				break
			}
		}

		if !found {
			log.Fatal("No supported sample format found")
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
	frame.SetNbSamples(codecContext.FrameSize())
	frame.SetSampleFormat(codecContext.SampleFormat())
	frame.SetChannelLayout(codecContext.ChannelLayout())
	frame.SetSampleRate(codecContext.SampleRate())

	// Allocate frame buffer
	if err := frame.AllocBuffer(0); err != nil {
		log.Fatal(fmt.Errorf("failed to allocate frame buffer: %w", err))
	}

	// Print encoding information
	log.Printf("Encoding audio:")
	log.Printf("  Codec: %s", encoder.Name())
	log.Printf("  Sample rate: %d Hz", codecContext.SampleRate())
	log.Printf("  Channels: %d", codecContext.ChannelLayout().Channels())
	log.Printf("  Sample format: %s", codecContext.SampleFormat().Name())
	log.Printf("  Frame size: %d samples", codecContext.FrameSize())
	log.Printf("  Bit rate: %d bps", codecContext.BitRate())

	// Generate and encode audio
	duration := 5.0 // 5 seconds
	sampleRate := float64(codecContext.SampleRate())
	totalSamples := int(duration * sampleRate)
	samplesGenerated := 0

	for samplesGenerated < totalSamples {
		// Make frame writable
		if err := frame.MakeWritable(); err != nil {
			log.Fatal(fmt.Errorf("failed to make frame writable: %w", err))
		}

		// Generate audio data (sine wave)
		if err := generateAudioData(frame, samplesGenerated, codecContext.SampleRate()); err != nil {
			log.Fatal(fmt.Errorf("failed to generate audio data: %w", err))
		}

		// Set frame timestamp
		frame.SetPts(int64(samplesGenerated))

		// Encode frame
		if err := encodeFrame(codecContext, frame, packet, outputFile); err != nil {
			log.Fatal(fmt.Errorf("failed to encode frame: %w", err))
		}

		samplesGenerated += frame.NbSamples()
	}

	// Flush encoder
	if err := encodeFrame(codecContext, nil, packet, outputFile); err != nil {
		log.Fatal(fmt.Errorf("failed to flush encoder: %w", err))
	}

	log.Println("Audio encoding completed successfully")
}

func checkSampleFormat(codec *astiav.Codec, sampleFormat astiav.SampleFormat) bool {
	formats := codec.SampleFormats()
	for _, sf := range formats {
		if sf == sampleFormat {
			return true
		}
	}
	return false
}

func generateAudioData(frame *astiav.Frame, startSample int, sampleRate int) error {
	channels := frame.ChannelLayout().Channels()
	samplesPerChannel := frame.NbSamples()
	sampleFormat := frame.SampleFormat()

	// Generate sine wave at 440 Hz (A4 note)
	frequency := 440.0

	if sampleFormat.IsPlanar() {
		// Planar format: each channel has its own data array
		for ch := 0; ch < channels; ch++ {
			if err := generateChannelData(frame, ch, startSample, samplesPerChannel, sampleRate, frequency, sampleFormat); err != nil {
				return err
			}
		}
	} else {
		// Packed format: all channels interleaved in one array
		if err := generatePackedData(frame, startSample, samplesPerChannel, channels, sampleRate, frequency, sampleFormat); err != nil {
			return err
		}
	}

	return nil
}

func generateChannelData(frame *astiav.Frame, channel, startSample, samplesPerChannel, sampleRate int, frequency float64, sampleFormat astiav.SampleFormat) error {
	bytesPerSample := sampleFormat.BytesPerSample()
	data := frame.DataSlice(channel, samplesPerChannel*bytesPerSample)
	if data == nil {
		return fmt.Errorf("failed to get data slice for channel %d", channel)
	}

	for i := 0; i < samplesPerChannel; i++ {
		sampleIndex := startSample + i
		t := float64(sampleIndex) / float64(sampleRate)
		value := math.Sin(2 * math.Pi * frequency * t)

		// Apply different phase for stereo effect
		if channel == 1 {
			value = math.Sin(2*math.Pi*frequency*t + math.Pi/4)
		}

		if err := writeSampleToData(data, i*bytesPerSample, value, sampleFormat); err != nil {
			return err
		}
	}

	return nil
}

func generatePackedData(frame *astiav.Frame, startSample, samplesPerChannel, channels, sampleRate int, frequency float64, sampleFormat astiav.SampleFormat) error {
	bytesPerSample := sampleFormat.BytesPerSample()
	totalBytes := samplesPerChannel * channels * bytesPerSample
	data := frame.DataSlice(0, totalBytes)
	if data == nil {
		return fmt.Errorf("failed to get data slice")
	}

	for i := 0; i < samplesPerChannel; i++ {
		sampleIndex := startSample + i
		t := float64(sampleIndex) / float64(sampleRate)

		for ch := 0; ch < channels; ch++ {
			value := math.Sin(2 * math.Pi * frequency * t)

			// Apply different phase for stereo effect
			if ch == 1 {
				value = math.Sin(2*math.Pi*frequency*t + math.Pi/4)
			}

			offset := (i*channels + ch) * bytesPerSample
			if err := writeSampleToData(data, offset, value, sampleFormat); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeSampleToData(data []byte, offset int, value float64, sampleFormat astiav.SampleFormat) error {
	switch sampleFormat {
	case astiav.SampleFormatS16, astiav.SampleFormatS16P:
		// 16-bit signed integer
		sample := int16(value * 32767)
		data[offset] = byte(sample)
		data[offset+1] = byte(sample >> 8)
	case astiav.SampleFormatS32, astiav.SampleFormatS32P:
		// 32-bit signed integer
		sample := int32(value * 2147483647)
		data[offset] = byte(sample)
		data[offset+1] = byte(sample >> 8)
		data[offset+2] = byte(sample >> 16)
		data[offset+3] = byte(sample >> 24)
	case astiav.SampleFormatFlt, astiav.SampleFormatFltp:
		// 32-bit float
		sample := float32(value)
		bits := math.Float32bits(sample)
		data[offset] = byte(bits)
		data[offset+1] = byte(bits >> 8)
		data[offset+2] = byte(bits >> 16)
		data[offset+3] = byte(bits >> 24)
	default:
		return fmt.Errorf("unsupported sample format: %s", sampleFormat.Name())
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
