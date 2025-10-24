.PHONY: build run test clean docker-build docker-run help

# 项目配置
PROJECT_NAME := boreas
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | awk '{print $$3}')

# 构建标志
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GoVersion=$(GO_VERSION)"

# 默认目标
.DEFAULT_GOAL := help

# 构建所有服务
build:
	@echo "Building $(PROJECT_NAME) services..."
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/management-service ./cmd/management-service
	go build $(LDFLAGS) -o bin/deploy-service ./cmd/deploy-service
	go build $(LDFLAGS) -o bin/webhook-service ./cmd/webhook-service
	@echo "Build completed!"

# 运行开发环境
run-dev:
	@echo "Starting development environment..."
	@make -j3 run-management-service run-deploy-service run-webhook-service

# 运行各个服务
run-management-service:
	@echo "Starting management service..."
	./bin/management-service

run-deploy-service:
	@echo "Starting deploy service..."
	./bin/deploy-service

run-webhook-service:
	@echo "Starting webhook service..."
	./bin/webhook-service

# 测试
test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 代码检查
lint:
	@echo "Running linters..."
	golangci-lint run

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...

# 清理
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker 构建
docker-build:
	@echo "Building Docker images..."
	docker build -t $(PROJECT_NAME)/management-service:$(VERSION) -f docker/management-service.Dockerfile .
	docker build -t $(PROJECT_NAME)/deploy-service:$(VERSION) -f docker/deploy-service.Dockerfile .
	docker build -t $(PROJECT_NAME)/webhook-service:$(VERSION) -f docker/webhook-service.Dockerfile .

# Docker 运行
docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

# 停止 Docker 服务
docker-stop:
	@echo "Stopping Docker services..."
	docker-compose down

# 数据库迁移
migrate-up:
	@echo "Running database migrations..."
	migrate -path migrations -database "postgres://localhost/boreas?sslmode=disable" up

migrate-down:
	@echo "Rolling back database migrations..."
	migrate -path migrations -database "postgres://localhost/boreas?sslmode=disable" down

# 生成 API 文档
docs:
	@echo "Generating API documentation..."
	swag init -g cmd/management-service/main.go -o docs/api

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# 安装开发工具
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 帮助信息
help:
	@echo "Available commands:"
	@echo "  build              - Build all services"
	@echo "  run-dev            - Run all services in development mode"
	@echo "  run-management-service - Run management service only"
	@echo "  run-deploy-service - Run deploy service only"
	@echo "  run-webhook-service - Run webhook service only"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  lint               - Run linters"
	@echo "  fmt                - Format code"
	@echo "  clean              - Clean build artifacts"
	@echo "  docker-build       - Build Docker images"
	@echo "  docker-run         - Run services with Docker Compose"
	@echo "  docker-stop        - Stop Docker services"
	@echo "  migrate-up         - Run database migrations"
	@echo "  migrate-down       - Rollback database migrations"
	@echo "  docs               - Generate API documentation"
	@echo "  deps               - Install dependencies"
	@echo "  install-tools      - Install development tools"
	@echo "  help               - Show this help message"
