package operator

import (
	"fmt"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
)

// Config Operator 配置
type Config struct {
	// K8S Operator 配置
	K8SOperatorURL string `mapstructure:"k8s_operator_url"`

	// PM Operator 配置
	PMOperatorURL string `mapstructure:"pm_operator_url"`

	// Mock Operator 配置
	UseMock bool `mapstructure:"use_mock"`
}

// CreateOperatorFromEnvironment 根据环境配置创建对应的 Operator 客户端
func CreateOperatorFromEnvironment(env *models.Environment, config *Config) (interfaces.Operator, error) {
	switch env.Type {
	case "kubernetes":
		if config.UseMock {
			return NewMockClient(), nil
		}
		if config.K8SOperatorURL == "" {
			return nil, fmt.Errorf("k8s operator URL not configured")
		}
		return NewK8sClient(config.K8SOperatorURL), nil

	case "physical":
		if config.UseMock {
			return NewMockClient(), nil
		}
		if config.PMOperatorURL == "" {
			return nil, fmt.Errorf("pm operator URL not configured")
		}
		return NewPMClient(config.PMOperatorURL), nil

	default:
		return nil, fmt.Errorf("unsupported environment type: %s", env.Type)
	}
}

// CreateOperatorByType 根据类型直接创建 Operator 客户端
func CreateOperatorByType(envType string, baseURL string, useMock bool) (interfaces.Operator, error) {
	if useMock {
		return NewMockClient(), nil
	}

	switch envType {
	case "kubernetes":
		if baseURL == "" {
			return nil, fmt.Errorf("base URL is required for k8s operator")
		}
		return NewK8sClient(baseURL), nil

	case "physical":
		if baseURL == "" {
			return nil, fmt.Errorf("base URL is required for pm operator")
		}
		return NewPMClient(baseURL), nil

	default:
		return nil, fmt.Errorf("unsupported environment type: %s", envType)
	}
}

// InitializeOperators 初始化所有环境的 Operator 客户端
func InitializeOperators(environments []*models.Environment, config *Config) (*Manager, error) {
	manager := NewManager()

	for _, env := range environments {
		operator, err := CreateOperatorFromEnvironment(env, config)
		if err != nil {
			return nil, fmt.Errorf("failed to create operator for environment %s: %w", env.Name, err)
		}

		manager.RegisterOperator(env.ID, operator)
	}

	return manager, nil
}
