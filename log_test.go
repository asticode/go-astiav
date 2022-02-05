package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

type logItem struct {
	l   astiav.LogLevel
	msg string
}

func TestLog(t *testing.T) {
	var lis []logItem
	astiav.SetLogCallback(func(l astiav.LogLevel, msg, parent string) {
		lis = append(lis, logItem{
			l:   l,
			msg: msg,
		})
	})
	astiav.SetLogLevel(astiav.LogLevelWarning)
	astiav.Log(astiav.LogLevelInfo, "info")
	astiav.Log(astiav.LogLevelWarning, "warning")
	astiav.Log(astiav.LogLevelError, "error")
	astiav.Log(astiav.LogLevelFatal, "fatal")
	require.Equal(t, []logItem{
		{
			l:   astiav.LogLevelWarning,
			msg: "warning",
		},
		{
			l:   astiav.LogLevelError,
			msg: "error",
		},
		{
			l:   astiav.LogLevelFatal,
			msg: "fatal",
		},
	}, lis)
	astiav.ResetLogCallback()
	lis = []logItem{}
	astiav.Log(astiav.LogLevelError, "test error log\n")
	require.Equal(t, []logItem{}, lis)
}
