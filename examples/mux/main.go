package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/asticode/go-astiav"
)

var (
	output = flag.String("o", "", "the output media file path")
	format = flag.String("f", "", "the output format (mp4, avi, mkv, etc.)")
)

type OutputStream struct {
	codecContext *astiav.CodecContext
	frame        *astiav.Frame
	stream       *astiav.Stream
	nextPts      int64
}

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
		log.Println("Usage: <binary path> -o <output media file> [-f <format>]")
		return
	}

	// Allocate output format context
	formatContext, err := astiav.AllocOutputFormatContext(nil, *format, *output)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to allocate output format context: %w", err))
	}
	defer formatContext.Free()

	outputFormat := formatContext.OutputFormat()
	if outputFormat == nil {
		log.Fatal("failed to get output format")
	}

	// Create video stream
	var videoStream *OutputStream
	if outputFormat.VideoCodec() != astiav.CodecIDNone {
		videoStream = addVideoStream(formatContext, outputFormat.VideoCodec())
		if videoStream == nil {
			log.Fatal("failed to add video stream")
		}
		defer videoStream.free()
	}

	// Create audio stream
	var audioStream *OutputStream
	if outputFormat.AudioCodec() != astiav.CodecIDNone {
		audioStream = addAudioStream(formatContext, outputFormat.AudioCodec())
		if audioStream == nil {
			log.Fatal("failed to add audio stream")
		}
		defer audioStream.free()
	}

	// Open output file
	if !outputFormat.Flags().Has(astiav.IOFormatFlagNofile) {
		ioContext, err := astiav.OpenIOContext(*output, astiav.NewIOContextFlags(astiav.IOContextFlagWrite), nil, nil)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to open output file: %w", err))
		}
		defer ioContext.Close()
		formatContext.SetPb(ioContext)
	}

	// Write file header
	if err := formatContext.WriteHeader(nil); err != nil {
		log.Fatal(fmt.Errorf("failed to write header: %w", err))
	}

	// Print muxing information
	log.Printf("Muxing to: %s", *output)
	log.Printf("Format: %s", outputFormat.Name())
	if videoStream != nil {
		log.Printf("Video codec: %s", videoStream.codecContext.Codec().Name())
	}
	if audioStream != nil {
		log.Printf("Audio codec: %s", audioStream.codecContext.Codec().Name())
	}

	// Generate and mux frames
	packet := astiav.AllocPacket()
	defer packet.Free()

	duration := 5.0 // 5 seconds
	frameCount := 0

	for {
		// Determine which stream to write next
		writeVideo := false
		writeAudio := false

		if videoStream != nil && audioStream != nil {
			// Compare timestamps to decide which stream to write
			videoPts := float64(videoStream.nextPts) * videoStream.codecContext.TimeBase().Float64()
			audioPts := float64(audioStream.nextPts) * audioStream.codecContext.TimeBase().Float64()

			if videoPts <= audioPts {
				writeVideo = true
			} else {
				writeAudio = true
			}
		} else if videoStream != nil {
			writeVideo = true
		} else if audioStream != nil {
			writeAudio = true
		}

		// Check if we've reached the duration
		if videoStream != nil {
			videoPts := float64(videoStream.nextPts) * videoStream.codecContext.TimeBase().Float64()
			if videoPts >= duration {
				break
			}
		}
		if audioStream != nil {
			audioPts := float64(audioStream.nextPts) * audioStream.codecContext.TimeBase().Float64()
			if audioPts >= duration {
				break
			}
		}

		// Write video frame
		if writeVideo && videoStream != nil {
			if err := writeVideoFrame(formatContext, videoStream, packet); err != nil {
				log.Fatal(fmt.Errorf("failed to write video frame: %w", err))
			}
			frameCount++
			if frameCount%25 == 0 {
				log.Printf("Muxed %d frames", frameCount)
			}
		}

		// Write audio frame
		if writeAudio && audioStream != nil {
			if err := writeAudioFrame(formatContext, audioStream, packet); err != nil {
				log.Fatal(fmt.Errorf("failed to write audio frame: %w", err))
			}
		}
	}

	// Flush encoders
	if videoStream != nil {
		flushEncoder(formatContext, videoStream, packet)
	}
	if audioStream != nil {
		flushEncoder(formatContext, audioStream, packet)
	}

	// Write file trailer
	if err := formatContext.WriteTrailer(); err != nil {
		log.Fatal(fmt.Errorf("failed to write trailer: %w", err))
	}

	log.Println("Muxing completed successfully")
}

func addVideoStream(formatContext *astiav.FormatContext, codecId astiav.CodecID) *OutputStream {
	// Find encoder
	codec := astiav.FindEncoder(codecId)
	if codec == nil {
		log.Printf("Could not find encoder for codec ID: %d", codecId)
		return nil
	}

	// Create stream
	stream := formatContext.NewStream(codec)
	if stream == nil {
		log.Printf("Could not allocate stream")
		return nil
	}

	// Allocate codec context
	codecContext := astiav.AllocCodecContext(codec)
	if codecContext == nil {
		log.Printf("Could not allocate codec context")
		return nil
	}

	// Set codec parameters
	codecContext.SetCodecId(codecId)
	codecContext.SetBitRate(400000)
	codecContext.SetWidth(320)
	codecContext.SetHeight(240)
	codecContext.SetTimeBase(astiav.NewRational(1, 25))
	codecContext.SetFramerate(astiav.NewRational(25, 1))
	codecContext.SetPixelFormat(astiav.PixelFormatYuv420P)
	codecContext.SetGopSize(10)
	codecContext.SetMaxBFrames(1)

	// Some formats want stream headers to be separate
	if formatContext.OutputFormat().Flags().Has(astiav.IOFormatFlagGlobalheader) {
		codecContext.SetFlags(codecContext.Flags().Add(astiav.CodecContextFlagGlobalHeader))
	}

	// Open codec
	if err := codecContext.Open(codec, nil); err != nil {
		log.Printf("Could not open codec: %v", err)
		codecContext.Free()
		return nil
	}

	// Copy codec parameters to stream
	if err := stream.CodecParameters().FromCodecContext(codecContext); err != nil {
		log.Printf("Could not copy codec parameters: %v", err)
		codecContext.Free()
		return nil
	}

	// Allocate frame
	frame := astiav.AllocFrame()
	if frame == nil {
		log.Printf("Could not allocate frame")
		codecContext.Free()
		return nil
	}

	frame.SetWidth(codecContext.Width())
	frame.SetHeight(codecContext.Height())
	frame.SetPixelFormat(codecContext.PixelFormat())

	if err := frame.AllocBuffer(32); err != nil {
		log.Printf("Could not allocate frame buffer: %v", err)
		frame.Free()
		codecContext.Free()
		return nil
	}

	return &OutputStream{
		codecContext: codecContext,
		frame:        frame,
		stream:       stream,
		nextPts:      0,
	}
}

func addAudioStream(formatContext *astiav.FormatContext, codecId astiav.CodecID) *OutputStream {
	// Find encoder
	codec := astiav.FindEncoder(codecId)
	if codec == nil {
		log.Printf("Could not find encoder for codec ID: %d", codecId)
		return nil
	}

	// Create stream
	stream := formatContext.NewStream(codec)
	if stream == nil {
		log.Printf("Could not allocate stream")
		return nil
	}

	// Allocate codec context
	codecContext := astiav.AllocCodecContext(codec)
	if codecContext == nil {
		log.Printf("Could not allocate codec context")
		return nil
	}

	// Set codec parameters
	codecContext.SetCodecId(codecId)
	codecContext.SetBitRate(64000)
	codecContext.SetSampleFormat(astiav.SampleFormatFltp)
	codecContext.SetSampleRate(44100)
	codecContext.SetChannelLayout(astiav.ChannelLayoutStereo)

	// Check if encoder supports the sample format
	if !checkSampleFormat(codec, codecContext.SampleFormat()) {
		// Try alternatives
		supportedFormats := []astiav.SampleFormat{
			astiav.SampleFormatS16,
			astiav.SampleFormatS16P,
			astiav.SampleFormatFlt,
		}

		found := false
		for _, sf := range supportedFormats {
			if checkSampleFormat(codec, sf) {
				codecContext.SetSampleFormat(sf)
				found = true
				break
			}
		}

		if !found {
			log.Printf("No supported sample format found")
			codecContext.Free()
			return nil
		}
	}

	// Some formats want stream headers to be separate
	if formatContext.OutputFormat().Flags().Has(astiav.IOFormatFlagGlobalheader) {
		codecContext.SetFlags(codecContext.Flags().Add(astiav.CodecContextFlagGlobalHeader))
	}

	// Open codec
	if err := codecContext.Open(codec, nil); err != nil {
		log.Printf("Could not open codec: %v", err)
		codecContext.Free()
		return nil
	}

	// Copy codec parameters to stream
	if err := stream.CodecParameters().FromCodecContext(codecContext); err != nil {
		log.Printf("Could not copy codec parameters: %v", err)
		codecContext.Free()
		return nil
	}

	// Allocate frame
	frame := astiav.AllocFrame()
	if frame == nil {
		log.Printf("Could not allocate frame")
		codecContext.Free()
		return nil
	}

	frame.SetNbSamples(codecContext.FrameSize())
	frame.SetSampleFormat(codecContext.SampleFormat())
	frame.SetChannelLayout(codecContext.ChannelLayout())
	frame.SetSampleRate(codecContext.SampleRate())

	if err := frame.AllocBuffer(0); err != nil {
		log.Printf("Could not allocate frame buffer: %v", err)
		frame.Free()
		codecContext.Free()
		return nil
	}

	return &OutputStream{
		codecContext: codecContext,
		frame:        frame,
		stream:       stream,
		nextPts:      0,
	}
}

func (os *OutputStream) free() {
	if os.frame != nil {
		os.frame.Free()
	}
	if os.codecContext != nil {
		os.codecContext.Free()
	}
}

func writeVideoFrame(formatContext *astiav.FormatContext, os *OutputStream, packet *astiav.Packet) error {
	// Make frame writable
	if err := os.frame.MakeWritable(); err != nil {
		return err
	}

	// Generate video data
	generateVideoFrame(os.frame, int(os.nextPts))

	// Set frame timestamp
	os.frame.SetPts(os.nextPts)

	// Encode frame
	if err := os.codecContext.SendFrame(os.frame); err != nil {
		return err
	}

	for {
		if err := os.codecContext.ReceivePacket(packet); err != nil {
			if err == astiav.ErrEagain || err == astiav.ErrEof {
				break
			}
			return err
		}

		// Rescale packet timestamp
		packet.RescaleTs(os.codecContext.TimeBase(), os.stream.TimeBase())
		packet.SetStreamIndex(os.stream.Index())

		// Write packet
		if err := formatContext.WriteInterleavedFrame(packet); err != nil {
			return err
		}

		packet.Unref()
	}

	os.nextPts++
	return nil
}

func writeAudioFrame(formatContext *astiav.FormatContext, os *OutputStream, packet *astiav.Packet) error {
	// Make frame writable
	if err := os.frame.MakeWritable(); err != nil {
		return err
	}

	// Generate audio data
	generateAudioFrame(os.frame, int(os.nextPts), os.codecContext.SampleRate())

	// Set frame timestamp
	os.frame.SetPts(os.nextPts)

	// Encode frame
	if err := os.codecContext.SendFrame(os.frame); err != nil {
		return err
	}

	for {
		if err := os.codecContext.ReceivePacket(packet); err != nil {
			if err == astiav.ErrEagain || err == astiav.ErrEof {
				break
			}
			return err
		}

		// Rescale packet timestamp
		packet.RescaleTs(os.codecContext.TimeBase(), os.stream.TimeBase())
		packet.SetStreamIndex(os.stream.Index())

		// Write packet
		if err := formatContext.WriteInterleavedFrame(packet); err != nil {
			return err
		}

		packet.Unref()
	}

	os.nextPts += int64(os.frame.NbSamples())
	return nil
}

func flushEncoder(formatContext *astiav.FormatContext, os *OutputStream, packet *astiav.Packet) error {
	// Send NULL frame to flush encoder
	if err := os.codecContext.SendFrame(nil); err != nil {
		return err
	}

	for {
		if err := os.codecContext.ReceivePacket(packet); err != nil {
			if err == astiav.ErrEagain || err == astiav.ErrEof {
				break
			}
			return err
		}

		// Rescale packet timestamp
		packet.RescaleTs(os.codecContext.TimeBase(), os.stream.TimeBase())
		packet.SetStreamIndex(os.stream.Index())

		// Write packet
		if err := formatContext.WriteInterleavedFrame(packet); err != nil {
			return err
		}

		packet.Unref()
	}

	return nil
}

func generateVideoFrame(frame *astiav.Frame, frameIndex int) {
	width := frame.Width()
	height := frame.Height()

	// Generate Y plane
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

	// Generate U and V planes
	for plane := 1; plane <= 2; plane++ {
		data := frame.DataSlice(plane, (width/2)*(height/2))
		if data != nil {
			linesize := frame.Linesize()[plane]
			for y := 0; y < height/2; y++ {
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
}

func generateAudioFrame(frame *astiav.Frame, startSample int, sampleRate int) {
	channels := frame.ChannelLayout().Channels()
	samplesPerChannel := frame.NbSamples()
	sampleFormat := frame.SampleFormat()

	frequency := 440.0 // A4 note

	if sampleFormat.IsPlanar() {
		// Planar format
		for ch := 0; ch < channels; ch++ {
			generateAudioChannelData(frame, ch, startSample, samplesPerChannel, sampleRate, frequency, sampleFormat)
		}
	} else {
		// Packed format
		generateAudioPackedData(frame, startSample, samplesPerChannel, channels, sampleRate, frequency, sampleFormat)
	}
}

func generateAudioChannelData(frame *astiav.Frame, channel, startSample, samplesPerChannel, sampleRate int, frequency float64, sampleFormat astiav.SampleFormat) {
	bytesPerSample := sampleFormat.BytesPerSample()
	data := frame.DataSlice(channel, samplesPerChannel*bytesPerSample)
	if data == nil {
		return
	}

	for i := 0; i < samplesPerChannel; i++ {
		sampleIndex := startSample + i
		t := float64(sampleIndex) / float64(sampleRate)
		value := math.Sin(2 * math.Pi * frequency * t)

		// Apply different phase for stereo effect
		if channel == 1 {
			value = math.Sin(2*math.Pi*frequency*t + math.Pi/4)
		}

		writeAudioSample(data, i*bytesPerSample, value, sampleFormat)
	}
}

func generateAudioPackedData(frame *astiav.Frame, startSample, samplesPerChannel, channels, sampleRate int, frequency float64, sampleFormat astiav.SampleFormat) {
	bytesPerSample := sampleFormat.BytesPerSample()
	totalBytes := samplesPerChannel * channels * bytesPerSample
	data := frame.DataSlice(0, totalBytes)
	if data == nil {
		return
	}

	for i := 0; i < samplesPerChannel; i++ {
		sampleIndex := startSample + i
		t := float64(sampleIndex) / float64(sampleRate)

		for ch := 0; ch < channels; ch++ {
			value := math.Sin(2 * math.Pi * frequency * t)

			if ch == 1 {
				value = math.Sin(2*math.Pi*frequency*t + math.Pi/4)
			}

			offset := (i*channels + ch) * bytesPerSample
			writeAudioSample(data, offset, value, sampleFormat)
		}
	}
}

func writeAudioSample(data []byte, offset int, value float64, sampleFormat astiav.SampleFormat) {
	switch sampleFormat {
	case astiav.SampleFormatS16, astiav.SampleFormatS16P:
		sample := int16(value * 32767)
		data[offset] = byte(sample)
		data[offset+1] = byte(sample >> 8)
	case astiav.SampleFormatFlt, astiav.SampleFormatFltp:
		sample := float32(value)
		bits := math.Float32bits(sample)
		data[offset] = byte(bits)
		data[offset+1] = byte(bits >> 8)
		data[offset+2] = byte(bits >> 16)
		data[offset+3] = byte(bits >> 24)
	}
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
