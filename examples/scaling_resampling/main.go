package main

import (
	"fmt"
	"log"

	"github.com/asticode/go-astiav"
)

func main() {
	fmt.Println("=== 演示缩放和重采样功能 ===")

	// 演示图像缩放功能
	demonstrateScaling()

	// 演示音频重采样功能
	demonstrateResampling()

	// 演示尺寸对齐功能
	demonstrateDimensionAlignment()

	fmt.Println("=== 所有功能演示完成 ===")
}

func demonstrateScaling() {
	fmt.Println("\n--- 图像缩放功能演示 ---")

	// 创建缩放上下文
	srcW, srcH := 1920, 1080
	dstW, dstH := 1280, 720
	srcFormat := astiav.PixelFormatYuv420P
	dstFormat := astiav.PixelFormatYuv420P

	swsCtx, err := astiav.CreateSoftwareScaleContext(
		srcW, srcH, srcFormat,
		dstW, dstH, dstFormat,
		astiav.NewSoftwareScaleContextFlags(astiav.SoftwareScaleContextFlagBicubic),
	)
	if err != nil {
		log.Printf("✗ 创建缩放上下文失败: %v", err)
		return
	}
	defer swsCtx.Free()

	fmt.Printf("✓ 成功创建缩放上下文: %dx%d -> %dx%d\n", srcW, srcH, dstW, dstH)

	// 测试上下文属性
	fmt.Printf("✓ 源分辨率: %dx%d\n", swsCtx.SourceWidth(), swsCtx.SourceHeight())
	fmt.Printf("✓ 目标分辨率: %dx%d\n", swsCtx.DestinationWidth(), swsCtx.DestinationHeight())
	fmt.Printf("✓ 源像素格式: %s\n", swsCtx.SourcePixelFormat().String())
	fmt.Printf("✓ 目标像素格式: %s\n", swsCtx.DestinationPixelFormat().String())

	// 创建测试帧
	srcFrame := astiav.AllocFrame()
	defer srcFrame.Free()
	srcFrame.SetWidth(srcW)
	srcFrame.SetHeight(srcH)
	srcFrame.SetPixelFormat(srcFormat)
	if err := srcFrame.AllocBuffer(1); err != nil {
		log.Printf("✗ 分配源帧缓冲区失败: %v", err)
		return
	}

	dstFrame := astiav.AllocFrame()
	defer dstFrame.Free()

	// 执行缩放
	if err := swsCtx.ScaleFrame(srcFrame, dstFrame); err != nil {
		log.Printf("✗ 缩放帧失败: %v", err)
		return
	}

	fmt.Printf("✓ 成功缩放帧: %dx%d -> %dx%d\n", 
		srcFrame.Width(), srcFrame.Height(),
		dstFrame.Width(), dstFrame.Height())
}

func demonstrateDimensionAlignment() {
	fmt.Println("\n--- 尺寸对齐功能演示 ---")

	// 创建编解码器上下文用于测试
	codec := astiav.FindEncoderByName("libx264")
	if codec == nil {
		log.Printf("✗ 找不到libx264编码器")
		return
	}

	codecCtx := astiav.AllocCodecContext(codec)
	if codecCtx == nil {
		log.Printf("✗ 分配编解码器上下文失败")
		return
	}
	defer codecCtx.Free()

	// 设置基本参数
	codecCtx.SetWidth(1920)
	codecCtx.SetHeight(1080)
	codecCtx.SetPixelFormat(astiav.PixelFormatYuv420P)

	// 测试尺寸对齐
	width, height := 1920, 1080
	fmt.Printf("原始尺寸: %dx%d\n", width, height)

	codecCtx.AlignDimensions(&width, &height)
	fmt.Printf("✓ 对齐后尺寸: %dx%d\n", width, height)

	// 测试高级尺寸对齐
	width2, height2 := 1920, 1080
	linesize := make([]int, 4)
	codecCtx.AlignDimensions2(&width2, &height2, linesize)
	fmt.Printf("✓ 高级对齐后尺寸: %dx%d, linesize: %v\n", width2, height2, linesize)
}

func demonstrateResampling() {
	fmt.Println("\n--- 音频重采样功能演示 ---")

	// 创建重采样上下文
	swrCtx := astiav.AllocSoftwareResampleContext()
	if swrCtx == nil {
		log.Printf("✗ 创建重采样上下文失败")
		return
	}
	defer swrCtx.Free()

	fmt.Printf("✓ 成功创建重采样上下文\n")

	// 创建输入帧
	inputFrame := astiav.AllocFrame()
	defer inputFrame.Free()
	inputFrame.SetChannelLayout(astiav.ChannelLayoutStereo)
	inputFrame.SetSampleFormat(astiav.SampleFormatS16)
	inputFrame.SetSampleRate(44100)
	inputFrame.SetNbSamples(1024)
	if err := inputFrame.AllocBuffer(1); err != nil {
		log.Printf("✗ 分配输入帧缓冲区失败: %v", err)
		return
	}

	// 创建输出帧
	outputFrame := astiav.AllocFrame()
	defer outputFrame.Free()
	outputFrame.SetChannelLayout(astiav.ChannelLayoutStereo)
	outputFrame.SetSampleFormat(astiav.SampleFormatFltp)
	outputFrame.SetSampleRate(48000)
	outputFrame.SetNbSamples(1200) // 大约 1024 * 48000 / 44100

	fmt.Printf("✓ 输入格式: %dHz, %s, %d样本\n", 
		inputFrame.SampleRate(), 
		inputFrame.SampleFormat().String(),
		inputFrame.NbSamples())
	fmt.Printf("✓ 输出格式: %dHz, %s, %d样本\n", 
		outputFrame.SampleRate(), 
		outputFrame.SampleFormat().String(),
		outputFrame.NbSamples())

	// 执行重采样
	if err := swrCtx.ConvertFrame(inputFrame, outputFrame); err != nil {
		log.Printf("✗ 重采样失败: %v", err)
		return
	}

	fmt.Printf("✓ 成功重采样: %dHz -> %dHz\n", 
		inputFrame.SampleRate(), outputFrame.SampleRate())

	// 测试延迟计算
	delay := swrCtx.Delay(int64(outputFrame.SampleRate()))
	fmt.Printf("✓ 重采样延迟: %d 样本\n", delay)
}