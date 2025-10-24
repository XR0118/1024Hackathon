package interfaces

import (
	"context"

	"github.com/boreas/internal/models"
)

// VersionRepository 版本仓库接口
type VersionRepository interface {
	Create(ctx context.Context, version *models.Version) error
	GetByID(ctx context.Context, id string) (*models.Version, error)
	List(ctx context.Context, filter *models.VersionFilter) ([]*models.Version, int, error)
	Delete(ctx context.Context, id string) error
}

// ApplicationRepository 应用仓库接口
type ApplicationRepository interface {
	Create(ctx context.Context, app *models.Application) error
	GetByID(ctx context.Context, id string) (*models.Application, error)
	List(ctx context.Context, filter *models.ApplicationFilter) ([]*models.Application, int, error)
	Update(ctx context.Context, app *models.Application) error
	Delete(ctx context.Context, id string) error
}

// EnvironmentRepository 环境仓库接口
type EnvironmentRepository interface {
	Create(ctx context.Context, env *models.Environment) error
	GetByID(ctx context.Context, id string) (*models.Environment, error)
	List(ctx context.Context, filter *models.EnvironmentFilter) ([]*models.Environment, int, error)
	Update(ctx context.Context, env *models.Environment) error
	Delete(ctx context.Context, id string) error
}

// DeploymentRepository 部署仓库接口
type DeploymentRepository interface {
	Create(ctx context.Context, deployment *models.Deployment) error
	GetByID(ctx context.Context, id string) (*models.Deployment, error)
	List(ctx context.Context, filter *models.DeploymentFilter) ([]*models.Deployment, int, error)
	Update(ctx context.Context, deployment *models.Deployment) error
}

// TaskRepository 任务仓库接口
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	List(ctx context.Context, filter *models.TaskFilter) ([]*models.Task, int, error)
	Update(ctx context.Context, task *models.Task) error
	GetByDeploymentID(ctx context.Context, deploymentID string) ([]*models.Task, error)
}

// WorkflowRepository 工作流仓库接口
type WorkflowRepository interface {
	Create(ctx context.Context, workflow *models.Workflow) error
	GetByID(ctx context.Context, id string) (*models.Workflow, error)
	GetByDeploymentID(ctx context.Context, deploymentID string) (*models.Workflow, error)
	Update(ctx context.Context, workflow *models.Workflow) error
}
