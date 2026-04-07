# Proxy Support Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add proxy support for HTTP, HTTPS, and SOCKS5 protocols. Read proxy from environment variables (HTTP_PROXY, HTTPS_PROXY, ALL_PROXY, SOCKS_PROXY) with --proxy flag taking priority.

**Architecture:** Add a --proxy flag in main.go. Add a `GetProxy()` helper in httpclient.go that checks --proxy first, then falls back to environment variables (HTTP_PROXY, HTTPS_PROXY, ALL_PROXY, SOCKS_PROXY). Configure the HTTP transport with the appropriate proxy based on protocol (http vs socks5).

**Tech Stack:** Go stdlib (net/http, net, golang.org/x/net/proxy)

---

## File Structure

- Modify: `main.go` — add `--proxy` flag
- Modify: `internal/httpclient.go` — add proxy configuration to HTTP transport

---

## Task 1: Add --proxy flag to main.go

- [ ] **Step 1: Add proxy flag variable**

```go
var (
	method      = kingpin.Flag("method", "HTTP 方法").Short('X').Default("GET").String()
	headers     = kingpin.Flag("header", "自定义请求头").Short('H').StringMap()
	data        = kingpin.Flag("data", "请求体数据").Short('d').String()
	contentType = kingpin.Flag("content-type", "Content-Type").Short('T').Default("application/json").String()
	timeout     = kingpin.Flag("timeout", "超时时间秒").Short('t').Default("30").Int()
	verbose     = kingpin.Flag("verbose", "打印请求头").Short('v').Bool()
	insecure    = kingpin.Flag("insecure", "忽略证书验证").Short('k').Bool()
	proxy       = kingpin.Flag("proxy", "代理地址，如 http://127.0.0.1:7890 或 socks5://127.0.0.1:1080").String()
	url         = kingpin.Arg("url", "目标 URL").Required().String()
)
```

- [ ] **Step 2: Pass proxy to DoRequest**

```go
resp, err := httpclient.DoRequest(*method, *url, headerMap, bodyReader, *timeout, *verbose, *insecure, *proxy)
```

- [ ] **Step 3: Build to verify it compiles (will fail until Task 2)**

---

## Task 2: Update httpclient.go to support proxy

- [ ] **Step 1: Add golang.org/x/net/proxy import**

```go
import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/proxy"
)
```

- [ ] **Step 2: Add GetProxy helper function**

```go
// GetProxy returns the proxy URL to use, checking --proxy flag first then environment variables.
// Environment variables checked: HTTP_PROXY, HTTPS_PROXY, ALL_PROXY, SOCKS_PROXY
func GetProxy(explicitProxy string) string {
	if explicitProxy != "" {
		return explicitProxy
	}
	// Check environment variables in order of priority
	if p := os.Getenv("SOCKS_PROXY"); p != "" {
		return p
	}
	if p := os.Getenv("ALL_PROXY"); p != "" {
		return p
	}
	if p := os.Getenv("HTTPS_PROXY"); p != "" {
		return p
	}
	if p := os.Getenv("HTTP_PROXY"); p != "" {
		return p
	}
	return ""
}
```

- [ ] **Step 3: Update NewClient to accept proxy and configure transport**

```go
func NewClient(timeout int, insecure bool, proxyURL string) *http.Client {
	tr := &http.Transport{}

	// Configure proxy
	if proxyURL != "" {
		u, err := url.Parse(proxyURL)
		if err == nil {
			if u.Scheme == "socks5" {
				// SOCKS5: use golang.org/x/net/proxy
				dialer, err := proxy.SOCKS5("tcp", u.Host, nil, proxy.Direct)
				if err == nil {
					tr.Dial = dialer.Dial
				}
			} else if u.Scheme == "http" || u.Scheme == "https" {
				// HTTP/HTTPS: use http.ProxyURL
				tr.Proxy = http.ProxyURL(u)
			}
		}
	}

	if insecure {
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{}
		}
		tr.TLSClientConfig.InsecureSkipVerify = true
	}

	return &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: tr,
	}
}
```

- [ ] **Step 4: Update DoRequest signature and use proxy**

```go
func DoRequest(method, url string, headers map[string]string, body io.Reader, timeout int, verbose, insecure bool, proxy string) (*Response, error) {
	client := NewClient(timeout, insecure, proxy)
	...
}
```

---

## Task 3: Build and verify

- [ ] **Step 1: Install proxy dependency**

```bash
go get golang.org/x/net/proxy
```

- [ ] **Step 2: Build**

```bash
go build -o gocurl.exe .
```

- [ ] **Step 3: Test with HTTP proxy**

```bash
./gocurl.exe --proxy http://127.0.0.1:7890 https://httpbin.org/get
```

- [ ] **Step 4: Test with SOCKS5 proxy**

```bash
./gocurl.exe --proxy socks5://127.0.0.1:1080 https://httpbin.org/get
```

- [ ] **Step 5: Test environment variable fallback** (set HTTP_PROXY and test without --proxy)

- [ ] **Step 6: Commit**

```bash
git add main.go internal/httpclient.go go.mod go.sum
git commit -m "feat: add proxy support for HTTP and SOCKS5"
```
