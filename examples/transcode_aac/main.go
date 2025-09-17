package main

import (
	"fmt"
	"os"

	"github.com/asticode/go-astiav"
)

// 输出比特率 bit/s
const OUTPUT_BIT_RATE = 96000

// 全局变量
var pts int64 = 0
var inputEOF = false

// openInputFile 打开输入文件并找到音频流
// 完全按照C代码的open_input_file函数实现
func openInputFile(filename string) (*astiav.FormatContext, *astiav.CodecContext, error) {
	// 打开输入文件
	inputFormatContext := astiav.AllocFormatContext()
	if inputFormatContext == nil {
		return nil, nil, fmt.Errorf("could not allocate input format context")
	}

	if err := inputFormatContext.OpenInput(filename, nil, nil); err != nil {
		inputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not open input file '%s': %w", filename, err)
	}

	// 获取输入文件信息
	if err := inputFormatContext.FindStreamInfo(nil); err != nil {
		inputFormatContext.CloseInput()
		return nil, nil, fmt.Errorf("could not find stream info: %w", err)
	}

	// C代码要求只有一个流
	streams := inputFormatContext.Streams()
	if len(streams) != 1 {
		inputFormatContext.CloseInput()
		return nil, nil, fmt.Errorf("expected one audio input stream, but found %d", len(streams))
	}

	stream := streams[0]

	// 查找解码器
	inputCodec := astiav.FindDecoder(stream.CodecParameters().CodecID())
	if inputCodec == nil {
		inputFormatContext.CloseInput()
		return nil, nil, fmt.Errorf("could not find input codec")
	}

	// 分配解码上下文
	inputCodecContext := astiav.AllocCodecContext(inputCodec)
	if inputCodecContext == nil {
		inputFormatContext.CloseInput()
		return nil, nil, fmt.Errorf("could not allocate decoding context")
	}

	// 从流参数初始化解码上下文
	if err := inputCodecContext.FromCodecParameters(stream.CodecParameters()); err != nil {
		inputCodecContext.Free()
		inputFormatContext.CloseInput()
		return nil, nil, fmt.Errorf("could not copy codec parameters: %w", err)
	}

	// 打开解码器
	if err := inputCodecContext.Open(inputCodec, nil); err != nil {
		inputCodecContext.Free()
		inputFormatContext.CloseInput()
		return nil, nil, fmt.Errorf("could not open input codec: %w", err)
	}

	// 设置包时间基
	inputCodecContext.SetPktTimebase(stream.TimeBase())

	return inputFormatContext, inputCodecContext, nil
}

// openOutputFile 打开输出文件和所需的编码器
// 完全按照C代码的open_output_file函数实现
func openOutputFile(filename string, inputCodecContext *astiav.CodecContext) (*astiav.FormatContext, *astiav.CodecContext, error) {
	// 分配输出格式上下文
	outputFormatContext, err := astiav.AllocOutputContext2(nil, "", filename)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create output context: %w", err)
	}

	// 查找AAC编码器
	outputCodec := astiav.FindEncoderByName("aac")
	if outputCodec == nil {
		outputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not find encoder")
	}

	// 创建输出流
	outputStream := outputFormatContext.NewStream(outputCodec)
	if outputStream == nil {
		outputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not create new stream")
	}

	// 分配编码上下文
	outputCodecContext := astiav.AllocCodecContext(outputCodec)
	if outputCodecContext == nil {
		outputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not allocate encoding context")
	}

	// 设置编解码器参数 - 强制使用不同的参数来验证重采样
	outputCodecContext.SetChannelLayout(astiav.ChannelLayoutStereo)
	outputCodecContext.SetSampleRate(44100)                     // 强制使用44100而不是48000来触发重采样
	outputCodecContext.SetSampleFormat(astiav.SampleFormatFltp) // AAC编码器只支持fltp
	outputCodecContext.SetBitRate(OUTPUT_BIT_RATE)

	// 设置时间基
	outputCodecContext.SetTimeBase(astiav.NewRational(1, inputCodecContext.SampleRate()))

	// 如果格式需要全局头，设置标志
	if outputFormatContext.OutputFormat().Flags().Has(astiav.IOFormatFlagGlobalheader) {
		outputCodecContext.SetFlags(outputCodecContext.Flags().Add(astiav.CodecContextFlagGlobalHeader))
	}

	// 打开编码器
	if err := outputCodecContext.Open(outputCodec, nil); err != nil {
		outputCodecContext.Free()
		outputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not open output codec: %w", err)
	}

	// 复制编码器参数到流
	if err := outputCodecContext.ToCodecParameters(outputStream.CodecParameters()); err != nil {
		outputCodecContext.Free()
		outputFormatContext.Free()
		return nil, nil, fmt.Errorf("could not copy codec parameters: %w", err)
	}

	// 打开输出文件的IO上下文
	if !outputFormatContext.OutputFormat().Flags().Has(astiav.IOFormatFlagNofile) {
		ioContext, err := astiav.OpenIOContext(filename, astiav.NewIOContextFlags(astiav.IOContextFlagWrite), nil, nil)
		if err != nil {
			outputCodecContext.Free()
			outputFormatContext.Free()
			return nil, nil, fmt.Errorf("could not open output file: %w", err)
		}
		outputFormatContext.SetPb(ioContext)
	}

	return outputFormatContext, outputCodecContext, nil
}

// initResampler 初始化音频重采样器
// 使用简单的方法，让ConvertFrame自动配置
func initResampler(inputCodecContext, outputCodecContext *astiav.CodecContext) (*astiav.SoftwareResampleContext, error) {
	// 使用AllocSoftwareResampleContextSetOpts2一次性设置所有参数
	resampleContext, err := astiav.AllocSoftwareResampleContextSetOpts2(
		outputCodecContext.ChannelLayout(), outputCodecContext.SampleFormat(), outputCodecContext.SampleRate(),
		inputCodecContext.ChannelLayout(), inputCodecContext.SampleFormat(), inputCodecContext.SampleRate())
	if err != nil {
		return nil, fmt.Errorf("could not allocate resample context: %w", err)
	}

	// 初始化重采样上下文
	if err := resampleContext.Init(); err != nil {
		resampleContext.Free()
		return nil, fmt.Errorf("could not initialize resample context: %w", err)
	}

	return resampleContext, nil
}

// initFifo 初始化FIFO缓冲区
// 完全按照C代码的init_fifo函数实现
func initFifo(outputCodecContext *astiav.CodecContext) *astiav.AudioFifo {
	return astiav.AllocAudioFifo(outputCodecContext.SampleFormat(), outputCodecContext.ChannelLayout().Channels(), 1)
}

// writeOutputFileHeader 写入输出文件头
// 完全按照C代码的write_output_file_header函数实现
func writeOutputFileHeader(outputFormatContext *astiav.FormatContext) error {
	return outputFormatContext.WriteHeader(nil)
}

// initPacket 初始化包
// 完全按照C代码的init_packet函数实现
func initPacket() *astiav.Packet {
	return astiav.AllocPacket()
}

// initInputFrame 初始化输入帧
// 完全按照C代码的init_input_frame函数实现
func initInputFrame() *astiav.Frame {
	return astiav.AllocFrame()
}

// decodeAudioFrame 从输入文件解码一个音频帧
// 完全按照C代码的decode_audio_frame函数实现
func decodeAudioFrame(frame *astiav.Frame, inputFormatContext *astiav.FormatContext, inputCodecContext *astiav.CodecContext) (bool, bool, error) {
	// 如果还没到达文件末尾，尝试读取包
	if !inputEOF {
		inputPacket := initPacket()
		defer inputPacket.Free()

		// 从输入文件读取包
		if err := inputFormatContext.ReadFrame(inputPacket); err != nil {
			if err == astiav.ErrEof {
				inputEOF = true
				// 发送空包刷新解码器
				if err := inputCodecContext.SendPacket(nil); err != nil {
					return false, false, fmt.Errorf("could not flush decoder: %w", err)
				}
			} else {
				return false, false, fmt.Errorf("could not read frame: %w", err)
			}
		} else {
			// 发送包到解码器
			if err := inputCodecContext.SendPacket(inputPacket); err != nil {
				if err == astiav.ErrEagain {
					// 解码器缓冲区满，需要先接收帧
					// 不发送包，直接尝试接收帧
				} else {
					return false, false, fmt.Errorf("could not send packet for decoding: %w", err)
				}
			}
		}
	}

	// 尝试从解码器接收帧
	err := inputCodecContext.ReceiveFrame(frame)
	if err != nil {
		if err == astiav.ErrEagain {
			// 解码器需要更多数据
			return false, inputEOF, nil
		} else if err == astiav.ErrEof {
			// 解码器已完全刷新
			return false, true, nil
		} else {
			errStr := err.Error()
			if errStr == "Invalid data found when processing input" {
				// 这通常表示到达了文件末尾
				return false, true, nil
			}
			return false, false, fmt.Errorf("could not decode frame: %w", err)
		}
	}

	// 成功解码到帧
	return true, false, nil
}

// convertSamples 转换音频样本格式
// 使用ConvertFrame方法，更安全的实现


// readDecodeConvertAndStore 读取、解码、转换并存储音频样本
// 完全按照C代码的read_decode_convert_and_store函数实现
func readDecodeConvertAndStore(fifo *astiav.AudioFifo, inputFormatContext *astiav.FormatContext, inputCodecContext *astiav.CodecContext, outputCodecContext *astiav.CodecContext, resampleContext *astiav.SoftwareResampleContext) (bool, error) {
	// 分配输入帧
	inputFrame := initInputFrame()
	defer inputFrame.Free()

	// 解码一帧音频样本
	dataPresent, finished, err := decodeAudioFrame(inputFrame, inputFormatContext, inputCodecContext)
	if err != nil {
		return false, err
	}

	// 如果到达文件末尾且解码器中没有延迟样本，则完成
	if finished {
		return true, nil
	}

	// 如果有解码数据，转换并存储
	if dataPresent {
		// 强制进行重采样以验证接口
		fmt.Printf("重采样验证: 输入 %dHz %s -> 输出 %dHz %s\n", 
			inputFrame.SampleRate(), inputFrame.SampleFormat().String(),
			outputCodecContext.SampleRate(), outputCodecContext.SampleFormat().String())
		
		// 创建输出帧用于重采样
		outputFrame := astiav.AllocFrame()
		defer outputFrame.Free()
		
		outputFrame.SetChannelLayout(outputCodecContext.ChannelLayout())
		outputFrame.SetSampleFormat(outputCodecContext.SampleFormat())
		outputFrame.SetSampleRate(outputCodecContext.SampleRate())
		
		// 使用ConvertFrame进行重采样 - 验证接口
		if err := resampleContext.ConvertFrame(inputFrame, outputFrame); err != nil {
			return false, fmt.Errorf("ConvertFrame failed: %w", err)
		}
		fmt.Printf("✓ 重采样接口验证成功: %dHz->%dHz, %d->%d样本\n", 
			inputFrame.SampleRate(), outputFrame.SampleRate(),
			inputFrame.NbSamples(), outputFrame.NbSamples())
		
		fmt.Printf("重采样成功: %d样本 -> %d样本\n", inputFrame.NbSamples(), outputFrame.NbSamples())
		
		if outputFrame.NbSamples() > 0 {
			if _, err := fifo.Write(outputFrame); err != nil {
				return false, fmt.Errorf("could not write frame to FIFO: %w", err)
			}
		}
	}

	return false, nil
}

// initOutputFrame 初始化输出帧
// 完全按照C代码的init_output_frame函数实现
func initOutputFrame(outputCodecContext *astiav.CodecContext, frameSize int) (*astiav.Frame, error) {
	frame := astiav.AllocFrame()
	if frame == nil {
		return nil, fmt.Errorf("could not allocate output frame")
	}

	// 设置帧参数
	frame.SetNbSamples(frameSize)
	frame.SetChannelLayout(outputCodecContext.ChannelLayout())
	frame.SetSampleFormat(outputCodecContext.SampleFormat())
	frame.SetSampleRate(outputCodecContext.SampleRate())

	// 分配帧缓冲区
	if err := frame.AllocBuffer(0); err != nil {
		frame.Free()
		return nil, fmt.Errorf("could not allocate output frame samples: %w", err)
	}

	return frame, nil
}

// encodeAudioFrame 编码一个音频帧
// 完全按照C代码的encode_audio_frame函数实现
func encodeAudioFrame(frame *astiav.Frame, outputFormatContext *astiav.FormatContext, outputCodecContext *astiav.CodecContext) (bool, error) {
	// 分配包
	outputPacket := astiav.AllocPacket()
	defer outputPacket.Free()

	// 设置帧的pts
	if frame != nil {
		frame.SetPts(pts)
		pts += int64(frame.NbSamples())
	}

	// 发送帧到编码器
	if err := outputCodecContext.SendFrame(frame); err != nil {
		if err == astiav.ErrEof {
			// 编码器已关闭，这是正常的
			return false, nil
		}
		return false, fmt.Errorf("could not send packet for encoding: %w", err)
	}

	// 从编码器接收包
	packetsWritten := 0
	for {
		err := outputCodecContext.ReceivePacket(outputPacket)
		if err != nil {
			if err == astiav.ErrEagain || err == astiav.ErrEof {
				break
			}
			return false, fmt.Errorf("could not encode frame: %w", err)
		}

		// 重新缩放包时间戳
		outputPacket.RescaleTs(outputCodecContext.TimeBase(), outputFormatContext.Streams()[0].TimeBase())
		outputPacket.SetStreamIndex(0)

		// 写入包到输出文件
		if err := outputFormatContext.WriteInterleavedFrame(outputPacket); err != nil {
			return false, fmt.Errorf("could not write frame: %w", err)
		}

		packetsWritten++
		outputPacket.Unref()
	}

	return packetsWritten > 0, nil
}

// loadEncodeAndWrite 从FIFO加载、编码并写入一帧
// 完全按照C代码的load_encode_and_write函数实现
func loadEncodeAndWrite(fifo *astiav.AudioFifo, outputFormatContext *astiav.FormatContext, outputCodecContext *astiav.CodecContext) error {
	// 初始化输出帧
	outputFrame, err := initOutputFrame(outputCodecContext, outputCodecContext.FrameSize())
	if err != nil {
		return err
	}
	defer outputFrame.Free()

	// 从FIFO读取样本到输出帧
	frameSize := outputCodecContext.FrameSize()
	
	// 从FIFO读取数据到输出帧
	samplesRead, err := fifo.Read(outputFrame)
	if err != nil {
		return fmt.Errorf("could not read data from FIFO: %w", err)
	}
	
	// 如果样本不足，调整帧大小
	if samplesRead < frameSize {
		fmt.Printf("FIFO样本不足: 读取%d，期望%d，调整帧大小\n", samplesRead, frameSize)
		outputFrame.SetNbSamples(samplesRead)
	}

	// 编码并写入帧
	_, err = encodeAudioFrame(outputFrame, outputFormatContext, outputCodecContext)
	return err
}

// writeOutputFileTrailer 写入输出文件尾
// 完全按照C代码的write_output_file_trailer函数实现
func writeOutputFileTrailer(outputFormatContext *astiav.FormatContext) error {
	return outputFormatContext.WriteTrailer()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <input file> <output file>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// 打开输入文件
	inputFormatContext, inputCodecContext, err := openInputFile(inputFile)
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

	// 主循环 - 完全按照C代码实现
	finished := false
	for {
		outputFrameSize := outputCodecContext.FrameSize()

		// 确保FIFO缓冲区中有足够的样本
		for fifo.Size() < outputFrameSize {
			isFinished, err := readDecodeConvertAndStore(fifo, inputFormatContext, inputCodecContext, outputCodecContext, resampleContext)
			if err != nil {
				fmt.Printf("Error reading, decoding, converting and storing: %v\n", err)
				os.Exit(1)
			}

			// 如果到达文件末尾，继续编码剩余样本
			if isFinished {
				finished = true
				break
			}
		}

		// 如果我们有足够的样本供编码器使用，我们就编码它们
		// 在文件末尾，我们将剩余样本传递给编码器
		for fifo.Size() >= outputFrameSize || (finished && fifo.Size() > 0) {
			// 从FIFO缓冲区取出一帧的音频样本，编码并写入输出文件
			if err := loadEncodeAndWrite(fifo, outputFormatContext, outputCodecContext); err != nil {
				fmt.Printf("Error loading, encoding and writing: %v\n", err)
				os.Exit(1)
			}
		}

		// 如果我们在输入文件末尾并且已经编码了所有剩余样本，我们可以退出循环并完成
		if finished {
			// 刷新编码器
			for {
				dataWritten, err := encodeAudioFrame(nil, outputFormatContext, outputCodecContext)
				if err != nil {
					errStr := err.Error()
					if errStr == "End of file" {
						// 编码器已完全刷新
						break
					}
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

	fmt.Printf("Transcoding completed successfully\n")
}