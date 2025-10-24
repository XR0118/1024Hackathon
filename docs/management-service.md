# 应用/部署/工作流管理服务

## 概述

管理服务负责:
1. 应用和部署的管理
2. 从部署生成工作流
3. 工作流任务的整体调度

## 架构设计

```
API 层 → Service 层 → WorkflowManager → TaskScheduler → Deploy 服务
```

## HTTP API

### 版本管理 API

#### 创建版本

**Endpoint**: `POST /api/v1/versions`

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

**Endpoint**: `GET /api/v1/versions`

**查询参数**:
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

**Endpoint**: `GET /api/v1/versions/{id}`

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

**Endpoint**: `DELETE /api/v1/versions/{id}`

**响应体**:
```json
{
    "message": "Version deleted successfully"
}
```

### 应用管理 API

#### 创建应用

**Endpoint**: `POST /api/v1/applications`

**请求体**:
```json
{
    "name": "api-service",
    "repository": "https://github.com/org/repo",
    "type": "microservice",
    "config": {
        "dockerfile": "Dockerfile",
        "build_args": "arg1=value1",
        "port": "8080"
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
        "build_args": "arg1=value1",
        "port": "8080"
    },
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:00:00Z"
}
```

#### 获取应用列表

**Endpoint**: `GET /api/v1/applications`

**查询参数**:
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
                "build_args": "arg1=value1",
                "port": "8080"
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

**Endpoint**: `GET /api/v1/applications/{id}`

**响应体**:
```json
{
    "id": "app_123456",
    "name": "api-service",
    "repository": "https://github.com/org/repo",
    "type": "microservice",
    "config": {
        "dockerfile": "Dockerfile",
        "build_args": "arg1=value1",
        "port": "8080"
    },
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:00:00Z"
}
```

#### 更新应用

**Endpoint**: `PUT /api/v1/applications/{id}`

**请求体**:
```json
{
    "name": "api-service-v2",
    "type": "microservice",
    "config": {
        "dockerfile": "Dockerfile.prod",
        "build_args": "arg1=value2",
        "port": "8080"
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
        "build_args": "arg1=value2",
        "port": "8080"
    },
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T11:00:00Z"
}
```

#### 删除应用

**Endpoint**: `DELETE /api/v1/applications/{id}`

**响应体**:
```json
{
    "message": "Application deleted successfully"
}
```

### 环境管理 API

#### 创建环境

**Endpoint**: `POST /api/v1/environments`

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

**Endpoint**: `GET /api/v1/environments`

**查询参数**:
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

**Endpoint**: `GET /api/v1/environments/{id}`

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

**Endpoint**: `PUT /api/v1/environments/{id}`

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

**Endpoint**: `DELETE /api/v1/environments/{id}`

**响应体**:
```json
{
    "message": "Environment deleted successfully"
}
```

### 部署管理 API

#### 创建部署

**Endpoint**: `POST /api/v1/deployments`

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

**Endpoint**: `GET /api/v1/deployments`

**查询参数**:
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

**Endpoint**: `GET /api/v1/deployments/{id}`

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

**Endpoint**: `POST /api/v1/deployments/{id}/cancel`

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

**Endpoint**: `POST /api/v1/deployments/{id}/rollback`

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

**Endpoint**: `GET /api/v1/tasks`

**查询参数**:
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
            "payload": "{\"build_config\": {\"dockerfile\": \"Dockerfile\"}}",
            "result": "{\"artifact_url\": \"https://...\", \"image\": \"org/app:v1.0.0\"}",
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

**Endpoint**: `GET /api/v1/tasks/{id}`

**响应体**:
```json
{
    "id": "task_123456",
    "deployment_id": "dep_123456",
    "type": "build",
    "status": "success",
    "payload": "{\"build_config\": {\"dockerfile\": \"Dockerfile\"}}",
    "result": "{\"artifact_url\": \"https://...\", \"image\": \"org/app:v1.0.0\"}",
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:02:00Z",
    "started_at": "2024-10-24T10:00:10Z",
    "completed_at": "2024-10-24T10:02:00Z"
}
```

#### 重试任务

**Endpoint**: `POST /api/v1/tasks/{id}/retry`

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
    "payload": "{\"build_config\": {\"dockerfile\": \"Dockerfile\"}}",
    "result": "",
    "created_at": "2024-10-24T10:00:00Z",
    "updated_at": "2024-10-24T10:15:00Z",
    "started_at": null,
    "completed_at": null
}
```

## Go 接口设计

### Service 层接口

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
```

### WorkflowManager 接口

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
```

### Repository 接口

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
```

## 部署流程

### 1. 创建部署

用户通过 API 创建部署:

```
POST /api/v1/deployments
{
    "version_id": "ver_123456",
    "application_ids": ["app_123456"],
    "environment_id": "env_123456"
}
```

### 2. 生成工作流

WorkflowManager 接收到新部署后，生成工作流:

```go
workflow := workflowManager.CreateWorkflow(ctx, deployment)
```

工作流包含以下任务:
1. **Build Task**: 构建应用镜像或制品
2. **Test Task**: 运行测试（可选）
3. **Deploy Task**: 执行实际部署
4. **HealthCheck Task**: 健康检查

### 3. 调度任务

TaskScheduler 按依赖顺序调度任务:

```go
for _, task := range workflow.Tasks {
    taskScheduler.ScheduleTask(ctx, task)
}
```

### 4. 执行任务

Deploy 服务接收任务并执行:

```
每个任务完成后更新状态:
- TaskStatusPending → TaskStatusRunning → TaskStatusSuccess/Failed
```

### 5. 更新部署状态

WorkflowManager 汇总所有任务状态，更新部署:

```
所有任务成功 → DeploymentStatusSuccess
任意任务失败 → DeploymentStatusFailed
```

## 工作流编排示例

### 标准部署工作流

```json
{
    "workflow_id": "wf_123456",
    "deployment_id": "dep_123456",
    "tasks": [
        {
            "id": "task_001",
            "type": "build",
            "depends_on": [],
            "config": {
                "dockerfile": "Dockerfile",
                "context": ".",
                "tags": ["org/app:v1.0.0", "org/app:latest"]
            }
        },
        {
            "id": "task_002",
            "type": "test",
            "depends_on": ["task_001"],
            "config": {
                "test_command": "go test ./...",
                "coverage_threshold": 80
            }
        },
        {
            "id": "task_003",
            "type": "deploy",
            "depends_on": ["task_002"],
            "config": {
                "strategy": "rolling-update",
                "replicas": 3,
                "max_unavailable": 1
            }
        },
        {
            "id": "task_004",
            "type": "health_check",
            "depends_on": ["task_003"],
            "config": {
                "endpoint": "/health",
                "expected_status": 200,
                "timeout": 60
            }
        }
    ]
}
```

### 蓝绿部署工作流

```json
{
    "workflow_id": "wf_789012",
    "deployment_id": "dep_789012",
    "tasks": [
        {
            "id": "task_001",
            "type": "build",
            "depends_on": []
        },
        {
            "id": "task_002",
            "type": "deploy_green",
            "depends_on": ["task_001"],
            "config": {
                "environment": "green",
                "replicas": 3
            }
        },
        {
            "id": "task_003",
            "type": "health_check",
            "depends_on": ["task_002"],
            "config": {
                "target": "green"
            }
        },
        {
            "id": "task_004",
            "type": "switch_traffic",
            "depends_on": ["task_003"],
            "config": {
                "from": "blue",
                "to": "green"
            }
        },
        {
            "id": "task_005",
            "type": "cleanup_blue",
            "depends_on": ["task_004"],
            "config": {
                "delay": 300
            }
        }
    ]
}
```

## 认证与授权

所有 API 请求需要在 HTTP 头中携带认证令牌:

```
Authorization: Bearer <token>
```

## 错误处理

参见 [核心数据结构文档](./core-models.md) 中的错误响应定义。
