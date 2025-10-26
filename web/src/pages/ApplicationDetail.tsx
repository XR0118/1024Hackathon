import React, { useEffect, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { applicationApi } from "@/services/api";
import { formatDate } from "@/utils";
import type { Application, ApplicationVersionInfo } from "@/types";
import { IconArrowLeft, IconCircleCheck, IconAlertCircle, IconServer } from "@tabler/icons-react";
import { useErrorStore } from "@/store/error";

const ApplicationDetail: React.FC = () => {
  const { name } = useParams<{ name: string }>();
  const navigate = useNavigate();
  const { setError } = useErrorStore();
  const [application, setApplication] = useState<Application | null>(null);
  const [versions, setVersions] = useState<ApplicationVersionInfo[]>([]);
  const [loading, setLoading] = useState(true);

  const loadApplication = useCallback(
    async (isInitialLoad = false) => {
      if (!name) return;

      try {
        if (isInitialLoad) {
          setLoading(true);
        }
        const data = await applicationApi.get(name);
        setApplication(data);
      } catch (error) {
        setError("Failed to load application details.");
      } finally {
        if (isInitialLoad) {
          setLoading(false);
        }
      }
    },
    [name, setError]
  );

  const loadVersions = useCallback(async () => {
    if (!name) return;

    try {
      const response = await applicationApi.getVersions(name);
      setVersions(response.versions || []);
    } catch (error) {
      console.error("Failed to load versions:", error);
      // 版本信息加载失败不影响主体显示
    }
  }, [name]);

  useEffect(() => {
    loadApplication(true);
    loadVersions();

    // 每5秒自动刷新一次版本信息
    const interval = setInterval(() => {
      loadApplication(false);
      loadVersions();
    }, 5000);

    return () => clearInterval(interval);
  }, [loadApplication, loadVersions]);

  const getHealthColor = (health: number) => {
    if (health >= 80) return "success";
    if (health >= 50) return "warning";
    return "danger";
  };

  const getHealthIcon = (health: number) => {
    if (health >= 80) return <IconCircleCheck size={16} />;
    return <IconAlertCircle size={16} />;
  };

  if (loading) {
    return (
      <div className="page-body">
        <div className="container-xl">
          <div className="text-center">
            <div className="spinner-border" role="status">
              <span className="visually-hidden">Loading...</span>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!application) {
    return (
      <div className="page-body">
        <div className="container-xl">
          <div className="empty">
            <p className="empty-title">应用不存在</p>
            <div className="empty-action">
              <button className="btn btn-primary" onClick={() => navigate("/applications")}>
                返回应用列表
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <div className="d-flex align-items-center">
              <button className="btn btn-ghost-secondary me-3" onClick={() => navigate("/applications")}>
                <IconArrowLeft size={20} />
              </button>
              <div>
                <h2 className="page-title mb-0">{application.name}</h2>
                <p className="text-muted mb-0">{application.description || "暂无描述"}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="page-body">
        <div className="container-xl">
          <div className="row">
            <div className="col-12">
              <div className="card">
                <div className="card-header">
                  <h3 className="card-title">版本列表</h3>
                </div>
                <div className="card-body">
                  {versions && versions.length > 0 ? (
                    <div className="accordion" id="versionsAccordion">
                      {versions.map((versionInfo, index) => (
                        <div className="accordion-item" key={index}>
                          <h2 className="accordion-header" id={`heading${index}`}>
                            <button
                              className={`accordion-button ${index !== 0 ? "collapsed" : ""}`}
                              type="button"
                              data-bs-toggle="collapse"
                              data-bs-target={`#collapse${index}`}
                              aria-expanded={index === 0}
                              aria-controls={`collapse${index}`}
                            >
                              <div className="d-flex justify-content-between align-items-center w-100 me-3">
                                <div className="d-flex align-items-center">
                                  <span className={`badge bg-${versionInfo.status === "normal" ? "blue" : "yellow"} me-3`}>{versionInfo.version}</span>
                                  {versionInfo.status === "revert" && <span className="badge bg-yellow me-3">回滚版本</span>}
                                  <span className={`badge bg-${versionInfo.status === "revert" ? "warning" : "info"}-lt me-3`}>
                                    覆盖率: {versionInfo.coverage}%
                                  </span>
                                </div>
                                <div className="d-flex align-items-center">
                                  <span className={`badge bg-${getHealthColor(versionInfo.health)}-lt me-3`}>
                                    {getHealthIcon(versionInfo.health)}
                                    <span className="ms-1">健康度: {versionInfo.health}%</span>
                                  </span>
                                  <small className="text-muted">最后更新: {formatDate(versionInfo.last_updated_at)}</small>
                                </div>
                              </div>
                            </button>
                          </h2>
                          <div
                            id={`collapse${index}`}
                            className={`accordion-collapse collapse ${index === 0 ? "show" : ""}`}
                            aria-labelledby={`heading${index}`}
                            data-bs-parent="#versionsAccordion"
                          >
                            <div className="accordion-body">
                              <h4 className="mb-3">
                                <IconServer size={20} className="me-2" />
                                实例列表
                              </h4>
                              {versionInfo.nodes && versionInfo.nodes.length > 0 ? (
                                <div className="table-responsive">
                                  <table className="table table-vcenter card-table">
                                    <thead>
                                      <tr>
                                        <th>实例名称</th>
                                        <th>健康度</th>
                                        <th>最后更新时间</th>
                                      </tr>
                                    </thead>
                                    <tbody>
                                      {versionInfo.nodes.map((node, nodeIndex) => (
                                        <tr key={nodeIndex}>
                                          <td>
                                            <IconServer size={16} className="me-2" />
                                            {node.name}
                                          </td>
                                          <td>
                                            <span className={`badge bg-${getHealthColor(node.health)}-lt`}>
                                              {getHealthIcon(node.health)}
                                              <span className="ms-1">{node.health}%</span>
                                            </span>
                                          </td>
                                          <td className="text-muted">{formatDate(node.last_updated_at)}</td>
                                        </tr>
                                      ))}
                                    </tbody>
                                  </table>
                                </div>
                              ) : (
                                <div className="empty">
                                  <p className="empty-title">暂无实例信息</p>
                                </div>
                              )}
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="empty">
                      <p className="empty-title">暂无版本信息</p>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ApplicationDetail;
