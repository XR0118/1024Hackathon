# 部署平台前端

基于React + TypeScript + Ant Design的GitOps部署平台前端应用。

## 技术栈

- **React 18** - UI框架
- **TypeScript** - 类型安全
- **Ant Design 5** - UI组件库
- **React Router 6** - 路由管理
- **Zustand** - 状态管理
- **Axios** - HTTP客户端
- **Vite** - 构建工具

## 目录结构

```
web/
├── src/
│   ├── components/     # 通用组件
│   │   └── Layout.tsx  # 主布局组件
│   ├── pages/          # 页面组件
│   │   ├── Dashboard.tsx
│   │   ├── Versions.tsx
│   │   ├── Applications.tsx
│   │   ├── Environments.tsx
│   │   ├── Deployments.tsx
│   │   ├── DeploymentDetail.tsx
│   │   └── CreateDeployment.tsx
│   ├── services/       # API服务
│   │   └── api.ts
│   ├── store/          # 状态管理
│   │   └── index.ts
│   ├── types/          # TypeScript类型定义
│   │   └── index.ts
│   ├── utils/          # 工具函数
│   │   └── index.ts
│   ├── App.tsx         # 根组件
│   ├── main.tsx        # 入口文件
│   └── index.css       # 全局样式
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

## 开发指南

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

应用将运行在 http://localhost:3000

API代理配置: `/api` -> `http://localhost:8080`

### 构建生产版本

```bash
npm run build
```

构建产物将输出到 `dist/` 目录。

### 预览生产构建

```bash
npm run preview
```

## 核心功能

### 1. 仪表板 (Dashboard)
- 系统状态概览(活跃版本、进行中部署、应用数、环境数)
- 最近部署列表
- 实时数据更新

### 2. 版本管理 (Versions)
- 版本列表展示
- 搜索和筛选(正常版本/回滚版本)
- 版本详情查看
- Git Tag关联

### 3. 应用管理 (Applications)
- 应用卡片展示
- 当前版本显示(按环境)
- 快速创建部署
- 应用详情查看

### 4. 环境管理 (Environments)
- 环境列表(K8S/物理机)
- 环境状态监控
- 应用数量统计

### 5. 部署管理 (Deployments)
- 部署列表展示
- 状态筛选和日期筛选
- 实时进度更新
- 部署详情查看

### 6. 部署详情 (Deployment Detail)
- 部署信息展示
- 流程步骤可视化
- 实时日志输出
- 人工确认/回滚操作
- 灰度发布进度

### 7. 新建部署 (Create Deployment)
- 分步向导流程
- 版本选择
- 应用选择(多选)
- 环境选择(多选)
- 部署配置(灰度、确认、回滚)
- 提交确认

## API接口

所有API请求通过 `/api/v1` 前缀访问后端服务。

主要接口:
- `GET /api/v1/versions` - 获取版本列表
- `GET /api/v1/applications` - 获取应用列表
- `GET /api/v1/environments` - 获取环境列表
- `GET /api/v1/deployments` - 获取部署列表
- `POST /api/v1/deployments` - 创建部署
- `GET /api/v1/deployments/:id` - 获取部署详情
- `PUT /api/v1/deployments/:id` - 更新部署(确认/回滚)

详细API文档参见 `src/services/api.ts`

## 状态管理

使用Zustand进行轻量级状态管理:

- `useAppStore` - 应用数据状态(版本、应用、环境、部署)
- `useUIStore` - UI状态(侧边栏折叠等)

## 设计规范

### 颜色规范
- 成功: `#52c41a` (绿色)
- 进行中: `#1890ff` (蓝色)
- 失败/错误: `#ff4d4f` (红色)
- 警告/待确认: `#faad14` (橙色)
- 中性/未开始: `#d9d9d9` (灰色)

### 图标使用
- 成功: ✓ (CheckOutlined)
- 失败: ✗ (CloseOutlined)
- 进行中: ⟳ (SyncOutlined)
- 待确认: ⏸ (PauseCircleOutlined)
- 回滚: ↶ (RollbackOutlined)

## 浏览器支持

- Chrome (最新版)
- Firefox (最新版)
- Safari (最新版)
- Edge (最新版)

## 注意事项

1. **实时更新**: 部署列表和详情页面会自动轮询更新数据
2. **权限控制**: 需配合后端API实现权限验证
3. **错误处理**: API调用失败会在控制台输出错误,可根据需要添加用户提示
4. **响应式设计**: 支持桌面、平板和移动端访问

## 后续开发

- [ ] 添加用户认证和权限管理
- [ ] 实现应用详情页面
- [ ] 添加部署统计图表
- [ ] 支持批量部署操作
- [ ] 添加通知中心
- [ ] 优化移动端体验
