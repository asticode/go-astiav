package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChapter(t *testing.T) {
	fc, err := globalHelper.inputFormatContext("video_with_chapters.mp4", nil)
	require.NoError(t, err)
	cs := fc.Chapters()
	require.Len(t, cs, 2)

	c1 := cs[0]
	c2 := cs[1]

	require.Equal(t, int64(0), c1.ID())
	require.Equal(t, NewRational(1, 1000), c1.TimeBase())
	require.Equal(t, int64(0), c1.Start())
	require.Equal(t, int64(2501), c1.End())
	require.Equal(t, "Chapter 1 - First Half", c1.Metadata().Get("title", nil, NewDictionaryFlags()).Value())

	require.Equal(t, int64(1), c2.ID())
	require.Equal(t, NewRational(1, 1000), c2.TimeBase())
	require.Equal(t, int64(2501), c2.Start())
	require.Equal(t, int64(5000), c2.End())
	require.Equal(t, "Chapter 2 - Second Half", c2.Metadata().Get("title", nil, NewDictionaryFlags()).Value())

	c1.SetID(2)
	require.Equal(t, int64(2), c1.ID())
	c1.SetStart(1)
	require.Equal(t, int64(1), c1.Start())
	c1.SetEnd(1)
	require.Equal(t, int64(1), c1.End())

	d := NewDictionary()
	d.Set("k", "v", 0)
	c1.SetMetadata(d)
	e := c1.Metadata().Get("k", nil, 0)
	require.NotNil(t, e)
	require.Equal(t, "v", e.Value())
	c1.SetMetadata(nil)
	require.Nil(t, c1.Metadata())
}
