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
	builds := version.GetAppBuilds()

	if len(deployment.MustInOrder) == 0 {
		for _, appBuild := range version.GetAppBuilds() {
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
		mustInOrder := deployment.GetMustInOrder()
		for i, appID := range mustInOrder {
			var appBuild *models.AppBuild
			for _, ab := range builds {
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
				prevTask, err := wc.findTaskByAppID(ctx, deployment.ID, mustInOrder[i-1])
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
		if task.Deployment.Status == models.DeploymentStatusPending {
			continue
		}
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

	var block bool
	prevVersion, err := wc.versionRepo.GetPreviousByVersionAndApp(
		ctx, task.Deployment.VersionID, task.AppID)
	if err == nil {
		deployments, _, err_ := wc.deploymentRepo.List(ctx, &models.DeploymentFilter{
			VersionID: prevVersion.ID,
		})
		if err_ != nil {
			wc.log.Warn("Failed to list deployments",
				zap.Error(err_), zap.String("version_id", prevVersion.ID))
		}
		for _, deployment := range deployments {
			if deployment.Status != models.DeploymentStatusRunning {
				continue
			}
			for _, pt := range deployment.Tasks {
				if pt.AppID == task.AppID {
					task.Status = models.TaskStatusBlocked
					task.BlockBy = pt.ID
					block = true
					break
				}
			}
			if block {
				break
			}
		}
	}

	if block {
		wc.log.Info("Task blocked",
			zap.String("task_id", task.ID),
			zap.String("blocking_task_id", task.BlockBy))
		return wc.taskRepo.Update(ctx, task)
	}

	executor := NewSimpleDeployExecutor(*task, wc.deploymentRepo, wc.versionRepo)
	status, err := executor.Apply(ctx)
	task.Status = status
	task.UpdatedAt = time.Now()
	if task.Status.IsFinished() {
		task.CompletedAt = &[]time.Time{time.Now()}[0]
	}

	if err != nil {
		wc.log.Error("Task execution failed", zap.Error(err), zap.String("task_id", task.ID))
		task.Result = err.Error()
	}

	return wc.taskRepo.Update(ctx, task)
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
	elapsed := time.Since(task.UpdatedAt)
	if elapsed > wc.config.TaskTimeout {
		wc.log.Warn("Task heartbeat timeout detected, rescheduling",
			zap.String("task_id", task.ID),
			zap.String("elapsed", elapsed.String()),
			zap.String("timeout", wc.config.TaskTimeout.String()))

		task.Status = models.TaskStatusPending
		task.StartedAt = nil
		task.UpdatedAt = time.Now()
		task.Result = fmt.Sprintf("Task heartbeat timeout after %s, rescheduling", elapsed.String())

		if err := wc.taskRepo.Update(ctx, task); err != nil {
			return fmt.Errorf("failed to reschedule timed out task: %w", err)
		}
	}

	return nil
}
