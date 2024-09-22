[![GoReportCard](http://goreportcard.com/badge/github.com/asticode/go-astiav)](http://goreportcard.com/report/github.com/asticode/go-astiav)
[![GoDoc](https://godoc.org/github.com/asticode/go-astiav?status.svg)](https://godoc.org/github.com/asticode/go-astiav)
[![Test](https://github.com/asticode/go-astiav/actions/workflows/test.yml/badge.svg)](https://github.com/asticode/go-astiav/actions/workflows/test.yml)
[![Coveralls](https://coveralls.io/repos/github/asticode/go-astiav/badge.svg?branch=master)](https://coveralls.io/github/asticode/go-astiav)

`astiav` is a Golang library providing C bindings for [ffmpeg](https://github.com/FFmpeg/FFmpeg)

It's only compatible with `ffmpeg` `n7.0`.

Its main goals are to:

- [x] provide a better GO idiomatic API
    - standard error pattern
    - typed constants and flags
    - struct-based functions
    - ...
- [x] provide the GO version of [ffmpeg examples](https://github.com/FFmpeg/FFmpeg/tree/n7.0/doc/examples)
- [x] be fully tested

# Examples

Examples are located in the [examples](examples) directory and mirror as much as possible the [ffmpeg examples](https://github.com/FFmpeg/FFmpeg/tree/n7.0/doc/examples).

|name|astiav|ffmpeg|
|---|---|---|
|BitStream Filtering|[see](examples/bit_stream_filtering/main.go)|X
|Custom IO Demuxing|[see](examples/custom_io_demuxing/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/avio_reading.c)
|Demuxing/Decoding|[see](examples/demuxing_decoding/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/demuxing_decoding.c)
|Filtering|[see](examples/filtering/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/filtering_video.c)
|Hardware Decoding|[see](examples/hardware_decoding/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/hw_decode.c)
|Remuxing|[see](examples/remuxing/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/remuxing.c)
|Scaling|[see](examples/scaling/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/scaling_video.c)
|Transcoding|[see](examples/transcoding/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/transcoding.c)

*Tip: you can use the video sample located in the `testdata` directory for your tests*

# Install ffmpeg from source

If you don't know how to install `ffmpeg`, you can use the following to install it from source:

```sh
$ make install-ffmpeg
```

`ffmpeg` will be built from source in a directory named `tmp` and located in you working directory

For your GO code to pick up `ffmpeg` dependency automatically, you'll need to add the following environment variables:

(don't forget to replace `{{ path to your working directory }}` with the absolute path to your working directory)

```sh
export CGO_LDFLAGS="-L{{ path to your working directory }}/tmp/n7.0/lib/",
export CGO_CFLAGS="-I{{ path to your working directory }}/tmp/n7.0/include/",
export PKG_CONFIG_PATH="{{ path to your working directory }}/tmp/n7.0/lib/pkgconfig",
```

# Why astiav?

After maintaining for several years the most starred [fork](https://github.com/asticode/goav) of [goav](https://github.com/giorgisio/goav), I've decided to write from scratch my own C bindings to fix most of the problems I still encountered using `goav`.