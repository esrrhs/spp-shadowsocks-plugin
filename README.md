# spp-shadowsocks-plugin

[<img src="https://img.shields.io/github/license/esrrhs/spp-shadowsocks-plugin">](https://github.com/esrrhs/spp-shadowsocks-plugin)
[<img src="https://img.shields.io/github/languages/top/esrrhs/spp-shadowsocks-plugin">](https://github.com/esrrhs/spp-shadowsocks-plugin)
[<img src="https://img.shields.io/github/v/release/esrrhs/spp-shadowsocks-plugin">](https://github.com/esrrhs/spp-shadowsocks-plugin/releases)
[<img src="https://img.shields.io/github/downloads/esrrhs/spp-shadowsocks-plugin/total">](https://github.com/esrrhs/spp-shadowsocks-plugin/releases)
[<img src="https://img.shields.io/github/actions/workflow/status/esrrhs/spp-shadowsocks-plugin/go.yml?branch=master">](https://github.com/esrrhs/spp-shadowsocks-plugin/actions)

A [SIP003](https://shadowsocks.org/doc/sip003.html) compliant plugin for Shadowsocks, powered by [SPP](https://github.com/esrrhs/spp).

```
     +------------+                    +---------------------------+
     |  SS Client +-- Local Loopback --+  Plugin Client (Tunnel)   +--+
     +------------+                    +---------------------------+  |
                                                                      |
                 Public Internet (Obfuscated/Transformed traffic) ==> |
                                                                      |
     +------------+                    +---------------------------+  |
     |  SS Server +-- Local Loopback --+  Plugin Server (Tunnel)   +--+
     +------------+                    +---------------------------+
```

---

## Features

* **Rich Underlay Protocols**: Supports standard `tcp`, `kcp`, `quic`, as well as custom reliable protocols `rudp`, `rhttp`, and `ricmp`.
* **Optional Encryption & Compression**: Seamlessly integrate data compression and encryption under the transport layer (disabled by default).
* **Cross-Platform**: Runs on Linux, macOS, Windows, FreeBSD, and Android.
* **Shadowsocks Android Support**: Compatible with Shadowsocks Android via [spp-shadowsocks-plugin-android](https://github.com/esrrhs/spp-shadowsocks-plugin-android).
* **SIP003 Compliant**: Standard Shadowsocks plugin interface with support for environment-variable-based configuration and VPN mode routing protection.

---

## Installation & Download

Download the pre-compiled binary for your operating system and architecture from [GitHub Releases](https://github.com/esrrhs/spp-shadowsocks-plugin/releases).

Extract the archive and place the `spp-shadowsocks-plugin` binary in your system `PATH` or the same directory as your Shadowsocks binary.

---

## Building from Source

### Non-Android Platforms

Ensure you have [Go](https://go.dev/) (1.24+) installed:

```bash
# Direct build
go build -ldflags="-s -w" -o spp-shadowsocks-plugin .

# Or build cross-platform release packages
./pack.sh
```

### Android Platform

#### 1. Setup Java Environment
* Install JDK:
```bash
yum install java-1.8.0-openjdk-devel -y
# or on Debian/Ubuntu: apt install openjdk-17-jdk -y
```

* Configure environment variables in `~/.bashrc`:
```bash
export JAVA_HOME=/usr/lib/jvm/java-1.8.0-openjdk/
export PATH=$PATH:$JAVA_HOME/bin
source ~/.bashrc
```

#### 2. Setup Android Command-Line Tools & SDK
* Download Command Line Tools from [Android Developer Studio](https://developer.android.com/studio/#downloads):
```bash
mkdir -p /home/project/android/cmdline-tools
cd /home/project/android/cmdline-tools
wget https://dl.google.com/android/repository/commandlinetools-linux-6609375_latest.zip
unzip commandlinetools-linux-6609375_latest.zip
```

* Configure environment variables in `~/.bashrc`:
```bash
export ANDROID_SDK_ROOT=/home/project/android
source ~/.bashrc
```

* Install SDK packages:
```bash
cd /home/project/android/cmdline-tools
./tools/bin/sdkmanager "build-tools;28.0.3" "platform-tools" "platforms;android-28"
```

#### 3. Setup Android NDK
* Download NDK from [Android NDK Downloads](https://developer.android.com/ndk/downloads/index.html):
```bash
wget https://dl.google.com/android/repository/android-ndk-r21b-linux-x86_64.zip
unzip android-ndk-r21b-linux-x86_64.zip -d /home/project/
```

* Configure environment variables in `~/.bashrc`:
```bash
export ANDROID_NDK_HOME=/home/project/android-ndk-r21b
source ~/.bashrc
```

#### 4. Build Plugin for Android
Compile `spp-shadowsocks-plugin` for all Android architectures (arm, arm64, 386, amd64):
```bash
./pack_android.sh
```

---

## Usage

### Non-Android Platforms

1. Download or compile Shadowsocks (for example, [go-shadowsocks2](https://github.com/shadowsocks/go-shadowsocks2)) and place it in the same directory as `spp-shadowsocks-plugin`.

2. **Start the Server**:
   Launch Shadowsocks server with the plugin configured with `type=server`:
   ```bash
   ./shadowsocks2-linux -s 'ss://AEAD_CHACHA20_POLY1305:your-password@:8488' \
     -verbose \
     --plugin spp-shadowsocks-plugin \
     --plugin-opts "type=server;proto=rudp"
   ```

3. **Start the Client**:
   Launch Shadowsocks client with the plugin:
   ```bash
   ./shadowsocks2-linux -c 'ss://AEAD_CHACHA20_POLY1305:your-password@127.0.0.1:8488' \
     -verbose \
     --plugin spp-shadowsocks-plugin \
     --plugin-opts "proto=rudp" \
     -socks :1080
   ```

4. **Verify Connection**:
   Connect via the local SOCKS5 proxy on port 1080:
   ```bash
   curl --socks5 127.0.0.1:1080 https://www.google.com
   ```

### Android Platform

On Android devices, install Shadowsocks and the corresponding plugin:
* Install [spp-shadowsocks-plugin-android](https://github.com/esrrhs/spp-shadowsocks-plugin-android) (available via GitHub Releases or Google Play).
* Configure the plugin name and options in Shadowsocks profile settings (e.g., `proto=rudp`).

---

## Plugin Options

Options can be passed through `--plugin-opts` (or the `SS_PLUGIN_OPTIONS` environment variable) as a semicolon-separated `key=value` string:

| Option | Type | Default | Description |
|---|---|---|---|
| `proto` | string | `tcp` | Transport protocol: `tcp`, `rudp`, `ricmp`, `rhttp`, `kcp`, `quic` |
| `type` | string | `client` | Role type: `client` or `server` |
| `compress` | int | `0` | Compression level: `0` (disabled), `1` (fastest), etc. |
| `encrypt` | string | `""` | Encryption method: `""` (disabled), `aes-128-gcm`, etc. |
| `key` | string | `""` | Pre-shared key for encryption |
| `pinginter` | int | `10` | Heartbeat interval (seconds) |
| `maxmsgsize` | int | `10240` | Maximum message buffer size |
| `maxclient` | int | `128` | Maximum concurrent clients (server mode) |
| `__android_vpn`| bool | `false`| Enables Android VPN routing protection socket |

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
