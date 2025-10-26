import type {
  Version,
  Application,
  ApplicationVersionsResponse,
  Environment,
  Deployment,
  DeploymentDetail,
  DashboardStats,
  DeploymentTrend,
} from '@/types'

export const mockVersions: Version[] = [
  {
    id: 'version-001',
    version: 'v1.2.5',
    git_tag: 'v1.2.5',
    git_commit: 'a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0',
    repository: 'https://github.com/example/microservices',
    status: 'normal',
    created_by: 'admin@example.com',
    created_at: '2024-10-21T14:00:00Z',
    description: '新增用户认证功能',
    app_builds: [
      { app_name: 'user-service', docker_image: 'registry.example.com/user-service:v1.2.5' },
      { app_name: 'order-service', docker_image: 'registry.example.com/order-service:v1.2.5' },
      { app_name: 'payment-service', docker_image: 'registry.example.com/payment-service:v1.2.5' },
      { app_name: 'notification-service', docker_image: 'registry.example.com/notification-service:v1.2.5' },
      { app_name: 'analytics-service', docker_image: 'registry.example.com/analytics-service:v1.2.5' },
    ],
  },
  {
    id: 'version-002',
    version: 'v1.2.4',
    git_tag: 'v1.2.4',
    git_commit: 'b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1',
    repository: 'https://github.com/example/microservices',
    status: 'revert',
    created_by: 'admin@example.com',
    created_at: '2024-10-20T10:00:00Z',
    description: '回滚版本：修复支付bug',
    app_builds: [
      { app_name: 'user-service', docker_image: 'registry.example.com/user-service:v1.2.4' },
      { app_name: 'order-service', docker_image: 'registry.example.com/order-service:v1.2.4' },
      { app_name: 'payment-service', docker_image: 'registry.example.com/payment-service:v1.2.4' },
      { app_name: 'notification-service', docker_image: 'registry.example.com/notification-service:v1.2.4' },
      { app_name: 'analytics-service', docker_image: 'registry.example.com/analytics-service:v1.2.4' },
    ],
  },
  {
    id: 'version-003',
    version: 'v1.2.3',
    git_tag: 'v1.2.3',
    git_commit: 'c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2',
    repository: 'https://github.com/example/microservices',
    status: 'normal',
    created_by: 'developer@example.com',
    created_at: '2024-10-19T15:30:00Z',
    description: '优化订单处理性能',
    app_builds: [
      { app_name: 'user-service', docker_image: 'registry.example.com/user-service:v1.2.3' },
      { app_name: 'order-service', docker_image: 'registry.example.com/order-service:v1.2.3' },
    ],
  },
  {
    id: 'version-004',
    version: 'v1.2.2',
    git_tag: 'v1.2.2',
    git_commit: 'd4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3',
    repository: 'https://github.com/example/microservices',
    status: 'normal',
    created_by: 'developer@example.com',
    created_at: '2024-10-18T09:00:00Z',
    description: '修复用户服务bug',
    app_builds: [
      { app_name: 'user-service', docker_image: 'registry.example.com/user-service:v1.2.2' },
    ],
  },
  {
    id: 'version-005',
    version: 'v1.2.1',
    git_tag: 'v1.2.1',
    git_commit: 'e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4',
    repository: 'https://github.com/example/microservices',
    status: 'normal',
    created_by: 'developer@example.com',
    created_at: '2024-10-17T14:00:00Z',
    description: '增加支付功能',
    app_builds: [
      { app_name: 'payment-service', docker_image: 'registry.example.com/payment-service:v1.2.1' },
    ],
  },
]

export const mockApplications: Application[] = [
  {
    id: 'app-001',
    name: 'user-service',
    description: '用户管理微服务',
    repository: 'https://github.com/example/user-service',
    type: 'microservice',
    config: {
      dockerfile: 'Dockerfile',
      context: '.',
    },
    created_at: '2024-01-15T08:00:00Z',
    updated_at: '2024-10-20T10:00:00Z',
  },
  {
    id: 'app-002',
    name: 'order-service',
    description: '订单管理微服务',
    repository: 'https://github.com/example/order-service',
    type: 'microservice',
    config: {
      dockerfile: 'Dockerfile',
      context: '.',
    },
    created_at: '2024-02-10T08:00:00Z',
    updated_at: '2024-10-20T10:00:00Z',
  },
  {
    id: 'app-003',
    name: 'payment-service',
    description: '支付处理服务',
    repository: 'https://github.com/example/payment-service',
    type: 'microservice',
    config: {
      dockerfile: 'Dockerfile',
      context: '.',
    },
    created_at: '2024-03-05T08:00:00Z',
    updated_at: '2024-10-18T09:00:00Z',
  },
  {
    id: 'app-004',
    name: 'notification-service',
    description: '消息通知服务',
    repository: 'https://github.com/example/notification-service',
    type: 'microservice',
    created_at: '2024-04-01T08:00:00Z',
    updated_at: '2024-10-21T14:00:00Z',
  },
  {
    id: 'app-005',
    name: 'analytics-service',
    description: '数据分析服务',
    repository: 'https://github.com/example/analytics-service',
    type: 'microservice',
    created_at: '2024-05-10T08:00:00Z',
    updated_at: '2024-10-21T14:00:00Z',
  },
]

// 应用版本信息 Mock 数据（从 Operator 查询返回）
export const mockApplicationVersions: Record<string, ApplicationVersionsResponse> = {
  'user-service': {
    application_id: 'app-001',
    name: 'user-service',
    versions: [
      {
        version: 'v1.2.5',
        status: 'normal',
        health: 95,
        coverage: 80,
        last_updated_at: '2024-10-21T14:30:00Z',
        nodes: [
          { name: 'prod-node-1', health: 98, last_updated_at: '2024-10-21T14:30:00Z' },
          { name: 'prod-node-2', health: 92, last_updated_at: '2024-10-21T14:28:00Z' },
        ],
      },
      {
        version: 'v1.2.3',
        status: 'normal',
        health: 92,
        coverage: 20,
        last_updated_at: '2024-10-20T10:00:00Z',
        nodes: [
          { name: 'staging-node-1', health: 95, last_updated_at: '2024-10-20T10:00:00Z' },
          { name: 'staging-node-2', health: 89, last_updated_at: '2024-10-20T09:55:00Z' },
        ],
      },
    ],
  },
  'order-service': {
    application_id: 'app-002',
    name: 'order-service',
    versions: [
      {
        version: 'v1.2.5',
        status: 'normal',
        health: 88,
        coverage: 100,
        last_updated_at: '2024-10-21T14:30:00Z',
        nodes: [
          { name: 'prod-node-3', health: 90, last_updated_at: '2024-10-21T14:30:00Z' },
        ],
      },
      {
        version: 'v1.2.3',
        status: 'normal',
        health: 88,
        coverage: 100,
        last_updated_at: '2024-10-20T10:00:00Z',
        nodes: [
          { name: 'staging-node-3', health: 90, last_updated_at: '2024-10-20T10:00:00Z' },
        ],
      },
    ],
  },
  'payment-service': {
    application_id: 'app-003',
    name: 'payment-service',
    versions: [
      {
        version: 'v1.2.1',
        status: 'normal',
        health: 85,
        coverage: 72,
        last_updated_at: '2024-10-18T09:00:00Z',
        nodes: [
          { name: 'dev-node-4', health: 88, last_updated_at: '2024-10-18T09:00:00Z' },
          { name: 'dev-node-5', health: 82, last_updated_at: '2024-10-18T08:55:00Z' },
        ],
      },
    ],
  },
  'notification-service': {
    application_id: 'app-004',
    name: 'notification-service',
    versions: [
      {
        version: 'v1.2.5',
        status: 'normal',
        health: 90,
        coverage: 100,
        last_updated_at: '2024-10-21T14:30:00Z',
        nodes: [
          { name: 'prod-node-4', health: 90, last_updated_at: '2024-10-21T14:30:00Z' },
        ],
      },
    ],
  },
  'analytics-service': {
    application_id: 'app-005',
    name: 'analytics-service',
    versions: [
      {
        version: 'v1.2.5',
        status: 'normal',
        health: 78,
        coverage: 60,
        last_updated_at: '2024-10-21T14:30:00Z',
        nodes: [
          { name: 'prod-node-5', health: 80, last_updated_at: '2024-10-21T14:30:00Z' },
          { name: 'prod-node-6', health: 76, last_updated_at: '2024-10-21T14:25:00Z' },
        ],
      },
      {
        version: 'v1.2.4',
        status: 'revert',
        health: 65,
        coverage: 40,
        last_updated_at: '2024-10-21T10:00:00Z',
        nodes: [
          { name: 'staging-node-4', health: 65, last_updated_at: '2024-10-21T10:00:00Z' },
        ],
      },
    ],
  },
}

export const mockEnvironments: Environment[] = [
  {
    id: 'env1',
    name: 'production',
    type: 'kubernetes',
    is_active: true,
    config: {
      cluster: 'prod-cluster',
      namespace: 'production',
    },
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-10-20T10:00:00Z',
  },
  {
    id: 'env2',
    name: 'staging',
    type: 'kubernetes',
    is_active: true,
    config: {
      cluster: 'staging-cluster',
      namespace: 'staging',
    },
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-10-20T10:00:00Z',
  },
  {
    id: 'env3',
    name: 'development',
    type: 'physical',
    is_active: true,
    config: {
      host: '192.168.1.100',
      port: '8080',
    },
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-10-20T10:00:00Z',
  },
]

export const mockDeployments: Deployment[] = [
  {
    id: 'deploy-001',
    versionId: 'v1.2.5',
    version: 'v1.2.5',
    applicationIds: ['user-service', 'order-service', 'payment-service', 'notification-service', 'analytics-service', 'report-service', 'message-service'],
    applications: ['user-service', 'order-service', 'payment-service', 'notification-service', 'analytics-service', 'report-service', 'message-service'],
    environmentIds: ['env1'],
    environments: ['production'],

    status: 'running',
    progress: 65,
    createdAt: '2024-10-21T14:00:00Z',
    updatedAt: '2024-10-21T14:30:00Z',
    requireConfirm: false,
    grayscaleEnabled: true,
    grayscaleRatio: 30,
  },
  {
    id: 'deploy-002',
    versionId: 'v1.2.4',
    version: 'v1.2.4',
    applicationIds: ['user-service', 'order-service', 'payment-service', 'notification-service', 'analytics-service'],
    applications: ['user-service', 'order-service', 'payment-service', 'notification-service', 'analytics-service'],
    environmentIds: ['env2'],
    environments: ['staging'],

    status: 'success',
    progress: 100,
    createdAt: '2024-10-21T10:00:00Z',
    updatedAt: '2024-10-21T10:45:00Z',
    duration: 2700,
    requireConfirm: false,
    grayscaleEnabled: false,
  },
  {
    id: 'deploy-003',
    versionId: 'v1.2.3',
    version: 'v1.2.3',
    applicationIds: ['user-service', 'order-service', 'payment-service'],
    applications: ['user-service', 'order-service', 'payment-service'],
    environmentIds: ['env1'],
    environments: ['production'],

    status: 'success',
    progress: 100,
    createdAt: '2024-10-20T10:00:00Z',
    updatedAt: '2024-10-20T10:30:00Z',
    duration: 1800,
    requireConfirm: false,
    grayscaleEnabled: true,
    grayscaleRatio: 50,
  },
  {
    id: 'deploy-004',
    versionId: 'v1.2.2',
    version: 'v1.2.2',
    applicationIds: ['user-service', 'order-service'],
    applications: ['user-service', 'order-service'],
    environmentIds: ['env2'],
    environments: ['staging'],

    status: 'success',
    progress: 100,
    createdAt: '2024-10-19T15:30:00Z',
    updatedAt: '2024-10-19T16:00:00Z',
    duration: 1800,
    requireConfirm: false,
    grayscaleEnabled: false,
  },
  {
    id: 'deploy-005',
    versionId: 'v1.2.1',
    version: 'v1.2.1',
    applicationIds: ['payment-service'],
    applications: ['payment-service'],
    environmentIds: ['env3'],
    environments: ['development'],
    status: 'paused',
    progress: 30,
    createdAt: '2024-10-18T09:00:00Z',
    updatedAt: '2024-10-18T09:15:00Z',
    requireConfirm: false,
    grayscaleEnabled: false,
  },
  {
    id: 'deploy-006',
    versionId: 'v1.2.0',
    version: 'v1.2.0',
    applicationIds: ['user-service'],
    applications: ['user-service'],
    environmentIds: ['env1'],
    environments: ['production'],

    status: 'pending',
    progress: 0,
    createdAt: '2024-10-17T14:00:00Z',
    updatedAt: '2024-10-17T14:00:00Z',
    requireConfirm: true,
    grayscaleEnabled: false,
  },
]

// Mock 数据使用宽松的类型检查
export const mockDeploymentDetails = {
  'deploy-001': {
    id: 'deploy-001',
    versionId: 'v1.2.5',
    version: 'v1.2.5',
    applicationIds: ['user-service', 'order-service', 'payment-service', 'notification-service', 'analytics-service', 'report-service', 'message-service'],
    applications: ['user-service', 'order-service', 'payment-service', 'notification-service', 'analytics-service', 'report-service', 'message-service'],
    environmentIds: ['env1'],
    environments: ['production'],

    status: 'running',
    progress: 65,
    createdAt: '2024-10-21T14:00:00Z',
    updatedAt: '2024-10-21T14:30:00Z',
    requireConfirm: false,
    grayscaleEnabled: true,
    grayscaleRatio: 30,
    tasks: [
      // 顶点1：准备部署（无依赖）
      {
        id: 'task-1',
        name: '准备部署',
        type: 'prepare',
        step: 'completed',

        status: 'success',
        dependencies: [],
        duration: 300,
        logs: ['检查版本信息...', '验证配置文件...', '准备完成'],
      },
      // 顶点2：环境检查（并行开始，无依赖）
      {
        id: 'task-2',
        name: '环境检查',
        type: 'health_check',
        step: 'completed',

        status: 'success',
        dependencies: [],
        duration: 120,
        logs: ['检查生产环境状态...', '环境正常'],
      },
      // 依赖 task-1
      {
        id: 'task-3',
        name: '构建镜像',
        type: 'build',
        step: 'completed',

        status: 'success',
        dependencies: ['task-1'],
        duration: 450,
        appId: 'user-service',
        logs: ['构建 user-service 镜像', '构建 order-service 镜像', '镜像构建完成'],
      },
      // 依赖 task-3
      {
        id: 'task-4',
        name: '等待灰度窗口',
        type: 'sleep',
        step: 'completed',

        status: 'success',
        dependencies: ['task-3'],
        duration: 60,
        params: {
          sleepDuration: 60,
        },
      },
      // 依赖 task-4 和 task-2（多个上游）
      {
        id: 'task-5',
        name: '人工确认',
        type: 'approval',
        status: 'waiting_approval',
        dependencies: ['task-4', 'task-2'],
        params: {
          approvalNote: '请确认是否继续灰度部署到生产环境',
        },
      },
      // 依赖 task-5
      {
        id: 'task-6',
        name: '灰度部署',
        type: 'deploy',
        step: 'pending',

        status: 'pending',
        dependencies: ['task-5'],
        appId: 'user-service',
        logs: [],
      },
      // 依赖 task-6
      {
        id: 'task-7',
        name: '健康检查',
        type: 'health_check',
        step: 'pending',

        status: 'pending',
        dependencies: ['task-6'],
        appId: 'user-service',
      },
    ],
    logs: [
      { timestamp: '2024-10-20T10:00:00Z', level: 'info', message: '部署开始' },
      { timestamp: '2024-10-20T10:05:00Z', level: 'info', message: '准备完成' },
      { timestamp: '2024-10-20T10:15:00Z', level: 'info', message: '镜像拉取完成' },
      { timestamp: '2024-10-20T10:15:00Z', level: 'info', message: '开始更新服务' },
      { timestamp: '2024-10-20T10:20:00Z', level: 'info', message: '更新中...' },
    ],
  },
  'deploy-003': {
    id: 'deploy-003',
    versionId: 'v1.2.3',
    version: 'v1.2.3',
    applicationIds: ['user-service', 'order-service', 'payment-service'],
    applications: ['user-service', 'order-service', 'payment-service'],
    environmentIds: ['env1'],
    environments: ['production'],

    status: 'success',
    progress: 100,
    createdAt: '2024-10-20T10:00:00Z',
    updatedAt: '2024-10-20T10:30:00Z',
    duration: 1800,
    requireConfirm: false,
    grayscaleEnabled: true,
    grayscaleRatio: 50,
    tasks: [
      {
        id: 'task-1',
        name: '准备部署',
        type: 'prepare',
        step: 'completed',

        status: 'success',
        dependencies: [],
        duration: 300,
      },
      {
        id: 'task-2',
        name: '构建镜像',
        type: 'build',
        step: 'completed',

        status: 'success',
        dependencies: ['task-1'],
        duration: 600,
        appId: 'order-service',
      },
      {
        id: 'task-3',
        name: '部署服务',
        type: 'deploy',
        step: 'completed',

        status: 'success',
        dependencies: ['task-2'],
        duration: 600,
        appId: 'order-service',
      },
      {
        id: 'task-4',
        name: '健康检查',
        type: 'health_check',
        step: 'completed',

        status: 'success',
        dependencies: ['task-3'],
        duration: 300,
        appId: 'order-service',
      },
    ],
    logs: [
      { timestamp: '2024-10-20T10:00:00Z', level: 'info', message: '部署开始' },
      { timestamp: '2024-10-20T10:05:00Z', level: 'info', message: '准备完成' },
      { timestamp: '2024-10-20T10:15:00Z', level: 'info', message: '部署成功' },
    ],
  },
  'deploy-005': {
    id: 'deploy-005',
    versionId: 'v1.2.1',
    version: 'v1.2.1',
    applicationIds: ['payment-service'],
    applications: ['payment-service'],
    environmentIds: ['env3'],
    environments: ['development'],
    status: 'paused',
    progress: 30,
    createdAt: '2024-10-18T09:00:00Z',
    updatedAt: '2024-10-18T09:15:00Z',
    requireConfirm: false,
    grayscaleEnabled: false,
    tasks: [
      {
        id: 'task-1',
        name: '准备部署',
        type: 'prepare',
        step: 'completed',

        status: 'success',
        dependencies: [],
        duration: 180,
      },
      {
        id: 'task-2',
        name: '构建镜像',
        type: 'build',
        step: 'running',

        status: 'running',
        dependencies: ['task-1'],
        appId: 'payment-service',
        logs: ['构建 payment-service 镜像中...', '已完成 30%'],
      },
      {
        id: 'task-3',
        name: '部署服务',
        type: 'deploy',
        step: 'pending',

        status: 'pending',
        dependencies: ['task-2'],
        appId: 'payment-service',
      },
      {
        id: 'task-4',
        name: '健康检查',
        type: 'health_check',
        step: 'pending',

        status: 'pending',
        dependencies: ['task-3'],
        appId: 'payment-service',
      },
    ],
    logs: [
      { timestamp: '2024-10-18T09:00:00Z', level: 'info', message: '部署任务已创建' },
      { timestamp: '2024-10-18T09:03:00Z', level: 'info', message: '准备部署完成' },
      { timestamp: '2024-10-18T09:03:00Z', level: 'info', message: '开始构建镜像' },
      { timestamp: '2024-10-18T09:15:00Z', level: 'warn', message: '部署已暂停' },
    ],
  },
  'deploy-006': {
    id: 'deploy-006',
    versionId: 'v1.2.0',
    version: 'v1.2.0',
    applicationIds: ['user-service'],
    applications: ['user-service'],
    environmentIds: ['env1'],
    environments: ['production'],

    status: 'pending',
    progress: 0,
    createdAt: '2024-10-17T14:00:00Z',
    updatedAt: '2024-10-17T14:00:00Z',
    requireConfirm: true,
    grayscaleEnabled: false,
    tasks: [
      {
        id: 'task-1',
        name: '准备部署',
        type: 'prepare',
        step: 'pending',

        status: 'pending',
        dependencies: [],
      },
      {
        id: 'task-2',
        name: '构建镜像',
        type: 'build',
        step: 'pending',

        status: 'pending',
        dependencies: ['task-1'],
        appId: 'user-service',
      },
      {
        id: 'task-3',
        name: '部署服务',
        type: 'deploy',
        step: 'pending',

        status: 'pending',
        dependencies: ['task-2'],
        appId: 'user-service',
      },
      {
        id: 'task-4',
        name: '健康检查',
        type: 'health_check',
        step: 'pending',

        status: 'pending',
        dependencies: ['task-3'],
        appId: 'user-service',
      },
    ],
    logs: [
      { timestamp: '2024-10-17T14:00:00Z', level: 'info', message: '部署创建，等待确认' },
    ],
  },
  'v1.2.0': {
    id: 'v1.2.0',
    versionId: 'v1.2.0',
    version: 'v1.2.0',
    applicationIds: ['user-service'],
    applications: ['user-service'],
    environmentIds: ['env1'],
    environments: ['production'],

    status: 'pending',
    progress: 0,
    createdAt: '2024-10-17T14:00:00Z',
    updatedAt: '2024-10-17T14:00:00Z',
    requireConfirm: true,
    grayscaleEnabled: false,
    tasks: [
      {
        id: 'task-1',
        name: '准备部署',
        type: 'prepare',
        step: 'pending',

        status: 'pending',
      },
      {
        id: 'task-2',
        name: '构建镜像',
        type: 'build',
        step: 'pending',

        status: 'pending',
      },
      {
        id: 'task-3',
        name: '部署服务',
        type: 'deploy',
        step: 'pending',

        status: 'pending',
      },
      {
        id: 'task-4',
        name: '健康检查',
        type: 'health_check',
        step: 'pending',

        status: 'pending',
      },
    ],
    logs: [
      { timestamp: '2024-10-17T14:00:00Z', level: 'info', message: '部署创建，等待确认' },
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
  { date: '2024-10-14', count: 6, successCount: 5, failedCount: 1 },
  { date: '2024-10-15', count: 3, successCount: 3, failedCount: 0 },
  { date: '2024-10-16', count: 9, successCount: 7, failedCount: 2 },
  { date: '2024-10-17', count: 5, successCount: 4, failedCount: 1 },
  { date: '2024-10-18', count: 7, successCount: 6, failedCount: 1 },
  { date: '2024-10-19', count: 5, successCount: 5, failedCount: 0 },
  { date: '2024-10-20', count: 8, successCount: 8, failedCount: 0 },
]
