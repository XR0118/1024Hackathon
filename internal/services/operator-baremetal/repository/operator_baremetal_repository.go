package repository

import (
	"github.com/XR0118/1024Hackathon/internal/pkg/models"
	"gorm.io/gorm"
)

type OperatorBaremetalRepository struct {
	db *gorm.DB
}

func NewOperatorBaremetalRepository(db *gorm.DB) *OperatorBaremetalRepository {
	return &OperatorBaremetalRepository{
		db: db,
	}
}

// GetDeploymentByID 根据ID获取部署信息
func (r *OperatorBaremetalRepository) GetDeploymentByID(id string) (*models.Deployment, error) {
	var deployment models.Deployment
	if err := r.db.Where("id = ?", id).First(&deployment).Error; err != nil {
		return nil, err
	}
	return &deployment, nil
}

// UpdateDeployment 更新部署信息
func (r *OperatorBaremetalRepository) UpdateDeployment(deployment *models.Deployment) error {
	return r.db.Save(deployment).Error
}

// GetDeploymentsByStatus 根据状态获取部署列表
func (r *OperatorBaremetalRepository) GetDeploymentsByStatus(status models.DeploymentStatus) ([]*models.Deployment, error) {
	var deployments []*models.Deployment
	if err := r.db.Where("status = ?", status).Find(&deployments).Error; err != nil {
		return nil, err
	}
	return deployments, nil
}

// CreateDeploymentLog 创建部署日志
func (r *OperatorBaremetalRepository) CreateDeploymentLog(log *models.DeploymentLog) error {
	return r.db.Create(log).Error
}

// GetDeploymentLogs 获取部署日志
func (r *OperatorBaremetalRepository) GetDeploymentLogs(deploymentID string) ([]*models.DeploymentLog, error) {
	var logs []*models.DeploymentLog
	if err := r.db.Where("deployment_id = ?", deploymentID).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
