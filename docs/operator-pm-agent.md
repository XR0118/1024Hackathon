# Operator PM Agent

Operator PM Agent 是 Boreas 持续部署平台中运行在物理机节点上的代理组件，负责执行具体的部署操作。

## 功能特性

- 支持 Docker 容器部署
- 支持二进制应用部署
- 支持脚本应用部署
- 提供应用状态监控
- 提供健康检查接口

## 核心接口

### 1. 应用部署接口

```http
POST /v1/apply
Content-Type: application/json

{
  "app": "my-app",
  "version": "v1.0.0",
  "pkg": {
    "type": "docker",
    "image": "my-app:v1.0.0",
    "command": ["/app/my-app"],
    "args": ["--config", "/app/config.yaml"],
    "environment": {
      "ENV": "production"
    },
    "volumes": [
      {
        "host_path": "/data",
        "container_path": "/app/data",
        "read_only": false
      }
    ],
    "ports": [
      {
        "host_port": 8080,
        "container_port": 8080,
        "protocol": "tcp"
      }
    ]
  }
}
```

### 2. 获取所有应用状态

```http
GET /v1/status
```

响应：
```json
{
  "apps": [
    {
      "app": "my-app",
      "version": "v1.0.0",
      "healthy": {
        "level": 100,
        "msg": "Running"
      },
      "status": "running",
      "updated": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 3. 获取指定应用状态

```http
GET /v1/status/my-app
```

响应：
```json
{
  "app": "my-app",
  "version": "v1.0.0",
  "healthy": {
    "level": 100
  }
}
```

### 4. 健康检查

```http
GET /v1/health
```

响应：
```json
{
  "status": "healthy",
  "service": "operator-pm-agent"
}
```

## 部署方式

### 使用 Docker

```bash
# 构建镜像
docker build -f deployments/docker/operator-pm-agent.Dockerfile -t boreas/operator-pm-agent .

# 运行容器
docker run -d \
  --name boreas-agent \
  -p 8081:8081 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /var/lib/boreas-agent:/var/lib/boreas-agent \
  boreas/operator-pm-agent
```

### 直接运行

```bash
# 编译
go build -o operator-pm-agent ./cmd/operator-pm-agent

# 运行
./operator-pm-agent --work-dir /var/lib/boreas-agent --port 8081
```

## 配置

Agent 支持以下命令行参数：

- `--work-dir`: 工作目录，默认为 `/var/lib/boreas-agent`
- `--port`: 服务端口，默认为 `8081`
- `--host`: 服务地址，默认为 `0.0.0.0`
- `--version`: 显示版本信息

## 与 Operator-PM 集成

1. 在 Operator-PM 中注册 Agent：

```bash
curl -X POST http://operator-pm:8080/api/v1/agents \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "pm-node-1",
    "agent_url": "http://pm-node-1:8081"
  }'
```

2. 在环境配置中指定 Agent ID：

```json
{
  "name": "production-pm",
  "type": "physical",
  "config": {
    "agent_id": "pm-node-1"
  }
}
```

## 监控和日志

Agent 会将应用状态持久化到工作目录中，支持：

- 应用状态恢复
- 部署日志记录
- 健康状态监控

## 安全注意事项

- Agent 需要访问 Docker Socket 来管理容器
- 建议在受信任的网络环境中运行
- 可以通过防火墙限制访问端口
- 支持 TLS 加密（待实现）
