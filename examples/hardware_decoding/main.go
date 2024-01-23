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
	input       = flag.String("i", "", "the input path")
	device_type = flag.String("d", "", "the hardware device type like: cuda")
)

type stream struct {
	decCodec        *astiav.Codec
	decCodecContext *astiav.CodecContext
	hwDeviceContext *astiav.HardwareDeviceContext
	inputStream     *astiav.Stream
}

func hw_decoder_init(ctx *astiav.CodecContext, t astiav.HardwareDeviceType) (*astiav.HardwareDeviceContext, error) {
	hdc, err := astiav.CreateHardwareDeviceContext(t, "", &astiav.Dictionary{})
	if err != nil {
		return nil, fmt.Errorf("Unable to create hardware device context: %w", err)
	}
	ctx.SetHardwareDeviceContext(hdc)
	return hdc, nil
}

func main() {
	// Handle ffmpeg logs
	astiav.SetLogLevel(astiav.LogLevelDebug)
	astiav.SetLogCallback(func(l astiav.LogLevel, fmt, msg, parent string) {
		log.Printf("ffmpeg log: %s (level: %d)\n", strings.TrimSpace(msg), l)
	})

	// Parse flags
	flag.Parse()

	// Usage
	if *input == "" || *device_type == "" {
		log.Println("Usage: <binary path> -d <device type> -i <input path>")
		return
	}

	hw_device := astiav.FindHardwareDeviceTypeByName(*device_type)
	if hw_device == astiav.HardwareDeviceTypeNone {
		log.Fatal(errors.New("main: hardware device not found"))
	}

	// Alloc packet
	pkt := astiav.AllocPacket()
	defer pkt.Free()

	gpu_f := astiav.AllocFrame()
	defer gpu_f.Free()

	sw_f := astiav.AllocFrame()
	defer sw_f.Free()

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

	var hw_pix_fmt astiav.PixelFormat

	// Loop through streams
	streams := make(map[int]*stream) // Indexed by input stream index
	for _, is := range inputFormatContext.Streams() {
		var err error

		// Only process video
		if is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
			continue
		}

		// Create stream
		s := &stream{inputStream: is}

		// Find decoder
		if s.decCodec = astiav.FindDecoder(is.CodecParameters().CodecID()); s.decCodec == nil {
			log.Fatal(errors.New("main: codec is nil"))
		}

		// Alloc codec context
		if s.decCodecContext = astiav.AllocCodecContext(s.decCodec); s.decCodecContext == nil {
			log.Fatal(errors.New("main: codec context is nil"))
		}
		defer s.decCodecContext.Free()

		hw_pix_fmt, err = astiav.FindSuitableHardwareFormat(s.decCodec, hw_device)
		if err != nil {
			log.Fatal(fmt.Errorf("main: find decoder hw format fails: %w", err))
		} else {
			log.Printf("Using hw_pix_fmt: %s", hw_pix_fmt.Name())
		}

		// Update codec context
		if err := is.CodecParameters().ToCodecContext(s.decCodecContext); err != nil {
			log.Fatal(fmt.Errorf("main: updating codec context failed: %w", err))
		}

		if s.hwDeviceContext, err = hw_decoder_init(s.decCodecContext, hw_device); err != nil {
			log.Fatal(fmt.Errorf("main: init hardware device failed: %w", err))
		}

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
			if err := s.decCodecContext.ReceiveFrame(gpu_f); err != nil {
				if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
					break
				}
				log.Fatal(fmt.Errorf("main: receiving frame failed: %w", err))
			}

			if gpu_f.PixelFormat() == hw_pix_fmt {
				err := gpu_f.TransferHardwareData(sw_f)
				if err != nil {
					log.Fatal(fmt.Errorf("main: Unable to transfer frame from gpu: %w", err))
				}

				data, err := sw_f.Data().Bytes(1)

				if err != nil {
					log.Fatal(fmt.Errorf("main: Unable to get frame bytes: %w", err))
				}

				sw_f.SetPts(gpu_f.Pts())

				// Do something with decoded frame
				log.Printf("new frame: pts: %d gpu frame pix fmt: %s, sw frame (transfered) pix fmt: %s, size: %d", sw_f.Pts(), gpu_f.PixelFormat().Name(), sw_f.PixelFormat().Name(), len(data))
			} else {
				log.Fatal(fmt.Errorf("main: Mismatch in pixel format: gpu: %s sw (transfered): %s", gpu_f.PixelFormat().Name(), sw_f.PixelFormat().Name()))
			}

		}
	}

	// Success
	log.Println("success")
}
