package astiav

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuffersrcFilterContext(t *testing.T) {
	fg := AllocFilterGraph()
	filter := FindFilterByName("movie")
	bufferSrcCtx, err := fg.NewBuffersrcFilterContext(filter, "movie")
	require.NoError(t, err)
	d := NewDictionary()
	require.NoError(t, d.Set("filename", "testdata/video.mp4", 0))
	require.NoError(t, bufferSrcCtx.Initialize(d))
}
