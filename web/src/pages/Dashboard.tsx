import React, { useEffect, useState, useCallback } from "react";
import { Link } from "react-router-dom";
import { IconRocket, IconApps, IconCloud, IconTag, IconRefresh } from "@tabler/icons-react";
import { dashboardApi } from "@/services/api";
import { formatDate, getStatusColor, getStatusText } from "@/utils";
import type { Deployment, DashboardStats } from "@/types";
import { useErrorStore } from "@/store/error";

const Dashboard: React.FC = () => {
  const { setError } = useErrorStore();
  const [stats, setStats] = useState<DashboardStats>({
    activeVersions: 0,
    runningDeployments: 0,
    totalApplications: 0,
    totalEnvironments: 0,
  });
  const [recentDeployments, setRecentDeployments] = useState<Deployment[]>([]);
  const [loading, setLoading] = useState(false);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [statsData, deploymentsData] = await Promise.all([dashboardApi.getStats(), dashboardApi.getRecentDeployments(10)]);
      setStats(statsData);
      setRecentDeployments(deploymentsData);
    } catch (error) {
      setError("Failed to load dashboard data.");
    } finally {
      setLoading(false);
    }
  }, [setError]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title">首页</h2>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <div className="btn-list">
              <button className="btn btn-primary d-none d-sm-inline-block" onClick={loadData} disabled={loading}>
                <IconRefresh className="icon" />
                刷新
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="row row-deck row-cards">
        <div className="col-sm-6 col-lg-3">
          <div className="card">
            <div className="card-body">
              <div className="d-flex align-items-center">
                <div className="subheader">活跃版本</div>
              </div>
              <div className="h1 mb-3">{stats.activeVersions}</div>
              <IconTag size={24} className="text-primary" />
            </div>
          </div>
        </div>
        <div className="col-sm-6 col-lg-3">
          <div className="card">
            <div className="card-body">
              <div className="d-flex align-items-center">
                <div className="subheader">进行中的部署</div>
              </div>
              <div className="h1 mb-3">{stats.runningDeployments}</div>
              <IconRocket size={24} className="text-warning" />
            </div>
          </div>
        </div>
        <div className="col-sm-6 col-lg-3">
          <div className="card">
            <div className="card-body">
              <div className="d-flex align-items-center">
                <div className="subheader">应用总数</div>
              </div>
              <div className="h1 mb-3">{stats.totalApplications}</div>
              <IconApps size={24} className="text-success" />
            </div>
          </div>
        </div>
        <div className="col-sm-6 col-lg-3">
          <div className="card">
            <div className="card-body">
              <div className="d-flex align-items-center">
                <div className="subheader">环境总数</div>
              </div>
              <div className="h1 mb-3">{stats.totalEnvironments}</div>
              <IconCloud size={24} className="text-info" />
            </div>
          </div>
        </div>
      </div>

      <div className="card mt-4">
        <div className="card-header">
          <h3 className="card-title">最近部署</h3>
        </div>
        <div className="table-responsive">
          <table className="table card-table table-vcenter text-nowrap datatable">
            <thead>
              <tr>
                <th>版本号</th>
                <th>应用</th>
                <th>目标环境</th>
                <th>状态</th>
                <th>创建时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {recentDeployments.map((deployment) => (
                <tr key={deployment.id}>
                  <td>{deployment.version}</td>
                  <td>{deployment.applications.join(", ")}</td>
                  <td>{deployment.environments.join(", ")}</td>
                  <td>
                    <span className={`badge bg-${getStatusColor(deployment.status)}-lt`}>{getStatusText(deployment.status)}</span>
                  </td>
                  <td>{formatDate(deployment.createdAt)}</td>
                  <td>
                    <Link to={`/deployments/${deployment.id}`} className="btn btn-sm btn-ghost-primary">
                      查看详情
                    </Link>
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

export default Dashboard;
