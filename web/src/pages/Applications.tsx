import React, { useEffect, useState, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { applicationApi } from "@/services/api";
import type { Application } from "@/types";
import { IconPlus, IconRocket, IconCircleCheck, IconAlertCircle } from "@tabler/icons-react";
import { useErrorStore } from "@/store/error";

const Applications: React.FC = () => {
  const navigate = useNavigate();
  const { setError } = useErrorStore();
  const [applications, setApplications] = useState<Application[]>([]);

  const loadApplications = useCallback(async () => {
    try {
      const data = await applicationApi.list();
      setApplications(data);
    } catch (error) {
      setError("Failed to load applications.");
    }
  }, [setError]);

  useEffect(() => {
    loadApplications();
    // 每5秒自动刷新一次
    const interval = setInterval(loadApplications, 5000);
    return () => clearInterval(interval);
  }, [loadApplications]);

  const getHealthColor = (health: number) => {
    if (health >= 80) return "success";
    if (health >= 50) return "warning";
    return "danger";
  };

  const getHealthIcon = (health: number) => {
    if (health >= 80) return <IconCircleCheck size={16} />;
    return <IconAlertCircle size={16} />;
  };

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
          <div className="col-sm-6 col-lg-4" key={app.name}>
            <div className="card">
              <div className="card-body">
                <div className="d-flex align-items-center mb-3">
                  {app.icon && (
                    <div className="me-3">
                      <img src={app.icon} alt={app.name} style={{ width: "40px", height: "40px" }} />
                    </div>
                  )}
                  <div>
                    <h3 className="card-title mb-1">{app.name}</h3>
                    <p className="text-muted mb-0" style={{ fontSize: "14px" }}>
                      {app.description}
                    </p>
                  </div>
                </div>

                <div className="mt-3">
                  <strong className="mb-2 d-block">版本信息:</strong>
                  {app.versions && app.versions.length > 0 ? (
                    <div className="list-group list-group-flush">
                      {app.versions.slice(0, 3).map((versionInfo, index) => (
                        <div key={index} className="list-group-item px-0 py-2">
                          <div className="d-flex justify-content-between align-items-center">
                            <div className="d-flex align-items-center">
                              <a
                                href={`/versions?search=${versionInfo.version}`}
                                target="_blank"
                                rel="noopener noreferrer"
                                className={`badge bg-${versionInfo.status === "normal" ? "blue" : "yellow"}-lt me-2`}
                                style={{ textDecoration: "none", cursor: "pointer" }}
                              >
                                {versionInfo.version}
                              </a>
                              {versionInfo.status === "revert" && <span className="badge bg-yellow me-2">回滚</span>}
                            </div>
                            <div className="d-flex align-items-center">
                              <span className={`badge bg-${versionInfo.status === "revert" ? "warning" : "info"}-lt me-2`}>
                                覆盖率: {versionInfo.coverage}%
                              </span>
                              <span className={`badge bg-${getHealthColor(versionInfo.health)}-lt me-2`}>
                                {getHealthIcon(versionInfo.health)}
                                <span className="ms-1">健康度: {versionInfo.health}%</span>
                              </span>
                            </div>
                          </div>
                        </div>
                      ))}
                      {app.versions.length > 3 && (
                        <div className="list-group-item px-0 py-2 text-center">
                          <small className="text-muted">还有 {app.versions.length - 3} 个版本</small>
                        </div>
                      )}
                    </div>
                  ) : (
                    <p className="text-muted">暂无版本信息</p>
                  )}
                </div>
              </div>
              <div className="card-footer">
                <div className="d-flex">
                  <button className="btn btn-ghost-secondary" onClick={() => navigate(`/applications/${app.name}`)}>
                    查看详情
                  </button>
                </div>
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
