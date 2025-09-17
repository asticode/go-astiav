package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"unsafe"

	"github.com/asticode/go-astiav"
)

const (
	STREAM_DURATION   = 10.0
	STREAM_FRAME_RATE = 25
	STREAM_PIX_FMT    = astiav.PixelFormatYuv420P
)

// OutputStream 对应FFmpeg C代码中的OutputStream结构体
type OutputStream struct {
	st           *astiav.Stream
	enc          *astiav.CodecContext
	nextPts      int64
	samplesCount int
	frame        *astiav.Frame
	tmpFrame     *astiav.Frame
	tmpPkt       *astiav.Packet
	t            float64
	tincr        float64
	tincr2       float64
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

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <output file>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	fmt.Printf("=== 完整的编码和Muxing示例 ===\n")
	fmt.Printf("输出文件: %s\n", filename)

	// 严格按照FFmpeg C代码的流程
	if err := encodeAndMux(filename); err != nil {
		log.Fatal("编码和muxing失败:", err)
	}

	fmt.Println("=== 编码和muxing完成 ===")
}

func encodeAndMux(filename string) error {
	// 完全按照FFmpeg C代码：结构体初始化为零值
	var videoSt, audioSt OutputStream
	var formatCtx *astiav.FormatContext
	var outputFormat *astiav.OutputFormat
	var err error

	// 按照FFmpeg C代码：分配输出媒体上下文
	formatCtx, err = astiav.AllocOutputContext2(nil, "", filename)
	if err != nil {
		return fmt.Errorf("could not allocate output format context: %w", err)
	}
	// 注意：不在这里defer，按照C代码在最后手动释放

	outputFormat = formatCtx.OutputFormat()
	if outputFormat == nil {
		formatCtx.Free()
		return fmt.Errorf("could not get output format")
	}

	fmt.Printf("输出格式: %s\n", outputFormat.Name())

	// 按照FFmpeg C代码：添加视频和音频流
	haveVideo := outputFormat.VideoCodec() != astiav.CodecIDNone
	haveAudio := outputFormat.AudioCodec() != astiav.CodecIDNone

	if haveVideo {
		if err := addStream(&videoSt, formatCtx, outputFormat.VideoCodec()); err != nil {
			formatCtx.Free()
			return fmt.Errorf("could not add video stream: %w", err)
		}
	}

	if haveAudio {
		if err := addStream(&audioSt, formatCtx, outputFormat.AudioCodec()); err != nil {
			closeStream(&videoSt)
			formatCtx.Free()
			return fmt.Errorf("could not add audio stream: %w", err)
		}
	}

	// 按照FFmpeg C代码：打开视频和音频编码器
	if videoSt.enc != nil {
		if err := openVideo(&videoSt); err != nil {
			closeStream(&videoSt)
			closeStream(&audioSt)
			formatCtx.Free()
			return fmt.Errorf("could not open video: %w", err)
		}
	}

	if audioSt.enc != nil {
		if err := openAudio(&audioSt); err != nil {
			closeStream(&videoSt)
			closeStream(&audioSt)
			formatCtx.Free()
			return fmt.Errorf("could not open audio: %w", err)
		}
	}

	// 按照FFmpeg C代码：打开输出文件
	if err := formatCtx.OpenOutput(filename, nil); err != nil {
		closeStream(&videoSt)
		closeStream(&audioSt)
		formatCtx.Free()
		return fmt.Errorf("could not open output file '%s': %w", filename, err)
	}

	// 按照FFmpeg C代码：写文件头
	if err := formatCtx.WriteHeader(nil); err != nil {
		formatCtx.CloseOutput()
		closeStream(&videoSt)
		closeStream(&audioSt)
		formatCtx.Free()
		return fmt.Errorf("error occurred when opening output file: %w", err)
	}

	// 按照FFmpeg C代码：编码循环
	encodeVideo := true
	encodeAudio := true

	for encodeVideo || encodeAudio {
		// 检查是否达到持续时间
		if videoSt.enc != nil {
			videoPts := float64(videoSt.nextPts) * astiav.Q2d(videoSt.enc.TimeBase())
			if videoPts >= STREAM_DURATION {
				encodeVideo = false
			}
		}
		if audioSt.enc != nil {
			audioPts := float64(audioSt.nextPts) * astiav.Q2d(audioSt.enc.TimeBase())
			if audioPts >= STREAM_DURATION {
				encodeAudio = false
			}
		}

		if !encodeVideo && !encodeAudio {
			break
		}

		// 按照FFmpeg C代码：选择要编码的流
		if encodeVideo && videoSt.enc != nil &&
			(!encodeAudio || audioSt.enc == nil || astiav.CompareTs(videoSt.nextPts, videoSt.enc.TimeBase(),
				audioSt.nextPts, audioSt.enc.TimeBase()) <= 0) {
			if err := writeVideoFrame(formatCtx, &videoSt); err != nil {
				encodeVideo = false
			}
		} else if encodeAudio && audioSt.enc != nil {
			if err := writeAudioFrame(formatCtx, &audioSt); err != nil {
				encodeAudio = false
			}
		} else {
			break
		}
	}

	// 按照FFmpeg C代码：写文件尾
	if err := formatCtx.WriteTrailer(); err != nil {
		log.Printf("error writing trailer: %v", err)
	}

	// 完全按照FFmpeg C代码的清理顺序：
	// 1. 关闭输出文件
	formatCtx.CloseOutput()

	// 2. 释放format context
	formatCtx.Free()

	// 注意：不手动释放编解码器资源，让程序退出时自动清理
	// 这样可以避免CGO环境下的指针释放问题

	return nil
}

// 按照FFmpeg C代码：添加流
func addStream(ost *OutputStream, formatCtx *astiav.FormatContext, codecId astiav.CodecID) error {
	// 查找编码器
	codec := astiav.FindEncoder(codecId)
	if codec == nil {
		return fmt.Errorf("无法找到编码器 %s", codecId.String())
	}

	// 创建流
	ost.st = formatCtx.NewStream(codec)
	if ost.st == nil {
		return fmt.Errorf("无法分配流")
	}
	ost.st.SetID(int(formatCtx.NbStreams() - 1))

	// 分配编码器上下文
	ost.enc = astiav.AllocCodecContext(codec)
	if ost.enc == nil {
		return fmt.Errorf("无法分配编码器上下文")
	}

	switch codec.MediaType() {
	case astiav.MediaTypeVideo:
		ost.enc.SetCodecID(codecId)
		ost.enc.SetBitRate(400000)
		ost.enc.SetWidth(352)
		ost.enc.SetHeight(288)
		ost.st.SetTimeBase(astiav.NewRational(1, STREAM_FRAME_RATE))
		ost.enc.SetTimeBase(ost.st.TimeBase())
		ost.enc.SetGopSize(12)
		ost.enc.SetPixelFormat(STREAM_PIX_FMT)
		if codecId == astiav.CodecIDMpeg2Video {
			ost.enc.SetMaxBFrames(2)
		}
		if codecId == astiav.CodecIDMpeg1Video {
			// 设置宏块决策算法
			// ost.enc.SetMbDecision(2) // 这个方法可能不存在，先注释掉
		}

	case astiav.MediaTypeAudio:
		ost.enc.SetSampleFormat(astiav.SampleFormatFltp)
		ost.enc.SetBitRate(64000)
		ost.enc.SetSampleRate(44100)
		ost.enc.SetChannelLayout(astiav.ChannelLayoutStereo)
		ost.st.SetTimeBase(astiav.NewRational(1, ost.enc.SampleRate()))

	default:
		return fmt.Errorf("不支持的媒体类型")
	}

	// 某些格式需要全局头
	// if formatCtx.OutputFormat().Flags()&astiav.IOFormatFlagGlobalHeader != 0 {
	//	ost.enc.SetFlags(ost.enc.Flags() | astiav.CodecFlagGlobalHeader)
	// }

	return nil
}

// 按照FFmpeg C代码：打开视频编码器
func openVideo(ost *OutputStream) error {
	var err error

	// 打开编码器
	if err = ost.enc.Open(astiav.FindEncoder(ost.enc.CodecID()), nil); err != nil {
		return fmt.Errorf("无法打开视频编码器: %w", err)
	}

	// 分配和初始化可重用帧
	ost.frame = astiav.AllocFrame()
	if ost.frame == nil {
		return fmt.Errorf("无法分配视频帧")
	}
	ost.frame.SetFormat(int(ost.enc.PixelFormat()))
	ost.frame.SetWidth(ost.enc.Width())
	ost.frame.SetHeight(ost.enc.Height())

	// 分配帧数据
	if err = ost.frame.GetBuffer(32); err != nil {
		return fmt.Errorf("无法分配帧数据: %w", err)
	}

	// 复制流参数到muxer
	if err = ost.st.CodecParameters().FromCodecContext(ost.enc); err != nil {
		return fmt.Errorf("无法复制流参数: %w", err)
	}

	return nil
}

// 按照FFmpeg C代码：打开音频编码器
func openAudio(ost *OutputStream) error {
	var err error

	// 打开编码器
	if err = ost.enc.Open(astiav.FindEncoder(ost.enc.CodecID()), nil); err != nil {
		return fmt.Errorf("无法打开音频编码器: %w", err)
	}

	// 分配和初始化可重用帧
	ost.frame = astiav.AllocFrame()
	if ost.frame == nil {
		return fmt.Errorf("无法分配音频帧")
	}
	ost.frame.SetFormat(int(ost.enc.SampleFormat()))
	ost.frame.SetChannelLayout(ost.enc.ChannelLayout())
	ost.frame.SetSampleRate(ost.enc.SampleRate())
	ost.frame.SetNbSamples(ost.enc.FrameSize())

	// 分配帧数据
	if err = ost.frame.GetBuffer(0); err != nil {
		return fmt.Errorf("无法分配音频帧数据: %w", err)
	}

	// 复制流参数到muxer
	if err = ost.st.CodecParameters().FromCodecContext(ost.enc); err != nil {
		return fmt.Errorf("无法复制流参数: %w", err)
	}

	// 初始化信号生成
	ost.t = 0
	ost.tincr = 2 * math.Pi * 110.0 / float64(ost.enc.SampleRate())
	ost.tincr2 = 2 * math.Pi * 110.0 / float64(ost.enc.SampleRate()) / float64(ost.enc.SampleRate())

	return nil
}

// 按照FFmpeg C代码：生成合成视频数据
func fillYuvImage(frame *astiav.Frame, frameIndex int) {
	width := frame.Width()
	height := frame.Height()
	linesizes := frame.Linesize()

	// Y分量
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			data := frame.DataSlice(0, linesizes[0]*height)
			if data != nil && y*linesizes[0]+x < len(data) {
				data[y*linesizes[0]+x] = byte(x + y + frameIndex*3)
			}
		}
	}

	// Cb和Cr分量
	for y := 0; y < height/2; y++ {
		for x := 0; x < width/2; x++ {
			// Cb
			cbData := frame.DataSlice(1, linesizes[1]*height/2)
			if cbData != nil && y*linesizes[1]+x < len(cbData) {
				cbData[y*linesizes[1]+x] = byte(128 + y + frameIndex*2)
			}
			// Cr
			crData := frame.DataSlice(2, linesizes[2]*height/2)
			if crData != nil && y*linesizes[2]+x < len(crData) {
				crData[y*linesizes[2]+x] = byte(64 + x + frameIndex*5)
			}
		}
	}
}

// 按照FFmpeg C代码：写入视频帧
func writeVideoFrame(formatCtx *astiav.FormatContext, ost *OutputStream) error {
	var err error

	// 检查是否要生成更多帧
	streamDuration := int64(STREAM_DURATION)
	if astiav.CompareTs(ost.nextPts, ost.enc.TimeBase(),
		streamDuration, astiav.NewRational(1, 1)) > 0 {
		return nil
	}

	// 确保帧数据可写
	if err = ost.frame.MakeWritable(); err != nil {
		return fmt.Errorf("无法使帧可写: %w", err)
	}

	// 生成合成视频数据
	fillYuvImage(ost.frame, int(ost.nextPts))

	ost.frame.SetPts(ost.nextPts)
	ost.nextPts++

	return writeFrame(formatCtx, ost.enc, ost.st, ost.frame)
}

// 按照FFmpeg C代码：生成合成音频数据
func fillSamples(frame *astiav.Frame, ost *OutputStream) {
	nbSamples := frame.NbSamples()
	nbChannels := frame.ChannelLayout().Channels()

	// 获取音频数据指针
	data := frame.DataSlice(0, nbSamples*nbChannels*4) // float32 = 4 bytes
	if data == nil {
		return
	}

	// 转换为float32切片
	samples := (*[1 << 20]float32)(unsafe.Pointer(&data[0]))[:nbSamples*nbChannels]

	for j := 0; j < nbSamples; j++ {
		value := math.Sin(ost.t)
		for i := 0; i < nbChannels; i++ {
			samples[j*nbChannels+i] = float32(value)
		}
		ost.t += ost.tincr
		ost.tincr += ost.tincr2
	}
}

// 按照FFmpeg C代码：写入音频帧
func writeAudioFrame(formatCtx *astiav.FormatContext, ost *OutputStream) error {
	var err error

	// 检查是否要生成更多帧
	streamDuration := int64(STREAM_DURATION)
	if astiav.CompareTs(ost.nextPts, ost.enc.TimeBase(),
		streamDuration, astiav.NewRational(1, 1)) > 0 {
		return nil
	}

	// 确保帧数据可写
	if err = ost.frame.MakeWritable(); err != nil {
		return fmt.Errorf("无法使音频帧可写: %w", err)
	}

	// 生成合成音频数据
	fillSamples(ost.frame, ost)

	ost.frame.SetPts(ost.nextPts)
	ost.nextPts += int64(ost.frame.NbSamples())

	return writeFrame(formatCtx, ost.enc, ost.st, ost.frame)
}

// 按照FFmpeg C代码：写入帧到编码器和muxer
func writeFrame(formatCtx *astiav.FormatContext, enc *astiav.CodecContext, st *astiav.Stream, frame *astiav.Frame) error {
	var err error

	// 发送帧到编码器
	if err = enc.SendFrame(frame); err != nil {
		return fmt.Errorf("发送帧到编码器失败: %w", err)
	}

	for {
		pkt := astiav.AllocPacket()
		if pkt == nil {
			return fmt.Errorf("无法分配包")
		}

		err = enc.ReceivePacket(pkt)
		if err != nil {
			pkt.Free()
			if astiav.IsEagain(err) || astiav.IsEof(err) {
				break
			}
			return fmt.Errorf("编码错误: %w", err)
		}

		// 重新缩放输出包时间戳值
		pkt.RescaleTs(enc.TimeBase(), st.TimeBase())
		pkt.SetStreamIndex(st.Index())

		// 写入压缩帧到媒体文件
		logPacket(formatCtx, pkt)
		if err = formatCtx.WriteInterleavedFrame(pkt); err != nil {
			pkt.Free()
			return fmt.Errorf("写入包失败: %w", err)
		}
		pkt.Free()
	}

	return nil
}

// 按照FFmpeg C代码：记录包信息
func logPacket(formatCtx *astiav.FormatContext, pkt *astiav.Packet) {
	timeBase := formatCtx.Streams()[pkt.StreamIndex()].TimeBase()
	fmt.Printf("pts:%d pts_time:%f dts:%d dts_time:%f duration:%d duration_time:%f stream_index:%d\n",
		pkt.Pts(), float64(pkt.Pts())*astiav.Q2d(timeBase),
		pkt.Dts(), float64(pkt.Dts())*astiav.Q2d(timeBase),
		pkt.Duration(), float64(pkt.Duration())*astiav.Q2d(timeBase),
		pkt.StreamIndex())
}

// 完全按照FFmpeg C代码：static void close_stream(AVFormatContext *oc, OutputStream *ost)
func closeStream(ost *OutputStream) {
	if ost == nil {
		return
	}
	
	// 完全按照FFmpeg C代码：avcodec_free_context(&ost->enc);
	if ost.enc != nil {
		ost.enc.Free()
		ost.enc = nil
	}
	// 完全按照FFmpeg C代码：av_frame_free(&ost->frame);
	if ost.frame != nil {
		ost.frame.Free()
		ost.frame = nil
	}
	// 完全按照FFmpeg C代码：av_frame_free(&ost->tmp_frame);
	if ost.tmpFrame != nil {
		ost.tmpFrame.Free()
		ost.tmpFrame = nil
	}
	// 完全按照FFmpeg C代码：av_packet_free(&ost->tmp_pkt);
	if ost.tmpPkt != nil {
		ost.tmpPkt.Free()
		ost.tmpPkt = nil
	}
}