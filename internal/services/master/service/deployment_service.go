package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/pkg/utils"
	"github.com/google/uuid"
)

type deploymentService struct {
	deploymentRepo interfaces.DeploymentRepository
	versionRepo    interfaces.VersionRepository
	appRepo        interfaces.ApplicationRepository
	envRepo        interfaces.EnvironmentRepository

	workflow *workflowController
}

// NewDeploymentService 创建部署服务
func NewDeploymentService(
	deploymentRepo interfaces.DeploymentRepository,
	versionRepo interfaces.VersionRepository,
	appRepo interfaces.ApplicationRepository,
	envRepo interfaces.EnvironmentRepository,
	workflow *workflowController,
) interfaces.DeploymentService {
	return &deploymentService{
		deploymentRepo: deploymentRepo,
		versionRepo:    versionRepo,
		appRepo:        appRepo,
		envRepo:        envRepo,
		workflow:       workflow,
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
	_, err = s.envRepo.GetByID(ctx, req.EnvironmentID)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	// 验证应用是否存在
	for _, appID := range req.MustInOrder {
		// 验证应用是否在版本中
		appFound := false
		for _, app := range version.GetAppBuilds() {
			if app.AppID == appID {
				appFound = true
				break
			}
		}
		if !appFound {
			return nil, fmt.Errorf("application %s not found in version %s", appID, version.ID)
		}
	}

	if req.ManualApproval {
		defaultStat := false
		for i := range req.Strategy {
			req.Strategy[i].ManualApprovalStatus = &defaultStat
		}
	}

	// 创建部署
	deployment := &models.Deployment{
		ID:            uuid.New().String(),
		VersionID:     req.VersionID,
		EnvironmentID: req.EnvironmentID,
		Status:        models.DeploymentStatusPending,
		CreatedBy:     getCurrentUser(ctx),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if len(req.MustInOrder) > 0 {
		sbs, _ := json.Marshal(req.MustInOrder)
		deployment.MustInOrder = sbs
	}
	if len(req.Strategy) > 0 {
		sbs, _ := json.Marshal(req.Strategy)
		deployment.Strategy = sbs
	}

	if err := s.deploymentRepo.Create(ctx, deployment); err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	// TODO: 创建并执行部署任务

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

func (s *deploymentService) StartDeployment(ctx context.Context, id string) (*models.Deployment, error) {
	// 获取部署
	deployment, err := s.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("deployment not found: %w", err)
	}

	if err := s.workflow.CreateTasksFromDeployment(ctx, deployment); err != nil {
		return nil, fmt.Errorf("failed to create tasks from deployment: %w", err)
	}

	deployment.Status = models.DeploymentStatusRunning
	deployment.StartedAt = &[]time.Time{time.Now()}[0]
	deployment.UpdatedAt = time.Now()

	if err := s.deploymentRepo.Update(ctx, deployment); err != nil {
		return nil, fmt.Errorf("failed to update deployment: %w", err)
	}

	return deployment, nil
}

func (s *deploymentService) RollbackDeployment(ctx context.Context, id string, req *models.RollbackRequest) error {

	currentDeployment, err := s.deploymentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("deployment not found: %w", err)
	}

	currentDeployment.Status = models.DeploymentStatusRolledBack
	currentDeployment.UpdatedAt = time.Now()

	if err := s.deploymentRepo.Update(ctx, currentDeployment); err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	return nil
}
