import React, { useEffect, useState, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { applicationApi } from "@/services/api";
import type { Application, VersionSummary } from "@/types";
import { IconPlus, IconRocket, IconCloud } from "@tabler/icons-react";
import { useErrorStore } from "@/store/error";
import { getEnvironmentTypeBadgeColor } from "@/utils";

interface ApplicationWithVersions extends Application {
  versions?: VersionSummary[];
}

const Applications: React.FC = () => {
  const navigate = useNavigate();
  const { setError } = useErrorStore();
  const [applications, setApplications] = useState<ApplicationWithVersions[]>([]);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    description: "",
    repository: "",
    type: "microservice" as "microservice" | "monolith",
  });

  const loadApplications = useCallback(async () => {
    try {
      const data = await applicationApi.list();
      setApplications(data);

      // 为每个应用加载版本概要信息
      data.forEach(async (app) => {
        try {
          const versionData = await applicationApi.getVersionsSummary(app.name);
          setApplications((prev) => prev.map((a) => (a.id === app.id ? { ...a, versions: versionData.versions } : a)));
        } catch (error) {
          console.error(`Failed to load versions for ${app.name}:`, error);
        }
      });
    } catch (error) {
      setError("Failed to load applications.");
    }
  }, [setError]);

  useEffect(() => {
    loadApplications();
  }, [loadApplications]);

  const handleOpenModal = () => {
    setShowCreateModal(true);
    setFormData({
      name: "",
      description: "",
      repository: "",
      type: "microservice",
    });
  };

  const handleCloseModal = () => {
    setShowCreateModal(false);
    setFormData({
      name: "",
      description: "",
      repository: "",
      type: "microservice",
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.name.trim()) {
      setError("应用名称不能为空");
      return;
    }
    if (!formData.repository.trim()) {
      setError("Git 仓库地址不能为空");
      return;
    }

    setIsSubmitting(true);
    try {
      await applicationApi.create(formData);
      handleCloseModal();
      loadApplications();
    } catch (error) {
      setError(error instanceof Error ? error.message : "创建应用失败");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title">应用</h2>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <button className="btn btn-primary" onClick={handleOpenModal}>
              <IconPlus className="icon" />
              新建应用
            </button>
          </div>
        </div>
      </div>

      <div className="row row-cards">
        {applications.map((app) => (
          <div className="col-sm-6 col-lg-4" key={app.id}>
            <div className="card">
              <div className="card-body">
                <h3 className="card-title mb-2">{app.name}</h3>
                <p className="text-muted mb-2" style={{ fontSize: "14px" }}>
                  {app.description || "暂无描述"}
                </p>

                {/* 关联环境 */}
                {app.environments && app.environments.length > 0 && (
                  <div className="mb-3">
                    <div className="d-flex align-items-center gap-1 flex-wrap">
                      <IconCloud size={14} className="text-muted" />
                      {app.environments.map((env) => (
                        <span key={env.id} className={`badge bg-${getEnvironmentTypeBadgeColor(env.type)}-lt`} style={{ fontSize: "11px" }}>
                          {env.name}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* 版本信息（概要，包含运行时信息） */}
                {app.versions && app.versions.length > 0 ? (
                  <div className="mt-3">
                    <div className="divide-y">
                      {app.versions.slice(0, 3).map((version, index) => (
                        <div key={index} className="py-2">
                          <div className="d-flex justify-content-between align-items-center mb-2">
                            <div className="d-flex align-items-center">
                              <span className={`badge bg-${version.status === "normal" ? "blue" : "yellow"} me-2`}>{version.version}</span>
                              {version.status === "revert" && (
                                <span className="badge bg-yellow-lt" style={{ fontSize: "11px" }}>
                                  回滚
                                </span>
                              )}
                            </div>
                            {/* 运行时指标 */}
                            <div className="d-flex align-items-center gap-2" style={{ fontSize: "11px" }}>
                              <span className={`badge bg-${version.healthy.level >= 80 ? "success" : version.healthy.level >= 50 ? "warning" : "danger"}-lt`}>
                                健康度 {version.healthy.level}%
                              </span>
                              <span className="badge bg-info-lt">覆盖度 {version.coverage_percent.toFixed(0)}%</span>
                            </div>
                          </div>
                        </div>
                      ))}
                      {app.versions.length > 3 && (
                        <div className="py-2">
                          <small className="text-muted">还有 {app.versions.length - 3} 个版本...</small>
                        </div>
                      )}
                    </div>
                  </div>
                ) : (
                  <div className="text-muted" style={{ fontSize: "12px" }}>
                    暂无版本信息
                  </div>
                )}
              </div>
              <div className="card-footer">
                <button className="btn btn-ghost-secondary" onClick={() => navigate(`/applications/${app.name}`)}>
                  查看详情
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {applications.length === 0 && (
        <div className="empty">
          <div className="empty-icon">
            <IconRocket size={48} />
          </div>
          <p className="empty-title">暂无应用</p>
          <p className="empty-subtitle text-muted">点击上方"新建应用"按钮创建您的第一个应用</p>
        </div>
      )}

      {/* 新建应用模态框 */}
      {showCreateModal && (
        <div className="modal modal-blur fade show" style={{ display: "block" }} onClick={handleCloseModal}>
          <div className="modal-dialog modal-dialog-centered" onClick={(e) => e.stopPropagation()}>
            <div className="modal-content">
              <div className="modal-header">
                <h5 className="modal-title">新建应用</h5>
                <button type="button" className="btn-close" onClick={handleCloseModal}></button>
              </div>
              <form onSubmit={handleSubmit}>
                <div className="modal-body">
                  <div className="mb-3">
                    <label className="form-label required">应用名称</label>
                    <input
                      type="text"
                      className="form-control"
                      placeholder="例如: user-service"
                      value={formData.name}
                      onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                      required
                      disabled={isSubmitting}
                    />
                  </div>
                  <div className="mb-3">
                    <label className="form-label">应用描述</label>
                    <textarea
                      className="form-control"
                      rows={3}
                      placeholder="简要描述应用的功能..."
                      value={formData.description}
                      onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                      disabled={isSubmitting}
                    />
                  </div>
                  <div className="mb-3">
                    <label className="form-label required">Git 仓库地址</label>
                    <input
                      type="text"
                      className="form-control"
                      placeholder="例如: https://github.com/org/repo.git"
                      value={formData.repository}
                      onChange={(e) => setFormData({ ...formData, repository: e.target.value })}
                      required
                      disabled={isSubmitting}
                    />
                  </div>
                  <div className="mb-3">
                    <label className="form-label required">应用类型</label>
                    <select
                      className="form-select"
                      value={formData.type}
                      onChange={(e) => setFormData({ ...formData, type: e.target.value as "microservice" | "monolith" })}
                      disabled={isSubmitting}
                    >
                      <option value="microservice">微服务</option>
                      <option value="monolith">单体应用</option>
                    </select>
                  </div>
                </div>
                <div className="modal-footer">
                  <button type="button" className="btn btn-link link-secondary" onClick={handleCloseModal} disabled={isSubmitting}>
                    取消
                  </button>
                  <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
                    {isSubmitting ? "创建中..." : "创建应用"}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
      {showCreateModal && <div className="modal-backdrop fade show"></div>}
    </div>
  );
};

export default Applications;
