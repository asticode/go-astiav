# v0.25.0

- `CodecParameters`.`CodecType` and `CodecParameters`.`SetCodecType` have been removed, use `CodecParameters`.`MediaType` and `CodecParameters`.`SetMediaType` instead
- `HardwareFrameContext`.`SetPixelFormat` has been replaced with `HardwareFrameContext`.`SetHardwarePixelFormat`
- `FormatContext`.`SetInterruptCallback` has been replaced with `FormatContext`.`SetIOInterrupter`

# v0.24.0

- use `FilterGraph`.`NewBuffersinkFilterContext` and `FilterGraph`.`NewBuffersrcFilterContext` instead of `FilterGraph`.`NewFilterContext` when creating `buffersink` and `buffersrc` filter contexts and use `BuffersinkFilterContext`.`GetFrame` and `BuffersrcFilterContext`.`AddFrame` to manipulate them. Use `BuffersinkFilterContext`.`FilterContext` and `BuffersrcFilterContext`.`FilterContext` in `FilterInOut`.`SetFilterContext`.
- `FilterLink` has been removed and methods like `BuffersinkFilterContext`.`ChannelLayout` have been added instead