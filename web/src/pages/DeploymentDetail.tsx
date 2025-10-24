import React, { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Card,
  Descriptions,
  Tag,
  Steps,
  Button,
  Space,
  Modal,
  Input,
  Alert,
  Slider,
} from 'antd'
import {
  ArrowLeftOutlined,
  CheckOutlined,
  CloseOutlined,
  RollbackOutlined,
} from '@ant-design/icons'
import { deploymentApi } from '@/services/api'
import { formatDate, formatDuration, getStatusColor, getStatusText } from '@/utils'
import type { DeploymentDetail } from '@/types'

const DeploymentDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [deployment, setDeployment] = useState<DeploymentDetail | null>(null)
  const [loading, setLoading] = useState(false)
  const [confirmModalVisible, setConfirmModalVisible] = useState(false)
  const [rollbackModalVisible, setRollbackModalVisible] = useState(false)
  const [note, setNote] = useState('')
  const [reason, setReason] = useState('')

  const loadDeployment = async () => {
    if (!id) return
    setLoading(true)
    try {
      const data = await deploymentApi.get(id)
      setDeployment(data)
    } catch (error) {
      console.error('Failed to load deployment:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadDeployment()
    const interval = setInterval(loadDeployment, 3000)
    return () => clearInterval(interval)
  }, [id])

  const handleConfirm = async () => {
    if (!id) return
    try {
      await deploymentApi.confirm(id, note)
      setConfirmModalVisible(false)
      setNote('')
      loadDeployment()
    } catch (error) {
      console.error('Failed to confirm deployment:', error)
    }
  }

  const handleRollback = async () => {
    if (!id) return
    try {
      await deploymentApi.rollback(id, reason)
      setRollbackModalVisible(false)
      setReason('')
      loadDeployment()
    } catch (error) {
      console.error('Failed to rollback deployment:', error)
    }
  }

  if (!deployment) {
    return <Card loading={loading}>加载中...</Card>
  }

  const currentStep = deployment.steps.findIndex(
    (step) => step.status === 'running' || step.status === 'pending'
  )

  return (
    <div>
      <Space style={{ marginBottom: 24 }}>
        <Button 
          icon={<ArrowLeftOutlined />} 
          onClick={() => navigate('/deployments')}
          style={{ borderRadius: 6 }}
        >
          返回
        </Button>
        <h1 style={{ margin: 0 }}>部署详情</h1>
      </Space>

      <Card 
        style={{ 
          marginBottom: 16,
          borderRadius: 8,
          border: '1px solid #e5e7eb',
          boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        }}
      >
        <Descriptions column={2}>
          <Descriptions.Item label="部署ID">{deployment.id}</Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={getStatusColor(deployment.status)}>
              {getStatusText(deployment.status)}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="版本">{deployment.version}</Descriptions.Item>
          <Descriptions.Item label="应用">
            {deployment.applications.join(', ')}
          </Descriptions.Item>
          <Descriptions.Item label="环境">
            {deployment.environments.join(', ')}
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">
            {formatDate(deployment.createdAt)}
          </Descriptions.Item>
          {deployment.duration && (
            <Descriptions.Item label="执行时长">
              {formatDuration(deployment.duration)}
            </Descriptions.Item>
          )}
        </Descriptions>
      </Card>

      {deployment.status === 'waiting_confirm' && (
        <Alert
          message="需要人工确认"
          description="此部署需要人工确认后才能继续。请仔细检查部署状态后决定是否继续。"
          type="warning"
          showIcon
          style={{ marginBottom: 16 }}
          action={
            <Space>
              <Button
                size="small"
                type="primary"
                icon={<CheckOutlined />}
                onClick={() => setConfirmModalVisible(true)}
              >
                确认继续
              </Button>
              <Button
                size="small"
                danger
                icon={<RollbackOutlined />}
                onClick={() => setRollbackModalVisible(true)}
              >
                回滚
              </Button>
            </Space>
          }
        />
      )}

      {deployment.grayscaleEnabled && (
        <Card 
          title="灰度发布" 
          style={{ 
            marginBottom: 16,
            borderRadius: 8,
            border: '1px solid #e5e7eb',
            boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
          }}
        >
          <p>当前灰度比例: {deployment.grayscaleRatio}%</p>
          <Slider value={deployment.grayscaleRatio} disabled />
        </Card>
      )}

      <Card 
        title="部署流程" 
        style={{ 
          marginBottom: 16,
          borderRadius: 8,
          border: '1px solid #e5e7eb',
          boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        }}
      >
        <Steps
          current={currentStep >= 0 ? currentStep : deployment.steps.length}
          status={
            deployment.status === 'failed'
              ? 'error'
              : deployment.status === 'success'
              ? 'finish'
              : 'process'
          }
          items={deployment.steps.map((step) => ({
            title: step.name,
            description: step.duration ? formatDuration(step.duration) : undefined,
            status:
              step.status === 'success'
                ? 'finish'
                : step.status === 'failed'
                ? 'error'
                : step.status === 'running'
                ? 'process'
                : 'wait',
          }))}
        />
      </Card>

      <Card 
        title="实时日志"
        style={{
          borderRadius: 8,
          border: '1px solid #e5e7eb',
          boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        }}
      >
        <div
          style={{
            background: '#000',
            color: '#0f0',
            padding: 16,
            borderRadius: 4,
            fontFamily: 'monospace',
            fontSize: 12,
            maxHeight: 400,
            overflow: 'auto',
          }}
        >
          {deployment.logs.map((log, index) => (
            <div key={index}>
              <span style={{ color: '#666' }}>[{log.timestamp}]</span>{' '}
              <span
                style={{
                  color:
                    log.level === 'error'
                      ? '#f00'
                      : log.level === 'warn'
                      ? '#fa0'
                      : '#0f0',
                }}
              >
                [{log.level.toUpperCase()}]
              </span>{' '}
              {log.message}
            </div>
          ))}
        </div>
      </Card>

      <Modal
        title="确认部署"
        open={confirmModalVisible}
        onOk={handleConfirm}
        onCancel={() => setConfirmModalVisible(false)}
      >
        <p>确认继续此部署吗?</p>
        <Input.TextArea
          placeholder="备注(可选)"
          value={note}
          onChange={(e) => setNote(e.target.value)}
          rows={4}
        />
      </Modal>

      <Modal
        title="回滚部署"
        open={rollbackModalVisible}
        onOk={handleRollback}
        onCancel={() => setRollbackModalVisible(false)}
      >
        <p>确认回滚此部署吗?</p>
        <Input.TextArea
          placeholder="回滚原因(可选)"
          value={reason}
          onChange={(e) => setReason(e.target.value)}
          rows={4}
        />
      </Modal>
    </div>
  )
}

export default DeploymentDetailPage
