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

type applicationService struct {
	appRepo interfaces.ApplicationRepository
}

// NewApplicationService 创建应用服务
func NewApplicationService(appRepo interfaces.ApplicationRepository) interfaces.ApplicationService {
	return &applicationService{
		appRepo: appRepo,
	}
}

func (s *applicationService) CreateApplication(ctx context.Context, req *models.CreateApplicationRequest) (*models.Application, error) {
	// 验证请求
	if err := utils.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 检查应用名称是否已存在
	existingApps, _, err := s.appRepo.List(ctx, &models.ApplicationFilter{
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check existing applications: %w", err)
	}

	for _, app := range existingApps {
		if app.Name == req.Name {
			return nil, fmt.Errorf("application with name %s already exists", req.Name)
		}
	}

	cfg, _ := json.Marshal(req.Config)
	// 创建应用
	app := &models.Application{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Repository:  req.Repository,
		Type:        req.Type,
		Config:      cfg,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return app, nil
}

func (s *applicationService) GetApplicationList(ctx context.Context, req *models.ListApplicationsRequest) (*models.ApplicationListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	filter := &models.ApplicationFilter{
		Repository: req.Repository,
		Type:       req.Type,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	applications, total, err := s.appRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}

	return &models.ApplicationListResponse{
		Applications: applications,
		Total:        total,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}, nil
}

func (s *applicationService) GetApplication(ctx context.Context, id string) (*models.Application, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	return app, nil
}

func (s *applicationService) UpdateApplication(ctx context.Context, id string, req *models.UpdateApplicationRequest) (*models.Application, error) {
	// 验证请求
	if err := utils.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 获取现有应用
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// 更新字段
	if req.Name != "" {
		app.Name = req.Name
	}
	if req.Description != "" {
		app.Description = req.Description
	}
	if req.Type != "" {
		app.Type = req.Type
	}
	if req.Config != nil {
		bs, _ := json.Marshal(req.Config)
		app.Config = bs
	}
	app.UpdatedAt = time.Now()

	// 保存更新
	if err := s.appRepo.Update(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}

	return app, nil
}

func (s *applicationService) DeleteApplication(ctx context.Context, id string) error {
	// 检查应用是否存在
	_, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// 删除应用
	if err := s.appRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	return nil
}
