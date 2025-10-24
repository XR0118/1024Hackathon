# Boreas Backend Services

后端服务模块，采用多服务架构。

## 服务列表

### Master Service
核心服务，负责版本管理、任务编排和部署监控。

详见: [services/master/README.md](services/master/README.md)

### Operator Service
部署执行器，负责执行具体的部署操作。

详见: [services/operator/README.md](services/operator/README.md)

## 技术栈

- Go 1.21+
- Gin Web Framework
- PostgreSQL
- GORM

## 开发要求

参见根目录的 [CLAUDE.md](../CLAUDE.md)
