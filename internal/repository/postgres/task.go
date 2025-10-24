package postgres

import (
	"context"

	"github.com/guaguasong/1024Hackathon/internal/interfaces"
	"github.com/guaguasong/1024Hackathon/internal/models"
	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository 创建任务仓库
func NewTaskRepository(db *gorm.DB) interfaces.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *taskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).
		Preload("Deployment").
		Where("id = ?", id).
		First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) List(ctx context.Context, filter *models.TaskFilter) ([]*models.Task, int, error) {
	var tasks []*models.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Task{})

	// 应用过滤器
	if filter.DeploymentID != "" {
		query = query.Where("deployment_id = ?", filter.DeploymentID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
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
	query = query.Order("created_at ASC")

	// 预加载关联数据
	query = query.Preload("Deployment")

	// 查询数据
	if err := query.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, int(total), nil
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *taskRepository) GetByDeploymentID(ctx context.Context, deploymentID string) ([]*models.Task, error) {
	var tasks []*models.Task
	err := r.db.WithContext(ctx).
		Where("deployment_id = ?", deploymentID).
		Order("created_at ASC").
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
