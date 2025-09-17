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

const (
	filterDescr = "aresample=8000,aformat=sample_fmts=s16:channel_layouts=mono"
	player      = "ffplay -f s16le -ar 8000 -ac 1 -"
)

var (
	input = flag.String("i", "", "the input audio file path")
)

var (
	formatContext    *astiav.FormatContext
	codecContext     *astiav.CodecContext
	buffersinkCtx    *astiav.BuffersinkFilterContext
	buffersrcCtx     *astiav.BuffersrcFilterContext
	filterGraph      *astiav.FilterGraph
	audioStreamIndex = -1
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
		fmt.Fprintf(os.Stderr, "Usage: %s -i <input_file> | %s\n", os.Args[0], player)
		os.Exit(1)
	}

	// Allocate packet and frames
	packet := astiav.AllocPacket()
	if packet == nil {
		log.Fatal("Could not allocate packet")
	}
	defer packet.Free()

	frame := astiav.AllocFrame()
	if frame == nil {
		log.Fatal("Could not allocate frame")
	}
	defer frame.Free()

	filtFrame := astiav.AllocFrame()
	if filtFrame == nil {
		log.Fatal("Could not allocate filtered frame")
	}
	defer filtFrame.Free()

	if err := openInputFile(*input); err != nil {
		log.Fatal("Error opening input file:", err)
	}
	defer cleanup()

	if err := initFilters(filterDescr); err != nil {
		log.Fatal("Error initializing filters:", err)
	}

	// Read all packets
	for {
		err := formatContext.ReadFrame(packet)
		if err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			log.Fatal("Error reading frame:", err)
		}

		if packet.StreamIndex() == audioStreamIndex {
			if err := codecContext.SendPacket(packet); err != nil {
				log.Printf("Error while sending a packet to the decoder: %v", err)
				break
			}

			for {
				err := codecContext.ReceiveFrame(frame)
				if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
					break
				} else if err != nil {
					log.Fatal("Error while receiving a frame from the decoder:", err)
				}

				// Push the audio data from decoded frame into the filtergraph
				if err := buffersrcCtx.AddFrame(frame, astiav.NewBuffersrcFlags(astiav.BuffersrcFlagKeepRef)); err != nil {
					log.Printf("Error while feeding the audio filtergraph: %v", err)
					break
				}

				// Pull filtered audio from the filtergraph
				for {
					err := buffersinkCtx.GetFrame(filtFrame, astiav.NewBuffersinkFlags())
					if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
						break
					}
					if err != nil {
						log.Fatal("Error getting filtered frame:", err)
					}
					printFrame(filtFrame)
					filtFrame.Unref()
				}
				frame.Unref()
			}
		}
		packet.Unref()
	}

	// Signal EOF to the filtergraph
	if err := buffersrcCtx.AddFrame(nil, astiav.NewBuffersrcFlags()); err != nil {
		log.Printf("Error while closing the filtergraph: %v", err)
	} else {
		// Pull remaining frames from the filtergraph
		for {
			err := buffersinkCtx.GetFrame(filtFrame, astiav.NewBuffersinkFlags())
			if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
				break
			}
			if err != nil {
				log.Fatal("Error getting remaining filtered frame:", err)
			}
			printFrame(filtFrame)
			filtFrame.Unref()
		}
	}

	log.Println("Audio decode and filtering completed successfully")
}

func openInputFile(filename string) error {
	// Allocate format context
	formatContext = astiav.AllocFormatContext()
	if formatContext == nil {
		return fmt.Errorf("could not allocate format context")
	}

	// Open input file
	if err := formatContext.OpenInput(filename, nil, nil); err != nil {
		return fmt.Errorf("could not open input file '%s': %w", filename, err)
	}

	// Find stream info
	if err := formatContext.FindStreamInfo(nil); err != nil {
		return fmt.Errorf("could not find stream information: %w", err)
	}

	// Find the first audio stream
	for i, stream := range formatContext.Streams() {
		if stream.CodecParameters().MediaType() == astiav.MediaTypeAudio {
			audioStreamIndex = i
			break
		}
	}

	if audioStreamIndex < 0 {
		return fmt.Errorf("could not find audio stream in input file '%s'", filename)
	}

	// Get a pointer to the codec context for the audio stream
	stream := formatContext.Streams()[audioStreamIndex]
	
	// Find decoder for the stream
	codec := astiav.FindDecoder(stream.CodecParameters().CodecID())
	if codec == nil {
		return fmt.Errorf("failed to find codec")
	}

	// Allocate codec context
	codecContext = astiav.AllocCodecContext(codec)
	if codecContext == nil {
		return fmt.Errorf("failed to allocate codec context")
	}

	// Copy codec parameters from input stream to output codec context
	if err := codecContext.FromCodecParameters(stream.CodecParameters()); err != nil {
		return fmt.Errorf("failed to copy codec parameters to decoder context: %w", err)
	}

	// Open codec
	if err := codecContext.Open(codec, nil); err != nil {
		return fmt.Errorf("failed to open codec: %w", err)
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
	abuffersrc := astiav.FindFilterByName("abuffer")
	if abuffersrc == nil {
		return fmt.Errorf("could not find abuffer filter")
	}

	abuffersink := astiav.FindFilterByName("abuffersink")
	if abuffersink == nil {
		return fmt.Errorf("could not find abuffersink filter")
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

	// Buffer audio source: 按照C代码方式设置参数
	timeBase := formatContext.Streams()[audioStreamIndex].TimeBase()
	channelLayout := codecContext.ChannelLayout()
	if !channelLayout.Valid() {
		// 暂时使用默认的stereo layout
		channelLayout = astiav.ChannelLayoutStereo
	}

	// 使用avfilter_graph_create_filter的等价方式
	buffersrcCtx, err = filterGraph.NewBuffersrcFilterContext(abuffersrc, "in")
	if err != nil {
		return fmt.Errorf("cannot create audio buffer source: %w", err)
	}

	// 设置buffersrc参数
	buffersrcParams := astiav.AllocBuffersrcFilterContextParameters()
	defer buffersrcParams.Free()
	buffersrcParams.SetSampleRate(codecContext.SampleRate())
	buffersrcParams.SetSampleFormat(codecContext.SampleFormat())
	buffersrcParams.SetChannelLayout(channelLayout)
	buffersrcParams.SetTimeBase(timeBase)

	if err := buffersrcCtx.SetParameters(buffersrcParams); err != nil {
		return fmt.Errorf("could not set buffersrc parameters: %w", err)
	}

	if err := buffersrcCtx.Initialize(nil); err != nil {
		return fmt.Errorf("could not initialize buffersrc context: %w", err)
	}

	// Buffer audio sink: 按照C代码方式创建
	buffersinkCtx, err = filterGraph.NewBuffersinkFilterContext(abuffersink, "out")
	if err != nil {
		return fmt.Errorf("cannot create audio buffer sink: %w", err)
	}

	// 按照FFmpeg 8.0的方式，让filter chain自动协商格式
	// 不设置具体的sample_formats和channel_layouts，由aformat filter处理

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

func printFrame(frame *astiav.Frame) {
	nbSamples := frame.NbSamples()
	channels := frame.ChannelLayout().Channels()
	
	fmt.Printf("nb_samples:%d pts:%d\n", 
		nbSamples, frame.Pts())

	// Print sample data (first few samples for demonstration)
	data := frame.DataSlice(0, nbSamples*channels*2) // 2 bytes per s16 sample
	if data != nil && len(data) >= 4 {
		// Print first sample as s16le
		sample := int16(data[0]) | int16(data[1])<<8
		fmt.Printf("First sample: %d\n", sample)
	}
}

func cleanup() {
	if filterGraph != nil {
		filterGraph.Free()
	}
	if codecContext != nil {
		codecContext.Free()
	}
	if formatContext != nil {
		formatContext.CloseInput()
		formatContext.Free()
	}
}