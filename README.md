# go_srs
srs的go版本

## 目录说明
|  目录  |  说明  |
|:------|:------|
| app | 应用层 |
| codec | 编解码器解析层 |
| global | 全局变量存放 |
| main | main入口 |
| protocol | 包含amf0协议，rtmp协议 |
| utils | 存放工具类 |
------

## protocol目录：
|  目录  |  说明  |
|:------|:------|
| amf0 | amf0协议实现 |
| packet | 信令包封包解包 |
| rtmp | rtmp协议，chunk，message |
| skt | 网络层 |

## 运行方法
go run main.go

obs推流地址：
* rtmp://ip:port/app/live?vhost=srs.net
* vhost对应在配置文件中配置的vhost

拉流地址：
* ffplay rtmp://ip:port/app/live?vhost=srs.net
* ffplay http://ip:port/app/live.flv?vhost=srs.net
* ffplay http://ip:port/hls/app/live.hls/vhost=srs.net

录制文件目录：
* go_srs/srs/main/html/app/xxx.hls
* go_srs/srs/main/html/app/xxx.flv
