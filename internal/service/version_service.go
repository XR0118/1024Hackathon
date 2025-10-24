package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/models"
	"github.com/boreas/internal/utils"
)

type versionService struct {
	versionRepo interfaces.VersionRepository
}

// NewVersionService 创建版本服务
func NewVersionService(versionRepo interfaces.VersionRepository) interfaces.VersionService {
	return &versionService{
		versionRepo: versionRepo,
	}
}

func (s *versionService) CreateVersion(ctx context.Context, req *models.CreateVersionRequest) (*models.Version, error) {
	// 验证请求
	if err := utils.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 检查版本是否已存在
	existingVersions, _, err := s.versionRepo.List(ctx, &models.VersionFilter{
		Repository: req.Repository,
		Page:       1,
		PageSize:   1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to check existing versions: %w", err)
	}

	for _, version := range existingVersions {
		if version.GitTag == req.GitTag {
			return nil, fmt.Errorf("version with tag %s already exists", req.GitTag)
		}
	}

	// 创建版本
	version := &models.Version{
		ID:          uuid.New().String(),
		GitTag:      req.GitTag,
		GitCommit:   req.GitCommit,
		Repository:  req.Repository,
		CreatedBy:   getCurrentUser(ctx),
		CreatedAt:   time.Now(),
		Description: req.Description,
	}

	if err := s.versionRepo.Create(ctx, version); err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	return version, nil
}

func (s *versionService) GetVersionList(ctx context.Context, req *models.ListVersionsRequest) (*models.VersionListResponse, error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	filter := &models.VersionFilter{
		Repository: req.Repository,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	versions, total, err := s.versionRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	return &models.VersionListResponse{
		Versions: versions,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *versionService) GetVersion(ctx context.Context, id string) (*models.Version, error) {
	version, err := s.versionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}
	return version, nil
}

func (s *versionService) DeleteVersion(ctx context.Context, id string) error {
	// 检查版本是否存在
	_, err := s.versionRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}

	// 删除版本
	if err := s.versionRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete version: %w", err)
	}

	return nil
}

// getCurrentUser 获取当前用户（从上下文）
func getCurrentUser(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return "system"
}
