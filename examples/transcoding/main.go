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
	input  = flag.String("i", "", "the input path")
	output = flag.String("o", "", "the output path")
)

var (
	c                   = astikit.NewCloser()
	inputFormatContext  *astiav.FormatContext
	outputFormatContext *astiav.FormatContext
	streams             = make(map[int]*stream) // Indexed by input stream index
)

type stream struct {
	buffersinkContext *astiav.FilterContext
	buffersrcContext  *astiav.FilterContext
	decCodec          *astiav.Codec
	decCodecContext   *astiav.CodecContext
	decFrame          *astiav.Frame
	encCodec          *astiav.Codec
	encCodecContext   *astiav.CodecContext
	encPkt            *astiav.Packet
	filterFrame       *astiav.Frame
	filterGraph       *astiav.FilterGraph
	inputStream       *astiav.Stream
	outputStream      *astiav.Stream
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
	if *input == "" || *output == "" {
		log.Println("Usage: <binary path> -i <input path> -o <output path>")
		return
	}

	// We use an astikit.Closer to free all resources properly
	defer c.Close()

	// Open input file
	if err := openInputFile(); err != nil {
		log.Fatal(fmt.Errorf("main: opening input file failed: %w", err))
	}

	// Open output file
	if err := openOutputFile(); err != nil {
		log.Fatal(fmt.Errorf("main: opening output file failed: %w", err))
	}

	// Init filters
	if err := initFilters(); err != nil {
		log.Fatal(fmt.Errorf("main: initializing filters failed: %w", err))
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

		// Get stream
		s, ok := streams[pkt.StreamIndex()]
		if !ok {
			continue
		}

		// Update packet
		pkt.RescaleTs(s.inputStream.TimeBase(), s.decCodecContext.TimeBase())

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

			// Filter, encode and write frame
			if err := filterEncodeWriteFrame(s.decFrame, s); err != nil {
				log.Fatal(fmt.Errorf("main: filtering, encoding and writing frame failed: %w", err))
			}
		}
	}

	// Loop through streams
	for _, s := range streams {
		// Flush filter
		if err := filterEncodeWriteFrame(nil, s); err != nil {
			log.Fatal(fmt.Errorf("main: filtering, encoding and writing frame failed: %w", err))
		}

		// Flush encoder
		if err := encodeWriteFrame(nil, s); err != nil {
			log.Fatal(fmt.Errorf("main: encoding and writing frame failed: %w", err))
		}
	}

	// Write trailer
	if err := outputFormatContext.WriteTrailer(); err != nil {
		log.Fatal(fmt.Errorf("main: writing trailer failed: %w", err))
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
		// Only process audio or video
		if is.CodecParameters().MediaType() != astiav.MediaTypeAudio &&
			is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
			continue
		}

		// Create stream
		s := &stream{inputStream: is}

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

		// Set framerate
		if is.CodecParameters().MediaType() == astiav.MediaTypeVideo {
			s.decCodecContext.SetFramerate(inputFormatContext.GuessFrameRate(is, nil))
		}

		// Open codec context
		if err = s.decCodecContext.Open(s.decCodec, nil); err != nil {
			err = fmt.Errorf("main: opening codec context failed: %w", err)
			return
		}

		// Alloc frame
		s.decFrame = astiav.AllocFrame()
		c.Add(s.decFrame.Free)

		// Store stream
		streams[is.Index()] = s
	}
	return
}

func openOutputFile() (err error) {
	// Alloc output format context
	if outputFormatContext, err = astiav.AllocOutputFormatContext(nil, "", *output); err != nil {
		err = fmt.Errorf("main: allocating output format context failed: %w", err)
		return
	} else if outputFormatContext == nil {
		err = errors.New("main: output format context is nil")
		return
	}
	c.Add(outputFormatContext.Free)

	// Loop through streams
	for _, is := range inputFormatContext.Streams() {
		// Get stream
		s, ok := streams[is.Index()]
		if !ok {
			continue
		}

		// Create output stream
		if s.outputStream = outputFormatContext.NewStream(nil); s.outputStream == nil {
			err = errors.New("main: output stream is nil")
			return
		}

		// Get codec id
		codecID := astiav.CodecIDMpeg4
		if s.decCodecContext.MediaType() == astiav.MediaTypeAudio {
			codecID = astiav.CodecIDAac
		}

		// Find encoder
		if s.encCodec = astiav.FindEncoder(codecID); s.encCodec == nil {
			err = errors.New("main: codec is nil")
			return
		}

		// Alloc codec context
		if s.encCodecContext = astiav.AllocCodecContext(s.encCodec); s.encCodecContext == nil {
			err = errors.New("main: codec context is nil")
			return
		}
		c.Add(s.encCodecContext.Free)

		// Update codec context
		if s.decCodecContext.MediaType() == astiav.MediaTypeAudio {
			if v := s.encCodec.ChannelLayouts(); len(v) > 0 {
				s.encCodecContext.SetChannelLayout(v[0])
			} else {
				s.encCodecContext.SetChannelLayout(s.decCodecContext.ChannelLayout())
			}
			s.encCodecContext.SetSampleRate(s.decCodecContext.SampleRate())
			if v := s.encCodec.SampleFormats(); len(v) > 0 {
				s.encCodecContext.SetSampleFormat(v[0])
			} else {
				s.encCodecContext.SetSampleFormat(s.decCodecContext.SampleFormat())
			}
			s.encCodecContext.SetTimeBase(s.decCodecContext.TimeBase())
		} else {
			s.encCodecContext.SetHeight(s.decCodecContext.Height())
			if v := s.encCodec.PixelFormats(); len(v) > 0 {
				s.encCodecContext.SetPixelFormat(v[0])
			} else {
				s.encCodecContext.SetPixelFormat(s.decCodecContext.PixelFormat())
			}
			s.encCodecContext.SetSampleAspectRatio(s.decCodecContext.SampleAspectRatio())
			s.encCodecContext.SetTimeBase(s.decCodecContext.TimeBase())
			s.encCodecContext.SetWidth(s.decCodecContext.Width())
		}

		// Update flags
		if s.decCodecContext.Flags().Has(astiav.CodecContextFlagGlobalHeader) {
			s.encCodecContext.SetFlags(s.encCodecContext.Flags().Add(astiav.CodecContextFlagGlobalHeader))
		}

		// Open codec context
		if err = s.encCodecContext.Open(s.encCodec, nil); err != nil {
			err = fmt.Errorf("main: opening codec context failed: %w", err)
			return
		}

		// Update codec parameters
		if err = s.outputStream.CodecParameters().FromCodecContext(s.encCodecContext); err != nil {
			err = fmt.Errorf("main: updating codec parameters failed: %w", err)
			return
		}

		// Update stream
		s.outputStream.SetTimeBase(s.encCodecContext.TimeBase())
	}

	// If this is a file, we need to use an io context
	if !outputFormatContext.OutputFormat().Flags().Has(astiav.IOFormatFlagNofile) {
		// Open io context
		var ioContext *astiav.IOContext
		if ioContext, err = astiav.OpenIOContext(*output, astiav.NewIOContextFlags(astiav.IOContextFlagWrite)); err != nil {
			err = fmt.Errorf("main: opening io context failed: %w", err)
			return
		}
		c.AddWithError(ioContext.Close)

		// Update output format context
		outputFormatContext.SetPb(ioContext)
	}

	// Write header
	if err = outputFormatContext.WriteHeader(nil); err != nil {
		err = fmt.Errorf("main: writing header failed: %w", err)
		return
	}
	return
}

func initFilters() (err error) {
	// Loop through output streams
	for _, s := range streams {
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

		// Switch on media type
		var args astiav.FilterArgs
		var buffersrc, buffersink *astiav.Filter
		var content string
		switch s.decCodecContext.MediaType() {
		case astiav.MediaTypeAudio:
			args = astiav.FilterArgs{
				"channel_layout": s.decCodecContext.ChannelLayout().String(),
				"sample_fmt":     s.decCodecContext.SampleFormat().Name(),
				"sample_rate":    strconv.Itoa(s.decCodecContext.SampleRate()),
				"time_base":      s.decCodecContext.TimeBase().String(),
			}
			buffersrc = astiav.FindFilterByName("abuffer")
			buffersink = astiav.FindFilterByName("abuffersink")
			content = fmt.Sprintf("aformat=sample_fmts=%s:channel_layouts=%s", s.encCodecContext.SampleFormat().Name(), s.encCodecContext.ChannelLayout().String())
		default:
			args = astiav.FilterArgs{
				"pix_fmt":      strconv.Itoa(int(s.decCodecContext.PixelFormat())),
				"pixel_aspect": s.decCodecContext.SampleAspectRatio().String(),
				"time_base":    s.decCodecContext.TimeBase().String(),
				"video_size":   strconv.Itoa(s.decCodecContext.Width()) + "x" + strconv.Itoa(s.decCodecContext.Height()),
			}
			buffersrc = astiav.FindFilterByName("buffer")
			buffersink = astiav.FindFilterByName("buffersink")
			content = fmt.Sprintf("format=pix_fmts=%s", s.encCodecContext.PixelFormat().Name())
		}

		// Check filters
		if buffersrc == nil {
			err = errors.New("main: buffersrc is nil")
			return
		}
		if buffersink == nil {
			err = errors.New("main: buffersink is nil")
			return
		}

		// Create filter contexts
		if s.buffersrcContext, err = s.filterGraph.NewFilterContext(buffersrc, "in", args); err != nil {
			err = fmt.Errorf("main: creating buffersrc context failed: %w", err)
			return
		}
		if s.buffersinkContext, err = s.filterGraph.NewFilterContext(buffersink, "out", nil); err != nil {
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
		if err = s.filterGraph.Parse(content, inputs, outputs); err != nil {
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

		// Alloc packet
		s.encPkt = astiav.AllocPacket()
		c.Add(s.encPkt.Free)
	}
	return
}

func filterEncodeWriteFrame(f *astiav.Frame, s *stream) (err error) {
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

		// Reset picture type
		s.filterFrame.SetPictureType(astiav.PictureTypeNone)

		// Encode and write frame
		if err = encodeWriteFrame(s.filterFrame, s); err != nil {
			err = fmt.Errorf("main: encoding and writing frame failed: %w", err)
			return
		}
	}
	return
}

func encodeWriteFrame(f *astiav.Frame, s *stream) (err error) {
	// Unref packet
	s.encPkt.Unref()

	// Send frame
	if err = s.encCodecContext.SendFrame(f); err != nil {
		err = fmt.Errorf("main: sending frame failed: %w", err)
		return
	}

	// Loop
	for {
		// Receive packet
		if err = s.encCodecContext.ReceivePacket(s.encPkt); err != nil {
			if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
				err = nil
				break
			}
			err = fmt.Errorf("main: receiving packet failed: %w", err)
			return
		}

		// Update pkt
		s.encPkt.SetStreamIndex(s.outputStream.Index())
		s.encPkt.RescaleTs(s.encCodecContext.TimeBase(), s.outputStream.TimeBase())

		// Write frame
		if err = outputFormatContext.WriteInterleavedFrame(s.encPkt); err != nil {
			err = fmt.Errorf("main: writing frame failed: %w", err)
			return
		}
	}
	return
}
