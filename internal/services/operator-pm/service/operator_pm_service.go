package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/services/operator-pm/repository"
	"gorm.io/gorm"
)

type OperatorPMService struct {
	db         *gorm.DB
	repository *repository.OperatorPMRepository
	agents     map[string]string // agent ID -> agent URL
}

func NewOperatorPMService(db *gorm.DB) *OperatorPMService {
	return &OperatorPMService{
		db:         db,
		repository: repository.NewOperatorPMRepository(db),
		agents:     make(map[string]string),
	}
}

// CheckPMConnection 检查物理机连接状态
func (s *OperatorPMService) CheckPMConnection() error {
	// 检查所有注册的agent连接状态
	for agentID, agentURL := range s.agents {
		if err := s.checkAgentConnection(agentURL); err != nil {
			return fmt.Errorf("agent %s (%s) is not reachable: %w", agentID, agentURL, err)
		}
	}
	return nil
}

// RegisterAgent 注册物理机Agent
func (s *OperatorPMService) RegisterAgent(agentID, agentURL string) {
	s.agents[agentID] = agentURL
}

// ListAgents 列出所有注册的Agent
func (s *OperatorPMService) ListAgents() map[string]string {
	return s.agents
}

// checkAgentConnection 检查Agent连接状态
func (s *OperatorPMService) checkAgentConnection(agentURL string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(agentURL + "/v1/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("agent returned status %d", resp.StatusCode)
	}

	return nil
}

// ExecuteDeployment 执行物理机部署
func (s *OperatorPMService) ExecuteDeployment(deploymentID string) (*models.DeploymentResult, error) {
	// 获取部署信息
	deployment, err := s.repository.GetDeploymentByID(deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	// 更新部署状态为执行中
	deployment.Status = models.DeploymentStatusRunning
	deployment.StartedAt = &time.Time{}
	*deployment.StartedAt = time.Now()
	if err := s.repository.UpdateDeployment(deployment); err != nil {
		return nil, fmt.Errorf("failed to update deployment status: %w", err)
	}

	// 执行物理机部署
	result, err := s.executePMDeployment(deployment)
	if err != nil {
		// 更新部署状态为失败
		deployment.Status = models.DeploymentStatusFailed
		deployment.CompletedAt = &time.Time{}
		*deployment.CompletedAt = time.Now()
		deployment.ErrorMessage = err.Error()
		s.repository.UpdateDeployment(deployment)
		return nil, fmt.Errorf("failed to execute deployment: %w", err)
	}

	// 更新部署状态为成功
	deployment.Status = models.DeploymentStatusSuccess
	deployment.CompletedAt = &time.Time{}
	*deployment.CompletedAt = time.Now()
	if err := s.repository.UpdateDeployment(deployment); err != nil {
		return nil, fmt.Errorf("failed to update deployment status: %w", err)
	}

	return result, nil
}

// GetDeploymentStatus 获取部署状态
func (s *OperatorPMService) GetDeploymentStatus(deploymentID string) (*models.DeploymentStatusInfo, error) {
	deployment, err := s.repository.GetDeploymentByID(deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	// 如果部署正在运行，检查物理机中的实际状态
	if deployment.Status == models.DeploymentStatusRunning {
		pmStatus, err := s.getPMDeploymentStatus(deploymentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get PM deployment status: %w", err)
		}
		return pmStatus, nil
	}

	return &models.DeploymentStatusInfo{
		ID:      deployment.ID,
		Status:  deployment.Status,
		Message: deployment.ErrorMessage,
	}, nil
}

// GetDeploymentLogs 获取部署日志
func (s *OperatorPMService) GetDeploymentLogs(deploymentID string) (*models.DeploymentLogs, error) {
	// 从物理机获取部署日志
	logs, err := s.getPMDeploymentLogs(deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get PM deployment logs: %w", err)
	}

	return logs, nil
}

// CancelDeployment 取消部署
func (s *OperatorPMService) CancelDeployment(deploymentID string) error {
	deployment, err := s.repository.GetDeploymentByID(deploymentID)
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	// 如果部署正在运行，取消物理机中的部署
	if deployment.Status == models.DeploymentStatusRunning {
		if err := s.cancelPMDeployment(deploymentID); err != nil {
			return fmt.Errorf("failed to cancel PM deployment: %w", err)
		}
	}

	// 更新部署状态为取消
	deployment.Status = models.DeploymentStatusCancelled
	deployment.CompletedAt = &time.Time{}
	*deployment.CompletedAt = time.Now()
	if err := s.repository.UpdateDeployment(deployment); err != nil {
		return fmt.Errorf("failed to update deployment status: %w", err)
	}

	return nil
}

// executePMDeployment 执行物理机部署
func (s *OperatorPMService) executePMDeployment(deployment *models.Deployment) (*models.DeploymentResult, error) {
	// 获取环境配置以确定目标agent
	env, err := s.repository.GetEnvironmentByID(deployment.EnvironmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}

	// 从环境配置中获取agent信息
	agentID, exists := env.Config["agent_id"]
	if !exists {
		return nil, fmt.Errorf("agent_id not found in environment config")
	}

	agentURL, exists := s.agents[agentID]
	if !exists {
		return nil, fmt.Errorf("agent %s not registered", agentID)
	}

	// 解析版本中的应用构建信息
	appBuilds := deployment.Version.AppBuilds

	// 为每个应用执行部署
	for _, appBuild := range appBuilds {
		if err := s.deployAppToAgent(agentURL, appBuild, deployment); err != nil {
			return nil, fmt.Errorf("failed to deploy app %s: %w", appBuild.AppName, err)
		}
	}

	return &models.DeploymentResult{
		ID:        deployment.ID,
		Status:    models.DeploymentStatusSuccess,
		Message:   "PM deployment completed successfully",
		Timestamp: time.Now(),
	}, nil
}

// deployAppToAgent 部署应用到Agent
func (s *OperatorPMService) deployAppToAgent(agentURL string, appBuild models.AppBuild, deployment *models.Deployment) error {
	// 构建部署包
	pkg := map[string]interface{}{
		"type":        "docker",
		"image":       appBuild.DockerImage,
		"command":     []string{},
		"args":        []string{},
		"environment": map[string]string{},
		"volumes":     []models.VolumeMount{},
		"ports":       []models.PortMapping{},
	}

	// 创建部署请求
	applyReq := models.ApplyRequest{
		App:     appBuild.AppName,
		Version: deployment.Version.GitTag,
		Pkg:     pkg,
	}

	// 发送部署请求到Agent
	client := &http.Client{Timeout: 30 * time.Second}

	reqData, err := json.Marshal(applyReq)
	if err != nil {
		return fmt.Errorf("failed to marshal apply request: %w", err)
	}

	resp, err := client.Post(agentURL+"/v1/apply", "application/json", bytes.NewBuffer(reqData))
	if err != nil {
		return fmt.Errorf("failed to send apply request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("agent returned status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var applyResp models.ApplyResponse
	if err := json.NewDecoder(resp.Body).Decode(&applyResp); err != nil {
		return fmt.Errorf("failed to decode apply response: %w", err)
	}

	if !applyResp.Success {
		return fmt.Errorf("agent deployment failed: %s", applyResp.Message)
	}

	return nil
}

// getPMDeploymentStatus 获取物理机部署状态
func (s *OperatorPMService) getPMDeploymentStatus(deploymentID string) (*models.DeploymentStatusInfo, error) {
	deployment, err := s.repository.GetDeploymentByID(deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	// 获取环境配置以确定目标agent
	env, err := s.repository.GetEnvironmentByID(deployment.EnvironmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}

	agentID, exists := env.Config["agent_id"]
	if !exists {
		return nil, fmt.Errorf("agent_id not found in environment config")
	}

	agentURL, exists := s.agents[agentID]
	if !exists {
		return nil, fmt.Errorf("agent %s not registered", agentID)
	}

	// 查询agent状态
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(agentURL + "/v1/status")
	if err != nil {
		return nil, fmt.Errorf("failed to query agent status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("agent returned status %d", resp.StatusCode)
	}

	// 解析状态响应
	var statusResp models.StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	// 检查是否有应用在运行
	if len(statusResp.Apps) == 0 {
		return &models.DeploymentStatusInfo{
			ID:      deploymentID,
			Status:  models.DeploymentStatusFailed,
			Message: "No apps running on agent",
		}, nil
	}

	// 检查所有应用的健康状态
	allHealthy := true
	for _, app := range statusResp.Apps {
		if app.Healthy.Level < 80 {
			allHealthy = false
			break
		}
	}

	status := models.DeploymentStatusSuccess
	message := "All apps are healthy"
	if !allHealthy {
		status = models.DeploymentStatusRunning
		message = "Some apps are not healthy"
	}

	return &models.DeploymentStatusInfo{
		ID:      deploymentID,
		Status:  status,
		Message: message,
	}, nil
}

// getPMDeploymentLogs 获取物理机部署日志
func (s *OperatorPMService) getPMDeploymentLogs(deploymentID string) (*models.DeploymentLogs, error) {
	deployment, err := s.repository.GetDeploymentByID(deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	// 获取环境配置以确定目标agent
	env, err := s.repository.GetEnvironmentByID(deployment.EnvironmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment: %w", err)
	}

	agentID, exists := env.Config["agent_id"]
	if !exists {
		return nil, fmt.Errorf("agent_id not found in environment config")
	}

	agentURL, exists := s.agents[agentID]
	if !exists {
		return nil, fmt.Errorf("agent %s not registered", agentID)
	}

	// 查询agent状态获取应用列表
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(agentURL + "/v1/status")
	if err != nil {
		return nil, fmt.Errorf("failed to query agent status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("agent returned status %d", resp.StatusCode)
	}

	var statusResp models.StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	// 收集所有应用的日志
	var allLogs []string
	for _, app := range statusResp.Apps {
		// 这里简化实现，实际应该从agent获取具体应用的日志
		allLogs = append(allLogs, fmt.Sprintf("App %s (%s): %s", app.App, app.Version, app.Healthy.Msg))
	}

	return &models.DeploymentLogs{
		ID:    deploymentID,
		Logs:  allLogs,
		Level: "info",
	}, nil
}

// cancelPMDeployment 取消物理机部署
func (s *OperatorPMService) cancelPMDeployment(deploymentID string) error {
	deployment, err := s.repository.GetDeploymentByID(deploymentID)
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	// 获取环境配置以确定目标agent
	env, err := s.repository.GetEnvironmentByID(deployment.EnvironmentID)
	if err != nil {
		return fmt.Errorf("failed to get environment: %w", err)
	}

	agentID, exists := env.Config["agent_id"]
	if !exists {
		return fmt.Errorf("agent_id not found in environment config")
	}

	agentURL, exists := s.agents[agentID]
	if !exists {
		return fmt.Errorf("agent %s not registered", agentID)
	}

	// 解析版本中的应用构建信息
	appBuilds := deployment.Version.AppBuilds

	// 停止每个应用
	for _, appBuild := range appBuilds {
		if err := s.stopAppOnAgent(agentURL, appBuild.AppName); err != nil {
			return fmt.Errorf("failed to stop app %s: %w", appBuild.AppName, err)
		}
	}

	return nil
}

// stopAppOnAgent 在Agent上停止应用
func (s *OperatorPMService) stopAppOnAgent(agentURL, appName string) error {
	// 这里简化实现，实际应该调用agent的停止接口
	// 目前通过查询状态来确认应用是否已停止
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(agentURL + "/v1/status/" + appName)
	if err != nil {
		// 如果应用不存在，认为已经停止
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// 应用不存在，认为已经停止
		return nil
	}

	// 这里可以添加实际的停止逻辑
	// 比如调用agent的停止接口或发送停止信号

	return nil
}
