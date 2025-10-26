package operator

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/boreas/internal/pkg/models"
)

// MockClient Mock Operator 客户端（用于测试和演示）
type MockClient struct {
	deployments map[string]*mockDeployment
}

type mockDeployment struct {
	app       string
	versions  []models.VersionStatus
	updatedAt time.Time
}

// NewMockClient 创建 Mock Operator 客户端
func NewMockClient() *MockClient {
	return &MockClient{
		deployments: make(map[string]*mockDeployment),
	}
}

// Apply 应用部署（模拟）
func (c *MockClient) Apply(ctx context.Context, req *models.ApplyDeploymentRequest) (*models.ApplyDeploymentResponse, error) {
	// 模拟部署延迟
	time.Sleep(500 * time.Millisecond)

	// 构建版本状态
	versions := make([]models.VersionStatus, 0, len(req.Versions))
	for _, v := range req.Versions {
		// 模拟节点状态
		nodes := []models.NodeStatus{
			{
				Node:    fmt.Sprintf("node-%d", rand.Intn(100)),
				Healthy: models.HealthInfo{Level: 85 + rand.Intn(15)},
			},
			{
				Node:    fmt.Sprintf("node-%d", rand.Intn(100)),
				Healthy: models.HealthInfo{Level: 85 + rand.Intn(15)},
			},
		}

		versions = append(versions, models.VersionStatus{
			Version: v.Version,
			Percent: v.Percent,
			Healthy: models.HealthInfo{Level: 90 + rand.Intn(10)},
			Nodes:   nodes,
		})
	}

	// 保存模拟的部署状态
	c.deployments[req.App] = &mockDeployment{
		app:       req.App,
		versions:  versions,
		updatedAt: time.Now(),
	}

	return &models.ApplyDeploymentResponse{
		App:     req.App,
		Message: "Deployment applied successfully (mock)",
		Success: true,
	}, nil
}

// GetApplicationStatus 获取应用状态（模拟）
func (c *MockClient) GetApplicationStatus(ctx context.Context, appName string) (*models.ApplicationStatusResponse, error) {
	deployment, ok := c.deployments[appName]
	if !ok {
		// 返回默认的模拟状态
		return &models.ApplicationStatusResponse{
			App: appName,
			Healthy: models.HealthInfo{
				Level: 0,
			},
			Versions: []models.VersionStatus{},
		}, nil
	}

	// 计算总体健康度
	totalHealth := 0
	for _, v := range deployment.versions {
		totalHealth += v.Healthy.Level
	}
	avgHealth := 0
	if len(deployment.versions) > 0 {
		avgHealth = totalHealth / len(deployment.versions)
	}

	return &models.ApplicationStatusResponse{
		App: appName,
		Healthy: models.HealthInfo{
			Level: avgHealth,
		},
		Versions: deployment.versions,
	}, nil
}

// HealthCheck 健康检查（模拟）
func (c *MockClient) HealthCheck(ctx context.Context) error {
	// Mock 客户端总是健康
	return nil
}

// GetType 获取 Operator 类型
func (c *MockClient) GetType() string {
	return "mock"
}
