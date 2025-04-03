FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod ./
# 只有当 go.sum 存在时才复制它
COPY go.sum* ./

# 初始化依赖（如果 go.sum 不存在）
RUN go mod tidy

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# 使用最小化的镜像
FROM alpine:latest

WORKDIR /app

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/server .

EXPOSE 8080

# 启动应用，默认不自动调整GOMAXPROCS
CMD ["./server"]