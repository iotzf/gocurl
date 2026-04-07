package httpclient

import (
	"bytes"
	"crypto/tls"
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

// DoRequest 执行 HTTP 请求
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