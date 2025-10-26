package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/boreas/internal/pkg/models"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type ClientOption func(*Client)

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

func NewClient(baseURL string, opts ...ClientOption) *Client {
	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (c *Client) Apply(ctx context.Context, app, version string, pkg models.DeploymentPackage) (*models.ApplyResponse, error) {
	reqBody := map[string]interface{}{
		"app":     app,
		"version": version,
		"pkg":     pkg,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/apply", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("server returned error code %d: %s", response.Code, response.Message)
	}

	var applyResp models.ApplyResponse
	if err := json.Unmarshal(response.Data, &applyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal apply response: %w", err)
	}

	return &applyResp, nil
}

func (c *Client) GetStatus(ctx context.Context, appName string) (*models.StatusResponse, error) {
	url := c.baseURL + "/v1/status"
	if appName != "" {
		url = url + "?app=" + appName
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("server returned error code %d: %s", response.Code, response.Message)
	}

	var statusResp models.StatusResponse
	if err := json.Unmarshal(response.Data, &statusResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status response: %w", err)
	}

	return &statusResp, nil
}
