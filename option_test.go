package astiav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOption(t *testing.T) {
	fc, err := AllocOutputFormatContext(nil, "mp4", "")
	require.NoError(t, err)
	pd := fc.PrivateData()
	require.NotNil(t, pd)
	os := pd.Options()
	require.NotNil(t, os)
	require.Equal(t, []Option{{Name: "brand"}, {Name: "empty_hdlr_name"}, {Name: "encryption_key"}, {Name: "encryption_kid"}, {Name: "encryption_scheme"}, {Name: "frag_duration"}, {Name: "frag_interleave"}, {Name: "frag_size"}, {Name: "fragment_index"}, {Name: "iods_audio_profile"}, {Name: "iods_video_profile"}, {Name: "ism_lookahead"}, {Name: "movflags"}, {Name: "cmaf"}, {Name: "dash"}, {Name: "default_base_moof"}, {Name: "delay_moov"}, {Name: "disable_chpl"}, {Name: "empty_moov"}, {Name: "faststart"}, {Name: "frag_custom"}, {Name: "frag_discont"}, {Name: "frag_every_frame"}, {Name: "frag_keyframe"}, {Name: "global_sidx"}, {Name: "isml"}, {Name: "moov_size"}, {Name: "negative_cts_offsets"}, {Name: "omit_tfhd_offset"}, {Name: "prefer_icc"}, {Name: "rtphint"}, {Name: "separate_moof"}, {Name: "skip_sidx"}, {Name: "skip_trailer"}, {Name: "use_metadata_tags"}, {Name: "write_colr"}, {Name: "write_gama"}, {Name: "min_frag_duration"}, {Name: "mov_gamma"}, {Name: "movie_timescale"}, {Name: "rtpflags"}, {Name: "latm"}, {Name: "rfc2190"}, {Name: "skip_rtcp"}, {Name: "h264_mode0"}, {Name: "send_bye"}, {Name: "skip_iods"}, {Name: "use_editlist"}, {Name: "use_stream_ids_as_track_ids"}, {Name: "video_track_timescale"}, {Name: "write_btrt"}, {Name: "write_prft"}, {Name: "pts"}, {Name: "wallclock"}, {Name: "write_tmcd"}}, os.List())
	_, err = os.Get("invalid", NewOptionSearchFlags())
	require.Error(t, err)
	v, err := os.Get("brand", NewOptionSearchFlags())
	require.NoError(t, err)
	require.Equal(t, "", v)
	require.Error(t, os.Set("invalid", "", NewOptionSearchFlags()))
	const brand = "test"
	require.NoError(t, os.Set("brand", brand, NewOptionSearchFlags()))
	v, err = os.Get("brand", NewOptionSearchFlags())
	require.NoError(t, err)
	require.Equal(t, brand, v)
}
