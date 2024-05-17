package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitStreamFilterContext(t *testing.T) {
	bsf := FindBitStreamFilterByName("null")
	require.NotNil(t, bsf)

	bsfc, err := AllocBitStreamFilterContext(bsf)
	require.NotNil(t, bsfc)
	require.NoError(t, err)
	defer bsfc.Free()

	cl := bsfc.Class()
	require.NotNil(t, cl)
	require.Equal(t, "AVBSFContext", cl.Name())

	bsfc.SetInputTimeBase(NewRational(15, 1))
	require.Equal(t, NewRational(15, 1), bsfc.InputTimeBase())

	cp1 := AllocCodecParameters()
	require.NotNil(t, cp1)
	defer cp1.Free()
	cp1.SetCodecID(CodecIDH264)
	require.NoError(t, cp1.Copy(bsfc.InputCodecParameters()))
	require.Equal(t, CodecIDH264, bsfc.InputCodecParameters().CodecID())

	require.NoError(t, bsfc.Initialize())

	// TODO Test SendPacket
	// TODO Test ReceivePacket
}
