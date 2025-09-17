package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/asticode/go-astiav"
)

func processClient(client *astiav.IOContext, inURI string) {
	var input *astiav.IOContext
	buf := make([]byte, 1024)
	var resource string
	
	defer func() {
		fmt.Fprintf(os.Stderr, "Flushing client\n")
		client.Flush()
		fmt.Fprintf(os.Stderr, "Closing client\n")
		client.Close()
		fmt.Fprintf(os.Stderr, "Closing input\n")
		if input != nil {
			input.Close()
		}
	}()
	
	// 执行握手
	for {
		ret := client.Handshake()
		if ret <= 0 {
			if ret < 0 {
				fmt.Printf("Handshake failed: %v\n", ret)
				return
			}
			break
		}
		
		// 获取资源路径
		resourceOpt, err := client.GetOption("resource")
		if err == nil && resourceOpt != "" {
			resource = resourceOpt
			break
		}
	}
	
	if resource == "" {
		fmt.Printf("No resource specified\n")
		return
	}
	
	fmt.Printf("Resource requested: %s\n", resource)
	
	var replyCode int
	if strings.HasPrefix(resource, "/") && resource[1:] == inURI {
		replyCode = 200
	} else {
		replyCode = 404 // HTTP_NOT_FOUND
	}
	
	// 设置回复代码
	err := client.SetOption("reply_code", fmt.Sprintf("%d", replyCode))
	if err != nil {
		fmt.Printf("Failed to set reply_code: %v\n", err)
		return
	}
	
	fmt.Printf("Set reply code to %d\n", replyCode)
	
	// 完成握手
	for {
		ret := client.Handshake()
		if ret <= 0 {
			break
		}
	}
	
	fmt.Fprintf(os.Stderr, "Handshake performed.\n")
	
	if replyCode != 200 {
		return
	}
	
	fmt.Fprintf(os.Stderr, "Opening input file.\n")
	
	// 打开输入文件
	input, err = astiav.OpenIOContext(inURI, astiav.NewIOContextFlags(astiav.IOContextFlagRead), nil, nil)
	if err != nil {
		fmt.Printf("Failed to open input: %s: %v\n", inURI, err)
		return
	}
	
	// 读取并发送文件内容
	for {
		n, err := input.Read(buf)
		if err != nil {
			if err.Error() == "end of file" {
				break
			}
			fmt.Printf("Error reading from input: %v\n", err)
			break
		}
		if n == 0 {
			break
		}
		
		client.Write(buf[:n])
		client.Flush()
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s input http://hostname[:port]\n"+
			"API example program to serve http to multiple clients.\n"+
			"\n", os.Args[0])
		os.Exit(1)
	}
	
	inURI := os.Args[1]
	outURI := os.Args[2]
	
	// 设置日志级别
	astiav.SetLogLevel(astiav.LogLevelDebug)
	
	// 初始化网络
	astiav.NetworkInit()
	defer astiav.NetworkDeinit()
	
	// 创建服务器选项
	options := astiav.NewDictionary()
	defer options.Free()
	
	err := options.Set("listen", "2", astiav.NewDictionaryFlags())
	if err != nil {
		fmt.Printf("Failed to set listen mode for server: %v\n", err)
		os.Exit(1)
	}
	
	// 打开服务器
	server, err := astiav.OpenIOContext(outURI, astiav.NewIOContextFlags(astiav.IOContextFlagWrite), nil, options)
	if err != nil {
		fmt.Printf("Failed to open server: %v\n", err)
		os.Exit(1)
	}
	defer server.Close()
	
	fmt.Fprintf(os.Stderr, "Entering main loop.\n")
	
	// 主循环 - 接受客户端连接
	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Printf("Failed to accept client: %v\n", err)
			break
		}
		
		fmt.Fprintf(os.Stderr, "Accepted client, processing request.\n")
		
		// 在Go中我们不使用fork，而是使用goroutine
		go func() {
			fmt.Fprintf(os.Stderr, "In goroutine.\n")
			processClient(client, inURI)
		}()
	}
	
	fmt.Printf("Server shutting down\n")
}