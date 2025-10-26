package interfaces

import (
	"context"

	"github.com/boreas/internal/pkg/models"
)

// VersionService 版本服务接口
type VersionService interface {
	CreateVersion(ctx context.Context, req *models.CreateVersionRequest) (*models.Version, error)
	GetVersionList(ctx context.Context, req *models.ListVersionsRequest) (*models.VersionListResponse, error)
	GetVersion(ctx context.Context, version string) (*models.Version, error) // 通过版本号查询
	DeleteVersion(ctx context.Context, version string) error                 // 通过版本号删除
}

// ApplicationService 应用服务接口
type ApplicationService interface {
	CreateApplication(ctx context.Context, req *models.CreateApplicationRequest) (*models.Application, error)
	GetApplicationList(ctx context.Context, req *models.ListApplicationsRequest) (*models.ApplicationListResponse, error)
	GetApplication(ctx context.Context, id string) (*models.Application, error)
	GetApplicationByName(ctx context.Context, name string) (*models.Application, error)
	UpdateApplication(ctx context.Context, id string, req *models.UpdateApplicationRequest) (*models.Application, error)
	DeleteApplication(ctx context.Context, id string) error
	// 版本相关
	GetApplicationVersionsSummary(ctx context.Context, appName string) (*models.ApplicationVersionsSummaryResponse, error)
	GetApplicationVersionsDetail(ctx context.Context, appName string) (*models.ApplicationVersionsDetailResponse, error)
}

// EnvironmentService 环境服务接口
type EnvironmentService interface {
	CreateEnvironment(ctx context.Context, req *models.CreateEnvironmentRequest) (*models.Environment, error)
	GetEnvironmentList(ctx context.Context, req *models.ListEnvironmentsRequest) (*models.EnvironmentListResponse, error)
	GetEnvironment(ctx context.Context, id string) (*models.Environment, error)
	UpdateEnvironment(ctx context.Context, id string, req *models.UpdateEnvironmentRequest) (*models.Environment, error)
	DeleteEnvironment(ctx context.Context, id string) error
}

// DeploymentService 部署服务接口
type DeploymentService interface {
	CreateDeployment(ctx context.Context, req *models.CreateDeploymentRequest) (*models.Deployment, error)
	GetDeploymentList(ctx context.Context, req *models.ListDeploymentsRequest) (*models.DeploymentListResponse, error)
	GetDeployment(ctx context.Context, id string) (*models.Deployment, error)
	StartDeployment(ctx context.Context, id string) (*models.Deployment, error)
	RollbackDeployment(ctx context.Context, id string, req *models.RollbackRequest) error
}

// TaskService 任务服务接口
type TaskService interface {
	GetTaskList(ctx context.Context, req *models.ListTasksRequest) (*models.TaskListResponse, error)
	GetTask(ctx context.Context, id string) (*models.Task, error)
	RetryTask(ctx context.Context, id string) (*models.Task, error)
}

// WebhookService Webhook服务接口
type WebhookService interface {
	HandleGitHubWebhook(ctx context.Context, event string, payload []byte) (*models.WebhookResponse, error)
	VerifySignature(payload []byte, signature string) error
}

// TaskScheduler 任务调度器接口
type TaskScheduler interface {
	ScheduleTask(ctx context.Context, task *models.Task) error
	GetNextTask(ctx context.Context) (*models.Task, error)
	UpdateTaskStatus(ctx context.Context, taskID string, status models.TaskStatus, result string) error
	GetTasksByDeployment(ctx context.Context, deploymentID string) ([]*models.Task, error)
}

// DeploymentManager 部署管理器接口
type DeploymentManager interface {
	ProcessWebhookEvent(ctx context.Context, event *models.GitHubEvent) (*models.ProcessResult, error)
	CreateVersionFromTag(ctx context.Context, tag *models.GitTag) (*models.Version, error)
	CreateVersionFromPR(ctx context.Context, pr *models.PullRequest) (*models.Version, error)
	CreateVersionFromRelease(ctx context.Context, release *models.Release) (*models.Version, error)
	TriggerAutoDeployments(ctx context.Context, versionID string) ([]*models.Deployment, error)
	StartDeployment(ctx context.Context, deploymentID string) error
	CancelDeployment(ctx context.Context, deploymentID string) error
	UpdateDeploymentStatus(ctx context.Context, deploymentID string, status models.DeploymentStatus, errorMsg string) error
	GetDeploymentProgress(ctx context.Context, deploymentID string) (*models.DeploymentProgress, error)
}

// Operator 运行环境操作器统一接口
// Master 通过此接口与各类 Operator（K8S、PM、Mock）进行通信
type Operator interface {
	// Apply 应用部署
	Apply(ctx context.Context, req *models.ApplyDeploymentRequest) (*models.ApplyDeploymentResponse, error)

	// GetApplicationStatus 获取应用状态
	GetApplicationStatus(ctx context.Context, appName string) (*models.ApplicationStatusResponse, error)

	// HealthCheck 健康检查
	HealthCheck(ctx context.Context) error

	// GetType 获取 Operator 类型 (kubernetes/physical/mock)
	GetType() string
}

// OperatorManager Operator 管理器接口
// 负责管理所有运行环境的 Operator 客户端
type OperatorManager interface {
	// RegisterOperator 注册 Operator
	RegisterOperator(environmentID string, operator Operator)

	// GetOperator 获取指定环境的 Operator
	GetOperator(environmentID string) (Operator, error)

	// GetOperatorByEnvironment 根据环境对象获取 Operator
	GetOperatorByEnvironment(env *models.Environment) (Operator, error)

	// ApplyDeployment 在指定环境应用部署
	ApplyDeployment(ctx context.Context, environmentID string, req *models.ApplyDeploymentRequest) (*models.ApplyDeploymentResponse, error)

	// GetApplicationStatus 获取指定环境中应用的状态
	GetApplicationStatus(ctx context.Context, environmentID string, appName string) (*models.ApplicationStatusResponse, error)

	// HealthCheckAll 检查所有 Operator 的健康状态
	HealthCheckAll(ctx context.Context) map[string]error

	// RemoveOperator 移除指定环境的 Operator
	RemoveOperator(environmentID string)

	// ListOperators 列出所有已注册的 Operator
	ListOperators() []string
}

// Deployer 部署器基础接口（已废弃，请使用 Operator）
// Deprecated: 使用 Operator 接口替代
type Deployer interface {
	Deploy(ctx context.Context, req *models.DeployRequest) (*models.DeployResult, error)
	Rollback(ctx context.Context, req *models.RollbackDeployRequest) (*models.DeployResult, error)
	GetDeploymentInfo(ctx context.Context, deploymentID string) (*models.DeploymentInfo, error)
	HealthCheck(ctx context.Context, deploymentID string) (*models.HealthCheckResult, error)
}

// KubernetesDeployer Kubernetes部署器接口
type KubernetesDeployer interface {
	Deployer
	ApplyManifest(ctx context.Context, namespace string, manifest []byte) error
	GetPodStatus(ctx context.Context, namespace, selector string) ([]*models.PodStatus, error)
	ScaleDeployment(ctx context.Context, namespace, name string, replicas int32) error
	GetLogs(ctx context.Context, namespace, podName string, lines int) (string, error)
	DeleteResources(ctx context.Context, namespace string, labels map[string]string) error
}

// PhysicalDeployer 物理机部署器接口
type PhysicalDeployer interface {
	Deployer
	UploadArtifact(ctx context.Context, hosts []string, artifact *models.Artifact) error
	ExecuteCommand(ctx context.Context, hosts []string, command string) (*models.CommandResult, error)
	RestartService(ctx context.Context, hosts []string, serviceName string) error
	CheckServiceStatus(ctx context.Context, hosts []string, serviceName string) ([]*models.ServiceStatus, error)
	CleanupOldVersions(ctx context.Context, hosts []string, keepVersions int) error
}
