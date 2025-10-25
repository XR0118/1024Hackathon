import React, { useEffect, useState, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { deploymentApi } from "@/services/api";
import { formatDate, getStatusColor, getStatusText } from "@/utils";
import type { Deployment } from "@/types";
import { IconPlus, IconRefresh } from "@tabler/icons-react";
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

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title">部署任务</h2>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <button className="btn btn-primary" onClick={() => navigate("/deployments/new")}>
              <IconPlus className="icon" />
              新建任务
            </button>
          </div>
        </div>
      </div>

      <div className="card">
        <div className="card-header">
          <div className="d-flex">
            <select className="form-select" value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}>
              <option value="">所有状态</option>
              <option value="pending">待开始</option>
              <option value="running">进行中</option>
              <option value="success">成功</option>
              <option value="failed">失败</option>
              <option value="waiting_confirm">待确认</option>
            </select>
            <input type="date" className="form-control ms-2" value={startDate} onChange={(e) => setStartDate(e.target.value)} />
            <span className="mx-2">to</span>
            <input type="date" className="form-control" value={endDate} onChange={(e) => setEndDate(e.target.value)} />
            <button className="btn btn-primary ms-2" onClick={loadDeployments} disabled={loading}>
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
              {deployments.map((deployment) => (
                <tr key={deployment.id}>
                  <td>{deployment.id}</td>
                  <td>{deployment.version}</td>
                  <td>{deployment.applications.join(", ")}</td>
                  <td>{deployment.environments.join(", ")}</td>
                  <td>
                    <span className={`badge bg-${getStatusColor(deployment.status)}-lt`}>{getStatusText(deployment.status)}</span>
                  </td>
                  <td>
                    {deployment.status === "running" && (
                      <div className="progress">
                        <div
                          className="progress-bar"
                          style={{ width: `${deployment.progress}%` }}
                          role="progressbar"
                          aria-valuenow={deployment.progress}
                          aria-valuemin={0}
                          aria-valuemax={100}
                        >
                          <span className="visually-hidden">{deployment.progress}% Complete</span>
                        </div>
                      </div>
                    )}
                  </td>
                  <td>{formatDate(deployment.createdAt)}</td>
                  <td>
                    <button className="btn btn-sm btn-ghost-primary" onClick={() => navigate(`/deployments/${deployment.id}`)}>
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
  );
};

export default Deployments;
