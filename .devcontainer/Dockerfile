FROM debian:latest

RUN apt-get update

##
# git
##

RUN apt-get install -y git nano

##
# ffmpeg
##

ARG FFMPEG_VERSION
                
RUN apt-get install -y \
    build-essential \
    nasm \
    pkg-config \
    libx264-dev \
    libpng-dev

RUN mkdir -p /opt/ffmpeg \
    && git clone https://github.com/FFmpeg/FFmpeg /opt/ffmpeg/src \
    && cd /opt/ffmpeg/src \
    && git checkout ${FFMPEG_VERSION}

RUN cd /opt/ffmpeg/src \
    && ./configure --prefix=.. --enable-libx264 --enable-gpl \
    && make \
    && make install

##
# go
##

ARG GO_VERSION

RUN apt-get install -y wget
RUN <<_EOF_ sh
arch=$(dpkg --print-architecture)
goArch="amd64"
case \${arch} in
    arm64) goArch="arm64" ;;
esac
wget -O /tmp/go.tar.gz https://dl.google.com/go/go${GO_VERSION}.linux-\${goArch}.tar.gz
tar -C /opt -xzf /tmp/go.tar.gz
_EOF_
ENV PATH="$PATH:/opt/go/bin"