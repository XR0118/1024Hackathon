# Boreas - 基于 GitOps 的多服务持续部署平台

Boreas 是一个基于 GitOps 的持续部署平台，支持 Kubernetes 和物理机等多种部署环境，提供自动化的版本管理、应用管理和部署编排功能。

## 📚 文档

- [项目概述](docs/summary.md) - 核心概念、架构概述、业务流程图
- [架构设计](docs/architecture.md) - 系统整体架构、核心组件、数据模型、业务流程
- [部署工作流设计](docs/deployment-workflow-design.md) - 部署任务流程、工作流可视化、前后端交互
- [Operator-PM 设计](docs/operator-pm-design.md) - 物理机部署组件架构、API 设计、部署指南
- [Master Service API](docs/api/master-service.md) - RESTful API 接口文档

## ✨ 核心特性

- **GitOps 驱动**: Git 事件自动触发部署流程，保持代码与部署状态一致
- **多环境支持**: 支持 Kubernetes、物理机等多种运行环境
- **智能编排**: 支持顺序部署、依赖管理、自动回滚
- **可视化管理**: Web 界面展示部署流程和状态
- **灵活扩展**: 插件化的 Operator 架构

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
   
   # 启动PM Operator
   make run-operator-pm
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

## 系统架构

```
Git 仓库 ──Webhook──> Master Service ──部署指令──> Operator-K8s/PM/Mock
                         │                              │
                         │                              │
                    PostgreSQL                    目标环境 (K8s/PM)
                      Redis
                         │
                         │
                   Web Dashboard (React)
```

### 核心组件

| 组件 | 端口 | 职责 |
|------|------|------|
| **Master Service** | 8080 | 核心业务逻辑、任务调度 |
| **Web Dashboard** | 3000 | 用户界面、工作流可视化 |
| **Operator-K8s** | 8081 | Kubernetes 部署执行 |
| **Operator-PM** | 8082 | 物理机部署主控 |
| **Operator-PM-Agent** | 8081 | 物理机节点代理 |
| **Operator-Mock** | 8083 | 模拟部署（测试） |

> ⚠️ **注意**：Operator 服务的实际端口需要与 Master Service 配置文件中的 `operator.k8s_operator_url` 和 `operator.pm_operator_url` 保持一致。

> 📖 详细架构设计请参阅 [架构设计文档](docs/architecture.md)

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

各服务配置文件位置：
- Master Service: `cmd/master-service/configs/master.yaml`
- Operator-K8s: `cmd/operator-k8s/configs/operator-k8s.yaml`
- Operator-PM: `cmd/operator-pm/configs/operator-pm.yaml`
- Operator-PM-Agent: `cmd/operator-pm-agent/configs/agent.yaml`

## 核心概念

- **Version (版本)**: 对应 Git Tag/Commit，包含应用构建信息
- **Application (应用)**: 部署的最小单元，可关联多个环境
- **Environment (环境)**: 部署目标 (K8s/物理机)
- **Deployment (部署)**: 将版本部署到环境的任务
- **Task (任务)**: 部署的执行单元，支持依赖关系

> 详细说明请参阅 [项目概述](docs/summary.md) 和 [架构设计](docs/architecture.md)

## 开发指南

### 项目结构

```
boreas/
├── cmd/                    # 应用入口
│   ├── master-service/    # 核心服务 + Webhook
│   ├── operator-k8s/      # K8s Operator
│   ├── operator-pm/       # PM Operator 主控
│   ├── operator-pm-agent/ # PM Agent
│   └── operator-mock/     # Mock Operator
├── internal/
│   ├── pkg/               # 共享库 (models, database, logger, client)
│   ├── services/          # 服务逻辑 (handler, service, repository)
│   └── interfaces/        # 接口定义
├── web/                   # React 前端
├── migrations/            # 数据库迁移
└── docs/                  # 文档
```

### 开发流程

1. 定义数据模型 (`internal/pkg/models/`)
2. 定义接口 (`internal/interfaces/`)
3. 实现数据访问 (`service/repository/`)
4. 实现业务逻辑 (`service/service/`)
5. 实现 HTTP 处理 (`service/handler/`)
6. 注册路由 (`cmd/*/main.go`)

### 代码检查

```bash
# 运行测试
make test-all

# 代码格式化
make fmt-all

# 运行 linter
make lint-all
```

## 部署指南

### 生产环境部署

1. **准备环境**
   - 安装PostgreSQL和Redis
   - 配置Kubernetes集群（如果使用K8s部署）
   - 配置物理机环境（如果使用Baremetal部署）

2. **配置应用**
   - 修改各服务配置文件（位于 `cmd/*/configs/`）或设置环境变量
   - 配置 GitHub Webhook 密钥
   - 配置 Kubernetes 认证信息（kubeconfig）
   - 配置物理机节点映射（Operator-PM）

3. **部署服务**
   ```bash
   # 使用Docker Compose
   make docker-run-all
   
   # 或使用Kubernetes
   kubectl apply -f deployments/k8s/
   ```

### 监控和日志

- 健康检查端点: `/v1/health`
- 就绪检查端点: `/v1/ready`
- 日志格式: JSON
- 日志级别: 可配置 (debug, info, warn, error)

## 部署流程

```
1. Git 事件触发 → Webhook → 创建 Version
2. 创建 Deployment → 生成 Task 列表（按 MustInOrder 顺序）
3. Workflow Controller 调度 Task 执行
   - Pending Task Scheduler: 执行无依赖的任务
   - Blocked Task Scheduler: 检查依赖，解除阻塞
4. 调用 Operator 执行具体部署
5. 更新 Task 和 Deployment 状态
```

**任务依赖处理示例**:
```yaml
MustInOrder: ["user-service", "order-service", "payment-service"]

生成任务:
- Task 1: user-service    (无依赖)
- Task 2: order-service   (依赖 Task 1)
- Task 3: payment-service (依赖 Task 2)
```

> 详细流程请参阅 [部署工作流设计](docs/deployment-workflow-design.md)

## 许可证

MIT License

## 联系

- 问题反馈: GitHub Issues