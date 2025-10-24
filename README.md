# Boreas - 基于 GitOps 的持续部署平台

Boreas 是一个基于 GitOps 的持续部署平台，支持 Kubernetes 和物理机部署，提供完整的版本管理、应用管理、环境管理和部署管理功能。

## 项目结构

```
boreas/
├── cmd/                    # 应用程序入口
│   ├── management-service/ # 管理服务
│   ├── deploy-service/    # 部署服务
│   └── webhook-service/   # Webhook 服务
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── database/         # 数据库连接
│   ├── handler/          # HTTP 处理器
│   ├── interfaces/       # 接口定义
│   ├── logger/           # 日志管理
│   ├── middleware/       # 中间件
│   ├── models/           # 数据模型
│   ├── repository/       # 数据访问层
│   └── service/          # 业务逻辑层
├── configs/              # 配置文件
├── docker/               # Docker 文件
├── migrations/           # 数据库迁移
├── nginx/                # Nginx 配置
├── docs/                 # 文档
└── web/                  # 前端应用
```

## 功能特性

### 核心功能

- **版本管理**: 基于 Git Tag 的版本管理，支持自动版本创建
- **应用管理**: 支持微服务和单体应用的管理
- **环境管理**: 支持 Kubernetes 和物理机环境
- **部署管理**: 完整的部署生命周期管理
- **工作流管理**: 基于任务的工作流编排
- **Webhook 集成**: 支持 GitHub Webhook 自动触发

### 部署支持

- **Kubernetes 部署**: 支持滚动更新、蓝绿部署等策略
- **物理机部署**: 支持 SSH 部署和系统服务管理
- **健康检查**: 多种健康检查方式
- **回滚支持**: 快速回滚到历史版本

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose

### 本地开发

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd boreas
   ```

2. **安装依赖**
   ```bash
   make deps
   make install-tools
   ```

3. **启动数据库**
   ```bash
   docker-compose up -d postgres redis
   ```

4. **运行数据库迁移**
   ```bash
   make migrate-up
   ```

5. **启动服务**
   ```bash
   # 启动管理服务
   make run-management-service
   
   # 或启动所有服务
   make run-dev
   ```

### Docker 部署

1. **构建镜像**
   ```bash
   make docker-build
   ```

2. **启动服务**
   ```bash
   make docker-run
   ```

3. **停止服务**
   ```bash
   make docker-stop
   ```

## API 文档

### 管理服务 API (端口 8080)

#### 版本管理
- `POST /api/v1/versions` - 创建版本
- `GET /api/v1/versions` - 获取版本列表
- `GET /api/v1/versions/{id}` - 获取版本详情
- `DELETE /api/v1/versions/{id}` - 删除版本

#### 应用管理
- `POST /api/v1/applications` - 创建应用
- `GET /api/v1/applications` - 获取应用列表
- `GET /api/v1/applications/{id}` - 获取应用详情
- `PUT /api/v1/applications/{id}` - 更新应用
- `DELETE /api/v1/applications/{id}` - 删除应用

#### 环境管理
- `POST /api/v1/environments` - 创建环境
- `GET /api/v1/environments` - 获取环境列表
- `GET /api/v1/environments/{id}` - 获取环境详情
- `PUT /api/v1/environments/{id}` - 更新环境
- `DELETE /api/v1/environments/{id}` - 删除环境

#### 部署管理
- `POST /api/v1/deployments` - 创建部署
- `GET /api/v1/deployments` - 获取部署列表
- `GET /api/v1/deployments/{id}` - 获取部署详情
- `POST /api/v1/deployments/{id}/cancel` - 取消部署
- `POST /api/v1/deployments/{id}/rollback` - 回滚部署

#### 任务管理
- `GET /api/v1/tasks` - 获取任务列表
- `GET /api/v1/tasks/{id}` - 获取任务详情
- `POST /api/v1/tasks/{id}/retry` - 重试任务

### Webhook 服务 API (端口 8082)

#### GitHub Webhook
- `POST /api/v1/webhooks/github` - 接收 GitHub Webhook

### 部署服务 API (端口 8081)

#### 内部 API
- `GET /internal/v1/deploy/info/{deployment_id}` - 获取部署信息
- `GET /internal/v1/deploy/health/{deployment_id}` - 健康检查
- `GET /internal/v1/deploy/logs/{deployment_id}` - 获取日志

## 配置说明

### 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `SERVER_HOST` | 服务器地址 | `0.0.0.0` |
| `SERVER_PORT` | 服务器端口 | `8080` |
| `DB_HOST` | 数据库主机 | `localhost` |
| `DB_PORT` | 数据库端口 | `5432` |
| `DB_USER` | 数据库用户 | `boreas` |
| `DB_PASSWORD` | 数据库密码 | `boreas123` |
| `DB_NAME` | 数据库名称 | `boreas` |
| `REDIS_HOST` | Redis 主机 | `localhost` |
| `REDIS_PORT` | Redis 端口 | `6379` |
| `LOG_LEVEL` | 日志级别 | `info` |
| `GITHUB_WEBHOOK_SECRET` | GitHub Webhook 密钥 | - |

### 配置文件

配置文件位于 `configs/config.yaml`，支持 YAML 格式配置。

## 开发指南

### 代码结构

- `internal/models/` - 数据模型和类型定义
- `internal/interfaces/` - 接口定义
- `internal/repository/` - 数据访问层实现
- `internal/service/` - 业务逻辑层实现
- `internal/handler/` - HTTP 处理器
- `internal/middleware/` - 中间件

### 添加新功能

1. 在 `internal/models/` 中定义数据模型
2. 在 `internal/interfaces/` 中定义接口
3. 在 `internal/repository/` 中实现数据访问
4. 在 `internal/service/` 中实现业务逻辑
5. 在 `internal/handler/` 中实现 HTTP 处理
6. 在 `cmd/` 中注册路由

### 测试

```bash
# 运行测试
make test

# 运行测试并生成覆盖率报告
make test-coverage
```

### 代码检查

```bash
# 格式化代码
make fmt

# 运行 linter
make lint
```

## 部署指南

### 生产环境部署

1. **准备环境**
   - 安装 PostgreSQL 和 Redis
   - 配置 Kubernetes 集群（如果使用 K8s 部署）

2. **配置应用**
   - 修改 `configs/config.yaml` 或设置环境变量
   - 配置 GitHub Webhook 密钥

3. **部署服务**
   ```bash
   # 使用 Docker Compose
   docker-compose up -d
   
   # 或使用 Kubernetes
   kubectl apply -f k8s/
   ```

### 监控和日志

- 健康检查端点: `/health`
- 就绪检查端点: `/ready`
- 日志格式: JSON
- 日志级别: 可配置

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

- 项目地址: [GitHub Repository]
- 问题反馈: [GitHub Issues]
- 文档: [Project Documentation]