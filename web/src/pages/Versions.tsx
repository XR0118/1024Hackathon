import React, { useEffect, useState, useCallback } from "react";
import { versionApi, applicationApi } from "@/services/api";
import { formatDate } from "@/utils";
import type { Version, VersionCoverageResponse } from "@/types";
import { IconSearch, IconRefresh, IconRocket, IconArrowBackUp, IconLoader2 } from "@tabler/icons-react";
import { useErrorStore } from "@/store/error";

// 覆盖率缓存类型
interface CoverageCache {
  [versionAppKey: string]: VersionCoverageResponse | "loading" | "error";
}

const Versions: React.FC = () => {
  const { setError } = useErrorStore();
  const [versions, setVersions] = useState<Version[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchText, setSearchText] = useState("");
  const [selectedVersion, setSelectedVersion] = useState<Version | null>(null);
  const [rollbackVersion, setRollbackVersion] = useState<Version | null>(null);
  const [rollbackReason, setRollbackReason] = useState("");
  const [rollbackLoading, setRollbackLoading] = useState(false);
  const [coverageCache, setCoverageCache] = useState<CoverageCache>({});

  // 分页状态
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(10); // 每页显示10条

  // 加载单个应用版本的覆盖率
  const loadCoverage = useCallback(async (appName: string, version: string) => {
    const cacheKey = `${version}:${appName}`;

    // 使用函数式更新避免依赖 coverageCache
    setCoverageCache((prev) => {
      // 避免重复加载
      if (prev[cacheKey]) {
        return prev;
      }
      return { ...prev, [cacheKey]: "loading" };
    });

    try {
      const coverage = await applicationApi.getVersionCoverage(appName, version);
      setCoverageCache((prev) => ({ ...prev, [cacheKey]: coverage }));
    } catch (error) {
      setCoverageCache((prev) => ({ ...prev, [cacheKey]: "error" }));
    }
  }, []);

  const loadVersions = useCallback(async () => {
    setLoading(true);
    try {
      const data = await versionApi.list({
        repository: searchText || undefined,
        page: 1,
        page_size: 100,
      });
      setVersions(data);

      // 为每个版本的每个应用加载覆盖率
      data.forEach((version) => {
        if (version.app_builds && version.app_builds.length > 0) {
          version.app_builds.forEach((build) => {
            loadCoverage(build.app_name, version.version);
          });
        }
      });
    } catch (error) {
      setError("Failed to load versions.");
    } finally {
      setLoading(false);
    }
  }, [searchText, setError, loadCoverage]);

  useEffect(() => {
    const timer = setTimeout(() => {
      loadVersions();
      setCurrentPage(1); // 搜索时重置到第一页
    }, 500);
    return () => clearTimeout(timer);
  }, [loadVersions]);

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchText(e.target.value);
  };

  const handleRollback = async () => {
    if (!rollbackVersion) return;

    setRollbackLoading(true);
    try {
      await versionApi.rollback(rollbackVersion.version, rollbackReason);
      setRollbackVersion(null);
      setRollbackReason("");
      loadVersions();
    } catch (error) {
      setError("回滚操作失败，请重试。");
    } finally {
      setRollbackLoading(false);
    }
  };

  // 分页计算
  const totalPages = Math.ceil(versions.length / pageSize);
  const startIndex = (currentPage - 1) * pageSize;
  const endIndex = startIndex + pageSize;
  const paginatedVersions = versions.slice(startIndex, endIndex);

  // 页面切换处理
  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    // 滚动到顶部
    window.scrollTo({ top: 0, behavior: "smooth" });
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
                <th>Commit</th>
                <th>应用覆盖率</th>
                <th>创建时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {paginatedVersions.map((ver) => (
                <tr key={ver.version}>
                  <td>
                    <div className="d-flex flex-column">
                      <span className="fw-bold">{ver.version}</span>
                      {ver.status === "revert" && (
                        <span className="badge bg-yellow-lt mt-1" style={{ width: "fit-content" }}>
                          回滚版本
                        </span>
                      )}
                    </div>
                  </td>
                  <td>
                    <span className="text-muted" style={{ fontSize: "12px" }}>
                      {ver.git_commit.substring(0, 8)}
                    </span>
                  </td>
                  <td>
                    <div className="d-flex flex-column gap-1">
                      {ver.app_builds && ver.app_builds.length > 0 ? (
                        ver.app_builds.map((build, idx) => {
                          const cacheKey = `${ver.version}:${build.app_name}`;
                          const coverage = coverageCache[cacheKey];

                          return (
                            <div key={idx} className="d-flex align-items-center gap-2">
                              <a href={`/applications/${build.app_name}`} target="_blank" rel="noopener noreferrer" style={{ textDecoration: "none" }}>
                                <span className="text-secondary">{build.app_name}</span>
                              </a>

                              {/* 覆盖率显示 */}
                              {coverage === "loading" ? (
                                <span className="text-muted" style={{ fontSize: "12px" }}>
                                  <IconLoader2 size={12} className="me-1" style={{ display: "inline" }} />
                                </span>
                              ) : coverage === "error" ? (
                                <span className="text-danger" style={{ fontSize: "12px" }}>
                                  -
                                </span>
                              ) : coverage && typeof coverage === "object" ? (
                                <span
                                  className={coverage.coverage_percent >= 80 ? "text-green" : coverage.coverage_percent >= 50 ? "text-yellow" : "text-orange"}
                                  style={{ fontSize: "12px", fontWeight: 500 }}
                                  title={`覆盖 ${coverage.covered_environments}/${coverage.total_environments} 个环境`}
                                >
                                  {coverage.coverage_percent.toFixed(0)}%
                                </span>
                              ) : null}
                            </div>
                          );
                        })
                      ) : (
                        <span className="text-muted">暂无应用</span>
                      )}
                    </div>
                  </td>
                  <td>{formatDate(ver.created_at)}</td>
                  <td>
                    <div className="d-flex gap-2">
                      <button
                        className="btn btn-sm btn-ghost-primary"
                        data-bs-toggle="modal"
                        data-bs-target="#versionDetailModal"
                        onClick={() => setSelectedVersion(ver)}
                      >
                        详情
                      </button>
                      <a href={`/deployments?version=${ver.version}`} target="_blank" rel="noopener noreferrer" className="btn btn-sm btn-ghost-secondary">
                        <IconRocket size={16} className="me-1" />
                        查看任务
                      </a>
                      <button className="btn btn-sm btn-danger" data-bs-toggle="modal" data-bs-target="#rollbackModal" onClick={() => setRollbackVersion(ver)}>
                        <IconArrowBackUp size={16} className="me-1" />
                        回滚
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {totalPages > 1 && (
          <div className="card-footer">
            <div className="d-flex align-items-center">
              <div className="ms-auto text-muted">
                共 {versions.length} 条记录，第 {currentPage} / {totalPages} 页
              </div>
              <ul className="pagination m-0 ms-auto">
                <li className={`page-item ${currentPage === 1 ? "disabled" : ""}`}>
                  <button className="page-link" onClick={() => handlePageChange(1)} disabled={currentPage === 1}>
                    首页
                  </button>
                </li>
                <li className={`page-item ${currentPage === 1 ? "disabled" : ""}`}>
                  <button className="page-link" onClick={() => handlePageChange(currentPage - 1)} disabled={currentPage === 1}>
                    上一页
                  </button>
                </li>
                {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => {
                  // 只显示当前页附近的页码
                  if (page === 1 || page === totalPages || (page >= currentPage - 2 && page <= currentPage + 2)) {
                    return (
                      <li key={page} className={`page-item ${currentPage === page ? "active" : ""}`}>
                        <button className="page-link" onClick={() => handlePageChange(page)}>
                          {page}
                        </button>
                      </li>
                    );
                  } else if (page === currentPage - 3 || page === currentPage + 3) {
                    return (
                      <li key={page} className="page-item disabled">
                        <span className="page-link">...</span>
                      </li>
                    );
                  }
                  return null;
                })}
                <li className={`page-item ${currentPage === totalPages ? "disabled" : ""}`}>
                  <button className="page-link" onClick={() => handlePageChange(currentPage + 1)} disabled={currentPage === totalPages}>
                    下一页
                  </button>
                </li>
                <li className={`page-item ${currentPage === totalPages ? "disabled" : ""}`}>
                  <button className="page-link" onClick={() => handlePageChange(totalPages)} disabled={currentPage === totalPages}>
                    末页
                  </button>
                </li>
              </ul>
            </div>
          </div>
        )}
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
                    <strong>状态:</strong>{" "}
                    <span className={`badge ${selectedVersion.status === "normal" ? "bg-blue" : "bg-yellow"}`}>
                      {selectedVersion.status === "normal" ? "正常" : "回滚版本"}
                    </span>
                  </p>
                  <p>
                    <strong>Git Tag:</strong> {selectedVersion.git_tag}
                  </p>
                  <p>
                    <strong>Git Commit:</strong> {selectedVersion.git_commit}
                  </p>
                  <p>
                    <strong>仓库:</strong>{" "}
                    <a href={selectedVersion.repository} target="_blank" rel="noopener noreferrer">
                      {selectedVersion.repository}
                    </a>
                  </p>
                  <p>
                    <strong>创建者:</strong> {selectedVersion.created_by}
                  </p>
                  <p>
                    <strong>创建时间:</strong> {formatDate(selectedVersion.created_at)}
                  </p>
                  {selectedVersion.description && (
                    <p>
                      <strong>描述:</strong> {selectedVersion.description}
                    </p>
                  )}

                  <h3 style={{ marginTop: 24 }}>应用构建信息</h3>
                  {selectedVersion.app_builds && selectedVersion.app_builds.length > 0 ? (
                    <div className="list-group">
                      {selectedVersion.app_builds.map((build, idx) => {
                        const cacheKey = `${selectedVersion.version}:${build.app_name}`;
                        const coverage = coverageCache[cacheKey];

                        return (
                          <div key={idx} className="list-group-item">
                            <div className="row align-items-center">
                              <div className="col">
                                <strong>{build.app_name}</strong>
                                <br />
                                <span className="text-muted" style={{ fontSize: "12px" }}>
                                  {build.docker_image}
                                </span>
                              </div>
                              <div className="col-auto">
                                {coverage === "loading" ? (
                                  <span className="badge bg-azure-lt">
                                    <IconLoader2 size={14} className="me-1" />
                                    加载覆盖率...
                                  </span>
                                ) : coverage === "error" ? (
                                  <span className="badge bg-red-lt">覆盖率加载失败</span>
                                ) : coverage && typeof coverage === "object" ? (
                                  <div className="d-flex flex-column gap-1">
                                    <span
                                      className={`badge ${
                                        coverage.coverage_percent >= 80 ? "bg-green" : coverage.coverage_percent >= 50 ? "bg-yellow" : "bg-orange"
                                      }`}
                                    >
                                      环境覆盖率: {coverage.coverage_percent.toFixed(1)}%
                                    </span>
                                    <small className="text-muted">
                                      {coverage.covered_environments}/{coverage.total_environments} 个环境
                                    </small>
                                  </div>
                                ) : null}
                              </div>
                            </div>

                            {/* 环境详情 */}
                            {coverage && typeof coverage === "object" && coverage.environments && (
                              <div className="mt-3">
                                <strong style={{ fontSize: "12px" }}>环境详情:</strong>
                                <div className="mt-2" style={{ fontSize: "11px" }}>
                                  {coverage.environments.map((env, envIdx) => (
                                    <div
                                      key={envIdx}
                                      className="d-flex align-items-center justify-content-between mb-1 p-2"
                                      style={{
                                        backgroundColor: env.is_covered ? "#d4edda" : "#f8d7da",
                                        borderRadius: "4px",
                                      }}
                                    >
                                      <span>
                                        <strong>{env.environment.name}</strong>
                                        <span className="text-muted ms-2">({env.environment.type})</span>
                                      </span>
                                      <span>
                                        {env.is_covered ? "✓" : "✗"} {env.coverage_percent.toFixed(0)}%
                                        <span className="text-muted ms-1">
                                          ({env.covered_instances}/{env.total_instances})
                                        </span>
                                      </span>
                                    </div>
                                  ))}
                                </div>
                              </div>
                            )}
                          </div>
                        );
                      })}
                    </div>
                  ) : (
                    <p className="text-muted">暂无构建信息</p>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* 回滚确认对话框 */}
      <div className="modal" id="rollbackModal" tabIndex={-1}>
        <div className="modal-dialog modal-dialog-centered">
          <div className="modal-content">
            <div className="modal-header">
              <h5 className="modal-title">确认回滚操作</h5>
              <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div className="modal-body">
              {rollbackVersion && (
                <div>
                  <div className="alert alert-danger">
                    <strong>警告：此操作不可撤销！</strong>
                    <br />
                    确定要回滚到版本 <strong>{rollbackVersion.version}</strong> 吗？
                  </div>
                  <div className="mb-3">
                    <label className="form-label">回滚原因（必填）</label>
                    <textarea
                      className="form-control"
                      rows={3}
                      placeholder="请填写回滚原因..."
                      value={rollbackReason}
                      onChange={(e) => setRollbackReason(e.target.value)}
                      required
                    />
                  </div>
                </div>
              )}
            </div>
            <div className="modal-footer">
              <button type="button" className="btn btn-secondary" data-bs-dismiss="modal">
                取消
              </button>
              <button
                type="button"
                className="btn btn-danger"
                onClick={handleRollback}
                disabled={rollbackLoading || !rollbackReason.trim()}
                data-bs-dismiss="modal"
              >
                {rollbackLoading ? "回滚中..." : "确认回滚"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Versions;
