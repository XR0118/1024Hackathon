import React, { useEffect, useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { applicationApi } from '@/services/api'
import { formatDate } from '@/utils'
import type { Application } from '@/types'
import { IconPlus, IconRocket } from '@tabler/icons-react'
import { useErrorStore } from '@/store/error'

const Applications: React.FC = () => {
  const navigate = useNavigate()
  const { setError } = useErrorStore();
  const [applications, setApplications] = useState<Application[]>([])

  const loadApplications = useCallback(async () => {
    try {
      const data = await applicationApi.list()
      setApplications(data)
    } catch (error) {
      setError('Failed to load applications.')
    }
  }, [setError])

  useEffect(() => {
    loadApplications()
  }, [loadApplications])

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title">应用管理</h2>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <button className="btn btn-primary">
              <IconPlus className="icon" />
              添加应用
            </button>
          </div>
        </div>
      </div>

      <div className="row row-cards">
        {applications.map((app) => (
          <div className="col-sm-6 col-lg-4" key={app.id}>
            <div className="card">
              <div className="card-body">
                <h3 className="card-title">{app.name}</h3>
                <p className="text-muted">{app.description}</p>
                <div>
                  <strong>当前版本:</strong>
                  <div className="mt-2">
                    {Object.entries(app.currentVersions).map(([env, version]) => (
                      <span className="badge bg-secondary-lt me-1" key={env}>
                        {env}: {version}
                      </span>
                    ))}
                  </div>
                </div>
                <p className="mt-3 text-muted" style={{ fontSize: '12px' }}>
                  最近部署: {formatDate(app.lastDeployedAt)}
                </p>
              </div>
              <div className="card-footer">
                <div className="d-flex">
                  <button
                    className="btn btn-ghost-primary"
                    onClick={() => navigate(`/deployments/new?appId=${app.id}`)}
                  >
                    <IconRocket size={16} className="me-2" />
                    新建部署
                  </button>
                  <button
                    className="btn btn-ghost-secondary ms-auto"
                    onClick={() => navigate(`/applications/${app.id}`)}
                  >
                    查看详情
                  </button>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

export default Applications
