# API 设计文档

## 概述

本文档描述了基于 GitOps 的持续部署平台（CD）的后端 API 设计。系统使用 Go 语言实现，遵循 RESTful 风格的 HTTP API 设计原则。

## 系统架构

### 模块划分

1. **Service 层**：处理所有 HTTP 请求，提供 RESTful API
2. **DeploymentManager**：接收 GitHub 回调，处理所有部署管理工作
3. **WorkflowManager**：管理和编排所有部署产生的任务
4. **Deploy 模块**：执行实际的部署操作，支持物理机和 Kubernetes

## 核心数据结构

### Version (版本)

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

### Application (应用)

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

### Environment (目标环境)

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

### Deployment (部署)

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

### Task (任务)

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

## HTTP API 文档

### 版本管理 API

#### 创建版本

- **Endpoint**: `POST /api/v1/versions`
- **描述**: 创建新的版本

**请求体**:
```json
{
    "git_tag": "v1.0.0",
    "git_commit": "abc123def456",
    "repository": "https://github.com/org/repo",
    "description": "Release version 1.0.0"
}
```

**响应体**:
```json
{
    "id": "ver_123456",
    "git_tag": "v1.0.0",
    "git_commit": "abc123def456",
    "repository": "https://github.com/org/repo",
    "created_by": "user@example.com",
    "created_at": "2024-10-24T10:00:00Z",
    "description": "Release version 1.0.0"
}
```

#### 获取版本列表

- **Endpoint**: `GET /api/v1/versions`
- **描述**: 获取所有版本列表
- **查询参数**:
  - `repository` (可选): 按仓库过滤
  - `page` (可选): 页码，默认 1
  - `page_size` (可选): 每页数量，默认 20

**响应体**:
```json
{
    "versions": [
        {
            "id": "ver_123456",
            "git_tag": "v1.0.0",
            "git_commit": "abc123def456",
            "repository": "https://github.com/org/repo",
            "created_by": "user@example.com",
            "created_at": "2024-10-24T10:00:00Z",
            "description": "Release version 1.0.0"
        }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
}
```

#### 获取版本详情

- **Endpoint**: `GET /api/v1/versions/{id}`
- **描述**: 获取指定版本详情

**响应体**:
```json
{
    "id": "ver_123456",
    "git_tag": "v1.0.0",
    "git_commit": "abc123def456",
    "repository": "https://github.com/org/repo",
    "created_by": "user@example.com",
    "created_at": "2024-10-24T10:00:00Z",
    "description": "Release version 1.0.0"
}
```

#### 删除版本

- **Endpoint**: `DELETE /api/v1/versions/{id}`
- **描述**: 删除指定版本

**响应体**:
```json
{
    "message": "Version deleted successfully"
}
```

### 应用管理 API

#### 创建应用

- **Endpoint**: `POST /api/v1/applications`
- **描述**: 创建新应用

**请求体**:
```json
{
    "name": "api-service",
    "repository": "https://github.com/org/repo",
    "type": "microservice",
    "config": {
        "dockerfile": "Dockerfile",
        "build_args": "arg1=value1"
    }
}
```

**响应体**:
```json
{
    "id": "app_123456",
    "name": "api-service",
    "repository": "https://github.com/org/repo",
    "type": "microservice",
    "config": {
        "dockerfile": "Dockerfile",
        "build_args": "arg1=value1"
    },
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:00:00Z"
}
```

#### 获取应用列表

- **Endpoint**: `GET /api/v1/applications`
- **描述**: 获取所有应用列表
- **查询参数**:
  - `repository` (可选): 按仓库过滤
  - `type` (可选): 按类型过滤
  - `page` (可选): 页码，默认 1
  - `page_size` (可选): 每页数量，默认 20

**响应体**:
```json
{
    "applications": [
        {
            "id": "app_123456",
            "name": "api-service",
            "repository": "https://github.com/org/repo",
            "type": "microservice",
            "config": {
                "dockerfile": "Dockerfile",
                "build_args": "arg1=value1"
            },
            "created_at": "2024-10-24T10:00:00Z",
            "updated_at": "2024-10-24T10:00:00Z"
        }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
}
```

#### 获取应用详情

- **Endpoint**: `GET /api/v1/applications/{id}`
- **描述**: 获取指定应用详情

**响应体**:
```json
{
    "id": "app_123456",
    "name": "api-service",
    "repository": "https://github.com/org/repo",
    "type": "microservice",
    "config": {
        "dockerfile": "Dockerfile",
        "build_args": "arg1=value1"
    },
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:00:00Z"
}
```

#### 更新应用

- **Endpoint**: `PUT /api/v1/applications/{id}`
- **描述**: 更新应用信息

**请求体**:
```json
{
    "name": "api-service-v2",
    "type": "microservice",
    "config": {
        "dockerfile": "Dockerfile.prod",
        "build_args": "arg1=value2"
    }
}
```

**响应体**:
```json
{
    "id": "app_123456",
    "name": "api-service-v2",
    "repository": "https://github.com/org/repo",
    "type": "microservice",
    "config": {
        "dockerfile": "Dockerfile.prod",
        "build_args": "arg1=value2"
    },
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T11:00:00Z"
}
```

#### 删除应用

- **Endpoint**: `DELETE /api/v1/applications/{id}`
- **描述**: 删除指定应用

**响应体**:
```json
{
    "message": "Application deleted successfully"
}
```

### 环境管理 API

#### 创建环境

- **Endpoint**: `POST /api/v1/environments`
- **描述**: 创建新的部署环境

**请求体**:
```json
{
    "name": "production",
    "type": "kubernetes",
    "config": {
        "cluster": "prod-cluster",
        "namespace": "default",
        "kubeconfig": "base64_encoded_config"
    },
    "is_active": true
}
```

**响应体**:
```json
{
    "id": "env_123456",
    "name": "production",
    "type": "kubernetes",
    "config": {
        "cluster": "prod-cluster",
        "namespace": "default",
        "kubeconfig": "base64_encoded_config"
    },
    "is_active": true,
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:00:00Z"
}
```

#### 获取环境列表

- **Endpoint**: `GET /api/v1/environments`
- **描述**: 获取所有环境列表
- **查询参数**:
  - `type` (可选): 按类型过滤（kubernetes/physical）
  - `is_active` (可选): 按激活状态过滤
  - `page` (可选): 页码，默认 1
  - `page_size` (可选): 每页数量，默认 20

**响应体**:
```json
{
    "environments": [
        {
            "id": "env_123456",
            "name": "production",
            "type": "kubernetes",
            "config": {
                "cluster": "prod-cluster",
                "namespace": "default",
                "kubeconfig": "base64_encoded_config"
            },
            "is_active": true,
            "created_at": "2024-10-24T10:00:00Z",
            "updated_at": "2024-10-24T10:00:00Z"
        }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
}
```

#### 获取环境详情

- **Endpoint**: `GET /api/v1/environments/{id}`
- **描述**: 获取指定环境详情

**响应体**:
```json
{
    "id": "env_123456",
    "name": "production",
    "type": "kubernetes",
    "config": {
        "cluster": "prod-cluster",
        "namespace": "default",
        "kubeconfig": "base64_encoded_config"
    },
    "is_active": true,
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:00:00Z"
}
```

#### 更新环境

- **Endpoint**: `PUT /api/v1/environments/{id}`
- **描述**: 更新环境信息

**请求体**:
```json
{
    "name": "production",
    "type": "kubernetes",
    "config": {
        "cluster": "prod-cluster-v2",
        "namespace": "production",
        "kubeconfig": "base64_encoded_config_new"
    },
    "is_active": true
}
```

**响应体**:
```json
{
    "id": "env_123456",
    "name": "production",
    "type": "kubernetes",
    "config": {
        "cluster": "prod-cluster-v2",
        "namespace": "production",
        "kubeconfig": "base64_encoded_config_new"
    },
    "is_active": true,
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T11:00:00Z"
}
```

#### 删除环境

- **Endpoint**: `DELETE /api/v1/environments/{id}`
- **描述**: 删除指定环境

**响应体**:
```json
{
    "message": "Environment deleted successfully"
}
```

### 部署管理 API

#### 创建部署

- **Endpoint**: `POST /api/v1/deployments`
- **描述**: 创建新的部署任务

**请求体**:
```json
{
    "version_id": "ver_123456",
    "application_ids": ["app_123456", "app_789012"],
    "environment_id": "env_123456"
}
```

**响应体**:
```json
{
    "id": "dep_123456",
    "version_id": "ver_123456",
    "application_ids": ["app_123456", "app_789012"],
    "environment_id": "env_123456",
    "status": "pending",
    "created_by": "user@example.com",
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:00:00Z",
    "started_at": null,
    "completed_at": null,
    "error_message": ""
}
```

#### 获取部署列表

- **Endpoint**: `GET /api/v1/deployments`
- **描述**: 获取所有部署列表
- **查询参数**:
  - `status` (可选): 按状态过滤
  - `environment_id` (可选): 按环境过滤
  - `version_id` (可选): 按版本过滤
  - `page` (可选): 页码，默认 1
  - `page_size` (可选): 每页数量，默认 20

**响应体**:
```json
{
    "deployments": [
        {
            "id": "dep_123456",
            "version_id": "ver_123456",
            "application_ids": ["app_123456", "app_789012"],
            "environment_id": "env_123456",
            "status": "running",
            "created_by": "user@example.com",
            "created_at": "2024-10-24T10:00:00Z",
            "updated_at": "2024-10-24T10:05:00Z",
            "started_at": "2024-10-24T10:05:00Z",
            "completed_at": null,
            "error_message": ""
        }
    ],
    "total": 200,
    "page": 1,
    "page_size": 20
}
```

#### 获取部署详情

- **Endpoint**: `GET /api/v1/deployments/{id}`
- **描述**: 获取指定部署详情

**响应体**:
```json
{
    "id": "dep_123456",
    "version_id": "ver_123456",
    "application_ids": ["app_123456", "app_789012"],
    "environment_id": "env_123456",
    "status": "success",
    "created_by": "user@example.com",
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:10:00Z",
    "started_at": "2024-10-24T10:05:00Z",
    "completed_at": "2024-10-24T10:10:00Z",
    "error_message": ""
}
```

#### 取消部署

- **Endpoint**: `POST /api/v1/deployments/{id}/cancel`
- **描述**: 取消正在进行的部署

**请求体**:
```json
{}
```

**响应体**:
```json
{
    "id": "dep_123456",
    "status": "failed",
    "error_message": "Deployment cancelled by user"
}
```

#### 回滚部署

- **Endpoint**: `POST /api/v1/deployments/{id}/rollback`
- **描述**: 回滚到之前的部署版本

**请求体**:
```json
{
    "target_version_id": "ver_111111"
}
```

**响应体**:
```json
{
    "id": "dep_999999",
    "version_id": "ver_111111",
    "application_ids": ["app_123456", "app_789012"],
    "environment_id": "env_123456",
    "status": "pending",
    "created_by": "user@example.com",
    "created_at": "2024-10-24T11:00:00Z",
    "updated_at": "2024-10-24T11:00:00Z",
    "started_at": null,
    "completed_at": null,
    "error_message": ""
}
```

### 任务管理 API

#### 获取任务列表

- **Endpoint**: `GET /api/v1/tasks`
- **描述**: 获取所有任务列表
- **查询参数**:
  - `deployment_id` (可选): 按部署 ID 过滤
  - `status` (可选): 按状态过滤
  - `type` (可选): 按类型过滤
  - `page` (可选): 页码，默认 1
  - `page_size` (可选): 每页数量，默认 20

**响应体**:
```json
{
    "tasks": [
        {
            "id": "task_123456",
            "deployment_id": "dep_123456",
            "type": "build",
            "status": "success",
            "payload": "{\"build_config\": {...}}",
            "result": "{\"artifact_url\": \"...\"}",
            "created_at": "2024-10-24T10:00:00Z",
            "updated_at": "2024-10-24T10:02:00Z",
            "started_at": "2024-10-24T10:00:10Z",
            "completed_at": "2024-10-24T10:02:00Z"
        }
    ],
    "total": 500,
    "page": 1,
    "page_size": 20
}
```

#### 获取任务详情

- **Endpoint**: `GET /api/v1/tasks/{id}`
- **描述**: 获取指定任务详情

**响应体**:
```json
{
    "id": "task_123456",
    "deployment_id": "dep_123456",
    "type": "build",
    "status": "success",
    "payload": "{\"build_config\": {...}}",
    "result": "{\"artifact_url\": \"...\"}",
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:02:00Z",
    "started_at": "2024-10-24T10:00:10Z",
    "completed_at": "2024-10-24T10:02:00Z"
}
```

#### 重试任务

- **Endpoint**: `POST /api/v1/tasks/{id}/retry`
- **描述**: 重试失败的任务

**请求体**:
```json
{}
```

**响应体**:
```json
{
    "id": "task_123456",
    "deployment_id": "dep_123456",
    "type": "build",
    "status": "pending",
    "payload": "{\"build_config\": {...}}",
    "result": "",
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:15:00Z",
    "started_at": null,
    "completed_at": null
}
```

### GitHub 回调 API

#### 接收 GitHub Webhook

- **Endpoint**: `POST /api/v1/webhooks/github`
- **描述**: 接收 GitHub 的 webhook 回调（由 DeploymentManager 处理）

**请求头**:
```
X-GitHub-Event: push
X-Hub-Signature-256: sha256=...
Content-Type: application/json
```

**请求体（示例 - Push Event）**:
```json
{
    "ref": "refs/tags/v1.0.0",
    "repository": {
        "full_name": "org/repo",
        "clone_url": "https://github.com/org/repo.git"
    },
    "head_commit": {
        "id": "abc123def456",
        "message": "Release v1.0.0"
    },
    "pusher": {
        "name": "username",
        "email": "user@example.com"
    }
}
```

**响应体**:
```json
{
    "message": "Webhook received and processed",
    "version_created": true,
    "version_id": "ver_123456"
}
```

#### 接收 PR 合并事件

- **Endpoint**: `POST /api/v1/webhooks/github`
- **描述**: 处理 PR 合并后自动创建版本

**请求头**:
```
X-GitHub-Event: pull_request
X-Hub-Signature-256: sha256=...
Content-Type: application/json
```

**请求体（示例 - PR Merged Event）**:
```json
{
    "action": "closed",
    "pull_request": {
        "merged": true,
        "merge_commit_sha": "abc123def456",
        "base": {
            "ref": "main"
        },
        "head": {
            "ref": "feature/new-feature"
        },
        "title": "Add new feature"
    },
    "repository": {
        "full_name": "org/repo",
        "clone_url": "https://github.com/org/repo.git"
    }
}
```

**响应体**:
```json
{
    "message": "Pull request merge processed",
    "version_created": true,
    "version_id": "ver_789012",
    "auto_tag": "v1.0.1"
}
```

## 错误处理

所有 API 在发生错误时返回统一的错误格式：

**错误响应体**:
```json
{
    "error": {
        "code": "ERROR_CODE",
        "message": "Human-readable error message",
        "details": {
            "field": "Additional error context"
        }
    }
}
```

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

### 常见错误代码

- `INVALID_REQUEST`: 请求参数验证失败
- `RESOURCE_NOT_FOUND`: 资源不存在
- `RESOURCE_CONFLICT`: 资源已存在或状态冲突
- `UNAUTHORIZED`: 认证失败
- `FORBIDDEN`: 权限不足
- `INTERNAL_ERROR`: 内部服务错误
- `DEPLOYMENT_FAILED`: 部署失败
- `TASK_FAILED`: 任务执行失败

## 认证与授权

所有 API 请求需要在 HTTP 头中携带认证令牌：

```
Authorization: Bearer <token>
```

## 分页

支持分页的 API 使用以下查询参数：

- `page`: 页码（从 1 开始）
- `page_size`: 每页数量（默认 20，最大 100）

分页响应包含以下字段：

```json
{
    "data": [...],
    "total": 100,
    "page": 1,
    "page_size": 20
}
```

## 模块接口设计

本节定义各模块的 Go 接口，用于规范模块之间的交互和依赖关系。

### Service 层接口

Service 层负责处理所有 HTTP 请求，提供 RESTful API。

```go
type VersionService interface {
    CreateVersion(ctx context.Context, req *CreateVersionRequest) (*Version, error)
    GetVersionList(ctx context.Context, req *ListVersionsRequest) (*VersionListResponse, error)
    GetVersion(ctx context.Context, id string) (*Version, error)
    DeleteVersion(ctx context.Context, id string) error
}

type ApplicationService interface {
    CreateApplication(ctx context.Context, req *CreateApplicationRequest) (*Application, error)
    GetApplicationList(ctx context.Context, req *ListApplicationsRequest) (*ApplicationListResponse, error)
    GetApplication(ctx context.Context, id string) (*Application, error)
    UpdateApplication(ctx context.Context, id string, req *UpdateApplicationRequest) (*Application, error)
    DeleteApplication(ctx context.Context, id string) error
}

type EnvironmentService interface {
    CreateEnvironment(ctx context.Context, req *CreateEnvironmentRequest) (*Environment, error)
    GetEnvironmentList(ctx context.Context, req *ListEnvironmentsRequest) (*EnvironmentListResponse, error)
    GetEnvironment(ctx context.Context, id string) (*Environment, error)
    UpdateEnvironment(ctx context.Context, id string, req *UpdateEnvironmentRequest) (*Environment, error)
    DeleteEnvironment(ctx context.Context, id string) error
}

type DeploymentService interface {
    CreateDeployment(ctx context.Context, req *CreateDeploymentRequest) (*Deployment, error)
    GetDeploymentList(ctx context.Context, req *ListDeploymentsRequest) (*DeploymentListResponse, error)
    GetDeployment(ctx context.Context, id string) (*Deployment, error)
    CancelDeployment(ctx context.Context, id string) (*Deployment, error)
    RollbackDeployment(ctx context.Context, id string, req *RollbackRequest) (*Deployment, error)
}

type TaskService interface {
    GetTaskList(ctx context.Context, req *ListTasksRequest) (*TaskListResponse, error)
    GetTask(ctx context.Context, id string) (*Task, error)
    RetryTask(ctx context.Context, id string) (*Task, error)
}

type WebhookService interface {
    HandleGitHubWebhook(ctx context.Context, event string, payload []byte) (*WebhookResponse, error)
}
```

### DeploymentManager 接口

DeploymentManager 负责接收 GitHub 回调和管理部署生命周期。

```go
type DeploymentManager interface {
    ProcessWebhookEvent(ctx context.Context, event *GitHubEvent) error
    
    CreateVersionFromTag(ctx context.Context, tag *GitTag) (*Version, error)
    
    CreateVersionFromPR(ctx context.Context, pr *PullRequest) (*Version, error)
    
    StartDeployment(ctx context.Context, deploymentID string) error
    
    CancelDeployment(ctx context.Context, deploymentID string) error
    
    UpdateDeploymentStatus(ctx context.Context, deploymentID string, status DeploymentStatus, errorMsg string) error
    
    GetDeploymentProgress(ctx context.Context, deploymentID string) (*DeploymentProgress, error)
}

type GitHubEvent struct {
    Type       string
    Repository string
    Payload    interface{}
}

type GitTag struct {
    Name       string
    Commit     string
    Repository string
    Pusher     string
    Message    string
}

type PullRequest struct {
    Number       int
    Title        string
    MergeCommit  string
    Repository   string
    BaseBranch   string
    HeadBranch   string
    MergedBy     string
}

type DeploymentProgress struct {
    DeploymentID string
    Status       DeploymentStatus
    TotalTasks   int
    CompletedTasks int
    FailedTasks  int
    CurrentTask  *Task
}
```

### WorkflowManager 接口

WorkflowManager 负责管理和编排部署任务的执行。

```go
type WorkflowManager interface {
    CreateWorkflow(ctx context.Context, deployment *Deployment) (*Workflow, error)
    
    ExecuteWorkflow(ctx context.Context, workflowID string) error
    
    GetWorkflowStatus(ctx context.Context, workflowID string) (*WorkflowStatus, error)
    
    CancelWorkflow(ctx context.Context, workflowID string) error
    
    RetryFailedTasks(ctx context.Context, workflowID string) error
}

type TaskScheduler interface {
    ScheduleTask(ctx context.Context, task *Task) error
    
    GetNextTask(ctx context.Context) (*Task, error)
    
    UpdateTaskStatus(ctx context.Context, taskID string, status TaskStatus, result string) error
    
    GetTasksByDeployment(ctx context.Context, deploymentID string) ([]*Task, error)
}

type Workflow struct {
    ID           string
    DeploymentID string
    Tasks        []*Task
    Status       WorkflowStatus
    CreatedAt    time.Time
    UpdatedAt    time.Time
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

### Deploy 模块接口

Deploy 模块负责执行实际的部署操作，支持 Kubernetes 和物理机部署。

```go
type Deployer interface {
    Deploy(ctx context.Context, req *DeployRequest) (*DeployResult, error)
    
    Rollback(ctx context.Context, req *RollbackDeployRequest) (*DeployResult, error)
    
    GetDeploymentInfo(ctx context.Context, deploymentID string) (*DeploymentInfo, error)
    
    HealthCheck(ctx context.Context, deploymentID string) (*HealthCheckResult, error)
}

type KubernetesDeployer interface {
    Deployer
    
    ApplyManifest(ctx context.Context, namespace string, manifest []byte) error
    
    GetPodStatus(ctx context.Context, namespace, selector string) ([]*PodStatus, error)
    
    ScaleDeployment(ctx context.Context, namespace, name string, replicas int32) error
    
    GetLogs(ctx context.Context, namespace, podName string, lines int) (string, error)
}

type PhysicalDeployer interface {
    Deployer
    
    UploadArtifact(ctx context.Context, hosts []string, artifact *Artifact) error
    
    ExecuteCommand(ctx context.Context, hosts []string, command string) (*CommandResult, error)
    
    RestartService(ctx context.Context, hosts []string, serviceName string) error
    
    CheckServiceStatus(ctx context.Context, hosts []string, serviceName string) ([]*ServiceStatus, error)
}

type DeployRequest struct {
    DeploymentID  string
    Version       *Version
    Applications  []*Application
    Environment   *Environment
    Config        map[string]string
}

type RollbackDeployRequest struct {
    DeploymentID       string
    TargetDeploymentID string
    Environment        *Environment
}

type DeployResult struct {
    Success      bool
    Message      string
    DeploymentID string
    Details      map[string]interface{}
}

type DeploymentInfo struct {
    DeploymentID string
    Status       string
    Replicas     int32
    ReadyReplicas int32
    UpdatedAt    time.Time
}

type HealthCheckResult struct {
    Healthy   bool
    Message   string
    Checks    []*HealthCheck
}

type HealthCheck struct {
    Name    string
    Status  string
    Message string
}

type PodStatus struct {
    Name      string
    Phase     string
    Ready     bool
    Restarts  int32
    NodeName  string
}

type Artifact struct {
    Name    string
    Version string
    Path    string
    Size    int64
}

type CommandResult struct {
    Host     string
    Success  bool
    Output   string
    Error    string
}

type ServiceStatus struct {
    Host      string
    Running   bool
    Status    string
    PID       int
    Uptime    string
}
```

### 数据访问层接口

定义数据持久化相关的接口。

```go
type VersionRepository interface {
    Create(ctx context.Context, version *Version) error
    GetByID(ctx context.Context, id string) (*Version, error)
    List(ctx context.Context, filter *VersionFilter) ([]*Version, int, error)
    Delete(ctx context.Context, id string) error
}

type ApplicationRepository interface {
    Create(ctx context.Context, app *Application) error
    GetByID(ctx context.Context, id string) (*Application, error)
    List(ctx context.Context, filter *ApplicationFilter) ([]*Application, int, error)
    Update(ctx context.Context, app *Application) error
    Delete(ctx context.Context, id string) error
}

type EnvironmentRepository interface {
    Create(ctx context.Context, env *Environment) error
    GetByID(ctx context.Context, id string) (*Environment, error)
    List(ctx context.Context, filter *EnvironmentFilter) ([]*Environment, int, error)
    Update(ctx context.Context, env *Environment) error
    Delete(ctx context.Context, id string) error
}

type DeploymentRepository interface {
    Create(ctx context.Context, deployment *Deployment) error
    GetByID(ctx context.Context, id string) (*Deployment, error)
    List(ctx context.Context, filter *DeploymentFilter) ([]*Deployment, int, error)
    Update(ctx context.Context, deployment *Deployment) error
}

type TaskRepository interface {
    Create(ctx context.Context, task *Task) error
    GetByID(ctx context.Context, id string) (*Task, error)
    List(ctx context.Context, filter *TaskFilter) ([]*Task, int, error)
    Update(ctx context.Context, task *Task) error
    GetByDeploymentID(ctx context.Context, deploymentID string) ([]*Task, error)
}

type VersionFilter struct {
    Repository string
    Page       int
    PageSize   int
}

type ApplicationFilter struct {
    Repository string
    Type       string
    Page       int
    PageSize   int
}

type EnvironmentFilter struct {
    Type     string
    IsActive *bool
    Page     int
    PageSize int
}

type DeploymentFilter struct {
    Status        DeploymentStatus
    EnvironmentID string
    VersionID     string
    Page          int
    PageSize      int
}

type TaskFilter struct {
    DeploymentID string
    Status       TaskStatus
    Type         string
    Page         int
    PageSize     int
}
```

### 请求和响应类型

定义 API 请求和响应的数据结构。

```go
type CreateVersionRequest struct {
    GitTag      string `json:"git_tag" binding:"required"`
    GitCommit   string `json:"git_commit" binding:"required"`
    Repository  string `json:"repository" binding:"required"`
    Description string `json:"description"`
}

type ListVersionsRequest struct {
    Repository string `form:"repository"`
    Page       int    `form:"page" binding:"min=1"`
    PageSize   int    `form:"page_size" binding:"min=1,max=100"`
}

type VersionListResponse struct {
    Versions []*Version `json:"versions"`
    Total    int        `json:"total"`
    Page     int        `json:"page"`
    PageSize int        `json:"page_size"`
}

type CreateApplicationRequest struct {
    Name       string            `json:"name" binding:"required"`
    Repository string            `json:"repository" binding:"required"`
    Type       string            `json:"type" binding:"required"`
    Config     map[string]string `json:"config"`
}

type UpdateApplicationRequest struct {
    Name   string            `json:"name"`
    Type   string            `json:"type"`
    Config map[string]string `json:"config"`
}

type ListApplicationsRequest struct {
    Repository string `form:"repository"`
    Type       string `form:"type"`
    Page       int    `form:"page" binding:"min=1"`
    PageSize   int    `form:"page_size" binding:"min=1,max=100"`
}

type ApplicationListResponse struct {
    Applications []*Application `json:"applications"`
    Total        int            `json:"total"`
    Page         int            `json:"page"`
    PageSize     int            `json:"page_size"`
}

type CreateEnvironmentRequest struct {
    Name     string            `json:"name" binding:"required"`
    Type     string            `json:"type" binding:"required"`
    Config   map[string]string `json:"config" binding:"required"`
    IsActive bool              `json:"is_active"`
}

type UpdateEnvironmentRequest struct {
    Name     string            `json:"name"`
    Type     string            `json:"type"`
    Config   map[string]string `json:"config"`
    IsActive *bool             `json:"is_active"`
}

type ListEnvironmentsRequest struct {
    Type     string `form:"type"`
    IsActive *bool  `form:"is_active"`
    Page     int    `form:"page" binding:"min=1"`
    PageSize int    `form:"page_size" binding:"min=1,max=100"`
}

type EnvironmentListResponse struct {
    Environments []*Environment `json:"environments"`
    Total        int            `json:"total"`
    Page         int            `json:"page"`
    PageSize     int            `json:"page_size"`
}

type CreateDeploymentRequest struct {
    VersionID      string   `json:"version_id" binding:"required"`
    ApplicationIDs []string `json:"application_ids" binding:"required"`
    EnvironmentID  string   `json:"environment_id" binding:"required"`
}

type ListDeploymentsRequest struct {
    Status        string `form:"status"`
    EnvironmentID string `form:"environment_id"`
    VersionID     string `form:"version_id"`
    Page          int    `form:"page" binding:"min=1"`
    PageSize      int    `form:"page_size" binding:"min=1,max=100"`
}

type DeploymentListResponse struct {
    Deployments []*Deployment `json:"deployments"`
    Total       int           `json:"total"`
    Page        int           `json:"page"`
    PageSize    int           `json:"page_size"`
}

type RollbackRequest struct {
    TargetVersionID string `json:"target_version_id" binding:"required"`
}

type ListTasksRequest struct {
    DeploymentID string `form:"deployment_id"`
    Status       string `form:"status"`
    Type         string `form:"type"`
    Page         int    `form:"page" binding:"min=1"`
    PageSize     int    `form:"page_size" binding:"min=1,max=100"`
}

type TaskListResponse struct {
    Tasks    []*Task `json:"tasks"`
    Total    int     `json:"total"`
    Page     int     `json:"page"`
    PageSize int     `json:"page_size"`
}

type WebhookResponse struct {
    Message        string `json:"message"`
    VersionCreated bool   `json:"version_created,omitempty"`
    VersionID      string `json:"version_id,omitempty"`
    AutoTag        string `json:"auto_tag,omitempty"`
}

type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
    Code    string                 `json:"code"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}
```

## 模块交互流程

### 部署流程

1. **用户发起部署**
   - 用户通过 Service 层 API 创建部署：`POST /api/v1/deployments`
   - Service 层验证请求参数并创建 Deployment 记录

2. **DeploymentManager 处理部署**
   - 接收到新的 Deployment 记录
   - 验证版本、应用和环境的有效性
   - 将部署状态更新为 `running`

3. **WorkflowManager 编排任务**
   - 根据部署需求创建一系列任务（Task）
   - 任务类型包括：构建（build）、测试（test）、部署（deploy）
   - 按依赖关系编排任务执行顺序

4. **Deploy 模块执行部署**
   - 接收 WorkflowManager 分配的部署任务
   - 根据环境类型（Kubernetes/物理机）选择部署策略
   - 执行实际的部署操作
   - 更新任务状态和结果

5. **状态回传**
   - Deploy 模块更新 Task 状态
   - WorkflowManager 汇总所有任务状态
   - DeploymentManager 更新 Deployment 最终状态
   - Service 层向用户返回部署结果

### GitHub 回调处理流程

1. **GitHub 触发 Webhook**
   - PR 合并或创建 Tag 时触发 webhook
   - GitHub 发送 POST 请求到 `/api/v1/webhooks/github`

2. **Service 层接收并验证**
   - 验证 GitHub 签名
   - 解析事件类型和数据

3. **DeploymentManager 处理事件**
   - 根据事件类型执行相应操作：
     - PR 合并：自动创建新版本
     - Tag 创建：关联或创建版本记录
   - 可选：自动触发部署到指定环境

4. **返回处理结果**
   - 向 GitHub 返回处理状态
   - 更新相关 Issue 或 PR 状态

## 配置示例

### Kubernetes 环境配置

```json
{
    "type": "kubernetes",
    "config": {
        "cluster": "prod-cluster",
        "namespace": "production",
        "kubeconfig": "base64_encoded_kubeconfig",
        "deployment_strategy": "rolling-update",
        "health_check_enabled": true,
        "timeout": 600
    }
}
```

### 物理机环境配置

```json
{
    "type": "physical",
    "config": {
        "hosts": ["192.168.1.10", "192.168.1.11"],
        "ssh_key": "base64_encoded_ssh_key",
        "deploy_path": "/opt/applications",
        "service_name": "api-service",
        "restart_command": "systemctl restart api-service"
    }
}
```

## 最佳实践

1. **版本管理**
   - PR 合并后自动创建版本
   - 使用语义化版本号（Semantic Versioning）
   - 保持版本与 Git Tag 一致

2. **部署策略**
   - 生产环境使用滚动更新
   - 测试环境可以使用蓝绿部署
   - 支持灰度发布和金丝雀发布

3. **错误处理**
   - 部署失败时自动回滚
   - 保留最近 N 个版本用于快速回滚
   - 详细记录错误日志

4. **监控与告警**
   - 实时监控部署进度
   - 部署失败时发送告警通知
   - 记录部署历史和审计日志
