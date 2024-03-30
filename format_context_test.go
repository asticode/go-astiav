package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatContext(t *testing.T) {
	fc1, err := globalHelper.inputFormatContext("video.mp4")
	require.NoError(t, err)
	ss := fc1.Streams()
	require.Len(t, ss, 2)
	s1 := ss[0]

	require.Equal(t, int64(607664), fc1.BitRate())
	require.Equal(t, NewFormatContextCtxFlags(0), fc1.CtxFlags())
	require.Equal(t, int64(5013333), fc1.Duration())
	require.True(t, fc1.EventFlags().Has(FormatEventFlagMetadataUpdated))
	require.True(t, fc1.Flags().Has(FormatContextFlagAutoBsf))
	require.Equal(t, NewRational(24, 1), fc1.GuessFrameRate(s1, nil))
	require.Equal(t, NewRational(1, 1), fc1.GuessSampleAspectRatio(s1, nil))
	require.True(t, fc1.InputFormat().Flags().Has(IOFormatFlagNoByteSeek))
	require.Equal(t, IOContextFlags(0), fc1.IOFlags())
	require.Equal(t, int64(0), fc1.MaxAnalyzeDuration())
	require.Equal(t, "isom", fc1.Metadata().Get("major_brand", nil, NewDictionaryFlags()).Value())
	require.Equal(t, int64(0), fc1.StartTime())
	require.Equal(t, 2, fc1.NbStreams())
	require.Len(t, fc1.Streams(), 2)
	cl := fc1.Class()
	require.NotNil(t, cl)
	require.Equal(t, "AVFormatContext", cl.Name())

	sdp, err := fc1.SDPCreate()
	require.NoError(t, err)
	require.Equal(t, "v=0\r\no=- 0 0 IN IP4 127.0.0.1\r\ns=Big Buck Bunny\r\nt=0 0\r\na=tool:libavformat 60.16.100\r\nm=video 0 RTP/AVP 96\r\nb=AS:441\r\na=rtpmap:96 H264/90000\r\na=fmtp:96 packetization-mode=1; sprop-parameter-sets=Z0LADasgKDPz4CIAAAMAAgAAAwBhHihUkA==,aM48gA==; profile-level-id=42C00D\r\na=control:streamid=0\r\nm=audio 0 RTP/AVP 97\r\nb=AS:161\r\na=rtpmap:97 MPEG4-GENERIC/48000/2\r\na=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3; config=1190\r\na=control:streamid=1\r\n", sdp)

	fc2, err := AllocOutputFormatContext(nil, "", "/tmp/test.mp4")
	require.NoError(t, err)
	defer fc2.Free()
	require.True(t, fc2.OutputFormat().Flags().Has(IOFormatFlagGlobalheader))

	fc3 := AllocFormatContext()
	require.NotNil(t, fc3)
	defer fc3.Free()
	c, err := OpenIOContext("testdata/video.mp4", NewIOContextFlags(IOContextFlagRead))
	require.NoError(t, err)
	defer c.Close() //nolint:errcheck
	fc3.SetPb(c)
	fc3.SetStrictStdCompliance(StrictStdComplianceExperimental)
	fc3.SetFlags(NewFormatContextFlags(FormatContextFlagAutoBsf))
	require.NotNil(t, fc3.Pb())
	require.Equal(t, StrictStdComplianceExperimental, fc3.StrictStdCompliance())
	require.True(t, fc3.Flags().Has(FormatContextFlagAutoBsf))
	s2 := fc3.NewStream(nil)
	require.NotNil(t, s2)
	s3 := fc3.NewStream(nil)
	require.NotNil(t, s3)
	require.Equal(t, 1, s3.Index())

	fc4 := AllocFormatContext()
	require.NotNil(t, fc4)
	defer fc4.Free()
	fc4.SetInterruptCallback().Interrupt()
	require.ErrorIs(t, fc4.OpenInput("testdata/video.mp4", nil, nil), ErrExit)

	// TODO Test ReadFrame
	// TODO Test SeekFrame
	// TODO Test Flush
	// TODO Test WriteHeader
	// TODO Test WriteFrame
	// TODO Test WriteInterleavedFrame
	// TODO Test WriteTrailer
}
