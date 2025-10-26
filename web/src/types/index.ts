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
  coverage: number
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
  status: 'pending' | 'running' | 'paused' | 'completed'
  progress: number
  createdAt: string
  updatedAt: string
  duration?: number
  requireConfirm: boolean
  grayscaleEnabled: boolean
  grayscaleRatio?: number
}

export interface DeploymentDetail extends Deployment {
  tasks: Task[]
  logs: DeploymentLog[]
}

export interface Task {
  id: string
  deploymentId?: string
  appId?: string
  name: string
  type: 'build' | 'test' | 'deploy' | 'health_check' | 'prepare' | 'sleep' | 'approval' | 'custom'
  status: 'pending' | 'running' | 'success' | 'failed' | 'blocked' | 'cancelled' | 'waiting_approval'
  dependencies?: string[]  // 上游依赖的任务ID列表，为空或不存在表示是顶点
  duration?: number
  startedAt?: string
  completedAt?: string
  logs?: string[]
  params?: TaskParams
}

export interface TaskParams {
  sleepDuration?: number  // sleep 任务的等待时间（秒）
  approvalNote?: string   // approval 任务的审批说明
  [key: string]: any      // 允许其他自定义参数
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
