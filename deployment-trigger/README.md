# Deployment Trigger 模块

## 概述

Deployment Trigger 模块负责接收 GitHub webhook 事件,当新的 tag 被创建时自动触发构建流程,并将构建结果提交到发布系统。

## 功能

1. 接收 GitHub webhook (tag 创建事件)
2. 拉取新 tag 代码并与上一个 tag 进行对比
3. 识别变更的应用列表
4. 根据应用配置执行 Docker 构建
5. 将构建的镜像推送到 Docker Hub
6. 调用发布系统 API 创建新版本

## 架构

```
GitHub Webhook → WebhookHandler → VersionService → GitService/DockerService/ManagementClient
```

## 使用方法

### 环境变量

- `GITHUB_WEBHOOK_SECRET`: GitHub webhook 的密钥
- `WORK_DIR`: 工作目录,默认 `/tmp/deployment-trigger`
- `DOCKER_REGISTRY`: Docker 镜像仓库地址
- `MANAGEMENT_API`: 管理系统 API 地址,默认 `http://localhost:8080/api/v1`

### 启动服务

```bash
go build -o deployment-trigger ./deployment-trigger/main.go
./deployment-trigger
```

服务将在端口 8081 上启动。

### GitHub Webhook 配置

在 GitHub 仓库设置中配置 Webhook:

```
Payload URL: http://your-server:8081/webhook/github
Content type: application/json
Secret: <your-webhook-secret>
Events: Push events (选择 Just the push event)
```

## 业务流程

1. **接收 webhook**: 接收 GitHub 的 tag 创建事件
2. **验证签名**: 验证 webhook 签名确保请求来自 GitHub
3. **过滤事件**: 只处理 tag 创建事件(refs/tags/*),其他返回 200
4. **克隆代码**: 克隆或更新仓库,切换到新 tag
5. **对比差异**: 与上一个 tag 对比,获取变更的应用列表
6. **构建镜像**: 根据每个应用的构建配置执行 Docker 构建
7. **推送镜像**: 将构建好的镜像推送到 Docker Hub
8. **创建版本**: 调用管理系统 API 创建新版本记录

## 目录结构

```
deployment-trigger/
├── main.go                          # 入口文件
├── internal/
│   ├── handler/
│   │   └── webhook.go              # Webhook 处理器
│   └── service/
│       ├── version.go              # 版本服务(主逻辑)
│       ├── git.go                  # Git 操作
│       ├── docker.go               # Docker 操作
│       ├── management.go           # 管理系统客户端
│       └── version_test.go         # 测试
├── Dockerfile
└── README.md
```

注意: 本模块使用项目根目录的 `go.mod` 文件,不再维护独立的 Go 模块。

## 依赖

- `github.com/go-git/go-git/v5`: Git 操作
- `github.com/docker/docker`: Docker 操作
- `github.com/google/go-github/v56`: GitHub API (可选)

## 测试

```bash
go test ./...
```

## 构建

```bash
go build -o deployment-trigger ./deployment-trigger/main.go
```

## Docker 部署

```bash
docker build -t deployment-trigger -f deployment-trigger/Dockerfile .
docker run -d \
  -p 8081:8081 \
  -e GITHUB_WEBHOOK_SECRET=your-secret \
  -e DOCKER_REGISTRY=registry.example.com \
  -e MANAGEMENT_API=http://management-api:8080/api/v1 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  deployment-trigger
```

## 注意事项

1. 需要 Docker 环境支持
2. 需要配置 Docker 镜像仓库访问权限
3. 应用目录需要包含 Dockerfile
4. 管理系统 API 需要提前配置好应用信息
