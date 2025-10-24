import axios from 'axios'
import type {
  Version,
  Application,
  Environment,
  Deployment,
  DeploymentDetail,
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
  (response) => response.data,
  (error) => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export const versionApi = {
  list: (params?: { search?: string; isRevert?: boolean }) =>
    api.get<any, Version[]>('/versions', { params }),
  
  get: (id: string) =>
    api.get<any, Version>(`/versions/${id}`),
  
  create: (data: Partial<Version>) =>
    api.post<any, Version>('/versions', data),
}

export const applicationApi = {
  list: () =>
    api.get<any, Application[]>('/applications'),
  
  get: (id: string) =>
    api.get<any, Application>(`/applications/${id}`),
  
  create: (data: Partial<Application>) =>
    api.post<any, Application>('/applications', data),
  
  update: (id: string, data: Partial<Application>) =>
    api.put<any, Application>(`/applications/${id}`, data),
}

export const environmentApi = {
  list: () =>
    api.get<any, Environment[]>('/environments'),
  
  get: (id: string) =>
    api.get<any, Environment>(`/environments/${id}`),
  
  create: (data: Partial<Environment>) =>
    api.post<any, Environment>('/environments', data),
}

export const deploymentApi = {
  list: (params?: {
    status?: string
    environmentId?: string
    applicationId?: string
    startDate?: string
    endDate?: string
  }) =>
    api.get<any, Deployment[]>('/deployments', { params }),
  
  get: (id: string) =>
    api.get<any, DeploymentDetail>(`/deployments/${id}`),
  
  create: (data: CreateDeploymentRequest) =>
    api.post<any, Deployment>('/deployments', data),
  
  confirm: (id: string, note?: string) =>
    api.put<any, Deployment>(`/deployments/${id}`, {
      action: 'confirm',
      note,
    }),
  
  rollback: (id: string, reason?: string) =>
    api.put<any, Deployment>(`/deployments/${id}`, {
      action: 'rollback',
      reason,
    }),
  
  cancel: (id: string) =>
    api.put<any, Deployment>(`/deployments/${id}`, {
      action: 'cancel',
    }),
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
