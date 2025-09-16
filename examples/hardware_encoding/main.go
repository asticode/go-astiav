package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/asticode/go-astiav"
	"github.com/asticode/go-astikit"
)

var (
	width       int
	height      int
	hwDeviceCtx *astiav.HardwareDeviceContext
	c           = astikit.NewCloser()
)

func setHwframeCtx(avctx *astiav.CodecContext, hwDeviceCtx *astiav.HardwareDeviceContext) error {
	hwFramesRef := astiav.AllocHardwareFramesContext(hwDeviceCtx)
	if hwFramesRef == nil {
		return errors.New("failed to create VAAPI frame context")
	}

	hwFramesRef.SetHardwarePixelFormat(astiav.PixelFormatVaapi)
	hwFramesRef.SetSoftwarePixelFormat(astiav.PixelFormatNv12)
	hwFramesRef.SetWidth(width)
	hwFramesRef.SetHeight(height)
	hwFramesRef.SetInitialPoolSize(20)

	if err := hwFramesRef.Initialize(); err != nil {
		hwFramesRef.Free()
		return fmt.Errorf("failed to initialize VAAPI frame context: %w", err)
	}

	avctx.SetHardwareFramesContext(hwFramesRef)
	return nil
}

func encodeWrite(avctx *astiav.CodecContext, frame *astiav.Frame, fout *os.File) error {
	encPkt := astiav.AllocPacket()
	if encPkt == nil {
		return errors.New("failed to allocate packet")
	}
	defer encPkt.Free()

	if err := avctx.SendFrame(frame); err != nil {
		return fmt.Errorf("error code: %w", err)
	}

	for {
		err := avctx.ReceivePacket(encPkt)
		if err != nil {
			if errors.Is(err, astiav.ErrEagain) {
				return nil // Like C code: ret = ((ret == AVERROR(EAGAIN)) ? 0 : -1);
			}
			return fmt.Errorf("receive packet failed: %w", err)
		}

		encPkt.SetStreamIndex(0)
		data := encPkt.Data()
		n, err := fout.Write(data)
		encPkt.Unref()
		if err != nil {
			return fmt.Errorf("write packet data failed: %w", err)
		}
		if n != len(data) {
			return fmt.Errorf("write size mismatch: wrote %d, expected %d", n, len(data))
		}
	}
}

func main() {
	// Handle ffmpeg logs
	astiav.SetLogLevel(astiav.LogLevelInfo)
	astiav.SetLogCallback(func(cl astiav.Classer, l astiav.LogLevel, fmt, msg string) {
		var cs string
		if cl != nil {
			if class := cl.Class(); class != nil {
				cs = " - class: " + class.String()
			}
		}
		log.Printf("ffmpeg log: %s%s - level: %d\n", strings.TrimSpace(msg), cs, l)
	})

	if len(os.Args) < 5 {
		fmt.Fprintf(os.Stderr, "Usage: %s <width> <height> <input file> <output file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s 320 240 input.yuv output.h264\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Note: This example uses VAAPI hardware encoding\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Note: Input file should be in NV12 format (YUV420 planar)\n")
		os.Exit(1)
	}

	var err error
	width, err = strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Invalid width:", err)
	}

	height, err = strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("Invalid height:", err)
	}

	size := width * height

	// We use an astikit.Closer to free all resources properly
	defer c.Close()

	fin, err := os.Open(os.Args[3])
	if err != nil {
		log.Fatalf("Fail to open input file: %v", err)
	}
	defer fin.Close()

	fout, err := os.OpenFile(os.Args[4], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Fail to open output file: %v", err)
	}
	defer fout.Close()

	// Create VAAPI hardware device context
	hwDeviceCtx, err = astiav.CreateHardwareDeviceContext(astiav.HardwareDeviceTypeVAAPI, "", nil, 0)
	if err != nil {
		log.Fatalf("Failed to create a VAAPI device: %v", err)
	}
	c.Add(hwDeviceCtx.Free)

	// Use VAAPI H.264 encoder
	encName := "h264_vaapi"
	codec := astiav.FindEncoderByName(encName)
	if codec == nil {
		log.Fatal("Could not find h264_vaapi encoder")
	}

	avctx := astiav.AllocCodecContext(codec)
	if avctx == nil {
		log.Fatal("Failed to allocate codec context")
	}
	c.Add(avctx.Free)

	avctx.SetWidth(width)
	avctx.SetHeight(height)
	avctx.SetTimeBase(astiav.NewRational(1, 25))
	avctx.SetFramerate(astiav.NewRational(25, 1))
	avctx.SetSampleAspectRatio(astiav.NewRational(1, 1))
	avctx.SetPixelFormat(astiav.PixelFormatVaapi)

	// set hw_frames_ctx for encoder's AVCodecContext
	if err := setHwframeCtx(avctx, hwDeviceCtx); err != nil {
		log.Fatalf("Failed to set hwframe context: %v", err)
	}

	if err := avctx.Open(codec, nil); err != nil {
		log.Fatalf("Cannot open video encoder codec: %v", err)
	}

	log.Printf("Using VAAPI encoder: %s", encName)
	log.Printf("Video dimensions: %dx%d", width, height)

	for {
		swFrame := astiav.AllocFrame()
		if swFrame == nil {
			log.Fatal("Failed to allocate software frame")
		}

		// read data into software frame, and transfer them into hw frame
		swFrame.SetWidth(width)
		swFrame.SetHeight(height)
		swFrame.SetPixelFormat(astiav.PixelFormatNv12)
		if err := swFrame.AllocBuffer(0); err != nil {
			swFrame.Free()
			log.Fatalf("Failed to allocate software frame buffer: %v", err)
		}

		// Read data into software frame, just like C code fread((uint8_t*)(sw_frame->data[0]), size, 1, fin)
		ySlice := swFrame.DataSlice(0, size)
		if ySlice == nil {
			swFrame.Free()
			log.Fatal("Failed to get Y plane data slice")
		}
		n, err := fin.Read(ySlice)
		if err != nil || n <= 0 {
			swFrame.Free()
			break
		}

		// Read UV plane, just like C code fread((uint8_t*)(sw_frame->data[1]), size/2, 1, fin)
		uvSlice := swFrame.DataSlice(1, size/2)
		if uvSlice == nil {
			swFrame.Free()
			log.Fatal("Failed to get UV plane data slice")
		}
		n, err = fin.Read(uvSlice)
		if err != nil || n <= 0 {
			swFrame.Free()
			break
		}

		hwFrame := astiav.AllocFrame()
		if hwFrame == nil {
			swFrame.Free()
			log.Fatal("Failed to allocate hardware frame")
		}

		if err := hwFrame.AllocHardwareBuffer(avctx.HardwareFramesContext()); err != nil {
			swFrame.Free()
			hwFrame.Free()
			log.Fatalf("Failed to get hardware frame buffer: %v", err)
		}

		if hwFrame.HardwareFramesContext() == nil {
			swFrame.Free()
			hwFrame.Free()
			log.Fatal("Hardware frames context is nil")
		}

		if err := swFrame.TransferHardwareData(hwFrame); err != nil {
			swFrame.Free()
			hwFrame.Free()
			log.Fatalf("Error while transferring frame data to surface: %v", err)
		}

		if err := encodeWrite(avctx, hwFrame, fout); err != nil {
			swFrame.Free()
			hwFrame.Free()
			log.Fatalf("Failed to encode: %v", err)
		}

		swFrame.Free()
		hwFrame.Free()
	}

	// flush encoder
	err = encodeWrite(avctx, nil, fout)
	if errors.Is(err, astiav.ErrEof) {
		err = nil // Like C code: if (err == AVERROR_EOF) err = 0;
	}
	if err != nil {
		log.Fatalf("Failed to flush encoder: %v", err)
	}

	log.Println("VAAPI hardware encoding completed successfully")
}
