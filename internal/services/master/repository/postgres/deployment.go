package postgres

import (
	"context"

	"github.com/XR0118/1024Hackathon/internal/interfaces"
	"github.com/XR0118/1024Hackathon/internal/pkg/models"
	"gorm.io/gorm"
)

type deploymentRepository struct {
	db *gorm.DB
}

// NewDeploymentRepository 创建部署仓库
func NewDeploymentRepository(db *gorm.DB) interfaces.DeploymentRepository {
	return &deploymentRepository{db: db}
}

func (r *deploymentRepository) Create(ctx context.Context, deployment *models.Deployment) error {
	return r.db.WithContext(ctx).Create(deployment).Error
}

func (r *deploymentRepository) GetByID(ctx context.Context, id string) (*models.Deployment, error) {
	var deployment models.Deployment
	err := r.db.WithContext(ctx).
		Preload("Version").
		Preload("Environment").
		Preload("Applications").
		Where("id = ?", id).
		First(&deployment).Error
	if err != nil {
		return nil, err
	}
	return &deployment, nil
}

func (r *deploymentRepository) List(ctx context.Context, filter *models.DeploymentFilter) ([]*models.Deployment, int, error) {
	var deployments []*models.Deployment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Deployment{})

	// 应用过滤器
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.EnvironmentID != "" {
		query = query.Where("environment_id = ?", filter.EnvironmentID)
	}
	if filter.VersionID != "" {
		query = query.Where("version_id = ?", filter.VersionID)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 排序
	query = query.Order("created_at DESC")

	// 预加载关联数据
	query = query.Preload("Version").Preload("Environment").Preload("Applications")

	// 查询数据
	if err := query.Find(&deployments).Error; err != nil {
		return nil, 0, err
	}

	return deployments, int(total), nil
}

func (r *deploymentRepository) Update(ctx context.Context, deployment *models.Deployment) error {
	return r.db.WithContext(ctx).Save(deployment).Error
}
