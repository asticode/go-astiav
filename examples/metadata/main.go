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
	input = flag.String("i", "", "the input path")
)

func main() {
	// Handle ffmpeg logs
	astiav.SetLogLevel(astiav.LogLevelDebug)
	astiav.SetLogCallback(func(c astiav.Classer, l astiav.LogLevel, fmt, msg string) {
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
	if *input == "" {
		log.Println("Usage: <binary path> -i <input path>")
		return
	}

	// Allocate input format context
	inputFormatContext := astiav.AllocFormatContext()
	if inputFormatContext == nil {
		log.Println(errors.New("main: input format context is nil"))
		return
	}
	defer inputFormatContext.Free()

	// Open input
	if err := inputFormatContext.OpenInput(*input, nil, nil); err != nil {
		log.Println(fmt.Errorf("main: opening input failed: %w", err))
		return
	}
	defer inputFormatContext.CloseInput()

	// Find stream info
	if err := inputFormatContext.FindStreamInfo(nil); err != nil {
		log.Println(fmt.Errorf("main: finding stream info failed: %w", err))
		return
	}

	// Iterate over all streams in the container.
	for _, stream := range inputFormatContext.Streams() {
		params := stream.CodecParameters()

		// Look up the decoder to get human-readable codec name and description.
		var name, description string
		if codec := astiav.FindDecoder(params.CodecID()); codec != nil {
			name, description = codec.Name(), codec.LongName()
		}

		fmt.Printf("stream index: %d\n", stream.Index())
		fmt.Printf("\ttype: %s\n", params.MediaType())
		fmt.Printf("\tname: %s\n", name)
		fmt.Printf("\tdescription: %s\n", description)
		fmt.Printf("\tprofile: %s\n", params.ProfileName())

		fmt.Printf("\ttime base: %s\n", stream.TimeBase())
		fmt.Printf("\tstart time: %f\n", float64(stream.StartTime())/float64(astiav.TimeBaseQ.Den()))
		fmt.Printf("\tduration (stream timebase): %d\n", stream.Duration())
		fmt.Printf("\tduration (seconds): %f\n", float64(stream.Duration())/float64(stream.TimeBase().Den()))
		fmt.Printf("\tcodec id: %d\n", params.CodecID())

		flags := stream.DispositionFlags()
		fmt.Printf("\tdefault: %t\n", flags.Has(astiav.DispositionFlagDefault))
		fmt.Printf("\tforced: %t\n", flags.Has(astiav.DispositionFlagDefault))

		// Dispatch based on media type and build the appropriate stream struct.
		switch params.MediaType() {
		case astiav.MediaTypeVideo:
			sar := params.SampleAspectRatio()
			dar := sar.Mul(astiav.NewRational(params.Width(), params.Height()))

			fmt.Printf("\theight: %d\n", params.Height())
			fmt.Printf("\twidth: %d\n", params.Width())
			fmt.Printf("\tsample aspect ratio: %d:%d\n", sar.Num(), sar.Den())
			fmt.Printf("\tdisplay aspect ratio: %d:%d\n", dar.Num(), dar.Den())
			fmt.Printf("\treal frame rate: %f\n", stream.RFrameRate().Float64())
			fmt.Printf("\tavg frame rate: %f\n", stream.AvgFrameRate().Float64())
		case astiav.MediaTypeAudio:
			fmt.Printf("\tchannels: %d\n", params.ChannelLayout().Channels())
			fmt.Printf("\tchannel layout: %s\n", params.ChannelLayout().String())

			fmt.Printf("\tsample rate: %d\n", params.SampleRate())
			fmt.Printf("\tbit rate: %d\n", params.BitRate())
		case astiav.MediaTypeSubtitle:
		}

		printDictionary(stream.Metadata())
		fmt.Printf("\n")
	}

	// Iterate over all chapters in the container.
	for i, chapter := range inputFormatContext.Chapters() {
		fmt.Printf("chapter index: %d\n", i)
		fmt.Printf("\tid: %d\n", chapter.ID())
		fmt.Printf("\tstart time: %f\n", float64(chapter.Start())/float64(chapter.TimeBase().Den()))
		fmt.Printf("\tend time: %f\n", float64(chapter.End())/float64(chapter.TimeBase().Den()))
		printDictionary(chapter.Metadata())
		fmt.Printf("\n")
	}

	// Print the container format information.
	fmt.Printf("format:\n")
	fmt.Printf("\tname: %s\n", inputFormatContext.InputFormat().Name())
	fmt.Printf("\tdescription: %s\n", inputFormatContext.InputFormat().LongName())
	fmt.Printf("\tbit rate: %d\n", inputFormatContext.BitRate())
	fmt.Printf("\tstart time: %f\n", float64(inputFormatContext.StartTime())/float64(astiav.TimeBaseQ.Den()))
	fmt.Printf("\tduration (format timebase): %d\n", inputFormatContext.Duration())
	fmt.Printf("\tduration (seconds): %f\n", float64(inputFormatContext.Duration())/float64(astiav.TimeBaseQ.Den()))
	printDictionary(inputFormatContext.Metadata())
}

func printDictionary(dictionary *astiav.Dictionary) {
	if dictionary == nil {
		return
	}

	fmt.Printf("\tmetadata:\n")
	var entry *astiav.DictionaryEntry
	flags := astiav.NewDictionaryFlags(astiav.DictionaryFlagIgnoreSuffix)

	for {
		entry = dictionary.Get("", entry, flags) // flags = 0 is standard
		if entry == nil {
			return
		}
		fmt.Printf("\t\t%s=%s\n", entry.Key(), entry.Value())
	}
}
