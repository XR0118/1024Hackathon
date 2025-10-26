package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/services/operator-pm/config"
)

type OperatorPMService struct {
	cfg    *config.Config
	client *http.Client
}

func NewOperatorPMService(cfg *config.Config) *OperatorPMService {
	return &OperatorPMService{
		cfg:    cfg,
		client: &http.Client{Timeout: time.Duration(cfg.PM.AgentTimeout) * time.Second},
	}
}

// CheckPMConnection 检查物理机连接状态
func (s *OperatorPMService) CheckPMConnection() error {
	// 检查所有配置的节点连接状态
	for nodeName := range s.cfg.PM.NodeToIP {
		agentURL, exists := s.cfg.GetAgentURL(nodeName)
		if !exists {
			return fmt.Errorf("node %s not found in IP mapping", nodeName)
		}

		resp, err := s.client.Get(agentURL + "/health")
		if err != nil {
			return fmt.Errorf("node %s (%s) is not reachable: %w", nodeName, agentURL, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("node %s (%s) returned status %d", nodeName, agentURL, resp.StatusCode)
		}
	}
	return nil
}

// ApplyDeployment 应用部署 - 核心API
func (s *OperatorPMService) ApplyDeployment(req *models.ApplyDeploymentRequest) (*models.ApplyDeploymentResponse, error) {
	// 1. 获取应用对应的节点列表
	nodes, exists := s.cfg.PM.AppToNodes[req.App]
	if !exists || len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes configured for application %s", req.App)
	}

	// 2. 为每个版本选择合适的节点进行部署
	var successCount int
	var totalCount int

	for _, version := range req.Versions {
		// 计算需要部署的节点数量
		nodeCount := int(float64(len(nodes)) * version.Percent)
		if nodeCount == 0 {
			nodeCount = 1 // 至少部署一个节点
		}

		// 选择前nodeCount个节点
		selectedNodes := nodes[:nodeCount]
		totalCount += len(selectedNodes)

		// 向选中的节点发送部署请求
		for _, nodeName := range selectedNodes {
			agentURL, exists := s.cfg.GetAgentURL(nodeName)
			if !exists {
				continue // 跳过无效节点
			}

			// 构建Agent的apply请求
			agentReq := map[string]interface{}{
				"app":     req.App,
				"version": version.Version,
				"package": version.Package,
			}

			// 发送请求到Agent
			if err := s.sendToAgent(agentURL+"/apply", agentReq); err != nil {
				// 记录错误但继续处理其他节点
				continue
			}

			successCount++
		}
	}

	// 3. 返回结果
	success := successCount > 0
	message := fmt.Sprintf("Deployed to %d/%d nodes", successCount, totalCount)
	if !success {
		message = "Failed to deploy to any nodes"
	}

	return &models.ApplyDeploymentResponse{
		App:     req.App,
		Message: message,
		Success: success,
	}, nil
}

// GetApplicationStatus 获取应用状态 - 核心API
func (s *OperatorPMService) GetApplicationStatus(appName string) (*models.ApplicationStatusResponse, error) {
	// 1. 获取应用对应的节点列表
	nodes, exists := s.cfg.PM.AppToNodes[appName]
	if !exists || len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes configured for application %s", appName)
	}

	// 2. 从所有节点收集状态信息
	var allNodeStatuses []models.NodeStatus
	var healthSum int = 0
	var healthCount int = 0

	for _, nodeName := range nodes {
		agentURL, exists := s.cfg.GetAgentURL(nodeName)
		if !exists {
			continue
		}

		// 获取节点状态
		nodeStatus, err := s.getNodeStatus(agentURL, nodeName)
		if err != nil {
			// 节点不可用，健康度为0
			allNodeStatuses = append(allNodeStatuses, models.NodeStatus{
				Node:    nodeName,
				Healthy: models.HealthInfo{Level: 0, Msg: "Node unreachable"},
			})
			healthSum += 0
			healthCount++
		} else {
			allNodeStatuses = append(allNodeStatuses, *nodeStatus)
			healthSum += nodeStatus.Healthy.Level
			healthCount++
		}
	}

	// 3. 计算平均健康度（加权平均，这里所有节点权重相同）
	var overallHealthy int = 0
	if healthCount > 0 {
		overallHealthy = healthSum / healthCount
	}

	// 4. 构建响应
	response := &models.ApplicationStatusResponse{
		App: appName,
		Healthy: models.HealthInfo{
			Level: overallHealthy,
			Msg:   fmt.Sprintf("Average health: %d%% (%d/%d nodes)", overallHealthy, healthCount, len(nodes)),
		},
		Versions: []models.VersionStatus{
			{
				Version: "latest", // 简化版本，实际应该从Agent获取
				Healthy: models.HealthInfo{Level: overallHealthy},
				Nodes:   allNodeStatuses,
			},
		},
	}

	return response, nil
}

// sendToAgent 发送请求到Agent
func (s *OperatorPMService) sendToAgent(url string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request to agent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("agent returned status %d", resp.StatusCode)
	}

	return nil
}

// getNodeStatus 获取节点状态
func (s *OperatorPMService) getNodeStatus(agentURL, nodeName string) (*models.NodeStatus, error) {
	resp, err := s.client.Get(agentURL + "/status")
	if err != nil {
		return nil, fmt.Errorf("failed to get status from node %s: %w", nodeName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("node %s returned status %d", nodeName, resp.StatusCode)
	}

	// 解析响应
	var status struct {
		Healthy models.HealthInfo `json:"healthy"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	return &models.NodeStatus{
		Node:    nodeName,
		Healthy: status.Healthy,
	}, nil
}
