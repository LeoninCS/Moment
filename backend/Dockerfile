# ========= 构建阶段 =========
FROM golang:1.24-alpine AS builder

# 安装编译需要的工具
RUN apk add --no-cache git

WORKDIR /app

# 配置 Go 模块代理，避免访问不到 proxy.golang.org
ENV GOPROXY=https://goproxy.cn,direct

# 先复制 go.mod / go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 再复制剩下代码
COPY . .

# 编译可执行文件（路径按你实际 main.go 调整）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server ./cmd/main.go

# ========= 运行阶段 =========
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
