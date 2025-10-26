import React, { useEffect, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { applicationApi } from "@/services/api";
import { formatDate, getEnvironmentTypeBadgeColor } from "@/utils";
import type { Application, ApplicationVersionsDetailResponse } from "@/types";
import { IconArrowLeft, IconCircleCheck, IconAlertCircle, IconServer, IconCloud } from "@tabler/icons-react";
import { useErrorStore } from "@/store/error";

const ApplicationDetail: React.FC = () => {
  const { name } = useParams<{ name: string }>();
  const navigate = useNavigate();
  const { setError } = useErrorStore();
  const [application, setApplication] = useState<Application | null>(null);
  const [versionDetail, setVersionDetail] = useState<ApplicationVersionsDetailResponse | null>(null);
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
      setVersionDetail(response);
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
                {/* 关联环境 */}
                {application.environments && application.environments.length > 0 && (
                  <div className="mt-2">
                    <div className="d-flex align-items-center gap-1 flex-wrap">
                      <IconCloud size={14} className="text-muted" />
                      {application.environments.map((env) => (
                        <span key={env.id} className={`badge bg-${getEnvironmentTypeBadgeColor(env.type)}-lt`} style={{ fontSize: "11px" }}>
                          {env.name}
                        </span>
                      ))}
                    </div>
                  </div>
                )}
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
                  <h3 className="card-title">版本部署信息（按环境）</h3>
                </div>
                <div className="card-body">
                  {versionDetail && versionDetail.environments && versionDetail.environments.length > 0 ? (
                    <div className="row row-cards">
                      {versionDetail.environments.map((envVersions, envIndex) => (
                        <div className="col-12" key={envIndex}>
                          <div className="card mb-3">
                            <div className="card-header">
                              <h4 className="card-title d-flex align-items-center">
                                <IconCloud size={20} className="me-2" />
                                {envVersions.environment.name}
                                <span className={`badge bg-${getEnvironmentTypeBadgeColor(envVersions.environment.type)}-lt ms-2`}>
                                  {envVersions.environment.type}
                                </span>
                              </h4>
                            </div>
                            <div className="card-body">
                              {envVersions.versions && envVersions.versions.length > 0 ? (
                                <div className="accordion" id={`versionsAccordion${envIndex}`}>
                                  {envVersions.versions.map((versionInfo, versionIndex) => {
                                    const accordionId = `env${envIndex}version${versionIndex}`;
                                    return (
                                      <div className="accordion-item" key={versionIndex}>
                                        <h2 className="accordion-header" id={`heading${accordionId}`}>
                                          <button
                                            className={`accordion-button ${versionIndex !== 0 ? "collapsed" : ""}`}
                                            type="button"
                                            data-bs-toggle="collapse"
                                            data-bs-target={`#collapse${accordionId}`}
                                            aria-expanded={versionIndex === 0}
                                            aria-controls={`collapse${accordionId}`}
                                          >
                                            <div className="d-flex justify-content-between align-items-center w-100 me-3">
                                              <div className="d-flex align-items-center">
                                                <span
                                                  className={`badge bg-${
                                                    versionInfo.status === "normal" ? "blue" : versionInfo.status === "abnormal" ? "red" : "yellow"
                                                  } me-3`}
                                                >
                                                  {versionInfo.version}
                                                </span>
                                                {versionInfo.status === "revert" && <span className="badge bg-yellow me-3">回滚版本</span>}
                                                {versionInfo.status === "abnormal" && <span className="badge bg-red me-3">异常</span>}
                                                <span className="badge bg-info-lt me-3">覆盖率: {versionInfo.coverage}%</span>
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
                                          id={`collapse${accordionId}`}
                                          className={`accordion-collapse collapse ${versionIndex === 0 ? "show" : ""}`}
                                          aria-labelledby={`heading${accordionId}`}
                                          data-bs-parent={`#versionsAccordion${envIndex}`}
                                        >
                                          <div className="accordion-body">
                                            <h5 className="mb-3">
                                              <IconServer size={18} className="me-2" />
                                              实例列表
                                            </h5>
                                            {versionInfo.instances && versionInfo.instances.length > 0 ? (
                                              <div className="table-responsive">
                                                <table className="table table-vcenter card-table">
                                                  <thead>
                                                    <tr>
                                                      <th>实例名称</th>
                                                      <th>状态</th>
                                                      <th>健康度</th>
                                                      <th>最后更新时间</th>
                                                    </tr>
                                                  </thead>
                                                  <tbody>
                                                    {versionInfo.instances.map((instance, instanceIndex) => (
                                                      <tr key={instanceIndex}>
                                                        <td>
                                                          <IconServer size={16} className="me-2" />
                                                          {instance.node_name}
                                                        </td>
                                                        <td>
                                                          <span className={`badge bg-${instance.status === "running" ? "success" : "warning"}-lt`}>
                                                            {instance.status}
                                                          </span>
                                                        </td>
                                                        <td>
                                                          <span className={`badge bg-${getHealthColor(instance.health)}-lt`}>
                                                            {getHealthIcon(instance.health)}
                                                            <span className="ms-1">{instance.health}%</span>
                                                          </span>
                                                        </td>
                                                        <td className="text-muted">{formatDate(instance.last_updated_at)}</td>
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
                                    );
                                  })}
                                </div>
                              ) : (
                                <div className="empty">
                                  <p className="empty-title">该环境暂无部署版本</p>
                                </div>
                              )}
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="empty">
                      <p className="empty-title">暂无版本部署信息</p>
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
