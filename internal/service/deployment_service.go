package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guaguasong/1024Hackathon/internal/interfaces"
	"github.com/guaguasong/1024Hackathon/internal/models"
	"github.com/guaguasong/1024Hackathon/internal/utils"
)

type deploymentService struct {
	deploymentRepo interfaces.DeploymentRepository
	versionRepo    interfaces.VersionRepository
	appRepo        interfaces.ApplicationRepository
	envRepo        interfaces.EnvironmentRepository
	workflowMgr    interfaces.WorkflowManager
}

// NewDeploymentService 创建部署服务
func NewDeploymentService(
	deploymentRepo interfaces.DeploymentRepository,
	versionRepo interfaces.VersionRepository,
	appRepo interfaces.ApplicationRepository,
	envRepo interfaces.EnvironmentRepository,
	workflowMgr interfaces.WorkflowManager,
) interfaces.DeploymentService {
	return &deploymentService{
		deploymentRepo: deploymentRepo,
		versionRepo:    versionRepo,
		appRepo:        appRepo,
		envRepo:        envRepo,
		workflowMgr:    workflowMgr,
	}
}

func (s *deploymentService) CreateDeployment(ctx context.Context, req *models.CreateDeploymentRequest) (*models.Deployment, error) {
	// 验证请求
	if err := utils.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 验证版本是否存在
	version, err := s.versionRepo.GetByID(ctx, req.VersionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	// 验证环境是否存在
	environment, err := s.envRepo.GetByID(ctx, req.EnvironmentID)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	// 验证应用是否存在
	for _, appID := range req.ApplicationIDs {
		_, err := s.appRepo.GetByID(ctx, appID)
		if err != nil {
			return nil, fmt.Errorf("application %s not found: %w", appID, err)
		}
	}

	// 创建部署
	deployment := &models.Deployment{
		ID:             uuid.New().String(),
		VersionID:      req.VersionID,
		ApplicationIDs: req.ApplicationIDs,
		EnvironmentID:  req.EnvironmentID,
		Status:         models.DeploymentStatusPending,
		CreatedBy:      getCurrentUser(ctx),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.deploymentRepo.Create(ctx, deployment); err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	// 创建工作流
	workflow, err := s.workflowMgr.CreateWorkflow(ctx, deployment)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	// 异步执行工作流
	go func() {
		ctx := context.Background()
		if err := s.workflowMgr.ExecuteWorkflow(ctx, workflow.ID); err != nil {
			// 更新部署状态为失败
			s.deploymentRepo.Update(ctx, &models.Deployment{
				ID:           deployment.ID,
				Status:       models.DeploymentStatusFailed,
				ErrorMessage: err.Error(),
				CompletedAt:  &[]time.Time{time.Now()}[0],
			})
		}
	}()

	return deployment, nil
}

func (s *deploymentService) GetDeploymentList(ctx context.Context, req *models.ListDeploymentsRequest) (*models.DeploymentListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	filter := &models.DeploymentFilter{
		Status:        models.DeploymentStatus(req.Status),
		EnvironmentID: req.EnvironmentID,
		VersionID:     req.VersionID,
		Page:          req.Page,
		PageSize:      req.PageSize,
	}

	deployments, total, err := s.deploymentRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	return &models.DeploymentListResponse{
		Deployments: deployments,
		Total:       total,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}, nil
}

func (s *deploymentService) GetDeployment(ctx context.Context, id string) (*models.Deployment, error) {
	deployment, err := s.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return deployment, nil
}

func (s *deploymentService) CancelDeployment(ctx context.Context, id string) (*models.Deployment, error) {
	// 获取部署
	deployment, err := s.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("deployment not found: %w", err)
	}

	// 检查状态是否可以取消
	if deployment.Status != models.DeploymentStatusPending && deployment.Status != models.DeploymentStatusRunning {
		return nil, fmt.Errorf("deployment cannot be cancelled in status %s", deployment.Status)
	}

	// 取消工作流
	if err := s.workflowMgr.CancelWorkflow(ctx, id); err != nil {
		return nil, fmt.Errorf("failed to cancel workflow: %w", err)
	}

	// 更新部署状态
	deployment.Status = models.DeploymentStatusFailed
	deployment.ErrorMessage = "Deployment cancelled by user"
	deployment.CompletedAt = &[]time.Time{time.Now()}[0]

	if err := s.deploymentRepo.Update(ctx, deployment); err != nil {
		return nil, fmt.Errorf("failed to update deployment: %w", err)
	}

	return deployment, nil
}

func (s *deploymentService) RollbackDeployment(ctx context.Context, id string, req *models.RollbackRequest) (*models.Deployment, error) {
	// 验证请求
	if err := utils.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取当前部署
	currentDeployment, err := s.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("deployment not found: %w", err)
	}

	// 验证目标版本是否存在
	_, err = s.versionRepo.GetByID(ctx, req.TargetVersionID)
	if err != nil {
		return nil, fmt.Errorf("target version not found: %w", err)
	}

	// 创建回滚部署
	rollbackDeployment := &models.Deployment{
		ID:             uuid.New().String(),
		VersionID:      req.TargetVersionID,
		ApplicationIDs: currentDeployment.ApplicationIDs,
		EnvironmentID:  currentDeployment.EnvironmentID,
		Status:         models.DeploymentStatusPending,
		CreatedBy:      getCurrentUser(ctx),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.deploymentRepo.Create(ctx, rollbackDeployment); err != nil {
		return nil, fmt.Errorf("failed to create rollback deployment: %w", err)
	}

	// 创建工作流
	workflow, err := s.workflowMgr.CreateWorkflow(ctx, rollbackDeployment)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	// 异步执行工作流
	go func() {
		ctx := context.Background()
		if err := s.workflowMgr.ExecuteWorkflow(ctx, workflow.ID); err != nil {
			// 更新部署状态为失败
			s.deploymentRepo.Update(ctx, &models.Deployment{
				ID:           rollbackDeployment.ID,
				Status:       models.DeploymentStatusFailed,
				ErrorMessage: err.Error(),
				CompletedAt:  &[]time.Time{time.Now()}[0],
			})
		}
	}()

	return rollbackDeployment, nil
}
