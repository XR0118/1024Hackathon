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

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title">应用</h2>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <button className="btn btn-primary">
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
                              <span className={`badge bg-${version.health_percent >= 80 ? "success" : version.health_percent >= 50 ? "warning" : "danger"}-lt`}>
                                健康度 {version.health_percent.toFixed(0)}%
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
    </div>
  );
};

export default Applications;
