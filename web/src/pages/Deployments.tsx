import React, { useEffect, useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, CardHeader, CardBody, Button, Table, TableHeader, TableColumn, TableBody, TableRow, TableCell, Chip, Select, SelectItem, Input, Progress } from '@heroui/react'
import { deploymentApi } from '@/services/api'
import { formatDate, getStatusColor, getStatusText } from '@/utils'
import type { Deployment } from '@/types'
import { Plus, RefreshCw } from 'lucide-react'
import { useErrorStore } from '@/store/error'

const Deployments: React.FC = () => {
  const navigate = useNavigate()
  const [deployments, setDeployments] = useState<Deployment[]>([])
  const [loading, setLoading] = useState(false)
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [startDate, setStartDate] = useState<string>('')
  const [endDate, setEndDate] = useState<string>('')

  const loadDeployments = useCallback(async () => {
    setLoading(true)
    try {
      const data = await deploymentApi.list({
        status: statusFilter || undefined,
        startDate: startDate || undefined,
        endDate: endDate || undefined,
      })
      setDeployments(data)
    } catch (error) {
      useErrorStore.getState().setError('Failed to load deployments.')
    } finally {
      setLoading(false)
    }
  }, [statusFilter, startDate, endDate])

  useEffect(() => {
    loadDeployments()
    const interval = setInterval(loadDeployments, 5000)
    return () => clearInterval(interval)
  }, [loadDeployments])

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold">部署管理</h2>
        <Button
          color="primary"
          startContent={<Plus size={16} />}
          onClick={() => navigate('/deployments/new')}
        >
          新建部署
        </Button>
      </div>

      <Card>
        <CardHeader>
          <div className="flex gap-2 w-full">
            <Select
              label="状态"
              placeholder="所有状态"
              selectedKeys={statusFilter ? [statusFilter] : []}
              onSelectionChange={(keys) => setStatusFilter(Array.from(keys)[0] as string || '')}
              className="max-w-xs"
            >
              <SelectItem key="">所有状态</SelectItem>
              <SelectItem key="pending">待开始</SelectItem>
              <SelectItem key="running">进行中</SelectItem>
              <SelectItem key="success">成功</SelectItem>
              <SelectItem key="failed">失败</SelectItem>
              <SelectItem key="waiting_confirm">待确认</SelectItem>
            </Select>
            <Input
              type="date"
              label="开始日期"
              value={startDate}
              onChange={e => setStartDate(e.target.value)}
              className="max-w-xs"
            />
            <Input
              type="date"
              label="结束日期"
              value={endDate}
              onChange={e => setEndDate(e.target.value)}
              className="max-w-xs"
            />
            <Button
              color="primary"
              startContent={<RefreshCw size={16} />}
              onClick={loadDeployments}
              isLoading={loading}
            >
              刷新
            </Button>
          </div>
        </CardHeader>
        <CardBody>
          <Table aria-label="部署列表">
            <TableHeader>
              <TableColumn>部署ID</TableColumn>
              <TableColumn>版本</TableColumn>
              <TableColumn>应用</TableColumn>
              <TableColumn>环境</TableColumn>
              <TableColumn>状态</TableColumn>
              <TableColumn>进度</TableColumn>
              <TableColumn>创建时间</TableColumn>
              <TableColumn>操作</TableColumn>
            </TableHeader>
            <TableBody>
              {deployments.map((deployment) => (
                <TableRow key={deployment.id}>
                  <TableCell>{deployment.id}</TableCell>
                  <TableCell>{deployment.version}</TableCell>
                  <TableCell>{deployment.applications.join(', ')}</TableCell>
                  <TableCell>{deployment.environments.join(', ')}</TableCell>
                  <TableCell>
                    <Chip color={getStatusColor(deployment.status) as any} variant="flat">
                      {getStatusText(deployment.status)}
                    </Chip>
                  </TableCell>
                  <TableCell>
                    {deployment.status === 'running' && (
                      <Progress
                        value={deployment.progress}
                        className="max-w-md"
                        color="primary"
                      />
                    )}
                  </TableCell>
                  <TableCell>{formatDate(deployment.createdAt)}</TableCell>
                  <TableCell>
                    <Button
                      size="sm"
                      color="primary"
                      variant="light"
                      onClick={() => navigate(`/deployments/${deployment.id}`)}
                    >
                      查看详情
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardBody>
      </Card>
    </div>
  )
}

export default Deployments
