import React, { useEffect, useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, CardBody, CardFooter, Button, Chip } from '@heroui/react'
import { applicationApi } from '@/services/api'
import { formatDate } from '@/utils'
import type { Application } from '@/types'
import { Plus, Rocket, CheckCircle, AlertCircle } from 'lucide-react'
import { useErrorStore } from '@/store/error'

const Applications: React.FC = () => {
  const navigate = useNavigate()
  const { setError } = useErrorStore()
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

  const getHealthColor = (health: number) => {
    if (health >= 80) return 'success'
    if (health >= 50) return 'warning'
    return 'danger'
  }

  const getHealthIcon = (health: number) => {
    if (health >= 80) return <CheckCircle size={16} />
    return <AlertCircle size={16} />
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold">应用管理</h2>
        <Button color="primary" startContent={<Plus size={16} />}>
          添加应用
        </Button>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {applications.map((app) => (
          <Card key={app.name} className="w-full">
            <CardBody className="space-y-4">
              <div className="flex items-center gap-3">
                {app.icon && (
                  <div>
                    <img src={app.icon} alt={app.name} className="w-10 h-10 rounded" />
                  </div>
                )}
                <div className="flex-1">
                  <h3 className="text-lg font-semibold">{app.name}</h3>
                  <p className="text-sm text-default-500">{app.description}</p>
                </div>
              </div>

              <div className="space-y-2">
                <p className="font-semibold text-sm">版本信息:</p>
                {app.versions && app.versions.length > 0 ? (
                  <div className="space-y-2">
                    {app.versions.slice(0, 3).map((versionInfo, index) => (
                      <div key={index} className="flex justify-between items-center py-2 border-b border-default-200 last:border-0">
                        <div className="flex items-center gap-2">
                          <Chip size="sm" color={versionInfo.status === 'normal' ? 'primary' : 'warning'} variant="flat">
                            {versionInfo.version}
                          </Chip>
                          {versionInfo.status === 'revert' && (
                            <Chip size="sm" color="warning">回滚</Chip>
                          )}
                        </div>
                        <div className="flex items-center gap-2">
                          <Chip
                            size="sm"
                            color={getHealthColor(versionInfo.health) as any}
                            variant="flat"
                            startContent={getHealthIcon(versionInfo.health)}
                          >
                            {versionInfo.health}%
                          </Chip>
                          <span className="text-xs text-default-500">{formatDate(versionInfo.lastUpdatedAt)}</span>
                        </div>
                      </div>
                    ))}
                    {app.versions.length > 3 && (
                      <p className="text-center text-sm text-default-500">
                        还有 {app.versions.length - 3} 个版本
                      </p>
                    )}
                  </div>
                ) : (
                  <p className="text-sm text-default-500">暂无版本信息</p>
                )}
              </div>
            </CardBody>
            <CardFooter className="gap-2">
              <Button
                color="primary"
                variant="light"
                startContent={<Rocket size={16} />}
                onClick={() => navigate(`/deployments/new?appName=${app.name}`)}
              >
                新建部署
              </Button>
              <Button
                variant="flat"
                onClick={() => navigate(`/applications/${app.name}`)}
              >
                查看详情
              </Button>
            </CardFooter>
          </Card>
        ))}
      </div>

      {applications.length === 0 && (
        <div className="flex flex-col items-center justify-center py-16">
          <Rocket size={64} className="text-default-300 mb-4" />
          <p className="text-xl font-semibold mb-2">暂无应用</p>
          <p className="text-default-500">
            点击上方"添加应用"按钮创建您的第一个应用
          </p>
        </div>
      )}
    </div>
  )
}

export default Applications
