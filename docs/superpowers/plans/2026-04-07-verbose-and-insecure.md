# Verbose and Insecure TLS Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `-v` flag to print request headers before sending, and `-k` flag to skip TLS certificate verification.

**Architecture:** Modify `main.go` to add two new kingpin flags, and modify `internal/httpclient.go` to accept verbose and insecure options when creating the HTTP client.

**Tech Stack:** Go stdlib `crypto/tls`, `net/http`

---

## File Structure

- Modify: `main.go` — add `-v` and `-k` flags
- Modify: `internal/httpclient.go` — add verbose printing and insecure TLS support

---

## Task 1: Add `-v` verbose flag to main.go

- [ ] **Step 1: Modify main.go to add verbose flag**

```go
var (
	method      = kingpin.Flag("method", "HTTP 方法").Short('X').Default("GET").String()
	headers     = kingpin.Flag("header", "自定义请求头").Short('H').StringMap()
	data        = kingpin.Flag("data", "请求体数据").Short('d').String()
	contentType = kingpin.Flag("content-type", "Content-Type").Short('T').Default("application/json").String()
	timeout     = kingpin.Flag("timeout", "超时时间秒").Short('t').Default("30").Int()
	verbose     = kingpin.Flag("verbose", "打印请求头").Short('v').Bool()
	insecure    = kingpin.Flag("insecure", "忽略证书验证").Short('k').Bool()
	url         = kingpin.Arg("url", "目标 URL").Required().String()
)
```

- [ ] **Step 2: Pass verbose and insecure to DoRequest**

```go
resp, err := httpclient.DoRequest(*method, *url, headerMap, bodyReader, *timeout, *verbose, *insecure)
```

---

## Task 2: Update httpclient.go to support verbose and insecure

- [ ] **Step 1: Add verbose and insecure parameters to NewClient**

```go
func NewClient(timeout int, insecure bool) *http.Client {
	tr := &http.Transport{}
	if insecure {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: tr,
	}
}
```

- [ ] **Step 2: Add verbose parameter to DoRequest**

```go
func DoRequest(method, url string, headers map[string]string, body io.Reader, timeout int, verbose, insecure bool) (*Response, error) {
	client := NewClient(timeout, insecure)

	// 打印请求头
	if verbose {
		fmt.Printf("> %s %s\n", method, url)
		for k, v := range headers {
			fmt.Printf("> %s: %s\n", k, v)
		}
		fmt.Println()
	}
	...
}
```

- [ ] **Step 3: Add crypto/tls import**

```go
import (
	"crypto/tls"
	...
)
```

---

## Task 3: Build and verify

- [ ] **Step 1: Build**

```bash
go build -o gocurl.exe .
```

- [ ] **Step 2: Test verbose flag**

```bash
./gocurl.exe -v https://httpbin.org/get
```

Expected: 输出请求头（`> GET https://httpbin.org/get` 和 `> Content-Type: application/json`）

- [ ] **Step 3: Test insecure flag**

```bash
./gocurl.exe -k https://self-signed.badssl.com/
```

Expected: 正常返回（不报证书错误）

- [ ] **Step 4: Commit**

```bash
git add main.go internal/httpclient.go
git commit -m "feat: add -v verbose and -k insecure flags"
```
