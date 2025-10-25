import type { AxiosInstance } from 'axios'
import {
  mockVersions,
  mockApplications,
  mockEnvironments,
  mockDeployments,
  mockDeploymentDetails,
  mockDashboardStats,
  mockDeploymentTrends,
} from './data'
import type { CreateDeploymentRequest } from '@/types'

export function setupMockHandlers(apiInstance: AxiosInstance) {
  const originalRequest = apiInstance.request.bind(apiInstance)

  apiInstance.request = function (config: any) {
    const method = config.method?.toUpperCase()
    const url = config.url

    console.log(`[Mock API] ${method} ${url}`)

    if (url?.startsWith('/api/v1/versions')) {
      if (method === 'GET' && url === '/api/v1/versions') {
        const search = config.params?.search?.toLowerCase()
        let results = mockVersions
        if (search) {
          results = mockVersions.filter(
            (v) =>
              v.version.toLowerCase().includes(search) ||
              v.message.toLowerCase().includes(search) ||
              v.author.toLowerCase().includes(search)
          )
        }
        return Promise.resolve({ data: results })
      }
      
      if (method === 'GET' && url.match(/^\/api\/v1\/versions\/.+$/)) {
        const version = url.split('/')[4]
        const found = mockVersions.find((v) => v.version === version)
        if (found) {
          return Promise.resolve({ data: found })
        }
        return Promise.reject({ response: { status: 404, data: { message: 'Version not found' } } })
      }
      
      if (method === 'POST' && url === '/api/v1/versions') {
        const newVersion = {
          id: String(mockVersions.length + 1),
          ...config.data,
          createdAt: new Date().toISOString(),
        }
        return Promise.resolve({ data: newVersion })
      }
    }

    if (url?.startsWith('/api/v1/applications')) {
      if (method === 'GET' && url === '/api/v1/applications') {
        return Promise.resolve({ data: mockApplications })
      }
      
      if (method === 'GET' && url.match(/^\/api\/v1\/applications\/.+$/)) {
        const id = url.split('/')[4]
        const found = mockApplications.find((a) => a.id === id)
        if (found) {
          return Promise.resolve({ data: found })
        }
        return Promise.reject({ response: { status: 404, data: { message: 'Application not found' } } })
      }
      
      if (method === 'POST' && url === '/api/v1/applications') {
        const newApp = {
          id: `app${mockApplications.length + 1}`,
          ...config.data,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        }
        return Promise.resolve({ data: newApp })
      }
      
      if (method === 'PUT' && url.match(/^\/api\/v1\/applications\/.+$/)) {
        const id = url.split('/')[4]
        const found = mockApplications.find((a) => a.id === id)
        if (found) {
          const updated = {
            ...found,
            ...config.data,
            updatedAt: new Date().toISOString(),
          }
          return Promise.resolve({ data: updated })
        }
        return Promise.reject({ response: { status: 404, data: { message: 'Application not found' } } })
      }
    }

    if (url?.startsWith('/api/v1/environments')) {
      if (method === 'GET' && url === '/api/v1/environments') {
        return Promise.resolve({ data: mockEnvironments })
      }
      
      if (method === 'GET' && url.match(/^\/api\/v1\/environments\/.+$/)) {
        const id = url.split('/')[4]
        const found = mockEnvironments.find((e) => e.id === id)
        if (found) {
          return Promise.resolve({ data: found })
        }
        return Promise.reject({ response: { status: 404, data: { message: 'Environment not found' } } })
      }
      
      if (method === 'POST' && url === '/api/v1/environments') {
        const newEnv = {
          id: `env${mockEnvironments.length + 1}`,
          ...config.data,
          createdAt: new Date().toISOString(),
        }
        return Promise.resolve({ data: newEnv })
      }
    }

    if (url?.startsWith('/api/v1/deployments')) {
      if (method === 'GET' && url === '/api/v1/deployments') {
        let results = mockDeployments
        const { status, environmentId, applicationId } = config.params || {}
        
        if (status) {
          results = results.filter((d) => d.status === status)
        }
        if (environmentId) {
          results = results.filter((d) => d.environments.includes(environmentId))
        }
        if (applicationId) {
          results = results.filter((d) => d.applications.includes(applicationId))
        }
        
        return Promise.resolve({ data: results })
      }
      
      if (method === 'GET' && url.match(/^\/api\/v1\/deployments\/.+$/)) {
        const id = url.split('/')[4]
        const found = mockDeploymentDetails[id]
        if (found) {
          return Promise.resolve({ data: found })
        }
        return Promise.reject({ response: { status: 404, data: { message: 'Deployment not found' } } })
      }
      
      if (method === 'POST' && url === '/api/v1/deployments') {
        const data = config.data as CreateDeploymentRequest
        const newDeployment = {
          id: `deploy${mockDeployments.length + 1}`,
          version: data.version,
          applications: data.applicationIds,
          environments: data.environmentIds,
          status: 'pending' as const,
          progress: 0,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          createdBy: '当前用户',
        }
        return Promise.resolve({ data: newDeployment })
      }
      
      if (method === 'PUT' && url.match(/^\/api\/v1\/deployments\/.+$/)) {
        const id = url.split('/')[4]
        const found = mockDeployments.find((d) => d.id === id)
        if (found) {
          const action = config.data?.action
          let newStatus = found.status
          
          if (action === 'confirm') {
            newStatus = 'running'
          } else if (action === 'rollback') {
            newStatus = 'rolling_back'
          } else if (action === 'cancel') {
            newStatus = 'cancelled'
          }
          
          const updated = {
            ...found,
            status: newStatus,
            updatedAt: new Date().toISOString(),
          }
          return Promise.resolve({ data: updated })
        }
        return Promise.reject({ response: { status: 404, data: { message: 'Deployment not found' } } })
      }
    }

    if (url?.startsWith('/api/v1/dashboard')) {
      if (url === '/api/v1/dashboard/stats') {
        return Promise.resolve({ data: mockDashboardStats })
      }
      
      if (url === '/api/v1/dashboard/trends') {
        const days = config.params?.days || 7
        const trends = mockDeploymentTrends.slice(-days)
        return Promise.resolve({ data: trends })
      }
      
      if (url === '/api/v1/dashboard/recent-deployments') {
        const limit = config.params?.limit || 10
        const recent = mockDeployments.slice(0, limit)
        return Promise.resolve({ data: recent })
      }
    }

    return originalRequest(config)
  }

  console.log('[Mock API] Mock handlers initialized')
}

export function disableMockHandlers(apiInstance: AxiosInstance) {
  console.log('[Mock API] Mock handlers disabled')
}
