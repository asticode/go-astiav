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
	require.Nil(t, err)
	defer bsfc.Free()

	cl := bsfc.Class()
	require.NotNil(t, cl)
	require.Equal(t, "AVBSFContext", cl.Name())

	bsfc.SetTimeBaseIn(NewRational(15, 1))
	require.Equal(t, NewRational(15, 1), bsfc.TimeBaseIn())

	fc, err := globalHelper.inputFormatContext("video.mp4")
	require.NoError(t, err)
	ss := fc.Streams()
	require.Len(t, ss, 2)
	s1 := ss[0]

	cp1 := s1.CodecParameters()
	bsfc.SetCodecParametersIn(cp1)
	require.Equal(t, int64(441324), bsfc.CodecParametersIn().BitRate())

	// video.mp4 bit stream h264 format is avcc
	pkt1, err := globalHelper.inputFirstPacket("video.mp4")
	pkt1Bsf, errBsf := globalHelper.inputFirstPacketWithBitStreamFilter("video.mp4", "h264_mp4toannexb")
	require.NoError(t, err)
	require.NoError(t, errBsf)

	require.NotEqual(t, pkt1.Data(), pkt1Bsf.Data())
}
