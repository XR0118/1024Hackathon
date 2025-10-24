# CLAUDE.md

## 项目概述
这是七牛云 2024/10/24 Hackathon 的项目仓库。

## 开发环境
- 编程语言: Go
- 项目类型: Hackathon 项目

## 代码规范
- 遵循 Go 语言官方代码规范
- 使用 gofmt 格式化代码
- 遵循 Go 的命名约定

## 开发指南
1. 确保安装了 Go 开发环境
2. 克隆项目后，使用 `go mod tidy` 安装依赖
3. 提交前请确保代码通过 `go fmt` 格式化
4. 提交前请确保代码通过 `go vet` 检查

## 测试
- 使用 `go test ./...` 运行所有测试
- 测试覆盖率报告: `go test -coverprofile=coverage.out ./...`

## 构建
- 构建项目: `go build`
- 交叉编译示例: `GOOS=linux GOARCH=amd64 go build`

## Git 提交规范
- feat: 新功能
- fix: 修复问题
- docs: 文档更新
- style: 代码格式调整
- refactor: 重构代码
- test: 测试相关
- chore: 构建或辅助工具变动

## 注意事项
- 不要提交敏感信息（如 API 密钥、密码等）
- 遵循 .gitignore 中的规则
- 保持代码简洁、可读
