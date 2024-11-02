FROM ubuntu:latest
LABEL authors="norato"
# 使用官方的 Golang 镜像创建构建产物。
FROM golang:1.23.2 AS builder

RUN mkdir /app
# 将本地代码复制到容器镜像中。
WORKDIR /app
COPY . .

# 在容器内构建命令。
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o openai-proxy .

# 使用一个新的阶段创建一个最小的镜像。
FROM alpine:3.20
LABEL authors="norato"
COPY --from=builder /app/openai-proxy /usr/local/bin/openai-proxy
# 更新文件权限以确保它是可执行的。
RUN chmod +x /usr/local/bin/openai-proxy
# 设置容器的默认端口
EXPOSE 6081

# 设置容器的默认命令。
ENTRYPOINT ["openai-proxy"]