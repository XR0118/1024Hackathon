package service

import (
	"context"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
)

type DeployExecutor interface {
	Apply(ctx context.Context) error
}

type DeployClient interface {
	Apply(ctx context.Context, app string, version string, pkg models.DeploymentPackage) (models.ApplyResponse, error)
	AppStatus(ctx context.Context, app string) ([]models.AgentAppStatus, error)
}

type SimpleDeployExecutor struct {
	task           models.Task
	deploymentRepo interfaces.DeploymentRepository
	client         DeployClient
}

func NewSimpleDeployExecutor(task models.Task, deploymentRepo interfaces.DeploymentRepository, client DeployClient) *SimpleDeployExecutor {
	return &SimpleDeployExecutor{
		task:           task,
		deploymentRepo: deploymentRepo,
		client:         client,
	}
}
