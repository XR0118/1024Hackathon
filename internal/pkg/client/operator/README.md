# Operator Client 包

## 概述

这个包提供了与各类 Operator（K8s、PM、Mock）通信的统一客户端实现。

## 目录结构

```
operator/
├── README.md              # 本文档
├── manager.go            # Operator Manager 管理器
├── factory.go            # Operator Factory 工厂
├── k8s_client.go         # K8s Operator 客户端
├── pm_client.go          # PM Operator 客户端
└── mock_client.go        # Mock Operator 客户端
```

## 核心组件

### 1. Operator Manager

管理所有环境的 Operator 客户端。

**主要功能：**
- 维护环境 ID 到 Operator 的映射
- 提供统一的 Operator 访问接口
- 支持并发安全访问

**使用示例：**

```go
import "github.com/boreas/internal/pkg/client/operator"

// 创建 manager
manager := operator.NewManager()

// 注册 operator
manager.RegisterOperator("env-001", k8sOperator)
manager.RegisterOperator("env-002", pmOperator)

// 使用 operator
status, err := manager.GetApplicationStatus(ctx, "env-001", "my-app")
```

### 2. Operator Factory

根据环境类型自动创建对应的 Operator 客户端。

**配置：**

```go
config := &operator.Config{
    K8SOperatorURL: "http://k8s-operator:8081",
    PMOperatorURL:  "http://pm-operator:8082",
    UseMock:        false,
}
```

**初始化所有环境的 Operator：**

```go
// 从数据库查询所有环境
environments, _, err := envRepo.List(ctx, nil)

// 初始化 operators
manager, err := operator.InitializeOperators(environments, config)
```

**手动创建单个 Operator：**

```go
// 根据环境对象创建
op, err := operator.CreateOperatorFromEnvironment(env, config)

// 根据类型创建
op, err := operator.CreateOperatorByType("kubernetes", "http://k8s-operator:8081", false)
```

### 3. K8s Operator Client

与 Kubernetes Operator 通信的客户端。

**接口端点：**
- `POST /api/apply` - 应用部署
- `GET /api/status/:appName` - 查询应用状态
- `GET /health` - 健康检查

**使用示例：**

```go
client := operator.NewK8sClient("http://k8s-operator:8081")

// 应用部署
resp, err := client.Apply(ctx, &models.ApplyDeploymentRequest{
    App:     "my-app",
    Version: "v1.2.3",
    // ...
})

// 查询状态
status, err := client.GetApplicationStatus(ctx, "my-app")
```

### 4. PM Operator Client

与物理机 Operator 通信的客户端。

**接口端点：**
- `POST /api/deploy` - 应用部署
- `GET /api/app/:appName` - 查询应用状态
- `GET /health` - 健康检查

**使用示例：**

```go
client := operator.NewPMClient("http://pm-operator:8082")

// 应用部署
resp, err := client.Apply(ctx, &models.ApplyDeploymentRequest{
    App:     "my-app",
    Version: "v1.2.3",
    // ...
})

// 查询状态
status, err := client.GetApplicationStatus(ctx, "my-app")
```

### 5. Mock Operator Client

用于开发和测试的 Mock 客户端。

**特点：**
- 返回模拟数据
- 不依赖真实的 Operator 服务
- 适合单元测试和本地开发

**使用示例：**

```go
client := operator.NewMockClient()

// 返回模拟的部署响应
resp, err := client.Apply(ctx, &models.ApplyDeploymentRequest{
    App:     "my-app",
    Version: "v1.2.3",
})

// 返回模拟的应用状态
status, err := client.GetApplicationStatus(ctx, "my-app")
```

## 配置

### 配置文件

在 `cmd/master-service/configs/master.yaml` 中配置：

```yaml
operator:
  k8s_operator_url: "http://localhost:8081"
  pm_operator_url: "http://localhost:8082"
  use_mock: false  # 设置为 true 使用 Mock Operator
```

### 环境变量

```bash
# K8s Operator URL
export MASTER_OPERATOR_K8S_OPERATOR_URL="http://k8s-operator:8081"

# PM Operator URL
export MASTER_OPERATOR_PM_OPERATOR_URL="http://pm-operator:8082"

# 是否使用 Mock
export MASTER_OPERATOR_USE_MOCK="false"
```

## 开发和测试

### 使用 Mock Operator

开发环境中，建议使用 Mock Operator：

```yaml
operator:
  use_mock: true
```

或者：

```bash
export MASTER_OPERATOR_USE_MOCK="true"
```

### 单元测试

```go
func TestApplicationService(t *testing.T) {
    // 创建 mock operator
    mockOp := operator.NewMockClient()
    
    // 创建 manager
    manager := operator.NewManager()
    manager.RegisterOperator("test-env", mockOp)
    
    // 创建 service
    svc := service.NewApplicationService(
        appRepo,
        versionRepo,
        deploymentRepo,
        manager,
    )
    
    // 测试
    result, err := svc.GetApplicationVersionsDetail(ctx, "test-app")
    assert.NoError(t, err)
}
```

## 接口规范

### Operator 统一接口

所有 Operator 客户端必须实现 `interfaces.Operator` 接口：

```go
type Operator interface {
    // Apply 应用部署
    Apply(ctx context.Context, req *models.ApplyDeploymentRequest) (*models.ApplyDeploymentResponse, error)
    
    // GetApplicationStatus 获取应用状态
    GetApplicationStatus(ctx context.Context, appName string) (*models.ApplicationStatusResponse, error)
    
    // HealthCheck 健康检查
    HealthCheck(ctx context.Context) error
    
    // GetType 获取 Operator 类型
    GetType() string
}
```

### 请求/响应模型

**ApplyDeploymentRequest：**

```go
type ApplyDeploymentRequest struct {
    App         string             `json:"app"`
    Version     string             `json:"version"`
    Deployments []VersionDeployment `json:"deployments"`
}
```

**ApplyDeploymentResponse：**

```go
type ApplyDeploymentResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}
```

**ApplicationStatusResponse：**

```go
type ApplicationStatusResponse struct {
    App      string          `json:"app"`
    Healthy  HealthInfo      `json:"healthy"`
    Versions []VersionStatus `json:"versions"`
}

type VersionStatus struct {
    Version string       `json:"version"`
    Percent float64      `json:"percent"`
    Healthy HealthInfo   `json:"healthy"`
    Nodes   []NodeStatus `json:"nodes"`
}

type NodeStatus struct {
    Node    string     `json:"node"`
    Healthy HealthInfo `json:"healthy"`
}

type HealthInfo struct {
    Level int `json:"level"` // 0-100
}
```

## 错误处理

### 超时

建议为每个请求设置超时：

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

status, err := operator.GetApplicationStatus(ctx, "my-app")
```

### 重试

对于网络错误，可以考虑重试：

```go
var status *models.ApplicationStatusResponse
var err error

for i := 0; i < 3; i++ {
    status, err = operator.GetApplicationStatus(ctx, "my-app")
    if err == nil {
        break
    }
    time.Sleep(time.Second * time.Duration(i+1))
}
```

### 降级

当 Operator 不可用时，可以返回默认状态：

```go
status, err := operator.GetApplicationStatus(ctx, "my-app")
if err != nil {
    // 返回默认状态
    return &models.ApplicationStatusResponse{
        App:      "my-app",
        Healthy:  models.HealthInfo{Level: 0},
        Versions: []models.VersionStatus{},
    }, nil
}
```

## 最佳实践

1. **使用连接池**：HTTP 客户端使用连接池提高性能
2. **设置超时**：为每个请求设置合理的超时时间
3. **错误处理**：优雅处理网络错误和超时
4. **日志记录**：记录关键操作和错误信息
5. **监控指标**：收集 Operator 调用的监控数据
6. **健康检查**：定期检查 Operator 的健康状态

## 扩展

### 添加新的 Operator 类型

1. 实现 `interfaces.Operator` 接口：

```go
type NewOperatorClient struct {
    baseURL string
    client  *http.Client
}

func (c *NewOperatorClient) Apply(ctx context.Context, req *models.ApplyDeploymentRequest) (*models.ApplyDeploymentResponse, error) {
    // 实现部署逻辑
}

func (c *NewOperatorClient) GetApplicationStatus(ctx context.Context, appName string) (*models.ApplicationStatusResponse, error) {
    // 实现状态查询逻辑
}

func (c *NewOperatorClient) HealthCheck(ctx context.Context) error {
    // 实现健康检查逻辑
}

func (c *NewOperatorClient) GetType() string {
    return "new-operator"
}
```

2. 在 Factory 中添加支持：

```go
func CreateOperatorFromEnvironment(env *models.Environment, config *Config) (interfaces.Operator, error) {
    switch env.Type {
    case "new-type":
        return NewNewOperatorClient(config.NewOperatorURL), nil
    // ... 其他类型
    }
}
```

## 相关文档

- [Operator Manager 集成文档](../../../docs/operator-manager-integration.md)
- [Operator PM 设计](../../../docs/operator-pm-design.md)
- [部署工作流设计](../../../docs/deployment-workflow-design.md)
