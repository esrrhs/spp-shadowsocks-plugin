# spp-shadowsocks-plugin

[<img src="https://img.shields.io/github/license/esrrhs/spp-shadowsocks-plugin">](https://github.com/esrrhs/spp-shadowsocks-plugin)
[<img src="https://img.shields.io/github/languages/top/esrrhs/spp-shadowsocks-plugin">](https://github.com/esrrhs/spp-shadowsocks-plugin)
[![Go Report Card](https://goreportcard.com/badge/github.com/esrrhs/spp-shadowsocks-plugin)](https://goreportcard.com/report/github.com/esrrhs/spp-shadowsocks-plugin)
[<img src="https://img.shields.io/github/v/release/esrrhs/spp-shadowsocks-plugin">](https://github.com/esrrhs/spp-shadowsocks-plugin/releases)
[<img src="https://img.shields.io/github/downloads/esrrhs/spp-shadowsocks-plugin/total">](https://github.com/esrrhs/spp-shadowsocks-plugin/releases)
[<img src="https://img.shields.io/github/actions/workflow/status/esrrhs/spp-shadowsocks-plugin/go.yml?branch=master">](https://github.com/esrrhs/spp-shadowsocks-plugin/actions)

[spp](https://github.com/esrrhs/spp)针对shadowsocks的插件
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

# 特性
* 支持协议tcp、kcp、quic，自定义协议rudp、rhttp、ricmp
* 支持加密压缩，默认关闭
* 支持Shadowsocks Android插件，[spp-shadowsocks-plugin-android](https://github.com/esrrhs/spp-shadowsocks-plugin-android)

# 编译
### 编译到非Android平台
```
# go build
```
### 或者
```
# ./pack.sh
```

### 编译到Android平台
#### 准备Java环境
* 安装jdk
```
# yum install java-1.8.0-openjdk-devel -y
```
* 设置环境变量，修改～/.bashrc
```
export JAVA_HOME=/usr/lib/jvm/java-1.8.0-openjdk/
export CLASSPATH=.:$JAVA_HOME/jre/lib/rt.jar:$JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar
export PATH=$PATH:$JAVA_HOME/bin
```
* 然后让它生效
```
# source ~/.bashrc
```

#### 准备Android SDK环境
* 在官方网站下载Command line tools only工具，[地址](https://developer.android.com/studio/#downloads)
```
# mkdir -p /home/project/android/cmdline-tools
# cd /home/project/android/cmdline-tools
# wget https://dl.google.com/android/repository/commandlinetools-linux-6609375_latest.zip
# unzip commandlinetools-linux-6609375_latest.zip 
```
* 设置环境变量，修改～/.bashrc
```
export ANDROID_SDK_ROOT=/home/project/android
```
* 然后让它生效
```
# source ~/.bashrc
```
* 安装android-sdk
```
# cd /home/project/android/cmdline-tools
# ./tools/bin/sdkmanager  "build-tools;28.0.3"  "platform-tools"  "platforms;android-28"
```

#### 准备Android NDK环境
* 在官方网站下载NDK，[地址](https://developer.android.com/ndk/downloads/index.html)
```
# wget https://dl.google.com/android/repository/android-ndk-r21b-linux-x86_64.zip
# unzip android-ndk-r21b-linux-x86_64.zip
# mv android-ndk-r21b /home/project
```
* 设置环境变量，修改～/.bashrc
```
export ANDROID_NDK_HOME=/home/project/android-ndk-r21b
```
* 然后让它生效
```
# source ~/.bashrc
```

#### 编译插件
* 编译spp-shadowsocks-plugin到Android各个架构，使用./pack_android.sh
```
# ./pack_android.sh
```

# 使用
### 非Android平台
* 下载shadowsocks的go版本[go-shadowsocks2](https://github.com/shadowsocks/go-shadowsocks2) 然后解压，放到和spp-shadowsocks-plugin同级目录
```
# gunzip shadowsocks2-linux.gz
```
* 执行shadowsocks2-linux客户端，附带spp-shadowsocks-plugin插件，采用rudp协议
```
# ./shadowsocks2-linux -c 'ss://AEAD_CHACHA20_POLY1305:your-password@127.0.0.1:8488' -verbose --plugin spp-shadowsocks-plugin --plugin-opts "proto=rudp" -socks :1080 
```
* 执行shadowsocks2-linux服务端，同样附带spp-shadowsocks-plugin插件，采用rudp协议，只不过加一个type=server表明是服务器
```
# ./shadowsocks2-linux -s 'ss://AEAD_CHACHA20_POLY1305:your-password@:8488' -verbose --plugin spp-shadowsocks-plugin --plugin-opts "type=server;proto=rudp"
```
* 完成，现在以socks5协议访问本地1080端口即可
```
# proxychains4 curl google.com
[proxychains] config file found: /usr/local/etc/proxychains.conf
[proxychains] preloading /usr/local/lib/libproxychains4.so
[proxychains] DLL init: proxychains-ng 4.14-git-23-g7fe8139
[proxychains] Strict chain  ...  127.0.0.1:1080  ...  google.com:80  ...  OK
<HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8">
<TITLE>301 Moved</TITLE></HEAD><BODY>
<H1>301 Moved</H1>
The document has moved
<A HREF="http://www.google.com/">here</A>.
</BODY></HTML>
```
### Android平台
* 客户端使用shadowsocks插件，[spp-shadowsocks-plugin-android](https://github.com/esrrhs/spp-shadowsocks-plugin-android)
