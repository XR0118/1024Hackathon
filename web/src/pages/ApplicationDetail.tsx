import React, { useEffect, useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { applicationApi } from '@/services/api'
import { formatDate } from '@/utils'
import type { Application } from '@/types'
import { ArrowLeft, CheckCircle, AlertCircle, Server } from 'lucide-react'
import { useErrorStore } from '@/store/error'
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  Chip,
  Spinner,
  Accordion,
  AccordionItem,
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
} from '@heroui/react'

const ApplicationDetail: React.FC = () => {
  const { name } = useParams<{ name: string }>()
  const navigate = useNavigate()
  const { setError } = useErrorStore()
  const [application, setApplication] = useState<Application | null>(null)
  const [loading, setLoading] = useState(true)

  const loadApplication = useCallback(async () => {
    if (!name) return
    
    try {
      setLoading(true)
      const data = await applicationApi.get(name)
      setApplication(data)
    } catch (error) {
      setError('Failed to load application details.')
    } finally {
      setLoading(false)
    }
  }, [name, setError])

  useEffect(() => {
    loadApplication()
  }, [loadApplication])

  const getHealthColor = (health: number) => {
    if (health >= 80) return 'success'
    if (health >= 50) return 'warning'
    return 'danger'
  }

  const getHealthIcon = (health: number) => {
    if (health >= 80) return <CheckCircle className="w-4 h-4" />
    return <AlertCircle className="w-4 h-4" />
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <Spinner size="lg" />
      </div>
    )
  }

  if (!application) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen gap-4">
        <p className="text-xl font-semibold">应用不存在</p>
        <Button color="primary" onPress={() => navigate('/applications')}>
          返回应用列表
        </Button>
      </div>
    )
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <div className="flex items-center gap-3">
          <Button
            isIconOnly
            variant="light"
            onPress={() => navigate('/applications')}
          >
            <ArrowLeft className="w-5 h-5" />
          </Button>
          {application.icon && (
            <img src={application.icon} alt={application.name} className="w-10 h-10" />
          )}
          <div>
            <h2 className="text-3xl font-bold">{application.name}</h2>
            <p className="text-gray-500">{application.description}</p>
          </div>
        </div>
      </div>

      <Card>
        <CardHeader>
          <h3 className="text-xl font-semibold">版本列表</h3>
        </CardHeader>
        <CardBody>
          {application.versions && application.versions.length > 0 ? (
            <Accordion defaultExpandedKeys={["0"]}>
              {application.versions.map((versionInfo, index) => (
                <AccordionItem
                  key={index}
                  aria-label={`Version ${versionInfo.version}`}
                  title={
                    <div className="flex justify-between items-center w-full">
                      <div className="flex items-center gap-3">
                        <Chip
                          color={versionInfo.status === 'normal' ? 'primary' : 'warning'}
                          variant="flat"
                        >
                          {versionInfo.version}
                        </Chip>
                        {versionInfo.status === 'revert' && (
                          <Chip color="warning" variant="flat">回滚版本</Chip>
                        )}
                      </div>
                      <div className="flex items-center gap-3">
                        <Chip
                          color={getHealthColor(versionInfo.health)}
                          variant="flat"
                          startContent={getHealthIcon(versionInfo.health)}
                        >
                          健康度: {versionInfo.health}%
                        </Chip>
                        <small className="text-gray-500">
                          最后更新: {formatDate(versionInfo.lastUpdatedAt)}
                        </small>
                      </div>
                    </div>
                  }
                >
                  <div className="space-y-4">
                    <div className="flex items-center gap-2">
                      <Server className="w-5 h-5" />
                      <h4 className="text-lg font-semibold">节点列表</h4>
                    </div>
                    {versionInfo.nodes && versionInfo.nodes.length > 0 ? (
                      <Table aria-label="节点列表">
                        <TableHeader>
                          <TableColumn>节点名称</TableColumn>
                          <TableColumn>健康度</TableColumn>
                          <TableColumn>最后更新时间</TableColumn>
                        </TableHeader>
                        <TableBody>
                          {versionInfo.nodes.map((node, nodeIndex) => (
                            <TableRow key={nodeIndex}>
                              <TableCell>
                                <div className="flex items-center gap-2">
                                  <Server className="w-4 h-4" />
                                  {node.name}
                                </div>
                              </TableCell>
                              <TableCell>
                                <Chip
                                  color={getHealthColor(node.health)}
                                  variant="flat"
                                  startContent={getHealthIcon(node.health)}
                                >
                                  {node.health}%
                                </Chip>
                              </TableCell>
                              <TableCell className="text-gray-500">
                                {formatDate(node.lastUpdatedAt)}
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    ) : (
                      <div className="text-center py-8">
                        <p className="text-gray-500">暂无节点信息</p>
                      </div>
                    )}
                  </div>
                </AccordionItem>
              ))}
            </Accordion>
          ) : (
            <div className="text-center py-8">
              <p className="text-gray-500">暂无版本信息</p>
            </div>
          )}
        </CardBody>
      </Card>
    </div>
  )
}

export default ApplicationDetail
