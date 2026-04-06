# gocurl Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建一个 Go 语言 HTTP 请求 CLI 工具，功能类似 curl

**Architecture:** 单二进制 CLI，使用 kingpin 做参数解析，internal 包分离 HTTP 客户端和格式化逻辑

**Tech Stack:** Go 1.21+, github.com/alecthomas/kingpin/v2

---

## File Structure

```
.
├── go.mod                           # 模块定义，依赖 kingpin
├── main.go                          # CLI 入口，参数解析
├── internal/
│   ├── httpclient.go                # HTTP 请求执行和格式化输出
│   └── httpclient_test.go           # httpclient 单元测试
└── docs/superpowers/specs/
    └── 2026-04-05-go-curl-design.md
```

---

## Task 1: 初始化项目

**Files:**
- Create: `go.mod`

- [ ] **Step 1: 初始化 go.mod**

Run:
```bash
cd C:/Users/cedar/Desktop/demo22 && go mod init gocurl
```

- [ ] **Step 2: 添加 kingpin 依赖**

Run:
```bash
cd C:/Users/cedar/Desktop/demo22 && go get github.com/alecthomas/kingpin/v2@latest
```

- [ ] **Step 3: 提交**

```bash
git add go.mod
git commit -m "init: add go.mod with kingpin dependency"
```

---

## Task 2: 实现 httpclient（含 formatter）

**Files:**
- Create: `internal/httpclient.go`
- Create: `internal/httpclient_test.go`

- [ ] **Step 1: 编写 Response 结构体和接口测试**

```go
// internal/httpclient_test.go
package httpclient

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestNewClient(t *testing.T) {
    client := NewClient(30)
    if client.Timeout.Seconds() != 30 {
        t.Errorf("expected timeout 30s, got %v", client.Timeout)
    }
}

func TestDoRequest_Get(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        w.Write([]byte(`{"method":"GET","url":"/test"}`))
    }))
    defer ts.Close()

    resp, err := DoRequest("GET", ts.URL, nil, nil)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if resp.StatusCode != 200 {
        t.Errorf("expected status 200, got %d", resp.StatusCode)
    }
    if string(resp.Body) != `{"method":"GET","url":"/test"}` {
        t.Errorf("unexpected body: %s", resp.Body)
    }
}
```

- [ ] **Step 2: 运行测试验证失败**

Run:
```bash
cd C:/Users/cedar/Desktop/demo22 && go test ./internal/... -v
```
Expected: FAIL — functions not defined

- [ ] **Step 3: 实现 minimal code**

```go
// internal/httpclient.go
package httpclient

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

// Response 封装 HTTP 响应
type Response struct {
    StatusCode int
    Headers    http.Header
    Body       []byte
}

// NewClient 创建配置好超时时间的 HTTP 客户端
func NewClient(timeout int) *http.Client {
    return &http.Client{
        Timeout: time.Duration(timeout) * time.Second,
    }
}

// DoRequest 执行 HTTP 请求
func DoRequest(method, url string, headers map[string]string, body []byte) (*Response, error) {
    client := NewClient(30)

    var bodyReader io.Reader
    if body != nil {
        bodyReader = bytes.NewReader(body)
    }

    req, err := http.NewRequest(method, url, bodyReader)
    if err != nil {
        return nil, err
    }

    for k, v := range headers {
        req.Header.Set(k, v)
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return &Response{
        StatusCode: resp.StatusCode,
        Headers:    resp.Header,
        Body:       bodyBytes,
    }, nil
}

// FormatResponse 格式化输出响应
func FormatResponse(resp *Response) string {
    var buf bytes.Buffer

    // 状态行
    buf.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode)))

    // 响应头
    for k, values := range resp.Headers {
        for _, v := range values {
            buf.WriteString(fmt.Sprintf("%s: %s\n", k, v))
        }
    }
    buf.WriteString("\n")

    // 响应体 - JSON 格式化
    body := string(resp.Body)
    if isJSON(resp.Headers.Get("Content-Type")) {
        var prettyJSON bytes.Buffer
        if err := json.Indent(&prettyJSON, resp.Body, "", "  "); err == nil {
            body = prettyJSON.String()
        }
    }
    buf.WriteString(body)

    return buf.String()
}

func isJSON(contentType string) bool {
    return strings.Contains(contentType, "application/json")
}
```

- [ ] **Step 4: 运行测试验证通过**

Run:
```bash
cd C:/Users/cedar/Desktop/demo22 && go test ./internal/... -v
```
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add internal/httpclient.go internal/httpclient_test.go
git commit -m "feat: add httpclient package with DoRequest and FormatResponse"
```

---

## Task 3: 实现 main.go CLI

**Files:**
- Create: `main.go`

- [ ] **Step 1: 编写 CLI 入口**

```go
// main.go
package main

import (
    "fmt"
    "os"

    "gocurl/internal/httpclient"

    "github.com/alecthomas/kingpin/v2"
)

func main() {
    app := kingpin.New("gocurl", "HTTP request CLI tool like curl")

    method := app.Flag("method", "HTTP method").Short('X').Default("GET").String()
    headers := app.Flag("header", "Custom header").Short('H').StringMap()
    data := app.Flag("data", "Request body").Short('d').String()
    contentType := app.Flag("content-type", "Content-Type").Short('T').Default("application/json").String()
    app.Flag("timeout", "Timeout in seconds").Short('t').Default("30").Int()

    url := app.Arg("url", "Target URL").String()

    app.HelpFlag('h', "Show help")

    kingpin.Parse()

    if *url == "" {
        app.Usage()
        os.Exit(1)
    }

    // 构建请求头
    headerMap := make(map[string]string)
    for k, v := range *headers {
        headerMap[k] = v
    }
    if *contentType != "" {
        headerMap["Content-Type"] = *contentType
    }

    // 解析请求体
    var body []byte
    if *data != "" {
        body = []byte(*data)
    }

    // 发送请求
    resp, err := httpclient.DoRequest(*method, *url, headerMap, body)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // 输出响应
    fmt.Print(httpclient.FormatResponse(resp))
}
```

- [ ] **Step 2: 编译验证**

Run:
```bash
cd C:/Users/cedar/Desktop/demo22 && go build -o gocurl.exe .
```
Expected: 编译成功，生成 gocurl.exe

- [ ] **Step 3: 提交**

```bash
git add main.go
git commit -m "feat: add CLI entry point with kingpin"
```

---

## Task 4: 验收测试

**验证 spec 中的验收标准：**

- [ ] **Test 1: GET 请求**

Run:
```bash
./gocurl.exe https://httpbin.org/get
```
Expected: 输出完整响应信息（状态码、头、body）

- [ ] **Test 2: POST JSON**

Run:
```bash
./gocurl.exe -X POST -d "{\"name\":\"test\"}" https://httpbin.org/post
```
Expected: 发送 JSON 并查看响应

- [ ] **Test 3: POST Form**

Run:
```bash
./gocurl.exe -X POST -T "application/x-www-form-urlencoded" -d "name=test" https://httpbin.org/post
```
Expected: 发送 Form 数据

- [ ] **Test 4: 自定义 Header**

Run:
```bash
./gocurl.exe -H "Authorization: Bearer token" https://httpbin.org/get
```
Expected: 携带自定义 Header

- [ ] **Test 5: 超时控制**

Run:
```bash
./gocurl.exe -t 1 https://httpbin.org/delay/5
```
Expected: 1秒超时后报错退出

---

## 自查清单

1. **Spec coverage:** 所有 spec 中的功能都有对应实现
2. **Placeholder scan:** 无 TBD/TODO/占位符
3. **Type consistency:** Response 结构体只定义一次，在 httpclient 包中
