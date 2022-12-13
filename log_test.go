package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

type logItem struct {
	fmt string
	l   astiav.LogLevel
	msg string
}

func TestLog(t *testing.T) {
	var lis []logItem
	astiav.SetLogCallback(func(l astiav.LogLevel, fmt, msg, parent string) {
		lis = append(lis, logItem{
			fmt: fmt,
			l:   l,
			msg: msg,
		})
	})
	astiav.SetLogLevel(astiav.LogLevelWarning)
	astiav.Log(astiav.LogLevelDebug, "debug")
	astiav.Log(astiav.LogLevelVerbose, "verbose")
	astiav.Log(astiav.LogLevelInfo, "info")
	astiav.Log(astiav.LogLevelWarning, "warning")
	astiav.Log(astiav.LogLevelError, "error")
	astiav.Log(astiav.LogLevelFatal, "fatal")
	require.Equal(t, []logItem{
		{
			fmt: "warning",
			l:   astiav.LogLevelWarning,
			msg: "warning",
		},
		{
			fmt: "error",
			l:   astiav.LogLevelError,
			msg: "error",
		},
		{
			fmt: "fatal",
			l:   astiav.LogLevelFatal,
			msg: "fatal",
		},
	}, lis)
	astiav.ResetLogCallback()
	lis = []logItem{}
	astiav.Log(astiav.LogLevelError, "test error log\n")
	require.Equal(t, []logItem{}, lis)
}

func TestLogf(t *testing.T) {
	var lis []logItem
	astiav.SetLogCallback(func(l astiav.LogLevel, fmt, msg, parent string) {
		lis = append(lis, logItem{
			fmt: fmt,
			l:   l,
			msg: msg,
		})
	})
	astiav.SetLogLevel(astiav.LogLevelWarning)
	astiav.Logf(astiav.LogLevelDebug, "debug %s %d %.3f", "s", 1, 2.0)
	astiav.Logf(astiav.LogLevelVerbose, "verbose %s %d %.3f", "s", 1, 2.0)
	astiav.Logf(astiav.LogLevelInfo, "info %s %d %.3f", "s", 1, 2.0)
	astiav.Logf(astiav.LogLevelWarning, "warning %s %d %.3f", "s", 1, 2.0)
	astiav.Logf(astiav.LogLevelError, "error %s %d %.3f", "s", 1, 2.0)
	astiav.Logf(astiav.LogLevelFatal, "fatal %s %d %.3f", "s", 1, 2.0)
	for i, l := range []logItem{
		{
			fmt: "warning s 1 2.000",
			l:   astiav.LogLevelWarning,
			msg: "warning s 1 2.000",
		},
		{
			fmt: "error s 1 2.000",
			l:   astiav.LogLevelError,
			msg: "error s 1 2.000",
		},
		{
			fmt: "fatal s 1 2.000",
			l:   astiav.LogLevelFatal,
			msg: "fatal s 1 2.000",
		},
	} {
		require.Equal(t, l, lis[i])
	}
	astiav.ResetLogCallback()
	lis = []logItem{}
	astiav.Log(astiav.LogLevelError, "test error log\n")
	require.Equal(t, []logItem{}, lis)
}
