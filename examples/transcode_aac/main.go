package main

import (
	"fmt"
	"os"

	"github.com/asticode/go-astiav"
)

// 输出比特率 (bit/s)
const OUTPUT_BIT_RATE = 96000

// 全局时间戳
var pts int64 = 0

// openInputFile 打开输入文件和所需的解码器
func openInputFile(filename string) (*astiav.FormatContext, *astiav.CodecContext, int, error) {
	// 打开输入文件
	inputFormatContext := astiav.AllocFormatContext()
	if inputFormatContext == nil {
		return nil, nil, -1, fmt.Errorf("could not allocate input format context")
	}

	if err := inputFormatContext.OpenInput(filename, nil, nil); err != nil {
		inputFormatContext.Free()
		return nil, nil, -1, fmt.Errorf("could not open input file '%s': %w", filename, err)
	}

	// 获取输入文件信息
	if err := inputFormatContext.FindStreamInfo(nil); err != nil {
		inputFormatContext.CloseInput()
		return nil, nil, -1, fmt.Errorf("could not find stream info: %w", err)
	}

	// 查找音频流
	streams := inputFormatContext.Streams()
	var stream *astiav.Stream
	var audioStreamIndex int = -1
	for i, s := range streams {
		if s.CodecParameters().MediaType() == astiav.MediaTypeAudio {
			stream = s
			audioStreamIndex = i
			break
		}
	}

	if stream == nil {
		inputFormatContext.CloseInput()
		return nil, nil, -1, fmt.Errorf("no audio stream found in input file")
	}

	// 查找解码器
	inputCodec := astiav.FindDecoder(stream.CodecParameters().CodecID())
	if inputCodec == nil {
		inputFormatContext.CloseInput()
		return nil, nil, -1, fmt.Errorf("could not find input codec")
	}

	// 分配解码上下文
	inputCodecContext := astiav.AllocCodecContext(inputCodec)
	if inputCodecContext == nil {
		inputFormatContext.CloseInput()
		return nil, nil, -1, fmt.Errorf("could not allocate decoding context")
	}

	// 从流参数初始化解码上下文
	if err := inputCodecContext.FromCodecParameters(stream.CodecParameters()); err != nil {
		inputCodecContext.Free()
		inputFormatContext.CloseInput()
		return nil, nil, -1, fmt.Errorf("could not copy codec parameters: %w", err)
	}

	// 打开解码器
	if err := inputCodecContext.Open(inputCodec, nil); err != nil {
		inputCodecContext.Free()
		inputFormatContext.CloseInput()
		return nil, nil, -1, fmt.Errorf("could not open input codec: %w", err)
	}

	// 设置包时间基
	inputCodecContext.SetPktTimebase(stream.TimeBase())

	return inputFormatContext, inputCodecContext, audioStreamIndex, nil
}

// openOutputFile 打开输出文件和所需的编码器
func openOutputFile(filename string, inputCodecContext *astiav.CodecContext) (*astiav.FormatContext, *astiav.CodecContext, error) {
	// 使用AllocOutputContext2创建输出上下文，自动设置格式
	outputFormatContext, err := astiav.AllocOutputContext2(nil, "", filename)
	if err != nil {
		return nil, nil, fmt.Errorf("could not allocate output context: %w", err)
	}

	// 打开输出文件
	outputIOContext, err := astiav.OpenIOContext(filename, astiav.NewIOContextFlags(astiav.IOContextFlagWrite), nil, nil)
	if err != nil {
		outputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not open output file '%s': %w", filename, err)
	}

	// 关联输出文件和容器格式上下文
	outputFormatContext.SetPb(outputIOContext)

	// 查找AAC编码器
	outputCodec := astiav.FindEncoder(astiav.CodecIDAac)
	if outputCodec == nil {
		outputIOContext.Close()
		outputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not find an AAC encoder")
	}

	// 创建新的音频流
	stream := outputFormatContext.NewStream(nil)
	if stream == nil {
		outputIOContext.Close()
		outputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not create new stream")
	}

	// 分配编码上下文
	outputCodecContext := astiav.AllocCodecContext(outputCodec)
	if outputCodecContext == nil {
		outputIOContext.Close()
		outputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not allocate encoding context")
	}

	// 设置基本编码参数
	channelLayout := astiav.ChannelLayoutStereo
	outputCodecContext.SetChannelLayout(channelLayout)
	outputCodecContext.SetSampleRate(inputCodecContext.SampleRate())
	outputCodecContext.SetSampleFormat(outputCodec.SampleFormats()[0])
	outputCodecContext.SetBitRate(OUTPUT_BIT_RATE)

	// 设置容器的采样率
	stream.SetTimeBase(astiav.NewRational(1, inputCodecContext.SampleRate()))

	// 某些容器格式需要全局头
	if outputFormatContext.OutputFormat().Flags().Has(astiav.IOFormatFlagGlobalheader) {
		outputCodecContext.SetFlags(outputCodecContext.Flags().Add(astiav.CodecContextFlagGlobalHeader))
	}

	// 打开编码器
	if err := outputCodecContext.Open(outputCodec, nil); err != nil {
		outputCodecContext.Free()
		outputIOContext.Close()
		outputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not open output codec: %w", err)
	}

	// 从编码上下文复制参数到流
	if err := stream.CodecParameters().FromCodecContext(outputCodecContext); err != nil {
		outputCodecContext.Free()
		outputIOContext.Close()
		outputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not initialize stream parameters: %w", err)
	}

	return outputFormatContext, outputCodecContext, nil
}

// initResampler 初始化音频重采样器
func initResampler(inputCodecContext, outputCodecContext *astiav.CodecContext) (*astiav.SoftwareResampleContext, error) {
	// 使用swr_alloc_set_opts2创建并设置重采样上下文
	resampleContext, err := astiav.AllocSoftwareResampleContextSetOpts2(
		outputCodecContext.ChannelLayout(), outputCodecContext.SampleFormat(), outputCodecContext.SampleRate(),
		inputCodecContext.ChannelLayout(), inputCodecContext.SampleFormat(), inputCodecContext.SampleRate())
	if err != nil {
		return nil, fmt.Errorf("could not allocate and set resample context: %w", err)
	}

	// 初始化重采样器
	if err := resampleContext.Init(); err != nil {
		resampleContext.Free()
		return nil, fmt.Errorf("could not initialize resample context: %w", err)
	}

	return resampleContext, nil
}

// initFifo 初始化音频FIFO缓冲区
func initFifo(outputCodecContext *astiav.CodecContext) *astiav.AudioFifo {
	return astiav.AllocAudioFifo(outputCodecContext.SampleFormat(), outputCodecContext.ChannelLayout().Channels(), 1)
}

// writeOutputFileHeader 写入输出文件头
func writeOutputFileHeader(outputFormatContext *astiav.FormatContext) error {
	if err := outputFormatContext.WriteHeader(nil); err != nil {
		return fmt.Errorf("could not write output file header: %w", err)
	}
	return nil
}

// decodeAudioFrame 从输入文件解码一个音频帧
func decodeAudioFrame(frame *astiav.Frame, inputFormatContext *astiav.FormatContext, inputCodecContext *astiav.CodecContext, audioStreamIndex int) (bool, bool, error) {
	// 分配包
	inputPacket := astiav.AllocPacket()
	defer inputPacket.Free()

	dataPresent := false
	finished := false

	// 循环读取直到找到音频包或到达文件末尾
	for {
		// 从输入文件读取一个帧
		if err := inputFormatContext.ReadFrame(inputPacket); err != nil {
			if err.Error() == "end of file" {
				finished = true
				break
			} else {
				return false, false, fmt.Errorf("could not read frame: %w", err)
			}
		}

		// 只处理音频流的包
		if inputPacket.StreamIndex() == audioStreamIndex {
			break
		}
	}

	// 发送包到解码器
	if !finished {
		if err := inputCodecContext.SendPacket(inputPacket); err != nil {
			return false, false, fmt.Errorf("could not send packet for decoding: %w", err)
		}
	} else {
		// 发送空包以刷新解码器
		if err := inputCodecContext.SendPacket(nil); err != nil {
			return false, false, fmt.Errorf("could not send flush packet: %w", err)
		}
	}

	// 从解码器接收帧
	err := inputCodecContext.ReceiveFrame(frame)
	if err != nil {
		if err.Error() == "resource temporarily unavailable" {
			// 解码器需要更多数据，这是正常的
			return false, finished, nil
		} else if err.Error() == "end of file" {
			finished = true
			return false, finished, nil
		} else {
			return false, false, fmt.Errorf("could not decode frame: %w", err)
		}
	} else {
		// 成功解码到帧
		dataPresent = true
		return dataPresent, finished, nil
	}
}

// initConvertedSamples 初始化转换样本的临时存储
func initConvertedSamples(outputCodecContext *astiav.CodecContext, frameSize int) ([][]byte, error) {
	samples, _, err := astiav.SamplesAllocArrayAndSamples(
		outputCodecContext.ChannelLayout().Channels(),
		frameSize,
		outputCodecContext.SampleFormat(),
		0)
	if err != nil {
		return nil, fmt.Errorf("could not allocate converted input samples: %w", err)
	}
	return samples, nil
}

// convertSamples 转换音频样本格式
func convertSamples(inputData [][]byte, convertedData [][]byte, frameSize int, resampleContext *astiav.SoftwareResampleContext) error {
	// 使用重采样器转换样本
	_, err := resampleContext.Convert(convertedData, frameSize, inputData, frameSize)
	if err != nil {
		return fmt.Errorf("could not convert input samples: %w", err)
	}
	return nil
}

// addSamplesToFifo 将转换后的音频样本添加到FIFO缓冲区
func addSamplesToFifo(fifo *astiav.AudioFifo, convertedInputSamples [][]byte, frameSize int) error {
	// 扩大FIFO以容纳新样本
	if err := fifo.Realloc(fifo.Size() + frameSize); err != nil {
		return fmt.Errorf("could not reallocate FIFO: %w", err)
	}

	// 创建临时帧来存储转换后的样本
	tempFrame := astiav.AllocFrame()
	defer tempFrame.Free()

	// 设置帧参数
	tempFrame.SetNbSamples(frameSize)
	tempFrame.SetChannelLayout(astiav.ChannelLayoutStereo) // 假设立体声
	tempFrame.SetSampleFormat(astiav.SampleFormatFltp)     // 假设浮点格式

	// 分配帧缓冲区
	if err := tempFrame.AllocBuffer(0); err != nil {
		return fmt.Errorf("could not allocate temp frame buffer: %w", err)
	}

	// 复制数据到帧
	for i, channelData := range convertedInputSamples {
		if len(channelData) > 0 {
			planeData, err := tempFrame.Data().BytesForPlane(i)
			if err == nil && len(planeData) >= len(channelData) {
				copy(planeData, channelData)
			}
		}
	}

	// 将样本写入FIFO
	if written, err := fifo.Write(tempFrame); err != nil || written < frameSize {
		return fmt.Errorf("could not write data to FIFO")
	}

	return nil
}

// readDecodeConvertAndStore 读取、解码、转换并存储音频样本
func readDecodeConvertAndStore(fifo *astiav.AudioFifo, inputFormatContext *astiav.FormatContext, inputCodecContext *astiav.CodecContext, outputCodecContext *astiav.CodecContext, resampleContext *astiav.SoftwareResampleContext, audioStreamIndex int) (bool, error) {
	// 分配输入帧
	inputFrame := astiav.AllocFrame()
	defer inputFrame.Free()

	// 解码一帧音频样本
	dataPresent, finished, err := decodeAudioFrame(inputFrame, inputFormatContext, inputCodecContext, audioStreamIndex)
	if err != nil {
		return false, err
	}

	// 如果到达文件末尾且解码器中没有延迟样本，则完成
	if finished {
		return true, nil
	}

	// 如果有解码数据，转换并存储
	if dataPresent {
		// 初始化转换样本的临时存储
		convertedInputSamples, err := initConvertedSamples(outputCodecContext, inputFrame.NbSamples())
		if err != nil {
			return false, err
		}

		// 获取输入帧数据
		inputData := make([][]byte, inputCodecContext.ChannelLayout().Channels())
		bytesPerSample := inputCodecContext.SampleFormat().BytesPerSample()
		sampleSize := inputFrame.NbSamples() * bytesPerSample

		for i := 0; i < len(inputData); i++ {
			planeData, err := inputFrame.Data().BytesForPlane(i)
			if err == nil && len(planeData) >= sampleSize {
				inputData[i] = planeData[:sampleSize]
			}
		}

		// 转换输入样本到所需的输出样本格式
		if err := convertSamples(inputData, convertedInputSamples, inputFrame.NbSamples(), resampleContext); err != nil {
			return false, err
		}

		// 将转换后的输入样本添加到FIFO缓冲区
		if err := addSamplesToFifo(fifo, convertedInputSamples, inputFrame.NbSamples()); err != nil {
			return false, err
		}
	}

	return false, nil
}

// initOutputFrame 初始化输出帧
func initOutputFrame(outputCodecContext *astiav.CodecContext, frameSize int) (*astiav.Frame, error) {
	// 创建新帧
	frame := astiav.AllocFrame()
	if frame == nil {
		return nil, fmt.Errorf("could not allocate output frame")
	}

	// 设置帧参数
	frame.SetNbSamples(frameSize)
	frame.SetChannelLayout(outputCodecContext.ChannelLayout())
	frame.SetSampleFormat(outputCodecContext.SampleFormat())
	frame.SetSampleRate(outputCodecContext.SampleRate())

	// 分配样本缓冲区
	if err := frame.AllocBuffer(0); err != nil {
		frame.Free()
		return nil, fmt.Errorf("could not allocate output frame samples: %w", err)
	}

	return frame, nil
}

// encodeAudioFrame 编码一个音频帧
func encodeAudioFrame(frame *astiav.Frame, outputFormatContext *astiav.FormatContext, outputCodecContext *astiav.CodecContext) (bool, error) {
	// 分配输出包
	outputPacket := astiav.AllocPacket()
	defer outputPacket.Free()

	dataPresent := false

	// 设置时间戳
	if frame != nil {
		frame.SetPts(pts)
		pts += int64(frame.NbSamples())
	}

	// 发送帧到编码器
	err := outputCodecContext.SendFrame(frame)
	if err != nil && err.Error() != "end of file" {
		return false, fmt.Errorf("could not send packet for encoding: %w", err)
	}

	// 从编码器接收包
	err = outputCodecContext.ReceivePacket(outputPacket)
	if err != nil {
		if err.Error() == "resource temporarily unavailable" {
			// 编码器需要更多数据
			return false, nil
		} else if err.Error() == "end of file" {
			// 最后一帧已编码
			return false, nil
		} else {
			return false, fmt.Errorf("could not encode frame: %w", err)
		}
	}

	dataPresent = true

	// 写入一个音频帧到输出文件
	if dataPresent {
		if err := outputFormatContext.WriteFrame(outputPacket); err != nil {
			return false, fmt.Errorf("could not write frame: %w", err)
		}
	}

	return dataPresent, nil
}

// loadEncodeAndWrite 从FIFO缓冲区加载、编码并写入音频帧
func loadEncodeAndWrite(fifo *astiav.AudioFifo, outputFormatContext *astiav.FormatContext, outputCodecContext *astiav.CodecContext) error {
	// 使用编码器的帧大小
	frameSize := outputCodecContext.FrameSize()
	if frameSize <= 0 {
		frameSize = 1024 // 默认帧大小
	}

	// 初始化输出帧
	outputFrame, err := initOutputFrame(outputCodecContext, frameSize)
	if err != nil {
		return err
	}
	defer outputFrame.Free()

	// 从FIFO缓冲区读取样本
	if read, err := fifo.Read(outputFrame); err != nil || read < frameSize {
		return fmt.Errorf("could not read data from FIFO")
	}

	// 编码一帧音频样本
	if _, err := encodeAudioFrame(outputFrame, outputFormatContext, outputCodecContext); err != nil {
		return err
	}

	return nil
}

// writeOutputFileTrailer 写入输出文件尾
func writeOutputFileTrailer(outputFormatContext *astiav.FormatContext) error {
	if err := outputFormatContext.WriteTrailer(); err != nil {
		return fmt.Errorf("could not write output file trailer: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <input file> <output file>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// 打开输入文件
	inputFormatContext, inputCodecContext, audioStreamIndex, err := openInputFile(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		inputCodecContext.Free()
		inputFormatContext.CloseInput()
	}()

	// 打开输出文件
	outputFormatContext, outputCodecContext, err := openOutputFile(outputFile, inputCodecContext)
	if err != nil {
		fmt.Printf("Error opening output file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		outputCodecContext.Free()
		if outputFormatContext.Pb() != nil {
			outputFormatContext.Pb().Close()
		}
		outputFormatContext.Free()
	}()

	// 初始化重采样器
	resampleContext, err := initResampler(inputCodecContext, outputCodecContext)
	if err != nil {
		fmt.Printf("Error initializing resampler: %v\n", err)
		os.Exit(1)
	}
	defer resampleContext.Free()

	// 初始化FIFO缓冲区
	fifo := initFifo(outputCodecContext)
	defer fifo.Free()

	// 写入输出文件头
	if err := writeOutputFileHeader(outputFormatContext); err != nil {
		fmt.Printf("Error writing output file header: %v\n", err)
		os.Exit(1)
	}

	// 主循环 - 完全按照C代码的逻辑
	for {
		outputFrameSize := outputCodecContext.FrameSize()
		if outputFrameSize <= 0 {
			outputFrameSize = 1024 // 默认帧大小
		}
		finished := false

		// 确保FIFO缓冲区中有足够的样本
		for fifo.Size() < outputFrameSize {
			finished, err := readDecodeConvertAndStore(fifo, inputFormatContext, inputCodecContext, outputCodecContext, resampleContext, audioStreamIndex)
			if err != nil {
				fmt.Printf("Error reading, decoding, converting and storing: %v\n", err)
				os.Exit(1)
			}

			// 如果到达文件末尾，继续编码剩余样本
			if finished {
				break
			}
		}

		// 如果有足够的样本或到达文件末尾，进行编码
		for fifo.Size() >= outputFrameSize || (finished && fifo.Size() > 0) {
			if err := loadEncodeAndWrite(fifo, outputFormatContext, outputCodecContext); err != nil {
				fmt.Printf("Error loading, encoding and writing: %v\n", err)
				os.Exit(1)
			}
		}

		// 如果到达文件末尾且已编码所有剩余样本，退出循环
		if finished {
			// 刷新编码器
			for {
				dataWritten, err := encodeAudioFrame(nil, outputFormatContext, outputCodecContext)
				if err != nil {
					fmt.Printf("Error flushing encoder: %v\n", err)
					os.Exit(1)
				}
				if !dataWritten {
					break
				}
			}
			break
		}
	}

	// 写入输出文件尾
	if err := writeOutputFileTrailer(outputFormatContext); err != nil {
		fmt.Printf("Error writing output file trailer: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Transcoding completed successfully!\n")
}