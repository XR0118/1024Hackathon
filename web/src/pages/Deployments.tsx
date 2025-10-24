import React, { useEffect, useState } from 'react'
import { Table, Tag, Button, Select, DatePicker, Space, Progress } from 'antd'
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { deploymentApi } from '@/services/api'
import { formatDate, getStatusColor, getStatusText } from '@/utils'
import type { Deployment } from '@/types'

const { RangePicker } = DatePicker

const Deployments: React.FC = () => {
  const navigate = useNavigate()
  const [deployments, setDeployments] = useState<Deployment[]>([])
  const [loading, setLoading] = useState(false)
  const [statusFilter, setStatusFilter] = useState<string | undefined>(undefined)

  const loadDeployments = async () => {
    setLoading(true)
    try {
      const data = await deploymentApi.list({
        status: statusFilter,
      })
      setDeployments(data)
    } catch (error) {
      console.error('Failed to load deployments:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadDeployments()
    const interval = setInterval(loadDeployments, 5000)
    return () => clearInterval(interval)
  }, [statusFilter])

  const columns = [
    {
      title: '部署ID',
      dataIndex: 'id',
      key: 'id',
      width: 100,
    },
    {
      title: '版本',
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
      title: '环境',
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
      title: '进度',
      dataIndex: 'progress',
      key: 'progress',
      render: (progress: number, record: Deployment) =>
        record.status === 'running' ? (
          <Progress percent={progress} size="small" />
        ) : null,
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
        <h1>部署管理</h1>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => navigate('/deployments/new')}
          style={{ borderRadius: 6 }}
        >
          新建部署
        </Button>
      </div>

      <Space style={{ marginBottom: 20 }} wrap>
        <Select
          placeholder="筛选状态"
          value={statusFilter}
          onChange={setStatusFilter}
          style={{ width: 150, borderRadius: 6 }}
          allowClear
        >
          <Select.Option value="pending">待开始</Select.Option>
          <Select.Option value="running">进行中</Select.Option>
          <Select.Option value="success">成功</Select.Option>
          <Select.Option value="failed">失败</Select.Option>
          <Select.Option value="waiting_confirm">待确认</Select.Option>
        </Select>
        <RangePicker style={{ borderRadius: 6 }} />
        <Button 
          icon={<ReloadOutlined />} 
          onClick={loadDeployments} 
          loading={loading}
          style={{ borderRadius: 6 }}
        >
          刷新数据
        </Button>
      </Space>

      <Table
        columns={columns}
        dataSource={deployments}
        rowKey="id"
        loading={loading}
      />
    </div>
  )
}

export default Deployments
