# CLAUDE.md

## 项目概述
这是七牛云 2024/10/24 Hackathon 的项目仓库。

## 开发环境
- 后端: Go
- 前端: TypeScript + React
- 项目类型: Hackathon 项目

## 代码规范

### 后端 (Go)
- 遵循 Go 语言官方代码规范
- 使用 gofmt 格式化代码
- 遵循 Go 的命名约定

### 前端 (TypeScript + React)
- 使用 TypeScript 严格模式
- 遵循 React 官方最佳实践
- 使用 ESLint 进行代码检查
- 使用 Prettier 进行代码格式化
- 组件命名使用 PascalCase
- 文件命名使用 kebab-case 或 PascalCase
- Hook 命名以 use 开头
- 类型定义文件使用 .d.ts 或在文件内定义
- 优先使用函数组件和 Hooks

## 开发指南

### 后端 (Go)
1. 确保安装了 Go 开发环境
2. 克隆项目后，使用 `go mod tidy` 安装依赖
3. 提交前请确保代码通过 `go fmt` 格式化
4. 提交前请确保代码通过 `go vet` 检查

### 前端 (TypeScript + React)
1. 确保安装了 Node.js 18+ 和 npm/yarn/pnpm
2. 克隆项目后，在前端目录运行 `npm install` 或 `yarn install` 安装依赖
3. 启动开发服务器: `npm run dev` 或 `yarn dev`
4. 提交前运行 `npm run lint` 进行代码检查
5. 提交前运行 `npm run format` 进行代码格式化
6. 提交前运行 `npm run type-check` 进行类型检查

## 测试

### 后端 (Go)
- 使用 `go test ./...` 运行所有测试
- 测试覆盖率报告: `go test -coverprofile=coverage.out ./...`

### 前端 (TypeScript + React)
- 使用 Jest + React Testing Library 进行单元测试
- 运行测试: `npm test` 或 `yarn test`
- 运行测试覆盖率: `npm run test:coverage` 或 `yarn test:coverage`
- E2E 测试推荐使用 Playwright 或 Cypress

## 构建

### 后端 (Go)
- 构建项目: `go build`
- 交叉编译示例: `GOOS=linux GOARCH=amd64 go build`

### 前端 (TypeScript + React)
- 开发环境: `npm run dev` 或 `yarn dev`
- 生产构建: `npm run build` 或 `yarn build`
- 预览构建结果: `npm run preview` 或 `yarn preview`
- 推荐使用 Vite 作为构建工具

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
