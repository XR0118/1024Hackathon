import axios from 'axios'
import type {
  Version,
  Application,
  ApplicationVersionsResponse,
  Environment,
  Deployment,
  DeploymentDetail,
  Task,
  DashboardStats,
  DeploymentTrend,
  CreateDeploymentRequest,
} from '@/types'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.response.use(
  (response) => {
    // 处理 mock 响应
    if (response.config?.headers?.['X-Mock-Response']) {
      return response.data
    }

    // 处理后端统一响应格式: { code: 0, message: "success", data: {...} }
    const data = response.data
    if (data && typeof data === 'object' && 'code' in data) {
      if (data.code === 0) {
        // 成功：返回 data 字段（如果存在），否则返回整个响应
        return data.data !== undefined ? data.data : data
      } else {
        // 业务错误
        return Promise.reject(new Error(data.message || 'Unknown error'))
      }
    }

    // 直接返回原始数据（如果不是标准格式）
    return data
  },
  (error) => {
    if (error.isMockResponse) {
      return Promise.reject(error)
    }

    console.error('API Error:', error)

    // 处理错误响应
    if (error.response?.data) {
      const errData = error.response.data
      // 处理后端错误格式: { error: { code: "...", message: "..." } }
      if (errData.error && errData.error.message) {
        return Promise.reject(new Error(errData.error.message))
      }
      // 处理统一响应格式的错误
      if (errData.message) {
        return Promise.reject(new Error(errData.message))
      }
    }

    return Promise.reject(error)
  }
)

export const versionApi = {
  list: (params?: { repository?: string; page?: number; page_size?: number }) =>
    api.get<any, { versions: Version[]; total: number; page: number; page_size: number }>('/versions', { params }).then(res => res.versions || []),

  get: (version: string) =>
    api.get<any, Version>(`/versions/${version}`),

  create: (data: Partial<Version>) =>
    api.post<any, Version>('/versions', data),

  delete: (version: string) =>
    api.delete<any, void>(`/versions/${version}`),

  rollback: (version: string, reason: string) =>
    api.post<any, any>(`/versions/${version}/rollback`, { reason }),
}

export const applicationApi = {
  list: async (params?: { repository?: string; type?: string }) => {
    const response = await api.get<any, any>('/applications', { params })
    // 后端返回分页格式: {applications: [...], total, page, page_size}
    // 提取 applications 数组
    if (response && response.applications) {
      return response.applications as Application[]
    }
    return response as Application[]
  },

  get: (id: string) =>
    api.get<any, Application>(`/applications/${id}`),

  // 获取应用的版本信息（从 Operator 查询，使用应用名称）
  getVersions: (name: string) =>
    api.get<any, ApplicationVersionsResponse>(`/applications/${name}/versions`),

  create: (data: Partial<Application>) =>
    api.post<any, Application>('/applications', data),

  update: (id: string, data: Partial<Application>) =>
    api.put<any, Application>(`/applications/${id}`, data),

  delete: (id: string) =>
    api.delete<any, void>(`/applications/${id}`),
}

export const environmentApi = {
  list: async () => {
    const response = await api.get<any, any>('/environments')
    // 后端返回分页格式: {environments: [...], total, page, page_size}
    // 提取 environments 数组
    if (response && response.environments) {
      return response.environments as Environment[]
    }
    return response as Environment[]
  },

  get: (id: string) =>
    api.get<any, Environment>(`/environments/${id}`),

  create: (data: Partial<Environment>) =>
    api.post<any, Environment>('/environments', data),

  update: (id: string, data: Partial<Environment>) =>
    api.put<any, Environment>(`/environments/${id}`, data),

  delete: (id: string) =>
    api.delete<any, void>(`/environments/${id}`),
}

export const deploymentApi = {
  list: (params?: {
    status?: string
    environment_id?: string
    version_id?: string
    page?: number
    page_size?: number
  }) =>
    api.get<any, { deployments: Deployment[]; total: number; page: number; page_size: number }>('/deployments', { params })
      .then(res => res.deployments || []),

  get: (id: string) =>
    api.get<any, DeploymentDetail>(`/deployments/${id}`),

  create: (data: CreateDeploymentRequest) =>
    api.post<any, Deployment>('/deployments', data),

  start: (id: string) =>
    api.post<any, Deployment>(`/deployments/${id}/start`, {}),

  pause: (id: string) =>
    api.post<any, Deployment>(`/deployments/${id}/pause`, {}),

  resume: (id: string) =>
    api.post<any, Deployment>(`/deployments/${id}/resume`, {}),
}

export const taskApi = {
  list: (params?: {
    deployment_id?: string
    status?: string
    type?: string
    page?: number
    page_size?: number
  }) =>
    api.get<any, { tasks: Task[]; total: number; page: number; page_size: number }>('/tasks', { params })
      .then(res => res.tasks || []),

  get: (id: string) =>
    api.get<any, Task>(`/tasks/${id}`),

  retry: (id: string) =>
    api.post<any, Task>(`/tasks/${id}/retry`, {}),
}

export const dashboardApi = {
  getStats: () =>
    api.get<any, DashboardStats>('/dashboard/stats'),

  getTrends: (days: number = 7) =>
    api.get<any, DeploymentTrend[]>('/dashboard/trends', {
      params: { days },
    }),

  getRecentDeployments: (limit: number = 10) =>
    api.get<any, Deployment[]>('/dashboard/recent-deployments', {
      params: { limit },
    }),
}

export default api
