import React, { useEffect, useState, useCallback } from "react";
import { environmentApi } from "@/services/api";
import type { Environment } from "@/types";
import { IconCloud, IconPlus } from "@tabler/icons-react";
import { useErrorStore } from "@/store/error";
import { getEnvironmentTypeDisplay, getEnvironmentTypeBadgeColor, getEnvironmentStatusDisplay, getEnvironmentStatusBadgeColor } from "@/utils";

const Environments: React.FC = () => {
  const { setError } = useErrorStore();
  const [environments, setEnvironments] = useState<Environment[]>([]);

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

  return (
    <div>
      <div className="page-header d-print-none">
        <div className="row align-items-center">
          <div className="col">
            <h2 className="page-title">运行环境</h2>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <button className="btn btn-primary">
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
    </div>
  );
};

export default Environments;
