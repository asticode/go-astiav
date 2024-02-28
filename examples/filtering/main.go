package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/asticode/go-astiav"
	"github.com/asticode/go-astikit"
)

var (
	input = flag.String("i", "", "the input path")
)

var (
	c                  = astikit.NewCloser()
	inputFormatContext *astiav.FormatContext
	s                  *stream
)

type stream struct {
	buffersinkContext *astiav.FilterContext
	buffersrcContext  *astiav.FilterContext
	decCodec          *astiav.Codec
	decCodecContext   *astiav.CodecContext
	decFrame          *astiav.Frame
	filterFrame       *astiav.Frame
	filterGraph       *astiav.FilterGraph
	inputStream       *astiav.Stream
	lastPts           int64
}

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

	// We use an astikit.Closer to free all resources properly
	defer c.Close()

	// Open input file
	if err := openInputFile(); err != nil {
		log.Fatal(fmt.Errorf("main: opening input file failed: %w", err))
	}

	// Init filter
	if err := initFilter(); err != nil {
		log.Fatal(fmt.Errorf("main: initializing filter failed: %w", err))
	}

	// Alloc packet
	pkt := astiav.AllocPacket()
	c.Add(pkt.Free)

	// Loop through packets
	for {
		// Read frame
		if err := inputFormatContext.ReadFrame(pkt); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			log.Fatal(fmt.Errorf("main: reading frame failed: %w", err))
		}

		// Invalid stream
		if pkt.StreamIndex() != s.inputStream.Index() {
			continue
		}

		// Send packet
		if err := s.decCodecContext.SendPacket(pkt); err != nil {
			log.Fatal(fmt.Errorf("main: sending packet failed: %w", err))
		}

		// Loop
		for {
			// Receive frame
			if err := s.decCodecContext.ReceiveFrame(s.decFrame); err != nil {
				if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
					break
				}
				log.Fatal(fmt.Errorf("main: receiving frame failed: %w", err))
			}

			// Filter frame
			if err := filterFrame(s.decFrame, s); err != nil {
				log.Fatal(fmt.Errorf("main: filtering frame failed: %w", err))
			}
		}
	}

	// Flush filter
	if err := filterFrame(nil, s); err != nil {
		log.Fatal(fmt.Errorf("main: filtering frame failed: %w", err))
	}

	// Success
	log.Println("success")
}

func openInputFile() (err error) {
	// Alloc input format context
	if inputFormatContext = astiav.AllocFormatContext(); inputFormatContext == nil {
		err = errors.New("main: input format context is nil")
		return
	}
	c.Add(inputFormatContext.Free)

	// Open input
	if err = inputFormatContext.OpenInput(*input, nil, nil); err != nil {
		err = fmt.Errorf("main: opening input failed: %w", err)
		return
	}
	c.Add(inputFormatContext.CloseInput)

	// Find stream info
	if err = inputFormatContext.FindStreamInfo(nil); err != nil {
		err = fmt.Errorf("main: finding stream info failed: %w", err)
		return
	}

	// Loop through streams
	for _, is := range inputFormatContext.Streams() {
		// Only process video
		if is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
			continue
		}

		// Create stream
		s = &stream{
			inputStream: is,
			lastPts:     astiav.NoPtsValue,
		}

		// Find decoder
		if s.decCodec = astiav.FindDecoder(is.CodecParameters().CodecID()); s.decCodec == nil {
			err = errors.New("main: codec is nil")
			return
		}

		// Alloc codec context
		if s.decCodecContext = astiav.AllocCodecContext(s.decCodec); s.decCodecContext == nil {
			err = errors.New("main: codec context is nil")
			return
		}
		c.Add(s.decCodecContext.Free)

		// Update codec context
		if err = is.CodecParameters().ToCodecContext(s.decCodecContext); err != nil {
			err = fmt.Errorf("main: updating codec context failed: %w", err)
			return
		}

		// Open codec context
		if err = s.decCodecContext.Open(s.decCodec, nil); err != nil {
			err = fmt.Errorf("main: opening codec context failed: %w", err)
			return
		}

		// Alloc frame
		s.decFrame = astiav.AllocFrame()
		c.Add(s.decFrame.Free)

		break
	}

	// No video stream
	if s == nil {
		err = errors.New("main: no video stream")
		return
	}
	return
}

func initFilter() (err error) {
	// Alloc graph
	if s.filterGraph = astiav.AllocFilterGraph(); s.filterGraph == nil {
		err = errors.New("main: graph is nil")
		return
	}
	c.Add(s.filterGraph.Free)

	// Alloc outputs
	outputs := astiav.AllocFilterInOut()
	if outputs == nil {
		err = errors.New("main: outputs is nil")
		return
	}
	c.Add(outputs.Free)

	// Alloc inputs
	inputs := astiav.AllocFilterInOut()
	if inputs == nil {
		err = errors.New("main: inputs is nil")
		return
	}
	c.Add(inputs.Free)

	// Create buffersrc
	buffersrc := astiav.FindFilterByName("buffer")
	if buffersrc == nil {
		err = errors.New("main: buffersrc is nil")
		return
	}

	// Create buffersink
	buffersink := astiav.FindFilterByName("buffersink")
	if buffersink == nil {
		err = errors.New("main: buffersink is nil")
		return
	}

	// Create filter contexts
	if s.buffersrcContext, err = s.filterGraph.NewFilterContext(buffersrc, "in", astiav.FilterArgs{
		"pix_fmt":      strconv.Itoa(int(s.decCodecContext.PixelFormat())),
		"pixel_aspect": s.decCodecContext.SampleAspectRatio().String(),
		"time_base":    s.inputStream.TimeBase().String(),
		"video_size":   strconv.Itoa(s.decCodecContext.Width()) + "x" + strconv.Itoa(s.decCodecContext.Height()),
	}); err != nil {
		err = fmt.Errorf("main: creating buffersrc context failed: %w", err)
		return
	}
	if s.buffersinkContext, err = s.filterGraph.NewFilterContext(buffersink, "in", nil); err != nil {
		err = fmt.Errorf("main: creating buffersink context failed: %w", err)
		return
	}

	// Update outputs
	outputs.SetName("in")
	outputs.SetFilterContext(s.buffersrcContext)
	outputs.SetPadIdx(0)
	outputs.SetNext(nil)

	// Update inputs
	inputs.SetName("out")
	inputs.SetFilterContext(s.buffersinkContext)
	inputs.SetPadIdx(0)
	inputs.SetNext(nil)

	// Parse
	if err = s.filterGraph.Parse("transpose=cclock", inputs, outputs); err != nil {
		err = fmt.Errorf("main: parsing filter failed: %w", err)
		return
	}

	// Configure
	if err = s.filterGraph.Configure(); err != nil {
		err = fmt.Errorf("main: configuring filter failed: %w", err)
		return
	}

	// Alloc frame
	s.filterFrame = astiav.AllocFrame()
	c.Add(s.filterFrame.Free)
	return
}

func filterFrame(f *astiav.Frame, s *stream) (err error) {
	// Add frame
	if err = s.buffersrcContext.BuffersrcAddFrame(f, astiav.NewBuffersrcFlags(astiav.BuffersrcFlagKeepRef)); err != nil {
		err = fmt.Errorf("main: adding frame failed: %w", err)
		return
	}

	// Loop
	for {
		// Unref frame
		s.filterFrame.Unref()

		// Get frame
		if err = s.buffersinkContext.BuffersinkGetFrame(s.filterFrame, astiav.NewBuffersinkFlags()); err != nil {
			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
				err = nil
				break
			}
			err = fmt.Errorf("main: getting frame failed: %w", err)
			return
		}

		// Do something with filtered frame
		log.Printf("new filtered frame: %dx%d\n", s.filterFrame.Width(), s.filterFrame.Height())
	}
	return
}
