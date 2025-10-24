# Boreas 部署系统架构文档

## 系统概述

Boreas 是一个多服务架构的自动化部署平台，支持 Kubernetes 和裸金属环境的应用部署。

## 架构设计

### 服务模块

#### 1. Master Service (核心服务)

作为整个系统的中枢模块，负责：
- 接收上游时间通知（如 Git Tag 创建）
- 创建和管理版本
- 规划和编排部署任务
- 下发部署任务给 Operator
- 监控和跟踪部署状态

**技术栈**: Go + Gin + PostgreSQL

#### 2. Operator Service (部署执行器)

作为 Master 的下游服务，负责：
- 执行上游下发的部署操作
- 维持部署状态
- 提供部署状态查询接口

**技术栈**: Go + Gin

**扩展性**: 
- 支持 Kubernetes 环境部署
- 支持裸金属环境部署
- 可根据环境类型实现不同的 Operator

#### 3. Frontend (Web 管理界面)

为 Master 服务提供的管理界面，支持：
- 实时状态查看
- 半自动任务编排/编辑
- 人工复核和确认

**技术栈**: React + TypeScript + Vite

## 数据流

```
Git Webhook → Master Service → Operator Service → K8S/裸金属
                ↓
            Frontend (监控和管理)
```

## 部署架构

- Master Service: 单实例或主从模式
- Operator Service: 多实例，按环境类型部署
- Frontend: Nginx 静态托管

## 目录结构

```
.
├── backend/
│   └── services/
│       ├── master/          # 核心服务
│       └── operator/        # 部署执行器
├── frontend/                # Web 管理界面
└── documentation/           # 项目文档
```
