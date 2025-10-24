# GitHub Webhook 回调服务

## 概述

GitHub Webhook 回调服务负责处理 GitHub webhook 回调事件，自动生成 Version 和对应的 Deployment。

## 架构设计

```
GitHub → Webhook Endpoint → DeploymentManager → Version/Deployment 创建
```

## HTTP API

### 接收 GitHub Webhook

接收 GitHub 的 webhook 回调并处理版本创建和部署触发。

**Endpoint**: `POST /api/v1/webhooks/github`

**请求头**:
```
X-GitHub-Event: push | pull_request | release
X-Hub-Signature-256: sha256=<signature>
X-GitHub-Delivery: <delivery-id>
Content-Type: application/json
```

### Push Event (Tag 创建)

当新 Tag 被推送时自动创建版本。

**请求体**:
```json
{
    "ref": "refs/tags/v1.0.0",
    "repository": {
        "full_name": "org/repo",
        "clone_url": "https://github.com/org/repo.git",
        "default_branch": "main"
    },
    "head_commit": {
        "id": "abc123def456",
        "message": "Release v1.0.0",
        "timestamp": "2024-10-24T10:00:00Z",
        "author": {
            "name": "John Doe",
            "email": "john@example.com"
        }
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
    "message": "Tag webhook processed successfully",
    "version_created": true,
    "version_id": "ver_123456",
    "deployments_triggered": [
        {
            "deployment_id": "dep_123456",
            "environment_id": "env_prod",
            "status": "pending"
        }
    ]
}
```

### Pull Request Event (PR 合并)

当 PR 被合并到主分支时自动创建版本。

**请求体**:
```json
{
    "action": "closed",
    "number": 123,
    "pull_request": {
        "merged": true,
        "merge_commit_sha": "abc123def456",
        "merged_at": "2024-10-24T10:00:00Z",
        "merged_by": {
            "login": "username",
            "email": "user@example.com"
        },
        "base": {
            "ref": "main",
            "sha": "def456abc789"
        },
        "head": {
            "ref": "feature/new-feature",
            "sha": "abc123def456"
        },
        "title": "Add new feature",
        "body": "This PR adds a new feature\n\nCloses #100"
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
    "message": "Pull request merge processed successfully",
    "version_created": true,
    "version_id": "ver_789012",
    "auto_tag": "v1.0.1-pr123",
    "deployments_triggered": []
}
```

### Release Event

当 GitHub Release 创建时处理版本同步。

**请求体**:
```json
{
    "action": "published",
    "release": {
        "tag_name": "v1.0.0",
        "name": "Version 1.0.0",
        "body": "Release notes...",
        "draft": false,
        "prerelease": false,
        "created_at": "2024-10-24T10:00:00Z",
        "published_at": "2024-10-24T10:05:00Z",
        "author": {
            "login": "username",
            "email": "user@example.com"
        },
        "target_commitish": "main"
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
    "message": "Release webhook processed successfully",
    "version_created": true,
    "version_id": "ver_123456",
    "deployments_triggered": [
        {
            "deployment_id": "dep_789012",
            "environment_id": "env_prod",
            "status": "pending"
        }
    ]
}
```

## Go 接口设计

### WebhookService 接口

```go
type WebhookService interface {
    HandleGitHubWebhook(ctx context.Context, event string, payload []byte) (*WebhookResponse, error)
    
    VerifySignature(payload []byte, signature string) error
}

type WebhookResponse struct {
    Message              string                 `json:"message"`
    VersionCreated       bool                   `json:"version_created,omitempty"`
    VersionID            string                 `json:"version_id,omitempty"`
    AutoTag              string                 `json:"auto_tag,omitempty"`
    DeploymentsTriggered []DeploymentReference  `json:"deployments_triggered,omitempty"`
}

type DeploymentReference struct {
    DeploymentID  string           `json:"deployment_id"`
    EnvironmentID string           `json:"environment_id"`
    Status        DeploymentStatus `json:"status"`
}
```

### DeploymentManager 接口

```go
type DeploymentManager interface {
    ProcessWebhookEvent(ctx context.Context, event *GitHubEvent) (*ProcessResult, error)
    
    CreateVersionFromTag(ctx context.Context, tag *GitTag) (*Version, error)
    
    CreateVersionFromPR(ctx context.Context, pr *PullRequest) (*Version, error)
    
    CreateVersionFromRelease(ctx context.Context, release *Release) (*Version, error)
    
    TriggerAutoDeployments(ctx context.Context, versionID string) ([]*Deployment, error)
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
    Timestamp  time.Time
}

type PullRequest struct {
    Number       int
    Title        string
    Body         string
    MergeCommit  string
    Repository   string
    BaseBranch   string
    HeadBranch   string
    MergedBy     string
    MergedAt     time.Time
}

type Release struct {
    TagName     string
    Name        string
    Body        string
    Commit      string
    Repository  string
    Author      string
    PublishedAt time.Time
    IsPrerelease bool
}

type ProcessResult struct {
    Success              bool
    VersionCreated       bool
    Version              *Version
    DeploymentsTriggered []*Deployment
    Error                error
}
```

## 处理流程

### Tag 创建流程

1. **接收 Webhook**
   - GitHub 推送 Tag 触发 webhook
   - 系统接收 POST 请求到 `/api/v1/webhooks/github`

2. **验证签名**
   - 验证 `X-Hub-Signature-256` 签名
   - 确保请求来自 GitHub

3. **解析事件**
   - 解析 Push Event
   - 验证是否为 Tag（`refs/tags/*`）

4. **创建版本**
   - 提取 Tag 信息（名称、提交哈希等）
   - 创建 Version 记录
   - 关联代码仓库

5. **触发部署（可选）**
   - 检查是否配置自动部署
   - 根据配置触发相应环境的部署
   - 返回部署信息

### PR 合并流程

1. **接收 Webhook**
   - GitHub PR 合并触发 webhook
   - 系统接收 POST 请求

2. **验证签名**
   - 验证请求签名

3. **解析事件**
   - 解析 Pull Request Event
   - 验证 `action=closed` 且 `merged=true`

4. **生成版本**
   - 提取合并提交信息
   - 生成自动 Tag（如 `v1.0.1-pr123`）
   - 创建 Version 记录

5. **可选操作**
   - 自动部署到测试环境
   - 通知相关人员
   - 更新 Issue 状态

## 配置示例

### Webhook 配置

在 GitHub 仓库设置中配置 Webhook：

```
Payload URL: https://your-domain.com/api/v1/webhooks/github
Content type: application/json
Secret: <your-webhook-secret>
Events: Push events, Pull requests, Releases
```

### 自动部署配置

```json
{
    "auto_deployment": {
        "enabled": true,
        "rules": [
            {
                "event_type": "tag",
                "tag_pattern": "^v[0-9]+\\.[0-9]+\\.[0-9]+$",
                "environments": ["env_prod"],
                "applications": ["app_*"]
            },
            {
                "event_type": "pr_merge",
                "base_branch": "main",
                "environments": ["env_staging"],
                "applications": ["app_*"]
            },
            {
                "event_type": "release",
                "prerelease": false,
                "environments": ["env_prod"],
                "applications": ["app_*"]
            }
        ]
    }
}
```

## 错误处理

### 签名验证失败

```json
{
    "error": {
        "code": "INVALID_SIGNATURE",
        "message": "GitHub webhook signature verification failed"
    }
}
```

HTTP 状态码: `401 Unauthorized`

### 不支持的事件类型

```json
{
    "error": {
        "code": "UNSUPPORTED_EVENT",
        "message": "Event type 'issues' is not supported",
        "details": {
            "event_type": "issues"
        }
    }
}
```

HTTP 状态码: `400 Bad Request`

### 版本创建失败

```json
{
    "error": {
        "code": "VERSION_CREATION_FAILED",
        "message": "Failed to create version from tag",
        "details": {
            "tag": "v1.0.0",
            "reason": "Version with this tag already exists"
        }
    }
}
```

HTTP 状态码: `409 Conflict`

## 安全考虑

1. **签名验证**
   - 必须验证 `X-Hub-Signature-256` 签名
   - 使用配置的 Secret 进行 HMAC 验证

2. **IP 白名单**
   - 可选：限制只接受来自 GitHub IP 的请求
   - GitHub IP 范围: https://api.github.com/meta

3. **重放攻击防护**
   - 使用 `X-GitHub-Delivery` 作为幂等性标识
   - 短期内相同 Delivery ID 的请求应被忽略

4. **错误信息**
   - 不要在响应中泄露敏感信息
   - 详细错误记录到日志系统

## 监控指标

- Webhook 接收数量
- 签名验证成功/失败率
- 版本创建成功/失败率
- 自动部署触发数量
- 处理延迟时间
- 错误类型分布
