package service

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/boreas/internal/pkg/logger"
	"github.com/boreas/internal/pkg/models"
	"go.uber.org/zap"
)

var MockAgent = NewMockDeploymentClient()

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

	logger.GetLogger().Info("apply deployment package",
		zap.Any("pkg", pkg), zap.Any("app", app), zap.Any("version", version),
	)
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
		Status:  "running",
		Updated: time.Now(),
	}

	return models.ApplyResponse{
		Success: true,
	}, nil
}

func (c *MockDeploymentClient) AppStatus(ctx context.Context, app string) ([]models.AgentAppStatus, error) {
	c.RLock()
	defer c.RUnlock()

	logger.GetLogger().Info("instance status",
		zap.Any("instances", c.Instances),
	)

	var statuses []models.AgentAppStatus
	for _, status := range c.Instances {
		if status.App == app {
			statuses = append(statuses, status)
		}
	}

	return statuses, nil
}
