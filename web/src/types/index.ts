export interface GitInfo {
  tag: string
}

export interface VersionApplication {
  name: string
  coverage: number
  health: number
  lastUpdatedAt: string
}

export interface Version {
  version: string
  git: GitInfo
  applications: VersionApplication[]
  createdAt: string
}

export interface ApplicationVersionInfo {
  version: string
  status: 'normal' | 'revert'
  health: number
  lastUpdatedAt: string
  nodes?: ApplicationNodeInfo[]
}

export interface ApplicationNodeInfo {
  name: string
  health: number
  lastUpdatedAt: string
}

export interface Application {
  name: string
  description: string
  icon?: string
  versions: ApplicationVersionInfo[]
}

export interface Environment {
  id: string
  name: string
  type: 'k8s' | 'physical'
  status: 'active' | 'inactive'
  applicationCount: number
}

export interface Deployment {
  id: string
  versionId: string
  version: string
  applicationIds: string[]
  applications: string[]
  environmentIds: string[]
  environments: string[]
  status: 'pending' | 'running' | 'success' | 'failed' | 'waiting_confirm'
  progress: number
  createdAt: string
  updatedAt: string
  duration?: number
  requireConfirm: boolean
  grayscaleEnabled: boolean
  grayscaleRatio?: number
}

export interface DeploymentDetail extends Deployment {
  steps: DeploymentStep[]
  logs: DeploymentLog[]
}

export interface DeploymentStep {
  id: string
  name: string
  status: 'pending' | 'running' | 'success' | 'failed'
  duration?: number
  logs?: string[]
}

export interface DeploymentLog {
  timestamp: string
  level: 'info' | 'warn' | 'error'
  message: string
}

export interface DashboardStats {
  activeVersions: number
  runningDeployments: number
  totalApplications: number
  totalEnvironments: number
}

export interface DeploymentTrend {
  date: string
  count: number
  successCount: number
  failedCount: number
}

export interface CreateDeploymentRequest {
  versionId: string
  applicationIds: string[]
  environmentIds: string[]
  requireConfirm: boolean
  grayscaleEnabled: boolean
  grayscaleRatio?: number
  autoRollback: boolean
}
