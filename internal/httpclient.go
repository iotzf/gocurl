package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// Response 封装 HTTP 响应
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

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

// NewClient 创建配置好超时时间的 HTTP 客户端
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

// DoRequest 执行 HTTP 请求
func DoRequest(method, url string, headers map[string]string, body io.Reader, timeout int, verbose, insecure bool, proxy string) (*Response, error) {
	resolvedProxy := GetProxy(proxy)
	client := NewClient(timeout, insecure, resolvedProxy)

	// 打印请求头
	if verbose {
		fmt.Printf("> %s %s\n", method, url)
		for k, v := range headers {
			fmt.Printf("> %s: %s\n", k, v)
		}
		fmt.Println()
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = body
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