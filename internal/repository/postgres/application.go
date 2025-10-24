package postgres

import (
	"context"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/models"
	"gorm.io/gorm"
)

type applicationRepository struct {
	db *gorm.DB
}

// NewApplicationRepository 创建应用仓库
func NewApplicationRepository(db *gorm.DB) interfaces.ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(ctx context.Context, app *models.Application) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *applicationRepository) GetByID(ctx context.Context, id string) (*models.Application, error) {
	var app models.Application
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) List(ctx context.Context, filter *models.ApplicationFilter) ([]*models.Application, int, error) {
	var applications []*models.Application
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Application{})

	// 应用过滤器
	if filter.Repository != "" {
		query = query.Where("repository = ?", filter.Repository)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
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
	if err := query.Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, int(total), nil
}

func (r *applicationRepository) Update(ctx context.Context, app *models.Application) error {
	return r.db.WithContext(ctx).Save(app).Error
}

func (r *applicationRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Application{}).Error
}
