# gocurl

Go 编写的命令行 HTTP 客户端工具，功能类似 curl。

## 安装

下载对应平台的压缩包并解压：

- Windows: [gocurl_windows_amd64_v1.zip](https://github.com/iotzf/gocurl/releases/download/v1.1.0/gocurl_windows_amd64_v1.zip)
- Linux: [gocurl_linux_amd64_v1.zip](https://github.com/iotzf/gocurl/releases/download/v1.1.0/gocurl_linux_amd64_v1.zip)
- macOS: [gocurl_darwin_amd64_v1.zip](https://github.com/iotzf/gocurl/releases/download/v1.1.0/gocurl_darwin_amd64_v1.zip)

或使用 goreleaser 构建：

```bash
goreleaser build --clean
```

## 使用方法

```
gocurl [选项] <URL>
```

### 选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `-X, --method <method>` | HTTP 方法 | GET |
| `-H, --header <key:value>` | 自定义请求头 | |
| `-d, --data <data>` | 请求体数据 | |
| `-T, --content-type <type>` | Content-Type | application/json |
| `-t, --timeout <seconds>` | 超时时间（秒） | 30 |
| `-v, --verbose` | 打印请求头 | |
| `-k, --insecure` | 忽略 TLS 证书验证 | |
| `--proxy <url>` | 代理地址 | |

### 示例

发送 GET 请求：

```bash
gocurl https://httpbin.org/get
```

发送 POST 请求：

```bash
gocurl -X POST -d '{"name":"test"}' https://httpbin.org/post
```

携带自定义 Header：

```bash
gocurl -H "Authorization: Bearer token" https://httpbin.org/get
```

使用代理：

```bash
gocurl --proxy http://127.0.0.1:7890 https://httpbin.org/get
gocurl --proxy socks5://127.0.0.1:1080 https://httpbin.org/get
```

忽略证书验证：

```bash
gocurl -k https://self-signed.badssl.com/
```

打印请求头：

```bash
gocurl -v https://httpbin.org/get
```

## 环境变量

支持通过环境变量设置代理（`--proxy` 优先）：

- `HTTP_PROXY` / `HTTPS_PROXY` / `ALL_PROXY` — HTTP/HTTPS 代理
- `SOCKS_PROXY` — SOCKS5 代理

## 构建

```bash
go build -o gocurl.exe .
```

## License

MIT
