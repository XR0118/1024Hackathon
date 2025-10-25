import React, { useEffect, useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  versionApi,
  applicationApi,
  environmentApi,
  deploymentApi,
} from '@/services/api'
import type { Version, Application, Environment, CreateDeploymentRequest } from '@/types'
import { IconArrowLeft } from '@tabler/icons-react'
import { useErrorStore } from '@/store/error'

const CreateDeployment: React.FC = () => {
  const navigate = useNavigate()
  const { setError } = useErrorStore();
  const [currentStep, setCurrentStep] = useState(0)
  
  const [versions, setVersions] = useState<Version[]>([])
  const [applications, setApplications] = useState<Application[]>([])
  const [environments, setEnvironments] = useState<Environment[]>([])
  const [loading, setLoading] = useState(false)
  
  const [selectedVersion, setSelectedVersion] = useState<string>()
  const [selectedApps, setSelectedApps] = useState<string[]>([])
  const [selectedEnvs, setSelectedEnvs] = useState<string[]>([])
  const [requireConfirm, setRequireConfirm] = useState(false)
  const [grayscaleEnabled, setGrayscaleEnabled] = useState(false)
  const [grayscaleRatio, setGrayscaleRatio] = useState(50)
  const [autoRollback, setAutoRollback] = useState(true)

  const loadData = useCallback(async () => {
    try {
      const [versionsData, appsData, envsData] = await Promise.all([
        versionApi.list(),
        applicationApi.list(),
        environmentApi.list(),
      ])
      setVersions(versionsData)
      setApplications(appsData)
      setEnvironments(envsData)
    } catch (error) {
      setError('Failed to load data for creating deployment.')
    }
  }, [setError])

  useEffect(() => {
    loadData()
  }, [loadData])

  const handleNext = () => {
    setCurrentStep(currentStep + 1)
  }

  const handlePrev = () => {
    setCurrentStep(currentStep - 1)
  }

  const handleSubmit = async () => {
    if (!selectedVersion || selectedApps.length === 0 || selectedEnvs.length === 0) {
      return
    }

    setLoading(true)
    try {
      const request: CreateDeploymentRequest = {
        versionId: selectedVersion,
        applicationIds: selectedApps,
        environmentIds: selectedEnvs,
        requireConfirm,
        grayscaleEnabled,
        grayscaleRatio: grayscaleEnabled ? grayscaleRatio : undefined,
        autoRollback,
      }
      const deployment = await deploymentApi.create(request)
      navigate(`/deployments/${deployment.id}`)
    } catch (error) {
      setError('Failed to create deployment.')
    } finally {
      setLoading(false)
    }
  }

  const toggleSelection = (id: string, selected: string[], setSelected: React.Dispatch<React.SetStateAction<string[]>>) => {
    if (selected.includes(id)) {
      setSelected(selected.filter(item => item !== id))
    } else {
      setSelected([...selected, id])
    }
  }

  const steps = [
    {
      title: '选择版本',
      content: (
        <div className="list-group">
          {versions.map((version) => (
            <button
              key={version.version}
              type="button"
              className={`list-group-item list-group-item-action ${selectedVersion === version.version ? 'active' : ''}`}
              onClick={() => setSelectedVersion(version.version)}
            >
              <strong>{version.version}</strong>
              <small className="d-block text-muted">{version.git.tag} - {version.createdAt}</small>
            </button>
          ))}
        </div>
      ),
    },
    {
      title: '选择应用',
      content: (
        <div className="list-group">
          {applications.map((app) => (
            <button
              key={app.name}
              type="button"
              className={`list-group-item list-group-item-action ${selectedApps.includes(app.name) ? 'active' : ''}`}
              onClick={() => toggleSelection(app.name, selectedApps, setSelectedApps)}
            >
              {app.name}
            </button>
          ))}
        </div>
      ),
    },
    {
      title: '选择环境',
      content: (
        <div className="list-group">
          {environments.map((env) => (
            <button
              key={env.id}
              type="button"
              className={`list-group-item list-group-item-action ${selectedEnvs.includes(env.id) ? 'active' : ''}`}
              onClick={() => toggleSelection(env.id, selectedEnvs, setSelectedEnvs)}
            >
              {env.name} ({env.type})
            </button>
          ))}
        </div>
      ),
    },
    {
      title: '配置选项',
      content: (
        <form>
          <div className="mb-3">
            <label className="form-check form-switch">
              <input className="form-check-input" type="checkbox" checked={requireConfirm} onChange={e => setRequireConfirm(e.target.checked)} />
              <span className="form-check-label">是否需要人工确认</span>
            </label>
          </div>
          <div className="mb-3">
            <label className="form-check form-switch">
              <input className="form-check-input" type="checkbox" checked={grayscaleEnabled} onChange={e => setGrayscaleEnabled(e.target.checked)} />
              <span className="form-check-label">启用灰度发布</span>
            </label>
          </div>
          {grayscaleEnabled && (
            <div className="mb-3">
              <label className="form-label">灰度比例: {grayscaleRatio}%</label>
              <input type="range" className="form-range" min="0" max="100" value={grayscaleRatio} onChange={e => setGrayscaleRatio(parseInt(e.target.value, 10))} />
            </div>
          )}
          <div className="mb-3">
            <label className="form-check form-switch">
              <input className="form-check-input" type="checkbox" checked={autoRollback} onChange={e => setAutoRollback(e.target.checked)} />
              <span className="form-check-label">失败自动回滚</span>
            </label>
          </div>
        </form>
      ),
    },
    {
      title: '确认提交',
      content: (
        <div>
          <div className="alert alert-info">请确认部署信息，提交后将立即开始部署流程</div>
          <div className="card">
            <div className="card-body">
              <p><strong>版本:</strong> {versions.find((v) => v.version === selectedVersion)?.version}</p>
              <p>
                <strong>应用:</strong>{' '}
                {selectedApps
                  .map((name) => applications.find((a) => a.name === name)?.name)
                  .join(', ')}
              </p>
              <p>
                <strong>环境:</strong>{' '}
                {selectedEnvs
                  .map((id) => environments.find((e) => e.id === id)?.name)
                  .join(', ')}
              </p>
              <p><strong>需要人工确认:</strong> {requireConfirm ? '是' : '否'}</p>
              <p><strong>灰度发布:</strong> {grayscaleEnabled ? `是 (${grayscaleRatio}%)` : '否'}</p>
              <p><strong>自动回滚:</strong> {autoRollback ? '是' : '否'}</p>
            </div>
          </div>
        </div>
      ),
    },
  ]

  const isStepValid = () => {
    switch (currentStep) {
      case 0:
        return !!selectedVersion
      case 1:
        return selectedApps.length > 0
      case 2:
        return selectedEnvs.length > 0
      default:
        return true
    }
  }

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <a href="javascript:void(0)" className="btn btn-ghost-secondary" onClick={(e) => { e.preventDefault(); navigate('/deployments')}}>
              <IconArrowLeft />
              返回
            </a>
            <h2 className="page-title ms-2 d-inline-block">新建部署</h2>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <ul className="nav nav-tabs card-header-tabs">
            {steps.map((step, index) => (
              <li className="nav-item" key={index}>
                <a href="javascript:void(0)" className={`nav-link ${currentStep === index ? 'active' : ''}`} onClick={(e) => { e.preventDefault(); setCurrentStep(index); }}>
                  {step.title}
                </a>
              </li>
            ))}
          </ul>
        </div>
        <div className="card-body">
          {steps[currentStep].content}
        </div>
        <div className="card-footer d-flex justify-content-between">
          <button className="btn btn-secondary" disabled={currentStep === 0} onClick={handlePrev}>
            上一步
          </button>
          <div>
            {currentStep < steps.length - 1 && (
              <button className="btn btn-primary" onClick={handleNext} disabled={!isStepValid()}>
                下一步
              </button>
            )}
            {currentStep === steps.length - 1 && (
              <button
                className="btn btn-primary"
                onClick={handleSubmit}
                disabled={loading || !isStepValid()}
              >
                {loading ? '提交中...' : '提交部署'}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default CreateDeployment
