import React, { useEffect, useState, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { Card, CardHeader, CardBody, Button, Table, TableHeader, TableColumn, TableBody, TableRow, TableCell, Chip } from '@heroui/react'
import {
  Rocket,
  AppWindow,
  Cloud,
  Tag,
  RefreshCw,
} from 'lucide-react'
import { dashboardApi } from '@/services/api'
import { formatDate, getStatusColor, getStatusText } from '@/utils'
import type { Deployment, DashboardStats } from '@/types'
import { useErrorStore } from '@/store/error'

const Dashboard: React.FC = () => {
  const { setError } = useErrorStore();
  const [stats, setStats] = useState<DashboardStats>({
    activeVersions: 0,
    runningDeployments: 0,
    totalApplications: 0,
    totalEnvironments: 0,
  })
  const [recentDeployments, setRecentDeployments] = useState<Deployment[]>([])
  const [loading, setLoading] = useState(false)

  const loadData = useCallback(async () => {
    setLoading(true)
    try {
      const [statsData, deploymentsData] = await Promise.all([
        dashboardApi.getStats(),
        dashboardApi.getRecentDeployments(10),
      ])
      setStats(statsData)
      setRecentDeployments(deploymentsData)
    } catch (error) {
      setError('Failed to load dashboard data.')
    } finally {
      setLoading(false)
    }
  }, [setError])

  useEffect(() => {
    loadData()
  }, [loadData])

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold">仪表板</h2>
        <Button
          color="primary"
          startContent={<RefreshCw size={16} />}
          onClick={loadData}
          isLoading={loading}
        >
          刷新
        </Button>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardBody className="flex flex-row items-center justify-between">
            <div>
              <p className="text-sm text-default-500">活跃版本</p>
              <p className="text-3xl font-bold mt-2">{stats.activeVersions}</p>
            </div>
            <Tag size={32} className="text-primary" />
          </CardBody>
        </Card>
        <Card>
          <CardBody className="flex flex-row items-center justify-between">
            <div>
              <p className="text-sm text-default-500">进行中的部署</p>
              <p className="text-3xl font-bold mt-2">{stats.runningDeployments}</p>
            </div>
            <Rocket size={32} className="text-warning" />
          </CardBody>
        </Card>
        <Card>
          <CardBody className="flex flex-row items-center justify-between">
            <div>
              <p className="text-sm text-default-500">应用总数</p>
              <p className="text-3xl font-bold mt-2">{stats.totalApplications}</p>
            </div>
            <AppWindow size={32} className="text-success" />
          </CardBody>
        </Card>
        <Card>
          <CardBody className="flex flex-row items-center justify-between">
            <div>
              <p className="text-sm text-default-500">环境总数</p>
              <p className="text-3xl font-bold mt-2">{stats.totalEnvironments}</p>
            </div>
            <Cloud size={32} className="text-secondary" />
          </CardBody>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <h3 className="text-xl font-semibold">最近部署</h3>
        </CardHeader>
        <CardBody>
          <Table aria-label="最近部署列表">
            <TableHeader>
              <TableColumn>版本号</TableColumn>
              <TableColumn>应用</TableColumn>
              <TableColumn>目标环境</TableColumn>
              <TableColumn>状态</TableColumn>
              <TableColumn>创建时间</TableColumn>
              <TableColumn>操作</TableColumn>
            </TableHeader>
            <TableBody>
              {recentDeployments.map((deployment) => (
                <TableRow key={deployment.id}>
                  <TableCell>{deployment.version}</TableCell>
                  <TableCell>{deployment.applications.join(', ')}</TableCell>
                  <TableCell>{deployment.environments.join(', ')}</TableCell>
                  <TableCell>
                    <Chip color={getStatusColor(deployment.status) as any} variant="flat">
                      {getStatusText(deployment.status)}
                    </Chip>
                  </TableCell>
                  <TableCell>{formatDate(deployment.createdAt)}</TableCell>
                  <TableCell>
                    <Button
                      as={Link}
                      to={`/deployments/${deployment.id}`}
                      size="sm"
                      color="primary"
                      variant="light"
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

export default Dashboard
