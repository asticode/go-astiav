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
	input  = flag.String("i", "", "the input path")
	output = flag.String("o", "", "the output path")
)

type FilteringContext struct {
	buffersinkCtx *astiav.BuffersinkFilterContext
	buffersrcCtx  *astiav.BuffersrcFilterContext
	filterGraph   *astiav.FilterGraph
	encPkt        *astiav.Packet
	filteredFrame *astiav.Frame
}

type StreamContext struct {
	decCtx   *astiav.CodecContext
	encCtx   *astiav.CodecContext
	decFrame *astiav.Frame
}

var (
	ifmtCtx    *astiav.FormatContext
	ofmtCtx    *astiav.FormatContext
	filterCtx  []*FilteringContext
	streamCtx  []*StreamContext
	c          = astikit.NewCloser()
)

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
	if err := openInputFile(*input); err != nil {
		log.Fatal(fmt.Errorf("opening input file failed: %w", err))
	}

	// Open output file
	if err := openOutputFile(*output); err != nil {
		log.Fatal(fmt.Errorf("opening output file failed: %w", err))
	}

	// Initialize filters
	if err := initFilters(); err != nil {
		log.Fatal(fmt.Errorf("initializing filters failed: %w", err))
	}

	// Process packets
	if err := processPackets(); err != nil {
		log.Fatal(fmt.Errorf("processing packets failed: %w", err))
	}

	log.Println("Transcoding completed successfully")
}

func openInputFile(filename string) error {
	// Allocate input format context
	ifmtCtx = astiav.AllocFormatContext()
	if ifmtCtx == nil {
		return errors.New("failed to allocate input format context")
	}
	c.Add(ifmtCtx.Free)

	// Open input
	if err := ifmtCtx.OpenInput(filename, nil, nil); err != nil {
		return fmt.Errorf("opening input failed: %w", err)
	}
	c.Add(ifmtCtx.CloseInput)

	// Find stream info
	if err := ifmtCtx.FindStreamInfo(nil); err != nil {
		return fmt.Errorf("finding stream info failed: %w", err)
	}

	// Initialize stream contexts
	streamCtx = make([]*StreamContext, len(ifmtCtx.Streams()))

	// Process each stream
	for i, stream := range ifmtCtx.Streams() {
		codecPar := stream.CodecParameters()
		
		// Find decoder
		dec := astiav.FindDecoder(codecPar.CodecID())
		if dec == nil {
			return fmt.Errorf("failed to find decoder for stream #%d", i)
		}

		// Allocate decoder context
		decCtx := astiav.AllocCodecContext(dec)
		if decCtx == nil {
			return fmt.Errorf("failed to allocate decoder context for stream #%d", i)
		}
		c.Add(decCtx.Free)

		// Copy codec parameters to decoder context
		if err := codecPar.ToCodecContext(decCtx); err != nil {
			return fmt.Errorf("failed to copy decoder parameters for stream #%d: %w", i, err)
		}

		// Set packet timebase (like in C code: codec_ctx->pkt_timebase = stream->time_base)
		decCtx.SetPktTimebase(stream.TimeBase())

		// Open decoder for video and audio streams
		if codecPar.MediaType() == astiav.MediaTypeVideo || codecPar.MediaType() == astiav.MediaTypeAudio {
			if codecPar.MediaType() == astiav.MediaTypeVideo {
				// Set frame rate from stream
				if stream.AvgFrameRate().Num() > 0 {
					decCtx.SetFramerate(stream.AvgFrameRate())
				} else if stream.RFrameRate().Num() > 0 {
					decCtx.SetFramerate(stream.RFrameRate())
				}
			}

			// Open decoder
			if err := decCtx.Open(dec, nil); err != nil {
				return fmt.Errorf("failed to open decoder for stream #%d: %w", i, err)
			}
		}

		// Allocate decode frame
		decFrame := astiav.AllocFrame()
		if decFrame == nil {
			return fmt.Errorf("failed to allocate decode frame for stream #%d", i)
		}
		c.Add(decFrame.Free)

		streamCtx[i] = &StreamContext{
			decCtx:   decCtx,
			decFrame: decFrame,
		}
	}

	// Log input format info
	log.Printf("Input format: %s, duration: %d", ifmtCtx.InputFormat().Name(), ifmtCtx.Duration())
	return nil
}

func openOutputFile(filename string) error {
	// Allocate output format context
	var err error
	ofmtCtx, err = astiav.AllocOutputFormatContext(nil, "", filename)
	if err != nil {
		return fmt.Errorf("failed to allocate output format context: %w", err)
	}
	c.Add(ofmtCtx.Free)

	// Process each stream
	for i, inStream := range ifmtCtx.Streams() {
		// Create output stream
		outStream := ofmtCtx.NewStream(nil)
		if outStream == nil {
			return fmt.Errorf("failed to allocate output stream #%d", i)
		}

		decCtx := streamCtx[i].decCtx
		codecPar := inStream.CodecParameters()

		if codecPar.MediaType() == astiav.MediaTypeVideo || codecPar.MediaType() == astiav.MediaTypeAudio {
			// Find encoder (use same codec as decoder for simplicity)
			enc := astiav.FindEncoder(codecPar.CodecID())
			if enc == nil {
				return fmt.Errorf("necessary encoder not found for stream #%d", i)
			}

			// Allocate encoder context
			encCtx := astiav.AllocCodecContext(enc)
			if encCtx == nil {
				return fmt.Errorf("failed to allocate encoder context for stream #%d", i)
			}
			c.Add(encCtx.Free)

			// Configure encoder based on media type
			if codecPar.MediaType() == astiav.MediaTypeVideo {
				encCtx.SetHeight(decCtx.Height())
				encCtx.SetWidth(decCtx.Width())
				encCtx.SetSampleAspectRatio(decCtx.SampleAspectRatio())
				
				// Use decoder's pixel format or a safe default
				if decCtx.PixelFormat() != astiav.PixelFormatNone {
					encCtx.SetPixelFormat(decCtx.PixelFormat())
				} else {
					encCtx.SetPixelFormat(astiav.PixelFormatYuv420P)
				}
				
				// Set time base
				if decCtx.Framerate().Num() > 0 {
					encCtx.SetTimeBase(decCtx.Framerate().Invert())
				} else {
					encCtx.SetTimeBase(astiav.NewRational(1, 25))
				}
			} else { // Audio
				encCtx.SetSampleRate(decCtx.SampleRate())
				encCtx.SetChannelLayout(decCtx.ChannelLayout())
				
				// Use decoder's sample format or a safe default
				if decCtx.SampleFormat() != astiav.SampleFormatNone {
					encCtx.SetSampleFormat(decCtx.SampleFormat())
				} else {
					encCtx.SetSampleFormat(astiav.SampleFormatFltp)
				}
				
				encCtx.SetTimeBase(astiav.NewRational(1, encCtx.SampleRate()))
			}

			// Set global header flag if needed
			if ofmtCtx.OutputFormat().Flags().Has(astiav.IOFormatFlagGlobalheader) {
				encCtx.SetFlags(encCtx.Flags().Add(astiav.CodecContextFlagGlobalHeader))
			}

			// Open encoder
			if err := encCtx.Open(enc, nil); err != nil {
				return fmt.Errorf("cannot open encoder for stream #%d: %w", i, err)
			}

			// Copy encoder parameters to output stream
			if err := encCtx.ToCodecParameters(outStream.CodecParameters()); err != nil {
				return fmt.Errorf("failed to copy encoder parameters to output stream #%d: %w", i, err)
			}

			outStream.SetTimeBase(encCtx.TimeBase())
			streamCtx[i].encCtx = encCtx
		} else if codecPar.MediaType() == astiav.MediaTypeUnknown {
			return fmt.Errorf("elementary stream #%d is of unknown type, cannot proceed", i)
		} else {
			// Copy parameters for other streams (subtitles, etc.)
			if err := codecPar.Copy(outStream.CodecParameters()); err != nil {
				return fmt.Errorf("copying parameters for stream #%d failed: %w", i, err)
			}
			outStream.SetTimeBase(inStream.TimeBase())
		}
	}

	// Log output format info
	log.Printf("Output format: %s", ofmtCtx.OutputFormat().Name())

	// Open output file
	if !ofmtCtx.OutputFormat().Flags().Has(astiav.IOFormatFlagNofile) {
		ioCtx, err := astiav.OpenIOContext(filename, astiav.NewIOContextFlags(astiav.IOContextFlagWrite), nil, nil)
		if err != nil {
			return fmt.Errorf("could not open output file: %w", err)
		}
		c.AddWithError(ioCtx.Close)
		ofmtCtx.SetPb(ioCtx)
	}

	// Write header
	if err := ofmtCtx.WriteHeader(nil); err != nil {
		return fmt.Errorf("error occurred when opening output file: %w", err)
	}

	return nil
}

func initFilters() error {
	filterCtx = make([]*FilteringContext, len(ifmtCtx.Streams()))

	for i, stream := range ifmtCtx.Streams() {
		codecPar := stream.CodecParameters()
		
		// Only initialize filters for video and audio streams that need encoding
		if codecPar.MediaType() != astiav.MediaTypeVideo && codecPar.MediaType() != astiav.MediaTypeAudio {
			continue
		}

		if streamCtx[i].encCtx == nil {
			continue
		}

		filterCtx[i] = &FilteringContext{}
		
		// Allocate filter graph
		filterGraph := astiav.AllocFilterGraph()
		if filterGraph == nil {
			return fmt.Errorf("failed to allocate filter graph for stream #%d", i)
		}
		c.Add(filterGraph.Free)

		// Create filter spec based on media type
		var filterSpec string
		if codecPar.MediaType() == astiav.MediaTypeVideo {
			filterSpec = "null" // Pass-through filter for video
		} else {
			filterSpec = "anull" // Pass-through filter for audio
		}

		if err := initFilter(filterCtx[i], streamCtx[i].decCtx, streamCtx[i].encCtx, filterSpec, filterGraph); err != nil {
			return fmt.Errorf("failed to initialize filter for stream #%d: %w", i, err)
		}

		// Allocate packet and frame for filtering
		filterCtx[i].encPkt = astiav.AllocPacket()
		c.Add(filterCtx[i].encPkt.Free)
		
		filterCtx[i].filteredFrame = astiav.AllocFrame()
		c.Add(filterCtx[i].filteredFrame.Free)
	}

	return nil
}

func initFilter(fctx *FilteringContext, decCtx, encCtx *astiav.CodecContext, filterSpec string, filterGraph *astiav.FilterGraph) error {
	var buffersrc, buffersink *astiav.Filter
	var err error

	// Get appropriate source and sink filters
	if decCtx.MediaType() == astiav.MediaTypeVideo {
		buffersrc = astiav.FindFilterByName("buffer")
		buffersink = astiav.FindFilterByName("buffersink")
	} else {
		buffersrc = astiav.FindFilterByName("abuffer")
		buffersink = astiav.FindFilterByName("abuffersink")
	}

	if buffersrc == nil || buffersink == nil {
		return errors.New("filtering source or sink element not found")
	}

	// Create buffer source
	fctx.buffersrcCtx, err = filterGraph.NewBuffersrcFilterContext(buffersrc, "in")
	if err != nil {
		return fmt.Errorf("cannot create buffer source: %w", err)
	}

	// Set buffer source parameters
	params := astiav.AllocBuffersrcFilterContextParameters()
	defer params.Free()

	if decCtx.MediaType() == astiav.MediaTypeVideo {
		params.SetWidth(decCtx.Width())
		params.SetHeight(decCtx.Height())
		params.SetPixelFormat(decCtx.PixelFormat())
		params.SetTimeBase(decCtx.PktTimebase())
		params.SetSampleAspectRatio(decCtx.SampleAspectRatio())
	} else {
		params.SetChannelLayout(decCtx.ChannelLayout())
		params.SetSampleFormat(decCtx.SampleFormat())
		params.SetSampleRate(decCtx.SampleRate())
		params.SetTimeBase(decCtx.PktTimebase())
	}

	if err := fctx.buffersrcCtx.SetParameters(params); err != nil {
		return fmt.Errorf("cannot set buffer source parameters: %w", err)
	}

	if err := fctx.buffersrcCtx.Initialize(nil); err != nil {
		return fmt.Errorf("cannot initialize buffer source: %w", err)
	}

	// Create buffer sink
	fctx.buffersinkCtx, err = filterGraph.NewBuffersinkFilterContext(buffersink, "out")
	if err != nil {
		return fmt.Errorf("cannot create buffer sink: %w", err)
	}

	if err := fctx.buffersinkCtx.Initialize(); err != nil {
		return fmt.Errorf("cannot initialize buffer sink: %w", err)
	}

	// Set up filter graph endpoints
	outputs := astiav.AllocFilterInOut()
	defer outputs.Free()
	outputs.SetName("in")
	outputs.SetFilterContext(fctx.buffersrcCtx.FilterContext())
	outputs.SetPadIdx(0)
	outputs.SetNext(nil)

	inputs := astiav.AllocFilterInOut()
	defer inputs.Free()
	inputs.SetName("out")
	inputs.SetFilterContext(fctx.buffersinkCtx.FilterContext())
	inputs.SetPadIdx(0)
	inputs.SetNext(nil)

	// Parse filter graph
	if err := filterGraph.Parse(filterSpec, inputs, outputs); err != nil {
		return fmt.Errorf("cannot parse filter graph: %w", err)
	}

	// Configure filter graph
	if err := filterGraph.Configure(); err != nil {
		return fmt.Errorf("cannot configure filter graph: %w", err)
	}

	fctx.filterGraph = filterGraph
	return nil
}

func processPackets() error {
	pkt := astiav.AllocPacket()
	defer pkt.Free()

	// Process all packets
	for {
		if err := ifmtCtx.ReadFrame(pkt); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("reading frame failed: %w", err)
		}

		streamIndex := pkt.StreamIndex()
		if streamIndex >= len(streamCtx) {
			pkt.Unref()
			continue
		}

		if err := decodeAndEncode(pkt, streamIndex); err != nil {
			pkt.Unref()
			return fmt.Errorf("decode and encode failed: %w", err)
		}

		pkt.Unref()
	}

	// Flush decoders and encoders
	for i := range streamCtx {
		if err := flushDecoderEncoder(i); err != nil {
			return fmt.Errorf("flushing decoder/encoder failed: %w", err)
		}
	}

	// Write trailer
	if err := ofmtCtx.WriteTrailer(); err != nil {
		return fmt.Errorf("writing trailer failed: %w", err)
	}

	return nil
}

func decodeAndEncode(pkt *astiav.Packet, streamIndex int) error {
	sctx := streamCtx[streamIndex]
	
	// Send packet to decoder
	if err := sctx.decCtx.SendPacket(pkt); err != nil {
		return fmt.Errorf("sending packet to decoder failed: %w", err)
	}

	// Receive frames from decoder
	for {
		if err := sctx.decCtx.ReceiveFrame(sctx.decFrame); err != nil {
			if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("receiving frame from decoder failed: %w", err)
		}

		// Use frame's original pts or set a reasonable value
		if sctx.decFrame.Pts() == astiav.NoPtsValue {
			sctx.decFrame.SetPts(0)
		}

		if err := filterEncodeWriteFrame(sctx.decFrame, streamIndex); err != nil {
			return fmt.Errorf("filter encode write frame failed: %w", err)
		}

		sctx.decFrame.Unref()
	}

	return nil
}

func filterEncodeWriteFrame(frame *astiav.Frame, streamIndex int) error {
	fctx := filterCtx[streamIndex]
	if fctx == nil {
		// Stream doesn't need filtering/encoding, just copy
		return nil
	}

	// Send frame to filter
	if err := fctx.buffersrcCtx.AddFrame(frame, astiav.NewBuffersrcFlags()); err != nil {
		return fmt.Errorf("error while feeding the filtergraph: %w", err)
	}

	// Pull filtered frames from the filtergraph
	for {
		if err := fctx.buffersinkCtx.GetFrame(fctx.filteredFrame, astiav.NewBuffersinkFlags()); err != nil {
			if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("error while pulling from filtergraph: %w", err)
		}

		if err := encodeWriteFrame(fctx.filteredFrame, streamIndex, nil); err != nil {
			return fmt.Errorf("encode write frame failed: %w", err)
		}

		fctx.filteredFrame.Unref()
	}

	return nil
}

func encodeWriteFrame(frame *astiav.Frame, streamIndex int, gotFrame *bool) error {
	sctx := streamCtx[streamIndex]
	if sctx.encCtx == nil {
		return nil
	}

	// Send frame to encoder
	if err := sctx.encCtx.SendFrame(frame); err != nil {
		return fmt.Errorf("sending frame to encoder failed: %w", err)
	}

	// Receive packets from encoder
	for {
		fctx := filterCtx[streamIndex]
		if err := sctx.encCtx.ReceivePacket(fctx.encPkt); err != nil {
			if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("receiving packet from encoder failed: %w", err)
		}

		// Prepare packet for muxing
		fctx.encPkt.SetStreamIndex(streamIndex)
		
		// Rescale timestamps
		inTimeBase := streamCtx[streamIndex].decCtx.TimeBase()
		outTimeBase := ofmtCtx.Streams()[streamIndex].TimeBase()
		
		fctx.encPkt.RescaleTs(inTimeBase, outTimeBase)

		// Write packet
		if err := ofmtCtx.WriteFrame(fctx.encPkt); err != nil {
			return fmt.Errorf("writing packet failed: %w", err)
		}

		if gotFrame != nil {
			*gotFrame = true
		}

		fctx.encPkt.Unref()
	}

	return nil
}

func flushDecoderEncoder(streamIndex int) error {
	sctx := streamCtx[streamIndex]
	
	if sctx.decCtx.MediaType() == astiav.MediaTypeVideo || sctx.decCtx.MediaType() == astiav.MediaTypeAudio {
		// Flush decoder
		if err := sctx.decCtx.SendPacket(nil); err != nil {
			return fmt.Errorf("flushing decoder failed: %w", err)
		}

		for {
			if err := sctx.decCtx.ReceiveFrame(sctx.decFrame); err != nil {
				if errors.Is(err, astiav.ErrEof) {
					break
				}
				return fmt.Errorf("receiving frame during flush failed: %w", err)
			}

			// Use frame's original pts or set a reasonable value
			if sctx.decFrame.Pts() == astiav.NoPtsValue {
				sctx.decFrame.SetPts(0)
			}

			if err := filterEncodeWriteFrame(sctx.decFrame, streamIndex); err != nil {
				return fmt.Errorf("filter encode write frame during flush failed: %w", err)
			}

			sctx.decFrame.Unref()
		}

		// Flush filter
		fctx := filterCtx[streamIndex]
		if fctx != nil {
			if err := fctx.buffersrcCtx.AddFrame(nil, astiav.NewBuffersrcFlags()); err != nil {
				return fmt.Errorf("flushing filter failed: %w", err)
			}

			for {
				if err := fctx.buffersinkCtx.GetFrame(fctx.filteredFrame, astiav.NewBuffersinkFlags()); err != nil {
					if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
						break
					}
					return fmt.Errorf("error while pulling from filtergraph during flush: %w", err)
				}

				if err := encodeWriteFrame(fctx.filteredFrame, streamIndex, nil); err != nil {
					return fmt.Errorf("encode write frame during flush failed: %w", err)
				}

				fctx.filteredFrame.Unref()
			}
		}

		// Flush encoder
		if sctx.encCtx != nil {
			if err := encodeWriteFrame(nil, streamIndex, nil); err != nil {
				return fmt.Errorf("flushing encoder failed: %w", err)
			}
		}
	}

	return nil
}
