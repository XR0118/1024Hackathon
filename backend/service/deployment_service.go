package service

import (
	"github.com/XR0118/1024Hackathon/backend/model"
)

type DeploymentService interface {
	Create(deployment *model.Deployment) error
	GetByID(id string) (*model.Deployment, error)
	List(status, envID string, page int) ([]*model.Deployment, int, error)
	Start(id string) error
	Pause(id string) error
	Resume(id string) error
	Cancel(id string) error
	Approve(id string, approval *model.Approval) error
	Rollback(id string) (*model.Deployment, error)
	GetLogs(id, appID, envID string) ([]DeploymentLog, error)
}

type DeploymentLog struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}
