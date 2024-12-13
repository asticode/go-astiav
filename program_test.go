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
	s := fc.NewStream(nil)
	s.SetID(2)
	require.Equal(t, 0, p.NbStreams())
	p.AddStream(s)
	require.Equal(t, 1, p.NbStreams())
	ss := p.Streams()
	require.Equal(t, 1, len(ss))
	require.Equal(t, s.ID(), ss[0].ID())

	p.SetFlags(1)
	require.Equal(t, 1, p.Flags())
	p.SetDiscard(DiscardAll)
	require.Equal(t, DiscardAll, p.Discard())
	d := NewDictionary()
	_ = d.Set("service_name", "test_service_name", 0)
	p.SetMetadata(d)
	require.Equal(t, p.Metadata().Get("service_name", nil, 0).Value(), "test_service_name")
}
