package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// Configurable via env:
// - BASE_URL (default: http://localhost:8080)
// - TOKEN (optional, adds Authorization: Bearer <TOKEN>)

func TestMasterServiceE2E(t *testing.T) {
	baseURL := getenvDefault("BASE_URL", "http://localhost:8080")
	token := "test"

	client := &http.Client{Timeout: 15 * time.Second}

	// Unique suffix to avoid collisions
	uniq := time.Now().UnixNano()
	appName := fmt.Sprintf("e2e-app-%d", uniq)
	envName := fmt.Sprintf("e2e-env-%d", uniq)

	// 1) Create Application
	appID := createApplication(t, client, baseURL, token, ApplicationCreateRequest{
		Name:       appName,
		Repository: "https://github.com/example/repo",
		Type:       "microservice",
		Config: map[string]string{
			"port":         "8080",
			"health_check": "/health",
		},
	})
	t.Logf("created application: id=%s name=%s", appID, appName)

	// 2) Create Environment
	envID := createEnvironment(t, client, baseURL, token, EnvironmentCreateRequest{
		Name:     envName,
		Type:     "kubernetes",
		IsActive: true,
		Config: map[string]string{
			"namespace": "production",
			"cluster":   "prod-cluster",
		},
	})
	t.Logf("created environment: id=%s name=%s", envID, envName)

	// 3) Create Version (with app_builds)
	versionID := createVersion(t, client, baseURL, token, VersionCreateRequest{
		GitTag:      fmt.Sprintf("v-e2e-%d", uniq),
		GitCommit:   "e2e-commit-sha",
		Repository:  "https://github.com/example/repo",
		Description: "E2E test version",
		AppBuild: []AppBuildItem{
			{AppID: appID, AppName: appName, DockerImage: fmt.Sprintf("registry.local/%s:%d", appName, uniq)},
		},
	})
	t.Logf("created version: id=%s", versionID)

	// 4) Create Deployment
	deploymentID := createDeployment(t, client, baseURL, token, DeploymentCreateRequest{
		VersionID:      versionID,
		MustInOrder:    []string{appID},
		EnvironmentID:  envID,
		ManualApproval: false,
		Strategy: []StrategyItem{
			{BatchSize: 10, BatchInterval: 10, CanaryRatio: 0, AutoRollback: true, ManualApprovalStatus: nil},
		},
	})
	t.Logf("created deployment: id=%s", deploymentID)

	// 5) Start Deployment
	did := startDeployment(t, client, baseURL, token, deploymentID)
	t.Logf("started deployment: id=%s deployment=%s", deploymentID, did)

	// 6) Poll Deployment Status
	finalStatus := pollDeploymentStatus(t, client, baseURL, token, deploymentID, 30, 5*time.Second)
	t.Logf("deployment %s final status: %s", deploymentID, finalStatus)

	// Assert terminal state (best-effort)
	switch finalStatus {
	case "success", "failed", "rolled_back", "cancelled":
		// OK; terminal
		break
	case "":
		// Unrecognized/empty; treat as non-terminal
		fallthrough
	default:
		t.Fatalf("deployment %s did not reach a recognized terminal state; status=%q", deploymentID, finalStatus)
	}

	// 3) Create Version (with app_builds)
	versionID2 := createVersion(t, client, baseURL, token, VersionCreateRequest{
		GitTag:      fmt.Sprintf("v-e2e-%d", uniq+10),
		GitCommit:   "e2e-commit-sha",
		Repository:  "https://github.com/example/repo",
		Description: "E2E test version",
		AppBuild: []AppBuildItem{
			{AppID: appID, AppName: appName, DockerImage: fmt.Sprintf("registry.local/%s:%d", appName, uniq+10)},
		},
	})
	t.Logf("created version2: id=%s", versionID)

	// 4) Create Deployment
	deploymentID2 := createDeployment(t, client, baseURL, token, DeploymentCreateRequest{
		VersionID:      versionID2,
		EnvironmentID:  envID,
		ManualApproval: false,
		Strategy: []StrategyItem{
			{BatchSize: 5, BatchInterval: 30, CanaryRatio: 0, AutoRollback: true},
		},
	})
	t.Logf("created deployment2: id=%s", deploymentID2)

	// 5) Start Deployment
	did2 := startDeployment(t, client, baseURL, token, deploymentID2)
	t.Logf("started deployment2: id=%s deployment=%s", deploymentID2, did2)

	t.Fatal()

}

// --- Request/Response payloads ---
type ApplicationCreateRequest struct {
	Name       string            `json:"name"`
	Repository string            `json:"repository"`
	Type       string            `json:"type"`
	Config     map[string]string `json:"config"`
}

type EnvironmentCreateRequest struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Config   map[string]string `json:"config"`
	IsActive bool              `json:"is_active"`
}

type AppBuildItem struct {
	AppID       string `json:"app_id"`
	AppName     string `json:"app_name"`
	DockerImage string `json:"docker_image"`
}

type VersionCreateRequest struct {
	GitTag      string         `json:"git_tag"`
	GitCommit   string         `json:"git_commit"`
	Repository  string         `json:"repository"`
	Description string         `json:"description"`
	AppBuild    []AppBuildItem `json:"app_build"`
}

type StrategyItem struct {
	BatchSize            int     `json:"batch_size"`
	BatchInterval        int     `json:"batch_interval"`
	CanaryRatio          float64 `json:"canary_ratio"`
	AutoRollback         bool    `json:"auto_rollback"`
	ManualApprovalStatus *string `json:"manual_approval_status"`
}

type DeploymentCreateRequest struct {
	VersionID      string         `json:"version_id"`
	MustInOrder    []string       `json:"must_in_order"`
	EnvironmentID  string         `json:"environment_id"`
	ManualApproval bool           `json:"manual_approval"`
	Strategy       []StrategyItem `json:"strategy"`
}

// --- Helpers ---
func getenvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func createApplication(t *testing.T, client *http.Client, baseURL, token string, req ApplicationCreateRequest) string {
	url := fmt.Sprintf("%s/api/v1/applications", baseURL)
	status, body := httpJSON(t, client, http.MethodPost, url, token, req)
	if status == http.StatusUnauthorized && token == "" {
		t.Skipf("service requires authorization; set TOKEN env to run E2E")
	}
	if status != http.StatusCreated && status != http.StatusOK {
		t.Fatalf("create application failed: status=%d body=%s", status, string(body))
	}
	if id := extractID(body); id != "" {
		return id
	}
	// Fallback: list and match by name
	listURL := fmt.Sprintf("%s/api/v1/applications?page=1&page_size=100", baseURL)
	lstStatus, lstBody := httpGet(t, client, listURL, token)
	if lstStatus != http.StatusOK {
		t.Fatalf("list applications failed: status=%d body=%s", lstStatus, string(lstBody))
	}
	return findIDByNameFromList(lstBody, "applications", req.Name)
}

func createEnvironment(t *testing.T, client *http.Client, baseURL, token string, req EnvironmentCreateRequest) string {
	url := fmt.Sprintf("%s/api/v1/environments", baseURL)
	status, body := httpJSON(t, client, http.MethodPost, url, token, req)
	if status != http.StatusCreated && status != http.StatusOK {
		t.Fatalf("create environment failed: status=%d body=%s", status, string(body))
	}
	if id := extractID(body); id != "" {
		return id
	}
	// Fallback: list and match by name
	listURL := fmt.Sprintf("%s/api/v1/environments?page=1&page_size=100", baseURL)
	lstStatus, lstBody := httpGet(t, client, listURL, token)
	if lstStatus != http.StatusOK {
		t.Fatalf("list environments failed: status=%d body=%s", lstStatus, string(lstBody))
	}
	return findIDByNameFromList(lstBody, "environments", req.Name)
}

func createVersion(t *testing.T, client *http.Client, baseURL, token string, req VersionCreateRequest) string {
	url := fmt.Sprintf("%s/api/v1/versions", baseURL)
	status, body := httpJSON(t, client, http.MethodPost, url, token, req)
	if status != http.StatusCreated && status != http.StatusOK {
		t.Fatalf("create version failed: status=%d body=%s", status, string(body))
	}
	id := extractID(body)
	if id == "" {
		t.Fatalf("create version returned no id: body=%s", string(body))
	}
	return id
}

func createDeployment(t *testing.T, client *http.Client, baseURL, token string, req DeploymentCreateRequest) string {
	url := fmt.Sprintf("%s/api/v1/deployments", baseURL)
	status, body := httpJSON(t, client, http.MethodPost, url, token, req)
	if status != http.StatusCreated && status != http.StatusOK {
		t.Fatalf("create deployment failed: status=%d body=%s", status, string(body))
	}
	id := extractID(body)
	if id == "" {
		t.Fatalf("create deployment returned no id: body=%s", string(body))
	}
	return id
}

func startDeployment(t *testing.T, client *http.Client, baseURL, token string, id string) string {
	url := fmt.Sprintf("%s/api/v1/deployments/%s/start", baseURL, id)
	status, body := httpJSON(t, client, http.MethodPost, url, token, nil)
	if status != http.StatusCreated && status != http.StatusOK {
		t.Fatalf("create deployment failed: status=%d body=%s", status, string(body))
	}
	id = extractID(body)
	if id == "" {
		t.Fatalf("create deployment returned no id: body=%s", string(body))
	}
	return id
}

func getDeploymentStatus(t *testing.T, client *http.Client, baseURL, token, deploymentID string) string {
	url := fmt.Sprintf("%s/api/v1/deployments/%s", baseURL, deploymentID)
	status, body := httpGet(t, client, url, token)
	if status != http.StatusOK {
		t.Fatalf("get deployment status failed: status=%d body=%s", status, string(body))
	}
	return string(body)
}

func pollDeploymentStatus(t *testing.T, client *http.Client, baseURL, token, deploymentID string, maxAttempts int, interval time.Duration) string {
	url := fmt.Sprintf("%s/api/v1/deployments/%s", baseURL, deploymentID)
	var statusStr string
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		st, body := httpGet(t, client, url, token)
		if st != http.StatusOK {
			t.Logf("get deployment attempt %d: status=%d body=%s", attempt, st, string(body))
		} else {
			statusStr = extractStatus(body)
			t.Logf("deployment %s status at attempt %d: %s", deploymentID, attempt, statusStr)
			if isTerminalStatus(statusStr) {
				return statusStr
			}
		}
		time.Sleep(interval)
	}
	return statusStr
}

func isTerminalStatus(s string) bool {
	switch s {
	case "success", "failed", "rolled_back", "cancelled":
		return true
	default:
		return false
	}
}

// --- HTTP helpers ---
func httpJSON(t *testing.T, client *http.Client, method, url, token string, payload any) (int, []byte) {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload error: %v", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("new request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("http do error: %v", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, respBody
}

func httpGet(t *testing.T, client *http.Client, url, token string) (int, []byte) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("new request error: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Logf("http get error: %v", err)
		return 0, nil
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, respBody
}

func extractData(body []byte) json.RawMessage {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(body, &m); err != nil {
		return nil
	}
	if data, ok := m["data"]; ok {
		return data
	}
	return nil
}

// --- JSON parsing helpers ---
func extractID(body []byte) string {
	body = extractData(body)

	var m2 map[string]interface{}
	if err := json.Unmarshal(body, &m2); err != nil {
		return ""
	}
	if id, ok := m2["id"].(string); ok {
		return id
	}
	return ""
}

func extractStatus(body []byte) string {
	body = extractData(body)
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return ""
	}
	if s, ok := m["status"].(string); ok {
		return s
	}
	return ""
}

func findIDByNameFromList(body []byte, key, targetName string) string {
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return ""
	}
	// Expect list under m[key] as []any
	if arr, ok := m[key].([]any); ok {
		for _, it := range arr {
			if obj, ok := it.(map[string]any); ok {
				nameVal, _ := obj["name"].(string)
				if nameVal == targetName {
					if id, ok := obj["id"].(string); ok {
						return id
					}
				}
			}
		}
	}
	// Fallback: maybe API returns array directly
	if arrAny, ok := m["data"].([]any); ok { // common pattern
		for _, it := range arrAny {
			if obj, ok := it.(map[string]any); ok {
				nameVal, _ := obj["name"].(string)
				if nameVal == targetName {
					if id, ok := obj["id"].(string); ok {
						return id
					}
				}
			}
		}
	}
	return ""
}
