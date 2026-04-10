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
	"path"
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

// DownloadFile downloads a file with progress bar
func DownloadFile(method, rawURL string, headers map[string]string, body io.Reader, timeout int, verbose, insecure bool, proxy string, outputFilename string) error {
	resolvedProxy := GetProxy(proxy)
	client := NewClient(timeout, insecure, resolvedProxy)

	// Print request headers
	if verbose {
		fmt.Printf("> %s %s\n", method, rawURL)
		for k, v := range headers {
			fmt.Printf("> %s: %s\n", k, v)
		}
		fmt.Println()
	}

	var bodyReader io.Reader
	if body != nil {
		bodyReader = body
	}

	req, err := http.NewRequest(method, rawURL, bodyReader)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Get total size
	totalSize := resp.ContentLength

	// Create progress bar
	pb := NewProgressBar(totalSize, outputFilename)

	// Check for resume support
	fileInfo, _ := os.Stat(outputFilename)
	var startPos int64 = 0
	if fileInfo != nil && fileInfo.Size() > 0 {
		// Partial file exists, check if server supports resume
		if resp.Header.Get("Accept-Ranges") == "bytes" {
			startPos = fileInfo.Size()
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startPos))
			pb.SetDownloaded(startPos)
			// Re-do request with Range header
			resp.Body.Close()
			resp, err = client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
		}
	}

	// Open output file
	outFile, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("cannot open output file: %w", err)
	}

	// Ensure file is truncated if not resuming
	if startPos == 0 {
		outFile.Truncate(0)
		outFile.Seek(0, 0)
	}

	// Create a writer that updates progress
	writer := &progressWriter{
		file:     outFile,
		progress: pb,
	}

	// Copy with progress tracking
	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		nr, readErr := resp.Body.Read(buf)
		if nr > 0 {
			// Write to file and update progress
			_, writeErr := writer.Write(buf[0:nr])
			if writeErr != nil {
				outFile.Close()
				return fmt.Errorf("write error: %w", writeErr)
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			outFile.Close()
			return fmt.Errorf("read error: %w", readErr)
		}
	}

	pb.Finish()
	outFile.Close()

	return nil
}

type progressWriter struct {
	file     *os.File
	progress *ProgressBar
}

func (pw *progressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.file.Write(p)
	pw.progress.downloaded += int64(n)
	return
}

// ExtractFilenameFromURL extracts filename from URL path
func ExtractFilenameFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "gocurl_download"
	}
	filename := path.Base(u.Path)
	// URL decode
	filename, err = url.QueryUnescape(filename)
	if err != nil || filename == "" || filename == "." || filename == "/" {
		return "gocurl_download"
	}
	return filename
}