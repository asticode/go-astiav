#!/usr/bin/env bash
set -euo pipefail

echo "========================================="
echo " FFmpeg 7.0 build for MSYS2/mingw "
echo "========================================="
echo

# -------------------------------------------------
# 0. MSYSTEM 체크 (mingw64에서만 실행되게)
# -------------------------------------------------
if [ "${MSYSTEM-}" != "MINGW64" ]; then
    echo "[ERROR] 이 스크립트는 MSYS2의 mingw64 쉘에서 실행해야 합니다."
    echo "        (C:\\msys64\\mingw64.exe 또는 MSYS2 시작 메뉴의 MSYS2 MinGW 64-bit)"
    exit 1
fi

# -------------------------------------------------
# 1. 설정값
# -------------------------------------------------
FFMPEG_VER="7.0.3"                        # 원하는 7.x 버전으로 바꿔도 됨
PREFIX="$HOME/ffmpeg-7.0-shared-win64"    # 설치 경로
SRC_DIR="$HOME/src"                       # 소스 다운받을 디렉토리

MINGW_PREFIX="/mingw64"
PKG_PREFIX="mingw-w64-x86_64"

echo "[INFO] FFMPEG_VER = $FFMPEG_VER"
echo "[INFO] PREFIX     = $PREFIX"
echo "[INFO] SRC_DIR    = $SRC_DIR"
echo

mkdir -p "$SRC_DIR"
cd "$SRC_DIR"

# -------------------------------------------------
# 2. pacman 업데이트 + 필요한 패키지 설치
# -------------------------------------------------
echo "[STEP 1] pacman 패키지 설치"

# 전체 업데이트 (이미 되어 있으면 금방 끝남)
pacman -Syu --noconfirm

# 기본 개발 툴
pacman -S --needed --noconfirm \
    base-devel \
    git \
    yasm \
    nasm \
    ${PKG_PREFIX}-toolchain \
    ${PKG_PREFIX}-pkgconf \
    ${PKG_PREFIX}-cmake

# FFmpeg 옵션에 맞는 라이브러리들
pacman -S --needed --noconfirm \
    ${PKG_PREFIX}-zlib \
    ${PKG_PREFIX}-fdk-aac \
    ${PKG_PREFIX}-fontconfig \
    ${PKG_PREFIX}-freetype \
    ${PKG_PREFIX}-fribidi \
    ${PKG_PREFIX}-lame \
    ${PKG_PREFIX}-libvorbis \
    ${PKG_PREFIX}-libvpx \
    ${PKG_PREFIX}-libwebp \
    ${PKG_PREFIX}-x264 \
    ${PKG_PREFIX}-x265

echo "[STEP 1] 패키지 설치 완료"
echo

# -------------------------------------------------
# 3. nv-codec-headers 설치 (NVENC용, 옵션 아님 거의 필수)
# -------------------------------------------------
echo "[STEP 2] nv-codec-headers 설치 (NVENC)"

cd "$SRC_DIR"
if [ ! -d "nv-codec-headers" ]; then
    git clone https://github.com/FFmpeg/nv-codec-headers.git
fi

cd nv-codec-headers
make clean || true
make PREFIX=/mingw64
make PREFIX=/mingw64 install

echo "[STEP 2] nv-codec-headers 설치 완료"
echo

# -------------------------------------------------
# 4. FFmpeg 7.x 소스 다운로드
# -------------------------------------------------
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

# -------------------------------------------------
# 5. configure
# -------------------------------------------------
echo "[STEP 4] configure 실행"

export PKG_CONFIG_PATH="$PREFIX/lib/pkgconfig:${MINGW_PREFIX}/lib/pkgconfig"
export PATH="${MINGW_PREFIX}/bin:${PATH}"

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

# -------------------------------------------------
# 6. 빌드 & 설치
# -------------------------------------------------
echo "[STEP 5] make / make install"

make -j"$(nproc || echo 4)"
make install

echo
echo "========================================="
echo " FFmpeg 빌드 완료!"
echo " PREFIX: $PREFIX"
echo " - bin: $PREFIX/bin (ffmpeg.exe, ffprobe.exe 등)"
echo " - lib: $PREFIX/lib (libavcodec.a, libavformat.a 등)"
echo " - pkgconfig: $PREFIX/lib/pkgconfig"
echo "========================================="
echo
echo "go-astiav 를 mingw64에서 쓸 때 예:"
echo
echo "  export PKG_CONFIG_PATH=\"$PREFIX/lib/pkgconfig:\$PKG_CONFIG_PATH\""
echo "  go env -w CGO_ENABLED=1"
echo "  go env -w CC=x86_64-w64-mingw32-gcc"
echo "  go env -w CXX=x86_64-w64-mingw32-g++"
echo "  go get github.com/asticode/go-astiav"
echo "  go build ./..."
echo
