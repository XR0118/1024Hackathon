import React, { useEffect, useState, useCallback } from "react";
import { environmentApi } from "@/services/api";
import type { Environment } from "@/types";
import { IconCloud, IconPlus } from "@tabler/icons-react";
import { useErrorStore } from "@/store/error";
import { getEnvironmentTypeDisplay, getEnvironmentTypeBadgeColor, getEnvironmentStatusDisplay, getEnvironmentStatusBadgeColor } from "@/utils";

const Environments: React.FC = () => {
  const { setError } = useErrorStore();
  const [environments, setEnvironments] = useState<Environment[]>([]);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    type: "kubernetes" as "kubernetes" | "physical",
    is_active: true,
  });

  const loadEnvironments = useCallback(async () => {
    try {
      const data = await environmentApi.list();
      setEnvironments(data);
    } catch (error) {
      setError("Failed to load environments.");
    }
  }, [setError]);

  useEffect(() => {
    loadEnvironments();
  }, [loadEnvironments]);

  const handleOpenModal = () => {
    setShowCreateModal(true);
    setFormData({
      name: "",
      type: "kubernetes",
      is_active: true,
    });
  };

  const handleCloseModal = () => {
    setShowCreateModal(false);
    setFormData({
      name: "",
      type: "kubernetes",
      is_active: true,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.name.trim()) {
      setError("环境名称不能为空");
      return;
    }

    setIsSubmitting(true);
    try {
      await environmentApi.create(formData);
      handleCloseModal();
      loadEnvironments();
    } catch (error) {
      setError(error instanceof Error ? error.message : "创建环境失败");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title">运行环境</h2>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <button className="btn btn-primary" onClick={handleOpenModal}>
              <IconPlus className="icon" />
              新建环境
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
                    <span className={`badge bg-${getEnvironmentTypeBadgeColor(env.type)}-lt`}>{getEnvironmentTypeDisplay(env.type)}</span>
                  </td>
                  <td>
                    <span className={`badge bg-${getEnvironmentStatusBadgeColor(env.is_active)}-lt`}>{getEnvironmentStatusDisplay(env.is_active)}</span>
                  </td>
                  <td>
                    <button className="btn btn-sm btn-ghost-primary">查看详情</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* 新建环境模态框 */}
      {showCreateModal && (
        <div className="modal modal-blur fade show" style={{ display: "block" }} onClick={handleCloseModal}>
          <div className="modal-dialog modal-dialog-centered" onClick={(e) => e.stopPropagation()}>
            <div className="modal-content">
              <div className="modal-header">
                <h5 className="modal-title">新建运行环境</h5>
                <button type="button" className="btn-close" onClick={handleCloseModal}></button>
              </div>
              <form onSubmit={handleSubmit}>
                <div className="modal-body">
                  <div className="mb-3">
                    <label className="form-label required">环境名称</label>
                    <input
                      type="text"
                      className="form-control"
                      placeholder="例如: production, staging"
                      value={formData.name}
                      onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                      required
                      disabled={isSubmitting}
                    />
                  </div>
                  <div className="mb-3">
                    <label className="form-label required">环境类型</label>
                    <select
                      className="form-select"
                      value={formData.type}
                      onChange={(e) => setFormData({ ...formData, type: e.target.value as "kubernetes" | "physical" })}
                      disabled={isSubmitting}
                    >
                      <option value="kubernetes">Kubernetes</option>
                      <option value="physical">物理机</option>
                    </select>
                  </div>
                  <div className="mb-3">
                    <label className="form-check">
                      <input
                        type="checkbox"
                        className="form-check-input"
                        checked={formData.is_active}
                        onChange={(e) => setFormData({ ...formData, is_active: e.target.checked })}
                        disabled={isSubmitting}
                      />
                      <span className="form-check-label">启用环境</span>
                    </label>
                    <small className="form-hint">只有启用的环境才能用于部署</small>
                  </div>
                </div>
                <div className="modal-footer">
                  <button type="button" className="btn btn-link link-secondary" onClick={handleCloseModal} disabled={isSubmitting}>
                    取消
                  </button>
                  <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
                    {isSubmitting ? "创建中..." : "创建环境"}
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

export default Environments;
