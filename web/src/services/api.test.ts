import { describe, it, expect } from 'vitest'
import {
  mockVersions,
  mockApplications,
  mockEnvironments,
  mockDeployments,
  mockDeploymentDetail,
  mockDashboardStats,
  mockDeploymentTrends,
} from '@/test/mocks/api'

describe('Mock API Data', () => {
  describe('Version Mock Data', () => {
    it('has valid version structure', () => {
      expect(mockVersions).toBeDefined()
      expect(mockVersions.length).toBeGreaterThan(0)
      expect(mockVersions[0]).toHaveProperty('version')
      expect(mockVersions[0]).toHaveProperty('git')
      expect(mockVersions[0]).toHaveProperty('applications')
      expect(mockVersions[0]).toHaveProperty('createdAt')
    })

    it('has valid git info', () => {
      expect(mockVersions[0].git).toHaveProperty('tag')
      expect(mockVersions[0].git.tag).toBe('v1.2.0')
    })

    it('has valid application info in version', () => {
      expect(mockVersions[0].applications).toBeDefined()
      expect(mockVersions[0].applications.length).toBeGreaterThan(0)
      expect(mockVersions[0].applications[0]).toHaveProperty('name')
      expect(mockVersions[0].applications[0]).toHaveProperty('coverage')
      expect(mockVersions[0].applications[0]).toHaveProperty('health')
    })
  })

  describe('Application Mock Data', () => {
    it('has valid application structure', () => {
      expect(mockApplications).toBeDefined()
      expect(mockApplications.length).toBeGreaterThan(0)
      expect(mockApplications[0]).toHaveProperty('name')
      expect(mockApplications[0]).toHaveProperty('description')
      expect(mockApplications[0]).toHaveProperty('versions')
    })

    it('has valid version info in application', () => {
      expect(mockApplications[0].versions).toBeDefined()
      expect(mockApplications[0].versions.length).toBeGreaterThan(0)
      expect(mockApplications[0].versions[0]).toHaveProperty('version')
      expect(mockApplications[0].versions[0]).toHaveProperty('status')
      expect(mockApplications[0].versions[0]).toHaveProperty('health')
    })

    it('has valid node info when present', () => {
      const versionWithNodes = mockApplications[0].versions.find(v => v.nodes)
      if (versionWithNodes?.nodes) {
        expect(versionWithNodes.nodes.length).toBeGreaterThan(0)
        expect(versionWithNodes.nodes[0]).toHaveProperty('name')
        expect(versionWithNodes.nodes[0]).toHaveProperty('health')
      }
    })
  })

  describe('Environment Mock Data', () => {
    it('has valid environment structure', () => {
      expect(mockEnvironments).toBeDefined()
      expect(mockEnvironments.length).toBeGreaterThan(0)
      expect(mockEnvironments[0]).toHaveProperty('id')
      expect(mockEnvironments[0]).toHaveProperty('name')
      expect(mockEnvironments[0]).toHaveProperty('type')
      expect(mockEnvironments[0]).toHaveProperty('status')
      expect(mockEnvironments[0]).toHaveProperty('applicationCount')
    })

    it('has valid environment types', () => {
      const types = mockEnvironments.map(e => e.type)
      types.forEach(type => {
        expect(['k8s', 'physical']).toContain(type)
      })
    })

    it('has valid environment statuses', () => {
      const statuses = mockEnvironments.map(e => e.status)
      statuses.forEach(status => {
        expect(['active', 'inactive']).toContain(status)
      })
    })
  })

  describe('Deployment Mock Data', () => {
    it('has valid deployment structure', () => {
      expect(mockDeployments).toBeDefined()
      expect(mockDeployments.length).toBeGreaterThan(0)
      expect(mockDeployments[0]).toHaveProperty('id')
      expect(mockDeployments[0]).toHaveProperty('version')
      expect(mockDeployments[0]).toHaveProperty('status')
      expect(mockDeployments[0]).toHaveProperty('progress')
      expect(mockDeployments[0]).toHaveProperty('applications')
      expect(mockDeployments[0]).toHaveProperty('environments')
    })

    it('has valid deployment statuses', () => {
      const statuses = mockDeployments.map(d => d.status)
      statuses.forEach(status => {
        expect(['pending', 'running', 'success', 'failed', 'waiting_confirm']).toContain(status)
      })
    })

    it('has valid progress values', () => {
      mockDeployments.forEach(deployment => {
        expect(deployment.progress).toBeGreaterThanOrEqual(0)
        expect(deployment.progress).toBeLessThanOrEqual(100)
      })
    })
  })

  describe('Deployment Detail Mock Data', () => {
    it('has valid deployment detail structure', () => {
      expect(mockDeploymentDetail).toBeDefined()
      expect(mockDeploymentDetail).toHaveProperty('steps')
      expect(mockDeploymentDetail).toHaveProperty('logs')
    })

    it('has valid steps', () => {
      expect(mockDeploymentDetail.steps).toBeDefined()
      expect(mockDeploymentDetail.steps.length).toBeGreaterThan(0)
      expect(mockDeploymentDetail.steps[0]).toHaveProperty('id')
      expect(mockDeploymentDetail.steps[0]).toHaveProperty('name')
      expect(mockDeploymentDetail.steps[0]).toHaveProperty('status')
    })

    it('has valid logs', () => {
      expect(mockDeploymentDetail.logs).toBeDefined()
      expect(mockDeploymentDetail.logs.length).toBeGreaterThan(0)
      expect(mockDeploymentDetail.logs[0]).toHaveProperty('timestamp')
      expect(mockDeploymentDetail.logs[0]).toHaveProperty('level')
      expect(mockDeploymentDetail.logs[0]).toHaveProperty('message')
    })
  })

  describe('Dashboard Stats Mock Data', () => {
    it('has valid dashboard stats structure', () => {
      expect(mockDashboardStats).toBeDefined()
      expect(mockDashboardStats).toHaveProperty('activeVersions')
      expect(mockDashboardStats).toHaveProperty('runningDeployments')
      expect(mockDashboardStats).toHaveProperty('totalApplications')
      expect(mockDashboardStats).toHaveProperty('totalEnvironments')
    })

    it('has non-negative counts', () => {
      expect(mockDashboardStats.activeVersions).toBeGreaterThanOrEqual(0)
      expect(mockDashboardStats.runningDeployments).toBeGreaterThanOrEqual(0)
      expect(mockDashboardStats.totalApplications).toBeGreaterThanOrEqual(0)
      expect(mockDashboardStats.totalEnvironments).toBeGreaterThanOrEqual(0)
    })
  })

  describe('Deployment Trends Mock Data', () => {
    it('has valid deployment trends structure', () => {
      expect(mockDeploymentTrends).toBeDefined()
      expect(mockDeploymentTrends.length).toBeGreaterThan(0)
      expect(mockDeploymentTrends[0]).toHaveProperty('date')
      expect(mockDeploymentTrends[0]).toHaveProperty('count')
      expect(mockDeploymentTrends[0]).toHaveProperty('successCount')
      expect(mockDeploymentTrends[0]).toHaveProperty('failedCount')
    })

    it('has consistent trend counts', () => {
      mockDeploymentTrends.forEach(trend => {
        expect(trend.count).toBeGreaterThanOrEqual(trend.successCount + trend.failedCount)
      })
    })
  })
})
