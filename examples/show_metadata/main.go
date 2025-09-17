package main

import (
	"fmt"
	"os"

	"github.com/asticode/go-astiav"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <input_file>\n"+
			"example program to demonstrate the use of the libavformat metadata API.\n"+
			"\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// 分配格式上下文
	formatCtx := astiav.AllocFormatContext()
	if formatCtx == nil {
		fmt.Printf("Could not allocate format context\n")
		os.Exit(1)
	}
	defer formatCtx.Free()

	// 打开输入文件
	err := formatCtx.OpenInput(inputFile, nil, nil)
	if err != nil {
		fmt.Printf("Could not open input file '%s': %v\n", inputFile, err)
		os.Exit(1)
	}
	defer formatCtx.CloseInput()

	// 查找流信息
	err = formatCtx.FindStreamInfo(nil)
	if err != nil {
		astiav.Log(nil, astiav.LogLevelError, "Cannot find stream information")
		os.Exit(1)
	}

	// 获取格式上下文的元数据
	metadata := formatCtx.Metadata()
	if metadata != nil {
		fmt.Printf("Format metadata:\n")
		printMetadata(metadata)
	}

	// 遍历所有流并显示其元数据
	streams := formatCtx.Streams()
	for i, stream := range streams {
		if stream != nil {
			streamMetadata := stream.Metadata()
			if streamMetadata != nil {
				fmt.Printf("\nStream #%d metadata:\n", i)
				printMetadata(streamMetadata)
			}
		}
	}

	// 获取程序元数据
	programs := formatCtx.Programs()
	for i, program := range programs {
		if program != nil {
			programMetadata := program.Metadata()
			if programMetadata != nil {
				fmt.Printf("\nProgram #%d metadata:\n", i)
				printMetadata(programMetadata)
			}
		}
	}
}

// printMetadata 打印字典中的所有键值对
func printMetadata(dict *astiav.Dictionary) {
	var entry *astiav.DictionaryEntry
	for {
		entry = dict.Iterate(entry)
		if entry == nil {
			break
		}
		fmt.Printf("%s=%s\n", entry.Key(), entry.Value())
	}
}