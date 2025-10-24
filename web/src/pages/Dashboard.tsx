import React, { useEffect, useState } from 'react'
import { Card, Row, Col, Table, Tag, Statistic, Button } from 'antd'
import {
  RocketOutlined,
  AppstoreOutlined,
  CloudOutlined,
  TagOutlined,
  ReloadOutlined,
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { dashboardApi } from '@/services/api'
import { formatDate, getStatusColor, getStatusText } from '@/utils'
import type { Deployment, DashboardStats } from '@/types'

const Dashboard: React.FC = () => {
  const navigate = useNavigate()
  const [stats, setStats] = useState<DashboardStats>({
    activeVersions: 0,
    runningDeployments: 0,
    totalApplications: 0,
    totalEnvironments: 0,
  })
  const [recentDeployments, setRecentDeployments] = useState<Deployment[]>([])
  const [loading, setLoading] = useState(false)

  const loadData = async () => {
    setLoading(true)
    try {
      const [statsData, deploymentsData] = await Promise.all([
        dashboardApi.getStats(),
        dashboardApi.getRecentDeployments(10),
      ])
      setStats(statsData)
      setRecentDeployments(deploymentsData)
    } catch (error) {
      console.error('Failed to load dashboard data:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [])

  const columns = [
    {
      title: '版本号',
      dataIndex: 'version',
      key: 'version',
    },
    {
      title: '应用',
      dataIndex: 'applications',
      key: 'applications',
      render: (apps: string[]) => apps.join(', '),
    },
    {
      title: '目标环境',
      dataIndex: 'environments',
      key: 'environments',
      render: (envs: string[]) => envs.join(', '),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: Deployment['status']) => (
        <Tag color={getStatusColor(status)}>{getStatusText(status)}</Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (date: string) => formatDate(date),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Deployment) => (
        <Button
          type="link"
          onClick={() => navigate(`/deployments/${record.id}`)}
        >
          查看详情
        </Button>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0 }}>仪表板</h1>
        <Button icon={<ReloadOutlined />} onClick={loadData} loading={loading}>
          刷新
        </Button>
      </div>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="活跃版本"
              value={stats.activeVersions}
              prefix={<TagOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="进行中的部署"
              value={stats.runningDeployments}
              prefix={<RocketOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="应用总数"
              value={stats.totalApplications}
              prefix={<AppstoreOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="环境总数"
              value={stats.totalEnvironments}
              prefix={<CloudOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <Card title="最近部署" bordered={false}>
        <Table
          columns={columns}
          dataSource={recentDeployments}
          rowKey="id"
          loading={loading}
          pagination={false}
        />
      </Card>
    </div>
  )
}

export default Dashboard
