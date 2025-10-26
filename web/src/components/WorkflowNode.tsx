import React from "react";
import { Handle, Position } from "reactflow";
import { IconCircleCheck, IconAlertCircle, IconClock, IconLoader, IconArrowUp, IconArrowDown } from "@tabler/icons-react";
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
      case "pending":
        return <IconClock size={20} className="text-secondary" />;
      default:
        return <IconClock size={20} className="text-secondary" />;
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
      case "pending":
        return "待执行";
      default:
        return "未知";
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
      {!data.isFirst && <Handle type="target" position={Position.Top} style={{ background: "#555" }} />}

      <div className="d-flex align-items-center gap-2 mb-2">
        {getStatusIcon()}
        <div className="flex-grow-1">
          <div className="fw-bold">{data.name}</div>
          <small className="text-muted">{getStatusText()}</small>
        </div>
      </div>

      {data.duration && (
        <div className="text-muted" style={{ fontSize: "0.85rem" }}>
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

      {!data.isLast && <Handle type="source" position={Position.Bottom} style={{ background: "#555" }} />}
    </div>
  );
};

export default WorkflowNode;
