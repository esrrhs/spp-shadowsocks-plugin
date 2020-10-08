# spp-shadowsocks-plugin
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
# 编译
## 编译到非Android平台
```
# go build
```
## 或者其他平台
```
# ./pack.sh
```

## 编译到Android平台
### Java环境
* 下载Java，解压到某个目录
```
# wget http://javadl.oracle.com/webapps/download/AutoDL?BundleId=225345_090f390dda5b47b9b721c7dfaa008135 -O  jre-8u144-linux-x64.tar.gz
# tar zxvf jre-8u144-linux-x64.tar.gz
# mkdir -p /home/project/
# mv jre1.8.0_144 /home/project/
```
* 设置环境变量，修改～/.bashrc
```
export JAVA_HOME=/home/project/jre1.8.0_144
export CLASSPATH=.:$JAVA_HOME/jre/lib/rt.jar:$JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar
export PATH=$PATH:$JAVA_HOME/bin
```
* 然后让它生效
```
# source ~/.bashrc
```

### Android SDK环境
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

### Android NDK环境
* 在官方网站下载NDK，[地址](https://developer.android.com/ndk/downloads/index.html)
```
# wget https://dl.google.com/android/repository/android-ndk-r21b-linux-x86_64.zip
# unzip android-ndk-r21b-linux-x86_64.zip
# mv android-ndk-r21b /home/project
```
* 设置环境变量，修改～/.bashrc
```
export ANDROID_NDK_HOME=/home/project/android-ndk/android-ndk-r21b
```
* 然后让它生效
```
# source ~/.bashrc
```

### 编译插件
* 编译spp-shadowsocks-plugin到Android各个架构，使用./pack_android.sh
```
# ./pack_android.sh
```

# 使用
* 下载shadowsocks的go版本[go-shadowsocks2](https://github.com/shadowsocks/go-shadowsocks2) 然后解压，放到和spp-shadowsocks-plugin同级目录
```
# gunzip shadowsocks2-linux.gz
```
* 执行shadowsocks2-linux客户端，附带spp-shadowsocks-plugin插件，采用rudp协议
```
# ./shadowsocks2-linux -c 'ss://AEAD_CHACHA20_POLY1305:your-password@127.0.0.1:8488' -verbose --plugin spp-shadowsocks-plugin --plugin-opts "type=client;proto=rudp" -socks :1080 
```
* 执行shadowsocks2-linux服务端，同样附带spp-shadowsocks-plugin插件，只不过type为server
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
