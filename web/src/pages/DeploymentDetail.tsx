import React, { useEffect, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { deploymentApi } from "@/services/api";
import { formatDate, formatDuration, getStatusColor, getStatusText } from "@/utils";
import type { DeploymentDetail as DeploymentDetailType } from "@/types";
import { IconArrowLeft, IconCheck, IconArrowBackUp } from "@tabler/icons-react";
import { useErrorStore } from "@/store/error";
import WorkflowViewer from "@/components/WorkflowViewer";

const DeploymentDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [deployment, setDeployment] = useState<DeploymentDetailType | null>(null);
  const [loading, setLoading] = useState(false);
  const [note, setNote] = useState("");
  const [reason, setReason] = useState("");

  const loadDeployment = useCallback(async () => {
    if (!id) return;
    setLoading(true);
    try {
      const data = await deploymentApi.get(id);
      setDeployment(data);
    } catch (error) {
      useErrorStore.getState().setError("Failed to load deployment details.");
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    loadDeployment();
    const interval = setInterval(loadDeployment, 3000);
    return () => clearInterval(interval);
  }, [id, loadDeployment]);

  const handleConfirm = async () => {
    if (!id) return;
    try {
      await deploymentApi.confirm(id, note);
      setNote("");
      loadDeployment();
    } catch (error) {
      useErrorStore.getState().setError("Failed to confirm deployment.");
    }
  };

  const handleRollback = async () => {
    if (!id) return;
    try {
      await deploymentApi.rollback(id, reason);
      setReason("");
      loadDeployment();
    } catch (error) {
      useErrorStore.getState().setError("Failed to rollback deployment.");
    }
  };

  const handleSaveWorkflow = async (tasks: any[]) => {
    // TODO: 调用 API 保存工作流修改
    console.log("保存工作流:", tasks);
    useErrorStore.getState().setError("工作流已保存（演示模式）");
    // 实际项目中应该调用 API
    // await deploymentApi.updateWorkflow(id, tasks);
  };

  if (loading && !deployment) {
    return <div>Loading...</div>;
  }

  if (!deployment) {
    return <div>Deployment not found</div>;
  }

  return (
    <div>
      {/* 紧凑的页面头部 */}
      <div className="page-header d-print-none mb-2">
        <div className="row align-items-center">
          <div className="col">
            <div className="d-flex align-items-center gap-3">
              <button className="btn btn-ghost-secondary btn-sm" onClick={() => navigate("/deployments")}>
                <IconArrowLeft size={18} />
              </button>
              <h2 className="page-title mb-0">部署详情 #{deployment.id}</h2>
              <span className={`badge bg-${getStatusColor(deployment.status)}`}>{getStatusText(deployment.status)}</span>
              {deployment.grayscaleEnabled && <span className="badge bg-azure-lt">灰度 {deployment.grayscaleRatio}%</span>}
            </div>
          </div>
        </div>
      </div>

      {/* 紧凑的基本信息 */}
      <div className="card mb-2" style={{ boxShadow: "none", border: "1px solid #e6e7e9" }}>
        <div className="card-body py-2">
          <div className="d-flex flex-wrap gap-4 align-items-center">
            <div className="d-flex align-items-center gap-2">
              <span className="text-muted" style={{ fontSize: "0.875rem" }}>
                版本:
              </span>
              <strong>{deployment.version}</strong>
            </div>
            <div className="d-flex align-items-center gap-2">
              <span className="text-muted" style={{ fontSize: "0.875rem" }}>
                应用:
              </span>
              <strong>{deployment.applications.join(", ")}</strong>
            </div>
            <div className="d-flex align-items-center gap-2">
              <span className="text-muted" style={{ fontSize: "0.875rem" }}>
                环境:
              </span>
              <strong>{deployment.environments.join(", ")}</strong>
            </div>
            <div className="d-flex align-items-center gap-2">
              <span className="text-muted" style={{ fontSize: "0.875rem" }}>
                创建时间:
              </span>
              <span>{formatDate(deployment.createdAt)}</span>
            </div>
            {deployment.duration && (
              <div className="d-flex align-items-center gap-2">
                <span className="text-muted" style={{ fontSize: "0.875rem" }}>
                  执行时长:
                </span>
                <span>{formatDuration(deployment.duration)}</span>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* 人工确认提示 */}
      {deployment.status === "waiting_confirm" && (
        <div className="alert alert-warning d-flex justify-content-between align-items-center mb-2 py-2">
          <div>
            <strong>需要人工确认</strong> - 此部署需要人工确认后才能继续
          </div>
          <div className="d-flex gap-2">
            <button className="btn btn-success btn-sm" data-bs-toggle="modal" data-bs-target="#confirmModal">
              <IconCheck size={16} className="me-1" />
              确认继续
            </button>
            <button className="btn btn-danger btn-sm" data-bs-toggle="modal" data-bs-target="#rollbackModal">
              <IconArrowBackUp size={16} className="me-1" />
              回滚
            </button>
          </div>
        </div>
      )}

      {/* 工作流区域 - 占据更大空间 */}
      <div className="card" style={{ height: "calc(100vh - 220px)", minHeight: "600px" }}>
        <div className="card-header py-2">
          <h3 className="card-title mb-0">部署流程</h3>
        </div>
        <div className="card-body" style={{ height: "calc(100% - 50px)", overflow: "hidden" }}>
          <WorkflowViewer
            tasks={deployment.tasks}
            onSave={handleSaveWorkflow}
            allowEdit={deployment.status === "pending" || deployment.status === "waiting_confirm"}
          />
        </div>
      </div>

      {/* Modals */}
      <div className="modal" id="confirmModal" tabIndex={-1}>
        <div className="modal-dialog">
          <div className="modal-content">
            <div className="modal-header">
              <h5 className="modal-title">确认部署</h5>
              <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div className="modal-body">
              <p>确认继续此部署吗?</p>
              <textarea className="form-control" placeholder="备注(可选)" value={note} onChange={(e) => setNote(e.target.value)} rows={4}></textarea>
            </div>
            <div className="modal-footer">
              <button type="button" className="btn btn-secondary" data-bs-dismiss="modal">
                取消
              </button>
              <button type="button" className="btn btn-primary" onClick={handleConfirm} data-bs-dismiss="modal">
                确认
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="modal" id="rollbackModal" tabIndex={-1}>
        <div className="modal-dialog">
          <div className="modal-content">
            <div className="modal-header">
              <h5 className="modal-title">回滚部署</h5>
              <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
            </div>
            <div className="modal-body">
              <p>确认回滚此部署吗?</p>
              <textarea className="form-control" placeholder="回滚原因(可选)" value={reason} onChange={(e) => setReason(e.target.value)} rows={4}></textarea>
            </div>
            <div className="modal-footer">
              <button type="button" className="btn btn-secondary" data-bs-dismiss="modal">
                取消
              </button>
              <button type="button" className="btn btn-danger" onClick={handleRollback} data-bs-dismiss="modal">
                回滚
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DeploymentDetailPage;
