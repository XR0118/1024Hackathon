package repository

import (
	"github.com/boreas/internal/pkg/models"
	"gorm.io/gorm"
)

type OperatorPMRepository struct {
	db *gorm.DB
}

func NewOperatorPMRepository(db *gorm.DB) *OperatorPMRepository {
	return &OperatorPMRepository{
		db: db,
	}
}

// GetDeploymentByID 根据ID获取部署信息
func (r *OperatorPMRepository) GetDeploymentByID(id string) (*models.Deployment, error) {
	var deployment models.Deployment
	if err := r.db.Where("id = ?", id).First(&deployment).Error; err != nil {
		return nil, err
	}
	return &deployment, nil
}

// UpdateDeployment 更新部署信息
func (r *OperatorPMRepository) UpdateDeployment(deployment *models.Deployment) error {
	return r.db.Save(deployment).Error
}

// GetDeploymentsByStatus 根据状态获取部署列表
func (r *OperatorPMRepository) GetDeploymentsByStatus(status models.DeploymentStatus) ([]*models.Deployment, error) {
	var deployments []*models.Deployment
	if err := r.db.Where("status = ?", status).Find(&deployments).Error; err != nil {
		return nil, err
	}
	return deployments, nil
}

// CreateDeploymentLog 创建部署日志
func (r *OperatorPMRepository) CreateDeploymentLog(log *models.DeploymentLog) error {
	return r.db.Create(log).Error
}

// GetDeploymentLogs 获取部署日志
func (r *OperatorPMRepository) GetDeploymentLogs(deploymentID string) ([]*models.DeploymentLog, error) {
	var logs []*models.DeploymentLog
	if err := r.db.Where("deployment_id = ?", deploymentID).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// GetEnvironmentByID 根据ID获取环境信息
func (r *OperatorPMRepository) GetEnvironmentByID(id string) (*models.Environment, error) {
	var environment models.Environment
	if err := r.db.Where("id = ?", id).First(&environment).Error; err != nil {
		return nil, err
	}
	return &environment, nil
}
