package main

import (
	"fmt"
	"os"

	"github.com/asticode/go-astiav"
)

func typeString(entryType astiav.IODirEntryType) string {
	switch entryType {
	case astiav.IODirEntryTypeDirectory:
		return "<DIR>"
	case astiav.IODirEntryTypeFile:
		return "<FILE>"
	case astiav.IODirEntryTypeBlockDevice:
		return "<BLOCK DEVICE>"
	case astiav.IODirEntryTypeCharacterDevice:
		return "<CHARACTER DEVICE>"
	case astiav.IODirEntryTypeNamedPipe:
		return "<PIPE>"
	case astiav.IODirEntryTypeSymbolicLink:
		return "<LINK>"
	case astiav.IODirEntryTypeSocket:
		return "<SOCKET>"
	case astiav.IODirEntryTypeServer:
		return "<SERVER>"
	case astiav.IODirEntryTypeShare:
		return "<SHARE>"
	case astiav.IODirEntryTypeWorkgroup:
		return "<WORKGROUP>"
	case astiav.IODirEntryTypeUnknown:
		fallthrough
	default:
		return "<UNKNOWN>"
	}
}

func listDirectory(inputDir string) error {
	// 打开目录
	dirCtx, err := astiav.OpenDir(inputDir, nil)
	if err != nil {
		fmt.Printf("Cannot open directory: %s. Error: %v\n", inputDir, err)
		return err
	}
	defer dirCtx.CloseDir()

	cnt := 0
	for {
		// 读取目录条目
		entry, err := dirCtx.ReadDir()
		if err != nil {
			fmt.Printf("Cannot list directory: %s. Error: %v\n", inputDir, err)
			return err
		}
		if entry == nil {
			break // 没有更多条目
		}

		var filemode string
		if entry.Filemode() == -1 {
			filemode = "???"
		} else {
			filemode = fmt.Sprintf("%3o", entry.Filemode())
		}

		uidAndGid := fmt.Sprintf("%d(%d)", entry.UserID(), entry.GroupID())

		if cnt == 0 {
			fmt.Printf("%-9s %12s %30s %10s %s %16s %16s %16s\n",
				"TYPE", "SIZE", "NAME", "UID(GID)", "UGO", "MODIFIED",
				"ACCESSED", "STATUS_CHANGED")
		}

		fmt.Printf("%-9s %12d %30s %10s %s %16d %16d %16d\n",
			typeString(entry.Type()),
			entry.Size(),
			entry.Name(),
			uidAndGid,
			filemode,
			entry.ModificationTimestamp(),
			entry.AccessTimestamp(),
			entry.StatusChangeTimestamp())

		// 释放条目
		entry.Free()
		cnt++
	}

	return nil
}

func usage(programName string) {
	fmt.Fprintf(os.Stderr, "usage: %s input_dir\n"+
		"API example program to show how to list files in directory "+
		"accessed through AVIOContext.\n", programName)
}

func main() {
	// 设置日志级别
	astiav.SetLogLevel(astiav.LogLevelDebug)

	if len(os.Args) < 2 {
		usage(os.Args[0])
		os.Exit(1)
	}

	// 初始化网络
	astiav.NetworkInit()
	defer astiav.NetworkDeinit()

	err := listDirectory(os.Args[1])
	if err != nil {
		os.Exit(1)
	}
}