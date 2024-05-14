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

	bsfc.SetTimeBaseIn(NewRational(15, 1))
	require.Equal(t, NewRational(15, 1), bsfc.TimeBaseIn())

	cp1 := AllocCodecParameters()
	require.NotNil(t, cp1)
	defer cp1.Free()
	cp1.SetCodecID(CodecIDH264)

	bsfc.SetCodecParametersIn(cp1)
	require.Equal(t, CodecIDH264, bsfc.CodecParametersIn().CodecID())

	// TODO: add tests for send and receive packet flows
}
