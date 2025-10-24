import React from 'react'
import { Layout as AntLayout, Menu, theme } from 'antd'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import {
  DashboardOutlined,
  TagOutlined,
  AppstoreOutlined,
  CloudOutlined,
  RocketOutlined,
} from '@ant-design/icons'

const { Header, Sider, Content } = AntLayout

const Layout: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken()

  const menuItems = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: '仪表板',
    },
    {
      key: '/versions',
      icon: <TagOutlined />,
      label: '版本管理',
    },
    {
      key: '/applications',
      icon: <AppstoreOutlined />,
      label: '应用管理',
    },
    {
      key: '/environments',
      icon: <CloudOutlined />,
      label: '环境管理',
    },
    {
      key: '/deployments',
      icon: <RocketOutlined />,
      label: '部署管理',
    },
  ]

  return (
    <AntLayout style={{ minHeight: '100vh', background: '#f6f8fa' }}>
      <Sider
        breakpoint="lg"
        collapsedWidth="0"
        style={{
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
          background: '#24292f',
          boxShadow: '2px 0 8px rgba(0, 0, 0, 0.1)',
        }}
        width={220}
      >
        <div
          style={{
            height: 64,
            margin: '16px 0',
            padding: '0 20px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'flex-start',
            color: '#fff',
            fontSize: 18,
            fontWeight: 600,
            borderBottom: '1px solid rgba(255, 255, 255, 0.1)',
          }}
        >
          <RocketOutlined style={{ marginRight: 12, fontSize: 20 }} />
          部署平台
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          style={{
            background: 'transparent',
            border: 'none',
            padding: '8px 12px',
          }}
        />
      </Sider>
      <AntLayout style={{ marginLeft: 220, background: '#f6f8fa' }}>
        <Header
          style={{
            padding: '0 24px',
            background: '#ffffff',
            borderBottom: '1px solid #d1d5db',
            display: 'flex',
            alignItems: 'center',
            boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
            height: 64,
          }}
        >
          <div style={{ fontSize: 14, color: '#57606a' }}>
            七牛云 2024/10/24 Hackathon - 部署管理平台
          </div>
        </Header>
        <Content style={{ margin: '24px', overflow: 'initial' }}>
          <div
            style={{
              padding: 32,
              background: '#ffffff',
              borderRadius: 8,
              minHeight: 'calc(100vh - 112px)',
              border: '1px solid #e5e7eb',
              boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
            }}
          >
            <Outlet />
          </div>
        </Content>
      </AntLayout>
    </AntLayout>
  )
}

export default Layout
