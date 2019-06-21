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

## 握手协议
<script src="mermaid.min.js"></script>
<script>mermaid.initialize({startOnLoad:true});</script>
<div class="mermaid">
graph LR
    A --- B
    B-->C[fa:fa-ban forbidden]
    B-->D(fa:fa-spinner);
</div>