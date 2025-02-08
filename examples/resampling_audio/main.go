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

var (
	af             *astiav.AudioFifo
	decodedFrame   *astiav.Frame
	finalFrame     *astiav.Frame
	resampledFrame *astiav.Frame
	src            *astiav.SoftwareResampleContext
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
		log.Fatal(errors.New("main: input format context is nil"))
	}
	defer inputFormatContext.Free()

	// Open input
	if err := inputFormatContext.OpenInput(*input, nil, nil); err != nil {
		log.Fatal(fmt.Errorf("main: opening input failed: %w", err))
	}
	defer inputFormatContext.CloseInput()

	// Find stream info
	if err := inputFormatContext.FindStreamInfo(nil); err != nil {
		log.Fatal(fmt.Errorf("main: finding stream info failed: %w", err))
	}

	// Loop through streams
	var s *astiav.Stream
	var cc *astiav.CodecContext
	for _, is := range inputFormatContext.Streams() {
		// Only process audio
		if is.CodecParameters().MediaType() != astiav.MediaTypeAudio {
			continue
		}

		// Store stream
		s = is

		// Find decoder
		c := astiav.FindDecoder(is.CodecParameters().CodecID())
		if c == nil {
			log.Fatal(errors.New("main: codec is nil"))
		}

		// Allocate codec context
		if cc = astiav.AllocCodecContext(c); cc == nil {
			log.Fatal(errors.New("main: codec context is nil"))
		}
		defer cc.Free()

		// Update codec context
		if err := is.CodecParameters().ToCodecContext(cc); err != nil {
			log.Fatal(fmt.Errorf("main: updating codec context failed: %w", err))
		}

		// Open codec context
		if err := cc.Open(c, nil); err != nil {
			log.Fatal(fmt.Errorf("main: opening codec context failed: %w", err))
		}
		break
	}

	// No stream
	if s == nil {
		log.Fatal("main: no audio stream found")
	}

	// Alloc resample context
	src = astiav.AllocSoftwareResampleContext()
	defer src.Free()

	// Allocate packet
	pkt := astiav.AllocPacket()
	defer pkt.Free()

	// Allocate decoded frame
	decodedFrame = astiav.AllocFrame()
	defer decodedFrame.Free()

	// Allocate resampled frame
	resampledFrame = astiav.AllocFrame()
	defer resampledFrame.Free()

	// For the resampled frame we need to setup mandatory information
	resampledFrame.SetChannelLayout(astiav.ChannelLayoutStereo)
	resampledFrame.SetSampleFormat(astiav.SampleFormatFltp)
	resampledFrame.SetSampleRate(24000)

	// Do this only if you want to make sure the resampled frame's number of samples doesn't get
	// bigger than a custom value ("200" in our case)
	resampledFrame.SetNbSamples(200)
	const align = 0
	if err := resampledFrame.AllocBuffer(align); err != nil {
		log.Fatal(fmt.Errorf("main: allocating buffer failed: %w", err))
	}

	// For the sake of the example we use an audio FIFO to ensure final frames have an exact constant
	// number of samples except for the last one. However this is optional and it depends on your use case
	finalFrame = astiav.AllocFrame()
	defer finalFrame.Free()
	finalFrame.SetChannelLayout(resampledFrame.ChannelLayout())
	finalFrame.SetNbSamples(resampledFrame.NbSamples())
	finalFrame.SetSampleFormat(resampledFrame.SampleFormat())
	finalFrame.SetSampleRate(resampledFrame.SampleRate())
	if err := finalFrame.AllocBuffer(align); err != nil {
		log.Fatal(fmt.Errorf("main: allocating buffer failed: %w", err))
	}
	af = astiav.AllocAudioFifo(finalFrame.SampleFormat(), finalFrame.ChannelLayout().Channels(), finalFrame.NbSamples())
	defer af.Free()

	// Loop
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
			if pkt.StreamIndex() != s.Index() {
				return false
			}

			// Send packet
			if err := cc.SendPacket(pkt); err != nil {
				log.Fatal(fmt.Errorf("main: sending packet failed: %w", err))
			}

			// Loop
			for {
				// We use a closure to ease unreferencing the frame
				if stop := func() bool {
					// Receive frame
					if err := cc.ReceiveFrame(decodedFrame); err != nil {
						if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
							return true
						}
						log.Fatal(fmt.Errorf("main: receiving frame failed: %w", err))
					}

					// Make sure to unreference the frame
					defer decodedFrame.Unref()

					// Log
					log.Printf("new decoded frame: nb samples: %d", decodedFrame.NbSamples())

					// Resample decoded frame
					if err := src.ConvertFrame(decodedFrame, resampledFrame); err != nil {
						log.Fatal(fmt.Errorf("main: resampling decoded frame failed: %w", err))
					}

					// Something was resampled
					if nbSamples := resampledFrame.NbSamples(); nbSamples > 0 {
						// Log
						log.Printf("new resampled frame: nb samples: %d", nbSamples)

						// Add resampled frame to audio fifo
						if err := addResampledFrameToAudioFIFO(false); err != nil {
							log.Fatal(fmt.Errorf("main: adding resampled frame to audio fifo failed: %w", err))
						}

						// Flush software resample context
						if err := flushSoftwareResampleContext(false); err != nil {
							log.Fatal(fmt.Errorf("main: flushing software resample context failed: %w", err))
						}
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

	// Flush software resample context
	if err := flushSoftwareResampleContext(true); err != nil {
		log.Fatal(fmt.Errorf("main: flushing software resample context failed: %w", err))
	}

	// Success
	log.Println("success")
}

func flushSoftwareResampleContext(finalFlush bool) error {
	// Loop
	for {
		// We're making the final flush or there's enough data to flush the resampler
		if finalFlush || src.Delay(int64(resampledFrame.SampleRate())) >= int64(resampledFrame.NbSamples()) {
			// Flush resampler
			if err := src.ConvertFrame(nil, resampledFrame); err != nil {
				log.Fatal(fmt.Errorf("main: flushing resampler failed: %w", err))
			}

			// Log
			if resampledFrame.NbSamples() > 0 {
				log.Printf("new resampled frame: nb samples: %d", resampledFrame.NbSamples())
			}

			// Add resampled frame to audio fifo
			if err := addResampledFrameToAudioFIFO(finalFlush); err != nil {
				log.Fatal(fmt.Errorf("main: adding resampled frame to audio fifo failed: %w", err))
			}

			// Final flush is done
			if finalFlush && resampledFrame.NbSamples() == 0 {
				break
			}
			continue
		}
		break
	}
	return nil
}

func addResampledFrameToAudioFIFO(flush bool) error {
	// Write
	if resampledFrame.NbSamples() > 0 {
		if _, err := af.Write(resampledFrame); err != nil {
			return fmt.Errorf("main: writing failed: %w", err)
		}
	}

	// Loop
	for {
		// We're flushing or there's enough data to read
		if (flush && af.Size() > 0) || (!flush && af.Size() >= finalFrame.NbSamples()) {
			// Read
			n, err := af.Read(finalFrame)
			if err != nil {
				return fmt.Errorf("main: reading failed: %w", err)
			}
			finalFrame.SetNbSamples(n)

			// Log
			log.Printf("new final frame: nb samples: %d", finalFrame.NbSamples())
			continue
		}
		break
	}
	return nil
}
