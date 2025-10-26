package interfaces

import (
	"context"

	"github.com/boreas/internal/pkg/models"
)

// VersionRepository 版本仓库接口
type VersionRepository interface {
	Create(ctx context.Context, version *models.Version) error
	GetByID(ctx context.Context, id string) (*models.Version, error)
	GetByVersion(ctx context.Context, version string) (*models.Version, error) // 通过版本号查询
	List(ctx context.Context, filter *models.VersionFilter) ([]*models.Version, int, error)
	// Get the previous version that contains the given app_id and was created before the target version
	GetPreviousByVersionAndApp(ctx context.Context, targetVersionID string, appID string) (*models.Version, error)
	Delete(ctx context.Context, id string) error
}

// ApplicationRepository 应用仓库接口
type ApplicationRepository interface {
	Create(ctx context.Context, app *models.Application) error
	GetByID(ctx context.Context, id string) (*models.Application, error)
	GetByName(ctx context.Context, name string) (*models.Application, error)
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
