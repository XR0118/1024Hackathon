package models

import (
	"time"
)

// GitHubEvent GitHub事件
type GitHubEvent struct {
	Type       string
	Repository string
	Payload    interface{}
}

// GitTag Git标签
type GitTag struct {
	Name       string
	Commit     string
	Repository string
	Pusher     string
	Message    string
	Timestamp  time.Time
}

// PullRequest 拉取请求
type PullRequest struct {
	Number      int
	Title       string
	Body        string
	MergeCommit string
	Repository  string
	BaseBranch  string
	HeadBranch  string
	MergedBy    string
	MergedAt    time.Time
}

// Release 发布
type Release struct {
	TagName      string
	Name         string
	Body         string
	Commit       string
	Repository   string
	Author       string
	PublishedAt  time.Time
	IsPrerelease bool
}

// ProcessResult 处理结果
type ProcessResult struct {
	Success              bool
	VersionCreated       bool
	Version              *Version
	DeploymentsTriggered []*Deployment
	Error                error
}

// DeploymentProgress 部署进度
type DeploymentProgress struct {
	DeploymentID   string
	Status         DeploymentStatus
	TotalTasks     int
	CompletedTasks int
	FailedTasks    int
	CurrentTask    *Task
}

// DeployRequest 部署请求
type DeployRequest struct {
	DeploymentID string
	Version      *Version
	Applications []*Application
	Environment  *Environment
	Config       map[string]string
}

// RollbackDeployRequest 回滚部署请求
type RollbackDeployRequest struct {
	DeploymentID       string
	TargetDeploymentID string
	Environment        *Environment
}

// DeployResult 部署结果
type DeployResult struct {
	Success      bool
	Message      string
	DeploymentID string
	Details      map[string]interface{}
}

// DeploymentInfo 部署信息
type DeploymentInfo struct {
	DeploymentID  string
	Status        string
	Replicas      int32
	ReadyReplicas int32
	UpdatedAt     time.Time
	Details       map[string]interface{}
}

// HealthCheckResult 健康检查结果
type HealthCheckResult struct {
	Healthy bool
	Message string
	Checks  []*HealthCheck
}

// HealthCheck 健康检查
type HealthCheck struct {
	Name    string
	Status  string
	Message string
}

// PodStatus Pod状态
type PodStatus struct {
	Name      string
	Phase     string
	Ready     bool
	Restarts  int32
	NodeName  string
	CreatedAt time.Time
}

// Artifact 制品
type Artifact struct {
	Name        string
	Version     string
	Path        string
	Size        int64
	Checksum    string
	ContentType string
}

// CommandResult 命令结果
type CommandResult struct {
	Host     string
	Success  bool
	Output   string
	Error    string
	ExitCode int
}

// ServiceStatus 服务状态
type ServiceStatus struct {
	Host    string
	Running bool
	Status  string
	PID     int
	Uptime  string
	Version string
}
