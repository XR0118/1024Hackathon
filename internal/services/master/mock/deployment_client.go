package mock

import (
	"context"
	"strconv"
	"sync"

	"github.com/boreas/internal/pkg/models"
)

type MockDeploymentClient struct {
	sync.RWMutex
	Instances map[string]models.AgentAppStatus
}

func NewMockDeploymentClient() *MockDeploymentClient {
	return &MockDeploymentClient{
		Instances: make(map[string]models.AgentAppStatus),
	}
}

func (c *MockDeploymentClient) Apply(ctx context.Context, app string, version string, pkg models.DeploymentPackage) (models.ApplyResponse, error) {
	c.Lock()
	defer c.Unlock()

	level := 100
	if s, ok := pkg.Environment["health"]; ok {
		level, _ = strconv.Atoi(s)
	}

	key := app + ":" + version
	c.Instances[key] = models.AgentAppStatus{
		App:      app,
		Version:  version,
		Replicas: pkg.Replicas,
		Healthy: models.HealthStatus{
			Level: level,
		},
	}

	return models.ApplyResponse{
		Success: true,
	}, nil
}

func (c *MockDeploymentClient) AppStatus(ctx context.Context, app string) ([]models.AgentAppStatus, error) {
	c.RLock()
	defer c.RUnlock()

	var statuses []models.AgentAppStatus
	for _, status := range c.Instances {
		if status.App == app {
			statuses = append(statuses, status)
		}
	}
	return statuses, nil
}
