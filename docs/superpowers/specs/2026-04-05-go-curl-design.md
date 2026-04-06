# gocurl - Go 语言 HTTP 请求 CLI 工具

## 概述

用 Go 编写的命令行 HTTP 客户端工具，功能类似 curl，支持 GET/POST 请求、定制请求头、JSON/Form 数据、超时控制，用于调试和测试 HTTP API。

## 项目结构

```
.
├── main.go              # 入口，CLI 参数解析
├── internal/
│   ├── httpclient.go    # HTTP 请求执行
│   └── formatter.go     # 响应格式化输出
├── go.mod
└── docs/superpowers/specs/
    └── 2026-04-05-go-curl-design.md
```

## CLI 接口

```
用法: gocurl [选项] <URL>

选项:
  -X, --method string     HTTP 方法 (默认: GET)
  -H, --header strings    自定义请求头，格式: "Key: Value"
  -d, --data string       请求体数据
  -T, --content-type string  Content-Type (默认: application/json)
  -t, --timeout int       超时时间秒 (默认: 30)
  -h, --help              帮助信息
```

## 模块设计

### httpclient.go

- `NewClient(timeout int) *http.Client` — 创建配置好超时时间的 HTTP 客户端
- `DoRequest(method, url string, headers map[string]string, body io.Reader) (*Response, error)` — 执行 HTTP 请求
- `Response` 结构体包含：StatusCode (int)、Headers (http.Header)、Body ([]byte)

### formatter.go

- `FormatResponse(*Response) string` — 格式化输出响应
- 输出格式：状态行 + 响应头 + 空行 + 响应体
- JSON 自动美化/缩进显示
- 非 JSON 原样输出

## 错误处理

| 错误场景 | 退出码 | 处理方式 |
|----------|--------|----------|
| 网络连接失败 | 1 | 打印错误信息 |
| 请求超时 | 1 | 打印 timeout 错误 |
| HTTP 4xx/5xx | 0 | 输出完整响应（用于调试） |
| URL 解析失败 | 1 | 打印错误信息 |

## 依赖

- `github.com/alecthomas/kingpin/v2` — CLI 参数解析

## 验收标准

1. `gocurl https://httpbin.org/get` 输出完整响应信息
2. `gocurl -X POST -d '{"name":"test"}' https://httpbin.org/post` 发送 JSON 并查看响应
3. `gocurl -X POST -T "application/x-www-form-urlencoded" -d "name=test" https://httpbin.org/post` 发送 Form 数据
4. `gocurl -H "Authorization: Bearer token" https://httpbin.org/get` 携带自定义 Header
5. `gocurl -t 1 https://httpbin.org/delay/5` 超时控制生效（1秒超时）
