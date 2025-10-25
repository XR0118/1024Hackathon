import type { AxiosInstance, InternalAxiosRequestConfig } from 'axios'
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
  apiInstance.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      const method = config.method?.toUpperCase()
      const url = config.url || ''

      console.log(`[Mock API] ${method} ${url}`)

      if (url.startsWith('/versions')) {
        if (method === 'GET' && url === '/versions') {
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
          return Promise.reject({
            config,
            response: { data: results, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (method === 'GET' && url.match(/^\/versions\/.+$/)) {
          const version = url.split('/')[2]
          const found = mockVersions.find((v) => v.version === version)
          if (found) {
            return Promise.reject({
              config,
              response: { data: found, status: 200, statusText: 'OK', headers: {}, config },
              isMockResponse: true,
            })
          }
          return Promise.reject({
            config,
            response: { data: { message: 'Version not found' }, status: 404, statusText: 'Not Found', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (method === 'POST' && url === '/versions') {
          const newVersion = {
            id: String(mockVersions.length + 1),
            ...config.data,
            createdAt: new Date().toISOString(),
          }
          return Promise.reject({
            config,
            response: { data: newVersion, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }
      }

      if (url.startsWith('/applications')) {
        if (method === 'GET' && url === '/applications') {
          return Promise.reject({
            config,
            response: { data: mockApplications, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (method === 'GET' && url.match(/^\/applications\/.+$/)) {
          const id = url.split('/')[2]
          const found = mockApplications.find((a) => a.id === id)
          if (found) {
            return Promise.reject({
              config,
              response: { data: found, status: 200, statusText: 'OK', headers: {}, config },
              isMockResponse: true,
            })
          }
          return Promise.reject({
            config,
            response: { data: { message: 'Application not found' }, status: 404, statusText: 'Not Found', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (method === 'POST' && url === '/applications') {
          const newApp = {
            id: `app${mockApplications.length + 1}`,
            ...config.data,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          }
          return Promise.reject({
            config,
            response: { data: newApp, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (method === 'PUT' && url.match(/^\/applications\/.+$/)) {
          const id = url.split('/')[2]
          const found = mockApplications.find((a) => a.id === id)
          if (found) {
            const updated = {
              ...found,
              ...config.data,
              updatedAt: new Date().toISOString(),
            }
            return Promise.reject({
              config,
              response: { data: updated, status: 200, statusText: 'OK', headers: {}, config },
              isMockResponse: true,
            })
          }
          return Promise.reject({
            config,
            response: { data: { message: 'Application not found' }, status: 404, statusText: 'Not Found', headers: {}, config },
            isMockResponse: true,
          })
        }
      }

      if (url.startsWith('/environments')) {
        if (method === 'GET' && url === '/environments') {
          return Promise.reject({
            config,
            response: { data: mockEnvironments, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (method === 'GET' && url.match(/^\/environments\/.+$/)) {
          const id = url.split('/')[2]
          const found = mockEnvironments.find((e) => e.id === id)
          if (found) {
            return Promise.reject({
              config,
              response: { data: found, status: 200, statusText: 'OK', headers: {}, config },
              isMockResponse: true,
            })
          }
          return Promise.reject({
            config,
            response: { data: { message: 'Environment not found' }, status: 404, statusText: 'Not Found', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (method === 'POST' && url === '/environments') {
          const newEnv = {
            id: `env${mockEnvironments.length + 1}`,
            ...config.data,
            createdAt: new Date().toISOString(),
          }
          return Promise.reject({
            config,
            response: { data: newEnv, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }
      }

      if (url.startsWith('/deployments')) {
        if (method === 'GET' && url === '/deployments') {
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
          
          return Promise.reject({
            config,
            response: { data: results, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (method === 'GET' && url.match(/^\/deployments\/.+$/)) {
          const id = url.split('/')[2]
          const found = mockDeploymentDetails[id]
          if (found) {
            return Promise.reject({
              config,
              response: { data: found, status: 200, statusText: 'OK', headers: {}, config },
              isMockResponse: true,
            })
          }
          return Promise.reject({
            config,
            response: { data: { message: 'Deployment not found' }, status: 404, statusText: 'Not Found', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (method === 'POST' && url === '/deployments') {
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
          return Promise.reject({
            config,
            response: { data: newDeployment, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (method === 'PUT' && url.match(/^\/deployments\/.+$/)) {
          const id = url.split('/')[2]
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
            return Promise.reject({
              config,
              response: { data: updated, status: 200, statusText: 'OK', headers: {}, config },
              isMockResponse: true,
            })
          }
          return Promise.reject({
            config,
            response: { data: { message: 'Deployment not found' }, status: 404, statusText: 'Not Found', headers: {}, config },
            isMockResponse: true,
          })
        }
      }

      if (url.startsWith('/dashboard')) {
        if (url.includes('/dashboard/stats')) {
          return Promise.reject({
            config,
            response: { data: mockDashboardStats, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (url.includes('/dashboard/trends')) {
          const days = config.params?.days || 7
          const trends = mockDeploymentTrends.slice(-days)
          return Promise.reject({
            config,
            response: { data: trends, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }
        
        if (url.includes('/dashboard/recent-deployments')) {
          const limit = config.params?.limit || 10
          const recent = mockDeployments.slice(0, limit)
          return Promise.reject({
            config,
            response: { data: recent, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }
      }

      return config
    },
    (error) => Promise.reject(error)
  )

  apiInstance.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.isMockResponse && error.response) {
        if (error.response.status >= 200 && error.response.status < 300) {
          return Promise.resolve(error.response)
        }
        return Promise.reject(error)
      }
      return Promise.reject(error)
    }
  )

  console.log('[Mock API] Mock handlers initialized')
}

export function disableMockHandlers(apiInstance: AxiosInstance) {
  console.log('[Mock API] Mock handlers disabled')
}
