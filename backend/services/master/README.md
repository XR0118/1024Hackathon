# Master Service

核心服务，作为整个系统的中枢模块。

## 功能

- 上游承接时间通知，包含 Git Tag 等，即创建版本
- 内部实现规划、编排部署任务
- 运行编排好的部署任务，让下游具体操作
- 同时监控跟踪部署状态

## 构建

```bash
cd backend/services/master
go mod tidy
go build -o bin/master ./cmd
```

## 运行

```bash
./bin/master
```

## API 端口

默认端口: 8080
