# Boreas - 基于 GitOps 的多服务持续部署平台

Boreas 是一个基于 GitOps 的持续部署平台，采用单体仓库 + 共享库模式，支持 Kubernetes 和物理机部署，提供完整的版本管理、应用管理、环境管理和部署管理功能。

## 架构概述

Boreas 采用单体仓库 + 共享库模式，包含以下核心服务：

### 1. Master Service (核心服务)
作为整个系统的中枢模块，负责：
- 上游承接时间通知，包含Git Tag等，即创建版本
- 内部实现规划、编排部署任务
- 运行编排好的部署任务，让下游具体操作
- 同时监控跟踪部署状态

**端口**: 8080  
**入口**: `cmd/master-service/main.go`

### 2. Operator Services (部署执行器)
为master下游，由于K8S与裸金属环境的不同，拆分成不同服务：

#### Operator-K8s
- 执行Kubernetes环境的部署操作
- 维持Kubernetes部署状态
- 提供Kubernetes部署状态，供上游查询

**端口**: 8081  
**入口**: `cmd/operator-k8s/main.go`

#### Operator-Baremetal
- 执行物理机环境的部署操作
- 维持物理机部署状态
- 提供物理机部署状态，供上游查询

**端口**: 8082  
**入口**: `cmd/operator-baremetal/main.go`

### 3. Web Management (Web管理界面)
为master服务搭配的web管理界面，包含功能：
- 状态查看
- 半自动编排/编辑等
- 人工复核确认

**端口**: 3000  
**目录**: `web/`

## 项目结构

```
boreas/
├── cmd/                          # 应用入口
│   ├── master-service/          # 核心服务入口
│   │   ├── main.go             # 主服务入口
│   │   └── webhook/            # Webhook服务入口
│   ├── operator-k8s/            # K8s Operator入口
│   │   └── main.go
│   └── operator-baremetal/      # Baremetal Operator入口
│       └── main.go
├── internal/                     # 内部包
│   ├── pkg/                     # 共享包
│   │   ├── config/              # 配置管理
│   │   ├── database/            # 数据库连接
│   │   ├── logger/              # 日志管理
│   │   ├── middleware/          # 中间件
│   │   ├── models/              # 数据模型
│   │   └── utils/               # 工具函数
│   ├── services/                # 各服务特有逻辑
│   │   ├── master/              # Master服务逻辑
│   │   │   ├── handler/         # HTTP处理器
│   │   │   ├── service/         # 业务逻辑
│   │   │   └── repository/      # 数据访问层
│   │   ├── operator-k8s/        # K8s Operator逻辑
│   │   │   ├── handler/
│   │   │   ├── service/
│   │   │   └── repository/
│   │   └── operator-baremetal/  # Baremetal Operator逻辑
│   │       ├── handler/
│   │       ├── service/
│   │       └── repository/
│   └── interfaces/              # 接口定义
├── web/                         # 前端管理界面
├── api/                         # API定义
│   ├── proto/                   # gRPC定义
│   └── openapi/                 # REST API定义
├── configs/                     # 配置文件
├── deployments/                 # 部署配置
│   └── docker/                  # Docker配置
├── docs/                        # 文档
├── migrations/                  # 数据库迁移
├── scripts/                     # 脚本
├── docker-compose.yml           # Docker Compose配置
├── go.mod                       # Go模块定义
└── Makefile                     # 构建脚本
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose
- Kubernetes 集群 (如果使用K8s部署)

### 本地开发

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd boreas
   ```

2. **安装依赖**
   ```bash
   make deps
   ```

3. **启动基础设施**
   ```bash
   docker-compose up -d postgres redis
   ```

4. **运行数据库迁移**
   ```bash
   make migrate-up
   ```

5. **启动Master服务**
   ```bash
   make run-master
   ```

6. **启动Operator服务**
   ```bash
   # 启动K8s Operator
   make run-operator-k8s
   
   # 启动Baremetal Operator
   make run-operator-baremetal
   ```

7. **启动Web管理界面**
   ```bash
   make run-web
   ```

### Docker部署

1. **构建所有服务**
   ```bash
   make docker-build-all
   ```

2. **启动所有服务**
   ```bash
   make docker-run-all
   ```

3. **停止所有服务**
   ```bash
   make docker-stop-all
   ```

## 服务间通信

### Master Service API
- **端口**: 8080
- **功能**: 版本管理、应用管理、环境管理、部署编排
- **健康检查**: `GET /health`
- **就绪检查**: `GET /ready`

### Operator-K8s API
- **端口**: 8081
- **功能**: Kubernetes部署执行、状态查询
- **健康检查**: `GET /health`
- **就绪检查**: `GET /ready`
- **部署执行**: `POST /api/v1/deploy/{id}/execute`
- **状态查询**: `GET /api/v1/deploy/{id}/status`
- **日志获取**: `GET /api/v1/deploy/{id}/logs`
- **取消部署**: `POST /api/v1/deploy/{id}/cancel`

### Operator-Baremetal API
- **端口**: 8082
- **功能**: 物理机部署执行、状态查询
- **健康检查**: `GET /health`
- **就绪检查**: `GET /ready`
- **部署执行**: `POST /api/v1/deploy/{id}/execute`
- **状态查询**: `GET /api/v1/deploy/{id}/status`
- **日志获取**: `GET /api/v1/deploy/{id}/logs`
- **取消部署**: `POST /api/v1/deploy/{id}/cancel`

### Web Management
- **端口**: 3000
- **功能**: 管理界面、状态查看、人工复核

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
| `REDIS_HOST` | Redis主机 | `localhost` |
| `REDIS_PORT` | Redis端口 | `6379` |
| `LOG_LEVEL` | 日志级别 | `info` |
| `GITHUB_WEBHOOK_SECRET` | GitHub Webhook密钥 | - |

### 配置文件

配置文件位于 `configs/config.yaml`，支持 YAML 格式配置。

## 开发指南

### 代码结构

- `internal/pkg/` - 共享包，所有服务都可以使用
- `internal/services/` - 各服务特有的业务逻辑
- `internal/interfaces/` - 接口定义
- `cmd/` - 各服务的入口点

### 添加新功能

1. 在 `internal/pkg/models/` 中定义数据模型
2. 在 `internal/interfaces/` 中定义接口
3. 在对应服务的 `repository/` 中实现数据访问
4. 在对应服务的 `service/` 中实现业务逻辑
5. 在对应服务的 `handler/` 中实现HTTP处理
6. 在 `cmd/` 中注册路由

### 测试

```bash
# 运行所有测试
make test-all

# 运行特定服务测试
make test-master
make test-operator-k8s
make test-operator-baremetal
```

### 代码检查

```bash
# 格式化所有代码
make fmt-all

# 运行所有linter
make lint-all
```

## 部署指南

### 生产环境部署

1. **准备环境**
   - 安装PostgreSQL和Redis
   - 配置Kubernetes集群（如果使用K8s部署）
   - 配置物理机环境（如果使用Baremetal部署）

2. **配置应用**
   - 修改 `configs/config.yaml` 或设置环境变量
   - 配置GitHub Webhook密钥
   - 配置Kubernetes认证信息

3. **部署服务**
   ```bash
   # 使用Docker Compose
   make docker-run-all
   
   # 或使用Kubernetes
   kubectl apply -f deployments/k8s/
   ```

### 监控和日志

- 健康检查端点: `/health`
- 就绪检查端点: `/ready`
- 日志格式: JSON
- 日志级别: 可配置

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License

## 联系方式

- 项目地址: [GitHub Repository]
- 问题反馈: [GitHub Issues]
- 文档: [Project Documentation]