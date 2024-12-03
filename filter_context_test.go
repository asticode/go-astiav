package astiav

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuffersrcParameters(t *testing.T) {
	buffersrcParameters := AllocBuffersrcParameters()
	require.NotNil(t, buffersrcParameters)

	hardwareDeviceCtx, err := CreateHardwareDeviceContext(HardwareDeviceTypeCUDA, "0", nil, 0)
	require.NoError(t, err)
	hardwareFrameCtx := AllocHardwareFrameContext(hardwareDeviceCtx)
	require.NotNil(t, hardwareFrameCtx)
	buffersrcParameters.SetHardwareFrameContext(hardwareFrameCtx)
	require.Equal(t, hardwareFrameCtx, buffersrcParameters.HardwareFrameContext())

	args := FilterArgs{
		"pix_fmt":      "0",
		"pixel_aspect": "1/1",
		"time_base":    "1/1000",
		"video_size":   "1920x1080",
	}
	buffersrc := FindFilterByName("buffer")
	fg := AllocFilterGraph()
	buffersrcCtx, err := fg.NewBuffersrcFilterContext(buffersrc, "in", args)
	require.NoError(t, err)
	err = buffersrcCtx.SetBuffersrcParameters(buffersrcParameters)
	require.NoError(t, err)
	require.Equal(t, buffersrcParameters.HardwareFrameContext(), buffersrcCtx.BuffersrcParameters().HardwareFrameContext())
}
