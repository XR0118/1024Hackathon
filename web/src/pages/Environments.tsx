import React, { useEffect, useState } from 'react'
import { Table, Tag, Button, Badge } from 'antd'
import { CloudOutlined, PlusOutlined } from '@ant-design/icons'
import { environmentApi } from '@/services/api'
import type { Environment } from '@/types'

const Environments: React.FC = () => {
  const [environments, setEnvironments] = useState<Environment[]>([])
  const [loading, setLoading] = useState(false)

  const loadEnvironments = async () => {
    setLoading(true)
    try {
      const data = await environmentApi.list()
      setEnvironments(data)
    } catch (error) {
      console.error('Failed to load environments:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadEnvironments()
  }, [])

  const columns = [
    {
      title: '环境名称',
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => (
        <span>
          <CloudOutlined style={{ marginRight: 8 }} />
          {name}
        </span>
      ),
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color={type === 'k8s' ? 'blue' : 'green'}>
          {type === 'k8s' ? 'Kubernetes' : '物理机'}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Badge
          status={status === 'active' ? 'success' : 'default'}
          text={status === 'active' ? '运行中' : '已停止'}
        />
      ),
    },
    {
      title: '应用数量',
      dataIndex: 'applicationCount',
      key: 'applicationCount',
    },
    {
      title: '操作',
      key: 'action',
      render: () => (
        <Button type="link">查看详情</Button>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 24, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0 }}>环境管理</h1>
        <Button type="primary" icon={<PlusOutlined />}>
          添加环境
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={environments}
        rowKey="id"
        loading={loading}
      />
    </div>
  )
}

export default Environments
