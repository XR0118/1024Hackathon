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
    version: 'v1.2.0',
    git: { tag: 'v1.2.0' },
    applications: [
      {
        name: 'api-service',
        coverage: 85,
        health: 95,
        lastUpdatedAt: '2024-10-24T10:00:00Z',
      },
      {
        name: 'web-app',
        coverage: 90,
        health: 98,
        lastUpdatedAt: '2024-10-24T10:00:00Z',
      },
    ],
    createdAt: '2024-10-24T10:00:00Z',
  },
  {
    version: 'v1.1.0',
    git: { tag: 'v1.1.0' },
    applications: [
      {
        name: 'api-service',
        coverage: 80,
        health: 92,
        lastUpdatedAt: '2024-10-20T10:00:00Z',
      },
    ],
    createdAt: '2024-10-20T10:00:00Z',
  },
]

export const mockApplications: Application[] = [
  {
    name: 'api-service',
    description: 'Backend API Service',
    icon: '🚀',
    versions: [
      {
        version: 'v1.2.0',
        status: 'normal',
        health: 95,
        lastUpdatedAt: '2024-10-24T10:00:00Z',
        nodes: [
          {
            name: 'node-1',
            health: 95,
            lastUpdatedAt: '2024-10-24T10:00:00Z',
          },
          {
            name: 'node-2',
            health: 96,
            lastUpdatedAt: '2024-10-24T10:00:00Z',
          },
        ],
      },
      {
        version: 'v1.1.0',
        status: 'revert',
        health: 92,
        lastUpdatedAt: '2024-10-20T10:00:00Z',
      },
    ],
  },
  {
    name: 'web-app',
    description: 'Frontend Web Application',
    icon: '🌐',
    versions: [
      {
        version: 'v1.2.0',
        status: 'normal',
        health: 98,
        lastUpdatedAt: '2024-10-24T10:00:00Z',
      },
    ],
  },
]

export const mockEnvironments: Environment[] = [
  {
    id: 'env-1',
    name: 'Production',
    type: 'k8s',
    status: 'active',
    applicationCount: 5,
  },
  {
    id: 'env-2',
    name: 'Staging',
    type: 'k8s',
    status: 'active',
    applicationCount: 5,
  },
  {
    id: 'env-3',
    name: 'Testing',
    type: 'physical',
    status: 'inactive',
    applicationCount: 3,
  },
]

export const mockDeployments: Deployment[] = [
  {
    id: 'deploy-1',
    versionId: 'v1.2.0',
    version: 'v1.2.0',
    applicationIds: ['app-1', 'app-2'],
    applications: ['api-service', 'web-app'],
    environmentIds: ['env-1'],
    environments: ['Production'],
    status: 'success',
    progress: 100,
    createdAt: '2024-10-24T10:00:00Z',
    updatedAt: '2024-10-24T10:15:00Z',
    duration: 900,
    requireConfirm: true,
    grayscaleEnabled: false,
  },
  {
    id: 'deploy-2',
    versionId: 'v1.2.0',
    version: 'v1.2.0',
    applicationIds: ['app-1'],
    applications: ['api-service'],
    environmentIds: ['env-2'],
    environments: ['Staging'],
    status: 'running',
    progress: 65,
    createdAt: '2024-10-24T11:00:00Z',
    updatedAt: '2024-10-24T11:05:00Z',
    requireConfirm: false,
    grayscaleEnabled: true,
    grayscaleRatio: 20,
  },
  {
    id: 'deploy-3',
    versionId: 'v1.1.0',
    version: 'v1.1.0',
    applicationIds: ['app-1'],
    applications: ['api-service'],
    environmentIds: ['env-1'],
    environments: ['Production'],
    status: 'waiting_confirm',
    progress: 100,
    createdAt: '2024-10-23T14:00:00Z',
    updatedAt: '2024-10-23T14:10:00Z',
    duration: 600,
    requireConfirm: true,
    grayscaleEnabled: false,
  },
]

export const mockDeploymentDetail: DeploymentDetail = {
  ...mockDeployments[0],
  steps: [
    {
      id: 'step-1',
      name: 'Pre-deployment checks',
      status: 'success',
      duration: 120,
      logs: ['Checking dependencies...', 'All checks passed'],
    },
    {
      id: 'step-2',
      name: 'Build and package',
      status: 'success',
      duration: 300,
      logs: ['Building application...', 'Build completed successfully'],
    },
    {
      id: 'step-3',
      name: 'Deploy to environment',
      status: 'success',
      duration: 480,
      logs: ['Deploying to production...', 'Deployment completed'],
    },
  ],
  logs: [
    {
      timestamp: '2024-10-24T10:00:00Z',
      level: 'info',
      message: 'Starting deployment process',
    },
    {
      timestamp: '2024-10-24T10:05:00Z',
      level: 'info',
      message: 'Pre-deployment checks completed',
    },
    {
      timestamp: '2024-10-24T10:10:00Z',
      level: 'info',
      message: 'Build completed successfully',
    },
    {
      timestamp: '2024-10-24T10:15:00Z',
      level: 'info',
      message: 'Deployment completed successfully',
    },
  ],
}

export const mockDashboardStats: DashboardStats = {
  activeVersions: 3,
  runningDeployments: 2,
  totalApplications: 5,
  totalEnvironments: 3,
}

export const mockDeploymentTrends: DeploymentTrend[] = [
  {
    date: '2024-10-18',
    count: 5,
    successCount: 4,
    failedCount: 1,
  },
  {
    date: '2024-10-19',
    count: 3,
    successCount: 3,
    failedCount: 0,
  },
  {
    date: '2024-10-20',
    count: 4,
    successCount: 3,
    failedCount: 1,
  },
  {
    date: '2024-10-21',
    count: 6,
    successCount: 5,
    failedCount: 1,
  },
  {
    date: '2024-10-22',
    count: 2,
    successCount: 2,
    failedCount: 0,
  },
  {
    date: '2024-10-23',
    count: 7,
    successCount: 6,
    failedCount: 1,
  },
  {
    date: '2024-10-24',
    count: 4,
    successCount: 4,
    failedCount: 0,
  },
]
