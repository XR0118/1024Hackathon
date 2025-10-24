import React, { useEffect, useState } from 'react'
import { Card, Row, Col, Button, Tag, Space } from 'antd'
import { AppstoreAddOutlined, RocketOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { applicationApi } from '@/services/api'
import { formatDate } from '@/utils'
import type { Application } from '@/types'

const Applications: React.FC = () => {
  const navigate = useNavigate()
  const [applications, setApplications] = useState<Application[]>([])
  const [loading, setLoading] = useState(false)

  const loadApplications = async () => {
    setLoading(true)
    try {
      const data = await applicationApi.list()
      setApplications(data)
    } catch (error) {
      console.error('Failed to load applications:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadApplications()
  }, [])

  return (
    <div>
      <div className="page-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1>应用管理</h1>
        <Button 
          type="primary" 
          icon={<AppstoreAddOutlined />}
          style={{ borderRadius: 6 }}
        >
          添加应用
        </Button>
      </div>

      <Row gutter={[16, 16]}>
        {applications.map((app) => (
          <Col xs={24} sm={12} lg={8} xl={6} key={app.id}>
            <Card
              hoverable
              loading={loading}
              style={{
                borderRadius: 8,
                border: '1px solid #e5e7eb',
                boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
              }}
              actions={[
                <Button
                  type="link"
                  icon={<RocketOutlined />}
                  onClick={() => navigate(`/deployments/new?appId=${app.id}`)}
                >
                  新建部署
                </Button>,
                <Button
                  type="link"
                  onClick={() => navigate(`/applications/${app.id}`)}
                >
                  查看详情
                </Button>,
              ]}
            >
              <Card.Meta
                title={
                  <Space>
                    <AppstoreAddOutlined style={{ fontSize: 20 }} />
                    {app.name}
                  </Space>
                }
                description={
                  <div>
                    <p style={{ marginBottom: 8 }}>{app.description}</p>
                    <div style={{ marginTop: 12 }}>
                      <strong>当前版本:</strong>
                      <div style={{ marginTop: 8 }}>
                        {Object.entries(app.currentVersions).map(([env, version]) => (
                          <Tag key={env} style={{ marginBottom: 4 }}>
                            {env}: {version}
                          </Tag>
                        ))}
                      </div>
                    </div>
                    <p style={{ marginTop: 12, fontSize: 12, color: '#999' }}>
                      最近部署: {formatDate(app.lastDeployedAt)}
                    </p>
                  </div>
                }
              />
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  )
}

export default Applications
