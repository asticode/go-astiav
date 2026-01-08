package astiav

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterGraph(t *testing.T) {
	fg1 := AllocFilterGraph()
	require.NotNil(t, fg1)
	defer fg1.Free()
	cl := fg1.Class()
	require.NotNil(t, cl)
	require.Equal(t, "AVFilterGraph", cl.Name())
	fg1.SetThreadCount(2)
	require.Equal(t, 2, fg1.ThreadCount())
	fg1.SetThreadType(ThreadTypeSlice)
	require.Equal(t, ThreadTypeSlice, fg1.ThreadType())

	type command struct {
		args      string
		cmd       string
		flags     FilterCommandFlags
		resp      string
		target    string
		withError bool
	}
	type buffersink struct {
		channelLayout     ChannelLayout
		colorRange        ColorRange
		colorSpace        ColorSpace
		frameRate         Rational
		height            int
		mediaType         MediaType
		name              string
		pixelFormat       PixelFormat
		sampleAspectRatio Rational
		sampleFormat      SampleFormat
		sampleRate        int
		timeBase          Rational
		width             int
	}
	type buffersrc struct {
		name string
	}
	type buffersrcParameters struct {
		channelLayout     ChannelLayout
		height            int
		mediaType         MediaType
		pixelFormat       PixelFormat
		sampleAspectRatio Rational
		sampleFormat      SampleFormat
		sampleRate        int
		timeBase          Rational
		width             int
	}
	type graph struct {
		buffersink buffersink
		buffersrc  buffersrc
		commands   []command
		content    string
		s          string
		sources    []buffersrcParameters
	}
	for _, v := range []graph{
		{
			buffersink: buffersink{
				colorRange:        ColorRangeUnspecified,
				colorSpace:        ColorSpaceUnspecified,
				frameRate:         NewRational(4, 1),
				height:            8,
				mediaType:         MediaTypeVideo,
				name:              "buffersink",
				pixelFormat:       PixelFormatYuv420P,
				sampleAspectRatio: NewRational(2, 1),
				timeBase:          NewRational(1, 4),
				width:             4,
			},
			buffersrc: buffersrc{name: "buffer"},
			commands: []command{
				{
					args:      "a",
					cmd:       "invalid",
					flags:     NewFilterCommandFlags(),
					target:    "scale",
					withError: true,
				},
				{
					args:   "4",
					cmd:    "width",
					flags:  NewFilterCommandFlags().Add(FilterCommandFlagOne),
					target: "scale",
				},
			},
			content: "[input_1]scale=4x8,settb=1/4,fps=fps=4/1,format=pix_fmts=yuv420p,setsar=2/1",
			s:       "                                                   +--------------+\nParsed_setsar_4:default--[4x8 2:1 yuv420p]--default|  filter_out  |\n                                                   | (buffersink) |\n                                                   +--------------+\n\n+-------------+\n| filter_in_1 |default--[2x4 1:2 rgba]--Parsed_scale_0:default\n|  (buffer)   |\n+-------------+\n\n                                            +----------------+\nfilter_in_1:default--[2x4 1:2 rgba]--default| Parsed_scale_0 |default--[4x8 1:2 yuv420p]--Parsed_settb_1:default\n                                            |    (scale)     |\n                                            +----------------+\n\n                                                  +----------------+\nParsed_scale_0:default--[4x8 1:2 yuv420p]--default| Parsed_settb_1 |default--[4x8 1:2 yuv420p]--Parsed_fps_2:default\n                                                  |    (settb)     |\n                                                  +----------------+\n\n                                                  +--------------+\nParsed_settb_1:default--[4x8 1:2 yuv420p]--default| Parsed_fps_2 |default--[4x8 1:2 yuv420p]--Parsed_format_3:default\n                                                  |    (fps)     |\n                                                  +--------------+\n\n                                                +-----------------+\nParsed_fps_2:default--[4x8 1:2 yuv420p]--default| Parsed_format_3 |default--[4x8 1:2 yuv420p]--Parsed_setsar_4:default\n                                                |    (format)     |\n                                                +-----------------+\n\n                                                   +-----------------+\nParsed_format_3:default--[4x8 1:2 yuv420p]--default| Parsed_setsar_4 |default--[4x8 2:1 yuv420p]--filter_out:default\n                                                   |    (setsar)     |\n                                                   +-----------------+\n\n",
			sources: []buffersrcParameters{
				{
					height:            4,
					mediaType:         MediaTypeVideo,
					pixelFormat:       PixelFormatRgba,
					sampleAspectRatio: NewRational(1, 2),
					timeBase:          NewRational(1, 2),
					width:             2,
				},
			},
		},
		{
			buffersink: buffersink{
				channelLayout: ChannelLayoutStereo,
				mediaType:     MediaTypeAudio,
				name:          "abuffersink",
				sampleFormat:  SampleFormatS16,
				sampleRate:    3,
				timeBase:      NewRational(1, 4),
			},
			buffersrc: buffersrc{name: "abuffer"},
			content:   "[input_1]aformat=sample_fmts=s16:channel_layouts=stereo:sample_rates=3,asettb=1/4",
			s:         "                                                  +---------------+\nParsed_asettb_1:default--[3Hz s16:stereo]--default|  filter_out   |\n                                                  | (abuffersink) |\n                                                  +---------------+\n\n+-------------+\n| filter_in_1 |default--[2Hz fltp:mono]--auto_aresample_0:default\n|  (abuffer)  |\n+-------------+\n\n                                                   +------------------+\nauto_aresample_0:default--[3Hz s16:stereo]--default| Parsed_aformat_0 |default--[3Hz s16:stereo]--Parsed_asettb_1:default\n                                                   |    (aformat)     |\n                                                   +------------------+\n\n                                                   +-----------------+\nParsed_aformat_0:default--[3Hz s16:stereo]--default| Parsed_asettb_1 |default--[3Hz s16:stereo]--filter_out:default\n                                                   |    (asettb)     |\n                                                   +-----------------+\n\n                                             +------------------+\nfilter_in_1:default--[2Hz fltp:mono]--default| auto_aresample_0 |default--[3Hz s16:stereo]--Parsed_aformat_0:default\n                                             |   (aresample)    |\n                                             +------------------+\n\n",
			sources: []buffersrcParameters{
				{
					channelLayout: ChannelLayoutMono,
					mediaType:     MediaTypeAudio,
					sampleFormat:  SampleFormatFltp,
					sampleRate:    2,
					timeBase:      NewRational(1, 2),
				},
			},
		},
	} {
		fg := AllocFilterGraph()
		require.NotNil(t, fg)
		defer fg.Free()

		buffersrc := FindFilterByName(v.buffersrc.name)
		require.NotNil(t, buffersrc)
		buffersink := FindFilterByName(v.buffersink.name)
		require.NotNil(t, buffersink)

		buffersinkContext, err := fg.NewBuffersinkFilterContext(buffersink, "filter_out")
		require.NoError(t, err)
		require.Equal(t, buffersink, buffersinkContext.FilterContext().Filter())
		cl = buffersinkContext.FilterContext().Class()
		require.NotNil(t, cl)
		require.Equal(t, "AVFilter", cl.Name())

		inputs := AllocFilterInOut()
		defer inputs.Free()
		inputs.SetName("out")
		inputs.SetFilterContext(buffersinkContext.FilterContext())
		inputs.SetPadIdx(0)
		inputs.SetNext(nil)

		var outputs *FilterInOut
		defer func() {
			if outputs != nil {
				outputs.Free()
			}
		}()

		var buffersrcContexts []*BuffersrcFilterContext
		for idx, src := range v.sources {
			buffersrcContext, err := fg.NewBuffersrcFilterContext(buffersrc, fmt.Sprintf("filter_in_%d", idx+1))
			require.NoError(t, err)
			buffersrcContextParameters := AllocBuffersrcFilterContextParameters()
			defer buffersrcContextParameters.Free()
			switch src.mediaType {
			case MediaTypeAudio:
				buffersrcContextParameters.SetChannelLayout(src.channelLayout)
				buffersrcContextParameters.SetSampleFormat(src.sampleFormat)
				buffersrcContextParameters.SetSampleRate(src.sampleRate)
				buffersrcContextParameters.SetTimeBase(src.timeBase)
			default:
				buffersrcContextParameters.SetHeight(src.height)
				buffersrcContextParameters.SetPixelFormat(src.pixelFormat)
				buffersrcContextParameters.SetSampleAspectRatio(src.sampleAspectRatio)
				buffersrcContextParameters.SetTimeBase(src.timeBase)
				buffersrcContextParameters.SetWidth(src.width)
			}
			buffersrcContext.SetParameters(buffersrcContextParameters)
			require.NoError(t, buffersrcContext.Initialize(nil))
			buffersrcContexts = append(buffersrcContexts, buffersrcContext)

			o := AllocFilterInOut()
			o.SetName(fmt.Sprintf("input_%d", idx+1))
			o.SetFilterContext(buffersrcContext.FilterContext())
			o.SetPadIdx(0)
			o.SetNext(outputs)

			outputs = o
		}

		require.Equal(t, len(buffersrcContexts)+1, fg.NbFilters())
		fs := fg.Filters()
		require.Equal(t, len(buffersrcContexts)+1, len(fs))
		require.Equal(t, buffersinkContext.FilterContext(), fs[0])
		for idx, c := range fs[1:] {
			require.Equal(t, buffersrcContexts[idx].FilterContext(), c)
		}

		require.NoError(t, fg.Parse(v.content, inputs, outputs))
		require.NoError(t, fg.Configure())

		require.Equal(t, v.buffersink.frameRate, buffersinkContext.FrameRate())
		require.Equal(t, v.buffersink.mediaType, buffersinkContext.MediaType())
		require.Equal(t, v.buffersink.timeBase, buffersinkContext.TimeBase())
		switch v.buffersink.mediaType {
		case MediaTypeAudio:
			require.True(t, v.buffersink.channelLayout.Equal(buffersinkContext.ChannelLayout()))
			require.Equal(t, v.buffersink.sampleFormat, buffersinkContext.SampleFormat())
			require.Equal(t, v.buffersink.sampleRate, buffersinkContext.SampleRate())
		default:
			require.Equal(t, v.buffersink.colorRange, buffersinkContext.ColorRange())
			require.Equal(t, v.buffersink.colorSpace, buffersinkContext.ColorSpace())
			require.Equal(t, v.buffersink.height, buffersinkContext.Height())
			require.Equal(t, v.buffersink.pixelFormat, buffersinkContext.PixelFormat())
			require.Equal(t, v.buffersink.sampleAspectRatio, buffersinkContext.SampleAspectRatio())
			require.Equal(t, v.buffersink.width, buffersinkContext.Width())
		}

		require.Equal(t, v.s, fg.String())

		for _, command := range v.commands {
			resp, err := fg.SendCommand(command.target, command.cmd, command.args, command.flags)
			if command.withError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, command.resp, resp)
		}
	}

	fg2 := AllocFilterGraph()
	require.NotNil(t, fg2)
	defer fg2.Free()
	fgs, err := fg2.ParseSegment("anullsrc")
	require.NoError(t, err)
	defer fgs.Free()
	require.Equal(t, 1, fgs.NbChains())
	cs := fgs.Chains()
	require.Equal(t, 1, len(cs))
	require.Equal(t, 1, cs[0].NbFilters())
	fs := cs[0].Filters()
	require.Equal(t, 1, len(fs))
	f := FindFilterByName(fs[0].FilterName())
	require.NotNil(t, f)
	require.Equal(t, 0, f.NbInputs())
	require.Equal(t, 1, f.NbOutputs())
	os := f.Outputs()
	require.Equal(t, 1, len(os))
	require.Equal(t, MediaTypeAudio, os[0].MediaType())

	// TODO Test BuffersrcAddFrame
	// TODO Test BuffersinkGetFrame
}
