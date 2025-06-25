package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// 创建自定义TLS配置
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		fmt.Println("加载系统根证书失败:", err)
		return
	}

	tlsConfig := &tls.Config{
		RootCAs: rootCAs,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Get("https://www.example.com")
	if err != nil {
		fmt.Println("HTTPS请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	fmt.Println("响应状态:", resp.Status)
	fmt.Println("响应内容前100字节:", string(body[:100]))
}
