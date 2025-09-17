package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/asticode/go-astiav"
)

const (
	InputSampleRate     = 48000
	InputFormat         = astiav.SampleFormatFltp
	VolumeVal           = 0.90
	FrameSize           = 1024
)

var (
	duration = flag.Float64("d", 1.0, "duration in seconds")
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

	if *duration <= 0 {
		log.Fatal("Invalid duration")
	}

	nbFrames := int(*duration * InputSampleRate / FrameSize)
	if nbFrames <= 0 {
		log.Fatal("Invalid duration: too short")
	}

	// Allocate the frame we will be using to store the data
	frame := astiav.AllocFrame()
	if frame == nil {
		log.Fatal("Error allocating the frame")
	}
	defer frame.Free()

	// Set up the filtergraph
	graph, src, sink, err := initFilterGraph()
	if err != nil {
		log.Fatal("Unable to init filter graph:", err)
	}
	defer graph.Free()

	// The main filtering loop
	for i := 0; i < nbFrames; i++ {
		// Get an input frame to be filtered
		if err := getInput(frame, i); err != nil {
			log.Fatal("Error generating input frame:", err)
		}

		// Send the frame to the input of the filtergraph
		if err := src.AddFrame(frame, astiav.NewBuffersrcFlags()); err != nil {
			frame.Unref()
			log.Fatal("Error submitting the frame to the filtergraph:", err)
		}

		// Get all the filtered output that is available
		for {
			err := sink.GetFrame(frame, astiav.NewBuffersinkFlags())
			if err != nil {
				if errors.Is(err, astiav.ErrEagain) {
					// Need to feed more frames in
					break
				} else if errors.Is(err, astiav.ErrEof) {
					// Nothing more to do, finish
					goto done
				} else {
					// An error occurred
					log.Fatal("Error filtering the data:", err)
				}
			}

			// Now do something with our filtered frame
			if err := processOutput(frame); err != nil {
				log.Fatal("Error processing the filtered frame:", err)
			}
			frame.Unref()
		}
	}

done:
	log.Println("Audio filtering completed successfully")
}

func initFilterGraph() (*astiav.FilterGraph, *astiav.BuffersrcFilterContext, *astiav.BuffersinkFilterContext, error) {
	// Create a new filtergraph, which will contain all the filters
	graph := astiav.AllocFilterGraph()
	if graph == nil {
		return nil, nil, nil, fmt.Errorf("unable to create filter graph")
	}

	// Create the abuffer filter; it will be used for feeding the data into the graph
	abuffer := astiav.FindFilterByName("abuffer")
	if abuffer == nil {
		return nil, nil, nil, fmt.Errorf("could not find the abuffer filter")
	}

	abufferCtx, err := graph.NewBuffersrcFilterContext(abuffer, "src")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not allocate the abuffer instance: %w", err)
	}

	// Set the filter options through the AVOptions API
	channelLayout := astiav.ChannelLayout5Point0
	chLayoutStr := channelLayout.String()
	
	if err := abufferCtx.FilterContext().SetOption("channel_layout", chLayoutStr); err != nil {
		return nil, nil, nil, fmt.Errorf("could not set channel_layout: %w", err)
	}
	
	if err := abufferCtx.FilterContext().SetOption("sample_fmt", InputFormat.Name()); err != nil {
		return nil, nil, nil, fmt.Errorf("could not set sample_fmt: %w", err)
	}
	
	timeBase := astiav.NewRational(1, InputSampleRate)
	if err := abufferCtx.FilterContext().SetOption("time_base", timeBase.String()); err != nil {
		return nil, nil, nil, fmt.Errorf("could not set time_base: %w", err)
	}
	
	if err := abufferCtx.FilterContext().SetOption("sample_rate", fmt.Sprintf("%d", InputSampleRate)); err != nil {
		return nil, nil, nil, fmt.Errorf("could not set sample_rate: %w", err)
	}

	// Initialize the filter
	if err := abufferCtx.Initialize(nil); err != nil {
		return nil, nil, nil, fmt.Errorf("could not initialize the abuffer filter: %w", err)
	}

	// Create volume filter
	volume := astiav.FindFilterByName("volume")
	if volume == nil {
		return nil, nil, nil, fmt.Errorf("could not find the volume filter")
	}

	volumeCtx, err := graph.NewFilterContext(volume, "volume")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not allocate the volume instance: %w", err)
	}

	// Set volume options using dictionary
	dict := astiav.NewDictionary()
	defer dict.Free()
	dict.Set("volume", fmt.Sprintf("%.2f", VolumeVal), astiav.NewDictionaryFlags())
	
	if err := volumeCtx.Initialize(dict); err != nil {
		return nil, nil, nil, fmt.Errorf("could not initialize the volume filter: %w", err)
	}

	// Create the aformat filter; it ensures that the output is of the format we want
	aformat := astiav.FindFilterByName("aformat")
	if aformat == nil {
		return nil, nil, nil, fmt.Errorf("could not find the aformat filter")
	}

	aformatCtx, err := graph.NewFilterContext(aformat, "aformat")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not allocate the aformat instance: %w", err)
	}

	// Set aformat options
	aformatDict := astiav.NewDictionary()
	defer aformatDict.Free()
	aformatDict.Set("sample_fmts", InputFormat.Name(), astiav.NewDictionaryFlags())
	aformatDict.Set("sample_rates", fmt.Sprintf("%d", InputSampleRate), astiav.NewDictionaryFlags())
	aformatDict.Set("channel_layouts", chLayoutStr, astiav.NewDictionaryFlags())
	
	if err := aformatCtx.Initialize(aformatDict); err != nil {
		return nil, nil, nil, fmt.Errorf("could not initialize the aformat filter: %w", err)
	}

	// Create the abuffersink filter; it will be used for getting the filtered data out of the graph
	abuffersink := astiav.FindFilterByName("abuffersink")
	if abuffersink == nil {
		return nil, nil, nil, fmt.Errorf("could not find the abuffersink filter")
	}

	abuffersinkCtx, err := graph.NewBuffersinkFilterContext(abuffersink, "sink")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not allocate the abuffersink instance: %w", err)
	}

	if err := abuffersinkCtx.Initialize(); err != nil {
		return nil, nil, nil, fmt.Errorf("could not initialize the abuffersink filter: %w", err)
	}

	// Connect the filters manually like in FFmpeg C code
	// abuffer -> volume -> aformat -> abuffersink
	
	// Connect abuffer to volume
	if err := graph.Link(abufferCtx.FilterContext(), 0, volumeCtx, 0); err != nil {
		return nil, nil, nil, fmt.Errorf("could not link abuffer to volume: %w", err)
	}

	// Connect volume to aformat
	if err := graph.Link(volumeCtx, 0, aformatCtx, 0); err != nil {
		return nil, nil, nil, fmt.Errorf("could not link volume to aformat: %w", err)
	}

	// Connect aformat to abuffersink
	if err := graph.Link(aformatCtx, 0, abuffersinkCtx.FilterContext(), 0); err != nil {
		return nil, nil, nil, fmt.Errorf("could not link aformat to abuffersink: %w", err)
	}

	// Configure the graph
	if err := graph.Configure(); err != nil {
		return nil, nil, nil, fmt.Errorf("could not configure the filter graph: %w", err)
	}

	return graph, abufferCtx, abuffersinkCtx, nil
}

func processOutput(frame *astiav.Frame) error {
	planar := frame.SampleFormat().IsPlanar()
	channels := frame.ChannelLayout().Channels()
	planes := 1
	if planar {
		planes = channels
	}
	bps := frame.SampleFormat().BytesPerSample()
	planeSize := bps * frame.NbSamples()
	if !planar {
		planeSize *= channels
	}

	for i := 0; i < planes; i++ {
		// Get plane data
		data := frame.DataSlice(i, planeSize)
		if data == nil {
			continue
		}

		// Print plane info (instead of MD5)
		fmt.Printf("plane %d: %d samples, %d bytes\n", i, frame.NbSamples(), len(data))
	}
	fmt.Println()

	return nil
}

// Construct a frame of audio data to be filtered;
// this simple example just synthesizes a sine wave.
func getInput(frame *astiav.Frame, frameNum int) error {
	// Set up the frame properties and allocate the buffer for the data
	frame.SetSampleRate(InputSampleRate)
	frame.SetSampleFormat(InputFormat)
	frame.SetChannelLayout(astiav.ChannelLayout5Point0)
	frame.SetNbSamples(FrameSize)
	frame.SetPts(int64(frameNum * FrameSize))

	if err := frame.AllocBuffer(0); err != nil {
		return err
	}

	// Fill the data for each channel
	for i := 0; i < 5; i++ { // 5.0 channel layout has 5 channels
		data := frame.DataSlice(i, FrameSize*4) // 4 bytes per float sample
		if data == nil {
			continue
		}

		// Convert byte slice to float32 slice
		floatData := make([]float32, FrameSize)
		for j := 0; j < FrameSize; j++ {
			floatData[j] = float32(math.Sin(2 * math.Pi * float64(frameNum+j) * float64(i+1) / FrameSize))
		}

		// Copy float data to frame
		for j := 0; j < FrameSize && j*4+3 < len(data); j++ {
			bits := math.Float32bits(floatData[j])
			data[j*4] = byte(bits)
			data[j*4+1] = byte(bits >> 8)
			data[j*4+2] = byte(bits >> 16)
			data[j*4+3] = byte(bits >> 24)
		}
	}

	return nil
}