# 使用官方 Go 运行时作为父镜像
FROM golang:1.24-alpine as builder

# 设置 Go Proxy 使用国内代理
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 设置容器内的工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载所有依赖项。如果 go.mod 和 go.sum 没有更改，依赖项将被缓存
RUN go mod tidy

# 将源代码复制到容器中
COPY . .

# 构建 Go 应用
RUN go build -o main .

# 暴露 8080 端口
EXPOSE 8080

# 运行编译后的二进制文件
CMD ["./main"]
