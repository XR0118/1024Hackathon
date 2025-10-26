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

// PMClient Physical Machine Operator 客户端
type PMClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewPMClient 创建 PM Operator 客户端
func NewPMClient(baseURL string) *PMClient {
	return &PMClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Apply 应用部署
func (c *PMClient) Apply(ctx context.Context, req *models.ApplyDeploymentRequest) (*models.ApplyDeploymentResponse, error) {
	// PM 环境支持多版本灰度部署，直接使用请求
	body, err := json.Marshal(req)
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

	var applyResp models.ApplyDeploymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&applyResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &applyResp, nil
}

// GetApplicationStatus 获取应用状态
func (c *PMClient) GetApplicationStatus(ctx context.Context, appName string) (*models.ApplicationStatusResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/v1/status/"+appName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var statusResp models.ApplicationStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &statusResp, nil
}

// HealthCheck 健康检查
func (c *PMClient) HealthCheck(ctx context.Context) error {
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
func (c *PMClient) GetType() string {
	return "physical"
}
