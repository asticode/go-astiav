// TODO Fix https://github.com/asticode/go-astiav/actions/runs/5853322732/job/15867145888
//go:build !windows

package astiav

import (
	"fmt"
	"strconv"
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
	type link struct {
		channelLayout     ChannelLayout
		frameRate         Rational
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
		buffersinkExpectedInput link
		buffersinkName          string
		buffersrcName           string
		commands                []command
		content                 string
		s                       string
		sources                 []FilterArgs
	}
	for _, v := range []graph{
		{
			buffersinkExpectedInput: link{
				frameRate:         NewRational(4, 1),
				height:            4,
				mediaType:         MediaTypeVideo,
				pixelFormat:       PixelFormatRgba,
				sampleAspectRatio: NewRational(1, 4),
				timeBase:          NewRational(1, 4),
				width:             2,
			},
			buffersinkName: "buffersink",
			buffersrcName:  "buffer",
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
			content: "[input_1]scale=2x4,settb=1/4,fps=fps=4/1,format=pix_fmts=rgba,setsar=1/4",
			s:       "                                                +--------------+\nParsed_setsar_4:default--[2x4 1:4 rgba]--default|  filter_out  |\n                                                | (buffersink) |\n                                                +--------------+\n\n+-------------+\n| filter_in_1 |default--[1x2 1:2 yuv420p]--Parsed_scale_0:default\n|  (buffer)   |\n+-------------+\n\n                                               +----------------+\nfilter_in_1:default--[1x2 1:2 yuv420p]--default| Parsed_scale_0 |default--[2x4 1:2 rgba]--Parsed_settb_1:default\n                                               |    (scale)     |\n                                               +----------------+\n\n                                               +----------------+\nParsed_scale_0:default--[2x4 1:2 rgba]--default| Parsed_settb_1 |default--[2x4 1:2 rgba]--Parsed_fps_2:default\n                                               |    (settb)     |\n                                               +----------------+\n\n                                               +--------------+\nParsed_settb_1:default--[2x4 1:2 rgba]--default| Parsed_fps_2 |default--[2x4 1:2 rgba]--Parsed_format_3:default\n                                               |    (fps)     |\n                                               +--------------+\n\n                                             +-----------------+\nParsed_fps_2:default--[2x4 1:2 rgba]--default| Parsed_format_3 |default--[2x4 1:2 rgba]--Parsed_setsar_4:default\n                                             |    (format)     |\n                                             +-----------------+\n\n                                                +-----------------+\nParsed_format_3:default--[2x4 1:2 rgba]--default| Parsed_setsar_4 |default--[2x4 1:4 rgba]--filter_out:default\n                                                |    (setsar)     |\n                                                +-----------------+\n\n",
			sources: []FilterArgs{
				{
					"height":    "2",
					"pix_fmt":   strconv.Itoa(int(PixelFormatYuv420P)),
					"sar":       "1/2",
					"time_base": "1/2",
					"width":     "1",
				},
			},
		},
		{
			buffersinkExpectedInput: link{
				channelLayout: ChannelLayoutStereo,
				mediaType:     MediaTypeAudio,
				sampleFormat:  SampleFormatS16,
				sampleRate:    3,
				timeBase:      NewRational(1, 4),
			},
			buffersinkName: "abuffersink",
			buffersrcName:  "abuffer",
			content:        "[input_1]aformat=sample_fmts=s16:channel_layouts=stereo:sample_rates=3,asettb=1/4",
			s:              "                                                  +---------------+\nParsed_asettb_1:default--[3Hz s16:stereo]--default|  filter_out   |\n                                                  | (abuffersink) |\n                                                  +---------------+\n\n+-------------+\n| filter_in_1 |default--[2Hz fltp:mono]--auto_aresample_0:default\n|  (abuffer)  |\n+-------------+\n\n                                                   +------------------+\nauto_aresample_0:default--[3Hz s16:stereo]--default| Parsed_aformat_0 |default--[3Hz s16:stereo]--Parsed_asettb_1:default\n                                                   |    (aformat)     |\n                                                   +------------------+\n\n                                                   +-----------------+\nParsed_aformat_0:default--[3Hz s16:stereo]--default| Parsed_asettb_1 |default--[3Hz s16:stereo]--filter_out:default\n                                                   |    (asettb)     |\n                                                   +-----------------+\n\n                                             +------------------+\nfilter_in_1:default--[2Hz fltp:mono]--default| auto_aresample_0 |default--[3Hz s16:stereo]--Parsed_aformat_0:default\n                                             |   (aresample)    |\n                                             +------------------+\n\n",
			sources: []FilterArgs{
				{
					"channel_layout": ChannelLayoutMono.String(),
					"sample_fmt":     strconv.Itoa(int(SampleFormatFltp)),
					"sample_rate":    "2",
					"time_base":      "1/2",
				},
			},
		},
	} {
		fg := AllocFilterGraph()
		require.NotNil(t, fg)
		defer fg.Free()

		buffersrc := FindFilterByName(v.buffersrcName)
		require.NotNil(t, buffersrc)
		buffersink := FindFilterByName(v.buffersinkName)
		require.NotNil(t, buffersink)

		buffersinkContext, err := fg.NewFilterContext(buffersink, "filter_out", nil)
		require.NoError(t, err)
		cl = buffersinkContext.Class()
		require.NotNil(t, cl)
		require.Equal(t, "AVFilter", cl.Name())

		inputs := AllocFilterInOut()
		defer inputs.Free()
		inputs.SetName("out")
		inputs.SetFilterContext(buffersinkContext)
		inputs.SetPadIdx(0)
		inputs.SetNext(nil)

		var outputs *FilterInOut
		defer func() {
			if outputs != nil {
				outputs.Free()
			}
		}()

		var buffersrcContexts []*FilterContext
		for idx, src := range v.sources {
			buffersrcContext, err := fg.NewFilterContext(buffersrc, fmt.Sprintf("filter_in_%d", idx+1), src)
			require.NoError(t, err)
			buffersrcContexts = append(buffersrcContexts, buffersrcContext)

			o := AllocFilterInOut()
			o.SetName(fmt.Sprintf("input_%d", idx+1))
			o.SetFilterContext(buffersrcContext)
			o.SetPadIdx(0)
			o.SetNext(outputs)

			outputs = o
		}

		require.NoError(t, fg.Parse(v.content, inputs, outputs))
		require.NoError(t, fg.Configure())

		require.Equal(t, 1, buffersinkContext.NbInputs())
		links := buffersinkContext.Inputs()
		require.Equal(t, 1, len(links))
		e, g := v.buffersinkExpectedInput, links[0]
		require.Equal(t, e.frameRate, g.FrameRate())
		require.Equal(t, e.mediaType, g.MediaType())
		require.Equal(t, e.timeBase, g.TimeBase())
		switch e.mediaType {
		case MediaTypeAudio:
			require.True(t, e.channelLayout.Equal(g.ChannelLayout()))
			require.Equal(t, e.sampleFormat, g.SampleFormat())
			require.Equal(t, e.sampleRate, g.SampleRate())
		default:
			require.Equal(t, e.height, g.Height())
			require.Equal(t, e.pixelFormat, g.PixelFormat())
			require.Equal(t, e.sampleAspectRatio, g.SampleAspectRatio())
			require.Equal(t, e.width, g.Width())
		}

		for _, buffersrcContext := range buffersrcContexts {
			require.Equal(t, 0, buffersrcContext.NbInputs())
			require.Equal(t, 1, buffersrcContext.NbOutputs())
			links := buffersrcContext.Outputs()
			require.Equal(t, 1, len(links))
			require.Equal(t, v.buffersinkExpectedInput.mediaType, links[0].MediaType())
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
