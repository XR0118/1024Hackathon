import React, { useState } from "react";
import { Handle, Position } from "reactflow";
import { IconCircleCheck, IconAlertCircle, IconClock, IconLoader, IconArrowUp, IconArrowDown, IconMoon, IconUserCheck } from "@tabler/icons-react";
import type { Task } from "@/types";

interface WorkflowNodeProps {
  data: Task & {
    isFirst?: boolean;
    isLast?: boolean;
    onMoveUp?: () => void;
    onMoveDown?: () => void;
    isEditMode?: boolean;
  };
}

const WorkflowNode: React.FC<WorkflowNodeProps> = ({ data }) => {
  const [sleepDuration, setSleepDuration] = useState(data.params?.sleepDuration || 60);

  const getStatusIcon = () => {
    switch (data.status) {
      case "success":
        return <IconCircleCheck size={20} className="text-success" />;
      case "failed":
        return <IconAlertCircle size={20} className="text-danger" />;
      case "running":
        return <IconLoader size={20} className="text-primary icon-spin" />;
      case "blocked":
        return <IconAlertCircle size={20} className="text-warning" />;
      case "cancelled":
        return <IconAlertCircle size={20} className="text-secondary" />;
      case "waiting_approval":
        return <IconUserCheck size={20} className="text-warning" />;
      case "pending":
        return <IconClock size={20} className="text-secondary" />;
      default:
        return <IconClock size={20} className="text-secondary" />;
    }
  };

  const getTypeIcon = () => {
    switch (data.type) {
      case "sleep":
        return <IconMoon size={16} className="text-info" />;
      case "approval":
        return <IconUserCheck size={16} className="text-warning" />;
      default:
        return null;
    }
  };

  const getStatusColor = () => {
    switch (data.status) {
      case "success":
        return "bg-success-lt border-success";
      case "failed":
        return "bg-danger-lt border-danger";
      case "running":
        return "bg-primary-lt border-primary";
      case "blocked":
        return "bg-warning-lt border-warning";
      case "cancelled":
        return "bg-secondary-lt border-secondary";
      case "waiting_approval":
        return "bg-warning-lt border-warning";
      case "pending":
        return "bg-secondary-lt border-secondary";
      default:
        return "bg-secondary-lt border-secondary";
    }
  };

  const getStatusText = () => {
    switch (data.status) {
      case "success":
        return "成功";
      case "failed":
        return "失败";
      case "running":
        return "运行中";
      case "blocked":
        return "被阻塞";
      case "cancelled":
        return "已取消";
      case "waiting_approval":
        return "等待确认";
      case "pending":
        return "待执行";
      default:
        return "未知";
    }
  };

  const handleApprove = (e: React.MouseEvent) => {
    e.stopPropagation();
    // TODO: 调用 API 确认审批
    console.log("确认审批任务:", data.id);
  };

  const handleSleepDurationChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    e.stopPropagation();
    const value = parseInt(e.target.value) || 0;
    setSleepDuration(value);
    // TODO: 更新任务参数
    if (data.params) {
      data.params.sleepDuration = value;
    }
  };

  return (
    <div
      className={`workflow-node ${getStatusColor()}`}
      style={{
        padding: "12px 16px",
        borderRadius: "8px",
        border: "2px solid",
        minWidth: "200px",
        backgroundColor: "white",
        boxShadow: data.status === "running" ? "0 0 0 3px rgba(32, 107, 196, 0.1)" : "0 2px 4px rgba(0,0,0,0.1)",
      }}
    >
      {/* 编辑模式下始终显示输入连接点，非编辑模式下只在有依赖时显示 */}
      {(data.isEditMode || (data.dependencies && data.dependencies.length > 0)) && (
        <Handle type="target" position={Position.Left} style={{ background: "#555" }} />
      )}

      <div className="d-flex align-items-center gap-2 mb-2">
        {getStatusIcon()}
        <div className="flex-grow-1">
          <div className="d-flex align-items-center gap-1">
            <span className="fw-bold">{data.name}</span>
            {getTypeIcon()}
          </div>
          <small className="text-muted d-block">{getStatusText()}</small>
          {data.appId && (
            <small className="text-primary d-block" style={{ fontSize: "0.8rem" }}>
              📦 {data.appId}
            </small>
          )}
        </div>
      </div>

      {/* Sleep 任务参数 */}
      {data.type === "sleep" && (
        <div className="mt-2" style={{ fontSize: "0.85rem" }}>
          <div className="d-flex align-items-center gap-2">
            <span className="text-muted">等待时间:</span>
            {data.isEditMode ? (
              <input
                type="number"
                className="form-control form-control-sm"
                style={{ width: "80px" }}
                value={sleepDuration}
                onChange={handleSleepDurationChange}
                min="1"
              />
            ) : (
              <span>{sleepDuration}秒</span>
            )}
          </div>
        </div>
      )}

      {/* Approval 任务说明和操作 */}
      {data.type === "approval" && (
        <div className="mt-2">
          {data.params?.approvalNote && (
            <div className="text-muted mb-2" style={{ fontSize: "0.85rem" }}>
              {data.params.approvalNote}
            </div>
          )}
          {data.status === "waiting_approval" && (
            <button className="btn btn-sm btn-success w-100" onClick={handleApprove}>
              <IconUserCheck size={14} className="me-1" />
              确认
            </button>
          )}
        </div>
      )}

      {data.duration && (
        <div className="text-muted mt-2" style={{ fontSize: "0.85rem" }}>
          耗时: {data.duration}s
        </div>
      )}

      {data.status === "running" && (
        <div className="progress mt-2" style={{ height: "4px" }}>
          <div className="progress-bar progress-bar-striped progress-bar-animated" style={{ width: "100%" }} />
        </div>
      )}

      {data.isEditMode && (
        <div className="d-flex gap-1 mt-2">
          {!data.isFirst && (
            <button
              className="btn btn-sm btn-outline-primary"
              onClick={(e) => {
                e.stopPropagation();
                data.onMoveUp?.();
              }}
              style={{ fontSize: "0.75rem", padding: "2px 6px" }}
            >
              <IconArrowUp size={12} />
              上移
            </button>
          )}
          {!data.isLast && (
            <button
              className="btn btn-sm btn-outline-primary"
              onClick={(e) => {
                e.stopPropagation();
                data.onMoveDown?.();
              }}
              style={{ fontSize: "0.75rem", padding: "2px 6px" }}
            >
              <IconArrowDown size={12} />
              下移
            </button>
          )}
        </div>
      )}

      {/* 始终显示输出连接点，允许创建下游连接 */}
      <Handle type="source" position={Position.Right} style={{ background: "#555" }} />
    </div>
  );
};

export default WorkflowNode;
