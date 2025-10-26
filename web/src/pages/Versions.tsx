import React, { useEffect, useState, useCallback } from "react";
import { versionApi } from "@/services/api";
import { formatDate } from "@/utils";
import type { Version } from "@/types";
import { IconSearch, IconRefresh, IconRocket, IconArrowBackUp } from "@tabler/icons-react";
import { useErrorStore } from "@/store/error";

const Versions: React.FC = () => {
  const { setError } = useErrorStore();
  const [versions, setVersions] = useState<Version[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchText, setSearchText] = useState("");
  const [selectedVersion, setSelectedVersion] = useState<Version | null>(null);
  const [rollbackVersion, setRollbackVersion] = useState<Version | null>(null);
  const [rollbackReason, setRollbackReason] = useState("");
  const [rollbackLoading, setRollbackLoading] = useState(false);

  // 分页状态
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(10); // 每页显示10条

  const loadVersions = useCallback(async () => {
    setLoading(true);
    try {
      const data = await versionApi.list({
        repository: searchText || undefined,
        page: 1,
        page_size: 100,
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
                <th>Git</th>
                <th>应用信息</th>
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
                    <div className="d-flex flex-column">
                      <span className="text-muted" style={{ fontSize: "12px" }}>
                        Tag:{" "}
                        <a href={`${ver.repository}/releases/tag/${ver.git_tag}`} target="_blank" rel="noopener noreferrer">
                          {ver.git_tag}
                        </a>
                      </span>
                      <span className="text-muted" style={{ fontSize: "12px" }}>
                        Commit: {ver.git_commit.substring(0, 8)}
                      </span>
                      <span className="text-muted" style={{ fontSize: "12px" }}>
                        Repo: {ver.repository.split("/").slice(-2).join("/")}
                      </span>
                    </div>
                  </td>
                  <td>
                    <div className="d-flex flex-column gap-1">
                      {ver.app_builds && ver.app_builds.length > 0 ? (
                        ver.app_builds.map((build, idx) => (
                          <div key={idx} className="d-flex align-items-center gap-2">
                            <a href={`/applications/${build.app_name}`} target="_blank" rel="noopener noreferrer">
                              <span className="badge bg-secondary-lt">{build.app_name}</span>
                            </a>
                            <small className="text-muted" style={{ fontSize: "11px" }}>
                              {build.docker_image}
                            </small>
                          </div>
                        ))
                      ) : (
                        <span className="text-muted">暂无构建信息</span>
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
                      {selectedVersion.app_builds.map((build, idx) => (
                        <div key={idx} className="list-group-item">
                          <div className="row">
                            <div className="col">
                              <strong>{build.app_name}</strong>
                            </div>
                            <div className="col-auto">
                              <span className="text-muted" style={{ fontSize: "12px" }}>
                                {build.docker_image}
                              </span>
                            </div>
                          </div>
                        </div>
                      ))}
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
