import type {
  Version,
  Application,
  Environment,
  Deployment,
  DeploymentDetail,
  DashboardStats,
  DeploymentTrend,
} from '@/types'

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
  {
    id: '2',
    version: 'v1.2.2',
    commit: 'def456ghi789',
    branch: 'main',
    author: '李四',
    message: 'fix: 修复登录bug',
    createdAt: '2024-10-19T15:30:00Z',
    applications: ['app1'],
  },
  {
    id: '3',
    version: 'v1.2.1',
    commit: 'ghi789jkl012',
    branch: 'main',
    author: '王五',
    message: 'refactor: 重构数据库连接池',
    createdAt: '2024-10-18T09:00:00Z',
    applications: ['app3'],
  },
]

export const mockApplications: Application[] = [
  {
    id: 'app1',
    name: 'user-service',
    description: '用户服务',
    repository: 'https://github.com/example/user-service',
    currentVersion: 'v1.2.3',
    latestVersion: 'v1.2.3',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-10-20T10:00:00Z',
    versions: [
      { version: 'v1.2.3', deployedAt: '2024-10-20T10:00:00Z' },
      { version: 'v1.2.2', deployedAt: '2024-10-19T15:30:00Z' },
    ],
    nodes: [
      { name: 'node1', ip: '192.168.1.10', status: 'running' },
      { name: 'node2', ip: '192.168.1.11', status: 'running' },
    ],
  },
  {
    id: 'app2',
    name: 'order-service',
    description: '订单服务',
    repository: 'https://github.com/example/order-service',
    currentVersion: 'v1.2.3',
    latestVersion: 'v1.2.3',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-10-20T10:00:00Z',
    versions: [
      { version: 'v1.2.3', deployedAt: '2024-10-20T10:00:00Z' },
    ],
    nodes: [
      { name: 'node3', ip: '192.168.1.20', status: 'running' },
    ],
  },
  {
    id: 'app3',
    name: 'payment-service',
    description: '支付服务',
    repository: 'https://github.com/example/payment-service',
    currentVersion: 'v1.2.1',
    latestVersion: 'v1.2.1',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-10-18T09:00:00Z',
    versions: [
      { version: 'v1.2.1', deployedAt: '2024-10-18T09:00:00Z' },
    ],
    nodes: [
      { name: 'node4', ip: '192.168.1.30', status: 'running' },
      { name: 'node5', ip: '192.168.1.31', status: 'stopped' },
    ],
  },
]

export const mockEnvironments: Environment[] = [
  {
    id: 'env1',
    name: 'production',
    type: 'k8s',
    description: '生产环境',
    config: {
      namespace: 'prod',
      cluster: 'prod-cluster',
    },
    createdAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 'env2',
    name: 'staging',
    type: 'k8s',
    description: '预发布环境',
    config: {
      namespace: 'staging',
      cluster: 'staging-cluster',
    },
    createdAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 'env3',
    name: 'development',
    type: 'physical',
    description: '开发环境',
    config: {
      servers: ['dev-server-1', 'dev-server-2'],
    },
    createdAt: '2024-01-01T00:00:00Z',
  },
]

export const mockDeployments: Deployment[] = [
  {
    id: 'deploy1',
    version: 'v1.2.3',
    applications: ['user-service', 'order-service'],
    environments: ['production'],
    status: 'running',
    progress: 75,
    createdAt: '2024-10-20T10:00:00Z',
    updatedAt: '2024-10-20T10:30:00Z',
    createdBy: '张三',
  },
  {
    id: 'deploy2',
    version: 'v1.2.2',
    applications: ['user-service'],
    environments: ['staging'],
    status: 'completed',
    progress: 100,
    createdAt: '2024-10-19T15:30:00Z',
    updatedAt: '2024-10-19T16:00:00Z',
    createdBy: '李四',
  },
  {
    id: 'deploy3',
    version: 'v1.2.1',
    applications: ['payment-service'],
    environments: ['development'],
    status: 'failed',
    progress: 50,
    error: '节点连接超时',
    createdAt: '2024-10-18T09:00:00Z',
    updatedAt: '2024-10-18T09:30:00Z',
    createdBy: '王五',
  },
  {
    id: 'deploy4',
    version: 'v1.2.0',
    applications: ['user-service'],
    environments: ['production'],
    status: 'pending',
    progress: 0,
    createdAt: '2024-10-17T14:00:00Z',
    updatedAt: '2024-10-17T14:00:00Z',
    createdBy: '赵六',
  },
]

export const mockDeploymentDetails: Record<string, DeploymentDetail> = {
  deploy1: {
    id: 'deploy1',
    version: 'v1.2.3',
    applications: ['user-service', 'order-service'],
    environments: ['production'],
    status: 'running',
    progress: 75,
    createdAt: '2024-10-20T10:00:00Z',
    updatedAt: '2024-10-20T10:30:00Z',
    createdBy: '张三',
    steps: [
      {
        name: '准备部署',
        status: 'completed',
        startedAt: '2024-10-20T10:00:00Z',
        completedAt: '2024-10-20T10:05:00Z',
        logs: ['检查版本信息...', '验证配置文件...', '准备完成'],
      },
      {
        name: '拉取镜像',
        status: 'completed',
        startedAt: '2024-10-20T10:05:00Z',
        completedAt: '2024-10-20T10:15:00Z',
        logs: ['拉取镜像 user-service:v1.2.3', '拉取镜像 order-service:v1.2.3', '镜像拉取完成'],
      },
      {
        name: '更新服务',
        status: 'running',
        startedAt: '2024-10-20T10:15:00Z',
        logs: ['更新 user-service 配置', '重启 user-service pod', '等待服务就绪...'],
      },
      {
        name: '健康检查',
        status: 'pending',
      },
    ],
  },
  deploy2: {
    id: 'deploy2',
    version: 'v1.2.2',
    applications: ['user-service'],
    environments: ['staging'],
    status: 'completed',
    progress: 100,
    createdAt: '2024-10-19T15:30:00Z',
    updatedAt: '2024-10-19T16:00:00Z',
    createdBy: '李四',
    steps: [
      {
        name: '准备部署',
        status: 'completed',
        startedAt: '2024-10-19T15:30:00Z',
        completedAt: '2024-10-19T15:35:00Z',
        logs: ['检查版本信息...', '验证配置文件...', '准备完成'],
      },
      {
        name: '拉取镜像',
        status: 'completed',
        startedAt: '2024-10-19T15:35:00Z',
        completedAt: '2024-10-19T15:40:00Z',
        logs: ['拉取镜像 user-service:v1.2.2', '镜像拉取完成'],
      },
      {
        name: '更新服务',
        status: 'completed',
        startedAt: '2024-10-19T15:40:00Z',
        completedAt: '2024-10-19T15:55:00Z',
        logs: ['更新 user-service 配置', '重启 user-service pod', '服务更新完成'],
      },
      {
        name: '健康检查',
        status: 'completed',
        startedAt: '2024-10-19T15:55:00Z',
        completedAt: '2024-10-19T16:00:00Z',
        logs: ['检查服务健康状态', '所有健康检查通过', '部署成功'],
      },
    ],
  },
  deploy3: {
    id: 'deploy3',
    version: 'v1.2.1',
    applications: ['payment-service'],
    environments: ['development'],
    status: 'failed',
    progress: 50,
    error: '节点连接超时',
    createdAt: '2024-10-18T09:00:00Z',
    updatedAt: '2024-10-18T09:30:00Z',
    createdBy: '王五',
    steps: [
      {
        name: '准备部署',
        status: 'completed',
        startedAt: '2024-10-18T09:00:00Z',
        completedAt: '2024-10-18T09:05:00Z',
        logs: ['检查版本信息...', '验证配置文件...', '准备完成'],
      },
      {
        name: '拉取镜像',
        status: 'failed',
        startedAt: '2024-10-18T09:05:00Z',
        completedAt: '2024-10-18T09:30:00Z',
        logs: [
          '拉取镜像 payment-service:v1.2.1',
          '连接到 node4 (192.168.1.30)...',
          '错误: 连接超时',
          '重试 1/3...',
          '错误: 连接超时',
          '重试 2/3...',
          '错误: 连接超时',
          '部署失败: 节点连接超时',
        ],
      },
    ],
  },
  deploy4: {
    id: 'deploy4',
    version: 'v1.2.0',
    applications: ['user-service'],
    environments: ['production'],
    status: 'pending',
    progress: 0,
    createdAt: '2024-10-17T14:00:00Z',
    updatedAt: '2024-10-17T14:00:00Z',
    createdBy: '赵六',
    steps: [
      {
        name: '准备部署',
        status: 'pending',
      },
      {
        name: '拉取镜像',
        status: 'pending',
      },
      {
        name: '更新服务',
        status: 'pending',
      },
      {
        name: '健康检查',
        status: 'pending',
      },
    ],
  },
}

export const mockDashboardStats: DashboardStats = {
  activeVersions: 3,
  runningDeployments: 1,
  totalApplications: 3,
  totalEnvironments: 3,
}

export const mockDeploymentTrends: DeploymentTrend[] = [
  { date: '2024-10-14', success: 5, failed: 1 },
  { date: '2024-10-15', success: 3, failed: 0 },
  { date: '2024-10-16', success: 7, failed: 2 },
  { date: '2024-10-17', success: 4, failed: 1 },
  { date: '2024-10-18', success: 6, failed: 1 },
  { date: '2024-10-19', success: 5, failed: 0 },
  { date: '2024-10-20', success: 8, failed: 0 },
]
