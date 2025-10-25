import React, { useEffect, useState, useCallback } from "react";
import { versionApi } from "@/services/api";
import { formatDate } from "@/utils";
import type { Version } from "@/types";
import { IconSearch, IconRefresh } from "@tabler/icons-react";
import { useErrorStore } from "@/store/error";

const Versions: React.FC = () => {
  const { setError } = useErrorStore();
  const [versions, setVersions] = useState<Version[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchText, setSearchText] = useState("");
  const [selectedVersion, setSelectedVersion] = useState<Version | null>(null);

  const loadVersions = useCallback(async () => {
    setLoading(true);
    try {
      const data = await versionApi.list({
        search: searchText || undefined,
      });
      setVersions(data);
    } catch (error) {
      setError("Failed to load versions.");
    } finally {
      setLoading(false);
    }
  }, [searchText, setError]);

  useEffect(() => {
    const timer = setTimeout(() => {
      loadVersions();
    }, 500);
    return () => clearTimeout(timer);
  }, [loadVersions]);

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchText(e.target.value);
  };

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title">版本</h2>
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
              <input type="text" className="form-control" placeholder="搜索版本号或标签..." value={searchText} onChange={handleSearchChange} />
            </div>
            <button className="btn btn-primary ms-2" onClick={loadVersions} disabled={loading}>
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
                <th>应用信息</th>
                <th>创建时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {versions.map((version) => (
                <tr key={version.version}>
                  <td>{version.version}</td>
                  <td>
                    <a href={`https://github.com/XR0118/1024Hackathon/releases/tag/${version.git.tag}`} target="_blank" rel="noopener noreferrer">
                      {version.git.tag}
                    </a>
                  </td>
                  <td>
                    <div className="d-flex flex-column gap-1">
                      {version.applications.map((app) => (
                        <div key={app.name} className="d-flex align-items-center gap-2">
                          <span className="badge bg-secondary-lt">{app.name}</span>
                          <small className="text-muted">
                            覆盖度: {app.coverage}% | 健康度: {app.health}% | 更新: {formatDate(app.lastUpdatedAt)}
                          </small>
                        </div>
                      ))}
                    </div>
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
                  <p>
                    <strong>版本号:</strong> {selectedVersion.version}
                  </p>
                  <p>
                    <strong>Git Tag:</strong> {selectedVersion.git.tag}
                  </p>
                  <p>
                    <strong>创建时间:</strong> {formatDate(selectedVersion.createdAt)}
                  </p>

                  <h3 style={{ marginTop: 24 }}>应用信息</h3>
                  <div className="list-group">
                    {selectedVersion.applications.map((app) => (
                      <div key={app.name} className="list-group-item">
                        <div className="row align-items-center">
                          <div className="col">
                            <strong>{app.name}</strong>
                          </div>
                          <div className="col-auto">
                            <div className="d-flex flex-column gap-1">
                              <div>
                                <span className="text-muted">覆盖度:</span> <span className="badge bg-info-lt">{app.coverage}%</span>
                              </div>
                              <div>
                                <span className="text-muted">健康度:</span>{" "}
                                <span className={`badge ${app.health >= 80 ? "bg-success-lt" : app.health >= 50 ? "bg-warning-lt" : "bg-danger-lt"}`}>
                                  {app.health}%
                                </span>
                              </div>
                              <div>
                                <span className="text-muted">最后更新:</span> {formatDate(app.lastUpdatedAt)}
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Versions;
