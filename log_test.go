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
			lString: l.String(),
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
			fmt: "warning",
			l:   astiav.LogLevelWarning,
			lString: "warning",
			msg: "warning",
		},
		{
			fmt: "error",
			l:   astiav.LogLevelError,
			lString: "error",
			msg: "error",
		},
		{
			fmt: "fatal",
			l:   astiav.LogLevelFatal,
			lString: "fatal",
			msg: "fatal",
		},
	}, lis)
	astiav.ResetLogCallback()
	lis = []logItem{}
	astiav.Log(astiav.LogLevelError, "test error log\n")
	require.Equal(t, []logItem{}, lis)
}
