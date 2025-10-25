import React, { useEffect, useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  versionApi,
  applicationApi,
  environmentApi,
  deploymentApi,
} from '@/services/api'
import type { Version, Application, Environment, CreateDeploymentRequest } from '@/types'
import { ArrowLeft } from 'lucide-react'
import { useErrorStore } from '@/store/error'
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  CardFooter,
  Tabs,
  Tab,
  Checkbox,
  Slider,
} from '@heroui/react'

const CreateDeployment: React.FC = () => {
  const navigate = useNavigate()
  const { setError } = useErrorStore();
  const [currentStep, setCurrentStep] = useState('version')
  
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

  const isStepValid = (step: string) => {
    switch (step) {
      case 'version':
        return !!selectedVersion
      case 'apps':
        return selectedApps.length > 0
      case 'envs':
        return selectedEnvs.length > 0
      default:
        return true
    }
  }

  const canGoNext = () => {
    const steps = ['version', 'apps', 'envs', 'config', 'confirm']
    const currentIndex = steps.indexOf(currentStep)
    return currentIndex < steps.length - 1 && isStepValid(currentStep)
  }

  const canGoPrev = () => {
    const steps = ['version', 'apps', 'envs', 'config', 'confirm']
    return steps.indexOf(currentStep) > 0
  }

  const handleNext = () => {
    const steps = ['version', 'apps', 'envs', 'config', 'confirm']
    const currentIndex = steps.indexOf(currentStep)
    if (currentIndex < steps.length - 1) {
      setCurrentStep(steps[currentIndex + 1])
    }
  }

  const handlePrev = () => {
    const steps = ['version', 'apps', 'envs', 'config', 'confirm']
    const currentIndex = steps.indexOf(currentStep)
    if (currentIndex > 0) {
      setCurrentStep(steps[currentIndex - 1])
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Button
          isIconOnly
          variant="light"
          onPress={() => navigate('/deployments')}
        >
          <ArrowLeft className="w-5 h-5" />
        </Button>
        <h2 className="text-3xl font-bold">新建部署</h2>
      </div>

      <Card>
        <CardHeader>
          <Tabs
            selectedKey={currentStep}
            onSelectionChange={(key) => setCurrentStep(key as string)}
          >
            <Tab key="version" title="选择版本" />
            <Tab key="apps" title="选择应用" />
            <Tab key="envs" title="选择环境" />
            <Tab key="config" title="配置选项" />
            <Tab key="confirm" title="确认提交" />
          </Tabs>
        </CardHeader>
        <CardBody>
          {currentStep === 'version' && (
            <div className="space-y-2">
              {versions.map((version) => (
                <Card
                  key={version.version}
                  isPressable
                  isHoverable
                  className={selectedVersion === version.version ? 'border-2 border-primary' : ''}
                  onPress={() => setSelectedVersion(version.version)}
                >
                  <CardBody>
                    <div>
                      <strong>{version.version}</strong>
                      <small className="block text-gray-500">
                        {version.git.tag} - {version.createdAt}
                      </small>
                    </div>
                  </CardBody>
                </Card>
              ))}
            </div>
          )}

          {currentStep === 'apps' && (
            <div className="space-y-2">
              {applications.map((app) => (
                <Card
                  key={app.name}
                  isPressable
                  isHoverable
                  className={selectedApps.includes(app.name) ? 'border-2 border-primary' : ''}
                  onPress={() => toggleSelection(app.name, selectedApps, setSelectedApps)}
                >
                  <CardBody>{app.name}</CardBody>
                </Card>
              ))}
            </div>
          )}

          {currentStep === 'envs' && (
            <div className="space-y-2">
              {environments.map((env) => (
                <Card
                  key={env.id}
                  isPressable
                  isHoverable
                  className={selectedEnvs.includes(env.id) ? 'border-2 border-primary' : ''}
                  onPress={() => toggleSelection(env.id, selectedEnvs, setSelectedEnvs)}
                >
                  <CardBody>
                    {env.name} ({env.type})
                  </CardBody>
                </Card>
              ))}
            </div>
          )}

          {currentStep === 'config' && (
            <div className="space-y-4">
              <Checkbox
                isSelected={requireConfirm}
                onValueChange={setRequireConfirm}
              >
                是否需要人工确认
              </Checkbox>
              <Checkbox
                isSelected={grayscaleEnabled}
                onValueChange={setGrayscaleEnabled}
              >
                启用灰度发布
              </Checkbox>
              {grayscaleEnabled && (
                <Slider
                  label="灰度比例"
                  value={grayscaleRatio}
                  onChange={(value) => setGrayscaleRatio(value as number)}
                  maxValue={100}
                  minValue={0}
                  step={1}
                  formatOptions={{ style: 'percent', maximumFractionDigits: 0 }}
                  className="max-w-md"
                />
              )}
              <Checkbox
                isSelected={autoRollback}
                onValueChange={setAutoRollback}
              >
                失败自动回滚
              </Checkbox>
            </div>
          )}

          {currentStep === 'confirm' && (
            <div className="space-y-4">
              <div className="bg-primary-50 p-4 rounded">
                请确认部署信息，提交后将立即开始部署流程
              </div>
              <Card>
                <CardBody className="space-y-2">
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
                </CardBody>
              </Card>
            </div>
          )}
        </CardBody>
        <CardFooter>
          <div className="flex justify-between w-full">
            <Button
              variant="light"
              isDisabled={!canGoPrev()}
              onPress={handlePrev}
            >
              上一步
            </Button>
            <div>
              {currentStep !== 'confirm' ? (
                <Button
                  color="primary"
                  onPress={handleNext}
                  isDisabled={!canGoNext()}
                >
                  下一步
                </Button>
              ) : (
                <Button
                  color="primary"
                  onPress={handleSubmit}
                  isDisabled={loading || !isStepValid('version') || !isStepValid('apps') || !isStepValid('envs')}
                >
                  {loading ? '提交中...' : '提交部署'}
                </Button>
              )}
            </div>
          </div>
        </CardFooter>
      </Card>
    </div>
  )
}

export default CreateDeployment
