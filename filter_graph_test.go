// TODO Fix https://github.com/asticode/go-astiav/actions/runs/5853322732/job/15867145888
//go:build !windows

package astiav

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterGraph(t *testing.T) {
	fg := AllocFilterGraph()
	defer fg.Free()
	cl := fg.Class()
	require.NotNil(t, cl)
	require.Equal(t, "AVFilterGraph", cl.Name())
	fg.SetThreadCount(2)
	require.Equal(t, 2, fg.ThreadCount())
	fg.SetThreadType(ThreadTypeSlice)
	require.Equal(t, ThreadTypeSlice, fg.ThreadType())

	bufferSink := FindFilterByName("buffersink")
	require.NotNil(t, bufferSink)

	fcOut, err := fg.NewFilterContext(bufferSink, "filter_out", nil)
	require.NoError(t, err)
	defer fcOut.Free()
	cl = fcOut.Class()
	require.NotNil(t, cl)
	require.Equal(t, "AVFilter", cl.Name())

	inputs := AllocFilterInOut()
	defer inputs.Free()
	inputs.SetName("out")
	inputs.SetFilterContext(fcOut)
	inputs.SetPadIdx(0)
	inputs.SetNext(nil)

	var outputs *FilterInOut
	defer func() {
		if outputs != nil {
			outputs.Free()
		}
	}()
	var fcIns []*FilterContext
	for i := 0; i < 2; i++ {
		bufferSrc := FindFilterByName("buffer")
		require.NotNil(t, bufferSrc)

		fcIn, err := fg.NewFilterContext(bufferSrc, fmt.Sprintf("filter_in_%d", i+1), FilterArgs{
			"pix_fmt":      strconv.Itoa(int(PixelFormatYuv420P)),
			"pixel_aspect": "1/1",
			"time_base":    "1/1000",
			"video_size":   "1x1",
		})
		require.NoError(t, err)
		fcIns = append(fcIns, fcIn)
		defer fcIn.Free()

		o := AllocFilterInOut()
		o.SetName(fmt.Sprintf("input_%d", i+1))
		o.SetFilterContext(fcIn)
		o.SetPadIdx(0)
		o.SetNext(outputs)

		outputs = o
	}

	err = fg.Parse("[input_1]scale=2x2[scaled_1];[input_2]scale=3x3[scaled_2];[scaled_1][scaled_2]overlay", inputs, outputs)
	require.NoError(t, err)

	err = fg.Configure()
	require.NoError(t, err)

	require.Equal(t, 1, fcOut.NbInputs())
	require.Equal(t, 1, len(fcOut.Inputs()))
	require.Equal(t, NewRational(1, 1000), fcOut.Inputs()[0].TimeBase())
	require.Equal(t, 0, fcOut.NbOutputs())
	for _, fc := range fcIns {
		require.Equal(t, 0, fc.NbInputs())
		require.Equal(t, 1, fc.NbOutputs())
		require.Equal(t, 1, len(fc.Outputs()))
		require.Equal(t, NewRational(1, 1000), fc.Outputs()[0].TimeBase())
	}

	resp, err := fg.SendCommand("scale", "invalid", "a", NewFilterCommandFlags())
	require.Error(t, err)
	require.Empty(t, resp)
	resp, err = fg.SendCommand("scale", "width", "4", NewFilterCommandFlags().Add(FilterCommandFlagOne))
	require.NoError(t, err)
	require.Empty(t, resp)

	require.Equal(t, "                                                    +--------------+\nParsed_overlay_2:default--[2x2 1:1 yuv420p]--default|  filter_out  |\n                                                    | (buffersink) |\n                                                    +--------------+\n\n+-------------+\n| filter_in_1 |default--[1x1 1:1 yuv420p]--Parsed_scale_0:default\n|  (buffer)   |\n+-------------+\n\n+-------------+\n| filter_in_2 |default--[1x1 1:1 yuv420p]--Parsed_scale_1:default\n|  (buffer)   |\n+-------------+\n\n                                               +----------------+\nfilter_in_1:default--[1x1 1:1 yuv420p]--default| Parsed_scale_0 |default--[4x2 1:2 yuv420p]--Parsed_overlay_2:main\n                                               |    (scale)     |\n                                               +----------------+\n\n                                               +----------------+\nfilter_in_2:default--[1x1 1:1 yuv420p]--default| Parsed_scale_1 |default--[3x3 1:1 yuva420p]--Parsed_overlay_2:overlay\n                                               |    (scale)     |\n                                               +----------------+\n\n                                                   +------------------+\nParsed_scale_0:default--[4x2 1:2 yuv420p]------main| Parsed_overlay_2 |default--[2x2 1:1 yuv420p]--filter_out:default\nParsed_scale_1:default--[3x3 1:1 yuva420p]--overlay|    (overlay)     |\n                                                   +------------------+\n\n", fg.String())

	// TODO Test BuffersrcAddFrame
	// TODO Test BuffersinkGetFrame
}
