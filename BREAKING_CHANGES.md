# v0.31.0

- removed `AllocImage` and `AllocSamples` that are considered useless using CGO until proven otherwise

# v0.30.0

- `HardwareFrameContext` has been renamed to `HardwareFramesContext`

# v0.29.0

- `NewFilterContext` has been removed, use `NewBuffersinkFilterContext` or `NewBuffersrcFilterContext` instead 
- `args` has been removed from `NewBuffersinkFilterContext` and `NewBuffersrcFilterContext`. Instead, after calling `NewBuffersrcFilterContext`, you need to use `BuffersrcFilterContext`.`SetParameters` then `BuffersrcFilterContext`.`Initialize`. You don't need to use anything else after calling `NewBuffersinkFilterContext`

# v0.27.0

- make sure to call the `IOInterrupter`.`Free` method after using `NewIOInterrupter`

# v0.25.0

- `CodecParameters`.`CodecType` and `CodecParameters`.`SetCodecType` have been removed, use `CodecParameters`.`MediaType` and `CodecParameters`.`SetMediaType` instead
- `HardwareFrameContext`.`SetPixelFormat` has been replaced with `HardwareFrameContext`.`SetHardwarePixelFormat`
- `FormatContext`.`SetInterruptCallback` has been replaced with `FormatContext`.`SetIOInterrupter`

# v0.24.0

- use `FilterGraph`.`NewBuffersinkFilterContext` and `FilterGraph`.`NewBuffersrcFilterContext` instead of `FilterGraph`.`NewFilterContext` when creating `buffersink` and `buffersrc` filter contexts and use `BuffersinkFilterContext`.`GetFrame` and `BuffersrcFilterContext`.`AddFrame` to manipulate them. Use `BuffersinkFilterContext`.`FilterContext` and `BuffersrcFilterContext`.`FilterContext` in `FilterInOut`.`SetFilterContext`.
- `FilterLink` has been removed and methods like `BuffersinkFilterContext`.`ChannelLayout` have been added instead