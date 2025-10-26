package operator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/boreas/internal/pkg/models"
)

// K8sClient Kubernetes Operator 客户端
type K8sClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewK8sClient 创建 K8s Operator 客户端
func NewK8sClient(baseURL string) *K8sClient {
	return &K8sClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Apply 应用部署
func (c *K8sClient) Apply(ctx context.Context, req *models.ApplyDeploymentRequest) (*models.ApplyDeploymentResponse, error) {
	// K8S 环境通常只部署单一版本，取第一个版本
	if len(req.Versions) == 0 {
		return nil, fmt.Errorf("no versions specified")
	}

	version := req.Versions[0]

	// 构建 K8S Apply 请求
	k8sReq := models.ApplyRequest{
		App:     req.App,
		Version: version.Version,
		Package: version.Package,
	}

	// 发送 HTTP 请求
	body, err := json.Marshal(k8sReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/apply", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var k8sResp models.ApplyResponse
	if err := json.NewDecoder(resp.Body).Decode(&k8sResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 转换为统一的响应格式
	return &models.ApplyDeploymentResponse{
		App:     k8sResp.App,
		Message: k8sResp.Message,
		Success: k8sResp.Success,
	}, nil
}

// GetApplicationStatus 获取应用状态
func (c *K8sClient) GetApplicationStatus(ctx context.Context, appName string) (*models.ApplicationStatusResponse, error) {
	// 发送 HTTP 请求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/v1/status/"+appName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var statusResp models.AppStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 转换为统一的响应格式
	return &models.ApplicationStatusResponse{
		App: statusResp.App,
		Healthy: models.HealthInfo{
			Level: statusResp.Healthy.Level,
		},
		Versions: []models.VersionStatus{
			{
				Version: statusResp.Version,
				Healthy: models.HealthInfo{
					Level: statusResp.Healthy.Level,
				},
				Nodes: []models.NodeStatus{}, // K8S 的节点信息需要进一步查询
			},
		},
	}, nil
}

// HealthCheck 健康检查
func (c *K8sClient) HealthCheck(ctx context.Context) error {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/v1/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// GetType 获取 Operator 类型
func (c *K8sClient) GetType() string {
	return "kubernetes"
}
