import React, { useEffect, useState, useCallback } from "react";
import { Link } from "react-router-dom";
import { IconRocket, IconApps, IconCloud, IconTag, IconRefresh } from "@tabler/icons-react";
import { dashboardApi, deploymentApi } from "@/services/api";
import { formatDate, getStatusColor, getStatusText } from "@/utils";
import type { DeploymentDetail, DashboardStats, Task } from "@/types";
import { useErrorStore } from "@/store/error";
import WorkflowViewer from "@/components/WorkflowViewer";

const Dashboard: React.FC = () => {
  const { setError } = useErrorStore();
  const [stats, setStats] = useState<DashboardStats>({
    activeVersions: 0,
    runningDeployments: 0,
    totalApplications: 0,
    totalEnvironments: 0,
  });
  const [runningDeployments, setRunningDeployments] = useState<DeploymentDetail[]>([]);
  const [mergedTasks, setMergedTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(false);

  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const [statsData, deploymentsData] = await Promise.all([
        dashboardApi.getStats(),
        // 使用统一的 deploymentApi.list，通过参数过滤状态
        deploymentApi.list({
          status: "running,paused",
        }),
      ]);
      setStats(statsData);
      setRunningDeployments(deploymentsData as DeploymentDetail[]);

      // 合并所有部署的任务到一个 DAG 中
      const allTasks: Task[] = [];
      (deploymentsData as DeploymentDetail[]).forEach((deployment) => {
        deployment.tasks?.forEach((task) => {
          allTasks.push({
            ...task,
            // 为任务 ID 添加部署前缀以避免冲突
            id: `${deployment.id}-${task.id}`,
            // 更新依赖关系的 ID
            dependencies: task.dependencies?.map((depId) => `${deployment.id}-${depId}`),
            // 记录所属部署
            deployment_id: deployment.id,
            // 在任务名称中添加部署标识
            name: `[${typeof deployment.version === "string" ? deployment.version : deployment.version?.version || deployment.version_id}] ${task.name}`,
          });
        });
      });
      setMergedTasks(allTasks);
    } catch (error) {
      setError("Failed to load dashboard data.");
    } finally {
      setLoading(false);
    }
  }, [setError]);

  useEffect(() => {
    loadData();
    // 每 3 秒刷新一次，实时监控部署状态
    const interval = setInterval(loadData, 3000);
    return () => clearInterval(interval);
  }, [loadData]);

  return (
    <div>
      {/* 紧凑的页面头部 */}
      <div className="page-header d-print-none mb-2">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title mb-0">首页</h2>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <button className="btn btn-primary btn-sm" onClick={loadData} disabled={loading}>
              <IconRefresh size={16} className="me-1" />
              刷新
            </button>
          </div>
        </div>
      </div>

      {/* 紧凑的统计卡片 */}
      <div className="row row-deck row-cards g-2">
        <div className="col-sm-6 col-lg-3">
          <div className="card">
            <div className="card-body py-2">
              <div className="d-flex align-items-center justify-content-between">
                <div>
                  <div className="text-muted small">活跃版本</div>
                  <div className="h2 mb-0">{stats.activeVersions}</div>
                </div>
                <IconTag size={32} className="text-primary opacity-50" />
              </div>
            </div>
          </div>
        </div>
        <div className="col-sm-6 col-lg-3">
          <div className="card">
            <div className="card-body py-2">
              <div className="d-flex align-items-center justify-content-between">
                <div>
                  <div className="text-muted small">进行中的部署</div>
                  <div className="h2 mb-0">{stats.runningDeployments}</div>
                </div>
                <IconRocket size={32} className="text-warning opacity-50" />
              </div>
            </div>
          </div>
        </div>
        <div className="col-sm-6 col-lg-3">
          <div className="card">
            <div className="card-body py-2">
              <div className="d-flex align-items-center justify-content-between">
                <div>
                  <div className="text-muted small">应用总数</div>
                  <div className="h2 mb-0">{stats.totalApplications}</div>
                </div>
                <IconApps size={32} className="text-success opacity-50" />
              </div>
            </div>
          </div>
        </div>
        <div className="col-sm-6 col-lg-3">
          <div className="card">
            <div className="card-body py-2">
              <div className="d-flex align-items-center justify-content-between">
                <div>
                  <div className="text-muted small">环境总数</div>
                  <div className="h2 mb-0">{stats.totalEnvironments}</div>
                </div>
                <IconCloud size={32} className="text-info opacity-50" />
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 进行中的部署 */}
      {runningDeployments.length > 0 ? (
        <>
          {/* 部署列表 */}
          <div className="card mt-2">
            <div className="card-header py-2">
              <h3 className="card-title mb-0">进行中的部署</h3>
            </div>
            <div className="table-responsive">
              <table className="table table-sm card-table table-vcenter text-nowrap datatable mb-0">
                <thead>
                  <tr>
                    <th>版本号</th>
                    <th>应用</th>
                    <th>目标环境</th>
                    <th>状态</th>
                    <th>进度</th>
                    <th>创建时间</th>
                    <th>操作</th>
                  </tr>
                </thead>
                <tbody>
                  {runningDeployments.map((deployment) => {
                    // 计算进度：已完成任务数 / 总任务数
                    const totalTasks = deployment.tasks?.length || 0;
                    const completedTasks = deployment.tasks?.filter((t) => t.status === "success").length || 0;
                    const progress = totalTasks > 0 ? Math.round((completedTasks / totalTasks) * 100) : 0;

                    // 从 version.app_builds 中获取应用列表
                    const apps =
                      typeof deployment.version === "object" && deployment.version?.app_builds
                        ? deployment.version.app_builds.map((build) => build.app_name)
                        : [];

                    return (
                      <tr key={deployment.id}>
                        <td>
                          <strong>{typeof deployment.version === "string" ? deployment.version : deployment.version?.version || deployment.version_id}</strong>
                        </td>
                        <td>
                          <span className="text-muted" style={{ fontSize: "0.875rem" }}>
                            {apps.slice(0, 2).join(", ")}
                            {apps.length > 2 && ` +${apps.length - 2}`}
                          </span>
                        </td>
                        <td>{deployment.environment?.name || deployment.environment_id}</td>
                        <td>
                          <span className={`badge bg-${getStatusColor(deployment.status)}-lt`}>{getStatusText(deployment.status)}</span>
                        </td>
                        <td>
                          <div className="d-flex align-items-center gap-2">
                            <div className="progress" style={{ width: "80px", height: "6px" }}>
                              <div className="progress-bar" style={{ width: `${progress}%` }} />
                            </div>
                            <small className="text-muted">{progress}%</small>
                          </div>
                        </td>
                        <td>
                          <small className="text-muted">{formatDate(deployment.created_at)}</small>
                        </td>
                        <td>
                          <Link to={`/deployments/${deployment.id}`} className="btn btn-sm btn-ghost-primary">
                            详情
                          </Link>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>

          {/* 合并的工作流视图 */}
          <div className="card mt-2" style={{ height: "600px" }}>
            <div className="card-header py-2">
              <h3 className="card-title mb-0">实时工作流视图</h3>
              <div className="card-subtitle text-muted">所有进行中部署的任务合并展示</div>
            </div>
            <div className="card-body" style={{ height: "calc(100% - 50px)", overflow: "hidden" }}>
              <WorkflowViewer tasks={mergedTasks} allowEdit={false} />
            </div>
          </div>
        </>
      ) : (
        <div className="card mt-2">
          <div className="card-body text-center py-5">
            <div className="text-muted">
              <IconRocket size={48} className="mb-3 opacity-50" />
              <h3>暂无进行中的部署</h3>
              <p>所有部署任务已完成或等待开始</p>
              <Link to="/deployments" className="btn btn-primary mt-3">
                查看所有部署
              </Link>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;
