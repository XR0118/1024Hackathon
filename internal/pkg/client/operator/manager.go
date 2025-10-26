package operator

import (
	"context"
	"fmt"
	"sync"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
)

// 编译时检查，确保 Manager 实现了 OperatorManager 接口
var _ interfaces.OperatorManager = (*Manager)(nil)

// Manager Operator 管理器
// 负责管理所有类型的 Operator 客户端，并根据环境类型选择合适的 Operator
type Manager struct {
	operators map[string]interfaces.Operator // key: environment_id
	mu        sync.RWMutex
}

// NewManager 创建 Operator 管理器
func NewManager() *Manager {
	return &Manager{
		operators: make(map[string]interfaces.Operator),
	}
}

// RegisterOperator 注册 Operator
func (m *Manager) RegisterOperator(environmentID string, operator interfaces.Operator) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.operators[environmentID] = operator
}

// GetOperator 获取指定环境的 Operator
func (m *Manager) GetOperator(environmentID string) (interfaces.Operator, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	operator, ok := m.operators[environmentID]
	if !ok {
		return nil, fmt.Errorf("operator not found for environment: %s", environmentID)
	}

	return operator, nil
}

// GetOperatorByEnvironment 根据环境对象获取 Operator
func (m *Manager) GetOperatorByEnvironment(env *models.Environment) (interfaces.Operator, error) {
	return m.GetOperator(env.ID)
}

// ApplyDeployment 在指定环境应用部署
func (m *Manager) ApplyDeployment(ctx context.Context, environmentID string, req *models.ApplyDeploymentRequest) (*models.ApplyDeploymentResponse, error) {
	operator, err := m.GetOperator(environmentID)
	if err != nil {
		return nil, err
	}

	return operator.Apply(ctx, req)
}

// GetApplicationStatus 获取指定环境中应用的状态
func (m *Manager) GetApplicationStatus(ctx context.Context, environmentID string, appName string) (*models.ApplicationStatusResponse, error) {
	operator, err := m.GetOperator(environmentID)
	if err != nil {
		return nil, err
	}

	return operator.GetApplicationStatus(ctx, appName)
}

// HealthCheckAll 检查所有 Operator 的健康状态
func (m *Manager) HealthCheckAll(ctx context.Context) map[string]error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make(map[string]error)
	for envID, operator := range m.operators {
		results[envID] = operator.HealthCheck(ctx)
	}

	return results
}

// RemoveOperator 移除指定环境的 Operator
func (m *Manager) RemoveOperator(environmentID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.operators, environmentID)
}

// ListOperators 列出所有已注册的 Operator
func (m *Manager) ListOperators() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	envIDs := make([]string, 0, len(m.operators))
	for envID := range m.operators {
		envIDs = append(envIDs, envID)
	}

	return envIDs
}
