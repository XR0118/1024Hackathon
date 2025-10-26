import type { AxiosInstance, InternalAxiosRequestConfig } from 'axios'
import {
  mockVersions,
  mockApplications,
  mockApplicationVersionsSummary,
  mockApplicationVersionsDetail,
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
          const repository = config.params?.repository?.toLowerCase()
          let results = mockVersions
          if (repository) {
            results = mockVersions.filter((v) =>
              v.repository.toLowerCase().includes(repository) ||
              v.git_tag.toLowerCase().includes(repository)
            )
          }
          // 返回分页格式
          return Promise.reject({
            config,
            response: {
              data: {
                versions: results,
                total: results.length,
                page: config.params?.page || 1,
                page_size: config.params?.page_size || 100,
              },
              status: 200,
              statusText: 'OK',
              headers: {},
              config,
            },
            isMockResponse: true,
          })
        }

        if (method === 'GET' && url.match(/^\/versions\/[^/]+$/)) {
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
            id: `version-${Date.now()}`,
            ...config.data,
            status: 'normal',
            created_at: new Date().toISOString(),
          }
          return Promise.reject({
            config,
            response: { data: newVersion, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }

        if (method === 'DELETE' && url.match(/^\/versions\/[^/]+$/)) {
          const version = url.split('/')[2]
          const found = mockVersions.find((v) => v.version === version)
          if (found) {
            return Promise.reject({
              config,
              response: { data: { message: 'Version deleted successfully' }, status: 200, statusText: 'OK', headers: {}, config },
              isMockResponse: true,
            })
          }
          return Promise.reject({
            config,
            response: { data: { message: 'Version not found' }, status: 404, statusText: 'Not Found', headers: {}, config },
            isMockResponse: true,
          })
        }

        if (method === 'POST' && url.match(/^\/versions\/[^/]+\/rollback$/)) {
          const version = url.split('/')[2]
          const found = mockVersions.find((v) => v.version === version)
          if (found) {
            // 模拟回滚操作，返回成功消息
            return Promise.reject({
              config,
              response: {
                data: {
                  message: 'Rollback initiated',
                  version: version,
                  reason: config.data?.reason,
                },
                status: 200,
                statusText: 'OK',
                headers: {},
                config,
              },
              isMockResponse: true,
            })
          }
          return Promise.reject({
            config,
            response: { data: { message: 'Version not found' }, status: 404, statusText: 'Not Found', headers: {}, config },
            isMockResponse: true,
          })
        }
      }

      if (url.startsWith('/applications')) {
        // 应用列表
        if (method === 'GET' && url === '/applications') {
          return Promise.reject({
            config,
            response: { data: { applications: mockApplications }, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }

        // 应用版本概要（匹配 /applications/:name/versions/summary）
        if (method === 'GET' && url.match(/^\/applications\/[^/]+\/versions\/summary$/)) {
          const name = url.split('/')[2]
          const summaryInfo = mockApplicationVersionsSummary[name]
          if (summaryInfo) {
            return Promise.reject({
              config,
              response: { data: summaryInfo, status: 200, statusText: 'OK', headers: {}, config },
              isMockResponse: true,
            })
          }
          // 如果没有版本概要信息，返回空列表
          return Promise.reject({
            config,
            response: {
              data: {
                application_id: '',
                application_name: name,
                versions: [],
              },
              status: 200,
              statusText: 'OK',
              headers: {},
              config,
            },
            isMockResponse: true,
          })
        }

        // 应用版本详情（匹配 /applications/:name/versions）
        if (method === 'GET' && url.match(/^\/applications\/[^/]+\/versions$/)) {
          const name = url.split('/')[2]
          const detailInfo = mockApplicationVersionsDetail[name]
          if (detailInfo) {
            return Promise.reject({
              config,
              response: { data: detailInfo, status: 200, statusText: 'OK', headers: {}, config },
              isMockResponse: true,
            })
          }
          // 如果没有版本详情信息，返回空结构
          return Promise.reject({
            config,
            response: {
              data: {
                application_id: '',
                application_name: name,
                environments: [],
              },
              status: 200,
              statusText: 'OK',
              headers: {},
              config,
            },
            isMockResponse: true,
          })
        }

        // 应用详情（匹配 /applications/:id）
        if (method === 'GET' && url.match(/^\/applications\/[^/]+$/)) {
          const id = url.split('/')[2]
          const found = mockApplications.find((a) => a.id === id || a.name === id)
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

        // 创建应用
        if (method === 'POST' && url === '/applications') {
          const newApp = {
            id: `app-${Date.now()}`,
            ...config.data,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          }
          return Promise.reject({
            config,
            response: { data: newApp, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }

        // 更新应用
        if (method === 'PUT' && url.match(/^\/applications\/.+$/)) {
          const id = url.split('/')[2]
          const found = mockApplications.find((a) => a.id === id)
          if (found) {
            const updated = {
              ...found,
              ...config.data,
              updated_at: new Date().toISOString(),
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
            response: { data: { environments: mockEnvironments }, status: 200, statusText: 'OK', headers: {}, config },
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
          const { status, environment_id, version_id } = config.params || {}

          if (status) {
            const statusList = status.split(',').map((s: string) => s.trim())
            results = results.filter((d) => statusList.includes(d.status))
          }
          if (environment_id) {
            results = results.filter((d) => d.environment_id === environment_id)
          }
          if (version_id) {
            results = results.filter((d) => d.version_id === version_id)
          }

          // 如果查询的是 running 或 paused 状态，返回包含 tasks 的详细数据
          const needDetails = status && (status.includes('running') || status.includes('paused'))
          const deploymentsData = needDetails
            ? results.map(d => mockDeploymentDetails[d.id] || d)
            : results

          // 返回分页格式
          return Promise.reject({
            config,
            response: {
              data: {
                deployments: deploymentsData,
                total: deploymentsData.length,
                page: config.params?.page || 1,
                page_size: config.params?.page_size || 100,
              },
              status: 200,
              statusText: 'OK',
              headers: {},
              config,
            },
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
            id: `deploy-${Date.now()}`,
            version_id: data.version_id,
            must_in_order: data.must_in_order || [],
            environment_id: data.environment_id,
            status: 'pending' as const,
            step: 'pending' as const,
            created_by: 'admin',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            manual_approval: data.manual_approval || false,
            strategy: data.strategy,
          }
          return Promise.reject({
            config,
            response: { data: newDeployment, status: 200, statusText: 'OK', headers: {}, config },
            isMockResponse: true,
          })
        }

        if (method === 'POST' && url.match(/^\/deployments\/[^/]+\/start$/)) {
          const id = url.split('/')[2]
          const found = mockDeployments.find((d) => d.id === id)
          if (found) {
            const updated = {
              ...found,
              status: 'running',
              step: 'running',
              started_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            }
            return Promise.reject({
              config,
              response: { data: updated, status: 200, statusText: 'OK', headers: {}, config },
              isMockResponse: true,
            })
          }
        }

        if (method === 'POST' && url.match(/^\/deployments\/[^/]+\/pause$/)) {
          const id = url.split('/')[2]
          const found = mockDeployments.find((d) => d.id === id)
          if (found) {
            const updated = {
              ...found,
              status: 'paused',
              step: 'paused',
              updated_at: new Date().toISOString(),
            }
            return Promise.reject({
              config,
              response: { data: updated, status: 200, statusText: 'OK', headers: {}, config },
              isMockResponse: true,
            })
          }
        }

        if (method === 'POST' && url.match(/^\/deployments\/[^/]+\/resume$/)) {
          const id = url.split('/')[2]
          const found = mockDeployments.find((d) => d.id === id)
          if (found) {
            const updated = {
              ...found,
              status: 'running',
              step: 'running',
              updated_at: new Date().toISOString(),
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
          return error.response.data
        }
        return Promise.reject(error)
      }
      return Promise.reject(error)
    }
  )

  console.log('[Mock API] Mock handlers initialized')
}

export function disableMockHandlers(_apiInstance: AxiosInstance) {
  console.log('[Mock API] Mock handlers disabled')
}
