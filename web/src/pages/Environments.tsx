import React, { useEffect, useState, useCallback } from 'react'
import { environmentApi } from '@/services/api'
import type { Environment } from '@/types'
import { IconCloud, IconPlus } from '@tabler/icons-react'
import { useErrorStore } from '@/store/error'

const Environments: React.FC = () => {
  const { setError } = useErrorStore();
  const [environments, setEnvironments] = useState<Environment[]>([])

  const loadEnvironments = useCallback(async () => {
    try {
      const data = await environmentApi.list()
      setEnvironments(data)
    } catch (error) {
      setError('Failed to load environments.')
    }
  }, [setError])

  useEffect(() => {
    loadEnvironments()
  }, [loadEnvironments])

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title">环境管理</h2>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <button className="btn btn-primary">
              <IconPlus className="icon" />
              添加环境
            </button>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="table-responsive">
          <table className="table card-table table-vcenter text-nowrap datatable">
            <thead>
              <tr>
                <th>环境名称</th>
                <th>类型</th>
                <th>状态</th>
                <th>应用数量</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {environments.map((env) => (
                <tr key={env.id}>
                  <td>
                    <IconCloud size={16} className="me-2" />
                    {env.name}
                  </td>
                  <td>
                    <span
                      className={`badge bg-${
                        env.type === 'k8s' ? 'primary' : 'success'
                      }-lt`}
                    >
                      {env.type === 'k8s' ? 'Kubernetes' : '物理机'}
                    </span>
                  </td>
                  <td>
                    <span
                      className={`badge bg-${
                        env.status === 'active' ? 'success' : 'secondary'
                      }-lt`}
                    >
                      {env.status === 'active' ? '运行中' : '已停止'}
                    </span>
                  </td>
                  <td>{env.applicationCount}</td>
                  <td>
                    <button className="btn btn-sm btn-ghost-primary">
                      查看详情
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}

export default Environments
