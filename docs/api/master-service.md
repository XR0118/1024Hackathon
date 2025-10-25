# Master Service API 文档

Master Service 是 Boreas 平台的核心服务，负责版本管理、应用管理、环境管理和部署编排。

## 基础信息

- **服务名称**: Master Service
- **端口**: 8080
- **基础路径**: `/api/v1`

## 健康检查

### GET /health
检查服务健康状态

**响应示例**:
```json
{
  "status": "healthy",
  "service": "master-service",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### GET /ready
检查服务就绪状态

**响应示例**:
```json
{
  "status": "ready",
  "service": "master-service",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 版本管理

### POST /api/v1/versions
创建新版本

**请求体**:
```json
{
  "git_tag": "v1.0.0",
  "git_commit": "abc123def456",
  "repository": "https://github.com/example/repo",
  "description": "Initial release",
  "app_build": [
    {
      "app_id": "app-uuid",
      "app_name": "api-service",
      "docker_image": "registry.example.com/api-service:v1.0.0"
    }
  ]
}
```

**响应示例**:
```json
{
  "id": "version-uuid",
  "git_tag": "v1.0.0",
  "git_commit": "abc123def456",
  "repository": "https://github.com/example/repo",
  "created_by": "user@example.com",
  "created_at": "2024-01-01T00:00:00Z",
  "description": "Initial release",
  "app_build": [
    {
      "app_id": "app-uuid",
      "app_name": "api-service",
      "docker_image": "registry.example.com/api-service:v1.0.0"
    }
  ]
}
```

### GET /api/v1/versions
获取版本列表

**查询参数**:
- `repository` (string): 仓库地址过滤
- `page` (int): 页码，默认1
- `page_size` (int): 每页大小，默认20

**响应示例**:
```json
{
  "versions": [
    {
      "id": "version-uuid",
      "git_tag": "v1.0.0",
      "git_commit": "abc123def456",
      "repository": "https://github.com/example/repo",
      "created_by": "user@example.com",
      "created_at": "2024-01-01T00:00:00Z",
      "description": "Initial release"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 20
}
```

### GET /api/v1/versions/{id}
获取版本详情

**路径参数**:
- `id` (string): 版本ID

**响应示例**:
```json
{
  "id": "version-uuid",
  "git_tag": "v1.0.0",
  "git_commit": "abc123def456",
  "repository": "https://github.com/example/repo",
  "created_by": "user@example.com",
  "created_at": "2024-01-01T00:00:00Z",
  "description": "Initial release"
}
```

### DELETE /api/v1/versions/{id}
删除版本

**路径参数**:
- `id` (string): 版本ID

**响应示例**:
```json
{
  "message": "Version deleted successfully"
}
```

## 应用管理

### POST /api/v1/applications
创建应用

**请求体**:
```json
{
  "name": "my-app",
  "repository": "https://github.com/example/repo",
  "type": "microservice",
  "config": {
    "port": "8080",
    "health_check": "/health"
  }
}
```

### GET /api/v1/applications
获取应用列表

**查询参数**:
- `repository` (string): 仓库地址过滤
- `type` (string): 应用类型过滤
- `page` (int): 页码，默认1
- `page_size` (int): 每页大小，默认20

### GET /api/v1/applications/{id}
获取应用详情

### PUT /api/v1/applications/{id}
更新应用

### DELETE /api/v1/applications/{id}
删除应用

## 环境管理

### POST /api/v1/environments
创建环境

**请求体**:
```json
{
  "name": "production",
  "type": "kubernetes",
  "config": {
    "namespace": "production",
    "cluster": "prod-cluster"
  },
  "is_active": true
}
```

### GET /api/v1/environments
获取环境列表

### GET /api/v1/environments/{id}
获取环境详情

### PUT /api/v1/environments/{id}
更新环境

### DELETE /api/v1/environments/{id}
删除环境

## 部署管理

### POST /api/v1/deployments
创建部署

**请求体**:
```json
{
  "version_id": "version-uuid",
  "must_in_order": ["app-uuid-1", "app-uuid-2"],
  "environment_id": "env-uuid",
  "manual_approval": false,
  "strategy": [
    {
      "batch_size": 1,
      "batch_interval": 10,
      "canary_ratio": 0.1,
      "auto_rollback": true,
      "manual_approval_status": null
    }
  ]
}
```

### GET /api/v1/deployments
获取部署列表

### GET /api/v1/deployments/{id}
获取部署详情

### POST /api/v1/deployments/{id}/cancel
取消部署

### POST /api/v1/deployments/{id}/rollback
回滚部署

## 任务管理

### GET /api/v1/tasks
获取任务列表

### GET /api/v1/tasks/{id}
获取任务详情

### POST /api/v1/tasks/{id}/retry
重试任务

## Webhook

### POST /api/v1/webhooks/github
GitHub Webhook 处理

**请求体**: GitHub Webhook 事件

**响应示例**:
```json
{
  "message": "Webhook processed successfully",
  "version_created": true,
  "version_id": "version-uuid",
  "deployments_triggered": [
    {
      "deployment_id": "deployment-uuid",
      "environment_id": "env-uuid",
      "status": "pending"
    }
  ]
}
```

## 错误响应

所有错误响应都遵循以下格式：

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description",
    "details": {
      "field": "Additional error details"
    }
  }
}
```

## 状态码

- `200` - 成功
- `201` - 创建成功
- `400` - 请求错误
- `404` - 资源不存在
- `500` - 服务器内部错误
