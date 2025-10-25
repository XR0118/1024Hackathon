import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'
import '@tabler/core/dist/css/tabler.min.css'

if (import.meta.env.VITE_USE_MOCK === 'true') {
  const { enableMockMode } = await import('./mocks')
  enableMockMode()
  console.log('%c[Mock Mode] Mock API enabled', 'color: #0ea5e9; font-weight: bold')
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
