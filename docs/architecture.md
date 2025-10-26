# Boreas 平台架构设计文档

## 1. 项目概述

Boreas 是一个基于 GitOps 的持续部署平台，支持 Kubernetes 和物理机多种部署环境，提供自动化的版本管理、应用管理和部署编排功能。

### 1.1 核心特性

- **GitOps 驱动**: 基于 Git 事件自动触发部署流程
- **多环境支持**: 支持 K8s、物理机等多种运行环境
- **智能编排**: 支持顺序部署、并行部署、依赖管理
- **可视化管理**: Web 界面展示部署流程和状态
- **灵活扩展**: 插件化的 Operator 架构

### 1.2 技术栈

- **后端**: Go 1.21+, Gin, GORM
- **数据库**: PostgreSQL 15+
- **缓存**: Redis 7+
- **前端**: React 18, TypeScript, Vite
- **容器**: Docker, Kubernetes

---

## 2. 系统架构

### 2.1 整体架构

```
┌────────────────────────────────────────────────────────────────┐
│                        外部系统                                  │
│  ┌─────────┐      ┌──────────┐      ┌──────────┐              │
│  │   Git   │      │ 用户浏览器 │      │  CI/CD   │              │
│  └────┬────┘      └─────┬────┘      └────┬─────┘              │
│       │ Webhook          │ HTTP           │ HTTP               │
└───────┼──────────────────┼────────────────┼────────────────────┘
        │                  │                │
┌───────▼──────────────────▼────────────────▼────────────────────┐
│                   Boreas 平台核心层                              │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                 Master Service (Port 8080/8090)          │  │
│  │  ┌──────────────┐  ┌───────────────┐  ┌──────────────┐  │  │
│  │  │   Webhook    │  │   Management  │  │   Workflow   │  │  │
│  │  │   Handler    │  │     API       │  │  Controller  │  │  │
│  │  └──────────────┘  └───────────────┘  └──────────────┘  │  │
│  └───────────────────────┬──────────────────────────────────┘  │
│                          │                                     │
│  ┌───────────────────────┼─────────────────────────────────┐  │
│  │              Shared Storage Layer                       │  │
│  │  ┌────────────────┐     ┌──────────────────┐           │  │
│  │  │  PostgreSQL    │     │      Redis       │           │  │
│  │  └────────────────┘     └──────────────────┘           │  │
│  └─────────────────────────────────────────────────────────┘  │
│                                                                │
│  ┌─────────────────────────────────────────────────────────┐  │
│  │               Web Dashboard (Port 3000)                 │  │
│  └─────────────────────────────────────────────────────────┘  │
└────────────────────────┬───────────────┬───────────────────────┘
                         │               │
        ┌────────────────┴────┬──────────┴─────────────┐
        │                     │                        │
┌───────▼─────────┐  ┌────────▼────────┐  ┌───────────▼──────────┐
│  Operator-K8s   │  │  Operator-Mock  │  │   Operator-PM        │
│  (Port 8081)    │  │  (Port 8083)    │  │   (Port 8082)        │
└───────┬─────────┘  └─────────────────┘  └───────────┬──────────┘
        │                                              │
┌───────▼─────────┐                      ┌─────────────▼──────────┐
│  Kubernetes     │                      │  Operator-PM-Agent     │
│  集群           │                      │  (Port 8081)           │
└─────────────────┘                      └─────────────┬──────────┘
                                                       │
                                         ┌─────────────▼──────────┐
                                         │  物理机环境             │
                                         └────────────────────────┘
```

### 2.2 组件说明

| 组件 | 端口 | 职责 | 关键技术 |
|------|------|------|---------|
| **Master Service** | 8080 | 核心业务逻辑、API 服务 | Go, Gin, GORM |
| **Webhook Service** | 8090 | 接收 Git 事件 | Go, Gin |
| **Workflow Controller** | - | 任务调度和编排 | Go, Goroutine |
| **Web Dashboard** | 3000 | 用户界面 | React, TypeScript |
| **Operator-K8s** | 8081 | K8s 部署执行 | Go, client-go |
| **Operator-PM** | 8082 | 物理机部署主控 | Go, Gin |
| **Operator-PM-Agent** | 8081 | 物理机节点代理 | Go, Docker API |
| **Operator-Mock** | 8083 | 模拟部署执行 | Go, Gin |

---

## 3. 核心组件详解

### 3.1 Master Service

**路径**: `cmd/master-service/`

#### 3.1.1 功能职责

1. **版本管理**: 创建、查询、删除版本
2. **应用管理**: CRUD 操作、配置管理、环境关联
3. **环境管理**: CRUD 操作、环境配置
4. **部署管理**: 创建、启动、取消、回滚部署
5. **触发器管理**: 自动触发规则配置

#### 3.1.2 模块结构

```
master-service/
├── main.go                    # 服务入口
├── webhook/                   # Webhook 服务
│   └── main.go
└── configs/
    └── master.yaml            # 配置文件

internal/services/master/
├── handler/                   # HTTP 处理器
│   ├── application_handler.go
│   ├── deployment_handler.go
│   ├── environment_handler.go
│   ├── task_handler.go
│   ├── trigger_handler.go
│   └── version_handler.go
├── service/                   # 业务逻辑
│   ├── application_service.go
│   ├── deployment_service.go
│   ├── deploy_executor_service.go
│   ├── environment_service.go
│   ├── task_service.go
│   ├── trigger.go
│   ├── version_service.go
│   ├── workflow_service.go    # 工作流控制器
│   ├── git.go                 # Git 集成
│   └── docker.go              # Docker 集成
└── repository/                # 数据访问
    └── postgres/
        ├── application_repository.go
        ├── deployment_repository.go
        ├── environment_repository.go
        ├── task_repository.go
        └── version_repository.go
```

#### 3.1.3 关键 API

| 接口 | 方法 | 功能 |
|------|------|------|
| `/api/v1/versions` | POST | 创建版本 |
| `/api/v1/versions` | GET | 获取版本列表 |
| `/api/v1/applications` | POST | 创建应用 |
| `/api/v1/applications` | GET | 获取应用列表 |
| `/api/v1/environments` | POST | 创建环境 |
| `/api/v1/deployments` | POST | 创建部署 |
| `/api/v1/deployments/:id/start` | POST | 启动部署 |
| `/api/v1/tasks` | GET | 获取任务列表 |
| `/webhook/github` | POST | GitHub Webhook |

### 3.2 Workflow Controller

**路径**: `internal/services/master/service/workflow_service.go`

#### 3.2.1 功能职责

1. **任务生成**: 从 Deployment 生成 Task 列表
2. **任务调度**: 根据依赖关系调度任务执行
3. **状态管理**: 跟踪任务和部署状态
4. **依赖处理**: 处理任务间的阻塞关系

#### 3.2.2 工作原理

```
┌──────────────────────────────────────────────────────────┐
│              Workflow Controller 工作流程                 │
└──────────────────────────────────────────────────────────┘

1. CreateTasksFromDeployment()
   ├─ 读取 Deployment.MustInOrder（部署顺序）
   ├─ 读取 Version.AppBuilds（应用构建信息）
   └─ 生成 Task 列表，设置依赖关系

2. Pending Task Scheduler (5秒轮询)
   ├─ 查询 Step=pending, Status=pending 的任务
   ├─ 检查依赖是否满足
   ├─ 调用 executeTask() 执行任务
   └─ 更新任务状态

3. Blocked Task Scheduler (10秒轮询)
   ├─ 查询 Step=blocked 的任务
   ├─ 检查依赖任务状态
   └─ 依赖满足后，将 Step 改为 pending

4. executeTask()
   ├─ 检查依赖任务是否完成
   ├─ 更新任务状态为 running
   ├─ 调用 DeployExecutor.Apply()
   ├─ 更新任务结果
   └─ 检查是否所有任务完成，更新 Deployment 状态
```

#### 3.2.3 任务状态机

```
Task.Step (执行步骤):
┌─────────┐
│ Pending │ ◄──┐
└────┬────┘    │
     │         │ 依赖满足
     ▼         │
┌─────────┐    │
│ Blocked │────┘
└─────────┘
     │
     ▼ 依赖满足
┌─────────┐
│ Running │
└────┬────┘
     │
     ▼
┌───────────┐
│ Completed │
└───────────┘

Task.Status (执行状态):
┌─────────┐
│ Pending │
└────┬────┘
     │
     ▼
┌─────────┐      ┌─────────┐
│ Running │ ───► │ Success │
└────┬────┘      └─────────┘
     │
     ▼
┌─────────┐
│ Failed  │
└─────────┘
```

### 3.3 Operator 组件

#### 3.3.1 Operator-K8s

**路径**: `cmd/operator-k8s/`

**功能**:
- 执行 Kubernetes 部署
- 管理 Deployment、Service、Ingress
- 查询 Pod 状态
- 获取容器日志

**关键 API**:
- `POST /api/v1/deploy/:id/execute` - 执行部署
- `GET /api/v1/deploy/:id/status` - 查询状态
- `GET /api/v1/deploy/:id/logs` - 获取日志

#### 3.3.2 Operator-PM

**路径**: `cmd/operator-pm/`

**功能**:
- 接收物理机部署请求
- 根据配置选择目标节点
- 分发任务到 Agent
- 汇总部署状态

**配置文件**:
- `operator-pm.yaml` - 主配置
- `app-to-nodes.yaml` - 应用→节点映射
- `node-to-ip.yaml` - 节点→IP 映射

**关键 API**:
- `POST /v1/apply` - 部署应用
- `GET /v1/status/:app` - 查询状态

#### 3.3.3 Operator-PM-Agent

**路径**: `cmd/operator-pm-agent/`

**功能**:
- 接收主控指令
- 管理本机应用（Docker/进程）
- 执行健康检查
- 上报状态

**运行模式**:
- Docker 模式
- 二进制模式
- 脚本模式

### 3.4 Web Dashboard

**路径**: `web/`

#### 3.4.1 页面结构

| 页面 | 路径 | 功能 |
|------|------|------|
| Dashboard | `/` | 系统概览 |
| Applications | `/applications` | 应用列表和管理 |
| ApplicationDetail | `/applications/:id` | 应用详情 |
| Versions | `/versions` | 版本列表 |
| Environments | `/environments` | 环境管理 |
| Deployments | `/deployments` | 部署列表 |
| DeploymentDetail | `/deployments/:id` | 部署详情和工作流 |

#### 3.4.2 核心组件

- **WorkflowViewer**: 工作流可视化（使用 React Flow）
- **WorkflowNode**: 任务节点渲染
- **Layout**: 页面布局

---

## 4. 数据模型

### 4.1 核心实体

#### Version (版本)
```go
type Version struct {
    ID          string         // 唯一标识
    Version     string         // 版本号（唯一）
    GitTag      string         // Git Tag
    GitCommit   string         // Git Commit
    Repository  string         // 仓库地址
    Status      string         // normal, revert
    CreatedBy   string
    CreatedAt   time.Time
    Description string
    AppBuilds   datatypes.JSON // 应用构建信息
}
```

#### Application (应用)
```go
type Application struct {
    ID           string
    Name         string         // 唯一
    Description  string
    Repository   string
    Type         string         // microservice, monolith
    Config       datatypes.JSON // 应用配置
    CreatedAt    time.Time
    UpdatedAt    time.Time
    Environments []Environment  // 关联环境
}
```

#### Environment (环境)
```go
type Environment struct {
    ID             string
    Name           string         // 唯一
    Type           string         // k8s, pm, mock
    Region         string
    Config         datatypes.JSON
    Status         string         // active, inactive
    OperatorAPIURL string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

#### Deployment (部署)
```go
type Deployment struct {
    ID             string
    VersionID      string
    MustInOrder    datatypes.JSON   // 部署顺序 []string
    EnvironmentID  string
    Status         DeploymentStatus // pending, running, completed
    CreatedBy      string
    CreatedAt      time.Time
    UpdatedAt      time.Time
    StartedAt      *time.Time
    CompletedAt    *time.Time
    ErrorMessage   string
    ManualApproval bool
    Strategy       datatypes.JSON
    Tasks          []Task
}
```

#### Task (任务)
```go
type Task struct {
    ID           string
    DeploymentID string
    AppID        string
    Name         string
    Type         string         // deploy
    Step         TaskStep       // pending, blocked, running, completed
    Status       TaskStatus     // pending, running, success, failed
    Dependencies datatypes.JSON // 依赖任务 ID 列表
    Payload      string
    Result       string
    CreatedAt    time.Time
    UpdatedAt    time.Time
    StartedAt    *time.Time
    CompletedAt  *time.Time
}
```

### 4.2 实体关系

```
Version (1) ──── (N) Deployment ──── (N) Task
                         │
                         │
Application (N) ──── (N) Environment
                         │
                         │ (1)
                    Deployment
```

---

## 5. 业务流程

### 5.1 完整部署流程

```
1. Git 事件触发
   Git Push/Tag → Webhook Service → Master Service

2. 创建版本
   Master Service → 解析事件 → 创建 Version 记录

3. 创建部署任务
   用户/自动触发 → 创建 Deployment → 关联 Version 和 Environment

4. 生成任务列表
   WorkflowController → CreateTasksFromDeployment() → 生成 Task 列表

5. 启动部署
   用户确认 → Deployment.Status = running

6. 任务调度
   Pending Task Scheduler → 检查依赖 → executeTask()

7. 执行部署
   DeployExecutor → 调用 Operator API → 执行具体部署

8. 状态更新
   Operator 返回结果 → 更新 Task.Status → 检查 Deployment 完成

9. 完成部署
   所有 Task 成功 → Deployment.Status = completed
```

### 5.2 任务依赖处理

```go
// 示例：顺序部署 3 个应用
MustInOrder: ["user-service", "order-service", "payment-service"]

生成的任务:
Task 1: user-service    (Dependencies: [])
Task 2: order-service   (Dependencies: [Task 1])
Task 3: payment-service (Dependencies: [Task 2])

执行流程:
1. Task 1: Step=pending → running → completed, Status=success
2. Blocked Task Scheduler 检测到 Task 2 依赖满足
3. Task 2: Step=blocked → pending → running → completed
4. 同理 Task 3 开始执行
```

---

## 6. 配置管理

### 6.1 Master Service 配置

**文件**: `cmd/master-service/configs/master.yaml`

```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  host: "localhost"
  port: 5432
  user: "boreas"
  password: "password"
  dbname: "boreas"

redis:
  host: "localhost"
  port: 6379

log:
  level: "info"
  format: "json"
```

### 6.2 Operator-PM 配置

**主配置**: `cmd/operator-pm/configs/operator-pm.yaml`
**应用映射**: `cmd/operator-pm/configs/app-to-nodes.yaml`
**节点映射**: `cmd/operator-pm/configs/node-to-ip.yaml`

---

## 7. 部署架构

### 7.1 Docker Compose 部署

```yaml
services:
  postgres:
    image: postgres:15
    ports: ["5432:5432"]
  
  redis:
    image: redis:7
    ports: ["6379:6379"]
  
  master-service:
    build: ./deployments/docker
    ports: ["8080:8080", "8090:8090"]
    depends_on: [postgres, redis]
  
  web:
    build: ./web
    ports: ["3000:3000"]
```

### 7.2 端口规划

| 服务 | 端口 | 协议 |
|------|------|------|
| Master Service (API) | 8080 | HTTP |
| Master Service (Webhook) | 8090 | HTTP |
| Operator-K8s | 8081 | HTTP |
| Operator-PM | 8082 | HTTP |
| Operator-PM-Agent | 8081 | HTTP |
| Operator-Mock | 8083 | HTTP |
| Web Dashboard | 3000 | HTTP |
| PostgreSQL | 5432 | TCP |
| Redis | 6379 | TCP |

---

## 8. 扩展性设计

### 8.1 Operator 插件化

新增 Operator 只需实现统一接口：

```go
type OperatorClient interface {
    Apply(ctx context.Context, task Task) (Status, error)
    GetStatus(ctx context.Context, taskID string) (Status, error)
    Cancel(ctx context.Context, taskID string) error
}
```

### 8.2 触发器扩展

支持多种触发方式：
- Git Webhook (已实现)
- 定时触发 (规划中)
- 手动触发 (已实现)
- API 触发 (已实现)

---

## 9. 监控和运维

### 9.1 健康检查

所有服务提供：
- `/v1/health` - 健康检查
- `/v1/ready` - 就绪检查

### 9.2 日志管理

- 格式: JSON
- 级别: DEBUG, INFO, WARN, ERROR
- 字段: timestamp, level, service, message, context

### 9.3 指标监控

建议集成 Prometheus：
- 部署成功率
- 部署平均耗时
- API 请求延迟
- 任务队列长度

---

## 10. 安全设计

### 10.1 认证授权

- API 认证: JWT Token
- Webhook 验证: HMAC 签名

### 10.2 数据安全

- 敏感配置加密存储
- 传输层 TLS 加密
- 数据库定期备份

---

## 附录

### A. 目录结构

```
boreas/
├── cmd/                          # 应用入口
│   ├── master-service/
│   ├── operator-k8s/
│   ├── operator-pm/
│   ├── operator-pm-agent/
│   └── operator-mock/
├── internal/
│   ├── pkg/                      # 共享库
│   │   ├── models/
│   │   ├── config/
│   │   ├── database/
│   │   ├── logger/
│   │   ├── middleware/
│   │   └── client/
│   ├── services/                 # 服务逻辑
│   │   └── master/
│   │       ├── handler/
│   │       ├── service/
│   │       └── repository/
│   └── interfaces/               # 接口定义
├── web/                          # 前端
├── migrations/                   # 数据库迁移
├── deployments/                  # 部署配置
└── docs/                         # 文档
```

### B. 相关文档

- [项目概述](./summary.md)
- [核心数据模型](./core-models.md)
- [部署工作流设计](./deployment-workflow-design.md)
- [Operator-PM 设计](./operator-pm-design.md)
- [Master Service API](./api/master-service.md)

---

**文档版本**: v1.0  
**最后更新**: 2025-10-26  
**维护者**: Boreas 开发团队

