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
      <div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1>仪表板</h1>
        <Button 
          icon={<ReloadOutlined />} 
          onClick={loadData} 
          loading={loading}
          style={{ borderRadius: 6 }}
        >
          刷新数据
        </Button>
      </div>

      <Row gutter={[20, 20]} style={{ marginBottom: 32 }}>
        <Col xs={24} sm={12} lg={6}>
          <div className="stat-card">
            <div className="stat-icon" style={{ color: '#0969da' }}>
              <TagOutlined />
            </div>
            <div className="stat-value" style={{ color: '#0969da' }}>
              {stats.activeVersions}
            </div>
            <div className="stat-label">活跃版本</div>
          </div>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <div className="stat-card">
            <div className="stat-icon" style={{ color: '#bf8700' }}>
              <RocketOutlined />
            </div>
            <div className="stat-value" style={{ color: '#bf8700' }}>
              {stats.runningDeployments}
            </div>
            <div className="stat-label">进行中的部署</div>
          </div>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <div className="stat-card">
            <div className="stat-icon" style={{ color: '#1a7f37' }}>
              <AppstoreOutlined />
            </div>
            <div className="stat-value" style={{ color: '#1a7f37' }}>
              {stats.totalApplications}
            </div>
            <div className="stat-label">应用总数</div>
          </div>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <div className="stat-card">
            <div className="stat-icon" style={{ color: '#8250df' }}>
              <CloudOutlined />
            </div>
            <div className="stat-value" style={{ color: '#8250df' }}>
              {stats.totalEnvironments}
            </div>
            <div className="stat-label">环境总数</div>
          </div>
        </Col>
      </Row>

      <Card 
        title={
          <span style={{ fontSize: 16, fontWeight: 600 }}>
            最近部署
          </span>
        }
        bordered={false}
        style={{ 
          borderRadius: 8,
          border: '1px solid #e5e7eb',
          boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)'
        }}
      >
        <Table
          columns={columns}
          dataSource={recentDeployments}
          rowKey="id"
          loading={loading}
          pagination={false}
          style={{ fontSize: 14 }}
        />
      </Card>
    </div>
  )
}

export default Dashboard
