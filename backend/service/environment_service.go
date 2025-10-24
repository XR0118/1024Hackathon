package service

import (
	"github.com/XR0118/1024Hackathon/backend/model"
)

type EnvironmentService interface {
	Create(env *model.TargetEnvironment) error
	GetByID(id string) (*model.TargetEnvironment, error)
	List(envType, status string) ([]*model.TargetEnvironment, error)
	Update(env *model.TargetEnvironment) error
	HealthCheck(id string) (*HealthCheckResult, error)
}

type HealthCheckResult struct {
	Status string              `json:"status"`
	Checks []HealthCheckDetail `json:"checks"`
}

type HealthCheckDetail struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
