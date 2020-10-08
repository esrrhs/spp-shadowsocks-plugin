# spp-shadowsocks-plugin
[spp](https://github.com/esrrhs/spp)针对shadowsocks的插件

# 使用
* 下载shadowsocks的go版本[go-shadowsocks2](https://github.com/shadowsocks/go-shadowsocks2) 然后解压
```
# gunzip shadowsocks2-linux.gz
```
* 编译spp-shadowsocks-plugin，放到和shadowsocks2-linux同级目录
```
# go build
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
