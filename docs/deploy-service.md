# 部署执行服务

## 概述

部署执行服务负责:
1. 接收部署任务并执行
2. 整合物理机和 Kubernetes 部署

## 架构设计

```
TaskScheduler → Deploy Service → Kubernetes Deployer / Physical Deployer
```

## HTTP API

部署服务主要通过内部接口与 WorkflowManager 交互，不直接对外暴露 HTTP API。但为了管理和监控需要，提供以下端点:

### 获取部署信息

**Endpoint**: `GET /internal/v1/deploy/info/{deployment_id}`

**响应体**:
```json
{
    "deployment_id": "dep_123456",
    "status": "running",
    "replicas": 3,
    "ready_replicas": 2,
    "updated_at": "2024-10-24T10:05:00Z",
    "details": {
        "pods": [
            {
                "name": "app-abc123",
                "phase": "Running",
                "ready": true
            },
            {
                "name": "app-def456",
                "phase": "Running",
                "ready": true
            },
            {
                "name": "app-ghi789",
                "phase": "Pending",
                "ready": false
            }
        ]
    }
}
```

### 健康检查

**Endpoint**: `GET /internal/v1/deploy/health/{deployment_id}`

**响应体**:
```json
{
    "healthy": true,
    "message": "All replicas are healthy",
    "checks": [
        {
            "name": "pod_status",
            "status": "healthy",
            "message": "3/3 pods running"
        },
        {
            "name": "http_endpoint",
            "status": "healthy",
            "message": "HTTP 200 OK"
        },
        {
            "name": "resource_usage",
            "status": "healthy",
            "message": "CPU: 45%, Memory: 60%"
        }
    ]
}
```

### 获取日志

**Endpoint**: `GET /internal/v1/deploy/logs/{deployment_id}`

**查询参数**:
- `lines` (可选): 日志行数，默认 100
- `pod_name` (可选): 指定 Pod 名称

**响应体**:
```json
{
    "logs": [
        {
            "timestamp": "2024-10-24T10:05:00Z",
            "pod_name": "app-abc123",
            "container": "app",
            "message": "Server started on port 8080"
        },
        {
            "timestamp": "2024-10-24T10:05:01Z",
            "pod_name": "app-abc123",
            "container": "app",
            "message": "Connected to database"
        }
    ]
}
```

## Go 接口设计

### Deployer 基础接口

```go
type Deployer interface {
    Deploy(ctx context.Context, req *DeployRequest) (*DeployResult, error)
    
    Rollback(ctx context.Context, req *RollbackDeployRequest) (*DeployResult, error)
    
    GetDeploymentInfo(ctx context.Context, deploymentID string) (*DeploymentInfo, error)
    
    HealthCheck(ctx context.Context, deploymentID string) (*HealthCheckResult, error)
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
    DeploymentID  string
    Status        string
    Replicas      int32
    ReadyReplicas int32
    UpdatedAt     time.Time
    Details       map[string]interface{}
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
```

### Kubernetes Deployer 接口

```go
type KubernetesDeployer interface {
    Deployer
    
    ApplyManifest(ctx context.Context, namespace string, manifest []byte) error
    
    GetPodStatus(ctx context.Context, namespace, selector string) ([]*PodStatus, error)
    
    ScaleDeployment(ctx context.Context, namespace, name string, replicas int32) error
    
    GetLogs(ctx context.Context, namespace, podName string, lines int) (string, error)
    
    DeleteResources(ctx context.Context, namespace string, labels map[string]string) error
}

type PodStatus struct {
    Name      string
    Phase     string
    Ready     bool
    Restarts  int32
    NodeName  string
    CreatedAt time.Time
}
```

### Physical Deployer 接口

```go
type PhysicalDeployer interface {
    Deployer
    
    UploadArtifact(ctx context.Context, hosts []string, artifact *Artifact) error
    
    ExecuteCommand(ctx context.Context, hosts []string, command string) (*CommandResult, error)
    
    RestartService(ctx context.Context, hosts []string, serviceName string) error
    
    CheckServiceStatus(ctx context.Context, hosts []string, serviceName string) ([]*ServiceStatus, error)
    
    CleanupOldVersions(ctx context.Context, hosts []string, keepVersions int) error
}

type Artifact struct {
    Name        string
    Version     string
    Path        string
    Size        int64
    Checksum    string
    ContentType string
}

type CommandResult struct {
    Host     string
    Success  bool
    Output   string
    Error    string
    ExitCode int
}

type ServiceStatus struct {
    Host      string
    Running   bool
    Status    string
    PID       int
    Uptime    string
    Version   string
}
```

## Kubernetes 部署实现

### 部署流程

1. **构建 Manifest**
   ```go
   manifest := buildKubernetesManifest(deployment, environment)
   ```

2. **应用 Manifest**
   ```go
   err := k8sDeployer.ApplyManifest(ctx, namespace, manifest)
   ```

3. **等待 Pod 就绪**
   ```go
   err := waitForPodsReady(ctx, namespace, selector, timeout)
   ```

4. **健康检查**
   ```go
   result, err := k8sDeployer.HealthCheck(ctx, deploymentID)
   ```

### Manifest 示例

#### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-service
  namespace: production
  labels:
    app: api-service
    version: v1.0.0
    deployment-id: dep_123456
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  selector:
    matchLabels:
      app: api-service
  template:
    metadata:
      labels:
        app: api-service
        version: v1.0.0
    spec:
      containers:
      - name: api-service
        image: org/api-service:v1.0.0
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: ENVIRONMENT
          value: production
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

#### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: api-service
  namespace: production
spec:
  type: ClusterIP
  selector:
    app: api-service
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
    name: http
```

### 滚动更新策略

```go
type RollingUpdateConfig struct {
    MaxSurge       int32
    MaxUnavailable int32
    MinReadySeconds int32
    ProgressDeadlineSeconds int32
}
```

### 蓝绿部署策略

```go
type BlueGreenConfig struct {
    GreenNamespace  string
    BlueNamespace   string
    ServiceSelector map[string]string
    VerificationTime int
}
```

## 物理机部署实现

### 部署流程

1. **上传制品**
   ```go
   err := physicalDeployer.UploadArtifact(ctx, hosts, artifact)
   ```

2. **备份当前版本**
   ```go
   backupCmd := "cp -r /opt/app /opt/app.backup"
   physicalDeployer.ExecuteCommand(ctx, hosts, backupCmd)
   ```

3. **解压新版本**
   ```go
   extractCmd := "tar -xzf /tmp/app-v1.0.0.tar.gz -C /opt/app"
   physicalDeployer.ExecuteCommand(ctx, hosts, extractCmd)
   ```

4. **重启服务**
   ```go
   err := physicalDeployer.RestartService(ctx, hosts, "api-service")
   ```

5. **验证服务**
   ```go
   status, err := physicalDeployer.CheckServiceStatus(ctx, hosts, "api-service")
   ```

### 部署脚本示例

#### 部署脚本 (deploy.sh)

```bash
#!/bin/bash

set -e

APP_NAME="api-service"
VERSION="$1"
DEPLOY_DIR="/opt/applications/${APP_NAME}"
ARTIFACT_PATH="/tmp/${APP_NAME}-${VERSION}.tar.gz"
BACKUP_DIR="/opt/backups/${APP_NAME}"

echo "Deploying ${APP_NAME} version ${VERSION}"

echo "Creating backup..."
mkdir -p ${BACKUP_DIR}
if [ -d ${DEPLOY_DIR} ]; then
    tar -czf ${BACKUP_DIR}/backup-$(date +%Y%m%d-%H%M%S).tar.gz -C ${DEPLOY_DIR} .
fi

echo "Extracting artifact..."
mkdir -p ${DEPLOY_DIR}
tar -xzf ${ARTIFACT_PATH} -C ${DEPLOY_DIR}

echo "Stopping service..."
systemctl stop ${APP_NAME} || true

echo "Starting service..."
systemctl start ${APP_NAME}

echo "Checking service status..."
sleep 5
systemctl status ${APP_NAME}

echo "Verifying health endpoint..."
curl -f http://localhost:8080/health || (echo "Health check failed" && exit 1)

echo "Deployment completed successfully"
```

#### Systemd 服务配置 (api-service.service)

```ini
[Unit]
Description=API Service
After=network.target

[Service]
Type=simple
User=app
Group=app
WorkingDirectory=/opt/applications/api-service
ExecStart=/opt/applications/api-service/bin/server
Restart=on-failure
RestartSec=5s

Environment="ENVIRONMENT=production"
Environment="PORT=8080"

StandardOutput=journal
StandardError=journal
SyslogIdentifier=api-service

[Install]
WantedBy=multi-user.target
```

### SSH 连接配置

```go
type SSHConfig struct {
    Host           string
    Port           int
    User           string
    PrivateKey     string
    Timeout        time.Duration
    MaxRetries     int
}
```

### 并发执行

```go
func (p *PhysicalDeployer) ExecuteCommandConcurrently(
    ctx context.Context,
    hosts []string,
    command string,
) (map[string]*CommandResult, error) {
    results := make(map[string]*CommandResult)
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    for _, host := range hosts {
        wg.Add(1)
        go func(h string) {
            defer wg.Done()
            result, err := p.executeOnHost(ctx, h, command)
            mu.Lock()
            results[h] = result
            mu.Unlock()
        }(host)
    }
    
    wg.Wait()
    return results, nil
}
```

## 回滚机制

### Kubernetes 回滚

```go
func (k *KubernetesDeployer) Rollback(
    ctx context.Context,
    req *RollbackDeployRequest,
) (*DeployResult, error) {
    revision := getPreviousRevision(req.TargetDeploymentID)
    
    cmd := fmt.Sprintf("kubectl rollout undo deployment/%s --to-revision=%d",
        req.DeploymentID, revision)
    
    err := executeKubectl(ctx, cmd)
    if err != nil {
        return nil, err
    }
    
    return &DeployResult{
        Success: true,
        Message: fmt.Sprintf("Rolled back to revision %d", revision),
    }, nil
}
```

### 物理机回滚

```go
func (p *PhysicalDeployer) Rollback(
    ctx context.Context,
    req *RollbackDeployRequest,
) (*DeployResult, error) {
    hosts := getHostsFromEnvironment(req.Environment)
    
    rollbackCmd := `
        systemctl stop api-service
        rm -rf /opt/app
        mv /opt/app.backup /opt/app
        systemctl start api-service
    `
    
    results, err := p.ExecuteCommand(ctx, hosts, rollbackCmd)
    if err != nil {
        return nil, err
    }
    
    return &DeployResult{
        Success: true,
        Message: "Rollback completed",
        Details: map[string]interface{}{
            "results": results,
        },
    }, nil
}
```

## 健康检查

### HTTP 健康检查

```go
type HTTPHealthCheck struct {
    Endpoint       string
    ExpectedStatus int
    Timeout        time.Duration
    Retries        int
}

func (h *HTTPHealthCheck) Check(ctx context.Context) (*HealthCheck, error) {
    resp, err := http.Get(h.Endpoint)
    if err != nil {
        return &HealthCheck{
            Name:    "http_check",
            Status:  "unhealthy",
            Message: err.Error(),
        }, nil
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == h.ExpectedStatus {
        return &HealthCheck{
            Name:    "http_check",
            Status:  "healthy",
            Message: fmt.Sprintf("HTTP %d", resp.StatusCode),
        }, nil
    }
    
    return &HealthCheck{
        Name:    "http_check",
        Status:  "unhealthy",
        Message: fmt.Sprintf("Expected %d, got %d", h.ExpectedStatus, resp.StatusCode),
    }, nil
}
```

### TCP 健康检查

```go
type TCPHealthCheck struct {
    Host    string
    Port    int
    Timeout time.Duration
}

func (t *TCPHealthCheck) Check(ctx context.Context) (*HealthCheck, error) {
    addr := fmt.Sprintf("%s:%d", t.Host, t.Port)
    conn, err := net.DialTimeout("tcp", addr, t.Timeout)
    if err != nil {
        return &HealthCheck{
            Name:    "tcp_check",
            Status:  "unhealthy",
            Message: err.Error(),
        }, nil
    }
    defer conn.Close()
    
    return &HealthCheck{
        Name:    "tcp_check",
        Status:  "healthy",
        Message: "TCP connection successful",
    }, nil
}
```

### 进程健康检查

```go
type ProcessHealthCheck struct {
    ProcessName string
}

func (p *ProcessHealthCheck) Check(ctx context.Context, host string) (*HealthCheck, error) {
    cmd := fmt.Sprintf("pgrep -x %s", p.ProcessName)
    result, err := executeRemoteCommand(ctx, host, cmd)
    
    if err != nil || result.ExitCode != 0 {
        return &HealthCheck{
            Name:    "process_check",
            Status:  "unhealthy",
            Message: fmt.Sprintf("Process %s not running", p.ProcessName),
        }, nil
    }
    
    return &HealthCheck{
        Name:    "process_check",
        Status:  "healthy",
        Message: fmt.Sprintf("Process %s running (PID: %s)", p.ProcessName, result.Output),
    }, nil
}
```

## 日志收集

### Kubernetes 日志

```go
func (k *KubernetesDeployer) GetLogs(
    ctx context.Context,
    namespace, podName string,
    lines int,
) (string, error) {
    opts := &corev1.PodLogOptions{
        TailLines: int64Ptr(int64(lines)),
    }
    
    req := k.clientset.CoreV1().Pods(namespace).GetLogs(podName, opts)
    logs, err := req.Stream(ctx)
    if err != nil {
        return "", err
    }
    defer logs.Close()
    
    buf := new(bytes.Buffer)
    _, err = io.Copy(buf, logs)
    if err != nil {
        return "", err
    }
    
    return buf.String(), nil
}
```

### 物理机日志

```go
func (p *PhysicalDeployer) GetLogs(
    ctx context.Context,
    host string,
    logPath string,
    lines int,
) (string, error) {
    cmd := fmt.Sprintf("tail -n %d %s", lines, logPath)
    result, err := p.ExecuteCommand(ctx, []string{host}, cmd)
    if err != nil {
        return "", err
    }
    
    return result.Output, nil
}
```

## 监控指标

- 部署成功率
- 部署耗时
- 回滚次数
- 健康检查成功率
- Pod/服务可用性
- 资源使用情况

## 错误处理

### 部署失败处理

1. **记录详细错误信息**
2. **自动回滚到上一个稳定版本**
3. **通知相关人员**
4. **保留现场用于问题排查**

### 超时处理

```go
type DeployTimeout struct {
    BuildTimeout      time.Duration
    DeployTimeout     time.Duration
    HealthCheckTimeout time.Duration
}
```

### 重试机制

```go
type RetryConfig struct {
    MaxRetries     int
    InitialDelay   time.Duration
    MaxDelay       time.Duration
    BackoffFactor  float64
}
```
