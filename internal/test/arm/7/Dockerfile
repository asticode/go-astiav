FROM arm32v7/debian:12.5

RUN apt-get update

RUN apt-get install -y \
	build-essential \
	git \
	pkg-config \
	libpng-dev

RUN \
	mkdir -p /opt/ffmpeg/src
    
WORKDIR /opt/ffmpeg/src

RUN \
	git clone https://github.com/FFmpeg/FFmpeg /opt/ffmpeg/src && \
	git checkout n7.0

RUN \
	./configure --prefix=.. && \
	make && \
	make install

ADD https://dl.google.com/go/go1.22.0.linux-amd64.tar.gz /tmp/go.tar.gz
RUN tar -C /opt -xzf /tmp/go.tar.gz

ENV GOARCH=arm
ENV GOARM=7
ENV GOCACHE=/opt/gocache
ENV GOMODCACHE=/opt/gomodcache
ENV CGO_LDFLAGS=-L/opt/ffmpeg/lib/
ENV CGO_CFLAGS=-I/opt/ffmpeg/include/
ENV PKG_CONFIG_PATH=/opt/ffmpeg/lib/pkgconfig
ENV CGO_ENABLED=1
    
WORKDIR /opt/astiav

CMD ["/opt/go/bin/go", "test"]