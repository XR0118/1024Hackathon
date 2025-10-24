package postgres

import (
	"context"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/models"
	"gorm.io/gorm"
)

type workflowRepository struct {
	db *gorm.DB
}

// NewWorkflowRepository 创建工作流仓库
func NewWorkflowRepository(db *gorm.DB) interfaces.WorkflowRepository {
	return &workflowRepository{db: db}
}

func (r *workflowRepository) Create(ctx context.Context, workflow *models.Workflow) error {
	return r.db.WithContext(ctx).Create(workflow).Error
}

func (r *workflowRepository) GetByID(ctx context.Context, id string) (*models.Workflow, error) {
	var workflow models.Workflow
	err := r.db.WithContext(ctx).
		Preload("Deployment").
		Preload("Tasks").
		Where("id = ?", id).
		First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *workflowRepository) GetByDeploymentID(ctx context.Context, deploymentID string) (*models.Workflow, error) {
	var workflow models.Workflow
	err := r.db.WithContext(ctx).
		Preload("Deployment").
		Preload("Tasks").
		Where("deployment_id = ?", deploymentID).
		First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *workflowRepository) Update(ctx context.Context, workflow *models.Workflow) error {
	return r.db.WithContext(ctx).Save(workflow).Error
}
