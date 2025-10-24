import React from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Versions from './pages/Versions'
import Applications from './pages/Applications'
import Environments from './pages/Environments'
import Deployments from './pages/Deployments'
import DeploymentDetail from './pages/DeploymentDetail'
import CreateDeployment from './pages/CreateDeployment'

const App: React.FC = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<Dashboard />} />
            <Route path="versions" element={<Versions />} />
            <Route path="applications" element={<Applications />} />
            <Route path="environments" element={<Environments />} />
            <Route path="deployments" element={<Deployments />} />
            <Route path="deployments/new" element={<CreateDeployment />} />
            <Route path="deployments/:id" element={<DeploymentDetail />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App
