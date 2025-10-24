package postgres

import (
	"context"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/models"
	"gorm.io/gorm"
)

type environmentRepository struct {
	db *gorm.DB
}

// NewEnvironmentRepository 创建环境仓库
func NewEnvironmentRepository(db *gorm.DB) interfaces.EnvironmentRepository {
	return &environmentRepository{db: db}
}

func (r *environmentRepository) Create(ctx context.Context, env *models.Environment) error {
	return r.db.WithContext(ctx).Create(env).Error
}

func (r *environmentRepository) GetByID(ctx context.Context, id string) (*models.Environment, error) {
	var env models.Environment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&env).Error
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func (r *environmentRepository) List(ctx context.Context, filter *models.EnvironmentFilter) ([]*models.Environment, int, error) {
	var environments []*models.Environment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Environment{})

	// 应用过滤器
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
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

	// 查询数据
	if err := query.Find(&environments).Error; err != nil {
		return nil, 0, err
	}

	return environments, int(total), nil
}

func (r *environmentRepository) Update(ctx context.Context, env *models.Environment) error {
	return r.db.WithContext(ctx).Save(env).Error
}

func (r *environmentRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Environment{}).Error
}
