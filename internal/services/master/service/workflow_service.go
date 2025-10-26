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

	wc.wg.Add(2)
	go wc.pendingTaskScheduler()
	go wc.blockedTaskScheduler()
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

	// 创建构建任务
	buildID := uuid.New().String()
	buildTask := &models.Task{
		ID:           buildID,
		DeploymentID: deployment.ID,
		AppID:        builds[0].AppID,
		Type:         models.TaskTypeBuild,
		Step:         models.TaskStepPending,
		Status:       models.TaskStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := wc.taskRepo.Create(ctx, buildTask); err != nil {
		return fmt.Errorf("failed to create build task for deployment: %w", err)
	}

	// 创建审核任务
	// approvalID := uuid.New().String()
	// approvalTask := &models.Task{
	// 	ID:           approvalID,
	// 	DeploymentID: deployment.ID,
	// 	Type:         models.TaskTypeApproval,
	// 	Step:         models.TaskStepBlocked,
	// 	Status:       models.TaskStatusPending,
	// 	CreatedAt:    time.Now(),
	// 	UpdatedAt:    time.Now(),
	// }
	// _ = approvalTask.SetDependencies([]string{buildID})
	// if err := wc.taskRepo.Create(ctx, approvalTask); err != nil {
	// 	return fmt.Errorf("failed to create approval task for deployment: %w", err)
	// }

	if len(deployment.MustInOrder) == 0 {
		for _, appBuild := range version.GetAppBuilds() {
			task := &models.Task{
				ID:           uuid.New().String(),
				DeploymentID: deployment.ID,
				AppID:        appBuild.AppID,
				Type:         models.TaskTypeDeploy,
				Step:         models.TaskStepBlocked,
				Status:       models.TaskStatusPending,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			_ = task.SetDependencies([]string{buildID})
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

			// 所有任务初始状态都是 pending
			// 是否被阻塞通过检查 dependencies 判断，不需要专门的 blocked 状态
			var blockBy string
			if i > 0 {
				prevTask, err := wc.findTaskByAppID(ctx, deployment.ID, mustInOrder[i-1])
				if err != nil {
					return fmt.Errorf("failed to find previous task: %w", err)
				}
				if prevTask != nil {
					blockBy = prevTask.ID
				}
			} else {
				blockBy = buildID
			}

			// 如果有依赖，初始状态为 blocked；否则为 pending
			step := models.TaskStepPending
			if blockBy != "" {
				step = models.TaskStepBlocked
			}

			task := &models.Task{
				ID:           uuid.New().String(),
				DeploymentID: deployment.ID,
				AppID:        appID,
				Name:         "Deploy " + appID,
				Type:         models.TaskTypeDeploy,
				Step:         step,
				Status:       models.TaskStatusPending,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			// 设置依赖关系
			if blockBy != "" {
				_ = task.SetDependencies([]string{blockBy})
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

	// 查询 Step = pending 且 Status = pending 的任务（准备执行的任务）
	var tasks []*models.Task
	var steps = []models.TaskStep{
		models.TaskStepPending,
		models.TaskStepRunning,
	}
	for _, step := range steps {
		status := models.TaskStatusPending
		if step == models.TaskStepRunning {
			status = models.TaskStatusRunning
		}
		filter := &models.TaskFilter{
			Step:     step,
			Status:   status,
			Page:     1,
			PageSize: 100,
		}

		ts, _, err := wc.taskRepo.List(ctx, filter)
		if err != nil {
			wc.log.Error("Failed to list pending tasks", zap.Error(err))
			return
		}
		tasks = append(tasks, ts...)
	}

	for _, task := range tasks {
		// 跳过部署尚未开始的任务
		if task.Deployment.Status == models.DeploymentStatusPending {
			continue
		}
		if err := wc.executeTask(ctx, task); err != nil {
			wc.log.Error("Failed to execute task", zap.Error(err), zap.String("task_id", task.ID))
		}
	}
}

func (wc *workflowController) executeTask(ctx context.Context, task *models.Task) error {
	// 首先检查依赖是否都已完成
	deps := task.GetDependencies()
	if len(deps) > 0 {
		for _, depID := range deps {
			depTask, err := wc.taskRepo.GetByID(ctx, depID)
			if err != nil {
				return fmt.Errorf("failed to get dependency task: %w", err)
			}
			if depTask.Status != models.TaskStatusSuccess {
				// 依赖未完成，任务设置为 blocked
				if task.Step != models.TaskStepBlocked {
					task.Step = models.TaskStepBlocked
					task.UpdatedAt = time.Now()
					_ = wc.taskRepo.Update(ctx, task)
				}
				wc.log.Info("Task dependencies not satisfied, blocked",
					zap.String("task_id", task.ID),
					zap.String("waiting_for", depID))
				return nil
			}
		}
	}

	// 依赖都完成了，可以执行任务
	task.Step = models.TaskStepRunning
	task.Status = models.TaskStatusRunning
	now := time.Now()
	task.StartedAt = &now
	task.UpdatedAt = now

	// 检查前一个版本的同 app 任务是否还在执行
	// 如果在执行，则等待其完成（可选逻辑）
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
				if pt.Type != models.TaskTypeDeploy {
					continue
				}
				// 如果前一个版本的同 app 任务已被阻塞，记录日志
				if pt.AppID == task.AppID && pt.Step == models.TaskStepBlocked {
					wc.log.Info("Previous version task still running",
						zap.String("task_id", task.ID),
						zap.String("prev_task_id", pt.ID),
						zap.String("prev_version", prevVersion.ID))
					task.Step = models.TaskStepBlocked
					task.UpdatedAt = time.Now()
					return wc.taskRepo.Update(ctx, task)
				}
			}
		}
	}

	wc.log.Info("Starting task execution",
		zap.String("task_id", task.ID),
		zap.Any("dependencies", task.GetDependencies()))

	// 更新任务状态为 running
	if err := wc.taskRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to update task status to running: %w", err)
	}

	switch task.Type {
	case models.TaskTypeBuild:
		// 模拟构建过程
		time.Sleep(10 * time.Second)

		task.Status = models.TaskStatusSuccess
		now := time.Now()
		task.UpdatedAt = now
		task.CompletedAt = &[]time.Time{now}[0]
		task.Step = models.TaskStepCompleted
		return wc.taskRepo.Update(ctx, task)

	case models.TaskTypeApproval:
		var result models.ApprovalTaskResult
		err := task.GetResult(&result)
		if err != nil {
			wc.log.Warn("Invalid ApprovalTaskResult:", zap.Error(err), zap.String("payload", string(task.Payload)))
			return err
		}
		if result.Approved {
			task.Status = models.TaskStatusSuccess
			now := time.Now()
			task.UpdatedAt = now
			task.CompletedAt = &now
			return wc.taskRepo.Update(ctx, task)
		}
		return nil
	case models.TaskTypeDeploy:
	default:
		return fmt.Errorf("unsupported task type: %s", task.Type)
	}

	executor := NewSimpleDeployExecutor(*task, wc.deploymentRepo, wc.versionRepo)
	status, err := executor.Apply(ctx)
	task.Status = status
	task.UpdatedAt = time.Now()
	if task.Status.IsFinished() {
		task.Step = models.TaskStepCompleted
		task.CompletedAt = &[]time.Time{time.Now()}[0]
	}

	if err != nil {
		wc.log.Error("Task execution failed", zap.Error(err), zap.String("task_id", task.ID))
		// 将错误信息存入 Result
		_ = task.SetResult(map[string]interface{}{
			"error": err.Error(),
		})
	}

	err = wc.taskRepo.Update(ctx, task)
	if err != nil {
		wc.log.Error("Failed to update task", zap.Error(err), zap.String("task_id", task.ID))
		return err
	}

	if task.Status.IsFinished() {
		deployment, err := wc.deploymentRepo.GetByID(ctx, task.DeploymentID)
		if err != nil {
			wc.log.Error("Failed to get deployment", zap.Error(err), zap.String("deployment_id", task.DeploymentID))
		}
		if !deployment.Status.IsFinished() {
			var deploymentStatus models.DeploymentStatus
			var successCnt int
			for _, dt := range deployment.Tasks {
				var done bool
				switch dt.Status {
				case models.TaskStatusFailed:
					deploymentStatus = models.DeploymentStatusCompleted
					done = true
				case models.TaskStatusSuccess:
					successCnt++
				default:
					// pending 或 running
					deploymentStatus = models.DeploymentStatusRunning
					done = true
				}
				if done {
					break
				}
			}
			if successCnt == len(deployment.Tasks) {
				deploymentStatus = models.DeploymentStatusCompleted
			}
			if deploymentStatus.IsFinished() {
				wc.log.Info("Deployment completed",
					zap.String("deployment_id", deployment.ID),
					zap.String("status", string(deploymentStatus)))

				deployment.Status = deploymentStatus
				deployment.CompletedAt = &[]time.Time{time.Now()}[0]
				if err := wc.deploymentRepo.Update(ctx, deployment); err != nil {
					wc.log.Error("Failed to update deployment", zap.Error(err), zap.String("deployment_id", deployment.ID))
					return err
				}
			}
		}
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

// processBlockedTasks 处理 blocked 状态的任务，检查依赖是否满足
func (wc *workflowController) processBlockedTasks() {
	ctx := context.Background()

	// 查询所有 Step = blocked 的任务
	filter := &models.TaskFilter{
		Step:     models.TaskStepBlocked,
		Page:     1,
		PageSize: 100,
	}

	tasks, _, err := wc.taskRepo.List(ctx, filter)
	if err != nil {
		wc.log.Error("Failed to list blocked tasks", zap.Error(err))
		return
	}

	for _, task := range tasks {
		// 检查依赖是否满足，如果满足则可以将 Step 改为 pending，让 processPendingTasks 去执行
		if err := wc.checkAndUnblockTask(ctx, task); err != nil {
			wc.log.Error("Failed to check blocked task", zap.Error(err), zap.String("task_id", task.ID))
		}
	}
}

// checkAndUnblockTask 检查任务的依赖是否都已完成
// 如果依赖都完成了，任务就可以被执行（由 executeTask 处理）
func (wc *workflowController) checkAndUnblockTask(ctx context.Context, task *models.Task) error {
	deps := task.GetDependencies()

	// 检查所有依赖任务是否都已完成
	allCompleted := true
	for _, depID := range deps {
		depTask, err := wc.taskRepo.GetByID(ctx, depID)
		if err != nil {
			return fmt.Errorf("failed to get dependency task: %w", err)
		}
		// 依赖任务必须成功完成
		if depTask.Status != models.TaskStatusSuccess {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		// 依赖都完成了，将 Step 改为 pending，让 processPendingTasks 去执行
		// 注意：不清空 dependencies，保留依赖信息用于审计和追踪
		task.Step = models.TaskStepPending
		task.UpdatedAt = time.Now()
		if err := wc.taskRepo.Update(ctx, task); err != nil {
			return fmt.Errorf("failed to unblock task: %w", err)
		}
		wc.log.Info("Task unblocked, dependencies satisfied",
			zap.String("task_id", task.ID),
			zap.Any("dependencies", deps))
	}

	return nil
}
