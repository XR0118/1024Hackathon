# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o operator-pm-agent ./cmd/operator-pm-agent

# 运行阶段
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata docker-cli bash curl

# 创建非root用户
RUN addgroup -g 1001 -S boreas && \
	adduser -u 1001 -S boreas -G boreas

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/operator-pm-agent .

# 创建必要目录
RUN mkdir -p /var/lib/boreas-agent/apps && \
	mkdir -p /etc/boreas-agent && \
	chown -R boreas:boreas /var/lib/boreas-agent && \
	chown -R boreas:boreas /app

# 复制默认配置文件
COPY cmd/operator-pm-agent/configs/agent.yaml /etc/boreas-agent/agent.yaml
RUN chown boreas:boreas /etc/boreas-agent/agent.yaml

# 切换到非root用户
USER boreas

# 暴露端口
EXPOSE 8081

# 设置环境变量
ENV WORK_DIR=/var/lib/boreas-agent
ENV HOST=0.0.0.0
ENV PORT=8081

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
	CMD curl -f http://localhost:8081/v1/health || exit 1

# 启动命令
CMD ["./operator-pm-agent", "--config", "/etc/boreas-agent/agent.yaml"]
