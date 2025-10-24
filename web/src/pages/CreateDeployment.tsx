import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Steps,
  Button,
  Card,
  Form,
  Select,
  Switch,
  Slider,
  Space,
  Alert,
  List,
  Tag,
} from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import {
  versionApi,
  applicationApi,
  environmentApi,
  deploymentApi,
} from '@/services/api'
import type { Version, Application, Environment, CreateDeploymentRequest } from '@/types'

const CreateDeployment: React.FC = () => {
  const navigate = useNavigate()
  const [currentStep, setCurrentStep] = useState(0)
  const [form] = Form.useForm()
  
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

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
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
      console.error('Failed to load data:', error)
    }
  }

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
      console.error('Failed to create deployment:', error)
    } finally {
      setLoading(false)
    }
  }

  const steps = [
    {
      title: '选择版本',
      content: (
        <List
          dataSource={versions}
          renderItem={(version) => (
            <List.Item
              onClick={() => setSelectedVersion(version.id)}
              style={{
                cursor: 'pointer',
                background: selectedVersion === version.id ? '#e6f7ff' : undefined,
                padding: 16,
              }}
            >
              <List.Item.Meta
                title={version.version}
                description={
                  <Space>
                    <Tag>{version.gitTag}</Tag>
                    <span>{version.createdAt}</span>
                  </Space>
                }
              />
            </List.Item>
          )}
        />
      ),
    },
    {
      title: '选择应用',
      content: (
        <Select
          mode="multiple"
          style={{ width: '100%' }}
          placeholder="选择至少一个应用"
          value={selectedApps}
          onChange={setSelectedApps}
          options={applications.map((app) => ({
            label: app.name,
            value: app.id,
          }))}
        />
      ),
    },
    {
      title: '选择环境',
      content: (
        <Select
          mode="multiple"
          style={{ width: '100%' }}
          placeholder="选择至少一个环境"
          value={selectedEnvs}
          onChange={setSelectedEnvs}
          options={environments.map((env) => ({
            label: `${env.name} (${env.type})`,
            value: env.id,
          }))}
        />
      ),
    },
    {
      title: '配置选项',
      content: (
        <Form layout="vertical">
          <Form.Item label="是否需要人工确认">
            <Switch checked={requireConfirm} onChange={setRequireConfirm} />
          </Form.Item>
          <Form.Item label="启用灰度发布">
            <Switch checked={grayscaleEnabled} onChange={setGrayscaleEnabled} />
          </Form.Item>
          {grayscaleEnabled && (
            <Form.Item label={`灰度比例: ${grayscaleRatio}%`}>
              <Slider
                min={0}
                max={100}
                value={grayscaleRatio}
                onChange={setGrayscaleRatio}
              />
            </Form.Item>
          )}
          <Form.Item label="失败自动回滚">
            <Switch checked={autoRollback} onChange={setAutoRollback} />
          </Form.Item>
        </Form>
      ),
    },
    {
      title: '确认提交',
      content: (
        <div>
          <Alert
            message="请确认部署信息"
            description="提交后将立即开始部署流程"
            type="info"
            showIcon
            style={{ marginBottom: 16 }}
          />
          <Card 
            title="部署摘要" 
            size="small"
            style={{
              borderRadius: 8,
              border: '1px solid #e5e7eb',
            }}
          >
            <p><strong>版本:</strong> {versions.find((v) => v.id === selectedVersion)?.version}</p>
            <p>
              <strong>应用:</strong>{' '}
              {selectedApps
                .map((id) => applications.find((a) => a.id === id)?.name)
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
          </Card>
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
      <Space style={{ marginBottom: 24 }}>
        <Button 
          icon={<ArrowLeftOutlined />} 
          onClick={() => navigate('/deployments')}
          style={{ borderRadius: 6 }}
        >
          返回
        </Button>
        <h1 style={{ margin: 0 }}>新建部署</h1>
      </Space>

      <Card
        style={{
          borderRadius: 8,
          border: '1px solid #e5e7eb',
          boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        }}
      >
        <Steps current={currentStep} items={steps.map((s) => ({ title: s.title }))} />
        <div style={{ marginTop: 24, marginBottom: 24 }}>
          {steps[currentStep].content}
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <Button 
            disabled={currentStep === 0} 
            onClick={handlePrev}
            style={{ borderRadius: 6 }}
          >
            上一步
          </Button>
          <Space>
            {currentStep < steps.length - 1 && (
              <Button 
                type="primary" 
                onClick={handleNext} 
                disabled={!isStepValid()}
                style={{ borderRadius: 6 }}
              >
                下一步
              </Button>
            )}
            {currentStep === steps.length - 1 && (
              <Button
                type="primary"
                onClick={handleSubmit}
                loading={loading}
                disabled={!isStepValid()}
                style={{ borderRadius: 6 }}
              >
                提交部署
              </Button>
            )}
          </Space>
        </div>
      </Card>
    </div>
  )
}

export default CreateDeployment
