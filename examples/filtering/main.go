package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/asticode/go-astikit"
	"log"
	"strings"

	"github.com/asticode/go-astikit"
	"github.com/peakedshout/go-astiav"
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
	buffersinkContext *astiav.BuffersinkFilterContext
	buffersrcContext  *astiav.BuffersrcFilterContext
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

	// Allocate packet
	pkt := astiav.AllocPacket()
	c.Add(pkt.Free)

	// Loop through packets
	for {
		// We use a closure to ease unreferencing the packet
		if stop := func() bool {
			// Read frame
			if err := inputFormatContext.ReadFrame(pkt); err != nil {
				if errors.Is(err, astiav.ErrEof) {
					return true
				}
				log.Fatal(fmt.Errorf("main: reading frame failed: %w", err))
			}

			// Make sure to unreference the packet
			defer pkt.Unref()

			// Invalid stream
			if pkt.StreamIndex() != s.inputStream.Index() {
				return false
			}

			// Send packet
			if err := s.decCodecContext.SendPacket(pkt); err != nil {
				log.Fatal(fmt.Errorf("main: sending packet failed: %w", err))
			}

			// Loop
			for {
				// We use a closure to ease unreferencing the frame
				if stop := func() bool {
					// Receive frame
					if err := s.decCodecContext.ReceiveFrame(s.decFrame); err != nil {
						if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
							return true
						}
						log.Fatal(fmt.Errorf("main: receiving frame failed: %w", err))
					}

					// Make sure to unreference the frame
					defer s.decFrame.Unref()

					// Filter frame
					if err := filterFrame(s.decFrame, s); err != nil {
						log.Fatal(fmt.Errorf("main: filtering frame failed: %w", err))
					}
					return false
				}(); stop {
					break
				}
			}
			return false
		}(); stop {
			break
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
	// Allocate input format context
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

		// Allocate codec context
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

		// Allocate frame
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
	// Allocate graph
	if s.filterGraph = astiav.AllocFilterGraph(); s.filterGraph == nil {
		err = errors.New("main: graph is nil")
		return
	}
	c.Add(s.filterGraph.Free)

	// Allocate outputs
	outputs := astiav.AllocFilterInOut()
	if outputs == nil {
		err = errors.New("main: outputs is nil")
		return
	}
	c.Add(outputs.Free)

	// Allocate inputs
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
	if s.buffersrcContext, err = s.filterGraph.NewBuffersrcFilterContext(buffersrc, "in"); err != nil {
		err = fmt.Errorf("main: creating buffersrc context failed: %w", err)
		return
	}
	if s.buffersinkContext, err = s.filterGraph.NewBuffersinkFilterContext(buffersink, "in"); err != nil {
		err = fmt.Errorf("main: creating buffersink context failed: %w", err)
		return
	}

	// Create buffersrc context parameters
	buffersrcContextParameters := astiav.AllocBuffersrcFilterContextParameters()
	defer buffersrcContextParameters.Free()
	buffersrcContextParameters.SetHeight(s.decCodecContext.Height())
	buffersrcContextParameters.SetPixelFormat(s.decCodecContext.PixelFormat())
	buffersrcContextParameters.SetSampleAspectRatio(s.decCodecContext.SampleAspectRatio())
	buffersrcContextParameters.SetTimeBase(s.inputStream.TimeBase())
	buffersrcContextParameters.SetWidth(s.decCodecContext.Width())

	// Set buffersrc context parameters
	if err = s.buffersrcContext.SetParameters(buffersrcContextParameters); err != nil {
		err = fmt.Errorf("main: setting buffersrc context parameters failed: %w", err)
		return
	}

	// Initialize buffersrc context
	if err = s.buffersrcContext.Initialize(nil); err != nil {
		err = fmt.Errorf("main: initializing buffersrc context failed: %w", err)
		return
	}

	// Update outputs
	outputs.SetName("in")
	outputs.SetFilterContext(s.buffersrcContext.FilterContext())
	outputs.SetPadIdx(0)
	outputs.SetNext(nil)

	// Update inputs
	inputs.SetName("out")
	inputs.SetFilterContext(s.buffersinkContext.FilterContext())
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

	// Allocate frame
	s.filterFrame = astiav.AllocFrame()
	c.Add(s.filterFrame.Free)
	return
}

func filterFrame(f *astiav.Frame, s *stream) (err error) {
	// Add frame
	if err = s.buffersrcContext.AddFrame(f, astiav.NewBuffersrcFlags(astiav.BuffersrcFlagKeepRef)); err != nil {
		err = fmt.Errorf("main: adding frame failed: %w", err)
		return
	}

	// Loop
	for {
		// We use a closure to ease unreferencing the frame
		if stop, err := func() (bool, error) {
			// Get frame
			if err := s.buffersinkContext.GetFrame(s.filterFrame, astiav.NewBuffersinkFlags()); err != nil {
				if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
					return true, nil
				}
				return false, fmt.Errorf("main: getting frame failed: %w", err)
			}

			// Make sure to unrefernce the frame
			defer s.filterFrame.Unref()

			// Do something with filtered frame
			log.Printf("new filtered frame: %dx%d\n", s.filterFrame.Width(), s.filterFrame.Height())
			return false, nil
		}(); err != nil {
			return err
		} else if stop {
			break
		}

	}
	return
}
