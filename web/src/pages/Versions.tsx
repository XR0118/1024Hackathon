import React, { useEffect, useState, useCallback } from 'react'
import { versionApi } from '@/services/api'
import { formatDate } from '@/utils'
import type { Version } from '@/types'
import { IconSearch, IconRefresh, IconAlertTriangle } from '@tabler/icons-react'
import { useErrorStore } from '@/store/error'

const Versions: React.FC = () => {
  const { setError } = useErrorStore();
  const [versions, setVersions] = useState<Version[]>([])
  const [loading, setLoading] = useState(false)
  const [searchText, setSearchText] = useState('')
  const [filterRevert, setFilterRevert] = useState<string>('')
  const [selectedVersion, setSelectedVersion] = useState<Version | null>(null)

  const loadVersions = useCallback(async () => {
    setLoading(true)
    try {
      const data = await versionApi.list({
        search: searchText || undefined,
        isRevert: filterRevert ? filterRevert === 'true' : undefined,
      })
      setVersions(data)
    } catch (error) {
      setError('Failed to load versions.')
    } finally {
      setLoading(false)
    }
  }, [searchText, filterRevert, setError])

  useEffect(() => {
    const timer = setTimeout(() => {
      loadVersions()
    }, 500) // Debounce search input
    return () => clearTimeout(timer)
  }, [loadVersions])

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchText(e.target.value)
  }

  const handleFilterChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setFilterRevert(e.target.value)
  }

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title">版本管理</h2>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <div className="d-flex">
            <div className="input-icon">
              <span className="input-icon-addon">
                <IconSearch />
              </span>
              <input
                type="text"
                className="form-control"
                placeholder="搜索版本号或标签..."
                value={searchText}
                onChange={handleSearchChange}
              />
            </div>
            <select className="form-select ms-2" value={filterRevert} onChange={handleFilterChange}>
              <option value="">所有类型</option>
              <option value="false">正常版本</option>
              <option value="true">回滚版本</option>
            </select>
            <button
              className="btn btn-primary ms-2"
              onClick={loadVersions}
              disabled={loading}
            >
              <IconRefresh className="icon" />
              刷新
            </button>
          </div>
        </div>
        <div className="table-responsive">
          <table className="table card-table table-vcenter text-nowrap datatable">
            <thead>
              <tr>
                <th>版本号</th>
                <th>Git Tag</th>
                <th>关联PR</th>
                <th>状态</th>
                <th>创建时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {versions.map((version) => (
                <tr key={version.id}>
                  <td>
                    {version.isRevert && (
                      <IconAlertTriangle
                        className="text-warning"
                        style={{ marginRight: 8 }}
                      />
                    )}
                    {version.version}
                  </td>
                  <td>
                    <a
                      href={`https://github.com/your-org/your-repo/releases/tag/${version.gitTag}`}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      {version.gitTag}
                    </a>
                  </td>
                  <td>
                    {version.relatedPR ? (
                      <a
                        href={version.relatedPR}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        查看PR
                      </a>
                    ) : (
                      '-'
                    )}
                  </td>
                  <td>
                    <span className="badge bg-success-lt">{version.status}</span>
                  </td>
                  <td>{formatDate(version.createdAt)}</td>
                  <td>
                    <button
                      className="btn btn-sm btn-ghost-primary"
                      data-bs-toggle="modal"
                      data-bs-target="#versionDetailModal"
                      onClick={() => setSelectedVersion(version)}
                    >
                      详情
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="modal" id="versionDetailModal" tabIndex={-1}>
        <div className="modal-dialog">
          <div className="modal-content">
            <div className="modal-header">
              <h5 className="modal-title">版本详情</h5>
              <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div className="modal-body">
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
                    <span className="badge bg-secondary-lt me-1" key={app}>{app}</span>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Versions
