#! /bin/bash
set -x

CGO_ENABLED=0 go build
zip spp-shadowsocks-plugin_linux64.zip spp-shadowsocks-plugin

GOOS=darwin GOARCH=amd64 go build
zip spp-shadowsocks-plugin_mac.zip spp-shadowsocks-plugin

GOOS=windows GOARCH=amd64 go build
zip spp-shadowsocks-plugin_windows64.zip spp-shadowsocks-plugin.exe

GOOS=linux GOARCH=mipsle go build
zip spp-shadowsocks-plugin_mipsle.zip spp-shadowsocks-plugin

GOOS=linux GOARCH=arm go build
zip spp-shadowsocks-plugin_arm.zip spp-shadowsocks-plugin

GOOS=linux GOARCH=mips go build
zip spp-shadowsocks-plugin_mips.zip spp-shadowsocks-plugin

GOOS=windows GOARCH=386 go build
zip spp-shadowsocks-plugin_windows32.zip spp-shadowsocks-plugin.exe

GOOS=linux GOARCH=arm64 go build
zip spp-shadowsocks-plugin_arm64.zip spp-shadowsocks-plugin

GOOS=linux GOARCH=mips64 go build
zip spp-shadowsocks-plugin_mips64.zip spp-shadowsocks-plugin

GOOS=linux GOARCH=mips64le go build
zip spp-shadowsocks-plugin_mips64le.zip spp-shadowsocks-plugin

