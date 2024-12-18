package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProgram(t *testing.T) {
	fc := AllocFormatContext()
	require.NotNil(t, fc)
	defer fc.Free()
	p := fc.NewProgram(1)
	require.Equal(t, 1, p.ID())
	p.SetID(2)
	require.Equal(t, 2, p.ID())
	p.SetDiscard(DiscardAll)
	require.Equal(t, DiscardAll, p.Discard())
	d := NewDictionary()
	require.NoError(t, d.Set("service_name", "test_service_name", 0))
	p.SetMetadata(d)
	require.Equal(t, p.Metadata().Get("service_name", nil, 0).Value(), "test_service_name")
	p.SetProgramNumber(101)
	require.Equal(t, 101, p.ProgramNumber())
	p.SetPmtPid(201)
	require.Equal(t, 201, p.PmtPid())
	p.SetPcrPid(301)
	require.Equal(t, 301, p.PcrPid())
	require.Equal(t, p.StartTime(), int64(-9223372036854775808))
	require.Equal(t, p.EndTime(), int64(-9223372036854775808))
	require.Equal(t, p.PtsWrapReference(), int64(-9223372036854775808))
	require.Equal(t, p.PtsWrapBehavior(), 0)
	s := fc.NewStream(nil)
	s.SetID(2)
	require.Equal(t, 0, p.NbStreams())
	require.Nil(t, p.StreamIndex(), nil)
	p.AddStream(s)
	require.Equal(t, 1, p.NbStreams())
	require.Equal(t, uint(0), *p.StreamIndex())
	var streamIndex uint = 1
	p.SetStreamIndex(streamIndex)
	require.Equal(t, streamIndex, *p.StreamIndex())
	s.SetIndex(int(streamIndex))
	ss := p.Streams()
	require.Equal(t, 1, len(ss))
	require.Equal(t, s.ID(), ss[0].ID())
}
