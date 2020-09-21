# portscan
在做内网做渗透的时候，可以上传该工具进行tcp端口开放的发现工作，来查看网内是否有代理进入的必要。

## Usage

编译：
```go
go build fofa_api
```

```go
Usage of portscan:
  -i string
        输入ip地址 eg:192.0.0.1-192.0.0.255
  -p string
        端口列表 eg:22,80,1433,3306,3389 (default "22,80,1433,3306,3389")
  -t int
        线程数量 默认50 (default 300)
```
        ### V1.0.0

