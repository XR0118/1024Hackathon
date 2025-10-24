package workflow

import (
	"github.com/XR0118/1024Hackathon/backend/model"
	"github.com/XR0118/1024Hackathon/backend/service"
)

type DeploymentController interface {
	HandleGitWebhook(event *service.GitEvent) error
	CreateVersion(event *service.GitEvent) (*model.Version, error)
	PlanDeployments(version *model.Version) ([]*model.Deployment, error)
	ExecuteDeployment(deployment *model.Deployment) error
	RollbackDeployment(deployment *model.Deployment) error
	GetControllerStatus() (*ControllerStatus, error)
}

type ControllerStatus struct {
	Status            string `json:"status"`
	ActiveDeployments int    `json:"active_deployments"`
	PendingApprovals  int    `json:"pending_approvals"`
	LastSync          string `json:"last_sync"`
}

type VersionGenerator interface {
	CreateVersion(event *service.GitEvent) (*model.Version, error)
	DetectRevert(event *service.GitEvent) bool
	FindParentVersion(repository, commit string) string
}

type DeploymentPlanner interface {
	PlanDeployments(version *model.Version) ([]*model.Deployment, error)
	FindApplications(repository string) ([]*model.Application, error)
	GetEnvironmentSequence(repository string) []string
	GetDeployStrategy(env *model.TargetEnvironment, version *model.Version) model.DeployStrategy
}

type DeploymentExecutor interface {
	Execute(deployment *model.Deployment) error
	DeployBatch(apps []*model.Application, version *model.Version, env *model.TargetEnvironment) []DeployResult
	PreCheckBatch(apps []*model.Application, env *model.TargetEnvironment) error
	HealthCheckBatch(apps []*model.Application, env *model.TargetEnvironment, strategy model.DeployStrategy) bool
	WaitForApproval(deployment *model.Deployment) bool
	UpdateProgress(deployment *model.Deployment, results []DeployResult)
}

type DeployResult struct {
	AppID    string `json:"app_id"`
	EnvID    string `json:"env_id"`
	Status   string `json:"status"`
	ErrorMsg string `json:"error_msg,omitempty"`
}

type StatusMonitor interface {
	MonitorDeployment(deploymentID string) error
	SendNotification(deployment *model.Deployment, eventType string) error
	RecordHistory(deployment *model.Deployment) error
}
