package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boreas/internal/pkg/models"
)

type AgentRepository struct {
	workDir string
}

func NewAgentRepository(workDir string) *AgentRepository {
	return &AgentRepository{
		workDir: workDir,
	}
}

// SaveAppStatus 保存应用状态
func (r *AgentRepository) SaveAppStatus(app *models.AgentAppStatus) error {
	appDir := filepath.Join(r.workDir, "apps", app.App)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create app directory: %w", err)
	}

	statusFile := filepath.Join(appDir, "status.json")
	data, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal app status: %w", err)
	}

	if err := os.WriteFile(statusFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write status file: %w", err)
	}

	return nil
}

// LoadAppStatus 加载应用状态
func (r *AgentRepository) LoadAppStatus(appName string) (*models.AgentAppStatus, error) {
	statusFile := filepath.Join(r.workDir, "apps", appName, "status.json")

	data, err := os.ReadFile(statusFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("app %s not found", appName)
		}
		return nil, fmt.Errorf("failed to read status file: %w", err)
	}

	var app models.AgentAppStatus
	if err := json.Unmarshal(data, &app); err != nil {
		return nil, fmt.Errorf("failed to unmarshal app status: %w", err)
	}

	return &app, nil
}

// LoadAllAppStatus 加载所有应用状态
func (r *AgentRepository) LoadAllAppStatus() ([]*models.AgentAppStatus, error) {
	appsDir := filepath.Join(r.workDir, "apps")

	// 检查apps目录是否存在
	if _, err := os.Stat(appsDir); os.IsNotExist(err) {
		return []*models.AgentAppStatus{}, nil
	}

	// 读取所有应用目录
	entries, err := os.ReadDir(appsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read apps directory: %w", err)
	}

	var apps []*models.AgentAppStatus
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		appName := entry.Name()
		app, err := r.LoadAppStatus(appName)
		if err != nil {
			// 记录错误但继续处理其他应用
			fmt.Printf("Warning: failed to load status for app %s: %v\n", appName, err)
			continue
		}

		apps = append(apps, app)
	}

	return apps, nil
}

// DeleteAppStatus 删除应用状态
func (r *AgentRepository) DeleteAppStatus(appName string) error {
	appDir := filepath.Join(r.workDir, "apps", appName)
	return os.RemoveAll(appDir)
}

// SaveAgentConfig 保存Agent配置
func (r *AgentRepository) SaveAgentConfig(config *models.AgentConfig) error {
	configFile := filepath.Join(r.workDir, "agent-config.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal agent config: %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadAgentConfig 加载Agent配置
func (r *AgentRepository) LoadAgentConfig() (*models.AgentConfig, error) {
	configFile := filepath.Join(r.workDir, "agent-config.json")

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("agent config not found")
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config models.AgentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal agent config: %w", err)
	}

	return &config, nil
}
