package models

import (
	"encoding/json"
	"time"
)

// Version 版本信息
type Version struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	GitTag      string    `json:"git_tag" gorm:"uniqueIndex;not null"`
	GitCommit   string    `json:"git_commit" gorm:"not null"`
	Repository  string    `json:"repository" gorm:"not null"`
	CreatedBy   string    `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	Description string    `json:"description"`

	AppBuilds []AppBuild `json:"app_build,omitempty" gorm:"type:jsonb"`
}

type AppBuild struct {
	AppID       string `json:"app_id"`
	AppName     string `json:"app_name"`
	DockerImage string `json:"docker_image"`
}

// Application 应用信息
type Application struct {
	ID         string            `json:"id" gorm:"primaryKey"`
	Name       string            `json:"name" gorm:"uniqueIndex;not null"`
	Repository string            `json:"repository" gorm:"not null"`
	Type       string            `json:"type" gorm:"not null"` // microservice, monolith
	Config     map[string]string `json:"config" gorm:"type:jsonb"`
	CreatedAt  time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
}

type BuildConfig struct {
	Dockerfile string             `json:"dockerfile"`
	BuildArgs  map[string]*string `json:"build_args"`
	Context    string             `json:"context"`
}

func (a *Application) GetBuildConfig() *BuildConfig {
	var ret = BuildConfig{}
	if s, ok := a.Config["build_config"]; ok {
		_ = json.Unmarshal([]byte(s), &ret)
	}
	return &ret
}

// Environment 环境信息
type Environment struct {
	ID        string            `json:"id" gorm:"primaryKey"`
	Name      string            `json:"name" gorm:"uniqueIndex;not null"`
	Type      string            `json:"type" gorm:"not null"` // kubernetes, physical
	Config    map[string]string `json:"config" gorm:"type:jsonb"`
	IsActive  bool              `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
}

// DeploymentStatus 部署状态
type DeploymentStatus string

const (
	DeploymentStatusPending    DeploymentStatus = "pending"
	DeploymentStatusRunning    DeploymentStatus = "running"
	DeploymentStatusSuccess    DeploymentStatus = "success"
	DeploymentStatusFailed     DeploymentStatus = "failed"
	DeploymentStatusRolledBack DeploymentStatus = "rolled_back"
	DeploymentStatusCancelled  DeploymentStatus = "cancelled"
)

// Deployment 部署信息
type Deployment struct {
	ID             string           `json:"id" gorm:"primaryKey"`
	VersionID      string           `json:"version_id" gorm:"not null"`
	ApplicationIDs []string         `json:"application_ids" gorm:"type:jsonb"`
	EnvironmentID  string           `json:"environment_id" gorm:"not null"`
	Status         DeploymentStatus `json:"status" gorm:"default:'pending'"`
	CreatedBy      string           `json:"created_by" gorm:"not null"`
	CreatedAt      time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	StartedAt      *time.Time       `json:"started_at,omitempty"`
	CompletedAt    *time.Time       `json:"completed_at,omitempty"`
	ErrorMessage   string           `json:"error_message,omitempty"`

	// 关联关系
	Version      Version       `json:"version,omitempty" gorm:"foreignKey:VersionID"`
	Environment  Environment   `json:"environment,omitempty" gorm:"foreignKey:EnvironmentID"`
	Applications []Application `json:"applications,omitempty" gorm:"many2many:deployment_applications;"`
}

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending TaskStatus = "pending"
	TaskStatusRunning TaskStatus = "running"
	TaskStatusSuccess TaskStatus = "success"
	TaskStatusFailed  TaskStatus = "failed"
)

// Task 任务信息
type Task struct {
	ID           string     `json:"id" gorm:"primaryKey"`
	DeploymentID string     `json:"deployment_id" gorm:"not null"`
	Type         string     `json:"type" gorm:"not null"` // build, test, deploy, health_check
	Status       TaskStatus `json:"status" gorm:"default:'pending'"`
	Payload      string     `json:"payload" gorm:"type:text"`
	Result       string     `json:"result,omitempty" gorm:"type:text"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`

	// 关联关系
	Deployment Deployment `json:"deployment,omitempty" gorm:"foreignKey:DeploymentID"`
}

// WorkflowStatus 工作流状态
type WorkflowStatus string

const (
	WorkflowStatusPending   WorkflowStatus = "pending"
	WorkflowStatusRunning   WorkflowStatus = "running"
	WorkflowStatusSuccess   WorkflowStatus = "success"
	WorkflowStatusFailed    WorkflowStatus = "failed"
	WorkflowStatusCancelled WorkflowStatus = "cancelled"
)

// Workflow 工作流信息
type Workflow struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	DeploymentID string         `json:"deployment_id" gorm:"not null"`
	Status       WorkflowStatus `json:"status" gorm:"default:'pending'"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`

	// 关联关系
	Deployment Deployment `json:"deployment,omitempty" gorm:"foreignKey:DeploymentID"`
	Tasks      []Task     `json:"tasks,omitempty" gorm:"foreignKey:DeploymentID"`
}

// 请求和响应类型

// CreateVersionRequest 创建版本请求
type CreateVersionRequest struct {
	GitTag      string `json:"git_tag" binding:"required"`
	GitCommit   string `json:"git_commit" binding:"required"`
	Repository  string `json:"repository" binding:"required"`
	Description string `json:"description"`

	AppBuilds []AppBuild `json:"app_build,omitempty"`
}

// ListVersionsRequest 版本列表请求
type ListVersionsRequest struct {
	Repository string `form:"repository"`
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=100"`
}

// VersionListResponse 版本列表响应
type VersionListResponse struct {
	Versions []*Version `json:"versions"`
	Total    int        `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// CreateApplicationRequest 创建应用请求
type CreateApplicationRequest struct {
	Name       string            `json:"name" binding:"required"`
	Repository string            `json:"repository" binding:"required"`
	Type       string            `json:"type" binding:"required"`
	Config     map[string]string `json:"config"`
}

// UpdateApplicationRequest 更新应用请求
type UpdateApplicationRequest struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

// ListApplicationsRequest 应用列表请求
type ListApplicationsRequest struct {
	Repository string `form:"repository"`
	Type       string `form:"type"`
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=100"`
}

// ApplicationListResponse 应用列表响应
type ApplicationListResponse struct {
	Applications []*Application `json:"applications"`
	Total        int            `json:"total"`
	Page         int            `json:"page"`
	PageSize     int            `json:"page_size"`
}

// CreateEnvironmentRequest 创建环境请求
type CreateEnvironmentRequest struct {
	Name     string            `json:"name" binding:"required"`
	Type     string            `json:"type" binding:"required"`
	Config   map[string]string `json:"config" binding:"required"`
	IsActive bool              `json:"is_active"`
}

// UpdateEnvironmentRequest 更新环境请求
type UpdateEnvironmentRequest struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Config   map[string]string `json:"config"`
	IsActive *bool             `json:"is_active"`
}

// ListEnvironmentsRequest 环境列表请求
type ListEnvironmentsRequest struct {
	Type     string `form:"type"`
	IsActive *bool  `form:"is_active"`
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}

// EnvironmentListResponse 环境列表响应
type EnvironmentListResponse struct {
	Environments []*Environment `json:"environments"`
	Total        int            `json:"total"`
	Page         int            `json:"page"`
	PageSize     int            `json:"page_size"`
}

// CreateDeploymentRequest 创建部署请求
type CreateDeploymentRequest struct {
	VersionID      string   `json:"version_id" binding:"required"`
	ApplicationIDs []string `json:"application_ids" binding:"required"`
	EnvironmentID  string   `json:"environment_id" binding:"required"`
}

// ListDeploymentsRequest 部署列表请求
type ListDeploymentsRequest struct {
	Status        string `form:"status"`
	EnvironmentID string `form:"environment_id"`
	VersionID     string `form:"version_id"`
	Page          int    `form:"page" binding:"min=1"`
	PageSize      int    `form:"page_size" binding:"min=1,max=100"`
}

// DeploymentListResponse 部署列表响应
type DeploymentListResponse struct {
	Deployments []*Deployment `json:"deployments"`
	Total       int           `json:"total"`
	Page        int           `json:"page"`
	PageSize    int           `json:"page_size"`
}

// RollbackRequest 回滚请求
type RollbackRequest struct {
	TargetVersionID string `json:"target_version_id" binding:"required"`
}

// ListTasksRequest 任务列表请求
type ListTasksRequest struct {
	DeploymentID string `form:"deployment_id"`
	Status       string `form:"status"`
	Type         string `form:"type"`
	Page         int    `form:"page" binding:"min=1"`
	PageSize     int    `form:"page_size" binding:"min=1,max=100"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Tasks    []*Task `json:"tasks"`
	Total    int     `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}

// WebhookResponse Webhook响应
type WebhookResponse struct {
	Message              string                `json:"message"`
	VersionCreated       bool                  `json:"version_created,omitempty"`
	VersionID            string                `json:"version_id,omitempty"`
	AutoTag              string                `json:"auto_tag,omitempty"`
	DeploymentsTriggered []DeploymentReference `json:"deployments_triggered,omitempty"`
}

// DeploymentReference 部署引用
type DeploymentReference struct {
	DeploymentID  string           `json:"deployment_id"`
	EnvironmentID string           `json:"environment_id"`
	Status        DeploymentStatus `json:"status"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail 错误详情
type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// 分页响应
type PaginationResponse struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// 新增类型定义

// DeploymentResult 部署结果
type DeploymentResult struct {
	ID        string           `json:"id"`
	Status    DeploymentStatus `json:"status"`
	Message   string           `json:"message"`
	Timestamp time.Time        `json:"timestamp"`
}

// DeploymentLogs 部署日志
type DeploymentLogs struct {
	ID    string   `json:"id"`
	Logs  []string `json:"logs"`
	Level string   `json:"level"`
}

// DeploymentLog 部署日志记录
type DeploymentLog struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	DeploymentID string    `json:"deployment_id" gorm:"not null"`
	Level        string    `json:"level" gorm:"not null"` // info, warn, error
	Message      string    `json:"message" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// DeploymentStatusInfo 部署状态信息
type DeploymentStatusInfo struct {
	ID      string           `json:"id"`
	Status  DeploymentStatus `json:"status"`
	Message string           `json:"message"`
}
