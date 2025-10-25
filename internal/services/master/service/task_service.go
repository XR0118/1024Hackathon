package service

import (
	"context"
	"fmt"

	"github.com/XR0118/1024Hackathon/internal/interfaces"
	"github.com/XR0118/1024Hackathon/internal/pkg/models"
)

type taskService struct {
	taskRepo interfaces.TaskRepository
}

// NewTaskService 创建任务服务
func NewTaskService(taskRepo interfaces.TaskRepository) interfaces.TaskService {
	return &taskService{
		taskRepo: taskRepo,
	}
}

func (s *taskService) GetTaskList(ctx context.Context, req *models.ListTasksRequest) (*models.TaskListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	filter := &models.TaskFilter{
		DeploymentID: req.DeploymentID,
		Status:       models.TaskStatus(req.Status),
		Type:         req.Type,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}

	tasks, total, err := s.taskRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return &models.TaskListResponse{
		Tasks:    tasks,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *taskService) GetTask(ctx context.Context, id string) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return task, nil
}

func (s *taskService) RetryTask(ctx context.Context, id string) (*models.Task, error) {
	// 获取任务
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// 检查任务状态是否可以重试
	if task.Status != models.TaskStatusFailed {
		return nil, fmt.Errorf("task cannot be retried in status %s", task.Status)
	}

	// 重置任务状态
	task.Status = models.TaskStatusPending
	task.Result = ""
	task.StartedAt = nil
	task.CompletedAt = nil

	// 更新任务
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// TODO: 重新调度任务执行
	// 这里应该通知任务调度器重新执行任务

	return task, nil
}
