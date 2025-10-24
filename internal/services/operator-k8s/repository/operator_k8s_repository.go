package repository

import (
	"github.com/boreas/internal/pkg/models"
	"gorm.io/gorm"
)

type OperatorK8sRepository struct {
	db *gorm.DB
}

func NewOperatorK8sRepository(db *gorm.DB) *OperatorK8sRepository {
	return &OperatorK8sRepository{
		db: db,
	}
}

// GetDeploymentByID 根据ID获取部署信息
func (r *OperatorK8sRepository) GetDeploymentByID(id string) (*models.Deployment, error) {
	var deployment models.Deployment
	if err := r.db.Where("id = ?", id).First(&deployment).Error; err != nil {
		return nil, err
	}
	return &deployment, nil
}

// UpdateDeployment 更新部署信息
func (r *OperatorK8sRepository) UpdateDeployment(deployment *models.Deployment) error {
	return r.db.Save(deployment).Error
}

// GetDeploymentsByStatus 根据状态获取部署列表
func (r *OperatorK8sRepository) GetDeploymentsByStatus(status models.DeploymentStatus) ([]*models.Deployment, error) {
	var deployments []*models.Deployment
	if err := r.db.Where("status = ?", status).Find(&deployments).Error; err != nil {
		return nil, err
	}
	return deployments, nil
}

// CreateDeploymentLog 创建部署日志
func (r *OperatorK8sRepository) CreateDeploymentLog(log *models.DeploymentLog) error {
	return r.db.Create(log).Error
}

// GetDeploymentLogs 获取部署日志
func (r *OperatorK8sRepository) GetDeploymentLogs(deploymentID string) ([]*models.DeploymentLog, error) {
	var logs []*models.DeploymentLog
	if err := r.db.Where("deployment_id = ?", deploymentID).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
