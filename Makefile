# Boreas 多服务架构 Makefile

.PHONY: help build run test clean docker-build docker-run docker-stop migrate-up migrate-down deps install-tools

# 默认目标
help:
	@echo "Boreas 多服务架构构建工具"
	@echo ""
	@echo "可用命令:"
	@echo "  deps                    - 安装所有依赖"
	@echo "  install-tools           - 安装开发工具"
	@echo "  build-all               - 构建所有服务"
	@echo "  build-master            - 构建Master服务"
	@echo "  build-operator-k8s      - 构建K8s Operator服务"
	@echo "  build-operator-pm        - 构建PM Operator服务"
	@echo "  build-web               - 构建Web管理界面"
	@echo "  run-all                 - 运行所有服务"
	@echo "  run-master              - 运行Master服务"
	@echo "  run-operator-k8s        - 运行K8s Operator服务"
	@echo "  run-operator-pm         - 运行PM Operator服务"
	@echo "  run-web                 - 运行Web管理界面"
	@echo "  test-all                - 运行所有测试"
	@echo "  test-master             - 运行Master服务测试"
	@echo "  test-operator-k8s       - 运行K8s Operator测试"
	@echo "  test-operator-pm        - 运行PM Operator测试"
	@echo "  fmt-all                 - 格式化所有代码"
	@echo "  lint-all                - 运行所有linter"
	@echo "  clean-all               - 清理所有构建文件"
	@echo "  docker-build-all        - 构建所有Docker镜像"
	@echo "  docker-run-all          - 运行所有Docker容器"
	@echo "  docker-stop-all         - 停止所有Docker容器"
	@echo "  migrate-up              - 运行数据库迁移"
	@echo "  migrate-down            - 回滚数据库迁移"

# 依赖管理
deps:
	@echo "安装Go依赖..."
	go mod tidy
	@echo "安装Node.js依赖..."
	cd web && npm install

install-tools:
	@echo "安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# 构建
build-all: build-master build-operator-k8s build-operator-pm build-web

build-master:
	@echo "构建Master服务..."
	go build -o bin/master-service cmd/master-service/main.go

build-operator-k8s:
	@echo "构建K8s Operator服务..."
	go build -o bin/operator-k8s cmd/operator-k8s/main.go

build-operator-pm:
	@echo "构建PM Operator服务..."
	go build -o bin/operator-pm cmd/operator-pm/main.go

build-web:
	@echo "构建Web管理界面..."
	cd web && npm run build

# 运行
run-all: run-master run-operator-k8s run-operator-pm run-web

run-master:
	@echo "运行Master服务..."
	go run cmd/master-service/main.go

run-operator-k8s:
	@echo "运行K8s Operator服务..."
	go run cmd/operator-k8s/main.go

run-operator-pm:
	@echo "运行PM Operator服务..."
	go run cmd/operator-pm/main.go

run-web:
	@echo "运行Web管理界面..."
	cd web && npm run dev

# 测试
test-all: test-master test-operator-k8s test-operator-pm

test-master:
	@echo "运行Master服务测试..."
	go test ./internal/services/master/...

test-operator-k8s:
	@echo "运行K8s Operator测试..."
	go test ./internal/services/operator-k8s/...

test-operator-pm:
	@echo "运行PM Operator测试..."
	go test ./internal/services/operator-pm/...

# 代码格式化
fmt-all:
	@echo "格式化所有代码..."
	go fmt ./...

# Lint
lint-all:
	@echo "运行所有linter..."
	golangci-lint run

# 清理
clean-all:
	@echo "清理所有构建文件..."
	rm -rf bin/
	cd web && rm -rf dist/
	cd web && rm -rf node_modules/

# Docker
docker-build-all:
	@echo "构建所有Docker镜像..."
	docker build -f deployments/docker/master-service.Dockerfile -t boreas/master-service:latest .
	docker build -f deployments/docker/operator-k8s.Dockerfile -t boreas/operator-k8s:latest .
	docker build -f deployments/docker/operator-pm.Dockerfile -t boreas/operator-pm:latest .
	docker build -f deployments/docker/web-management.Dockerfile -t boreas/web-management:latest .

docker-run-all:
	@echo "运行所有Docker容器..."
	docker-compose up -d

docker-stop-all:
	@echo "停止所有Docker容器..."
	docker-compose down

# 数据库迁移
migrate-up:
	@echo "运行数据库迁移..."
	go run cmd/master-service/migrate.go up

migrate-down:
	@echo "回滚数据库迁移..."
	go run cmd/master-service/migrate.go down

# 开发环境
dev-setup: deps install-tools
	@echo "设置开发环境..."
	docker-compose up -d postgres redis
	@echo "等待数据库启动..."
	sleep 10
	make migrate-up
	@echo "开发环境设置完成！"

# 生产环境
prod-deploy: docker-build-all docker-run-all
	@echo "生产环境部署完成！"