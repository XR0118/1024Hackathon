# Operator-PM 物理机管理组件概要设计

## 1. 概述

Operator-PM 是 Boreas 平台中专门负责物理机应用部署和管理的组件，采用主从架构设计，包含 Operator-PM（主控服务）和 Operator-PM-Agent（代理服务）两个核心组件。

## 2. 系统架构

### 2.1 整体架构

```
┌─────────────────┐    HTTP API    ┌─────────────────┐
│   Operator-PM   │◄──────────────►│ Operator-PM-    │
│   (主控服务)     │                │ Agent (代理服务) │
│                 │                │                 │
│ - 应用部署管理   │                │ - 应用运行管理   │
│ - 状态监控      │                │ - 健康检查      │
│ - 节点选择      │                │ - 资源管理      │
└─────────────────┘                └─────────────────┘
         │                                   │
         │                                   │
    ┌─────────┐                         ┌─────────┐
    │ 配置管理 │                         │ 物理机  │
    │ - 节点映射│                         │ 环境    │
    │ - IP映射 │                         │         │
    └─────────┘                         └─────────┘
```

### 2.2 组件关系

- **Operator-PM**: 作为主控服务，负责任务调度、节点选择、状态汇总
- **Operator-PM-Agent**: 作为代理服务，部署在每台物理机上，负责具体的应用运行管理
- **配置驱动**: 通过配置文件管理应用-节点映射和节点-IP映射关系

## 3. Operator-PM 主控服务

### 3.1 核心功能

1. **应用部署管理**
   - 根据配置选择目标节点
   - 支持多版本并行部署
   - 支持百分比流量控制

2. **状态监控**
   - 实时监控所有节点状态
   - 汇总应用健康状态
   - 提供统一的状态查询接口

3. **节点管理**
   - 维护节点-IP映射关系
   - 支持节点健康检查
   - 动态节点选择算法

### 3.2 API 接口

#### 3.2.1 应用部署接口

```http
POST /v1/apply
Content-Type: application/json

{
  "app": "my-app",
  "versions": [
    {
      "version": "v1.0.0",
      "percent": 0.5,
      "pkg": {
        "type": "docker",
        "image": "my-app:v1.0.0"
      }
    }
  ]
}
```

**响应:**
```json
{
  "app": "my-app",
  "message": "Deployed to 3/5 nodes",
  "success": true
}
```

#### 3.2.2 应用状态查询接口

```http
GET /v1/status/:app
```

**响应:**
```json
{
  "app": "my-app",
  "healthy": {"level": 100},
  "versions": [
    {
      "version": "v1.0.0",
      "percent": 0.5,
      "healthy": {"level": 100},
      "nodes": [
        {"node": "node-1", "healthy": {"level": 100}},
        {"node": "node-2", "healthy": {"level": 100}}
      ]
    }
  ]
}
```

#### 3.2.3 健康检查接口

```http
GET /v1/health
GET /v1/ready
```

### 3.3 配置管理

#### 3.3.1 主配置文件 (operator-pm.yaml)

```yaml
server:
  host: "0.0.0.0"
  port: 8080

log:
  level: "info"
  format: "json"
  output: "stdout"

pm:
  agent_timeout: 30
  max_retries: 3
  
  health_check:
    interval: 60
    timeout: 10
  
  deployment:
    timeout: 300
    max_concurrent: 5
    retry_interval: 30
    status_check: 30

  agent:
    port: 8081
    path: "/v1"

  config_paths:
    app_to_nodes: "./cmd/operator-pm/configs/app-to-nodes.yaml"
    node_to_ip: "./cmd/operator-pm/configs/node-to-ip.yaml"
```

#### 3.3.2 应用-节点映射 (app-to-nodes.yaml)

```yaml
app_to_nodes:
  my-app:
    - node-1
    - node-2
    - node-3
  web-service:
    - node-2
    - node-3
    - node-5
```

#### 3.3.3 节点-IP映射 (node-to-ip.yaml)

```yaml
node_to_ip:
  node-1: "192.168.1.10"
  node-2: "192.168.1.11"
  node-3: "192.168.1.12"
  node-4: "192.168.1.13"
  node-5: "192.168.1.14"
```

### 3.4 Agent 配置管理

#### 3.4.1 Agent 配置文件 (agent.yaml)

```yaml
server:
  host: "0.0.0.0"
  port: 8081

log:
  level: "info"
  format: "json"
  output: "stdout"

agent:
  id: "pm-node-1"
  hostname: "pm-node-1.example.com"
  work_dir: "/var/lib/boreas-agent"
  
  docker:
    enabled: true
    socket_path: "/var/run/docker.sock"
    registry: ""
    network_mode: "bridge"
  
  health:
    check_interval: 30
    timeout: 10
    retry_count: 3
  
  config:
    # 自定义配置项
    custom_setting: "value"
```

### 3.4 技术实现

- **语言**: Go
- **Web框架**: Gin
- **配置管理**: Viper + 独立YAML文件
- **HTTP客户端**: 标准库 net/http
- **无状态设计**: 不依赖数据库，完全基于配置

## 4. Operator-PM-Agent 代理服务

### 4.1 核心功能

1. **应用运行管理**
   - 支持Docker容器运行
   - 支持二进制文件运行
   - 支持脚本执行

2. **健康检查**
   - 定期检查应用健康状态
   - 支持自定义健康检查逻辑
   - 实时状态上报

3. **资源管理**
   - 监控系统资源使用情况
   - 支持资源限制配置
   - 自动清理异常进程

### 4.2 API 接口

#### 4.2.1 应用部署接口

```http
POST /v1/apply
Content-Type: application/json

{
  "app": "my-app",
  "version": "v1.0.0",
  "pkg": {
    "type": "docker",
    "image": "my-app:v1.0.0",
    "ports": [{"host": 8080, "container": 8080}],
    "environment": {"ENV": "production"}
  }
}
```

#### 4.2.2 状态查询接口

```http
GET /v1/status
GET /v1/status/:app
```

**响应:**
```json
{
  "apps": [
    {
      "app": "my-app",
      "version": "v1.0.0",
      "healthy": {"level": 100, "msg": "OK"},
      "status": "running",
      "updated": "2024-01-01T12:00:00Z"
    }
  ]
}
```

#### 4.2.3 健康检查接口

```http
GET /v1/health
```

### 4.3 应用运行器 (Runner)

支持多种运行模式：

1. **Docker模式**
   - 使用Docker API管理容器
   - 支持镜像拉取和更新
   - 支持端口映射和卷挂载

2. **二进制模式**
   - 直接运行可执行文件
   - 支持进程管理和监控
   - 支持环境变量配置

3. **脚本模式**
   - 执行Shell脚本
   - 支持自定义启动逻辑
   - 支持日志收集

### 4.4 技术实现

- **语言**: Go
- **Web框架**: Gin
- **容器管理**: Docker API
- **进程管理**: 标准库 os/exec
- **配置管理**: YAML文件
- **状态存储**: 内存 + 文件持久化

## 5. 部署架构

### 5.1 部署模式

```
┌─────────────────────────────────────────────────────────┐
│                    Operator-PM                          │
│                  (主控服务)                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   Config    │  │   Service   │  │   Handler   │     │
│  │   管理      │  │   业务逻辑   │  │   API接口   │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────────────────────────────────────────┘
         │                    │                    │
         │                    │                    │
    ┌─────────┐         ┌─────────┐         ┌─────────┐
    │ Node-1  │         │ Node-2  │         │ Node-3  │
    │ Agent   │         │ Agent   │         │ Agent   │
    └─────────┘         └─────────┘         └─────────┘
```

### 5.2 配置文件结构

```
cmd/operator-pm/configs/
├── operator-pm.yaml      # 主配置文件
├── app-to-nodes.yaml     # 应用->节点映射
└── node-to-ip.yaml       # 节点->IP地址映射

cmd/operator-pm-agent/configs/
├── agent.yaml            # Agent配置文件
└── agent-example.yaml    # Agent配置示例
```

## 6. 核心特性

### 6.1 高可用性

- **无单点故障**: 主控服务可多实例部署
- **节点容错**: 单个节点故障不影响整体服务
- **自动恢复**: 支持应用自动重启和故障转移

### 6.2 可扩展性

- **水平扩展**: 支持动态添加物理节点
- **配置驱动**: 通过配置文件管理节点和应用关系
- **模块化设计**: 各组件独立部署和升级

### 6.3 易用性

- **RESTful API**: 标准化的HTTP接口
- **配置简单**: 基于YAML的配置文件
- **监控友好**: 提供丰富的状态和健康检查接口

### 6.4 安全性

- **网络隔离**: 支持内网部署
- **权限控制**: 可配置访问控制
- **日志审计**: 完整的操作日志记录

## 7. 监控和运维

### 7.1 健康检查

- **服务级别**: `/v1/health` - 服务基本健康状态
- **就绪检查**: `/v1/ready` - 服务是否就绪接收请求
- **应用级别**: 每个应用的健康状态监控

### 7.2 日志管理

- **结构化日志**: JSON格式日志输出
- **日志级别**: 支持DEBUG、INFO、WARN、ERROR
- **日志轮转**: 支持日志文件自动轮转

### 7.3 指标监控

- **部署指标**: 部署成功率、部署时间
- **运行指标**: 应用运行状态、资源使用情况
- **网络指标**: 节点间通信延迟、成功率

## 8. 部署指南

### 8.1 环境要求

#### 8.1.1 Operator-PM 主控服务

- **操作系统**: Linux (推荐 Ubuntu 20.04+)
- **内存**: 最少 512MB，推荐 1GB+
- **CPU**: 最少 1 核，推荐 2 核+
- **网络**: 能够访问所有 Agent 节点

#### 8.1.2 Operator-PM-Agent 代理服务

- **操作系统**: Linux (推荐 Ubuntu 20.04+)
- **内存**: 最少 1GB，推荐 2GB+
- **CPU**: 最少 2 核，推荐 4 核+
- **Docker**: 如果使用 Docker 模式，需要安装 Docker Engine
- **网络**: 能够被主控服务访问

### 8.2 部署步骤

#### 8.2.1 部署 Operator-PM 主控服务

1. **准备配置文件**
   ```bash
   # 创建配置目录
   mkdir -p /etc/boreas-operator-pm
   
   # 复制配置文件
   cp cmd/operator-pm/configs/*.yaml /etc/boreas-operator-pm/
   ```

2. **修改配置文件**
   ```bash
   # 编辑主配置文件
   vim /etc/boreas-operator-pm/operator-pm.yaml
   
   # 编辑应用-节点映射
   vim /etc/boreas-operator-pm/app-to-nodes.yaml
   
   # 编辑节点-IP映射
   vim /etc/boreas-operator-pm/node-to-ip.yaml
   ```

3. **启动服务**
   ```bash
   # 使用 systemd 管理服务
   sudo systemctl start boreas-operator-pm
   sudo systemctl enable boreas-operator-pm
   ```

#### 8.2.2 部署 Operator-PM-Agent 代理服务

1. **准备配置文件**
   ```bash
   # 创建配置目录
   mkdir -p /etc/boreas-agent
   
   # 复制配置文件
   cp cmd/operator-pm-agent/configs/agent.yaml /etc/boreas-agent/
   ```

2. **修改配置文件**
   ```bash
   # 编辑 Agent 配置
   vim /etc/boreas-agent/agent.yaml
   
   # 设置正确的 Agent ID 和主机名
   # 配置 Docker 相关设置
   ```

3. **启动服务**
   ```bash
   # 使用 systemd 管理服务
   sudo systemctl start boreas-agent
   sudo systemctl enable boreas-agent
   ```

### 8.3 验证部署

#### 8.3.1 验证主控服务

```bash
# 检查服务状态
curl http://localhost:8080/v1/health

# 检查就绪状态
curl http://localhost:8080/v1/ready

# 测试应用状态查询
curl http://localhost:8080/v1/status/my-app
```

#### 8.3.2 验证代理服务

```bash
# 检查 Agent 健康状态
curl http://localhost:8081/v1/health

# 检查所有应用状态
curl http://localhost:8081/v1/status

# 测试应用部署
curl -X POST http://localhost:8081/v1/apply \
  -H "Content-Type: application/json" \
  -d '{"app":"test-app","version":"v1.0.0","pkg":{"type":"docker","image":"nginx:latest"}}'
```

### 8.4 监控和维护

#### 8.4.1 日志管理

```bash
# 查看主控服务日志
journalctl -u boreas-operator-pm -f

# 查看代理服务日志
journalctl -u boreas-agent -f

# 查看应用运行日志
tail -f /var/lib/boreas-agent/logs/*.log
```

#### 8.4.2 配置更新

```bash
# 更新配置文件后重启服务
sudo systemctl restart boreas-operator-pm
sudo systemctl restart boreas-agent

# 或者发送信号重新加载配置
sudo systemctl reload boreas-operator-pm
```

## 9. 未来扩展

### 9.1 功能扩展

- **滚动更新**: 支持零停机应用更新
- **回滚机制**: 支持快速回滚到上一版本
- **资源调度**: 基于资源使用情况的智能调度

### 9.2 集成扩展

- **CI/CD集成**: 与持续集成系统集成
- **监控集成**: 与Prometheus、Grafana等监控系统集成
- **日志集成**: 与ELK、Fluentd等日志系统集成

## 10. 总结

Operator-PM 组件采用简洁的配置驱动架构，通过主从模式实现了物理机应用的统一管理和部署。系统设计注重简单性、可靠性和可扩展性，能够满足中小规模物理机环境的应用管理需求。

主要优势：
- **配置简单**: 基于YAML的配置管理
- **架构清晰**: 主从分离，职责明确
- **无状态设计**: 不依赖数据库，易于部署
- **API友好**: 标准化的RESTful接口
- **监控完善**: 丰富的健康检查和状态监控
