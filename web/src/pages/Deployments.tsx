import React, { useEffect, useState, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { deploymentApi } from "@/services/api";
import { formatDate, getStatusColor, getStatusText } from "@/utils";
import type { Deployment } from "@/types";
import { IconRefresh } from "@tabler/icons-react";
import { useErrorStore } from "@/store/error";

const Deployments: React.FC = () => {
  const navigate = useNavigate();
  const [deployments, setDeployments] = useState<Deployment[]>([]);
  const [loading, setLoading] = useState(false);
  const [statusFilter, setStatusFilter] = useState<string>("");
  const [startDate, setStartDate] = useState<string>("");
  const [endDate, setEndDate] = useState<string>("");

  const loadDeployments = useCallback(async () => {
    setLoading(true);
    try {
      const data = await deploymentApi.list({
        status: statusFilter || undefined,
        startDate: startDate || undefined,
        endDate: endDate || undefined,
      });
      setDeployments(data);
    } catch (error) {
      useErrorStore.getState().setError("Failed to load deployments.");
    } finally {
      setLoading(false);
    }
  }, [statusFilter, startDate, endDate]);

  useEffect(() => {
    loadDeployments();
    const interval = setInterval(loadDeployments, 5000);
    return () => clearInterval(interval);
  }, [loadDeployments]);

  const renderApplications = (applications: string[]) => {
    const maxDisplay = 3;
    const displayApps = applications.slice(0, maxDisplay);
    const remainingCount = applications.length - maxDisplay;

    return (
      <div className="d-flex flex-wrap gap-1" title={applications.join(", ")}>
        {displayApps.map((app, index) => (
          <a
            key={index}
            href={`/applications/${app}`}
            target="_blank"
            rel="noopener noreferrer"
            className="badge bg-blue-lt"
            style={{ textDecoration: "none", cursor: "pointer" }}
          >
            {app}
          </a>
        ))}
        {remainingCount > 0 && (
          <span className="badge bg-secondary-lt" title={applications.slice(maxDisplay).join(", ")}>
            +{remainingCount}
          </span>
        )}
      </div>
    );
  };

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title">部署任务</h2>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <div className="d-flex align-items-center gap-2">
            <select className="form-select" style={{ width: "auto" }} value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}>
              <option value="">所有状态</option>
              <option value="completed">完成</option>
              <option value="running">运行中</option>
              <option value="paused">暂停中</option>
              <option value="pending">待开始</option>
            </select>
            <input type="date" className="form-control" style={{ width: "auto" }} value={startDate} onChange={(e) => setStartDate(e.target.value)} />
            <span className="text-nowrap">至</span>
            <input type="date" className="form-control" style={{ width: "auto" }} value={endDate} onChange={(e) => setEndDate(e.target.value)} />
            <button className="btn btn-primary" onClick={loadDeployments} disabled={loading}>
              <IconRefresh className="icon" />
              刷新
            </button>
          </div>
        </div>
        <div className="table-responsive">
          <table className="table card-table table-vcenter text-nowrap datatable">
            <thead>
              <tr>
                <th>部署ID</th>
                <th>版本</th>
                <th>应用</th>
                <th>环境</th>
                <th>状态</th>
                <th>进度</th>
                <th>创建时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {deployments.map((deployment) => {
                const apps = deployment.must_in_order || [];
                return (
                  <tr key={deployment.id}>
                    <td>{deployment.id}</td>
                    <td>{deployment.version || deployment.version_id}</td>
                    <td>{renderApplications(apps)}</td>
                    <td>{deployment.environment?.name || deployment.environment_id}</td>
                    <td>
                      <span className={`badge bg-${getStatusColor(deployment.status)}-lt`}>{getStatusText(deployment.status)}</span>
                    </td>
                    <td>{deployment.status === "running" && <span className="badge bg-primary">运行中</span>}</td>
                    <td>{formatDate(deployment.created_at)}</td>
                    <td>
                      <button className="btn btn-sm btn-ghost-primary" onClick={() => navigate(`/deployments/${deployment.id}`)}>
                        查看详情
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default Deployments;
