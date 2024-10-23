# v0.24.0

- use `FilterGraph`.`NewBuffersinkFilterContext` and `FilterGraph`.`NewBuffersrcFilterContext` instead of `FilterGraph`.`NewFilterContext` when creating `buffersink` and `buffersrc` filter contexts and use `BuffersinkFilterContext`.`GetFrame` and `BuffersrcFilterContext`.`AddFrame` to manipulate them. Use `BuffersinkFilterContext`.`FilterContext` and `BuffersrcFilterContext`.`FilterContext` in `FilterInOut`.`SetFilterContext`.
- `FilterLink` has been removed and methods like `BuffersinkFilterContext`.`ChannelLayout` have been added instead