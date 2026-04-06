package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gocurl/internal"

	"github.com/alecthomas/kingpin/v2"
)

var (
	method      = kingpin.Flag("method", "HTTP 方法").Short('X').Default("GET").String()
	headers     = kingpin.Flag("header", "自定义请求头").Short('H').StringMap()
	data        = kingpin.Flag("data", "请求体数据").Short('d').String()
	contentType = kingpin.Flag("content-type", "Content-Type").Short('T').Default("application/json").String()
	timeout     = kingpin.Flag("timeout", "超时时间秒").Short('t').Default("30").Int()
	url         = kingpin.Arg("url", "目标 URL").Required().String()
)

func main() {
	kingpin.Parse()

	// 构建请求头
	headerMap := make(map[string]string)
	for k, v := range *headers {
		headerMap[k] = v
	}
	if *contentType != "" {
		headerMap["Content-Type"] = *contentType
	}

	// 解析请求体
	var bodyReader io.Reader
	if *data != "" {
		bodyReader = strings.NewReader(*data)
	}

	// 发送请求
	resp, err := httpclient.DoRequest(*method, *url, headerMap, bodyReader, *timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// 输出响应
	fmt.Print(httpclient.FormatResponse(resp))
}
