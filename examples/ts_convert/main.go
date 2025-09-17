/*
 * SPTS/MPTS Transport Stream Converter
 *
 * Go implementation for converting between Single Program Transport Stream (SPTS)
 * and Multiple Program Transport Stream (MPTS)
 *
 * This example demonstrates:
 * - SPTS to MPTS conversion (combining multiple SPTS into one MPTS)
 * - MPTS to SPTS conversion (extracting individual programs from MPTS)
 * - PAT/PMT table handling and PID remapping
 * - Transport Stream packet processing
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/asticode/go-astiav"
)

var (
	mode       = flag.String("mode", "", "conversion mode: spts2mpts or mpts2spts")
	input      = flag.String("i", "", "input file(s) - for spts2mpts: comma-separated SPTS files; for mpts2spts: single MPTS file")
	output     = flag.String("o", "", "output file(s) - for spts2mpts: single MPTS file; for mpts2spts: output directory")
	programIds = flag.String("pids", "", "program IDs for mpts2spts mode (comma-separated, e.g., '100,200')")
)

type TSConverter struct {
	inputFiles    []string
	outputPath    string
	mode          string
	programIds    []int
	inputContexts []*astiav.FormatContext
	outputContext *astiav.FormatContext
}

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

	// Usage
	if *mode == "" || *input == "" || *output == "" {
		printUsage()
		return
	}

	// Create converter
	converter := &TSConverter{
		mode:       *mode,
		outputPath: *output,
	}

	// Parse input files
	converter.inputFiles = strings.Split(*input, ",")
	for i := range converter.inputFiles {
		converter.inputFiles[i] = strings.TrimSpace(converter.inputFiles[i])
	}

	// Parse program IDs for mpts2spts mode
	if *mode == "mpts2spts" && *programIds != "" {
		pidStrs := strings.Split(*programIds, ",")
		for _, pidStr := range pidStrs {
			pid, err := strconv.Atoi(strings.TrimSpace(pidStr))
			if err != nil {
				log.Fatal(fmt.Errorf("invalid program ID: %s", pidStr))
			}
			converter.programIds = append(converter.programIds, pid)
		}
	}

	// Execute conversion
	if err := converter.convert(); err != nil {
		log.Fatal(fmt.Errorf("conversion failed: %w", err))
	}

	log.Println("Transport Stream conversion completed successfully")
}

func printUsage() {
	fmt.Printf("Usage: %s -mode <mode> -i <input> -o <output> [options]\n\n", os.Args[0])
	fmt.Println("Transport Stream Converter - Convert between SPTS and MPTS")
	fmt.Println()
	fmt.Println("Modes:")
	fmt.Println("  spts2mpts  Convert multiple SPTS files to single MPTS file")
	fmt.Println("  mpts2spts  Convert MPTS file to multiple SPTS files")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Convert multiple SPTS to MPTS")
	fmt.Printf("  %s -mode spts2mpts -i \"input1.ts,input2.ts,input3.ts\" -o output.ts\n", os.Args[0])
	fmt.Println()
	fmt.Println("  # Convert MPTS to SPTS (extract all programs)")
	fmt.Printf("  %s -mode mpts2spts -i input.ts -o output_dir/\n", os.Args[0])
	fmt.Println()
	fmt.Println("  # Convert MPTS to SPTS (extract specific programs)")
	fmt.Printf("  %s -mode mpts2spts -i input.ts -o output_dir/ -pids \"100,200\"\n", os.Args[0])
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -pids     Program IDs to extract (for mpts2spts mode)")
}

func (c *TSConverter) convert() error {
	switch c.mode {
	case "spts2mpts":
		return c.convertSPTSToMPTS()
	case "mpts2spts":
		return c.convertMPTSToSPTS()
	default:
		return fmt.Errorf("unsupported mode: %s", c.mode)
	}
}

func (c *TSConverter) convertSPTSToMPTS() error {
	log.Printf("Converting %d SPTS files to MPTS: %s", len(c.inputFiles), c.outputPath)

	// Open input files
	for _, inputFile := range c.inputFiles {
		inputCtx := astiav.AllocFormatContext()
		if inputCtx == nil {
			return errors.New("failed to allocate input format context")
		}

		if err := inputCtx.OpenInput(inputFile, nil, nil); err != nil {
			return fmt.Errorf("opening input file %s failed: %w", inputFile, err)
		}

		if err := inputCtx.FindStreamInfo(nil); err != nil {
			return fmt.Errorf("finding stream info for %s failed: %w", inputFile, err)
		}

		c.inputContexts = append(c.inputContexts, inputCtx)
		log.Printf("Opened SPTS file: %s (%d streams)", inputFile, inputCtx.NbStreams())
	}

	// Create output MPTS file
	outputCtx, err := astiav.AllocOutputFormatContext(nil, "mpegts", c.outputPath)
	if err != nil {
		return fmt.Errorf("allocating output format context failed: %w", err)
	}
	c.outputContext = outputCtx

	// Open output file
	ioCtx, err := astiav.OpenIOContext(c.outputPath, astiav.NewIOContextFlags(astiav.IOContextFlagWrite), nil, nil)
	if err != nil {
		return fmt.Errorf("opening output file failed: %w", err)
	}
	outputCtx.SetPb(ioCtx)

	// Create programs in the output context
	streamMapping := make(map[int]map[int]int) // inputIndex -> inputStreamIndex -> outputStreamIndex
	programId := 100                           // Starting program ID
	streamPID := 256                           // Starting PID for streams (0-255 are reserved)

	for inputIdx, inputCtx := range c.inputContexts {
		streamMapping[inputIdx] = make(map[int]int)
		currentProgramId := programId + inputIdx // Each input gets a unique program ID

		// Create a program in the output context
		program := outputCtx.NewProgram(currentProgramId)
		if program == nil {
			return errors.New("failed to create program")
		}

		// Add streams to this program
		for i := 0; i < inputCtx.NbStreams(); i++ {
			inputStream := inputCtx.Streams()[i]
			outputStream := outputCtx.NewStream(nil)
			if outputStream == nil {
				return errors.New("failed to create output stream")
			}

			// Copy stream parameters
			if err := outputStream.CodecParameters().Copy(inputStream.CodecParameters()); err != nil {
				return fmt.Errorf("copying stream parameters failed: %w", err)
			}

			// Set unique PID for each stream
			outputStream.SetID(streamPID)
			streamMapping[inputIdx][i] = outputStream.Index()

			// Add stream to the program
			program.AddStream(outputStream)

			streamPID++ // Each stream gets a unique PID
		}

		// Set program metadata
		metadata := astiav.NewDictionary()
		metadata.Set("service_name", fmt.Sprintf("Service%02d", currentProgramId), 0)
		metadata.Set("service_provider", "Go-AstiAV", 0)
		program.SetMetadata(metadata)

		// Set program number and PMT PID
		program.SetProgramNumber(currentProgramId)
		program.SetPmtPid(streamPID)
		streamPID++

		log.Printf("Created Program %d from input file %d (%d streams)",
			currentProgramId, inputIdx+1, inputCtx.NbStreams())
	}

	// Write header
	if err := outputCtx.WriteHeader(nil); err != nil {
		return fmt.Errorf("writing header failed: %w", err)
	}

	// Process packets from all inputs
	packet := astiav.AllocPacket()
	defer packet.Free()

	for inputIdx, inputCtx := range c.inputContexts {
		log.Printf("Processing packets from input %d...", inputIdx+1)

		for {
			if err := inputCtx.ReadFrame(packet); err != nil {
				if errors.Is(err, astiav.ErrEof) {
					break
				}
				return fmt.Errorf("reading packet from input %d failed: %w", inputIdx+1, err)
			}

			// Remap stream index
			if newStreamIndex, exists := streamMapping[inputIdx][packet.StreamIndex()]; exists {
				packet.SetStreamIndex(newStreamIndex)

				// Write packet to output
				if err := outputCtx.WriteInterleavedFrame(packet); err != nil {
					return fmt.Errorf("writing packet failed: %w", err)
				}
			}

			packet.Unref()
		}
	}

	// Write trailer
	if err := outputCtx.WriteTrailer(); err != nil {
		return fmt.Errorf("writing trailer failed: %w", err)
	}

	// Cleanup
	for _, inputCtx := range c.inputContexts {
		inputCtx.CloseInput()
		inputCtx.Free()
	}

	if outputCtx.Pb() != nil {
		outputCtx.Pb().Close()
	}
	outputCtx.Free()

	log.Printf("Successfully created MPTS file: %s", c.outputPath)
	return nil
}

func (c *TSConverter) convertMPTSToSPTS() error {
	log.Printf("Converting MPTS to SPTS: %s -> %s", c.inputFiles[0], c.outputPath)

	// Open input MPTS file
	inputCtx := astiav.AllocFormatContext()
	if inputCtx == nil {
		return errors.New("failed to allocate input format context")
	}
	defer inputCtx.Free()

	if err := inputCtx.OpenInput(c.inputFiles[0], nil, nil); err != nil {
		return fmt.Errorf("opening input file failed: %w", err)
	}
	defer inputCtx.CloseInput()

	if err := inputCtx.FindStreamInfo(nil); err != nil {
		return fmt.Errorf("finding stream info failed: %w", err)
	}

	log.Printf("Input MPTS has %d streams", inputCtx.NbStreams())

	// Create output directory
	if err := os.MkdirAll(c.outputPath, 0755); err != nil {
		return fmt.Errorf("creating output directory failed: %w", err)
	}

	// Group streams by program ID
	programStreams := make(map[int][]int) // programId -> streamIndexes

	for i := 0; i < inputCtx.NbStreams(); i++ {
		stream := inputCtx.Streams()[i]
		programId := stream.ID()
		if programId == 0 {
			programId = 100 // Default program ID
		}
		programStreams[programId] = append(programStreams[programId], i)
	}

	log.Printf("Found %d programs in MPTS", len(programStreams))

	// Filter programs if specific PIDs requested
	if len(c.programIds) > 0 {
		filteredPrograms := make(map[int][]int)
		for _, pid := range c.programIds {
			if streams, exists := programStreams[pid]; exists {
				filteredPrograms[pid] = streams
			}
		}
		programStreams = filteredPrograms
	}

	// Create SPTS files for each program
	outputContexts := make(map[int]*astiav.FormatContext)
	streamMappings := make(map[int]map[int]int) // programId -> inputStreamIndex -> outputStreamIndex

	for programId, streamIndexes := range programStreams {
		outputFile := filepath.Join(c.outputPath, fmt.Sprintf("program_%d.ts", programId))

		outputCtx, err := astiav.AllocOutputFormatContext(nil, "mpegts", outputFile)
		if err != nil {
			return fmt.Errorf("allocating output format context for program %d failed: %w", programId, err)
		}

		ioCtx, err := astiav.OpenIOContext(outputFile, astiav.NewIOContextFlags(astiav.IOContextFlagWrite), nil, nil)
		if err != nil {
			return fmt.Errorf("opening output file %s failed: %w", outputFile, err)
		}
		outputCtx.SetPb(ioCtx)

		streamMappings[programId] = make(map[int]int)

		// Create streams for this program
		for _, inputStreamIndex := range streamIndexes {
			inputStream := inputCtx.Streams()[inputStreamIndex]
			outputStream := outputCtx.NewStream(nil)
			if outputStream == nil {
				return fmt.Errorf("failed to create output stream for program %d", programId)
			}

			if err := outputStream.CodecParameters().Copy(inputStream.CodecParameters()); err != nil {
				return fmt.Errorf("copying stream parameters failed: %w", err)
			}

			streamMappings[programId][inputStreamIndex] = outputStream.Index()
		}

		if err := outputCtx.WriteHeader(nil); err != nil {
			return fmt.Errorf("writing header for program %d failed: %w", programId, err)
		}

		outputContexts[programId] = outputCtx
		log.Printf("Created SPTS file for program %d: %s (%d streams)",
			programId, outputFile, len(streamIndexes))
	}

	// Process packets
	packet := astiav.AllocPacket()
	defer packet.Free()

	packetCount := 0
	for {
		if err := inputCtx.ReadFrame(packet); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("reading packet failed: %w", err)
		}

		streamIndex := packet.StreamIndex()
		stream := inputCtx.Streams()[streamIndex]
		programId := stream.ID()
		if programId == 0 {
			programId = 100
		}

		// Find which output context this packet belongs to
		if outputCtx, exists := outputContexts[programId]; exists {
			if newStreamIndex, exists := streamMappings[programId][streamIndex]; exists {
				packet.SetStreamIndex(newStreamIndex)

				if err := outputCtx.WriteFrame(packet); err != nil {
					return fmt.Errorf("writing packet to program %d failed: %w", programId, err)
				}
			}
		}

		packet.Unref()
		packetCount++

		if packetCount%10000 == 0 {
			log.Printf("Processed %d packets...", packetCount)
		}
	}

	// Write trailers and cleanup
	for programId, outputCtx := range outputContexts {
		if err := outputCtx.WriteTrailer(); err != nil {
			log.Printf("Warning: writing trailer for program %d failed: %v", programId, err)
		}

		if outputCtx.Pb() != nil {
			outputCtx.Pb().Close()
		}
		outputCtx.Free()
	}

	log.Printf("Successfully extracted %d programs to SPTS files", len(outputContexts))
	log.Printf("Total packets processed: %d", packetCount)
	return nil
}
