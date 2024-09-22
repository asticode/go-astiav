#include <libavutil/channel_layout.h>

// Calling C.AV_CHANNEL_LAYOUT_* in Go gives a "could not determine kind of name for X" error
// therefore we need to bridge the channel layout values
AVChannelLayout *astiavChannelLayoutMono              = &(AVChannelLayout)AV_CHANNEL_LAYOUT_MONO;
AVChannelLayout *astiavChannelLayoutStereo            = &(AVChannelLayout)AV_CHANNEL_LAYOUT_STEREO;
AVChannelLayout *astiavChannelLayout2Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_2POINT1;
AVChannelLayout *astiavChannelLayout21                = &(AVChannelLayout)AV_CHANNEL_LAYOUT_2_1;
AVChannelLayout *astiavChannelLayoutSurround          = &(AVChannelLayout)AV_CHANNEL_LAYOUT_SURROUND;
AVChannelLayout *astiavChannelLayout3Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_3POINT1;
AVChannelLayout *astiavChannelLayout4Point0           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_4POINT0;
AVChannelLayout *astiavChannelLayout4Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_4POINT1;
AVChannelLayout *astiavChannelLayout22                = &(AVChannelLayout)AV_CHANNEL_LAYOUT_2_2;
AVChannelLayout *astiavChannelLayoutQuad              = &(AVChannelLayout)AV_CHANNEL_LAYOUT_QUAD;
AVChannelLayout *astiavChannelLayout5Point0           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT0;
AVChannelLayout *astiavChannelLayout5Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT1;
AVChannelLayout *astiavChannelLayout5Point0Back       = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT0_BACK;
AVChannelLayout *astiavChannelLayout5Point1Back       = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT1_BACK;
AVChannelLayout *astiavChannelLayout6Point0           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_6POINT0;
AVChannelLayout *astiavChannelLayout6Point0Front      = &(AVChannelLayout)AV_CHANNEL_LAYOUT_6POINT0_FRONT;
AVChannelLayout *astiavChannelLayoutHexagonal         = &(AVChannelLayout)AV_CHANNEL_LAYOUT_HEXAGONAL;
AVChannelLayout *astiavChannelLayout3Point1Point2     = &(AVChannelLayout)AV_CHANNEL_LAYOUT_3POINT1POINT2;
AVChannelLayout *astiavChannelLayout6Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_6POINT1;
AVChannelLayout *astiavChannelLayout6Point1Back       = &(AVChannelLayout)AV_CHANNEL_LAYOUT_6POINT1_BACK;
AVChannelLayout *astiavChannelLayout6Point1Front      = &(AVChannelLayout)AV_CHANNEL_LAYOUT_6POINT1_FRONT;
AVChannelLayout *astiavChannelLayout7Point0           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT0;
AVChannelLayout *astiavChannelLayout7Point0Front      = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT0_FRONT;
AVChannelLayout *astiavChannelLayout7Point1           = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1;
AVChannelLayout *astiavChannelLayout7Point1Wide       = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1_WIDE;
AVChannelLayout *astiavChannelLayout7Point1WideBack   = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1_WIDE_BACK;
AVChannelLayout *astiavChannelLayout5Point1Point2Back = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT1POINT2_BACK;
AVChannelLayout *astiavChannelLayoutOctagonal         = &(AVChannelLayout)AV_CHANNEL_LAYOUT_OCTAGONAL;
AVChannelLayout *astiavChannelLayoutCube              = &(AVChannelLayout)AV_CHANNEL_LAYOUT_CUBE;
AVChannelLayout *astiavChannelLayout5Point1Point4Back = &(AVChannelLayout)AV_CHANNEL_LAYOUT_5POINT1POINT4_BACK;
AVChannelLayout *astiavChannelLayout7Point1Point2     = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1POINT2;
AVChannelLayout *astiavChannelLayout7Point1Point4Back = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1POINT4_BACK;
AVChannelLayout *astiavChannelLayoutHexadecagonal     = &(AVChannelLayout)AV_CHANNEL_LAYOUT_HEXADECAGONAL;
AVChannelLayout *astiavChannelLayoutStereoDownmix     = &(AVChannelLayout)AV_CHANNEL_LAYOUT_STEREO_DOWNMIX;
AVChannelLayout *astiavChannelLayout22Point2          = &(AVChannelLayout)AV_CHANNEL_LAYOUT_22POINT2;
AVChannelLayout *astiavChannelLayout7Point1TopBack    = &(AVChannelLayout)AV_CHANNEL_LAYOUT_7POINT1_TOP_BACK;