# 部署任务流程设计文档

## 📋 文档概述

本文档详细描述了 Boreas 持续部署平台的部署任务流程设计，包括核心概念、数据模型、执行流程和前端展示。

**最后更新**: 2025-10-26

---

## 🎯 核心概念

### 1. 部署任务 (Deployment)

部署任务是平台的核心执行单元，每个版本(Version)在特定环境(Environment)中的部署会创建一个部署任务。

**关键特性**:
- 一个版本可以对应多个部署任务（不同环境）
- 包含多个应用(Application)的部署操作
- 支持部署编排（顺序、批次、灰度）
- 支持人工确认和自动回滚

### 2. 任务 (Task)

Task 是 Deployment 的子单元，表示对单个应用的具体操作。

**任务类型**:
- `build`: 构建镜像/制品
- `test`: 执行测试
- `deploy`: 部署应用
- `health_check`: 健康检查

### 3. 步骤/任务 (Task)

Task 是部署工作流的基本执行单元，前后端统一使用 Task 概念。前端将 Task 可视化为工作流节点。

---

## 📊 数据模型

### 后端模型 (Go)

#### Deployment 模型

```go
type Deployment struct {
    ID            string           // 部署任务唯一标识
    VersionID     string           // 关联的版本ID
    MustInOrder   datatypes.JSON   // 应用部署顺序 []string
    EnvironmentID string           // 目标环境ID
    Status        DeploymentStatus // 部署状态
    CreatedBy     string           // 创建人
    CreatedAt     time.Time        // 创建时间
    UpdatedAt     time.Time        // 更新时间
    StartedAt     *time.Time       // 开始时间
    CompletedAt   *time.Time       // 完成时间
    ErrorMessage  string           // 错误信息
    
    ManualApproval bool           // 是否需要人工审批
    Strategy       datatypes.JSON // 部署策略 []DeploySteps
    
    // 关联
    Version     Version
    Environment Environment
    Tasks       []Task
}

type DeploymentStatus string
const (
    DeploymentStatusPending    = "pending"      // 等待执行
    DeploymentStatusRunning    = "running"      // 执行中
    DeploymentStatusSuccess    = "success"      // 成功
    DeploymentStatusFailed     = "failed"       // 失败
    DeploymentStatusRolledBack = "rolled_back"  // 已回滚
    DeploymentStatusCancelled  = "cancelled"    // 已取消
)
```

#### Task 模型

```go
type Task struct {
    ID           string     // 任务唯一标识
    DeploymentID string     // 所属部署任务ID
    AppID        string     // 关联的应用ID
    Type         string     // 任务类型
    Status       TaskStatus // 任务状态
    BlockBy      string     // 阻塞依赖
    Payload      string     // 任务负载数据
    Result       string     // 执行结果
    CreatedAt    time.Time
    UpdatedAt    time.Time
    StartedAt    *time.Time
    CompletedAt  *time.Time
    
    // 关联
    Deployment  Deployment
    Application Application
}

type TaskStatus string
const (
    TaskStatusPending    = "pending"      // 等待执行
    TaskStatusRunning    = "running"      // 执行中
    TaskStatusSuccess    = "success"      // 成功
    TaskStatusFailed     = "failed"       // 失败
    TaskStatusBlocked    = "blocked"      // 被阻塞
    TaskStatusCancelled  = "cancelled"    // 已取消
    TaskStatusRolledBack = "rolled_back"  // 已回滚
)
```

#### DeploySteps 策略

```go
type DeploySteps struct {
    BatchSize            int     // 批次大小
    BatchInterval        int     // 批次间隔（秒）
    CanaryRatio          float64 // 金丝雀比例
    AutoRollback         bool    // 自动回滚
    ManualApprovalStatus *bool   // 人工审批状态
}
```

### 前端模型 (TypeScript)

#### Deployment 接口

```typescript
export interface Deployment {
  id: string
  versionId: string
  version: string
  applicationIds: string[]
  applications: string[]
  environmentIds: string[]
  environments: string[]
  status: 'pending' | 'running' | 'success' | 'failed' | 'waiting_confirm'
  progress: number
  createdAt: string
  updatedAt: string
  duration?: number
  requireConfirm: boolean
  grayscaleEnabled: boolean
  grayscaleRatio?: number
}

export interface DeploymentDetail extends Deployment {
  tasks: Task[]
  logs: DeploymentLog[]
}
```

#### Task 接口

```typescript
export interface Task {
  id: string
  deploymentId?: string
  appId?: string
  name: string
  type: 'build' | 'test' | 'deploy' | 'health_check' | 'prepare' | 'custom'
  status: 'pending' | 'running' | 'success' | 'failed' | 'blocked' | 'cancelled'
  blockBy?: string
  duration?: number
  startedAt?: string
  completedAt?: string
  logs?: string[]
}
```

**注意**: 前端 Task 接口与后端保持一致，便于数据同步和管理。

---

## 🔄 执行流程

### 1. 部署任务创建流程

```mermaid
sequenceDiagram
    participant User as 用户/Git
    participant Master as Master 服务
    participant DB as 数据库
    participant Operator as Operator

    User->>Master: 创建部署请求
    Master->>Master: 验证版本和环境
    Master->>Master: 解析应用列表
    Master->>Master: 生成部署策略
    Master->>DB: 创建 Deployment 记录
    Master->>DB: 创建 Task 记录
    Master-->>User: 返回部署任务ID
    
    opt 自动开始
        Master->>Master: 启动部署执行器
        Master->>Operator: 发送部署指令
    end
```

### 2. 部署任务执行流程

```mermaid
stateDiagram-v2
    [*] --> Pending: 创建部署任务
    Pending --> Running: 开始执行
    
    Running --> WaitingConfirm: 需要人工确认
    WaitingConfirm --> Running: 确认继续
    WaitingConfirm --> RolledBack: 选择回滚
    
    Running --> Success: 所有任务成功
    Running --> Failed: 任务失败
    
    Failed --> RolledBack: 自动/手动回滚
    
    Success --> [*]
    Failed --> [*]
    RolledBack --> [*]
    Running --> Cancelled: 取消部署
    Cancelled --> [*]
```

### 3. Task 执行顺序

部署任务中的 Task 按照以下规则执行:

1. **顺序约束**: `MustInOrder` 字段定义应用部署顺序
2. **类型顺序**: 同一应用内，Task 按类型顺序执行
   - build → test → deploy → health_check
3. **阻塞依赖**: `BlockBy` 字段定义任务间依赖关系
4. **并行执行**: 无依赖的任务可并行执行

**示例**:
```json
{
  "MustInOrder": ["user-service", "order-service", "payment-service"],
  "Tasks": [
    {"AppID": "user-service", "Type": "build"},
    {"AppID": "user-service", "Type": "deploy", "BlockBy": "user-service-build"},
    {"AppID": "order-service", "Type": "deploy", "BlockBy": "user-service-deploy"},
    {"AppID": "payment-service", "Type": "deploy", "BlockBy": "order-service-deploy"}
  ]
}
```

### 4. 部署策略执行

支持多种部署策略:

#### 蓝绿部署
```json
{
  "Strategy": [
    {
      "BatchSize": 0,
      "CanaryRatio": 0,
      "AutoRollback": true
    }
  ]
}
```

#### 金丝雀部署
```json
{
  "Strategy": [
    {
      "BatchSize": 1,
      "BatchInterval": 300,
      "CanaryRatio": 0.1,
      "AutoRollback": true,
      "ManualApprovalStatus": null
    },
    {
      "BatchSize": 0,
      "CanaryRatio": 1.0,
      "AutoRollback": false
    }
  ]
}
```

#### 滚动更新
```json
{
  "Strategy": [
    {
      "BatchSize": 3,
      "BatchInterval": 60,
      "CanaryRatio": 0,
      "AutoRollback": true
    }
  ]
}
```

---

## 🎨 前端工作流展示

### 当前实现 (v1.0)

前端使用 **React Flow** 实现 DAG 工作流可视化。

#### 组件架构

```
DeploymentDetail (页面)
  └── WorkflowViewer (工作流查看器)
        ├── WorkflowNode (自定义节点)
        └── ReactFlow (图表引擎)
```

#### Task 展示策略

**更新**: 前后端统一使用 Task 概念，不再需要额外的映射层。

**当前方案**: 前端直接展示后端返回的 Task 列表

**示例**:
```typescript
// 前端直接使用后端的 Task 数据
tasks: [
  { id: 'task-1', name: '准备部署', type: 'prepare', status: 'success' },
  { id: 'task-2', name: '构建镜像', type: 'build', status: 'success' },
  { id: 'task-3', name: '部署服务', type: 'deploy', status: 'running' },
  { id: 'task-4', name: '健康检查', type: 'health_check', status: 'pending' }
]
```

#### 编辑功能

支持编辑模式（仅 `pending` 和 `waiting_confirm` 状态）:

- ✅ 拖拽节点位置
- ✅ 上移/下移调整顺序
- ✅ 创建/删除连接线
- ✅ 添加新任务
- ✅ 删除任务（Delete/Backspace）

---

## 🔧 待优化事项

### 1. Task 与后端同步 ✅ **已完成**

**更新**: 前后端已统一使用 Task 概念，不再需要额外的聚合逻辑。

**当前实现**:
- 前端直接使用后端返回的 Task 列表
- Task 接口在前后端保持一致
- 通过 `type` 字段区分任务类型（build、test、deploy、health_check 等）

**优势**:
- 减少了前后端的概念差异
- 简化了数据映射逻辑
- 便于实时状态同步

### 2. 实时状态更新

**当前**: 前端每 3 秒轮询

**建议**: 实现 WebSocket 推送

```go
// 伪代码
func (s *Service) ExecuteDeployment(deploymentID string) {
    for _, task := range deployment.Tasks {
        s.executeTask(task)
        
        // 推送状态更新
        s.wsHub.Broadcast(deploymentID, StatusUpdate{
            TaskID: task.ID,
            Status: task.Status,
        })
    }
}
```

### 3. Task 依赖关系可视化

**当前**: 前端任务是线性顺序（A → B → C → D）

**建议**: 支持复杂 DAG（有向无环图）可视化

- 利用 Task 的 `blockBy` 字段展示依赖关系
- 支持并行任务的可视化展示
- 在编辑模式下允许创建复杂的依赖链

**示例**:
```typescript
// 示例：并行构建多个应用
tasks: [
  { id: 'task-1', name: '准备', type: 'prepare' },
  { id: 'task-2a', name: '构建服务A', type: 'build', blockBy: 'task-1' },
  { id: 'task-2b', name: '构建服务B', type: 'build', blockBy: 'task-1' },
  { id: 'task-3', name: '部署', type: 'deploy', blockBy: 'task-2a,task-2b' }
]
```

### 4. Task 日志关联

**当前**: `Task` 包含 `logs` 字段，但未实现详细展示

**建议**: 点击任务节点展开日志面板

```typescript
interface Task {
  id: string
  name: string
  type: string
  status: string
  logs: TaskLog[]
}

interface TaskLog {
  timestamp: string
  level: 'info' | 'warn' | 'error'
  message: string
}
```

**UI 改进**:
- 点击工作流节点展开侧边日志面板
- 实时流式显示日志
- 支持日志级别过滤

### 5. 任务级别人工确认

**当前**: `waiting_confirm` 是部署级别的状态

**建议**: 支持任务级别的确认

```typescript
interface Task {
  id: string
  name: string
  status: 'pending' | 'running' | 'waiting_confirm' | 'success' | 'failed' | 'blocked'
  requireConfirm: boolean
  confirmedBy?: string
  confirmedAt?: string
}
```

**UI 改进**:
- 在需要确认的任务节点上显示"等待确认"徽章
- 点击节点弹出确认对话框
- 记录确认人和确认时间

---

## 📝 API 设计建议

### 获取部署详情（含任务）

```http
GET /api/v1/deployments/:id

Response:
{
  "id": "deploy-001",
  "version": "v1.2.5",
  "status": "running",
  "tasks": [
    {
      "id": "task-1",
      "deploymentId": "deploy-001",
      "appId": "user-service",
      "name": "准备部署",
      "type": "prepare",
      "status": "success",
      "startedAt": "2024-10-21T14:00:00Z",
      "completedAt": "2024-10-21T14:05:00Z",
      "duration": 300,
      "logs": ["检查版本信息...", "验证配置文件...", "准备完成"]
    },
    {
      "id": "task-2",
      "deploymentId": "deploy-001",
      "appId": "user-service",
      "name": "构建镜像",
      "type": "build",
      "status": "running",
      "startedAt": "2024-10-21T14:05:00Z",
      "blockBy": "task-1",
      "logs": ["构建中..."]
    }
  ]
}
```

### 更新工作流编排

```http
PUT /api/v1/deployments/:id/workflow

Request:
{
  "tasks": [
    {
      "id": "task-1",
      "name": "准备部署",
      "type": "prepare",
      "order": 1,
      "blockBy": ""
    },
    {
      "id": "task-2",
      "name": "构建镜像",
      "type": "build",
      "order": 2,
      "blockBy": "task-1"
    }
  ]
}

Response:
{
  "success": true,
  "message": "工作流已更新"
}
```

### 获取任务日志

```http
GET /api/v1/deployments/:id/tasks/:taskId/logs

Response:
{
  "taskId": "task-2",
  "logs": [
    {
      "timestamp": "2024-10-21T14:05:00Z",
      "level": "info",
      "message": "开始构建 user-service"
    },
    {
      "timestamp": "2024-10-21T14:05:30Z",
      "level": "info",
      "message": "镜像构建成功: user-service:v1.2.5"
    }
  ]
}
```

---

## 🎯 下一步行动项

### 短期（本周）
- [x] ~~明确 Step 和 Task 的映射关系~~ - 已统一为 Task 概念
- [x] ~~确定步骤聚合逻辑~~ - 前后端直接使用 Task
- [ ] 设计任务日志展示 UI
- [ ] 实现点击任务节点查看日志功能

### 中期（本月）
- [ ] 后端实现完整的 Task CRUD API
- [ ] 支持基于 `blockBy` 的 DAG 依赖关系可视化
- [ ] 添加 WebSocket 实时推送任务状态
- [ ] 实现任务级别的人工确认

### 长期（下季度）
- [ ] 支持自定义工作流模板
- [ ] 实现工作流版本控制
- [ ] 添加工作流可视化编排器（拖拽设计）
- [ ] 支持工作流的保存和复用

---

## 📚 参考资料

- [GitOps 最佳实践](https://www.weave.works/blog/what-is-gitops-really)
- [Argo CD Workflow](https://argoproj.github.io/argo-workflows/)
- [React Flow 文档](https://reactflow.dev/)
- [Kubernetes Deployment Strategies](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)

---

**文档维护**: 该文档会随着系统演进持续更新。如有疑问或建议，请联系开发团队。

