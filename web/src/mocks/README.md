# Mock 数据使用说明

## 概述

本项目提供了完整的 Mock 数据功能,允许前端在后端未准备好的情况下独立运行和测试所有页面功能。

## 快速开始

### 方式一: 使用 npm script (推荐)

```bash
cd web
npm run dev:mock
```

这将启动开发服务器并自动启用 Mock 模式。

### 方式二: 手动设置环境变量

```bash
cd web
VITE_USE_MOCK=true npm run dev
```

### 方式三: 修改 .env.development

编辑 `web/.env.development` 文件:

```env
VITE_USE_MOCK=true
```

然后正常启动开发服务器:

```bash
npm run dev
```

## Mock 数据说明

### 包含的 Mock 数据

Mock 数据文件位于 `src/mocks/data.ts`,包含以下内容:

1. **版本数据 (mockVersions)**
   - 3个版本记录
   - 包含 commit、branch、author 等完整信息

2. **应用数据 (mockApplications)**
   - 3个应用: user-service、order-service、payment-service
   - 包含版本历史和节点信息

3. **环境数据 (mockEnvironments)**
   - 3个环境: production、staging、development
   - 包含 k8s 和物理机两种类型

4. **部署数据 (mockDeployments)**
   - 4个部署记录
   - 包含各种状态: running、completed、failed、pending

5. **部署详情 (mockDeploymentDetails)**
   - 每个部署的详细步骤和日志
   - 完整的部署流程模拟

6. **仪表板数据**
   - 统计数据 (mockDashboardStats)
   - 趋势数据 (mockDeploymentTrends)

### 支持的 API 端点

所有 API 端点都已实现 Mock 支持:

#### 版本管理
- `GET /api/v1/versions` - 获取版本列表(支持搜索)
- `GET /api/v1/versions/:version` - 获取版本详情
- `POST /api/v1/versions` - 创建版本

#### 应用管理
- `GET /api/v1/applications` - 获取应用列表
- `GET /api/v1/applications/:id` - 获取应用详情
- `POST /api/v1/applications` - 创建应用
- `PUT /api/v1/applications/:id` - 更新应用

#### 环境管理
- `GET /api/v1/environments` - 获取环境列表
- `GET /api/v1/environments/:id` - 获取环境详情
- `POST /api/v1/environments` - 创建环境

#### 部署管理
- `GET /api/v1/deployments` - 获取部署列表(支持过滤)
- `GET /api/v1/deployments/:id` - 获取部署详情
- `POST /api/v1/deployments` - 创建部署
- `PUT /api/v1/deployments/:id` - 更新部署状态(确认/回滚/取消)

#### 仪表板
- `GET /api/v1/dashboard/stats` - 获取统计数据
- `GET /api/v1/dashboard/trends` - 获取趋势数据
- `GET /api/v1/dashboard/recent-deployments` - 获取最近部署

## 功能特性

### 1. 完整的 CRUD 操作
Mock API 支持创建、读取、更新操作,模拟真实 API 行为。

### 2. 搜索和过滤
支持版本搜索、部署状态过滤等功能。

### 3. 实时日志
Mock 请求会在浏览器控制台输出日志,方便调试:

```
[Mock API] GET /api/v1/versions
[Mock API] GET /api/v1/applications
```

### 4. 错误处理
Mock API 支持 404 等错误响应,模拟真实场景。

### 5. 数据一致性
Mock 数据之间相互关联,保证数据一致性:
- 版本和应用的关联
- 部署和应用、环境的关联
- 部署详情的完整步骤流程

## 自定义 Mock 数据

### 修改现有数据

编辑 `src/mocks/data.ts` 文件,直接修改导出的数据对象:

```typescript
export const mockVersions: Version[] = [
  {
    id: '1',
    version: 'v1.2.3',
    commit: 'abc123def456',
    branch: 'main',
    author: '张三',
    message: 'feat: 添加用户管理功能',
    createdAt: '2024-10-20T10:00:00Z',
    applications: ['app1', 'app2'],
  },
]
```

### 添加新的 API 端点

在 `src/mocks/handlers.ts` 中添加新的处理逻辑:

```typescript
if (url?.startsWith('/your-new-endpoint')) {
  if (method === 'GET') {
    return Promise.resolve({ data: yourMockData })
  }
}
```

## 页面测试场景

### 1. 仪表板页面
访问 http://localhost:3000 可以看到:
- 统计卡片显示数据
- 最近部署列表
- 所有交互功能正常

### 2. 版本列表页面
访问 http://localhost:3000/versions 可以:
- 查看版本列表
- 使用搜索功能
- 查看版本详情

### 3. 应用管理页面
访问 http://localhost:3000/applications 可以:
- 查看应用列表
- 查看应用详情
- 查看版本历史和节点信息

### 4. 部署管理页面
访问 http://localhost:3000/deployments 可以:
- 查看部署列表
- 按状态、环境、应用过滤
- 查看部署详情和步骤
- 查看实时日志

### 5. 创建部署页面
访问 http://localhost:3000/deployments/create 可以:
- 选择版本、应用、环境
- 提交创建请求
- 接收模拟响应

## 切换回真实 API

### 方式一: 使用默认 dev 命令

```bash
npm run dev
```

### 方式二: 设置环境变量

```bash
VITE_USE_MOCK=false npm run dev
```

### 方式三: 修改 .env 文件

编辑 `.env.development`:

```env
VITE_USE_MOCK=false
```

## 调试技巧

1. **检查 Mock 模式是否启用**
   
   在浏览器控制台查看是否有蓝色的启动日志:
   ```
   [Mock Mode] Mock API enabled
   ```

2. **查看 API 调用**
   
   所有 Mock API 请求都会在控制台输出:
   ```
   [Mock API] GET /api/v1/versions
   [Mock API] POST /api/v1/deployments
   ```

3. **检查环境变量**
   
   在代码中打印环境变量:
   ```typescript
   console.log('VITE_USE_MOCK:', import.meta.env.VITE_USE_MOCK)
   ```

## 注意事项

1. Mock 模式下的数据修改不会持久化,刷新页面后会重置
2. Mock 模式只在开发环境使用,生产构建会自动禁用
3. Mock 数据的类型定义与真实 API 保持一致
4. 建议定期与后端 API 对齐,更新 Mock 数据结构

## 常见问题

### Q: Mock 模式没有生效?
A: 检查以下几点:
   - 确认运行了 `npm run dev:mock` 或设置了正确的环境变量
   - 查看浏览器控制台是否有 `[Mock Mode] Mock API enabled` 日志
   - 清除浏览器缓存后重试

### Q: 如何添加更多测试数据?
A: 编辑 `src/mocks/data.ts`,在对应的数组中添加新对象即可。

### Q: 能否只 Mock 部分 API?
A: 可以,在 `src/mocks/handlers.ts` 中只处理需要 Mock 的端点,其他请求会正常发送到后端。

### Q: 如何模拟 API 延迟?
A: 在 handlers 中添加延迟:
   ```typescript
   await new Promise(resolve => setTimeout(resolve, 1000))
   return Promise.resolve({ data: mockData })
   ```

## 技术实现

Mock 功能基于 Axios 拦截器实现:
- 在请求发送前拦截
- 根据 URL 和 Method 返回对应的 Mock 数据
- 保留原有的响应处理逻辑
- 不需要额外的 Mock 服务器

## 下一步

1. 根据实际使用情况丰富 Mock 数据
2. 添加更多边界场景的测试数据
3. 与后端 API 文档保持同步
4. 考虑使用 Mock Service Worker (MSW) 进行更强大的 Mock 功能
