package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
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
	astiav.SetLogLevel(astiav.LogLevelInfo) // 减少日志输出
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

	log.Println("Step 1: Opening input file...")
	// Open input file
	if err := openInputFile(); err != nil {
		log.Fatal(fmt.Errorf("main: opening input file failed: %w", err))
	}
	log.Println("Step 1: Input file opened successfully")

	log.Println("Step 2: Initializing filter...")
	// Init filter
	if err := initFilter(); err != nil {
		log.Fatal(fmt.Errorf("main: initializing filter failed: %w", err))
	}
	log.Println("Step 2: Filter initialized successfully")

	// Allocate packet
	pkt := astiav.AllocPacket()
	c.Add(pkt.Free)

	log.Println("Step 3: Starting packet processing...")
	frameCount := 0
	packetCount := 0
	// Loop through packets
	for {
		// Read frame
		if err := inputFormatContext.ReadFrame(pkt); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				log.Println("Reached end of file")
				break
			}
			log.Fatal(fmt.Errorf("main: reading frame failed: %w", err))
		}

		// Make sure to unreference the packet
		defer pkt.Unref()

		// Invalid stream
		if pkt.StreamIndex() != s.inputStream.Index() {
			continue
		}

		packetCount++
		// Send packet
		log.Printf("Packet %d: size=%d, pts=%d, dts=%d", packetCount, pkt.Size(), pkt.Pts(), pkt.Dts())
		if err := s.decCodecContext.SendPacket(pkt); err != nil {
			log.Fatal(fmt.Errorf("main: sending packet failed: %w", err))
		}

		// Try to receive frames after each packet
		for {
			// Receive frame
			if err := s.decCodecContext.ReceiveFrame(s.decFrame); err != nil {
				if errors.Is(err, astiav.ErrEof) {
					log.Println("Decoder EOF reached")
					break
				}
				if errors.Is(err, astiav.ErrEagain) {
					log.Println("Decoder needs more input (EAGAIN)")
					break
				}
				log.Fatal(fmt.Errorf("main: receiving frame failed: %w", err))
			}

			// Make sure to unreference the frame
			defer s.decFrame.Unref()

			frameCount++
			log.Printf("Processing frame %d: %dx%d, format: %s", frameCount, s.decFrame.Width(), s.decFrame.Height(), s.decFrame.PixelFormat())

			// Filter frame
			if err := filterFrame(s.decFrame, s); err != nil {
				log.Fatal(fmt.Errorf("main: filtering frame failed: %w", err))
			}
		}
	}

	log.Println("Step 3.5: Flushing decoder...")
	// Send NULL packet to flush decoder
	if err := s.decCodecContext.SendPacket(nil); err != nil {
		log.Fatal(fmt.Errorf("main: flushing decoder failed: %w", err))
	}

	// Receive remaining frames from decoder
	for {
		if err := s.decCodecContext.ReceiveFrame(s.decFrame); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				log.Println("Decoder flush complete")
				break
			}
			if errors.Is(err, astiav.ErrEagain) {
				log.Println("No more frames in decoder")
				break
			}
			log.Fatal(fmt.Errorf("main: receiving frame during flush failed: %w", err))
		}

		// Make sure to unreference the frame
		defer s.decFrame.Unref()

		frameCount++
		log.Printf("Flushed frame %d: %dx%d, format: %s", frameCount, s.decFrame.Width(), s.decFrame.Height(), s.decFrame.PixelFormat())

		// Filter frame
		if err := filterFrame(s.decFrame, s); err != nil {
			log.Fatal(fmt.Errorf("main: filtering flushed frame failed: %w", err))
		}
	}

	log.Println("Step 4: Flushing filter...")
	// Flush filter
	if err := filterFrame(nil, s); err != nil {
		log.Fatal(fmt.Errorf("main: filtering frame failed: %w", err))
	}

	// Success
	log.Printf("Success! Processed %d frames", frameCount)
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

		log.Printf("Video stream found: %dx%d, codec: %s, pixel format: %s", 
			s.decCodecContext.Width(), s.decCodecContext.Height(), 
			s.decCodec.Name(), s.decCodecContext.PixelFormat())

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
	// Alloc filter graph
	if s.filterGraph = astiav.AllocFilterGraph(); s.filterGraph == nil {
		err = errors.New("main: graph is nil")
		return
	}
	c.Add(s.filterGraph.Free)

	// Get filters
	buffersrcFilter := astiav.FindFilterByName("buffer")
	if buffersrcFilter == nil {
		err = errors.New("main: buffersrc is nil")
		return
	}
	buffersinkFilter := astiav.FindFilterByName("buffersink")
	if buffersinkFilter == nil {
		err = errors.New("main: buffersink is nil")
		return
	}

	// Create buffer source
	if s.buffersrcContext, err = s.filterGraph.NewBuffersrcFilterContext(buffersrcFilter, "in"); err != nil {
		err = fmt.Errorf("main: creating buffersrc context failed: %w", err)
		return
	}

	// Set buffersrc parameters - 基于C测试的成功经验
	buffersrcParams := astiav.AllocBuffersrcFilterContextParameters()
	defer buffersrcParams.Free()
	buffersrcParams.SetWidth(s.decCodecContext.Width())
	buffersrcParams.SetHeight(s.decCodecContext.Height())
	buffersrcParams.SetPixelFormat(s.decCodecContext.PixelFormat())
	buffersrcParams.SetTimeBase(s.inputStream.TimeBase())
	buffersrcParams.SetSampleAspectRatio(s.decCodecContext.SampleAspectRatio())

	log.Printf("Buffersrc params: %dx%d, format: %s, timebase: %s", 
		s.decCodecContext.Width(), s.decCodecContext.Height(), 
		s.decCodecContext.PixelFormat(), s.inputStream.TimeBase())

	if err = s.buffersrcContext.SetParameters(buffersrcParams); err != nil {
		err = fmt.Errorf("main: setting buffersrc parameters failed: %w", err)
		return
	}

	// Initialize buffersrc context
	if err = s.buffersrcContext.Initialize(nil); err != nil {
		err = fmt.Errorf("main: initializing buffersrc context failed: %w", err)
		return
	}

	// Create buffer sink
	if s.buffersinkContext, err = s.filterGraph.NewBuffersinkFilterContext(buffersinkFilter, "out"); err != nil {
		err = fmt.Errorf("main: creating buffersink context failed: %w", err)
		return
	}

	// Initialize buffersink - 不设置像素格式，让其自动协商
	if err = s.buffersinkContext.Initialize(); err != nil {
		err = fmt.Errorf("main: initializing buffersink context failed: %w", err)
		return
	}

	// Alloc outputs
	outputs := astiav.AllocFilterInOut()
	if outputs == nil {
		err = errors.New("main: outputs is nil")
		return
	}
	c.Add(outputs.Free)

	// Set outputs
	outputs.SetName("in")
	outputs.SetFilterContext(s.buffersrcContext.FilterContext())
	outputs.SetPadIdx(0)
	outputs.SetNext(nil)

	// Alloc inputs
	inputs := astiav.AllocFilterInOut()
	if inputs == nil {
		err = errors.New("main: inputs is nil")
		return
	}
	c.Add(inputs.Free)

	// Set inputs
	inputs.SetName("out")
	inputs.SetFilterContext(s.buffersinkContext.FilterContext())
	inputs.SetPadIdx(0)
	inputs.SetNext(nil)

	// Parse - 使用更复杂的滤镜链
	filterDesc := "scale=160:120,hflip,vflip"
	log.Printf("Parsing filter: %s", filterDesc)
	if err = s.filterGraph.Parse(filterDesc, inputs, outputs); err != nil {
		err = fmt.Errorf("main: parsing filter failed: %w", err)
		return
	}

	// Configure
	log.Println("Configuring filter graph...")
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
			log.Printf("Filtered frame: %dx%d, format: %s", s.filterFrame.Width(), s.filterFrame.Height(), s.filterFrame.PixelFormat())
			return false, nil
		}(); err != nil {
			return err
		} else if stop {
			break
		}

	}
	return
}