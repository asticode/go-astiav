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
|Custom IO Muxing|[see](examples/custom_io_muxing/main.go)|X
|Demuxing/Decoding|[see](examples/demuxing_decoding/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/demuxing_decoding.c)
|Filtering|[see](examples/filtering/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/filtering_video.c)
|Frame data manipulation|[see](examples/frame_data_manipulation/main.go)|X
|Hardware Decoding/Filtering|[see](examples/hardware_decoding_filtering/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/hw_decode.c)
|Hardware Encoding|[see](examples/hardware_encoding/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/vaapi_encode.c)
|Remuxing|[see](examples/remuxing/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/remuxing.c)
|Resampling audio|[see](examples/resampling_audio/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/resample_audio.c)
|Scaling video|[see](examples/scaling_video/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/scaling_video.c)
|Transcoding|[see](examples/transcoding/main.go)|[see](https://github.com/FFmpeg/FFmpeg/blob/n7.0/doc/examples/transcoding.c)

*Tip: you can use the video sample located in the `testdata` directory for your tests*

# Patterns

*NB: errors are not checked below for readibility purposes, however you should!*

First off, all use cases are different and it's impossible to provide patterns that work in every situation. That's why `ffmpeg`'s doc or source code should be your ultimate source of truth regarding how to use this library. That's why all methods of this library have been documented with a link referencing the documentation of the C function they use.

However it's possible to give rules of thumb and patterns that fit most use cases and can kickstart most people. Here's a few of them:

## When to call `Alloc()`, `.Unref()` and `.Free()`

Let's take the `FormatContext.ReadFrame()` pattern as an example. The pattern is similar with frames.

```go
// You can allocate the packet once and reuse the same object in the for loop below
pkt := astiav.AllocPacket()

// However, once you're done using the packet, you need to make sure to free it
defer pkt.Free()

// Loop
for {
    // We'll use a closure to ease unreferencing the packet
    func() {
        // Read frame using the same packet every time
        formatContext.ReadFrame(pkt)

        // However make sure to unreference the packet once you're done with what 
        // have been "injected" by the .ReadFrame() method
        defer pkt.Unref()

        // Here you can do whatever you feel like with your packet
    }()
}
```

# Breaking changes

You can see the list of breaking changes [here](BREAKING_CHANGES.md).

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

## Building on Windows

Building on Windows requires msys2 / mingw64 gcc toolchain. Read the [Quickstart guide](https://www.msys2.org) to install Msys2.

Once complete run the Mingw64 shell from the installation folder, run the below commands:

```shell
# Update Packages
pacman -Syu
# Install Requirements to Build
pacman -S --noconfirm --needed git diffutils mingw-w64-x86_64-toolchain pkg-config make yasm
# Clone the repository using git
git clone https://github.com/asticode/go-astiav
cd go-astiav
```

Then once you clone this repository, follow along the build instructions above.

> **Notes:**
> For `pkg-config` use `pkgconfiglite` from choco.
> Remember to set `CGO` and `PKG_CONFIG` env vars properly to point to the folder where ffmpeg was built.

# Why astiav?

After maintaining for several years the most starred [fork](https://github.com/asticode/goav) of [goav](https://github.com/giorgisio/goav), I've decided to write from scratch my own C bindings to fix most of the problems I still encountered using `goav`.