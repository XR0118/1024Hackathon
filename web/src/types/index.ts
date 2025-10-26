// ==================== 健康度相关类型 ====================

// 健康状态信息
export interface HealthInfo {
  level: number                       // 健康度等级 (0-100)
  msg?: string                        // 健康状态描述信息
}

// ==================== 基础类型 ====================

export interface GitInfo {
  tag: string
  commit: string
  repository: string
}

export interface AppBuild {
  app_name: string
  docker_image: string
}

export interface Version {
  id: string
  version: string                 // 版本号（唯一标识）
  git_tag: string
  git_commit: string
  repository: string
  status: 'normal' | 'revert'
  created_by: string
  created_at: string
  description: string
  app_builds?: AppBuild[]
}

export interface Application {
  id: string                          // 匹配后端：应用 ID
  name: string                        // 应用名称
  description: string                 // 应用描述
  repository: string                  // Git 仓库地址
  type: 'microservice' | 'monolith'   // 应用类型
  config?: Record<string, string>     // 可选配置
  created_at?: string                 // 创建时间
  updated_at?: string                 // 更新时间
  environments?: Environment[]        // 关联的运行环境列表
}

// ==================== 版本概要相关类型 ====================

// 版本概要信息（只包含核心运行时指标）
export interface VersionSummary {
  version: string                     // 版本号
  status: 'normal' | 'revert'         // 版本状态
  healthy: HealthInfo                 // 健康度 (0-100)
  coverage_percent: number            // 覆盖度百分比 (0-100)
}

// 应用版本概要响应
export interface ApplicationVersionsSummaryResponse {
  application_id: string
  application_name: string
  versions: VersionSummary[]
}

// ==================== 版本详情相关类型 ====================

// 版本实例信息
export interface VersionInstance {
  node_name: string                   // 节点名称
  healthy: HealthInfo                 // 健康度 (0-100)
  status: string                      // 实例状态
  last_updated_at: string             // 最后更新时间
}

// 环境下的版本详细信息
export interface EnvironmentVersionDetail {
  version: string                     // 版本号
  status: 'normal' | 'abnormal' | 'revert'  // 版本状态
  git_tag: string                     // Git 标签
  git_commit: string                  // Git 提交哈希
  instances: VersionInstance[]        // 实例列表
  healthy: HealthInfo                 // 该版本在此环境的健康度 (0-100)
  coverage: number                    // 该版本在此环境的覆盖率(%)
  last_updated_at: string             // 最后更新时间
}

// 环境维度的版本信息
export interface EnvironmentVersions {
  environment: Environment            // 环境信息
  versions: EnvironmentVersionDetail[] // 该环境下的版本列表
}

// 应用版本详情响应（按环境组织）
export interface ApplicationVersionsDetailResponse {
  application_id: string
  application_name: string
  environments: EnvironmentVersions[]
}

// ==================== 已废弃的类型（保留用于向后兼容） ====================

/**
 * @deprecated 请使用 EnvironmentVersionDetail
 */
export interface ApplicationVersionInfo {
  version: string
  status: 'normal' | 'revert'
  health: number
  coverage: number
  last_updated_at: string
  nodes?: ApplicationNodeInfo[]
}

/**
 * @deprecated 请使用 VersionInstance
 */
export interface ApplicationNodeInfo {
  name: string
  health: number
  last_updated_at: string
}

/**
 * @deprecated 请使用 ApplicationVersionsDetailResponse
 */
export interface ApplicationVersionsResponse {
  application_id: string
  name: string
  versions: ApplicationVersionInfo[]
}

export interface Environment {
  id: string
  name: string
  type: 'kubernetes' | 'physical'  // 匹配后端：kubernetes 或 physical
  is_active: boolean               // 匹配后端：布尔类型
  config?: Record<string, string>  // 可选配置
  created_at?: string             // 创建时间
  updated_at?: string             // 更新时间
}

export interface Deployment {
  id: string
  version_id: string
  version?: string // 从 Version 关联获取
  must_in_order?: string[] // 应用部署顺序（应用名称数组）
  environment_id: string
  environment?: Environment // 从 Environment 关联获取
  status: 'pending' | 'running' | 'paused' | 'completed'
  created_by: string
  created_at: string
  updated_at: string
  started_at?: string
  completed_at?: string
  error_message?: string
  manual_approval?: boolean
  strategy?: any // DeploySteps[]
}

export interface DeploymentDetail extends Deployment {
  tasks: Task[]
  logs: DeploymentLog[]
}

export interface Task {
  id: string
  deployment_id: string
  app_id: string // 应用名称
  name: string // 任务名称
  type: 'build' | 'sleep' | 'deploy' | 'test' | 'approval' // 任务类型：构建/等待/部署/测试/复核
  step: 'pending' | 'running' | 'blocked' | 'completed' // workflow 执行中的状态
  status: 'pending' | 'running' | 'success' | 'failed' // 最终结果状态
  dependencies?: string[] // 上游依赖任务ID列表（DAG结构）
  payload?: Record<string, any> | BuildTaskPayload | SleepTaskPayload | DeployTaskPayload | TestTaskPayload | ApprovalTaskPayload // 任务参数（通用结构体）
  result?: Record<string, any> | BuildTaskResult | SleepTaskResult | DeployTaskResult | TestTaskResult | ApprovalTaskResult // 任务结果（通用结构体）
  created_at: string
  updated_at: string
  started_at?: string
  completed_at?: string
  deployment?: Deployment // 关联的部署
  application?: Application // 关联的应用
}

// ==================== Task Payload 示例 ====================

// Build Task Payload
export interface BuildTaskPayload {
  dockerfile?: string
  context?: string
  build_args?: Record<string, string>
  target_image: string
}

// Build Task Result
export interface BuildTaskResult {
  image: string
  image_id: string
  size: number
  build_duration: number
  logs?: string[]
}

// Sleep Task Payload
export interface SleepTaskPayload {
  duration: number // 等待时间（秒）
  reason?: string // 等待原因
}

// Sleep Task Result
export interface SleepTaskResult {
  actual_duration: number // 实际等待时间（秒）
}

// Deploy Task Payload
export interface DeployTaskPayload {
  image: string
  replicas: number
  strategy: 'rolling' | 'blue-green' | 'canary'
  canary_ratio?: number
  health_check?: {
    endpoint: string
    interval: number
    timeout: number
  }
}

// Deploy Task Result
export interface DeployTaskResult {
  deployed_instances: number
  healthy_instances: number
  rollout_duration: number
  endpoints?: string[]
}

// Test Task Payload
export interface TestTaskPayload {
  test_suite: string
  test_cases?: string[]
  environment: string
  timeout?: number
}

// Test Task Result
export interface TestTaskResult {
  passed: number
  failed: number
  skipped: number
  duration: number
  coverage?: number
  failures?: Array<{
    test: string
    message: string
  }>
}

// Approval Task Payload
export interface ApprovalTaskPayload {
  note: string // 审批说明
  required_approvers?: string[] // 需要审批的人员
  auto_approve_after?: number // 自动批准的超时时间（秒）
}

// Approval Task Result
export interface ApprovalTaskResult {
  approved: boolean
  approver?: string // 审批人
  approved_at?: string
  rejection_reason?: string
  approval_state?: 'waiting' | 'approved' | 'rejected' | 'timeout' // 业务状态在 result 中
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
  version_id: string
  must_in_order?: string[] // 应用部署顺序（应用名称数组）
  environment_id: string
  manual_approval?: boolean
  strategy?: any[] // DeploySteps[]
}
