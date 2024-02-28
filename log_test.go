package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type logItem struct {
	c   Classer
	fmt string
	l   LogLevel
	msg string
}

func TestLog(t *testing.T) {
	var lis []logItem

	SetLogLevel(LogLevelWarning)
	require.Equal(t, LogLevelWarning, GetLogLevel())

	SetLogCallback(func(c Classer, l LogLevel, fmt, msg string) {
		lis = append(lis, logItem{
			c:   c,
			fmt: fmt,
			l:   l,
			msg: msg,
		})
	})
	f := AllocFilterGraph()
	defer f.Free()
	Log(f, LogLevelInfo, "info")
	Log(f, LogLevelWarning, "warning %s", "arg")
	Log(f, LogLevelError, "error")
	Log(f, LogLevelFatal, "fatal")
	require.Equal(t, []logItem{
		{
			c:   f,
			fmt: "warning %s",
			l:   LogLevelWarning,
			msg: "warning arg",
		},
		{
			c:   f,
			fmt: "error",
			l:   LogLevelError,
			msg: "error",
		},
		{
			c:   f,
			fmt: "fatal",
			l:   LogLevelFatal,
			msg: "fatal",
		},
	}, lis)

	ResetLogCallback()
	lis = []logItem{}
	Log(nil, LogLevelError, "test error log\n")
	require.Equal(t, []logItem{}, lis)

	lcs := []Classer{}
	SetLogCallback(func(c Classer, l LogLevel, fmt, msg string) {
		if c != nil {
			lcs = append(lcs, c)
		}
	})
	classers.del(f)
	lcs = []Classer{}
	Log(f, LogLevelWarning, "")
	require.Len(t, lcs, 1)
	require.IsType(t, &UnknownClasser{}, lcs[0])
}
