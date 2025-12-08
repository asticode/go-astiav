package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	fc, err := globalHelper.inputFormatContext("video.mp4", nil)
	require.NoError(t, err)
	ss := fc.Streams()
	require.Len(t, ss, 2)
	s1 := ss[0]
	s2 := ss[1]

	require.Equal(t, 0, s1.Index())
	require.Equal(t, NewRational(24, 1), s1.AvgFrameRate())
	require.True(t, s1.DispositionFlags().Has(DispositionFlagDefault))
	require.Equal(t, int64(61440), s1.Duration())
	require.True(t, s1.EventFlags().Has(StreamEventFlag(2)))
	require.Equal(t, 1, s1.ID())
	require.Equal(t, "und", s1.Metadata().Get("language", nil, NewDictionaryFlags()).Value())
	require.Equal(t, int64(120), s1.NbFrames())
	require.Equal(t, NewRational(24, 1), s1.RFrameRate())
	require.Equal(t, NewRational(1, 1), s1.SampleAspectRatio())
	require.Equal(t, int64(0), s1.StartTime())
	require.Equal(t, NewRational(1, 12288), s1.TimeBase())
	cl := s1.Class()
	require.NotNil(t, cl)
	require.Equal(t, "AVStream", cl.Name())

	require.Equal(t, 1, s2.Index())
	require.Equal(t, int64(240640), s2.Duration())
	require.Equal(t, 2, s2.ID())
	require.Equal(t, int64(235), s2.NbFrames())
	require.Equal(t, int64(0), s2.StartTime())
	require.Equal(t, NewRational(1, 48000), s2.TimeBase())

	s1.SetAvgFrameRate(NewRational(2, 1))
	require.Equal(t, NewRational(2, 1), s1.AvgFrameRate())
	s1.SetDiscard(DiscardAll)
	require.Equal(t, DiscardAll, s1.Discard())
	s1.SetDispositionFlags(2)
	require.Equal(t, DispositionFlags(2), s1.DispositionFlags())
	s1.SetEventFlags(1)
	require.Equal(t, StreamEventFlags(1), s1.EventFlags())
	s1.SetID(2)
	require.Equal(t, 2, s1.ID())
	s1.SetIndex(1)
	require.Equal(t, 1, s1.Index())
	s1.SetPTSWrapBits(2)
	require.Equal(t, 2, s1.PTSWrapBits())
	s1.SetRFrameRate(NewRational(2, 1))
	require.Equal(t, NewRational(2, 1), s1.RFrameRate())
	s1.SetSampleAspectRatio(NewRational(2, 1))
	require.Equal(t, NewRational(2, 1), s1.SampleAspectRatio())
	s1.SetStartTime(1)
	require.Equal(t, int64(1), s1.StartTime())
	s1.SetTimeBase(NewRational(1, 1))
	require.Equal(t, NewRational(1, 1), s1.TimeBase())

	d := NewDictionary()
	d.Set("k", "v", 0)
	s1.SetMetadata(d)
	e := s1.Metadata().Get("k", nil, 0)
	require.NotNil(t, e)
	require.Equal(t, "v", e.Value())
	s1.SetMetadata(nil)
	require.Nil(t, s1.Metadata())
}
