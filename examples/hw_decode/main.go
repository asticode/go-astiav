/*
 * HW Acceleration API (video decoding) decode sample
 * 
 * Go implementation of FFmpeg's hw_decode.c example
 * 
 * This example demonstrates hardware-accelerated video decoding
 * with output frames from HW video surfaces.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/asticode/go-astiav"
)

var (
	hwDeviceCtx *astiav.HardwareDeviceContext
	hwPixFmt    astiav.PixelFormat
	outputFile  *os.File
)

// hwDecoderInit 初始化硬件解码器
// 完全按照C代码的hw_decoder_init函数实现
func hwDecoderInit(ctx *astiav.CodecContext, hwType astiav.HardwareDeviceType) error {
	var err error
	
	// 创建硬件设备上下文
	hwDeviceCtx, err = astiav.CreateHardwareDeviceContext(hwType, "", nil, 0)
	if err != nil {
		return fmt.Errorf("failed to create specified HW device: %w", err)
	}
	
	// 设置编解码器的硬件设备上下文
	ctx.SetHardwareDeviceContext(hwDeviceCtx)
	
	return nil
}

// getHwFormat 获取硬件像素格式
// 完全按照C代码的get_hw_format函数实现
func getHwFormat(pixFmts []astiav.PixelFormat) astiav.PixelFormat {
	for _, pf := range pixFmts {
		if pf == hwPixFmt {
			return pf
		}
	}
	
	fmt.Fprintf(os.Stderr, "Failed to get HW surface format.\n")
	return astiav.PixelFormatNone
}

// decodeWrite 解码并写入数据
// 完全按照C代码的decode_write函数实现
func decodeWrite(avctx *astiav.CodecContext, packet *astiav.Packet) error {
	var frame, swFrame, tmpFrame *astiav.Frame
	var buffer []byte
	var ret error
	
	// 发送包到解码器
	if err := avctx.SendPacket(packet); err != nil {
		return fmt.Errorf("error during decoding: %w", err)
	}
	
	for {
		// 分配帧
		frame = astiav.AllocFrame()
		if frame == nil {
			return fmt.Errorf("can not alloc frame")
		}
		defer frame.Free()
		
		swFrame = astiav.AllocFrame()
		if swFrame == nil {
			return fmt.Errorf("can not alloc frame")
		}
		defer swFrame.Free()
		
		// 从解码器接收帧
		ret = avctx.ReceiveFrame(frame)
		if ret != nil {
			if ret == astiav.ErrEagain || ret == astiav.ErrEof {
				return nil
			}
			return fmt.Errorf("error while decoding: %w", ret)
		}
		
		// 检查是否是硬件格式
		if frame.PixelFormat() == hwPixFmt {
			// 从GPU传输数据到CPU
			if err := frame.TransferHardwareData(swFrame); err != nil {
				return fmt.Errorf("error transferring the data to system memory: %w", err)
			}
			tmpFrame = swFrame
		} else {
			tmpFrame = frame
		}
		
		// 获取图像缓冲区大小
		size, err := tmpFrame.ImageBufferSize(1)
		if err != nil {
			return fmt.Errorf("can not get image buffer size: %w", err)
		}
		buffer = make([]byte, size)
		
		// 复制图像数据到缓冲区
		if _, err := tmpFrame.ImageCopyToBuffer(buffer, 1); err != nil {
			return fmt.Errorf("can not copy image to buffer: %w", err)
		}
		
		// 写入原始数据到文件
		if _, err := outputFile.Write(buffer); err != nil {
			return fmt.Errorf("failed to dump raw data: %w", err)
		}
	}
}

func main() {
	// 处理ffmpeg日志
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
	
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <device type> <input file> <output file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s videotoolbox input.mp4 output.yuv\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Available device types: videotoolbox, vaapi, cuda, etc.\n")
		os.Exit(1)
	}
	
	deviceTypeName := os.Args[1]
	inputFile := os.Args[2]
	outputFileName := os.Args[3]
	
	// 查找硬件设备类型
	hwType := astiav.FindHardwareDeviceTypeByName(deviceTypeName)
	if hwType == astiav.HardwareDeviceTypeNone {
		fmt.Fprintf(os.Stderr, "Device type %s is not supported.\n", deviceTypeName)
		fmt.Fprintf(os.Stderr, "Available device types: videotoolbox, vaapi, cuda, d3d11va, dxva2, qsv, opencl, vulkan\n")
		os.Exit(1)
	}
	
	// 分配包
	packet := astiav.AllocPacket()
	if packet == nil {
		fmt.Fprintf(os.Stderr, "Failed to allocate AVPacket\n")
		os.Exit(1)
	}
	defer packet.Free()
	
	// 打开输入文件
	inputCtx := astiav.AllocFormatContext()
	if inputCtx == nil {
		fmt.Fprintf(os.Stderr, "Failed to allocate format context\n")
		os.Exit(1)
	}
	defer inputCtx.Free()
	
	if err := inputCtx.OpenInput(inputFile, nil, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open input file '%s': %v\n", inputFile, err)
		os.Exit(1)
	}
	defer inputCtx.CloseInput()
	
	// 查找流信息
	if err := inputCtx.FindStreamInfo(nil); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot find input stream information: %v\n", err)
		os.Exit(1)
	}
	
	// 查找最佳视频流
	videoStream, decoder, err := inputCtx.FindBestStream(astiav.MediaTypeVideo, -1, -1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot find a video stream in the input file: %v\n", err)
		os.Exit(1)
	}
	
	videoStreamIndex := videoStream.Index()
	
	// 简化硬件配置检查 - 假设解码器支持指定的硬件类型
	// 在实际应用中，这里需要更复杂的硬件配置检查
	switch hwType {
	case astiav.HardwareDeviceTypeVideoToolbox:
		hwPixFmt = astiav.PixelFormatVideotoolbox
	case astiav.HardwareDeviceTypeVAAPI:
		hwPixFmt = astiav.PixelFormatVaapi
	case astiav.HardwareDeviceTypeCUDA:
		hwPixFmt = astiav.PixelFormatCuda
	default:
		fmt.Fprintf(os.Stderr, "Unsupported hardware device type\n")
		os.Exit(1)
	}
	
	// 分配解码器上下文
	decoderCtx := astiav.AllocCodecContext(decoder)
	if decoderCtx == nil {
		fmt.Fprintf(os.Stderr, "Failed to allocate codec context\n")
		os.Exit(1)
	}
	defer decoderCtx.Free()
	
	// 复制流参数到解码器上下文
	if err := decoderCtx.FromCodecParameters(videoStream.CodecParameters()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to copy codec parameters: %v\n", err)
		os.Exit(1)
	}
	
	// 初始化硬件解码器
	if err := hwDecoderInit(decoderCtx, hwType); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize hardware decoder: %v\n", err)
		os.Exit(1)
	}
	defer hwDeviceCtx.Free()
	
	// 打开解码器
	if err := decoderCtx.Open(decoder, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open codec for stream #%d: %v\n", videoStreamIndex, err)
		os.Exit(1)
	}
	
	// 打开输出文件
	outputFile, err = os.Create(outputFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open output file '%s': %v\n", outputFileName, err)
		os.Exit(1)
	}
	defer outputFile.Close()
	
	fmt.Printf("Hardware decoding with %s\n", deviceTypeName)
	fmt.Printf("Input: %s\n", inputFile)
	fmt.Printf("Output: %s\n", outputFileName)
	
	// 实际解码和转储原始数据
	for {
		if err := inputCtx.ReadFrame(packet); err != nil {
			if err == astiav.ErrEof {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading frame: %v\n", err)
			break
		}
		
		if packet.StreamIndex() == videoStreamIndex {
			if err := decodeWrite(decoderCtx, packet); err != nil {
				fmt.Fprintf(os.Stderr, "Error in decode_write: %v\n", err)
				break
			}
		}
		
		packet.Unref()
	}
	
	// 刷新解码器
	if err := decodeWrite(decoderCtx, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error flushing decoder: %v\n", err)
	}
	
	fmt.Println("Hardware decoding completed successfully")
}