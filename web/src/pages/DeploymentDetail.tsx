import React, { useEffect, useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { deploymentApi } from '@/services/api'
import { formatDate, formatDuration, getStatusColor, getStatusText } from '@/utils'
import type { DeploymentDetail as DeploymentDetailType } from '@/types'
import { IconArrowLeft, IconCheck, IconArrowBackUp } from '@tabler/icons-react'
import { useErrorStore } from '@/store/error'

const DeploymentDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { setError } = useErrorStore();
  const [deployment, setDeployment] = useState<DeploymentDetailType | null>(null)
  const [loading, setLoading] = useState(false)
  const [note, setNote] = useState('')
  const [reason, setReason] = useState('')

  const loadDeployment = useCallback(async () => {
    if (!id) return
    setLoading(true)
    try {
      const data = await deploymentApi.get(id)
      setDeployment(data)
    } catch (error) {
      setError('Failed to load deployment details.')
    } finally {
      setLoading(false)
    }
  }, [id, setError])

  useEffect(() => {
    loadDeployment()
    const interval = setInterval(loadDeployment, 3000)
    return () => clearInterval(interval)
  }, [id, loadDeployment])

  const handleConfirm = async () => {
    if (!id) return
    try {
      await deploymentApi.confirm(id, note)
      setNote('')
      loadDeployment()
    } catch (error) {
      setError('Failed to confirm deployment.')
    }
  }

  const handleRollback = async () => {
    if (!id) return
    try {
      await deploymentApi.rollback(id, reason)
      setReason('')
      loadDeployment()
    } catch (error) {
      setError('Failed to rollback deployment.')
    }
  }

  if (loading && !deployment) {
    return <div>Loading...</div>
  }

  if (!deployment) {
    return <div>Deployment not found</div>
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
            <h2 className="page-title ms-2 d-inline-block">部署详情</h2>
          </div>
        </div>
      </div>

      <div className="card mb-3">
        <div className="card-body">
          <div className="row">
            <div className="col-md-6">
              <p><strong>部署ID:</strong> {deployment.id}</p>
              <p><strong>版本:</strong> {deployment.version}</p>
              <p><strong>应用:</strong> {deployment.applications.join(', ')}</p>
            </div>
            <div className="col-md-6">
              <p><strong>状态:</strong> <span className={`badge bg-${getStatusColor(deployment.status)}-lt`}>{getStatusText(deployment.status)}</span></p>
              <p><strong>环境:</strong> {deployment.environments.join(', ')}</p>
              <p><strong>创建时间:</strong> {formatDate(deployment.createdAt)}</p>
            </div>
            {deployment.duration && <div className="col-12"><p><strong>执行时长:</strong> {formatDuration(deployment.duration)}</p></div>}
          </div>
        </div>
      </div>

      {deployment.status === 'waiting_confirm' && (
        <div className="alert alert-warning d-flex justify-content-between align-items-center">
          <div>
            <strong>需要人工确认:</strong> 此部署需要人工确认后才能继续。
          </div>
          <div>
            <button className="btn btn-success me-2" data-bs-toggle="modal" data-bs-target="#confirmModal">
              <IconCheck size={16} className="me-2" />
              确认继续
            </button>
            <button className="btn btn-danger" data-bs-toggle="modal" data-bs-target="#rollbackModal">
              <IconArrowBackUp size={16} className="me-2" />
              回滚
            </button>
          </div>
        </div>
      )}

      {deployment.grayscaleEnabled && (
        <div className="card mb-3">
          <div className="card-header"><h3 className="card-title">灰度发布</h3></div>
          <div className="card-body">
            <label className="form-label">当前灰度比例: {deployment.grayscaleRatio}%</label>
            <input type="range" className="form-range" value={deployment.grayscaleRatio} disabled />
          </div>
        </div>
      )}

      <div className="card mb-3">
        <div className="card-header"><h3 className="card-title">部署流程</h3></div>
        <div className="card-body">
          <ul className="steps">
            {deployment.steps.map((step, index) => (
              <li key={index} className={`step-item ${
                step.status === 'success' ? 'active' : ''
              }`}>
                <a href="javascript:void(0)" onClick={(e) => e.preventDefault()}>{step.name}</a>
              </li>
            ))}
          </ul>
        </div>
      </div>

      <div className="card">
        <div className="card-header"><h3 className="card-title">实时日志</h3></div>
        <div className="card-body" style={{ background: '#000', color: '#0f0', fontFamily: 'monospace', fontSize: '12px', maxHeight: '400px', overflow: 'auto' }}>
          {deployment.logs.map((log, index) => (
            <div key={index}>
              <span style={{ color: '#666' }}>[{log.timestamp}]</span>{' '}
              <span style={{ color: log.level === 'error' ? '#f00' : log.level === 'warn' ? '#fa0' : '#0f0' }}>[{log.level.toUpperCase()}]</span>{' '}
              {log.message}
            </div>
          ))}
        </div>
      </div>

      {/* Modals */}
      <div className="modal" id="confirmModal" tabIndex={-1}>
        <div className="modal-dialog">
          <div className="modal-content">
            <div className="modal-header">
              <h5 className="modal-title">确认部署</h5>
              <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div className="modal-body">
              <p>确认继续此部署吗?</p>
              <textarea className="form-control" placeholder="备注(可选)" value={note} onChange={(e) => setNote(e.target.value)} rows={4}></textarea>
            </div>
            <div className="modal-footer">
              <button type="button" className="btn btn-secondary" data-bs-dismiss="modal">取消</button>
              <button type="button" className="btn btn-primary" onClick={handleConfirm} data-bs-dismiss="modal">确认</button>
            </div>
          </div>
        </div>
      </div>

      <div className="modal" id="rollbackModal" tabIndex={-1}>
        <div className="modal-dialog">
          <div className="modal-content">
            <div className="modal-header">
              <h5 className="modal-title">回滚部署</h5>
              <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div className="modal-body">
              <p>确认回滚此部署吗?</p>
              <textarea className="form-control" placeholder="回滚原因(可选)" value={reason} onChange={(e) => setReason(e.target.value)} rows={4}></textarea>
            </div>
            <div className="modal-footer">
              <button type="button" className="btn btn-secondary" data-bs-dismiss="modal">取消</button>
              <button type="button" className="btn btn-danger" onClick={handleRollback} data-bs-dismiss="modal">回滚</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default DeploymentDetailPage
