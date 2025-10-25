import React from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Versions from './pages/Versions'
import Applications from './pages/Applications'
import ApplicationDetail from './pages/ApplicationDetail'
import Environments from './pages/Environments'
import Deployments from './pages/Deployments'
import DeploymentDetail from './pages/DeploymentDetail'
import CreateDeployment from './pages/CreateDeployment'

const App: React.FC = () => {
  return (
    <BrowserRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Dashboard />} />
          <Route path="versions" element={<Versions />} />
          <Route path="applications" element={<Applications />} />
          <Route path="applications/:name" element={<ApplicationDetail />} />
          <Route path="environments" element={<Environments />} />
          <Route path="deployments" element={<Deployments />} />
          <Route path="deployments/new" element={<CreateDeployment />} />
          <Route path="deployments/:id" element={<DeploymentDetail />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
