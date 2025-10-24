# 核心数据结构

## 概述

本文档定义了持续部署平台的核心数据结构。这些数据结构被所有模块共享使用。

## Version (版本)

版本对应 Git Tag，记录发布版本信息。

```go
type Version struct {
    ID          string    `json:"id"`
    GitTag      string    `json:"git_tag"`
    GitCommit   string    `json:"git_commit"`
    Repository  string    `json:"repository"`
    CreatedBy   string    `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
    Description string    `json:"description"`
}
```

**字段说明**:
- `ID`: 版本唯一标识符
- `GitTag`: Git 标签名称（如 v1.0.0）
- `GitCommit`: Git 提交哈希值
- `Repository`: 代码仓库地址
- `CreatedBy`: 创建者
- `CreatedAt`: 创建时间
- `Description`: 版本描述

## Application (应用)

应用是运行的基础单位，对应单体服务。

```go
type Application struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Repository  string            `json:"repository"`
    Type        string            `json:"type"`
    Config      map[string]string `json:"config"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}
```

**字段说明**:
- `ID`: 应用唯一标识符
- `Name`: 应用名称
- `Repository`: 代码仓库地址
- `Type`: 应用类型（如 microservice, monolith）
- `Config`: 应用配置（构建参数、运行参数等）
- `CreatedAt`: 创建时间
- `UpdatedAt`: 更新时间

## Environment (目标环境)

目标环境是部署目标，支持逻辑隔离。

```go
type Environment struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Type        string            `json:"type"`
    Config      map[string]string `json:"config"`
    IsActive    bool              `json:"is_active"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}
```

**字段说明**:
- `ID`: 环境唯一标识符
- `Name`: 环境名称（如 production, staging, development）
- `Type`: 环境类型（kubernetes 或 physical）
- `Config`: 环境配置（集群信息、主机列表等）
- `IsActive`: 是否激活
- `CreatedAt`: 创建时间
- `UpdatedAt`: 更新时间

### Kubernetes 环境配置示例

```json
{
    "cluster": "prod-cluster",
    "namespace": "production",
    "kubeconfig": "base64_encoded_kubeconfig",
    "deployment_strategy": "rolling-update",
    "health_check_enabled": true,
    "timeout": 600
}
```

### 物理机环境配置示例

```json
{
    "hosts": ["192.168.1.10", "192.168.1.11"],
    "ssh_key": "base64_encoded_ssh_key",
    "deploy_path": "/opt/applications",
    "service_name": "api-service",
    "restart_command": "systemctl restart api-service"
}
```

## Deployment (部署)

部署是发布的基础单元，包含版本、应用和环境信息。

```go
type Deployment struct {
    ID             string           `json:"id"`
    VersionID      string           `json:"version_id"`
    ApplicationIDs []string         `json:"application_ids"`
    EnvironmentID  string           `json:"environment_id"`
    Status         DeploymentStatus `json:"status"`
    CreatedBy      string           `json:"created_by"`
    CreatedAt      time.Time        `json:"created_at"`
    UpdatedAt      time.Time        `json:"updated_at"`
    StartedAt      *time.Time       `json:"started_at,omitempty"`
    CompletedAt    *time.Time       `json:"completed_at,omitempty"`
    ErrorMessage   string           `json:"error_message,omitempty"`
}

type DeploymentStatus string

const (
    DeploymentStatusPending    DeploymentStatus = "pending"
    DeploymentStatusRunning    DeploymentStatus = "running"
    DeploymentStatusSuccess    DeploymentStatus = "success"
    DeploymentStatusFailed     DeploymentStatus = "failed"
    DeploymentStatusRolledBack DeploymentStatus = "rolled_back"
)
```

**字段说明**:
- `ID`: 部署唯一标识符
- `VersionID`: 版本 ID
- `ApplicationIDs`: 应用 ID 列表（支持多应用部署）
- `EnvironmentID`: 目标环境 ID
- `Status`: 部署状态
- `CreatedBy`: 创建者
- `CreatedAt`: 创建时间
- `UpdatedAt`: 更新时间
- `StartedAt`: 开始时间
- `CompletedAt`: 完成时间
- `ErrorMessage`: 错误信息

## Task (任务)

任务是部署过程中的具体执行任务。

```go
type Task struct {
    ID           string     `json:"id"`
    DeploymentID string     `json:"deployment_id"`
    Type         string     `json:"type"`
    Status       TaskStatus `json:"status"`
    Payload      string     `json:"payload"`
    Result       string     `json:"result,omitempty"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
    StartedAt    *time.Time `json:"started_at,omitempty"`
    CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

type TaskStatus string

const (
    TaskStatusPending   TaskStatus = "pending"
    TaskStatusRunning   TaskStatus = "running"
    TaskStatusSuccess   TaskStatus = "success"
    TaskStatusFailed    TaskStatus = "failed"
)
```

**字段说明**:
- `ID`: 任务唯一标识符
- `DeploymentID`: 所属部署 ID
- `Type`: 任务类型（如 build, test, deploy）
- `Status`: 任务状态
- `Payload`: 任务输入参数（JSON 格式）
- `Result`: 任务执行结果（JSON 格式）
- `CreatedAt`: 创建时间
- `UpdatedAt`: 更新时间
- `StartedAt`: 开始时间
- `CompletedAt`: 完成时间

## Workflow (工作流)

工作流定义了一组任务的执行流程。

```go
type Workflow struct {
    ID           string         `json:"id"`
    DeploymentID string         `json:"deployment_id"`
    Tasks        []*Task        `json:"tasks"`
    Status       WorkflowStatus `json:"status"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
}

type WorkflowStatus string

const (
    WorkflowStatusPending   WorkflowStatus = "pending"
    WorkflowStatusRunning   WorkflowStatus = "running"
    WorkflowStatusSuccess   WorkflowStatus = "success"
    WorkflowStatusFailed    WorkflowStatus = "failed"
    WorkflowStatusCancelled WorkflowStatus = "cancelled"
)
```

**字段说明**:
- `ID`: 工作流唯一标识符
- `DeploymentID`: 所属部署 ID
- `Tasks`: 任务列表
- `Status`: 工作流状态
- `CreatedAt`: 创建时间
- `UpdatedAt`: 更新时间

## 通用响应结构

### 分页响应

```go
type PaginationResponse struct {
    Total    int `json:"total"`
    Page     int `json:"page"`
    PageSize int `json:"page_size"`
}
```

### 错误响应

```go
type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
    Code    string                 `json:"code"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}
```

### 常见错误代码

- `INVALID_REQUEST`: 请求参数验证失败
- `RESOURCE_NOT_FOUND`: 资源不存在
- `RESOURCE_CONFLICT`: 资源已存在或状态冲突
- `UNAUTHORIZED`: 认证失败
- `FORBIDDEN`: 权限不足
- `INTERNAL_ERROR`: 内部服务错误
- `DEPLOYMENT_FAILED`: 部署失败
- `TASK_FAILED`: 任务执行失败

### HTTP 状态码

- `200 OK`: 请求成功
- `201 Created`: 资源创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未授权
- `403 Forbidden`: 无权限访问
- `404 Not Found`: 资源不存在
- `409 Conflict`: 资源冲突
- `500 Internal Server Error`: 服务器内部错误
- `503 Service Unavailable`: 服务不可用
