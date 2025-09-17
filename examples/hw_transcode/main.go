/*
 * Hardware-accelerated transcoding example
 * 
 * Go implementation combining FFmpeg's vaapi_transcode.c and hw_decode.c examples
 * 
 * This example demonstrates hardware-accelerated video transcoding
 * using hardware decoding and encoding with various hardware acceleration APIs.
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/asticode/go-astiav"
)

var (
	input      = flag.String("i", "", "the input path")
	output     = flag.String("o", "", "the output path")
	hwType     = flag.String("hwtype", "", "hardware device type (videotoolbox, vaapi, cuda, etc.)")
	encCodec   = flag.String("c", "", "encoder codec (h264_videotoolbox, h264_vaapi, h264_nvenc, etc.)")
	filterDesc = flag.String("filter", "", "filter description (e.g., 'scale=640:480', 'hwupload,scale_vt=1280:720,hwdownload')")
)

type HWTranscodeContext struct {
	// Input context
	inputFormatCtx *astiav.FormatContext
	decoderCtx     *astiav.CodecContext
	videoStreamIdx int
	
	// Output context  
	outputFormatCtx *astiav.FormatContext
	encoderCtx      *astiav.CodecContext
	outputStream    *astiav.Stream
	
	// Hardware context
	hwDeviceCtx *astiav.HardwareDeviceContext
	hwPixFmt    astiav.PixelFormat
	
	// Hardware filter context
	filterGraph       *astiav.FilterGraph
	buffersrcCtx      *astiav.BuffersrcFilterContext
	buffersinkCtx     *astiav.BuffersinkFilterContext
	hwFramesCtx       *astiav.HardwareFramesContext
	useHardwareFilter bool
	
	// Working frames and packets
	packet       *astiav.Packet
	frame        *astiav.Frame
	hwFrame      *astiav.Frame
	swFrame      *astiav.Frame
	filteredFrame *astiav.Frame
}

func main() {
	// Handle ffmpeg logs
	astiav.SetLogLevel(astiav.LogLevelInfo)
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
	if *input == "" || *output == "" || *hwType == "" || *encCodec == "" {
		fmt.Printf("usage: %s -i input -o output -hwtype hwtype -c encoder [-filter filter_desc]\n", os.Args[0])
		fmt.Printf("Hardware-accelerated transcoding example with optional hardware filtering.\n")
		fmt.Printf("Perform hardware-accelerated transcoding using specified hardware device.\n\n")
		fmt.Printf("Examples:\n")
		fmt.Printf("  # VideoToolbox (macOS)\n")
		fmt.Printf("  %s -i input.mp4 -o output.mp4 -hwtype videotoolbox -c h264_videotoolbox\n", os.Args[0])
		fmt.Printf("  # VideoToolbox with hardware scaling\n")
		fmt.Printf("  %s -i input.mp4 -o output.mp4 -hwtype videotoolbox -c h264_videotoolbox -filter 'scale_vt=1280:720'\n", os.Args[0])
		fmt.Printf("  # VAAPI (Linux)\n")
		fmt.Printf("  %s -i input.mp4 -o output.mp4 -hwtype vaapi -c h264_vaapi\n", os.Args[0])
		fmt.Printf("  # VAAPI with hardware scaling\n")
		fmt.Printf("  %s -i input.mp4 -o output.mp4 -hwtype vaapi -c h264_vaapi -filter 'scale_vaapi=1280:720'\n", os.Args[0])
		fmt.Printf("  # NVENC (NVIDIA)\n")
		fmt.Printf("  %s -i input.mp4 -o output.mp4 -hwtype cuda -c h264_nvenc\n", os.Args[0])
		fmt.Printf("  # NVENC with hardware scaling\n")
		fmt.Printf("  %s -i input.mp4 -o output.mp4 -hwtype cuda -c h264_nvenc -filter 'scale_cuda=1280:720'\n", os.Args[0])
		fmt.Printf("\nAvailable hardware device types:\n")
		
		// 显示可用的硬件设备类型
		fmt.Printf("  - videotoolbox (macOS)\n")
		fmt.Printf("  - vaapi (Linux)\n")
		fmt.Printf("  - cuda (NVIDIA)\n")
		fmt.Printf("  - qsv (Intel)\n")
		fmt.Printf("  - dxva2 (Windows)\n")
		fmt.Printf("  - d3d11va (Windows)\n")
		return
	}

	// 创建转码上下文
	ctx := &HWTranscodeContext{}
	defer ctx.cleanup()

	// 初始化硬件设备
	if err := ctx.initHardwareDevice(*hwType); err != nil {
		log.Fatal(fmt.Errorf("initializing hardware device failed: %w", err))
	}

	// 打开输入文件
	if err := ctx.openInputFile(*input); err != nil {
		log.Fatal(fmt.Errorf("opening input file failed: %w", err))
	}

	// 打开输出文件
	if err := ctx.openOutputFile(*output, *encCodec); err != nil {
		log.Fatal(fmt.Errorf("opening output file failed: %w", err))
	}

	// 初始化硬件filter（如果指定了filter）
	if *filterDesc != "" {
		if err := ctx.initHardwareFilter(*filterDesc); err != nil {
			log.Fatal(fmt.Errorf("initializing hardware filter failed: %w", err))
		}
		ctx.useHardwareFilter = true
		
		// 更新编码器分辨率以匹配filter输出
		if err := ctx.updateEncoderForFilter(); err != nil {
			log.Fatal(fmt.Errorf("updating encoder for filter failed: %w", err))
		}
		
		log.Printf("Hardware filter initialized: %s", *filterDesc)
	} else {
		// 没有filter的情况，直接打开编码器
		if err := ctx.openEncoder(); err != nil {
			log.Fatal(fmt.Errorf("opening encoder failed: %w", err))
		}
	}

	// 执行转码
	if err := ctx.transcode(); err != nil {
		log.Fatal(fmt.Errorf("transcoding failed: %w", err))
	}

	log.Println("Hardware transcoding completed successfully")
}

// initHardwareDevice 初始化硬件设备
func (ctx *HWTranscodeContext) initHardwareDevice(hwTypeName string) error {
	// 解析硬件设备类型
	hwDeviceType := astiav.FindHardwareDeviceTypeByName(hwTypeName)
	if hwDeviceType == astiav.HardwareDeviceTypeNone {
		return fmt.Errorf("unsupported hardware device type: %s", hwTypeName)
	}

	// 创建硬件设备上下文
	var err error
	ctx.hwDeviceCtx, err = astiav.CreateHardwareDeviceContext(hwDeviceType, "", nil, 0)
	if err != nil {
		return fmt.Errorf("failed to create hardware device context: %w", err)
	}

	// 设置硬件像素格式
	switch hwTypeName {
	case "videotoolbox":
		ctx.hwPixFmt = astiav.PixelFormatVideotoolbox
	case "vaapi":
		ctx.hwPixFmt = astiav.PixelFormatVaapi
	case "cuda":
		ctx.hwPixFmt = astiav.PixelFormatCuda
	case "qsv":
		ctx.hwPixFmt = astiav.PixelFormatQsv
	case "d3d11va":
		ctx.hwPixFmt = astiav.PixelFormatD3D11
	default:
		return fmt.Errorf("unsupported hardware pixel format for device type: %s", hwTypeName)
	}

	log.Printf("Hardware device initialized: %s", hwTypeName)
	return nil
}

// getHwFormat 获取硬件像素格式回调函数
func (ctx *HWTranscodeContext) getHwFormat(pixFmts []astiav.PixelFormat) astiav.PixelFormat {
	for _, pf := range pixFmts {
		if pf == ctx.hwPixFmt {
			return pf
		}
	}
	
	fmt.Fprintf(os.Stderr, "Failed to get HW surface format.\n")
	return astiav.PixelFormatNone
}

// openInputFile 打开输入文件并初始化解码器
func (ctx *HWTranscodeContext) openInputFile(filename string) error {
	// 分配输入格式上下文
	ctx.inputFormatCtx = astiav.AllocFormatContext()
	if ctx.inputFormatCtx == nil {
		return errors.New("failed to allocate input format context")
	}

	// 打开输入文件
	if err := ctx.inputFormatCtx.OpenInput(filename, nil, nil); err != nil {
		return fmt.Errorf("opening input failed: %w", err)
	}

	// 查找流信息
	if err := ctx.inputFormatCtx.FindStreamInfo(nil); err != nil {
		return fmt.Errorf("finding stream info failed: %w", err)
	}

	// 查找最佳视频流
	videoStream, decoder, err := ctx.inputFormatCtx.FindBestStream(astiav.MediaTypeVideo, -1, -1)
	if err != nil {
		return fmt.Errorf("finding best video stream failed: %w", err)
	}
	ctx.videoStreamIdx = videoStream.Index()

	// 分配解码器上下文
	ctx.decoderCtx = astiav.AllocCodecContext(decoder)
	if ctx.decoderCtx == nil {
		return errors.New("failed to allocate decoder context")
	}

	// 复制编解码器参数
	if err := videoStream.CodecParameters().ToCodecContext(ctx.decoderCtx); err != nil {
		return fmt.Errorf("copying codec parameters failed: %w", err)
	}

	// 设置硬件设备上下文和格式回调
	ctx.decoderCtx.SetHardwareDeviceContext(ctx.hwDeviceCtx)
	ctx.decoderCtx.SetPixelFormatCallback(func(pixFmts []astiav.PixelFormat) astiav.PixelFormat {
		return ctx.getHwFormat(pixFmts)
	})

	// 打开解码器
	if err := ctx.decoderCtx.Open(decoder, nil); err != nil {
		return fmt.Errorf("opening decoder failed: %w", err)
	}

	// 创建硬件帧上下文用于filter
	if err := ctx.createHardwareFramesContext(); err != nil {
		return fmt.Errorf("creating hardware frames context failed: %w", err)
	}

	log.Printf("Input file opened: %s", filename)
	log.Printf("Video stream: %dx%d, codec: %s", 
		ctx.decoderCtx.Width(), ctx.decoderCtx.Height(), decoder.Name())
	
	return nil
}

// openOutputFile 打开输出文件并初始化编码器
func (ctx *HWTranscodeContext) openOutputFile(filename, codecName string) error {
	// 查找编码器
	encoder := astiav.FindEncoderByName(codecName)
	if encoder == nil {
		return fmt.Errorf("encoder not found: %s", codecName)
	}

	// 分配输出格式上下文
	var err error
	ctx.outputFormatCtx, err = astiav.AllocOutputFormatContext(nil, "", filename)
	if err != nil {
		return fmt.Errorf("allocating output format context failed: %w", err)
	}

	// 创建输出流
	ctx.outputStream = ctx.outputFormatCtx.NewStream(encoder)
	if ctx.outputStream == nil {
		return errors.New("failed to create output stream")
	}

	// 分配编码器上下文
	ctx.encoderCtx = astiav.AllocCodecContext(encoder)
	if ctx.encoderCtx == nil {
		return errors.New("failed to allocate encoder context")
	}

	// 设置编码器参数 - 如果有filter，稍后会更新分辨率
	ctx.encoderCtx.SetWidth(ctx.decoderCtx.Width())
	ctx.encoderCtx.SetHeight(ctx.decoderCtx.Height())
	ctx.encoderCtx.SetPixelFormat(ctx.hwPixFmt)
	ctx.encoderCtx.SetTimeBase(astiav.NewRational(1, 25))
	ctx.encoderCtx.SetFramerate(astiav.NewRational(25, 1))
	ctx.encoderCtx.SetBitRate(2000000) // 2Mbps

	// 设置硬件设备上下文
	ctx.encoderCtx.SetHardwareDeviceContext(ctx.hwDeviceCtx)

	// 如果输出格式需要全局头部
	if ctx.outputFormatCtx.OutputFormat().Flags().Has(astiav.IOFormatFlagGlobalheader) {
		ctx.encoderCtx.SetFlags(ctx.encoderCtx.Flags().Add(astiav.CodecContextFlagGlobalHeader))
	}

	// 暂时不打开编码器，等filter初始化完成后再打开
	// 这样可以根据filter输出调整编码器参数

	// 打开输出文件
	if !ctx.outputFormatCtx.OutputFormat().Flags().Has(astiav.IOFormatFlagNofile) {
		ioContext, err := astiav.OpenIOContext(filename, astiav.NewIOContextFlags(astiav.IOContextFlagWrite), nil, nil)
		if err != nil {
			return fmt.Errorf("opening output file failed: %w", err)
		}
		ctx.outputFormatCtx.SetPb(ioContext)
	}

	log.Printf("Output file opened: %s", filename)
	log.Printf("Encoder: %s, %dx%d", encoder.Name(), 
		ctx.encoderCtx.Width(), ctx.encoderCtx.Height())

	return nil
}

// transcode 执行硬件转码
func (ctx *HWTranscodeContext) transcode() error {
	// 分配工作用的帧和包
	ctx.packet = astiav.AllocPacket()
	if ctx.packet == nil {
		return errors.New("failed to allocate packet")
	}

	ctx.frame = astiav.AllocFrame()
	if ctx.frame == nil {
		return errors.New("failed to allocate frame")
	}

	ctx.hwFrame = astiav.AllocFrame()
	if ctx.hwFrame == nil {
		return errors.New("failed to allocate hw frame")
	}

	ctx.swFrame = astiav.AllocFrame()
	if ctx.swFrame == nil {
		return errors.New("failed to allocate sw frame")
	}

	frameCount := 0
	
	// 主转码循环
	for {
		// 读取包
		if err := ctx.inputFormatCtx.ReadFrame(ctx.packet); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("reading frame failed: %w", err)
		}

		// 只处理视频流
		if ctx.packet.StreamIndex() != ctx.videoStreamIdx {
			ctx.packet.Unref()
			continue
		}

		// 解码和编码
		if err := ctx.decodeAndEncode(); err != nil {
			return fmt.Errorf("decode and encode failed: %w", err)
		}

		frameCount++
		if frameCount%100 == 0 {
			log.Printf("Processed %d frames", frameCount)
		}

		ctx.packet.Unref()
	}

	// 刷新解码器和编码器
	if err := ctx.flushCodecs(); err != nil {
		return fmt.Errorf("flushing codecs failed: %w", err)
	}

	// 写入文件尾
	if err := ctx.outputFormatCtx.WriteTrailer(); err != nil {
		return fmt.Errorf("writing trailer failed: %w", err)
	}

	log.Printf("Total frames processed: %d", frameCount)
	return nil
}

// decodeAndEncode 解码和编码单个包
func (ctx *HWTranscodeContext) decodeAndEncode() error {
	// 发送包到解码器
	if err := ctx.decoderCtx.SendPacket(ctx.packet); err != nil {
		return fmt.Errorf("sending packet to decoder failed: %w", err)
	}

	// 从解码器接收帧
	for {
		if err := ctx.decoderCtx.ReceiveFrame(ctx.frame); err != nil {
			if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("receiving frame from decoder failed: %w", err)
		}

		// 编码帧
		if err := ctx.encodeFrame(ctx.frame); err != nil {
			return fmt.Errorf("encoding frame failed: %w", err)
		}
	}

	return nil
}

// encodeFrame 编码单个帧
func (ctx *HWTranscodeContext) encodeFrame(frame *astiav.Frame) error {
	var encFrame *astiav.Frame = frame

	// 应用硬件filter（如果启用）
	if ctx.useHardwareFilter {
		filteredFrame, err := ctx.filterFrame(frame)
		if err != nil {
			return fmt.Errorf("filtering frame failed: %w", err)
		}
		if filteredFrame != nil {
			encFrame = filteredFrame
		}
	}

	// 发送帧到编码器
	if err := ctx.encoderCtx.SendFrame(encFrame); err != nil {
		return fmt.Errorf("sending frame to encoder failed: %w", err)
	}

	// 从编码器接收包
	for {
		encPacket := astiav.AllocPacket()
		if encPacket == nil {
			return errors.New("failed to allocate encode packet")
		}
		defer encPacket.Free()

		if err := ctx.encoderCtx.ReceivePacket(encPacket); err != nil {
			if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("receiving packet from encoder failed: %w", err)
		}

		// 重缩放时间戳
		encPacket.RescaleTs(ctx.encoderCtx.TimeBase(), ctx.outputStream.TimeBase())
		encPacket.SetStreamIndex(ctx.outputStream.Index())

		// 写入输出文件
		if err := ctx.outputFormatCtx.WriteInterleavedFrame(encPacket); err != nil {
			return fmt.Errorf("writing frame failed: %w", err)
		}
	}

	return nil
}

// flushCodecs 刷新解码器和编码器
func (ctx *HWTranscodeContext) flushCodecs() error {
	// 刷新解码器
	if err := ctx.decoderCtx.SendPacket(nil); err != nil {
		return fmt.Errorf("flushing decoder failed: %w", err)
	}

	for {
		if err := ctx.decoderCtx.ReceiveFrame(ctx.frame); err != nil {
			if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("receiving frame during flush failed: %w", err)
		}

		if err := ctx.encodeFrame(ctx.frame); err != nil {
			return fmt.Errorf("encoding frame during flush failed: %w", err)
		}
	}

	// 刷新编码器
	if err := ctx.encoderCtx.SendFrame(nil); err != nil {
		return fmt.Errorf("flushing encoder failed: %w", err)
	}

	for {
		encPacket := astiav.AllocPacket()
		if encPacket == nil {
			return errors.New("failed to allocate encode packet during flush")
		}
		defer encPacket.Free()

		if err := ctx.encoderCtx.ReceivePacket(encPacket); err != nil {
			if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
				break
			}
			return fmt.Errorf("receiving packet during flush failed: %w", err)
		}

		encPacket.RescaleTs(ctx.encoderCtx.TimeBase(), ctx.outputStream.TimeBase())
		encPacket.SetStreamIndex(ctx.outputStream.Index())

		if err := ctx.outputFormatCtx.WriteInterleavedFrame(encPacket); err != nil {
			return fmt.Errorf("writing frame during flush failed: %w", err)
		}
	}

	return nil
}

// createHardwareFramesContext 创建硬件帧上下文
func (ctx *HWTranscodeContext) createHardwareFramesContext() error {
	// 如果解码器已经有硬件帧上下文，直接使用
	if ctx.decoderCtx.HardwareFramesContext() != nil {
		return nil
	}

	// 创建新的硬件帧上下文
	hwFramesCtx := astiav.AllocHardwareFramesContext(ctx.hwDeviceCtx)
	if hwFramesCtx == nil {
		return errors.New("failed to allocate hardware frames context")
	}
	
	hwFramesCtx.SetHardwarePixelFormat(ctx.hwPixFmt)
	hwFramesCtx.SetSoftwarePixelFormat(astiav.PixelFormatYuv420P)
	hwFramesCtx.SetWidth(ctx.decoderCtx.Width())
	hwFramesCtx.SetHeight(ctx.decoderCtx.Height())
	hwFramesCtx.SetInitialPoolSize(20)
	
	if err := hwFramesCtx.Initialize(); err != nil {
		return fmt.Errorf("initializing hardware frames context failed: %w", err)
	}
	
	// 设置到解码器
	ctx.decoderCtx.SetHardwareFramesContext(hwFramesCtx)
	ctx.hwFramesCtx = hwFramesCtx
	
	return nil
}

// initHardwareFilter 初始化硬件filter
func (ctx *HWTranscodeContext) initHardwareFilter(filterDesc string) error {
	// 创建filter图
	ctx.filterGraph = astiav.AllocFilterGraph()
	if ctx.filterGraph == nil {
		return errors.New("failed to allocate filter graph")
	}

	// 查找buffer source和sink filter
	buffersrc := astiav.FindFilterByName("buffer")
	if buffersrc == nil {
		return errors.New("buffer filter not found")
	}

	buffersink := astiav.FindFilterByName("buffersink")
	if buffersink == nil {
		return errors.New("buffersink filter not found")
	}

	// 创建buffer source filter context
	var err error
	ctx.buffersrcCtx, err = ctx.filterGraph.NewBuffersrcFilterContext(buffersrc, "in")
	if err != nil {
		return fmt.Errorf("creating buffer source filter failed: %w", err)
	}

	// 设置buffer source参数
	buffersrcParams := astiav.AllocBuffersrcFilterContextParameters()
	if buffersrcParams == nil {
		return errors.New("failed to allocate buffersrc parameters")
	}
	defer buffersrcParams.Free()

	buffersrcParams.SetWidth(ctx.decoderCtx.Width())
	buffersrcParams.SetHeight(ctx.decoderCtx.Height())
	// 使用硬件像素格式 - 直接处理硬件帧
	buffersrcParams.SetPixelFormat(ctx.hwPixFmt)
	
	// 使用输入流的时间基数
	inputStream := ctx.inputFormatCtx.Streams()[ctx.videoStreamIdx]
	timeBase := inputStream.TimeBase()
	if timeBase.Num() == 0 || timeBase.Den() == 0 {
		timeBase = astiav.NewRational(1, 25) // 默认25fps
	}
	buffersrcParams.SetTimeBase(timeBase)
	buffersrcParams.SetSampleAspectRatio(ctx.decoderCtx.SampleAspectRatio())
	
	// 设置硬件帧上下文 - 关键！
	if ctx.decoderCtx.HardwareFramesContext() != nil {
		buffersrcParams.SetHardwareFramesContext(ctx.decoderCtx.HardwareFramesContext())
	}

	if err := ctx.buffersrcCtx.SetParameters(buffersrcParams); err != nil {
		return fmt.Errorf("setting buffersrc parameters failed: %w", err)
	}

	if err := ctx.buffersrcCtx.Initialize(nil); err != nil {
		return fmt.Errorf("initializing buffersrc failed: %w", err)
	}

	// 创建buffer sink filter context
	ctx.buffersinkCtx, err = ctx.filterGraph.NewBuffersinkFilterContext(buffersink, "out")
	if err != nil {
		return fmt.Errorf("creating buffer sink filter failed: %w", err)
	}

	if err := ctx.buffersinkCtx.Initialize(); err != nil {
		return fmt.Errorf("initializing buffersink failed: %w", err)
	}

	// 解析并创建filter链
	inputs := astiav.AllocFilterInOut()
	if inputs == nil {
		return errors.New("failed to allocate filter inputs")
	}
	defer inputs.Free()

	outputs := astiav.AllocFilterInOut()
	if outputs == nil {
		return errors.New("failed to allocate filter outputs")
	}
	defer outputs.Free()

	// 设置inputs和outputs
	outputs.SetName("in")
	outputs.SetFilterContext(ctx.buffersrcCtx.FilterContext())
	outputs.SetPadIdx(0)
	outputs.SetNext(nil)

	inputs.SetName("out")
	inputs.SetFilterContext(ctx.buffersinkCtx.FilterContext())
	inputs.SetPadIdx(0)
	inputs.SetNext(nil)

	// 解析filter描述
	if err := ctx.filterGraph.Parse(filterDesc, inputs, outputs); err != nil {
		return fmt.Errorf("parsing filter description failed: %w", err)
	}

	// 配置filter图
	if err := ctx.filterGraph.Configure(); err != nil {
		return fmt.Errorf("configuring filter graph failed: %w", err)
	}

	// 分配filtered frame
	ctx.filteredFrame = astiav.AllocFrame()
	if ctx.filteredFrame == nil {
		return errors.New("failed to allocate filtered frame")
	}

	return nil
}

// openEncoder 打开编码器
func (ctx *HWTranscodeContext) openEncoder() error {
	encoder := astiav.FindEncoderByName(*encCodec)
	if encoder == nil {
		return fmt.Errorf("encoder not found: %s", *encCodec)
	}

	if err := ctx.encoderCtx.Open(encoder, nil); err != nil {
		return fmt.Errorf("opening encoder failed: %w", err)
	}

	// 复制编码器参数到流
	if err := ctx.outputStream.CodecParameters().FromCodecContext(ctx.encoderCtx); err != nil {
		return fmt.Errorf("copying encoder parameters failed: %w", err)
	}

	// 写入文件头
	if err := ctx.outputFormatCtx.WriteHeader(nil); err != nil {
		return fmt.Errorf("writing header failed: %w", err)
	}

	return nil
}

// updateEncoderForFilter 更新编码器参数以匹配filter输出
func (ctx *HWTranscodeContext) updateEncoderForFilter() error {
	if !ctx.useHardwareFilter || ctx.buffersinkCtx == nil {
		return nil
	}

	// 获取filter输出的分辨率
	width := ctx.buffersinkCtx.Width()
	height := ctx.buffersinkCtx.Height()
	
	if width > 0 && height > 0 {
		// 重新设置编码器分辨率
		ctx.encoderCtx.SetWidth(width)
		ctx.encoderCtx.SetHeight(height)
		
		log.Printf("Updated encoder resolution to %dx%d based on filter output", width, height)
	}

	// 现在打开编码器
	encoder := astiav.FindEncoderByName(*encCodec)
	if encoder == nil {
		return fmt.Errorf("encoder not found: %s", *encCodec)
	}

	if err := ctx.encoderCtx.Open(encoder, nil); err != nil {
		return fmt.Errorf("opening encoder failed: %w", err)
	}

	// 复制编码器参数到流
	if err := ctx.outputStream.CodecParameters().FromCodecContext(ctx.encoderCtx); err != nil {
		return fmt.Errorf("copying encoder parameters failed: %w", err)
	}

	// 写入文件头
	if err := ctx.outputFormatCtx.WriteHeader(nil); err != nil {
		return fmt.Errorf("writing header failed: %w", err)
	}

	return nil
}

// filterFrame 使用硬件filter处理帧
func (ctx *HWTranscodeContext) filterFrame(frame *astiav.Frame) (*astiav.Frame, error) {
	if !ctx.useHardwareFilter {
		return frame, nil
	}

	// 直接发送硬件帧到buffer source - 不转换到CPU
	if err := ctx.buffersrcCtx.AddFrame(frame, astiav.NewBuffersrcFlags()); err != nil {
		return nil, fmt.Errorf("adding frame to buffersrc failed: %w", err)
	}

	// 从buffer sink获取过滤后的硬件帧
	if err := ctx.buffersinkCtx.GetFrame(ctx.filteredFrame, astiav.NewBuffersinkFlags()); err != nil {
		if errors.Is(err, astiav.ErrEagain) || errors.Is(err, astiav.ErrEof) {
			return nil, nil // 没有可用的帧
		}
		return nil, fmt.Errorf("getting frame from buffersink failed: %w", err)
	}

	return ctx.filteredFrame, nil
}

// cleanup 清理资源
func (ctx *HWTranscodeContext) cleanup() {
	if ctx.packet != nil {
		ctx.packet.Free()
	}
	if ctx.frame != nil {
		ctx.frame.Free()
	}
	if ctx.hwFrame != nil {
		ctx.hwFrame.Free()
	}
	if ctx.swFrame != nil {
		ctx.swFrame.Free()
	}
	if ctx.filteredFrame != nil {
		ctx.filteredFrame.Free()
	}
	if ctx.filterGraph != nil {
		ctx.filterGraph.Free()
	}
	if ctx.decoderCtx != nil {
		ctx.decoderCtx.Free()
	}
	if ctx.encoderCtx != nil {
		ctx.encoderCtx.Free()
	}
	if ctx.inputFormatCtx != nil {
		ctx.inputFormatCtx.CloseInput()
		ctx.inputFormatCtx.Free()
	}
	if ctx.outputFormatCtx != nil {
		if ctx.outputFormatCtx.Pb() != nil {
			ctx.outputFormatCtx.Pb().Close()
		}
		ctx.outputFormatCtx.Free()
	}
	if ctx.hwDeviceCtx != nil {
		ctx.hwDeviceCtx.Free()
	}
}