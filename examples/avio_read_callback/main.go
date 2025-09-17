package main

import (
	"fmt"
	"io"
	"os"

	"github.com/asticode/go-astiav"
)

// BufferData represents the buffer data for custom reading
type BufferData struct {
	data []byte
	pos  int
}

// Read implements the read callback for custom IO
func (bd *BufferData) Read(b []byte) (n int, err error) {
	if bd.pos >= len(bd.data) {
		return 0, io.EOF
	}
	
	n = copy(b, bd.data[bd.pos:])
	bd.pos += n
	
	fmt.Printf("ptr:%p size:%d\n", &bd.data[bd.pos-n], len(bd.data)-bd.pos+n)
	
	return n, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s input_file\n"+
			"API example program to show how to read from a custom buffer "+
			"accessed through AVIOContext.\n", os.Args[0])
		os.Exit(1)
	}
	
	inputFilename := os.Args[1]
	
	// 读取文件内容到缓冲区
	fileData, err := os.ReadFile(inputFilename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
	
	// 创建缓冲区数据结构
	bufferData := &BufferData{
		data: fileData,
		pos:  0,
	}
	
	// 分配格式上下文
	formatCtx := astiav.AllocFormatContext()
	if formatCtx == nil {
		fmt.Printf("Could not allocate format context\n")
		os.Exit(1)
	}
	defer formatCtx.Free()
	
	// 创建自定义IO上下文
	const bufferSize = 4096
	ioCtx, err := astiav.AllocIOContext(bufferSize, false, bufferData.Read, nil, nil)
	if err != nil {
		fmt.Printf("Could not allocate IO context: %v\n", err)
		os.Exit(1)
	}
	defer ioCtx.Free()
	
	// 设置格式上下文的IO上下文
	formatCtx.SetPb(ioCtx)
	
	// 打开输入
	err = formatCtx.OpenInput("", nil, nil)
	if err != nil {
		fmt.Printf("Could not open input: %v\n", err)
		os.Exit(1)
	}
	defer formatCtx.CloseInput()
	
	// 查找流信息
	err = formatCtx.FindStreamInfo(nil)
	if err != nil {
		fmt.Printf("Could not find stream information: %v\n", err)
		os.Exit(1)
	}
	
	// 打印格式信息
	formatCtx.Dump(0, inputFilename, false)
	
	fmt.Printf("Successfully processed file using custom IO callback\n")
}