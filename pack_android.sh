#! /bin/bash

set -e
#set -x

VERSION=${VERSION:-$(grep -E '^\s*Version\s*=\s*"[^"]+"' version/version.go 2>/dev/null | sed -E 's/.*"([^"]+)".*/\1/')}
VERSION=${VERSION:-"0.8.0"}
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

TOOLCHAIN=$(find $ANDROID_NDK_HOME/toolchains/llvm/prebuilt/* -maxdepth 1 -type d -print -quit)/bin
echo "Toolchain: $TOOLCHAIN"
MIN_API=35
LDFLAGS="-s -w -extldflags=-Wl,-z,max-page-size=16384 -X 'github.com/esrrhs/spp-shadowsocks-plugin/version.Version=${VERSION}' -X 'github.com/esrrhs/spp-shadowsocks-plugin/version.GitCommit=${GIT_COMMIT}' -X 'github.com/esrrhs/spp-shadowsocks-plugin/version.BuildTime=${BUILD_TIME}'"

rm -f spp-shadowsocks-plugin_android_arm7.zip
rm -f spp-shadowsocks-plugin_android_arm64.zip
rm -f spp-shadowsocks-plugin_android_386.zip
rm -f spp-shadowsocks-plugin_android_amd64.zip
rm -f spp-shadowsocks-plugin

env CGO_ENABLED=1 CC=$TOOLCHAIN/armv7a-linux-androideabi${MIN_API}-clang GOOS=android GOARCH=arm GOARM=7 go build -ldflags="$LDFLAGS" -o spp-shadowsocks-plugin .
zip -q spp-shadowsocks-plugin_android_arm7.zip spp-shadowsocks-plugin
rm -f spp-shadowsocks-plugin

env CGO_ENABLED=1 CC=$TOOLCHAIN/aarch64-linux-android${MIN_API}-clang GOOS=android GOARCH=arm64 go build -ldflags="$LDFLAGS" -o spp-shadowsocks-plugin .
zip -q spp-shadowsocks-plugin_android_arm64.zip spp-shadowsocks-plugin
rm -f spp-shadowsocks-plugin

env CGO_ENABLED=1 CC=$TOOLCHAIN/i686-linux-android${MIN_API}-clang GOOS=android GOARCH=386 go build -ldflags="$LDFLAGS" -o spp-shadowsocks-plugin .
zip -q spp-shadowsocks-plugin_android_386.zip spp-shadowsocks-plugin
rm -f spp-shadowsocks-plugin

env CGO_ENABLED=1 CC=$TOOLCHAIN/x86_64-linux-android${MIN_API}-clang GOOS=android GOARCH=amd64 go build -ldflags="$LDFLAGS" -o spp-shadowsocks-plugin .
zip -q spp-shadowsocks-plugin_android_amd64.zip spp-shadowsocks-plugin
rm -f spp-shadowsocks-plugin

echo "Android builds completed successfully."
