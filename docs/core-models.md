# 核心数据结构

## 概述

本文档定义了持续部署平台的核心数据结构。这些数据结构被所有模块共享使用。

## Version (版本)

版本对应 Git Tag，记录发布版本信息。

```go
type Version struct {
    ID            string            `json:"id"`
    Name          string            `json:"name"`
    GitTag        string            `json:"git_tag"`
    GitCommit     string            `json:"git_commit"`
    Repository    string            `json:"repository"`
    CreatedBy     string            `json:"created_by"`
    CreatedAt     time.Time         `json:"created_at"`
    Metadata      map[string]string `json:"metadata,omitempty"`
    Status        string            `json:"status"`
}
```

## Application (应用)

应用是运行的基础单位，对应单体服务。

```go
type Application struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Repository  string    `json:"repository"`
    Description string    `json:"description"`
    Owner       string    `json:"owner"`
    Config      AppConfig `json:"config"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type AppConfig struct {
    BuildConfig   BuildConfig   `json:"build_config"`
    RuntimeConfig RuntimeConfig `json:"runtime_config"`
    HealthCheck   HealthCheck   `json:"health_check"`
}

type BuildConfig struct {
    Dockerfile string            `json:"dockerfile"`
    BuildArgs  map[string]string `json:"build_args"`
    Context    string            `json:"context"`
}

type RuntimeConfig struct {
    Port      int               `json:"port"`
    Env       map[string]string `json:"env"`
    Resources Resources         `json:"resources"`
}

type Resources struct {
    CPU    string `json:"cpu"`
    Memory string `json:"memory"`
}

type HealthCheck struct {
    Path         string `json:"path"`
    Port         int    `json:"port"`
    InitialDelay int    `json:"initial_delay"`
    Period       int    `json:"period"`
}
```

## Environment (目标环境)

目标环境是部署目标，支持逻辑隔离。

```go
type TargetEnvironment struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Type      string    `json:"type"`
    Region    string    `json:"region"`
    Config    EnvConfig `json:"config"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}

type EnvConfig struct {
    K8SConfig      *K8SConfig      `json:"k8s_config,omitempty"`
    PhysicalConfig *PhysicalConfig `json:"physical_config,omitempty"`
}

type K8SConfig struct {
    KubeConfig  string `json:"kube_config"`
    Namespace   string `json:"namespace"`
    ClusterName string `json:"cluster_name"`
}

type PhysicalConfig struct {
    Hosts     []Host    `json:"hosts"`
    SSHConfig SSHConfig `json:"ssh_config"`
}

type Host struct {
    IP       string `json:"ip"`
    Hostname string `json:"hostname"`
    Role     string `json:"role"`
}

type SSHConfig struct {
    User    string `json:"user"`
    Port    int    `json:"port"`
    KeyPath string `json:"key_path"`
}
```

## Deployment (部署)

部署是发布的基础单元，包含版本、应用和环境信息。

```go
type Deployment struct {
    ID           string         `json:"id"`
    Name         string         `json:"name"`
    VersionID    string         `json:"version_id"`
    Applications []string       `json:"apps"`
    TargetEnvs   []string       `json:"target_envs"`
    Details      []TaskDetail   `json:"task_details"`
    Strategy     DeployStrategy `json:"strategy"`
    Status       string         `json:"status"`
    Progress     Progress       `json:"progress"`
    Approvals    []Approval     `json:"approvals"`
    CreatedBy    string         `json:"created_by"`
    CreatedAt    time.Time      `json:"created_at"`
    StartedAt    *time.Time     `json:"started_at,omitempty"`
    CompletedAt  *time.Time     `json:"completed_at,omitempty"`
}

type DeployStrategy struct {
    Type           string  `json:"type"`
    BatchSize      int     `json:"batch_size"`
    BatchInterval  int     `json:"batch_interval"`
    CanaryRatio    float64 `json:"canary_ratio"`
    AutoRollback   bool    `json:"auto_rollback"`
    ManualApproval bool    `json:"manual_approval"`
}

type Progress struct {
    Total        int          `json:"total"`
    Completed    int          `json:"completed"`
    Failed       int          `json:"failed"`
    CurrentBatch int          `json:"current_batch"`
}

type TaskDetail struct {
    ID           string     `json:"id"`
    DeploymentID string     `json:"deployment_id"`
    AppID        string     `json:"app_id"`
    EnvID        string     `json:"env_id"`
    Status       string     `json:"status"`
    StartedAt    *time.Time `json:"started_at,omitempty"`
    CompletedAt  *time.Time `json:"completed_at,omitempty"`
    ErrorMsg     string     `json:"error_msg,omitempty"`
}

type Approval struct {
    ApproverID string    `json:"approver_id"`
    Action     string    `json:"action"`
    Comment    string    `json:"comment"`
    Timestamp  time.Time `json:"timestamp"`
}
```

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
