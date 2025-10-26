package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/services/operator-pm-agent/repository"
)

type AgentService struct {
	repository *repository.AgentRepository
	runner     Runner // 应用运行器
	apps       map[string]*models.AgentAppStatus
	mutex      sync.RWMutex
	workDir    string
}

// NewAgentService 创建新的 Agent 服务
func NewAgentService(workDir string) (*AgentService, error) {
	return NewAgentServiceWithConfig(workDir)
}

// NewAgentServiceWithConfig 使用配置创建 Agent 服务
func NewAgentServiceWithConfig(workDir string) (*AgentService, error) {
	// 创建仓库实例
	repo := repository.NewAgentRepository(workDir)

	// 确保工作目录存在
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create work directory: %w", err)
	}

	// 创建运行器（目前只实现 SimpleRunner）
	runner := NewSimpleRunner()

	service := &AgentService{
		repository: repo,
		runner:     runner,
		apps:       make(map[string]*models.AgentAppStatus),
		workDir:    workDir,
	}

	// 启动时恢复已存在的应用状态
	if err := service.restoreAppStates(); err != nil {
		return nil, fmt.Errorf("failed to restore app states: %w", err)
	}

	return service, nil
}

// ApplyApp 应用部署
func (s *AgentService) ApplyApp(req models.ApplyRequest) (*models.ApplyResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 获取部署包
	pkg := &req.Package

	// 检查是否为有效的部署类型
	if pkg.Type != "binary" && pkg.Type != "script" {
		return nil, fmt.Errorf("unsupported deployment type: %s, only 'binary' and 'script' are supported", pkg.Type)
	}

	// 1. 创建应用目录（如果不存在）
	appDir := filepath.Join(s.workDir, req.App)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create app directory: %w", err)
	}

	// 2. 创建版本目录
	versionDir := filepath.Join(appDir, req.Version)
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create version directory: %w", err)
	}

	// 3. 保存应用配置文件
	configPath := filepath.Join(versionDir, "config.json")
	configData := map[string]interface{}{
		"type":        pkg.Type,
		"command":     pkg.Command,
		"args":        pkg.Args,
		"environment": pkg.Environment,
	}
	configJSON, _ := json.MarshalIndent(configData, "", "  ")
	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		return nil, fmt.Errorf("failed to write config file: %w", err)
	}

	// 4. 更新软链接指向新版本
	currentLink := filepath.Join(appDir, "current")
	currentDir := versionDir // current 指向新版本目录

	// 删除旧的软链接
	os.Remove(currentLink)

	// 创建新的软链接
	if err := os.Symlink(req.Version, currentLink); err != nil {
		return nil, fmt.Errorf("failed to create symlink: %w", err)
	}

	// 5. 重启应用（Restart 会自动停止旧进程并启动新进程）
	ctx := context.Background()
	if err := s.runner.Restart(ctx, currentDir, req.Version, pkg); err != nil {
		return nil, fmt.Errorf("failed to restart app: %w", err)
	}

	// 7. 更新应用状态
	appStatus := &models.AgentAppStatus{
		App:     req.App,
		Version: req.Version,
		Healthy: models.HealthStatus{Level: 100, Msg: "Deployed successfully"},
		Status:  "running",
		Updated: time.Now(),
	}
	s.apps[req.App] = appStatus

	// 8. 保存状态到文件
	if err := s.repository.SaveAppStatus(appStatus); err != nil {
		return nil, fmt.Errorf("failed to save app status: %w", err)
	}

	return &models.ApplyResponse{
		Success: true,
		Message: "App deployed successfully",
		App:     req.App,
		Version: req.Version,
	}, nil
}

// GetAllAppStatus 获取所有应用状态
func (s *AgentService) GetAllAppStatus() (*models.StatusResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 更新所有应用的实际状态
	ctx := context.Background()
	for appName := range s.apps {
		if err := s.updateAppHealth(ctx, appName); err != nil {
			// 记录错误但不中断
			fmt.Printf("Failed to update health for app %s: %v\n", appName, err)
		}
	}

	// 构建响应
	apps := make([]models.AgentAppStatus, 0, len(s.apps))
	for _, app := range s.apps {
		apps = append(apps, *app)
	}

	return &models.StatusResponse{Apps: apps}, nil
}

// GetAppStatus 获取指定应用状态
func (s *AgentService) GetAppStatus(appName string) (*models.AppStatusResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	app, exists := s.apps[appName]
	if !exists {
		return nil, fmt.Errorf("app %s not found", appName)
	}

	// 更新应用健康状态
	ctx := context.Background()
	if err := s.updateAppHealth(ctx, appName); err != nil {
		return nil, fmt.Errorf("failed to update app health: %w", err)
	}

	return &models.AppStatusResponse{
		App:     app.App,
		Version: app.Version,
		Healthy: app.Healthy,
	}, nil
}

// parseDeploymentPackage 解析部署包
func (s *AgentService) parseDeploymentPackage(pkg map[string]interface{}) (*models.DeploymentPackage, error) {
	data, err := json.Marshal(pkg)
	if err != nil {
		return nil, err
	}

	var deploymentPkg models.DeploymentPackage
	if err := json.Unmarshal(data, &deploymentPkg); err != nil {
		return nil, err
	}

	return &deploymentPkg, nil
}

// updateAppHealth 更新应用健康状态
func (s *AgentService) updateAppHealth(ctx context.Context, appName string) error {
	app, exists := s.apps[appName]
	if !exists {
		return fmt.Errorf("app %s not found", appName)
	}

	// 使用 runner 获取应用状态
	appDir := filepath.Join(s.workDir, appName)
	status, err := s.runner.Status(ctx, appDir)
	if err != nil {
		// 记录错误但不影响状态更新
		fmt.Printf("Warning: failed to check app status: %v\n", err)
		status = &AppStatus{Running: false, Message: "Status check failed"}
	}

	// 更新状态
	if status.Running {
		app.Status = "running"
		app.Healthy = models.HealthStatus{Level: 100, Msg: status.Message}
	} else {
		app.Status = "stopped"
		app.Healthy = models.HealthStatus{Level: 0, Msg: status.Message}
	}

	app.Updated = time.Now()
	return nil
}

// restoreAppStates 恢复应用状态
func (s *AgentService) restoreAppStates() error {
	apps, err := s.repository.LoadAllAppStatus()
	if err != nil {
		return fmt.Errorf("failed to load app states: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, app := range apps {
		s.apps[app.App] = app
		// 注意：进程不会被恢复，因为它们已经退出
	}

	return nil
}
