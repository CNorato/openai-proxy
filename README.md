# OpenAI 转发代理
将官方 OpenAI 请求转发到第三方 OpenAI 服务的 HTTP 代理。
## Usage
```shell
$ ./openai-proxy
Usage of openai-proxy:
  -host string
        主机名
  -openai string
        OpenAI域名 (default "api.openai.com")
  -port string
        端口号 (default "6081")
  -target string
        目标域名 (default "https://api.gpt.ge")
```