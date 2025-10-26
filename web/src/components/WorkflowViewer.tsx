import React, { useMemo, useState, useCallback } from "react";
import ReactFlow, {
  Node,
  Edge,
  Controls,
  Background,
  MarkerType,
  BackgroundVariant,
  addEdge,
  Connection,
  NodeChange,
  EdgeChange,
  applyNodeChanges,
  applyEdgeChanges,
} from "reactflow";
import "reactflow/dist/style.css";
import WorkflowNode from "./WorkflowNode";
import type { Task } from "@/types";
import { IconEdit, IconDeviceFloppy, IconX, IconPlus, IconTrash } from "@tabler/icons-react";

interface WorkflowViewerProps {
  tasks: Task[];
  onSave?: (tasks: Task[]) => void;
  allowEdit?: boolean;
}

const WorkflowViewer: React.FC<WorkflowViewerProps> = ({ tasks, onSave, allowEdit = true }) => {
  const [isEditMode, setIsEditMode] = useState(false);
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);
  const nodeTypes = useMemo(() => ({ workflowNode: WorkflowNode }), []);

  // 上移节点
  const handleMoveUp = useCallback((nodeId: string) => {
    setNodes((nds) => {
      const index = nds.findIndex((n) => n.id === nodeId);
      if (index <= 0) return nds;

      const newNodes = [...nds];
      [newNodes[index - 1], newNodes[index]] = [newNodes[index], newNodes[index - 1]];

      // 更新位置
      return newNodes.map((node, idx) => ({
        ...node,
        position: { x: 250, y: idx * 150 },
        data: {
          ...node.data,
          isFirst: idx === 0,
          isLast: idx === newNodes.length - 1,
        },
      }));
    });

    // 重建边
    rebuildEdges();
  }, []);

  // 下移节点
  const handleMoveDown = useCallback((nodeId: string) => {
    setNodes((nds) => {
      const index = nds.findIndex((n) => n.id === nodeId);
      if (index < 0 || index >= nds.length - 1) return nds;

      const newNodes = [...nds];
      [newNodes[index], newNodes[index + 1]] = [newNodes[index + 1], newNodes[index]];

      // 更新位置
      return newNodes.map((node, idx) => ({
        ...node,
        position: { x: 250, y: idx * 150 },
        data: {
          ...node.data,
          isFirst: idx === 0,
          isLast: idx === newNodes.length - 1,
        },
      }));
    });

    // 重建边
    rebuildEdges();
  }, []);

  // 重建顺序边
  const rebuildEdges = useCallback(() => {
    setNodes((nds) => {
      const newEdges: Edge[] = [];
      for (let i = 0; i < nds.length - 1; i++) {
        newEdges.push({
          id: `e${nds[i].id}-${nds[i + 1].id}`,
          source: nds[i].id,
          target: nds[i + 1].id,
          type: "smoothstep",
          animated: false,
          style: { stroke: "#8c8c8c", strokeWidth: 2 },
          markerEnd: { type: MarkerType.ArrowClosed, color: "#8c8c8c" },
        });
      }
      setEdges(newEdges);
      return nds;
    });
  }, []);

  // 初始化节点和边
  const initializeNodesAndEdges = useCallback(() => {
    // 简单的 DAG 布局：计算每个节点的层级
    const calculateLevel = (taskId: string, taskMap: Map<string, any>, visited = new Set<string>()): number => {
      if (visited.has(taskId)) return 0; // 避免循环依赖
      visited.add(taskId);

      const task = taskMap.get(taskId);
      if (!task || !task.dependencies || task.dependencies.length === 0) {
        return 0; // 顶点节点
      }

      // 当前节点的层级 = 所有依赖节点的最大层级 + 1
      const maxDependencyLevel = Math.max(...task.dependencies.map((depId: string) => calculateLevel(depId, taskMap, new Set(visited))));
      return maxDependencyLevel + 1;
    };

    const taskMap = new Map(tasks.map((task) => [task.id, task]));
    const levels = new Map<string, number>();

    // 计算每个任务的层级
    tasks.forEach((task) => {
      levels.set(task.id, calculateLevel(task.id, taskMap));
    });

    // 按层级分组
    const levelGroups = new Map<number, string[]>();
    levels.forEach((level, taskId) => {
      if (!levelGroups.has(level)) {
        levelGroups.set(level, []);
      }
      levelGroups.get(level)!.push(taskId);
    });

    // 生成节点位置（横向布局：从左到右）
    const initialNodes: Node[] = tasks.map((task) => {
      const level = levels.get(task.id) || 0;
      const tasksInLevel = levelGroups.get(level) || [];
      const indexInLevel = tasksInLevel.indexOf(task.id);
      const totalInLevel = tasksInLevel.length;

      // 垂直居中分布同层级的节点
      const yOffset = (indexInLevel - (totalInLevel - 1) / 2) * 180;

      return {
        id: task.id,
        type: "workflowNode",
        position: { x: level * 320, y: 150 + yOffset }, // x 代表层级，y 代表同层中的位置
        data: {
          ...task,
          isFirst: level === 0,
          isLast: false, // DAG 中不再有明确的 last 概念
          isEditMode: false,
          onMoveUp: () => handleMoveUp(task.id),
          onMoveDown: () => handleMoveDown(task.id),
        },
      };
    });

    // 根据 dependencies 生成边
    const initialEdges: Edge[] = [];
    tasks.forEach((task) => {
      if (task.dependencies && task.dependencies.length > 0) {
        task.dependencies.forEach((depId) => {
          const sourceTask = taskMap.get(depId);
          if (sourceTask) {
            initialEdges.push({
              id: `e${depId}-${task.id}`,
              source: depId,
              target: task.id,
              type: "smoothstep",
              animated: sourceTask.status === "running" || task.status === "running",
              style: {
                stroke: sourceTask.status === "success" ? "#52c41a" : sourceTask.status === "failed" ? "#ff4d4f" : "#8c8c8c",
                strokeWidth: 2,
              },
              markerEnd: {
                type: MarkerType.ArrowClosed,
                color: sourceTask.status === "success" ? "#52c41a" : sourceTask.status === "failed" ? "#ff4d4f" : "#8c8c8c",
              },
            });
          }
        });
      }
    });

    setNodes(initialNodes);
    setEdges(initialEdges);
  }, [tasks, handleMoveUp, handleMoveDown]);

  // 组件挂载时初始化
  React.useEffect(() => {
    initializeNodesAndEdges();
  }, [initializeNodesAndEdges]);

  // 节点变化处理
  const onNodesChange = useCallback(
    (changes: NodeChange[]) => {
      console.log("节点变化:", changes, "编辑模式:", isEditMode);
      if (isEditMode) {
        setNodes((nds) => {
          const result = applyNodeChanges(changes, nds);
          console.log("更新后的节点:", result);
          return result;
        });
      }
    },
    [isEditMode]
  );

  // 边变化处理
  const onEdgesChange = useCallback(
    (changes: EdgeChange[]) => {
      console.log("边变化:", changes, "编辑模式:", isEditMode);
      if (isEditMode) {
        setEdges((eds) => {
          const result = applyEdgeChanges(changes, eds);
          console.log("更新后的边:", result);
          return result;
        });
      }
    },
    [isEditMode]
  );

  // 连接处理
  const onConnect = useCallback(
    (connection: Connection) => {
      if (isEditMode) {
        setEdges((eds) =>
          addEdge(
            {
              ...connection,
              type: "smoothstep",
              animated: false,
              style: { stroke: "#8c8c8c", strokeWidth: 2 },
              markerEnd: { type: MarkerType.ArrowClosed, color: "#8c8c8c" },
            },
            eds
          )
        );
      }
    },
    [isEditMode]
  );

  // 添加新任务节点
  const handleAddNode = useCallback(() => {
    const newId = `task-${Date.now()}`;
    setNodes((nds) => {
      const newNode: Node = {
        id: newId,
        type: "workflowNode",
        position: { x: 100, y: 150 + nds.length * 100 }, // 横向布局：新节点默认在最左侧
        data: {
          id: newId,
          name: `新任务 ${nds.length + 1}`,
          type: "custom" as const,
          status: "pending" as const,
          dependencies: [], // 新任务默认无依赖（顶点节点）
          isFirst: nds.length === 0,
          isLast: true,
          isEditMode: isEditMode,
          onMoveUp: () => handleMoveUp(newId),
          onMoveDown: () => handleMoveDown(newId),
        },
      };

      // 更新其他节点的 isLast 标记
      const updatedNodes = nds.map((node, idx) => ({
        ...node,
        data: {
          ...node.data,
          isLast: false,
        },
      }));

      return [...updatedNodes, newNode];
    });
  }, [isEditMode, handleMoveUp, handleMoveDown]);

  // 进入编辑模式
  const handleEnterEditMode = useCallback(() => {
    setIsEditMode(true);
    // 更新所有节点的 isEditMode 标记
    setNodes((nds) =>
      nds.map((node) => ({
        ...node,
        data: { ...node.data, isEditMode: true },
      }))
    );
  }, []);

  // 取消编辑
  const handleCancelEdit = useCallback(() => {
    setIsEditMode(false);
    initializeNodesAndEdges();
  }, [initializeNodesAndEdges]);

  // 保存修改
  const handleSave = useCallback(() => {
    // 从边重新构建每个节点的 dependencies
    const dependenciesMap = new Map<string, string[]>();
    edges.forEach((edge) => {
      const targetId = edge.target;
      if (!dependenciesMap.has(targetId)) {
        dependenciesMap.set(targetId, []);
      }
      dependenciesMap.get(targetId)!.push(edge.source);
    });

    // 将节点和边转换回 Task 格式
    const updatedTasks: Task[] = nodes.map((node) => ({
      id: node.id,
      name: node.data.name,
      type: node.data.type,
      status: node.data.status,
      dependencies: dependenciesMap.get(node.id) || [], // 从边重建依赖关系
      duration: node.data.duration,
      logs: node.data.logs,
      params: node.data.params, // 保存任务参数
      deploymentId: node.data.deploymentId,
      appId: node.data.appId,
      blockBy: node.data.blockBy,
      startedAt: node.data.startedAt,
      completedAt: node.data.completedAt,
    }));

    if (onSave) {
      onSave(updatedTasks);
    }
    setIsEditMode(false);
  }, [nodes, edges, onSave]);

  return (
    <div style={{ height: "100%", display: "flex", flexDirection: "column" }}>
      {/* 工具栏 */}
      {allowEdit && (
        <div className="d-flex gap-2 mb-2" style={{ flexShrink: 0 }}>
          {!isEditMode ? (
            <button className="btn btn-primary btn-sm" onClick={handleEnterEditMode}>
              <IconEdit size={16} className="me-1" />
              编辑工作流
            </button>
          ) : (
            <>
              <button className="btn btn-success btn-sm" onClick={handleSave}>
                <IconDeviceFloppy size={16} className="me-1" />
                保存
              </button>
              <button className="btn btn-secondary btn-sm" onClick={handleCancelEdit}>
                <IconX size={16} className="me-1" />
                取消
              </button>
              <button className="btn btn-primary btn-sm" onClick={handleAddNode}>
                <IconPlus size={16} className="me-1" />
                添加任务
              </button>
              <span className="badge bg-info-lt align-self-center ms-auto" style={{ fontSize: "0.75rem" }}>
                <IconEdit size={14} className="me-1" />
                可拖拽节点、上移/下移调整顺序、创建连接、删除元素（Del/Backspace）
              </span>
            </>
          )}
        </div>
      )}

      <div style={{ flex: 1, border: "1px solid #e5e7eb", borderRadius: "8px", overflow: "hidden" }}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          nodeTypes={nodeTypes}
          fitView
          fitViewOptions={{ padding: 0.2 }}
          nodesDraggable={isEditMode}
          nodesConnectable={isEditMode}
          elementsSelectable={true}
          panOnScroll={true}
          zoomOnScroll={true}
          minZoom={0.5}
          maxZoom={1.5}
          deleteKeyCode={isEditMode ? ["Backspace", "Delete"] : []}
        >
          <Background variant={BackgroundVariant.Dots} gap={16} size={1} color="#e5e7eb" />
          <Controls showInteractive={false} />
        </ReactFlow>
      </div>
    </div>
  );
};

export default WorkflowViewer;
