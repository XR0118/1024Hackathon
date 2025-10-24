# Operator Service

部署执行器，为 master 下游。由于 K8S 与裸金属环境的不同，有不同的执行器实现。

## 功能

- 执行上游下发的部署操作
- 维持部署状态
- 提供部署状态，供上游查询

## 构建

```bash
cd backend/services/operator
go mod tidy
go build -o bin/operator ./cmd
```

## 运行

```bash
./bin/operator
```

## API 端口

默认端口: 8081
