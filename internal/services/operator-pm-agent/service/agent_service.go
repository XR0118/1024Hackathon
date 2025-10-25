package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/services/operator-pm-agent/repository"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type AgentService struct {
	repository   *repository.AgentRepository
	dockerClient *client.Client
	apps         map[string]*models.AgentAppStatus
	mutex        sync.RWMutex
	workDir      string
}

func NewAgentService(workDir string) (*AgentService, error) {
	// 初始化Docker客户端
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// 创建仓库实例
	repo := repository.NewAgentRepository(workDir)

	// 确保工作目录存在
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create work directory: %w", err)
	}

	service := &AgentService{
		repository:   repo,
		dockerClient: dockerClient,
		apps:         make(map[string]*models.AgentAppStatus),
		workDir:      workDir,
	}

	// 启动时恢复已存在的应用状态
	if err := service.restoreAppStates(); err != nil {
		return nil, fmt.Errorf("failed to restore app states: %w", err)
	}

	return service, nil
}

// NewAgentServiceWithConfig 使用配置创建Agent服务
func NewAgentServiceWithConfig(workDir string, dockerEnabled bool, dockerSocketPath string) (*AgentService, error) {
	var dockerClient *client.Client
	var err error

	// 只有在启用Docker时才初始化Docker客户端
	if dockerEnabled {
		dockerClient, err = client.NewClientWithOpts(
			client.FromEnv,
			client.WithAPIVersionNegotiation(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create docker client: %w", err)
		}
	}

	// 创建仓库实例
	repo := repository.NewAgentRepository(workDir)

	// 确保工作目录存在
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create work directory: %w", err)
	}

	service := &AgentService{
		repository:   repo,
		dockerClient: dockerClient,
		apps:         make(map[string]*models.AgentAppStatus),
		workDir:      workDir,
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

	// 解析部署包
	pkg, err := s.parseDeploymentPackage(req.Pkg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse deployment package: %w", err)
	}

	// 停止旧版本应用（如果存在）
	if _, exists := s.apps[req.App]; exists {
		if err := s.stopApp(req.App); err != nil {
			return nil, fmt.Errorf("failed to stop old app: %w", err)
		}
	}

	// 部署新应用
	if err := s.deployApp(req.App, req.Version, pkg); err != nil {
		return nil, fmt.Errorf("failed to deploy app: %w", err)
	}

	// 更新应用状态
	appStatus := &models.AgentAppStatus{
		App:     req.App,
		Version: req.Version,
		Healthy: models.HealthStatus{Level: 100, Msg: "Deployed successfully"},
		Status:  "running",
		Updated: time.Now(),
	}
	s.apps[req.App] = appStatus

	// 保存状态到文件
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
	for appName := range s.apps {
		if err := s.updateAppHealth(appName); err != nil {
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
	if err := s.updateAppHealth(appName); err != nil {
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

// deployApp 部署应用
func (s *AgentService) deployApp(appName, version string, pkg *models.DeploymentPackage) error {
	switch pkg.Type {
	case "docker":
		return s.deployDockerApp(appName, version, pkg)
	case "binary":
		return s.deployBinaryApp(appName, version, pkg)
	case "script":
		return s.deployScriptApp(appName, version, pkg)
	default:
		return fmt.Errorf("unsupported deployment type: %s", pkg.Type)
	}
}

// deployDockerApp 部署Docker应用
func (s *AgentService) deployDockerApp(appName, version string, pkg *models.DeploymentPackage) error {
	containerName := fmt.Sprintf("%s-%s", appName, version)

	// 停止并删除已存在的容器
	if err := s.stopDockerContainer(containerName); err != nil {
		fmt.Printf("Warning: failed to stop existing container: %v\n", err)
	}

	// 拉取镜像
	// TODO: 实现镜像拉取逻辑
	// if err := s.pullDockerImage(pkg.Image); err != nil {
	// 	return fmt.Errorf("failed to pull image: %w", err)
	// }

	// 创建容器配置
	containerConfig := &container.Config{
		Image: pkg.Image,
		Env:   s.convertEnvVars(pkg.Environment),
		Cmd:   pkg.Command,
	}

	hostConfig := &container.HostConfig{
		Binds: s.convertVolumeMounts(pkg.Volumes),
	}

	// 创建容器
	_, err := s.dockerClient.ContainerCreate(
		context.Background(),
		containerConfig,
		hostConfig,
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// 启动容器
	// TODO: 实现容器启动逻辑
	// if err := s.dockerClient.ContainerStart(
	// 	context.Background(),
	// 	resp.ID,
	// 	types.ContainerStartOptions{},
	// ); err != nil {
	// 	return fmt.Errorf("failed to start container: %w", err)
	// }

	return nil
}

// deployBinaryApp 部署二进制应用
func (s *AgentService) deployBinaryApp(appName, version string, pkg *models.DeploymentPackage) error {
	appDir := filepath.Join(s.workDir, "apps", appName, version)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create app directory: %w", err)
	}

	// 这里应该下载二进制文件，简化实现
	// 实际实现中需要从包管理器中下载或从镜像中提取

	// 创建启动脚本
	scriptPath := filepath.Join(appDir, "start.sh")
	script := fmt.Sprintf(`#!/bin/bash
cd %s
exec %s %s
`, appDir, strings.Join(pkg.Command, " "), strings.Join(pkg.Args, " "))

	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("failed to create start script: %w", err)
	}

	// 启动应用
	cmd := exec.Command("bash", scriptPath)
	cmd.Dir = appDir
	cmd.Env = s.convertEnvVars(pkg.Environment)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start binary app: %w", err)
	}

	return nil
}

// deployScriptApp 部署脚本应用
func (s *AgentService) deployScriptApp(appName, version string, pkg *models.DeploymentPackage) error {
	appDir := filepath.Join(s.workDir, "apps", appName, version)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create app directory: %w", err)
	}

	// 创建脚本文件
	scriptPath := filepath.Join(appDir, "script.sh")
	script := strings.Join(pkg.Command, "\n")

	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("failed to create script: %w", err)
	}

	// 启动脚本
	cmd := exec.Command("bash", scriptPath)
	cmd.Dir = appDir
	cmd.Env = s.convertEnvVars(pkg.Environment)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start script app: %w", err)
	}

	return nil
}

// stopApp 停止应用
func (s *AgentService) stopApp(appName string) error {
	app, exists := s.apps[appName]
	if !exists {
		return nil
	}

	containerName := fmt.Sprintf("%s-%s", appName, app.Version)
	if err := s.stopDockerContainer(containerName); err != nil {
		fmt.Printf("Warning: failed to stop docker container: %v\n", err)
	}

	// 更新状态
	app.Status = "stopped"
	app.Healthy = models.HealthStatus{Level: 0, Msg: "Stopped"}
	app.Updated = time.Now()

	return nil
}

// updateAppHealth 更新应用健康状态
func (s *AgentService) updateAppHealth(appName string) error {
	app, exists := s.apps[appName]
	if !exists {
		return fmt.Errorf("app %s not found", appName)
	}

	// 检查Docker容器状态
	containerName := fmt.Sprintf("%s-%s", appName, app.Version)
	containerInfo, err := s.dockerClient.ContainerInspect(context.Background(), containerName)
	if err != nil {
		// 容器不存在或出错
		app.Status = "stopped"
		app.Healthy = models.HealthStatus{Level: 0, Msg: "Container not found"}
	} else {
		switch containerInfo.State.Status {
		case "running":
			app.Status = "running"
			app.Healthy = models.HealthStatus{Level: 100, Msg: "Running"}
		case "exited":
			app.Status = "stopped"
			app.Healthy = models.HealthStatus{Level: 0, Msg: "Exited"}
		default:
			app.Status = "unknown"
			app.Healthy = models.HealthStatus{Level: 50, Msg: "Unknown status"}
		}
	}

	app.Updated = time.Now()
	return nil
}

// stopDockerContainer 停止Docker容器
func (s *AgentService) stopDockerContainer(containerName string) error {
	// 尝试停止容器
	if err := s.dockerClient.ContainerStop(context.Background(), containerName, container.StopOptions{}); err != nil {
		// 忽略容器不存在的错误
		if !strings.Contains(err.Error(), "No such container") {
			return err
		}
	}

	// 删除容器
	// TODO: 实现容器删除逻辑
	// if err := s.dockerClient.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{}); err != nil {
	// 	if !strings.Contains(err.Error(), "No such container") {
	// 		return err
	// 	}
	// }

	return nil
}

// pullDockerImage 拉取Docker镜像
func (s *AgentService) pullDockerImage(image string) error {
	// TODO: 实现镜像拉取逻辑
	// reader, err := s.dockerClient.ImagePull(context.Background(), image, types.ImagePullOptions{})
	// if err != nil {
	// 	return err
	// }
	// defer reader.Close()

	// 读取响应以确保拉取完成
	// _, err = io.Copy(io.Discard, reader)
	// return err
	return nil
}

// convertEnvVars 转换环境变量
func (s *AgentService) convertEnvVars(env map[string]string) []string {
	var envVars []string
	for k, v := range env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}
	return envVars
}

// convertPortMappings 转换端口映射
func (s *AgentService) convertPortMappings(ports []models.PortMapping) map[string][]string {
	portBindings := make(map[string][]string)
	for _, port := range ports {
		key := fmt.Sprintf("%d/%s", port.ContainerPort, port.Protocol)
		portBindings[key] = []string{fmt.Sprintf("0.0.0.0:%d", port.HostPort)}
	}
	return portBindings
}

// convertVolumeMounts 转换卷挂载
func (s *AgentService) convertVolumeMounts(volumes []models.VolumeMount) []string {
	var binds []string
	for _, vol := range volumes {
		bind := fmt.Sprintf("%s:%s", vol.HostPath, vol.ContainerPath)
		if vol.ReadOnly {
			bind += ":ro"
		}
		binds = append(binds, bind)
	}
	return binds
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
	}

	return nil
}
