# 第一阶段：构建阶段
FROM golang:1.24-alpine as builder

# 设置 Go Proxy 使用国内代理
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 设置容器内的工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载所有依赖项
RUN go mod tidy

# 复制源代码
COPY . .

# 构建可执行文件
RUN go build -o main .

# 第二阶段：运行阶段
FROM alpine:3.19

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件和配置文件
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/config.local.yaml .

# 暴露端口
EXPOSE 8080

# 运行应用程序
CMD ["./main"]