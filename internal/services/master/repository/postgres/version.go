package postgres

import (
	"context"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"gorm.io/gorm"
)

type versionRepository struct {
	db *gorm.DB
}

// NewVersionRepository 创建版本仓库
func NewVersionRepository(db *gorm.DB) interfaces.VersionRepository {
	return &versionRepository{db: db}
}

func (r *versionRepository) Create(ctx context.Context, version *models.Version) error {
	return r.db.WithContext(ctx).Create(version).Error
}

func (r *versionRepository) GetByID(ctx context.Context, id string) (*models.Version, error) {
	var version models.Version
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *versionRepository) List(ctx context.Context, filter *models.VersionFilter) ([]*models.Version, int, error) {
	var versions []*models.Version
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Version{})

	// 应用过滤器
	if filter.Repository != "" {
		query = query.Where("repository = ?", filter.Repository)
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
	if err := query.Find(&versions).Error; err != nil {
		return nil, 0, err
	}

	return versions, int(total), nil
}

// GetPreviousByVersionAndApp returns the latest version before the target version's creation time
// that contains the specified app_id within its app_builds JSONB array.
func (r *versionRepository) GetPreviousByVersionAndApp(ctx context.Context, targetVersionID string, appID string) (*models.Version, error) {
	// Fetch target version to get its created_at
	var target models.Version
	if err := r.db.WithContext(ctx).Where("id = ?", targetVersionID).First(&target).Error; err != nil {
		return nil, err
	}

	var prev models.Version
	// Use EXISTS with jsonb_array_elements on app_builds to filter by app_id
	err := r.db.WithContext(ctx).
		Where("created_at < ?", target.CreatedAt).
		Where("EXISTS (SELECT 1 FROM jsonb_array_elements(app_builds) AS elem WHERE elem->>'app_id' = ?)", appID).
		Order("created_at DESC").
		Limit(1).
		First(&prev).Error
	if err != nil {
		return nil, err
	}
	return &prev, nil
}

func (r *versionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Version{}).Error
}
