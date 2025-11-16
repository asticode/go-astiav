#!/usr/bin/env bash
set -euo pipefail

echo "========================================="
echo " FFmpeg 7.x build for Ubuntu      "
echo "========================================="
echo

# -----------------------------------------
# 0. 설정
# -----------------------------------------
FFMPEG_VER="7.0.3"                     # 필요하면 7.x 다른 버전으로 변경
PREFIX="$HOME/ffmpeg-7.0-shared-linux"       # 설치 경로
SRC_DIR="$HOME/src"                    # 소스 내려받는 위치

echo "[INFO] FFMPEG_VER = $FFMPEG_VER"
echo "[INFO] PREFIX     = $PREFIX"
echo "[INFO] SRC_DIR    = $SRC_DIR"
echo

mkdir -p "$SRC_DIR"
cd "$SRC_DIR"

# -----------------------------------------
# 1. 필수 패키지 설치 (sudo 필요)
# -----------------------------------------
echo "[STEP 1] apt 패키지 설치"

# multiverse 활성화 (libfdk-aac-dev 위해)
if ! grep -qi "multiverse" /etc/apt/sources.list /etc/apt/sources.list.d/*.list 2>/dev/null; then
    echo "[INFO] multiverse repo 추가"
    sudo apt-get update
    sudo apt-get install -y software-properties-common
    sudo add-apt-repository -y multiverse
fi

sudo apt-get update

sudo apt-get install -y \
    build-essential \
    pkg-config \
    cmake \
    git \
    curl \
    wget \
    yasm \
    nasm \
    autoconf \
    automake \
    libtool \
    texinfo \
    zlib1g-dev \
    libfontconfig1-dev \
    libfreetype6-dev \
    libfribidi-dev \
    libfdk-aac-dev \
    libmp3lame-dev \
    libvorbis-dev \
    libvpx-dev \
    libwebp-dev \
    libx264-dev \
    libx265-dev \
    libnuma-dev

echo "[STEP 1] 패키지 설치 완료"
echo

# -----------------------------------------
# 2. nv-codec-headers 설치 (NVENC용)
# -----------------------------------------
echo "[STEP 2] nv-codec-headers 설치 (NVENC)"

cd "$SRC_DIR"
if [ ! -d "nv-codec-headers" ]; then
    git clone https://github.com/FFmpeg/nv-codec-headers.git
fi

cd nv-codec-headers
make clean || true
make
sudo make install   # /usr/local/include/ffnvcodec 에 설치

echo "[STEP 2] nv-codec-headers 설치 완료"
echo

# -----------------------------------------
# 3. FFmpeg 7.x 소스 다운로드
# -----------------------------------------
echo "[STEP 3] FFmpeg ${FFMPEG_VER} 소스 다운로드"

cd "$SRC_DIR"

FFMPEG_TAR="ffmpeg-${FFMPEG_VER}.tar.xz"
FFMPEG_DIR="ffmpeg-${FFMPEG_VER}"

if [ ! -f "$FFMPEG_TAR" ]; then
    curl -LO "https://ffmpeg.org/releases/${FFMPEG_TAR}"
fi

if [ -d "$FFMPEG_DIR" ]; then
    echo "[INFO] 기존 디렉토리 $FFMPEG_DIR 이(가) 있어 그대로 사용합니다."
else
    tar xf "$FFMPEG_TAR"
fi

cd "$FFMPEG_DIR"

echo "[STEP 3] FFmpeg 소스 준비 완료"
echo

# -----------------------------------------
# 4. configure (build 설정)
# -----------------------------------------
echo "[STEP 4] configure 실행"

export PKG_CONFIG_PATH="$PREFIX/lib/pkgconfig:${PKG_CONFIG_PATH-}"

./configure \
  --prefix="$PREFIX" \
  --extra-libs="-lpthread -lm" \
  --enable-gpl \
  --enable-version3 \
  --enable-nonfree \
  --disable-debug \
  --disable-doc \
  --enable-ffmpeg \
  --enable-ffprobe \
  --disable-ffplay \
  --enable-shared \
  --enable-libfdk-aac \
  --enable-libfontconfig \
  --enable-libfreetype \
  --enable-libfribidi \
  --enable-libmp3lame \
  --enable-libvorbis \
  --enable-libvpx \
  --enable-libwebp \
  --enable-libx264 \
  --enable-libx265 \
  --enable-nvenc

echo "[STEP 4] configure 완료"
echo

# -----------------------------------------
# 5. 빌드 & 설치
# -----------------------------------------
echo "[STEP 5] make / make install"

make -j"$(nproc || echo 4)"
make install

echo
echo "========================================="
echo " FFmpeg 빌드 완료!"
echo " PREFIX: $PREFIX"
echo " - bin: $PREFIX/bin (ffmpeg, ffprobe 등)"
echo " - lib: $PREFIX/lib (libavcodec.a, libavformat.a 등)"
echo " - pkgconfig: $PREFIX/lib/pkgconfig"
echo "========================================="
echo
echo "go-astiav 쓸 때 예:"
echo
echo "  export PKG_CONFIG_PATH=\"$PREFIX/lib/pkgconfig:\$PKG_CONFIG_PATH\""
echo "  go env -w CGO_ENABLED=1"
echo "  go get github.com/asticode/go-astiav"
echo "  go build ./..."
echo

