package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/asticode/go-astiav"
	"github.com/asticode/go-astikit"
)

var (
	hwDeviceCtx     *astiav.HardwareDeviceContext
	hwPixFmt        astiav.PixelFormat
	outputFile      *os.File
	c               = astikit.NewCloser()
	
	// Filter related variables
	buffersinkCtx   *astiav.BuffersinkFilterContext
	buffersrcCtx    *astiav.BuffersrcFilterContext
	filterGraph     *astiav.FilterGraph
	videoStreamIndex int
)

const filterDescr = "scale_vaapi=78:24"

func hwDecoderInit(ctx *astiav.CodecContext, hwType astiav.HardwareDeviceType) error {
	var err error
	// 指定硬件设备路径，对于VAAPI通常是 /dev/dri/renderD128
	devicePath := ""
	if hwType == astiav.HardwareDeviceTypeVAAPI {
		devicePath = "/dev/dri/renderD128"
	}
	
	hwDeviceCtx, err = astiav.CreateHardwareDeviceContext(hwType, devicePath, nil, 0)
	if err != nil && hwType == astiav.HardwareDeviceTypeVAAPI {
		// 如果默认设备失败，尝试空路径
		hwDeviceCtx, err = astiav.CreateHardwareDeviceContext(hwType, "", nil, 0)
	}
	if err != nil {
		return fmt.Errorf("failed to create specified HW device: %w", err)
	}
	
	ctx.SetHardwareDeviceContext(hwDeviceCtx)
	return nil
}

func getHwFormat(pixFmts []astiav.PixelFormat) astiav.PixelFormat {
	for _, pf := range pixFmts {
		if pf == hwPixFmt {
			return pf
		}
	}
	
	log.Println("Failed to get HW surface format")
	return astiav.PixelFormatNone
}

func initFilters(filtersDescr string, decCtx *astiav.CodecContext, inputCtx *astiav.FormatContext) error {
	// Allocate filter graph
	filterGraph = astiav.AllocFilterGraph()
	if filterGraph == nil {
		return errors.New("cannot allocate filter graph")
	}
	c.Add(filterGraph.Free)

	// Allocate inputs and outputs
	outputs := astiav.AllocFilterInOut()
	if outputs == nil {
		return errors.New("cannot allocate filter outputs")
	}
	c.Add(outputs.Free)

	inputs := astiav.AllocFilterInOut()
	if inputs == nil {
		return errors.New("cannot allocate filter inputs")
	}
	c.Add(inputs.Free)

	// Get buffer source and sink filters
	buffersrc := astiav.FindFilterByName("buffer")
	if buffersrc == nil {
		return errors.New("cannot find buffer source")
	}

	buffersink := astiav.FindFilterByName("buffersink")
	if buffersink == nil {
		return errors.New("cannot find buffer sink")
	}

	// Create buffer source context
	timeBase := inputCtx.Streams()[videoStreamIndex].TimeBase()

	var err error
	buffersrcCtx, err = filterGraph.NewBuffersrcFilterContext(buffersrc, "in")
	if err != nil {
		return fmt.Errorf("cannot create buffer source: %w", err)
	}

	// Set buffer source parameters for hardware filtering
	params := astiav.AllocBuffersrcFilterContextParameters()
	defer params.Free()
	params.SetWidth(decCtx.Width())
	params.SetHeight(decCtx.Height())
	// Use hardware pixel format for VAAPI processing
	params.SetPixelFormat(decCtx.PixelFormat())
	params.SetTimeBase(timeBase)
	params.SetSampleAspectRatio(decCtx.SampleAspectRatio())
	// Set hardware frames context for hardware filtering - 这是关键！
	if decCtx.HardwareFramesContext() != nil {
		params.SetHardwareFramesContext(decCtx.HardwareFramesContext())
	}

	if err := buffersrcCtx.SetParameters(params); err != nil {
		return fmt.Errorf("cannot set buffer source parameters: %w", err)
	}

	if err := buffersrcCtx.Initialize(nil); err != nil {
		return fmt.Errorf("cannot initialize buffer source: %w", err)
	}

	// Create buffer sink context
	buffersinkCtx, err = filterGraph.NewBuffersinkFilterContext(buffersink, "out")
	if err != nil {
		return fmt.Errorf("cannot create buffer sink: %w", err)
	}

	if err := buffersinkCtx.Initialize(); err != nil {
		return fmt.Errorf("cannot initialize buffer sink: %w", err)
	}

	// Set the endpoints for the filter graph
	outputs.SetName("in")
	outputs.SetFilterContext(buffersrcCtx.FilterContext())
	outputs.SetPadIdx(0)
	outputs.SetNext(nil)

	inputs.SetName("out")
	inputs.SetFilterContext(buffersinkCtx.FilterContext())
	inputs.SetPadIdx(0)
	inputs.SetNext(nil)

	// Parse the filter graph
	if err := filterGraph.Parse(filtersDescr, inputs, outputs); err != nil {
		return fmt.Errorf("cannot parse filter graph: %w", err)
	}

	// Configure the filter graph
	if err := filterGraph.Configure(); err != nil {
		return fmt.Errorf("cannot configure filter graph: %w", err)
	}

	return nil
}

func filterFrame(frame *astiav.Frame) error {
	// Add frame to buffer source
	if err := buffersrcCtx.AddFrame(frame, astiav.NewBuffersrcFlags(astiav.BuffersrcFlagKeepRef)); err != nil {
		return fmt.Errorf("error while feeding the filtergraph: %w", err)
	}

	// Pull filtered frames from the filtergraph
	for {
		filtFrame := astiav.AllocFrame()
		if filtFrame == nil {
			return errors.New("cannot allocate filtered frame")
		}
		defer filtFrame.Free()

		err := buffersinkCtx.GetFrame(filtFrame, astiav.NewBuffersinkFlags())
		if err != nil {
			if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("error while getting frame from filtergraph: %w", err)
		}

		// Write filtered frame to output
		if err := displayFrame(filtFrame); err != nil {
			return fmt.Errorf("error displaying frame: %w", err)
		}

		filtFrame.Unref()
	}

	return nil
}

func displayFrame(frame *astiav.Frame) error {
	var outputFrame *astiav.Frame
	
	// Check if this is a hardware frame that needs to be transferred to system memory
	if frame.HardwareFramesContext() != nil {
		// Create a software frame for transfer
		swFrame := astiav.AllocFrame()
		if swFrame == nil {
			return errors.New("cannot allocate software frame")
		}
		defer swFrame.Free()
		
		// Transfer data from hardware frame to software frame
		if err := frame.TransferHardwareData(swFrame); err != nil {
			return fmt.Errorf("error transferring hardware frame to system memory: %w", err)
		}
		
		outputFrame = swFrame
	} else {
		outputFrame = frame
	}

	// Get frame buffer size and copy to output file
	size, err := outputFrame.ImageBufferSize(1)
	if err != nil {
		return fmt.Errorf("failed to get image buffer size: %w", err)
	}

	buffer := make([]byte, size)
	_, err = outputFrame.ImageCopyToBuffer(buffer, 1)
	if err != nil {
		return fmt.Errorf("failed to copy image to buffer: %w", err)
	}

	if _, err := outputFile.Write(buffer); err != nil {
		return fmt.Errorf("failed to write buffer: %w", err)
	}

	return nil
}

var filtersInitialized = false

func decodeWrite(avctx *astiav.CodecContext, packet *astiav.Packet, inputCtx *astiav.FormatContext) error {
	err := avctx.SendPacket(packet)
	if err != nil {
		return fmt.Errorf("error during decoding: %w", err)
	}

	for {
		frame := astiav.AllocFrame()
		if frame == nil {
			return errors.New("can not alloc frame")
		}
		defer frame.Free()

		err := avctx.ReceiveFrame(frame)
		if err != nil {
			if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("error during decoding: %w", err)
		}

		// Initialize filters after first frame when hardware frames context is available
		if !filtersInitialized {
			if avctx.HardwareFramesContext() == nil {
				return fmt.Errorf("hardware frames context not available after decoding first frame")
			}
			log.Printf("Hardware frames context available, initializing filters")
			if err := initFilters(filterDescr, avctx, inputCtx); err != nil {
				return fmt.Errorf("failed to initialize filters: %w", err)
			}
			filtersInitialized = true
		}

		// Set PTS if not already set
		if frame.Pts() == astiav.NoPtsValue {
			frame.SetPts(0)
		}

		// Push the hardware frame directly into the filtergraph for hardware processing
		if err := filterFrame(frame); err != nil {
			return fmt.Errorf("error filtering frame: %w", err)
		}

		frame.Unref()
	}

	return nil
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

	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <device type> <input file> <output file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s vaapi input.mp4 output.yuv\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Filter: %s\n", filterDescr)
		os.Exit(1)
	}

	// We use an astikit.Closer to free all resources properly
	defer c.Close()

	deviceTypeName := os.Args[1]
	inputFile := os.Args[2]
	outputFileName := os.Args[3]

	// Find hardware device type
	hwType := astiav.FindHardwareDeviceTypeByName(deviceTypeName)
	if hwType == astiav.HardwareDeviceTypeNone {
		log.Fatalf("Device type %s is not supported", deviceTypeName)
	}

	// Allocate packet
	packet := astiav.AllocPacket()
	if packet == nil {
		log.Fatal("Failed to allocate AVPacket")
	}
	c.Add(packet.Free)

	// Open the input file
	inputCtx := astiav.AllocFormatContext()
	if inputCtx == nil {
		log.Fatal("Failed to allocate format context")
	}
	c.Add(inputCtx.Free)
	
	if err := inputCtx.OpenInput(inputFile, nil, nil); err != nil {
		log.Fatalf("Cannot open input file '%s': %v", inputFile, err)
	}
	c.Add(inputCtx.CloseInput)

	if err := inputCtx.FindStreamInfo(nil); err != nil {
		log.Fatalf("Cannot find input stream information: %v", err)
	}

	// Find the video stream information
	videoStreamIndex = -1
	var decoder *astiav.Codec
	var videoStream *astiav.Stream

	for i, stream := range inputCtx.Streams() {
		if stream.CodecParameters().MediaType() == astiav.MediaTypeVideo {
			videoStreamIndex = i
			videoStream = stream
			decoder = astiav.FindDecoder(stream.CodecParameters().CodecID())
			break
		}
	}

	if videoStreamIndex == -1 {
		log.Fatal("Cannot find a video stream in the input file")
	}

	if decoder == nil {
		log.Fatal("Failed to find decoder")
	}

	// Find hardware config
	hwConfigFound := false
	configs := decoder.HardwareConfigs()
	for _, config := range configs {
		if config.MethodFlags().Has(astiav.CodecHardwareConfigMethodFlagHwDeviceCtx) &&
			config.HardwareDeviceType() == hwType {
			hwPixFmt = config.PixelFormat()
			hwConfigFound = true
			break
		}
	}

	if !hwConfigFound {
		log.Fatalf("Decoder %s does not support device type %s", decoder.Name(), deviceTypeName)
	}

	// Allocate decoder context
	decoderCtx := astiav.AllocCodecContext(decoder)
	if decoderCtx == nil {
		log.Fatal("Failed to allocate decoder context")
	}
	c.Add(decoderCtx.Free)

	if err := videoStream.CodecParameters().ToCodecContext(decoderCtx); err != nil {
		log.Fatalf("Failed to copy codec parameters to decoder context: %v", err)
	}

	decoderCtx.SetPixelFormatCallback(getHwFormat)

	if err := hwDecoderInit(decoderCtx, hwType); err != nil {
		log.Fatalf("Failed to initialize hardware decoder: %v", err)
	}

	if err := decoderCtx.Open(decoder, nil); err != nil {
		log.Fatalf("Failed to open codec for stream #%d: %v", videoStreamIndex, err)
	}

	// 等待第一帧解码后再初始化过滤器，因为需要硬件帧上下文
	log.Printf("Decoder opened successfully, hardware frames context will be available after first frame")

	// Open the file to dump raw data
	var err error
	outputFile, err = os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer outputFile.Close()

	log.Printf("Using hardware decoder: %s", decoder.Name())
	log.Printf("Hardware device type: %s", deviceTypeName)
	log.Printf("Hardware pixel format: %s", hwPixFmt.String())
	log.Printf("Filter description: %s", filterDescr)

	// Actual decoding, filtering and dump the raw data
	for {
		err := inputCtx.ReadFrame(packet)
		if err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			log.Fatalf("Error reading frame: %v", err)
		}

		if videoStreamIndex == packet.StreamIndex() {
			if err := decodeWrite(decoderCtx, packet, inputCtx); err != nil {
				log.Fatalf("Error in decode_write: %v", err)
			}
		} 

		packet.Unref()
	}

	// Flush the decoder
	if err := decodeWrite(decoderCtx, nil, inputCtx); err != nil && !errors.Is(err, astiav.ErrEof) {
		log.Fatalf("Error flushing decoder: %v", err)
	}

	log.Println("Hardware decoding and filtering completed successfully")
}