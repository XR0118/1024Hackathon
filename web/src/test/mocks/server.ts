import { vi } from 'vitest'
import {
  mockVersions,
  mockApplications,
  mockEnvironments,
  mockDeployments,
  mockDeploymentDetail,
  mockDashboardStats,
  mockDeploymentTrends,
} from './api'

export const createMockApi = () => {
  return {
    get: vi.fn((url: string) => {
      if (url.includes('/versions/v')) {
        const version = url.split('/').pop()
        return Promise.resolve(mockVersions.find(v => v.version === version))
      }
      if (url.includes('/versions')) {
        return Promise.resolve(mockVersions)
      }
      if (url.includes('/applications/')) {
        const appName = url.split('/').pop()
        return Promise.resolve(mockApplications.find(a => a.name === appName))
      }
      if (url.includes('/applications')) {
        return Promise.resolve(mockApplications)
      }
      if (url.includes('/environments/')) {
        const envId = url.split('/').pop()
        return Promise.resolve(mockEnvironments.find(e => e.id === envId))
      }
      if (url.includes('/environments')) {
        return Promise.resolve(mockEnvironments)
      }
      if (url.includes('/deployments/')) {
        const deployId = url.split('/').pop()
        if (deployId === 'deploy-1') {
          return Promise.resolve(mockDeploymentDetail)
        }
        return Promise.resolve(mockDeployments.find(d => d.id === deployId))
      }
      if (url.includes('/dashboard/stats')) {
        return Promise.resolve(mockDashboardStats)
      }
      if (url.includes('/dashboard/trends')) {
        return Promise.resolve(mockDeploymentTrends)
      }
      if (url.includes('/dashboard/recent-deployments')) {
        return Promise.resolve(mockDeployments.slice(0, 2))
      }
      if (url.includes('/deployments')) {
        return Promise.resolve(mockDeployments)
      }
      return Promise.reject(new Error('Not found'))
    }),
    post: vi.fn((url: string, data: any) => {
      if (url.includes('/versions')) {
        return Promise.resolve({ ...data, createdAt: new Date().toISOString() })
      }
      if (url.includes('/applications')) {
        return Promise.resolve({ ...data, versions: [] })
      }
      if (url.includes('/environments')) {
        return Promise.resolve({ ...data, id: 'env-new', applicationCount: 0 })
      }
      if (url.includes('/deployments')) {
        return Promise.resolve({
          id: 'deploy-new',
          ...data,
          status: 'pending',
          progress: 0,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        })
      }
      return Promise.reject(new Error('Not found'))
    }),
    put: vi.fn((url: string, data: any) => {
      if (url.includes('/applications/')) {
        return Promise.resolve({ ...mockApplications[0], ...data })
      }
      if (url.includes('/deployments/')) {
        const deployment = mockDeployments[0]
        return Promise.resolve({ ...deployment, status: 'success' })
      }
      return Promise.reject(new Error('Not found'))
    }),
  }
}
