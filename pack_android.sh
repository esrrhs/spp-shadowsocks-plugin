#! /bin/bash

TOOLCHAIN=$(find $ANDROID_NDK_HOME/toolchains/llvm/prebuilt/* -maxdepth 1 -type d -print -quit)/bin
echo $TOOLCHAIN
MIN_API=21

env CGO_ENABLED=1 CC=$TOOLCHAIN/armv7a-linux-androideabi${MIN_API}-clang GOOS=android GOARCH=arm GOARM=7 go build -ldflags="-s -w"
zip spp-shadowsocks-plugin_android_arm7.zip spp-shadowsocks-plugin

env CGO_ENABLED=1 CC=$TOOLCHAIN/aarch64-linux-android${MIN_API}-clang GOOS=android GOARCH=arm64 go build -ldflags="-s -w"
zip spp-shadowsocks-plugin_android_arm64.zip spp-shadowsocks-plugin

env CGO_ENABLED=1 CC=$TOOLCHAIN/i686-linux-android${MIN_API}-clang GOOS=android GOARCH=386 go build -ldflags="-s -w"
zip spp-shadowsocks-plugin_android_386.zip spp-shadowsocks-plugin

env CGO_ENABLED=1 CC=$TOOLCHAIN/x86_64-linux-android${MIN_API}-clang GOOS=android GOARCH=amd64 go build -ldflags="-s -w"
zip spp-shadowsocks-plugin_android_amd64.zip spp-shadowsocks-plugin
