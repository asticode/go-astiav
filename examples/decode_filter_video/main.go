package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/asticode/go-astiav"
)

const (
	filterDescr = "scale=78:24,transpose=cclock"
)

var (
	input = flag.String("i", "", "input file")
)

var (
	formatContext    *astiav.FormatContext
	codecContext     *astiav.CodecContext
	buffersinkCtx    *astiav.BuffersinkFilterContext
	buffersrcCtx     *astiav.BuffersrcFilterContext
	filterGraph      *astiav.FilterGraph
	videoStreamIndex = -1
	lastPts          int64 = astiav.NoPtsValue
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

	if *input == "" {
		log.Fatal("Input file is required")
	}

	// Open input file
	if err := openInputFile(*input); err != nil {
		log.Fatal("Error opening input file:", err)
	}
	defer cleanup()

	// Initialize filters
	if err := initFilters(filterDescr); err != nil {
		log.Fatal("Error initializing filters:", err)
	}

	// Main processing loop
	frame := astiav.AllocFrame()
	if frame == nil {
		log.Fatal("Could not allocate frame")
	}
	defer frame.Free()

	filtFrame := astiav.AllocFrame()
	if filtFrame == nil {
		log.Fatal("Could not allocate filter frame")
	}
	defer filtFrame.Free()

	packet := astiav.AllocPacket()
	if packet == nil {
		log.Fatal("Could not allocate packet")
	}
	defer packet.Free()

	// Read frames from the file
	for {
		if err := formatContext.ReadFrame(packet); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			log.Fatal("Error reading frame:", err)
		}

		if packet.StreamIndex() == videoStreamIndex {
			if err := codecContext.SendPacket(packet); err != nil {
				log.Fatal("Error sending packet:", err)
			}

			for {
				err := codecContext.ReceiveFrame(frame)
				if err != nil {
					if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
						break
					}
					log.Fatal("Error receiving frame:", err)
				}

				frame.SetPts(frame.Pts())

				// Push the decoded frame into the filtergraph
				if err := buffersrcCtx.AddFrame(frame, astiav.NewBuffersrcFlags()); err != nil {
					log.Fatal("Error while feeding the filtergraph:", err)
				}

				// Pull filtered frames from the filtergraph
				for {
					err := buffersinkCtx.GetFrame(filtFrame, astiav.NewBuffersinkFlags())
					if err != nil {
						if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
							break
						}
						log.Fatal("Error getting filtered frame:", err)
					}

					displayFrame(filtFrame, buffersinkCtx.TimeBase())
					filtFrame.Unref()
				}
				frame.Unref()
			}
		}
		packet.Unref()
	}

	log.Println("Video decode and filtering completed successfully")
}

func openInputFile(filename string) error {
	var err error

	// Open input file
	if formatContext = astiav.AllocFormatContext(); formatContext == nil {
		return fmt.Errorf("could not allocate format context")
	}
	if err = formatContext.OpenInput(filename, nil, nil); err != nil {
		return fmt.Errorf("cannot open input file: %w", err)
	}

	// Find stream information
	if err := formatContext.FindStreamInfo(nil); err != nil {
		return fmt.Errorf("cannot find stream information: %w", err)
	}

	// Find the video stream
	for i, stream := range formatContext.Streams() {
		if stream.CodecParameters().MediaType() == astiav.MediaTypeVideo {
			videoStreamIndex = i
			break
		}
	}

	if videoStreamIndex == -1 {
		return fmt.Errorf("cannot find a video stream in the input file")
	}

	// Find decoder for the video stream
	codecParameters := formatContext.Streams()[videoStreamIndex].CodecParameters()
	codec := astiav.FindDecoder(codecParameters.CodecID())
	if codec == nil {
		return fmt.Errorf("failed to find codec")
	}

	// Create decoding context
	codecContext = astiav.AllocCodecContext(codec)
	if codecContext == nil {
		return fmt.Errorf("failed to allocate codec context")
	}

	// Copy codec parameters to context
	if err := codecContext.FromCodecParameters(codecParameters); err != nil {
		return fmt.Errorf("failed to copy codec parameters to context: %w", err)
	}

	// Initialize the decoder
	if err := codecContext.Open(codec, nil); err != nil {
		return fmt.Errorf("cannot open video decoder: %w", err)
	}

	return nil
}

func initFilters(filtersDescr string) error {
	var err error

	// Create filter graph
	filterGraph = astiav.AllocFilterGraph()
	if filterGraph == nil {
		return fmt.Errorf("could not allocate filter graph")
	}

	// Get filters
	buffersrc := astiav.FindFilterByName("buffer")
	if buffersrc == nil {
		return fmt.Errorf("could not find the buffer filter")
	}

	buffersink := astiav.FindFilterByName("buffersink")
	if buffersink == nil {
		return fmt.Errorf("could not find the buffersink filter")
	}

	// Create filter in/out
	outputs := astiav.AllocFilterInOut()
	if outputs == nil {
		return fmt.Errorf("could not allocate filter outputs")
	}
	defer outputs.Free()

	inputs := astiav.AllocFilterInOut()
	if inputs == nil {
		return fmt.Errorf("could not allocate filter inputs")
	}
	defer inputs.Free()

	// Buffer video source: 按照C代码方式设置参数
	timeBase := formatContext.Streams()[videoStreamIndex].TimeBase()

	// 使用CreateFilter方法创建buffersrc
	buffersrcCtx, err = filterGraph.NewBuffersrcFilterContext(buffersrc, "in")
	if err != nil {
		return fmt.Errorf("cannot create buffer source: %w", err)
	}

	// 设置buffersrc参数
	buffersrcParams := astiav.AllocBuffersrcFilterContextParameters()
	defer buffersrcParams.Free()
	buffersrcParams.SetWidth(codecContext.Width())
	buffersrcParams.SetHeight(codecContext.Height())
	buffersrcParams.SetPixelFormat(codecContext.PixelFormat())
	buffersrcParams.SetTimeBase(timeBase)
	buffersrcParams.SetSampleAspectRatio(codecContext.SampleAspectRatio())

	if err := buffersrcCtx.SetParameters(buffersrcParams); err != nil {
		return fmt.Errorf("could not set buffersrc parameters: %w", err)
	}

	if err := buffersrcCtx.Initialize(nil); err != nil {
		return fmt.Errorf("could not initialize buffersrc context: %w", err)
	}

	// Buffer video sink: 按照C代码方式创建
	buffersinkCtx, err = filterGraph.NewBuffersinkFilterContext(buffersink, "out")
	if err != nil {
		return fmt.Errorf("cannot create buffer sink: %w", err)
	}

	// 按照C代码设置buffersink选项 - 设置为gray8像素格式
	if err := buffersinkCtx.FilterContext().SetOption("pixel_formats", "gray8"); err != nil {
		return fmt.Errorf("cannot set output pixel format: %w", err)
	}

	if err := buffersinkCtx.Initialize(); err != nil {
		return fmt.Errorf("could not initialize buffersink context: %w", err)
	}

	// 按照C代码设置filter graph的端点
	outputs.SetName("in")
	outputs.SetFilterContext(buffersrcCtx.FilterContext())
	outputs.SetPadIdx(0)
	outputs.SetNext(nil)

	inputs.SetName("out")
	inputs.SetFilterContext(buffersinkCtx.FilterContext())
	inputs.SetPadIdx(0)
	inputs.SetNext(nil)

	// 按照C代码解析filter graph
	if err := filterGraph.Parse(filtersDescr, inputs, outputs); err != nil {
		return fmt.Errorf("could not parse filter graph: %w", err)
	}

	// 按照C代码配置filter graph
	if err := filterGraph.Configure(); err != nil {
		return fmt.Errorf("could not configure filter graph: %w", err)
	}

	return nil
}

func displayFrame(frame *astiav.Frame, timeBase astiav.Rational) {
	if frame.Pts() != astiav.NoPtsValue {
		if lastPts != astiav.NoPtsValue {
			// Sleep roughly the right amount of time
			delay := astiav.RescaleQ(frame.Pts()-lastPts, timeBase, astiav.TimeBaseQ)
			if delay > 0 && delay < 1000000 {
				time.Sleep(time.Duration(delay) * time.Microsecond)
			}
		}
		lastPts = frame.Pts()
	}

	// Print frame info
	fmt.Printf("Frame: %dx%d, format: %s, pts: %d\n",
		frame.Width(), frame.Height(), frame.PixelFormat().Name(), frame.Pts())

	// Print first few pixels as ASCII art (simplified)
	if frame.PixelFormat() == astiav.PixelFormatGray8 {
		data := frame.DataSlice(0, frame.Width()*frame.Height())
		if data != nil && len(data) > 0 {
			fmt.Print("ASCII representation (first line): ")
			width := frame.Width()
			if width > 78 {
				width = 78
			}
			for i := 0; i < width; i++ {
				if i < len(data) {
					pixel := data[i]
					if pixel > 128 {
						fmt.Print("#")
					} else if pixel > 64 {
						fmt.Print("*")
					} else if pixel > 32 {
						fmt.Print(".")
					} else {
						fmt.Print(" ")
					}
				}
			}
			fmt.Println()
		}
	}
}

func cleanup() {
	if codecContext != nil {
		codecContext.Free()
	}
	if formatContext != nil {
		formatContext.CloseInput()
	}
	if filterGraph != nil {
		filterGraph.Free()
	}
}