package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/XR0118/1024Hackathon/internal/interfaces"
	"github.com/XR0118/1024Hackathon/internal/pkg/models"
	"github.com/XR0118/1024Hackathon/internal/pkg/utils"
)

type environmentService struct {
	envRepo interfaces.EnvironmentRepository
}

// NewEnvironmentService 创建环境服务
func NewEnvironmentService(envRepo interfaces.EnvironmentRepository) interfaces.EnvironmentService {
	return &environmentService{
		envRepo: envRepo,
	}
}

func (s *environmentService) CreateEnvironment(ctx context.Context, req *models.CreateEnvironmentRequest) (*models.Environment, error) {
	// 验证请求
	if err := utils.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 检查环境名称是否已存在
	existingEnvs, _, err := s.envRepo.List(ctx, &models.EnvironmentFilter{
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check existing environments: %w", err)
	}

	for _, env := range existingEnvs {
		if env.Name == req.Name {
			return nil, fmt.Errorf("environment with name %s already exists", req.Name)
		}
	}

	// 创建环境
	env := &models.Environment{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Type:      req.Type,
		Config:    req.Config,
		IsActive:  req.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.envRepo.Create(ctx, env); err != nil {
		return nil, fmt.Errorf("failed to create environment: %w", err)
	}

	return env, nil
}

func (s *environmentService) GetEnvironmentList(ctx context.Context, req *models.ListEnvironmentsRequest) (*models.EnvironmentListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	filter := &models.EnvironmentFilter{
		Type:     req.Type,
		IsActive: req.IsActive,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	environments, total, err := s.envRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}

	return &models.EnvironmentListResponse{
		Environments: environments,
		Total:        total,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}, nil
}

func (s *environmentService) GetEnvironment(ctx context.Context, id string) (*models.Environment, error) {
	env, err := s.envRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}
	return env, nil
}

func (s *environmentService) UpdateEnvironment(ctx context.Context, id string, req *models.UpdateEnvironmentRequest) (*models.Environment, error) {
	// 验证请求
	if err := utils.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有环境
	env, err := s.envRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("environment not found: %w", err)
	}

	// 更新字段
	if req.Name != "" {
		env.Name = req.Name
	}
	if req.Type != "" {
		env.Type = req.Type
	}
	if req.Config != nil {
		env.Config = req.Config
	}
	if req.IsActive != nil {
		env.IsActive = *req.IsActive
	}
	env.UpdatedAt = time.Now()

	// 保存更新
	if err := s.envRepo.Update(ctx, env); err != nil {
		return nil, fmt.Errorf("failed to update environment: %w", err)
	}

	return env, nil
}

func (s *environmentService) DeleteEnvironment(ctx context.Context, id string) error {
	// 检查环境是否存在
	_, err := s.envRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("environment not found: %w", err)
	}

	// 删除环境
	if err := s.envRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete environment: %w", err)
	}

	return nil
}
