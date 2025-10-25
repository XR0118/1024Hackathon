package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WorkflowConfig struct {
	PendingCheckInterval time.Duration
	BlockedCheckInterval time.Duration
	RunningCheckInterval time.Duration
	TaskTimeout          time.Duration
}

type workflowController struct {
	taskService    interfaces.TaskService
	taskRepo       interfaces.TaskRepository
	deploymentRepo interfaces.DeploymentRepository
	versionRepo    interfaces.VersionRepository
	config         WorkflowConfig
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	log            *zap.Logger
}

func NewWorkflowController(
	taskService interfaces.TaskService,
	taskRepo interfaces.TaskRepository,
	deploymentRepo interfaces.DeploymentRepository,
	versionRepo interfaces.VersionRepository,
	config WorkflowConfig,
	log *zap.Logger,
) *workflowController {
	if config.PendingCheckInterval == 0 {
		config.PendingCheckInterval = 5 * time.Second
	}
	if config.BlockedCheckInterval == 0 {
		config.BlockedCheckInterval = 10 * time.Second
	}
	if config.RunningCheckInterval == 0 {
		config.RunningCheckInterval = 30 * time.Second
	}
	if config.TaskTimeout == 0 {
		config.TaskTimeout = 30 * time.Minute
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &workflowController{
		taskService:    taskService,
		taskRepo:       taskRepo,
		deploymentRepo: deploymentRepo,
		versionRepo:    versionRepo,
		config:         config,
		ctx:            ctx,
		cancel:         cancel,
		log:            log,
	}
}

func (wc *workflowController) Start() {
	wc.log.Info("Starting workflow controller")

	wc.wg.Add(3)
	go wc.pendingTaskScheduler()
	go wc.blockedTaskScheduler()
	go wc.runningTaskScheduler()
}

func (wc *workflowController) Stop() {
	wc.log.Info("Stopping workflow controller")
	wc.cancel()
	wc.wg.Wait()
	wc.log.Info("Workflow controller stopped")
}

func (wc *workflowController) CreateTasksFromDeployment(ctx context.Context, deployment *models.Deployment) error {
	version, err := wc.versionRepo.GetByID(ctx, deployment.VersionID)
	if err != nil {
		return fmt.Errorf("failed to get version: %w", err)
	}

	if len(deployment.MustInOrder) == 0 {
		for _, appBuild := range version.AppBuilds {
			task := &models.Task{
				ID:           uuid.New().String(),
				DeploymentID: deployment.ID,
				AppID:        appBuild.AppID,
				Type:         "deploy",
				Status:       models.TaskStatusPending,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			if err := wc.taskRepo.Create(ctx, task); err != nil {
				return fmt.Errorf("failed to create task for app %s: %w", appBuild.AppID, err)
			}
		}
	} else {
		for i, appID := range deployment.MustInOrder {
			var appBuild *models.AppBuild
			for _, ab := range version.AppBuilds {
				if ab.AppID == appID {
					appBuild = &ab
					break
				}
			}
			if appBuild == nil {
				return fmt.Errorf("app %s not found in version %s", appID, version.ID)
			}

			var status models.TaskStatus
			var blockBy string
			if i == 0 {
				status = models.TaskStatusPending
			} else {
				status = models.TaskStatusBlocked
				prevTask, err := wc.findTaskByAppID(ctx, deployment.ID, deployment.MustInOrder[i-1])
				if err != nil {
					return fmt.Errorf("failed to find previous task: %w", err)
				}
				if prevTask != nil {
					blockBy = prevTask.ID
				}
			}

			task := &models.Task{
				ID:           uuid.New().String(),
				DeploymentID: deployment.ID,
				AppID:        appID,
				Type:         "deploy",
				Status:       status,
				BlockBy:      blockBy,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			if err := wc.taskRepo.Create(ctx, task); err != nil {
				return fmt.Errorf("failed to create task for app %s: %w", appID, err)
			}
		}
	}

	return nil
}

func (wc *workflowController) findTaskByAppID(ctx context.Context, deploymentID, appID string) (*models.Task, error) {
	tasks, err := wc.taskRepo.GetByDeploymentID(ctx, deploymentID)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		if task.AppID == appID {
			return task, nil
		}
	}
	return nil, nil
}

func (wc *workflowController) pendingTaskScheduler() {
	defer wc.wg.Done()
	ticker := time.NewTicker(wc.config.PendingCheckInterval)
	defer ticker.Stop()

	wc.log.Info("Pending task scheduler started")

	for {
		select {
		case <-wc.ctx.Done():
			wc.log.Info("Pending task scheduler stopped")
			return
		case <-ticker.C:
			wc.processPendingTasks()
		}
	}
}

func (wc *workflowController) processPendingTasks() {
	ctx := context.Background()

	filter := &models.TaskFilter{
		Status:   models.TaskStatusPending,
		Page:     1,
		PageSize: 100,
	}

	tasks, _, err := wc.taskRepo.List(ctx, filter)
	if err != nil {
		wc.log.Error("Failed to list pending tasks", zap.Error(err))
		return
	}

	for _, task := range tasks {
		if err := wc.executeTask(ctx, task); err != nil {
			wc.log.Error("Failed to execute task", zap.Error(err), zap.String("task_id", task.ID))
		}
	}
}

func (wc *workflowController) executeTask(ctx context.Context, task *models.Task) error {
	task.Status = models.TaskStatusRunning
	now := time.Now()
	task.StartedAt = &now
	task.UpdatedAt = now

	if err := wc.taskRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to update task status to running: %w", err)
	}

	return nil
}

func (wc *workflowController) blockedTaskScheduler() {
	defer wc.wg.Done()
	ticker := time.NewTicker(wc.config.BlockedCheckInterval)
	defer ticker.Stop()

	wc.log.Info("Blocked task scheduler started")

	for {
		select {
		case <-wc.ctx.Done():
			wc.log.Info("Blocked task scheduler stopped")
			return
		case <-ticker.C:
			wc.processBlockedTasks()
		}
	}
}

func (wc *workflowController) processBlockedTasks() {
	ctx := context.Background()

	filter := &models.TaskFilter{
		Status:   models.TaskStatusBlocked,
		Page:     1,
		PageSize: 100,
	}

	tasks, _, err := wc.taskRepo.List(ctx, filter)
	if err != nil {
		wc.log.Error("Failed to list blocked tasks", zap.Error(err))
		return
	}

	for _, task := range tasks {
		if err := wc.checkAndUnblockTask(ctx, task); err != nil {
			wc.log.Error("Failed to check blocked task", zap.Error(err), zap.String("task_id", task.ID))
		}
	}
}

func (wc *workflowController) checkAndUnblockTask(ctx context.Context, task *models.Task) error {
	if task.BlockBy == "" {
		return nil
	}

	blockingTask, err := wc.taskRepo.GetByID(ctx, task.BlockBy)
	if err != nil {
		return fmt.Errorf("failed to get blocking task: %w", err)
	}

	if blockingTask.Status == models.TaskStatusSuccess {
		task.Status = models.TaskStatusPending
		task.BlockBy = ""
		task.UpdatedAt = time.Now()

		if err := wc.taskRepo.Update(ctx, task); err != nil {
			return fmt.Errorf("failed to unblock task: %w", err)
		}
		wc.log.Info("Task unblocked",
			zap.String("task_id", task.ID),
			zap.String("blocking_task_id", blockingTask.ID))
	}

	return nil
}

func (wc *workflowController) runningTaskScheduler() {
	defer wc.wg.Done()
	ticker := time.NewTicker(wc.config.RunningCheckInterval)
	defer ticker.Stop()

	wc.log.Info("Running task scheduler started")

	for {
		select {
		case <-wc.ctx.Done():
			wc.log.Info("Running task scheduler stopped")
			return
		case <-ticker.C:
			wc.processRunningTasks()
		}
	}
}

func (wc *workflowController) processRunningTasks() {
	ctx := context.Background()

	filter := &models.TaskFilter{
		Status:   models.TaskStatusRunning,
		Page:     1,
		PageSize: 100,
	}

	tasks, _, err := wc.taskRepo.List(ctx, filter)
	if err != nil {
		wc.log.Error("Failed to list running tasks", zap.Error(err))
		return
	}

	for _, task := range tasks {
		if err := wc.checkTaskTimeout(ctx, task); err != nil {
			wc.log.Error("Failed to check task timeout", zap.Error(err), zap.String("task_id", task.ID))
		}
	}
}

func (wc *workflowController) checkTaskTimeout(ctx context.Context, task *models.Task) error {
	if task.StartedAt == nil {
		return nil
	}

	elapsed := time.Since(*task.StartedAt)
	if elapsed > wc.config.TaskTimeout {
		wc.log.Warn("Task timeout detected, rescheduling",
			zap.String("task_id", task.ID),
			zap.String("elapsed", elapsed.String()),
			zap.String("timeout", wc.config.TaskTimeout.String()))

		task.Status = models.TaskStatusPending
		task.StartedAt = nil
		task.UpdatedAt = time.Now()
		task.Result = fmt.Sprintf("Task timed out after %s, rescheduling", elapsed.String())

		if err := wc.taskRepo.Update(ctx, task); err != nil {
			return fmt.Errorf("failed to reschedule timed out task: %w", err)
		}
	}

	return nil
}
