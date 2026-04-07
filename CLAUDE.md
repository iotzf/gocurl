# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

gocurl 是一个用 Go 编写的命令行 HTTP 客户端工具，功能类似 curl，支持 GET/POST 请求、自定义请求头、JSON/Form 数据提交、超时控制。

## 构建和运行

```bash
go build -o gocurl.exe .
./gocurl.exe https://httpbin.org/get
```

## CLI 用法

```
gocurl [选项] <URL>

选项:
  -X, --method string     HTTP 方法 (默认: GET)
  -H, --header strings   自定义请求头，格式: "Key: Value"
  -d, --data string      请求体数据
  -T, --content-type string  Content-Type (默认: application/json)
  -t, --timeout int      超时时间秒 (默认: 30)
```

## 架构

- `main.go` — 入口文件，CLI 参数解析（使用 kingpin）
- `internal/httpclient.go` — HTTP 请求执行和响应格式化

核心函数：
- `httpclient.DoRequest(method, url string, headers map[string]string, body io.Reader, timeout int) (*Response, error)`
- `httpclient.FormatResponse(resp *Response) string` — 自动识别 JSON 并美化输出

## 依赖

- `github.com/alecthomas/kingpin/v2` — CLI 参数解析
