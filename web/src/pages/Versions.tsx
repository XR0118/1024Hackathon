import React, { useEffect, useState } from 'react'
import { Table, Tag, Button, Input, Select, Space, Drawer } from 'antd'
import { SearchOutlined, ReloadOutlined, WarningOutlined } from '@ant-design/icons'
import { versionApi } from '@/services/api'
import { formatDate } from '@/utils'
import type { Version } from '@/types'

const Versions: React.FC = () => {
  const [versions, setVersions] = useState<Version[]>([])
  const [loading, setLoading] = useState(false)
  const [searchText, setSearchText] = useState('')
  const [filterRevert, setFilterRevert] = useState<boolean | undefined>(undefined)
  const [selectedVersion, setSelectedVersion] = useState<Version | null>(null)
  const [drawerVisible, setDrawerVisible] = useState(false)

  const loadVersions = async () => {
    setLoading(true)
    try {
      const data = await versionApi.list({
        search: searchText || undefined,
        isRevert: filterRevert,
      })
      setVersions(data)
    } catch (error) {
      console.error('Failed to load versions:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadVersions()
  }, [searchText, filterRevert])

  const columns = [
    {
      title: '版本号',
      dataIndex: 'version',
      key: 'version',
      render: (version: string, record: Version) => (
        <span>
          {record.isRevert && <WarningOutlined style={{ color: '#faad14', marginRight: 8 }} />}
          {version}
        </span>
      ),
    },
    {
      title: 'Git Tag',
      dataIndex: 'gitTag',
      key: 'gitTag',
      render: (tag: string) => (
        <a href={`https://github.com/your-org/your-repo/releases/tag/${tag}`} target="_blank" rel="noopener noreferrer">
          {tag}
        </a>
      ),
    },
    {
      title: '关联PR',
      dataIndex: 'relatedPR',
      key: 'relatedPR',
      render: (pr?: string) =>
        pr ? (
          <a href={pr} target="_blank" rel="noopener noreferrer">
            查看PR
          </a>
        ) : (
          '-'
        ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          active: 'blue',
          deployed: 'green',
          deploying: 'orange',
        }
        return <Tag color={colorMap[status]}>{status}</Tag>
      },
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (date: string) => formatDate(date),
      sorter: (a: Version, b: Version) =>
        new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Version) => (
        <Button
          type="link"
          onClick={() => {
            setSelectedVersion(record)
            setDrawerVisible(true)
          }}
        >
          详情
        </Button>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <h1>版本管理</h1>
      </div>

      <Space style={{ marginBottom: 16 }} wrap>
        <Input
          placeholder="搜索版本号或标签"
          prefix={<SearchOutlined />}
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          style={{ width: 300 }}
          allowClear
        />
        <Select
          placeholder="筛选类型"
          value={filterRevert}
          onChange={setFilterRevert}
          style={{ width: 150 }}
          allowClear
        >
          <Select.Option value={false}>正常版本</Select.Option>
          <Select.Option value={true}>回滚版本</Select.Option>
        </Select>
        <Button icon={<ReloadOutlined />} onClick={loadVersions} loading={loading}>
          刷新
        </Button>
      </Space>

      <Table
        columns={columns}
        dataSource={versions}
        rowKey="id"
        loading={loading}
      />

      <Drawer
        title="版本详情"
        placement="right"
        onClose={() => setDrawerVisible(false)}
        open={drawerVisible}
        width={600}
      >
        {selectedVersion && (
          <div>
            <h3>基本信息</h3>
            <p><strong>版本号:</strong> {selectedVersion.version}</p>
            <p><strong>Git Tag:</strong> {selectedVersion.gitTag}</p>
            <p><strong>状态:</strong> {selectedVersion.status}</p>
            <p><strong>创建时间:</strong> {formatDate(selectedVersion.createdAt)}</p>
            <p><strong>回滚标记:</strong> {selectedVersion.isRevert ? '是' : '否'}</p>
            
            <h3 style={{ marginTop: 24 }}>包含的应用</h3>
            {selectedVersion.applications.map((app) => (
              <Tag key={app} style={{ marginBottom: 8 }}>{app}</Tag>
            ))}
          </div>
        )}
      </Drawer>
    </div>
  )
}

export default Versions
