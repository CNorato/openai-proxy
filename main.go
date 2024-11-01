package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
)

var DomainReplaceMap = make(map[string]*url.URL)

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	targetUrl := req.URL
	if targetDomain, ok := DomainReplaceMap[req.Host]; ok {
		targetUrl.Scheme = targetDomain.Scheme
		targetUrl.Host = targetDomain.Host
	}
	proxyReq, err := http.NewRequest(req.Method, targetUrl.String(), req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadGateway)
		return
	}

	proxyReq.Header.Set("Host", req.Host)
	proxyReq.Header.Set("X-Forwarded-For", req.RemoteAddr)

	for header, values := range req.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	response, err := client.Do(proxyReq)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadGateway)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadGateway)
		}
	}(response.Body)

	// 将响应的状态和头部复制到原始响应中
	for key, value := range response.Header {
		res.Header()[key] = value
	}
	res.WriteHeader(response.StatusCode)

	// 将响应内容复制到原始响应体中
	_, err = io.Copy(res, response.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadGateway)
		return
	}
}

func main() {
	var port string
	var host string
	var openAiDomain string
	var targetDomain string
	flag.StringVar(&port, "port", "6081", "端口号")
	flag.StringVar(&host, "host", "", "主机名")
	flag.StringVar(&openAiDomain, "openai", "api.openai.com", "OpenAI域名")
	flag.StringVar(&targetDomain, "target", "https://api.gpt.ge", "目标域名")
	flag.Parse()
	// 初始化域名映射
	target, _ := url.Parse(targetDomain)
	DomainReplaceMap[openAiDomain] = target
	// 设置监听的端口
	http.HandleFunc("/", handleRequestAndRedirect)
	log.Printf("Starting proxy server on %s:%s", host, port)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}
