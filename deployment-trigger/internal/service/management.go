package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ManagementClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewManagementClient(baseURL string) *ManagementClient {
	return &ManagementClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (m *ManagementClient) GetApplication(ctx context.Context, appName string) (*Application, error) {
	url := fmt.Sprintf("%s/applications?name=%s", m.baseURL, appName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var listResp struct {
		Applications []Application `json:"applications"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(listResp.Applications) == 0 {
		return nil, fmt.Errorf("application not found: %s", appName)
	}

	return &listResp.Applications[0], nil
}

func (m *ManagementClient) CreateVersion(ctx context.Context, req *CreateVersionRequest) (*Version, error) {
	url := fmt.Sprintf("%s/versions", m.baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(respBody))
	}

	var version Version
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &version, nil
}

type Application struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Repository  string    `json:"repository"`
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	Config      AppConfig `json:"config"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AppConfig struct {
	BuildConfig   *BuildConfig   `json:"build_config"`
	RuntimeConfig *RuntimeConfig `json:"runtime_config"`
	HealthCheck   *HealthCheck   `json:"health_check"`
}

type RuntimeConfig struct {
	Port      int               `json:"port"`
	Env       map[string]string `json:"env"`
	Resources *Resources        `json:"resources"`
}

type Resources struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type HealthCheck struct {
	Path         string `json:"path"`
	Port         int    `json:"port"`
	InitialDelay int    `json:"initial_delay"`
	Period       int    `json:"period"`
}

type Version struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	GitTag     string            `json:"git_tag"`
	GitCommit  string            `json:"git_commit"`
	Repository string            `json:"repository"`
	CreatedBy  string            `json:"created_by"`
	CreatedAt  time.Time         `json:"created_at"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Status     string            `json:"status"`
	Apps       []AppBuild        `json:"apps,omitempty"`
}

type AppBuild struct {
	AppID       string `json:"app_id"`
	AppName     string `json:"app_name"`
	DockerImage string `json:"docker_image"`
}

type CreateVersionRequest struct {
	Name       string            `json:"name"`
	GitTag     string            `json:"git_tag"`
	GitCommit  string            `json:"git_commit"`
	Repository string            `json:"repository"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Apps       []AppBuild        `json:"apps,omitempty"`
}
