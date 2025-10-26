package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// GitInfo Git 信息
type GitInfo struct {
	Tag        string `json:"tag"`
	Commit     string `json:"commit"`
	Repository string `json:"repository"`
}

// Version 版本信息
type Version struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Version     string         `json:"version" gorm:"uniqueIndex;not null"` // 版本号，作为唯一标识
	GitTag      string         `json:"git_tag" gorm:"not null"`
	GitCommit   string         `json:"git_commit" gorm:"not null"`
	Repository  string         `json:"repository" gorm:"not null"`
	Status      string         `json:"status" gorm:"not null;default:'normal'"` // normal, revert
	CreatedBy   string         `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	Description string         `json:"description"`
	AppBuilds   datatypes.JSON `json:"app_builds,omitempty" gorm:"type:jsonb"`
}

// GetGitInfo 获取 Git 信息
func (v *Version) GetGitInfo() GitInfo {
	return GitInfo{
		Tag:        v.GitTag,
		Commit:     v.GitCommit,
		Repository: v.Repository,
	}
}

// GetAppBuilds 获取应用构建信息
func (v *Version) GetAppBuilds() []AppBuild {
	var builds []AppBuild
	_ = json.Unmarshal(v.AppBuilds, &builds)
	return builds
}

// SetAppBuilds 设置应用构建信息
func (v *Version) SetAppBuilds(builds []AppBuild) error {
	data, err := json.Marshal(builds)
	if err != nil {
		return err
	}
	v.AppBuilds = data
	return nil
}

// AppBuild 应用构建信息
type AppBuild struct {
	AppID       string `json:"app_id"`
	AppName     string `json:"app_name"`     // 应用名称（唯一标识）
	DockerImage string `json:"docker_image"` // Docker 镜像地址
}

// Application 应用信息
type Application struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"uniqueIndex;not null"`
	Description  string         `json:"description"` // 应用描述
	Repository   string         `json:"repository" gorm:"not null"`
	Type         string         `json:"type" gorm:"not null"` // microservice, monolith
	Config       datatypes.JSON `json:"config" gorm:"type:jsonb"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	Environments []Environment  `json:"environments,omitempty" gorm:"many2many:application_environments;"` // 关联的环境列表
}

// ApplicationEnvironment 应用-环境映射表
type ApplicationEnvironment struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	ApplicationID string    `json:"application_id" gorm:"not null"`
	EnvironmentID string    `json:"environment_id" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type BuildConfig struct {
	Dockerfile string             `json:"dockerfile"`
	BuildArgs  map[string]*string `json:"build_args"`
	Context    string             `json:"context"`
}

func (a *Application) GetBuildConfig() *BuildConfig {
	var ret = BuildConfig{}
	var config map[string]string
	b, _ := a.Config.MarshalJSON()
	_ = json.Unmarshal(b, &config)
	if s, ok := config["build_config"]; ok {
		_ = json.Unmarshal([]byte(s), &ret)
	}
	return &ret
}

// Environment 环境信息
type Environment struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null"`
	Type      string         `json:"type" gorm:"not null"` // kubernetes, physical
	Config    datatypes.JSON `json:"config" gorm:"type:jsonb"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// DeploymentStatus 部署状态
// DeploymentStatus 部署状态（面向前端）
type DeploymentStatus string

const (
	DeploymentStatusPending   DeploymentStatus = "pending"   // 待开始
	DeploymentStatusRunning   DeploymentStatus = "running"   // 运行中
	DeploymentStatusPaused    DeploymentStatus = "paused"    // 已暂停
	DeploymentStatusCompleted DeploymentStatus = "completed" // 已完成
)

func (d DeploymentStatus) IsFinished() bool {
	return d == DeploymentStatusCompleted
}

// Deployment 部署信息
type Deployment struct {
	ID            string           `json:"id" gorm:"primaryKey"`
	VersionID     string           `json:"version_id" gorm:"not null"`
	MustInOrder   datatypes.JSON   `json:"must_in_order" gorm:"type:jsonb"` // version 中包含的 app 必须按照这个顺序部署（应用名称数组）
	EnvironmentID string           `json:"environment_id" gorm:"not null"`
	Status        DeploymentStatus `json:"status" gorm:"default:'pending'"` // 状态：pending/running/paused/completed
	CreatedBy     string           `json:"created_by" gorm:"not null"`
	CreatedAt     time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	StartedAt     *time.Time       `json:"started_at,omitempty"`
	CompletedAt   *time.Time       `json:"completed_at,omitempty"`
	ErrorMessage  string           `json:"error_message,omitempty"`

	Rollback bool           `json:"rollback"`
	Strategy datatypes.JSON `json:"strategy" gorm:"type:jsonb"`

	// 关联关系
	Version     Version     `json:"version,omitempty" gorm:"foreignKey:VersionID"`
	Environment Environment `json:"environment,omitempty" gorm:"foreignKey:EnvironmentID"`
	Tasks       []Task      `json:"tasks,omitempty" gorm:"foreignKey:DeploymentID"`
}

func (d *Deployment) GetMustInOrder() []string {
	var ret []string
	_ = json.Unmarshal(d.MustInOrder, &ret)
	return ret
}

func (d *Deployment) GetStrategy() []DeploySteps {
	var ret []DeploySteps
	_ = json.Unmarshal(d.Strategy, &ret)
	return ret
}

type DeploySteps struct {
	BatchSize            int     `json:"batch_size"`
	BatchInterval        int     `json:"batch_interval"`
	CanaryRatio          float64 `json:"canary_ratio"`
	AutoRollback         bool    `json:"auto_rollback"`
	ManualApprovalStatus *bool   `json:"manual_approval_status"`
}

// TaskType 任务类型
type TaskType string

const (
	TaskTypeBuild    TaskType = "build"    // 构建
	TaskTypeSleep    TaskType = "sleep"    // 等待
	TaskTypeDeploy   TaskType = "deploy"   // 部署
	TaskTypeTest     TaskType = "test"     // 测试
	TaskTypeApproval TaskType = "approval" // 复核/审批
)

// TaskStep workflow 执行中的状态
type TaskStep string

const (
	TaskStepPending   TaskStep = "pending"   // 待执行
	TaskStepRunning   TaskStep = "running"   // 执行中
	TaskStepBlocked   TaskStep = "blocked"   // 被阻塞（依赖未完成）
	TaskStepCompleted TaskStep = "completed" // 已完成
)

// TaskStatus 最终结果状态
type TaskStatus string

const (
	TaskStatusPending TaskStatus = "pending" // 待执行
	TaskStatusRunning TaskStatus = "running" // 执行中
	TaskStatusSuccess TaskStatus = "success" // 成功
	TaskStatusFailed  TaskStatus = "failed"  // 失败
)

func (s TaskStatus) IsFinished() bool {
	return s == TaskStatusSuccess || s == TaskStatusFailed
}

func (s TaskStep) IsCompleted() bool {
	return s == TaskStepCompleted
}

// 业务特定状态应该在 Payload/Result 中定义
// 例如：
// - Approval 任务的 "waiting_approval" 状态 -> 在 Result 中定义为 {"approval_state": "waiting"}

// Task 任务信息
type Task struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	DeploymentID string         `json:"deployment_id" gorm:"not null"`
	AppID        string         `json:"app_id" gorm:"not null"`              // 应用名称
	Name         string         `json:"name" gorm:"not null"`                // 任务名称
	Type         TaskType       `json:"type" gorm:"not null"`                // 任务类型：build/sleep/deploy/test/approval
	Step         TaskStep       `json:"step" gorm:"default:'pending'"`       // workflow 执行中的状态
	Status       TaskStatus     `json:"status" gorm:"default:'pending'"`     // 最终结果状态
	Dependencies datatypes.JSON `json:"dependencies" gorm:"type:jsonb"`      // 上游依赖任务ID列表（DAG结构）
	Payload      datatypes.JSON `json:"payload,omitempty" gorm:"type:jsonb"` // 任务参数（通用结构体）
	Result       datatypes.JSON `json:"result,omitempty" gorm:"type:jsonb"`  // 任务结果（通用结构体）
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	StartedAt    *time.Time     `json:"started_at,omitempty"`
	CompletedAt  *time.Time     `json:"completed_at,omitempty"`

	// 关联关系
	Deployment  Deployment  `json:"deployment,omitempty" gorm:"foreignKey:DeploymentID"`
	Application Application `json:"application,omitempty" gorm:"foreignKey:AppID"`
}

// GetDependencies 获取依赖列表
func (t *Task) GetDependencies() []string {
	var deps []string
	_ = json.Unmarshal(t.Dependencies, &deps)
	return deps
}

// SetDependencies 设置依赖列表
func (t *Task) SetDependencies(deps []string) error {
	data, err := json.Marshal(deps)
	if err != nil {
		return err
	}
	t.Dependencies = data
	return nil
}

// GetPayload 获取 Payload 并反序列化到指定类型
func (t *Task) GetPayload(v interface{}) error {
	if t.Payload == nil || len(t.Payload) == 0 {
		return nil
	}
	return json.Unmarshal(t.Payload, v)
}

// SetPayload 设置 Payload
func (t *Task) SetPayload(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	t.Payload = data
	return nil
}

// GetResult 获取 Result 并反序列化到指定类型
func (t *Task) GetResult(v interface{}) error {
	if t.Result == nil || len(t.Result) == 0 {
		return nil
	}
	return json.Unmarshal(t.Result, v)
}

// SetResult 设置 Result
func (t *Task) SetResult(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	t.Result = data
	return nil
}

// ==================== Task Payload 和 Result 类型定义 ====================

// BuildTaskPayload Build 任务参数
type BuildTaskPayload struct {
	Dockerfile  string            `json:"dockerfile,omitempty"`
	Context     string            `json:"context,omitempty"`
	BuildArgs   map[string]string `json:"build_args,omitempty"`
	TargetImage string            `json:"target_image"`
}

// BuildTaskResult Build 任务结果
type BuildTaskResult struct {
	Image         string   `json:"image"`
	ImageID       string   `json:"image_id"`
	Size          int64    `json:"size"`
	BuildDuration int      `json:"build_duration"` // 秒
	Logs          []string `json:"logs,omitempty"`
}

// SleepTaskPayload Sleep 任务参数
type SleepTaskPayload struct {
	Duration int    `json:"duration"` // 等待时间（秒）
	Reason   string `json:"reason,omitempty"`
}

// SleepTaskResult Sleep 任务结果
type SleepTaskResult struct {
	ActualDuration int `json:"actual_duration"` // 实际等待时间（秒）
}

// DeployTaskPayload Deploy 任务参数
type DeployTaskPayload struct {
	Image       string                 `json:"image"`
	Replicas    int                    `json:"replicas"`
	Strategy    string                 `json:"strategy"` // rolling/blue-green/canary
	CanaryRatio int                    `json:"canary_ratio,omitempty"`
	HealthCheck *HealthCheckConfig     `json:"health_check,omitempty"`
	ExtraConfig map[string]interface{} `json:"extra_config,omitempty"` // 额外配置
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Endpoint string `json:"endpoint"`
	Interval int    `json:"interval"` // 秒
	Timeout  int    `json:"timeout"`  // 秒
}

// DeployTaskResult Deploy 任务结果
type DeployTaskResult struct {
	DeployedInstances int      `json:"deployed_instances"`
	HealthyInstances  int      `json:"healthy_instances"`
	RolloutDuration   int      `json:"rollout_duration"` // 秒
	Endpoints         []string `json:"endpoints,omitempty"`
}

// TestTaskPayload Test 任务参数
type TestTaskPayload struct {
	TestSuite   string   `json:"test_suite"`
	TestCases   []string `json:"test_cases,omitempty"`
	Environment string   `json:"environment"`
	Timeout     int      `json:"timeout,omitempty"` // 秒
}

// TestTaskResult Test 任务结果
type TestTaskResult struct {
	Passed   int                 `json:"passed"`
	Failed   int                 `json:"failed"`
	Skipped  int                 `json:"skipped"`
	Duration int                 `json:"duration"` // 秒
	Coverage float64             `json:"coverage,omitempty"`
	Failures []TestFailureDetail `json:"failures,omitempty"`
}

// TestFailureDetail 测试失败详情
type TestFailureDetail struct {
	Test    string `json:"test"`
	Message string `json:"message"`
}

// ApprovalTaskPayload Approval 任务参数
type ApprovalTaskPayload struct {
	Note              string   `json:"note"`                         // 审批说明
	RequiredApprovers []string `json:"required_approvers,omitempty"` // 需要审批的人员
	AutoApproveAfter  int      `json:"auto_approve_after,omitempty"` // 自动批准的超时时间（秒）
}

// ApprovalTaskResult Approval 任务结果
type ApprovalTaskResult struct {
	Approved        bool   `json:"approved"`
	Approver        string `json:"approver,omitempty"`    // 审批人
	ApprovedAt      string `json:"approved_at,omitempty"` // 审批时间
	RejectionReason string `json:"rejection_reason,omitempty"`
	ApprovalState   string `json:"approval_state,omitempty"` // waiting/approved/rejected/timeout
}

// 请求和响应类型

// CreateVersionRequest 创建版本请求
type CreateVersionRequest struct {
	Version     string     `json:"version" binding:"required"` // 版本号
	GitTag      string     `json:"git_tag" binding:"required"`
	GitCommit   string     `json:"git_commit" binding:"required"`
	Repository  string     `json:"repository" binding:"required"`
	Description string     `json:"description"`
	AppBuilds   []AppBuild `json:"app_builds,omitempty"`
}

// ListVersionsRequest 版本列表请求
type ListVersionsRequest struct {
	Repository string `form:"repository"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}

// VersionListResponse 版本列表响应
type VersionListResponse struct {
	Versions []*Version `json:"versions"`
	Total    int        `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// RollbackVersionRequest 回滚版本请求
type RollbackVersionRequest struct {
	Reason string `json:"reason" binding:"required"` // 回滚原因
}

// CreateApplicationRequest 创建应用请求
type CreateApplicationRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	Repository  string            `json:"repository" binding:"required"`
	Type        string            `json:"type" binding:"required"`
	Config      map[string]string `json:"config"`
}

// UpdateApplicationRequest 更新应用请求
type UpdateApplicationRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"`
	Config      map[string]string `json:"config"`
}

// ListApplicationsRequest 应用列表请求
type ListApplicationsRequest struct {
	Repository string `form:"repository"`
	Type       string `form:"type"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}

// ApplicationListResponse 应用列表响应
type ApplicationListResponse struct {
	Applications []*Application `json:"applications"`
	Total        int            `json:"total"`
	Page         int            `json:"page"`
	PageSize     int            `json:"page_size"`
}

// ApplicationVersionInfo 应用版本信息（从 Operator 查询）
type ApplicationVersionInfo struct {
	Version       string                `json:"version"`
	Status        string                `json:"status"` // normal, revert
	Health        int                   `json:"health"`
	Coverage      int                   `json:"coverage"`
	LastUpdatedAt string                `json:"last_updated_at"`
	Nodes         []ApplicationNodeInfo `json:"nodes,omitempty"`
}

// ApplicationNodeInfo 应用节点信息
type ApplicationNodeInfo struct {
	Name          string `json:"name"`
	Health        int    `json:"health"`
	LastUpdatedAt string `json:"last_updated_at"`
}

// ApplicationVersionsResponse 应用版本列表响应
// Deprecated: 请使用 ApplicationVersionsDetailResponse
type ApplicationVersionsResponse struct {
	ApplicationID string                   `json:"application_id"`
	Name          string                   `json:"name"`
	Versions      []ApplicationVersionInfo `json:"versions"`
}

// ==================== 版本概要相关结构 ====================

// VersionSummary 版本概要信息（只包含核心运行时指标）
type VersionSummary struct {
	Version         string     `json:"version"`          // 版本号
	Status          string     `json:"status"`           // 版本状态: normal, revert
	Healthy         HealthInfo `json:"healthy"`          // 健康度 (0-100)
	CoveragePercent float64    `json:"coverage_percent"` // 覆盖度百分比 (0-100)
}

// ApplicationVersionsSummaryResponse 应用版本概要响应
type ApplicationVersionsSummaryResponse struct {
	ApplicationID   string           `json:"application_id"`
	ApplicationName string           `json:"application_name"`
	Versions        []VersionSummary `json:"versions"`
}

// ==================== 版本详情相关结构 ====================

// VersionInstance 版本实例信息
type VersionInstance struct {
	NodeName      string     `json:"node_name"`       // 节点名称
	Healthy       HealthInfo `json:"healthy"`         // 健康度 (0-100)
	Status        string     `json:"status"`          // 实例状态
	LastUpdatedAt time.Time  `json:"last_updated_at"` // 最后更新时间
}

// EnvironmentVersionDetail 环境下的版本详细信息
type EnvironmentVersionDetail struct {
	Version       string            `json:"version"`         // 版本号
	Status        string            `json:"status"`          // 版本状态
	GitTag        string            `json:"git_tag"`         // Git 标签
	GitCommit     string            `json:"git_commit"`      // Git 提交哈希
	Instances     []VersionInstance `json:"instances"`       // 实例列表
	Healthy       HealthInfo        `json:"healthy"`         // 该版本在此环境的健康度 (0-100)
	Coverage      int               `json:"coverage"`        // 该版本在此环境的覆盖率(%)
	LastUpdatedAt time.Time         `json:"last_updated_at"` // 最后更新时间
}

// EnvironmentVersions 环境维度的版本信息
type EnvironmentVersions struct {
	Environment Environment                `json:"environment"` // 环境信息
	Versions    []EnvironmentVersionDetail `json:"versions"`    // 该环境下的版本列表
}

// ApplicationVersionsDetailResponse 应用版本详情响应（按环境组织）
type ApplicationVersionsDetailResponse struct {
	ApplicationID   string                `json:"application_id"`
	ApplicationName string                `json:"application_name"`
	Environments    []EnvironmentVersions `json:"environments"` // 环境列表
}

// VersionCoverageResponse 版本覆盖率响应（累积覆盖率）
type VersionCoverageResponse struct {
	ApplicationID       string                       `json:"application_id"`
	ApplicationName     string                       `json:"application_name"`
	TargetVersion       string                       `json:"target_version"`       // 查询的目标版本
	TotalEnvironments   int                          `json:"total_environments"`   // 总环境数
	CoveredEnvironments int                          `json:"covered_environments"` // 已覆盖的环境数（版本 >= 目标版本）
	CoveragePercent     float64                      `json:"coverage_percent"`     // 覆盖率百分比
	Environments        []EnvironmentVersionCoverage `json:"environments"`         // 各环境的覆盖详情
}

// EnvironmentVersionCoverage 环境版本覆盖详情
type EnvironmentVersionCoverage struct {
	Environment         Environment            `json:"environment"`          // 环境信息
	IsCovered           bool                   `json:"is_covered"`           // 是否被覆盖（有版本 >= 目标版本）
	CurrentVersion      string                 `json:"current_version"`      // 当前最高版本
	TotalInstances      int                    `json:"total_instances"`      // 总实例数
	CoveredInstances    int                    `json:"covered_instances"`    // 被覆盖的实例数（版本 >= 目标版本）
	CoveragePercent     float64                `json:"coverage_percent"`     // 该环境内的覆盖率
	VersionDistribution []VersionInstanceCount `json:"version_distribution"` // 版本分布详情
}

// VersionInstanceCount 版本实例计数
type VersionInstanceCount struct {
	Version       string `json:"version"`        // 版本号
	InstanceCount int    `json:"instance_count"` // 实例数
	IsCovered     bool   `json:"is_covered"`     // 是否被目标版本覆盖
}

// CreateEnvironmentRequest 创建环境请求
type CreateEnvironmentRequest struct {
	Name     string            `json:"name" binding:"required"`
	Type     string            `json:"type" binding:"required"`
	Config   map[string]string `json:"config"` // 可选配置
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
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
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
	VersionID     string   `json:"version_id" binding:"required"`
	MustInOrder   []string `json:"must_in_order"`
	EnvironmentID string   `json:"environment_id" binding:"required"`

	Strategy       []DeploySteps `json:"strategy" binding:"required"`
	ManualApproval bool          `json:"manual_approval"`
}

// ListDeploymentsRequest 部署列表请求
type ListDeploymentsRequest struct {
	Status        string `form:"status"`
	EnvironmentID string `form:"environment_id"`
	VersionID     string `form:"version_id"`
	Page          int    `form:"page"`
	PageSize      int    `form:"page_size"`
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
	TargetVersionID string `json:"target_version_id"`
}

// ListTasksRequest 任务列表请求
type ListTasksRequest struct {
	DeploymentID string `form:"deployment_id"`
	Status       string `form:"status"`
	Type         string `form:"type"`
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
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

// 核心API数据结构

// ApplyDeploymentRequest 应用部署请求（Operator PM API）
type ApplyDeploymentRequest struct {
	App      string              `json:"app" binding:"required"`
	Versions []VersionDeployment `json:"versions" binding:"required"`
}

// VersionDeployment 版本部署信息
type VersionDeployment struct {
	Version string            `json:"version" binding:"required"`
	Percent float64           `json:"percent" binding:"required,min=0,max=1"`
	Package DeploymentPackage `json:"package" binding:"required"`
}

// ApplyDeploymentResponse 应用部署响应（Operator PM API）
type ApplyDeploymentResponse struct {
	App     string `json:"app"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// ApplicationStatusResponse 应用状态响应（Operator PM API）
type ApplicationStatusResponse struct {
	App      string          `json:"app"`
	Healthy  HealthInfo      `json:"healthy"`
	Versions []VersionStatus `json:"versions"`
}

// VersionStatus 版本状态
type VersionStatus struct {
	Version string       `json:"version"`
	Healthy HealthInfo   `json:"healthy"`
	Nodes   []NodeStatus `json:"nodes"` // 该版本的节点列表，上层根据节点数计算覆盖率
}

// NodeStatus 节点状态
type NodeStatus struct {
	Node    string     `json:"node"`
	Healthy HealthInfo `json:"healthy"`
}

// HealthInfo 健康状态信息
type HealthInfo struct {
	Level int    `json:"level"`         // 0-100
	Msg   string `json:"msg,omitempty"` // 健康状态描述信息
}
