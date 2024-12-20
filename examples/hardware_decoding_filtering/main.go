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
	decoderCodecName       = flag.String("c", "", "the decoder codec name (e.g. h264_cuvid)")
	filter                 = flag.String("f", "", "the hardware filter")
	hardwareDeviceName     = flag.String("n", "", "the hardware device name (e.g. 0)")
	hardwareDeviceTypeName = flag.String("t", "", "the hardware device type (e.g. cuda)")
	input                  = flag.String("i", "", "the input path")
)

var (
	buffersinkContext     *astiav.BuffersinkFilterContext
	buffersrcContext      *astiav.BuffersrcFilterContext
	c                     = astikit.NewCloser()
	decCodec              *astiav.Codec
	decCodecContext       *astiav.CodecContext
	decodedHardwareFrame  *astiav.Frame
	filterGraph           *astiav.FilterGraph
	filteredHardwareFrame *astiav.Frame
	inputStream           *astiav.Stream
	softwareFrame         *astiav.Frame
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
	if *input == "" || *hardwareDeviceTypeName == "" {
		log.Println("Usage: <binary path> -t <hardware device type> -i <input path> [-n <hardware device name> -c <decoder codec> -f <hardware filter>]")
		return
	}

	// We use an astikit.Closer to free all resources properly
	defer c.Close()

	// Get hardware device type
	hardwareDeviceType := astiav.FindHardwareDeviceTypeByName(*hardwareDeviceTypeName)
	if hardwareDeviceType == astiav.HardwareDeviceTypeNone {
		log.Fatal(errors.New("main: hardware device not found"))
	}

	// Allocate packet
	pkt := astiav.AllocPacket()
	c.Add(pkt.Free)

	// Allocate decoded hardware frame
	decodedHardwareFrame = astiav.AllocFrame()
	c.Add(decodedHardwareFrame.Free)

	// Allocate software frame
	softwareFrame = astiav.AllocFrame()
	c.Add(softwareFrame.Free)

	// Allocate input format context
	inputFormatContext := astiav.AllocFormatContext()
	if inputFormatContext == nil {
		log.Fatal(errors.New("main: input format context is nil"))
	}
	c.Add(inputFormatContext.Free)

	// Open input
	if err := inputFormatContext.OpenInput(*input, nil, nil); err != nil {
		log.Fatal(fmt.Errorf("main: opening input failed: %w", err))
	}
	c.Add(inputFormatContext.CloseInput)

	// Find stream info
	if err := inputFormatContext.FindStreamInfo(nil); err != nil {
		log.Fatal(fmt.Errorf("main: finding stream info failed: %w", err))
	}

	// Loop through streams
	var hdc *astiav.HardwareDeviceContext
	hardwarePixelFormat := astiav.PixelFormatNone
	for _, is := range inputFormatContext.Streams() {
		// Only process video
		if is.CodecParameters().MediaType() != astiav.MediaTypeVideo {
			continue
		}

		// Update input stream
		inputStream = is

		// Find decoder
		decCodec = astiav.FindDecoder(is.CodecParameters().CodecID())
		if *decoderCodecName != "" {
			decCodec = astiav.FindDecoderByName(*decoderCodecName)
		}

		// No codec
		if decCodec == nil {
			log.Fatal(errors.New("main: codec is nil"))
		}

		// Allocate codec context
		if decCodecContext = astiav.AllocCodecContext(decCodec); decCodecContext == nil {
			log.Fatal(errors.New("main: codec context is nil"))
		}
		c.Add(decCodecContext.Free)

		// Loop through codec hardware configs
		for _, p := range decCodec.HardwareConfigs() {
			// Valid hardware config
			if p.MethodFlags().Has(astiav.CodecHardwareConfigMethodFlagHwDeviceCtx) && p.HardwareDeviceType() == hardwareDeviceType {
				hardwarePixelFormat = p.PixelFormat()
				break
			}
		}

		// No valid hardware pixel format
		if hardwarePixelFormat == astiav.PixelFormatNone {
			log.Fatal(errors.New("main: hardware device type not supported by decoder"))
		}

		// Update codec context
		if err := is.CodecParameters().ToCodecContext(decCodecContext); err != nil {
			log.Fatal(fmt.Errorf("main: updating codec context failed: %w", err))
		}

		// Create hardware device context
		var err error
		if hdc, err = astiav.CreateHardwareDeviceContext(hardwareDeviceType, *hardwareDeviceName, nil, 0); err != nil {
			log.Fatal(fmt.Errorf("main: creating hardware device context failed: %w", err))
		}
		c.Add(hdc.Free)

		// Update decoder context
		decCodecContext.SetHardwareDeviceContext(hdc)
		decCodecContext.SetPixelFormatCallback(func(pfs []astiav.PixelFormat) astiav.PixelFormat {
			for _, pf := range pfs {
				if pf == hardwarePixelFormat {
					return pf
				}
			}
			log.Fatal(errors.New("main: using hardware pixel format failed"))
			return astiav.PixelFormatNone
		})

		// Open codec context
		if err := decCodecContext.Open(decCodec, nil); err != nil {
			log.Fatal(fmt.Errorf("main: opening codec context failed: %w", err))
		}
		break
	}

	// No video stream
	if inputStream == nil {
		log.Fatal("main: no video stream found")
	}

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
			if pkt.StreamIndex() != inputStream.Index() {
				return false
			}

			// Send packet
			if err := decCodecContext.SendPacket(pkt); err != nil {
				log.Fatal(fmt.Errorf("main: sending packet failed: %w", err))
			}

			// Loop
			for {
				// We use a closure to ease unreferencing frames
				if stop := func() bool {
					// Receive frame
					if err := decCodecContext.ReceiveFrame(decodedHardwareFrame); err != nil {
						if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
							return true
						}
						log.Fatal(fmt.Errorf("main: receiving frame failed: %w", err))
					}

					// Make sure to unreference hardware frame
					defer decodedHardwareFrame.Unref()

					// Invalid pixel format
					if decodedHardwareFrame.PixelFormat() != hardwarePixelFormat {
						log.Fatalf("main: invalid decoded pixel format %s, expected %s", decodedHardwareFrame.PixelFormat(), hardwarePixelFormat)
					}

					// No filter requested
					if *filter == "" {
						// Do something with hardware frame
						if err := doSomethingWithHardwareFrame(decodedHardwareFrame); err != nil {
							log.Fatal(fmt.Errorf("main: doing something with hardware frame failed: %w", err))
						}
						return false
					}

					// Make sure the filter is initialized
					// We need to wait for the first frame to be decoded before initializing the filter
					// since the decoder codec context doesn't have a valid hardware frame context until then
					if filterGraph == nil {
						if err := initFilter(); err != nil {
							log.Fatal(fmt.Errorf("main: initializing filter failed: %w", err))
						}
					}

					// Filter frame
					if err := filterFrame(); err != nil {
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

	// Success
	log.Println("success")
}

func initFilter() (err error) {
	// Allocate graph
	if filterGraph = astiav.AllocFilterGraph(); filterGraph == nil {
		err = errors.New("main: graph is nil")
		return
	}
	c.Add(filterGraph.Free)

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
	if buffersrcContext, err = filterGraph.NewBuffersrcFilterContext(buffersrc, "in"); err != nil {
		err = fmt.Errorf("main: creating buffersrc context failed: %w", err)
		return
	}
	if buffersinkContext, err = filterGraph.NewBuffersinkFilterContext(buffersink, "in"); err != nil {
		err = fmt.Errorf("main: creating buffersink context failed: %w", err)
		return
	}

	// Create buffersrc context parameters
	buffersrcContextParameters := astiav.AllocBuffersrcFilterContextParameters()
	defer buffersrcContextParameters.Free()
	buffersrcContextParameters.SetHardwareFramesContext(decCodecContext.HardwareFramesContext())
	buffersrcContextParameters.SetHeight(decCodecContext.Height())
	buffersrcContextParameters.SetPixelFormat(decCodecContext.PixelFormat())
	buffersrcContextParameters.SetSampleAspectRatio(decCodecContext.SampleAspectRatio())
	buffersrcContextParameters.SetTimeBase(inputStream.TimeBase())
	buffersrcContextParameters.SetWidth(decCodecContext.Width())

	// Set buffersrc context parameters
	if err = buffersrcContext.SetParameters(buffersrcContextParameters); err != nil {
		err = fmt.Errorf("main: setting buffersrc context parameters failed: %w", err)
		return
	}

	// Initialize buffersrc context
	if err = buffersrcContext.Initialize(); err != nil {
		err = fmt.Errorf("main: initializing buffersrc context failed: %w", err)
		return
	}

	// Update outputs
	outputs.SetName("in")
	outputs.SetFilterContext(buffersrcContext.FilterContext())
	outputs.SetPadIdx(0)
	outputs.SetNext(nil)

	// Update inputs
	inputs.SetName("out")
	inputs.SetFilterContext(buffersinkContext.FilterContext())
	inputs.SetPadIdx(0)
	inputs.SetNext(nil)

	// Parse
	if err = filterGraph.Parse(*filter, inputs, outputs); err != nil {
		err = fmt.Errorf("main: parsing filter failed: %w", err)
		return
	}

	// Configure
	if err = filterGraph.Configure(); err != nil {
		err = fmt.Errorf("main: configuring filter failed: %w", err)
		return
	}

	// Allocate frame
	filteredHardwareFrame = astiav.AllocFrame()
	c.Add(filteredHardwareFrame.Free)
	return
}

func filterFrame() (err error) {
	// Add frame
	if err = buffersrcContext.AddFrame(decodedHardwareFrame, astiav.NewBuffersrcFlags(astiav.BuffersrcFlagKeepRef)); err != nil {
		err = fmt.Errorf("main: adding frame failed: %w", err)
		return
	}

	// Loop
	for {
		// We use a closure to ease unreferencing the frame
		if stop, err := func() (bool, error) {
			// Get frame
			if err := buffersinkContext.GetFrame(filteredHardwareFrame, astiav.NewBuffersinkFlags()); err != nil {
				if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
					return true, nil
				}
				return false, fmt.Errorf("main: getting frame failed: %w", err)
			}

			// Make sure to unrefernce the frame
			defer filteredHardwareFrame.Unref()

			// Do something with hardware frame
			if err := doSomethingWithHardwareFrame(filteredHardwareFrame); err != nil {
				return false, fmt.Errorf("main: doing something with hardware frame failed: %w", err)
			}
			return false, nil
		}(); err != nil {
			return err
		} else if stop {
			break
		}

	}
	return
}

func doSomethingWithHardwareFrame(hardwareFrame *astiav.Frame) error {
	// Transfer hardware data
	if err := hardwareFrame.TransferHardwareData(softwareFrame); err != nil {
		return fmt.Errorf("main: transferring hardware data failed: %w", err)
	}

	// Make sure to unreference software frame
	defer softwareFrame.Unref()

	// Update pts
	softwareFrame.SetPts(hardwareFrame.Pts())

	// Do something with software frame
	log.Printf("new software frame: pts: %d", softwareFrame.Pts())
	return nil
}
