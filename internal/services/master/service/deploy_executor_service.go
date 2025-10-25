package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/services/master/mock"
)

type DeployExecutor interface {
	Apply(ctx context.Context) error
}

type DeployClient interface {
	Apply(ctx context.Context, app string, version string, pkg models.DeploymentPackage) (models.ApplyResponse, error)
	AppStatus(ctx context.Context, app string) ([]models.AgentAppStatus, error)
}

type SimpleDeployExecutor struct {
	task           models.Task
	deploymentRepo interfaces.DeploymentRepository
	versionRepo    interfaces.VersionRepository
	client         DeployClient
}

var mockClient = mock.NewMockDeploymentClient()

func NewSimpleDeployExecutor(task models.Task,
	deploymentRepo interfaces.DeploymentRepository, versionRepo interfaces.VersionRepository) *SimpleDeployExecutor {
	return &SimpleDeployExecutor{
		task:           task,
		deploymentRepo: deploymentRepo,
		client:         mockClient,
	}
}

func (e *SimpleDeployExecutor) Apply(ctx context.Context) (models.TaskStatus, error) {
	deployment, err := e.deploymentRepo.GetByID(ctx, e.task.DeploymentID)
	if err != nil {
		return models.TaskStatusFailed, err
	}

	if deployment.Status == models.DeploymentStatusCancelled {
		return models.TaskStatusCancelled, nil
	}
	if deployment.Status == models.DeploymentStatusFailed {
		return models.TaskStatusFailed, nil
	}

	appStatus, _ := e.client.AppStatus(ctx, e.task.AppID)
	total := 0
	versionMap := make(map[string]int)
	for _, s := range appStatus {
		total += s.Replicas
		versionMap[s.Version] += s.Replicas
	}

	pkg := models.DeploymentPackage{}
	_ = json.Unmarshal([]byte(e.task.Payload), &pkg)

	if deployment.Status == models.DeploymentStatusRolledBack {
		// 回滚到当前 app 的上一个版本
		prevVersion, err_ := e.versionRepo.GetPreviousByVersionAndApp(ctx, e.task.Deployment.VersionID, e.task.AppID)
		if err_ != nil {
			return models.TaskStatusFailed, err_
		}
		versionMap[prevVersion.ID] += versionMap[e.task.Deployment.VersionID]
		versionMap[e.task.Deployment.VersionID] = 0

		pkg.Replicas = versionMap[prevVersion.ID]
		_, err = e.client.Apply(ctx, e.task.AppID, prevVersion.ID, pkg)
		if err != nil {
			return models.TaskStatusFailed, err
		}
		pkg.Replicas = 0
		_, err = e.client.Apply(ctx, e.task.AppID, e.task.Deployment.VersionID, pkg)
		if err != nil {
			return models.TaskStatusFailed, err
		}
		return models.TaskStatusRolledBack, nil
	}

	if total == 0 {
		for _, s := range e.task.Deployment.GetStrategy() {
			total += s.BatchSize
		}
		pkg.Replicas = total
		_, err = e.client.Apply(ctx, e.task.AppID, e.task.Deployment.ID, pkg)
		if err != nil {
			return models.TaskStatusFailed, nil
		}
		return models.TaskStatusRunning, nil
	}

	if versionMap[e.task.Deployment.ID] == total {
		return models.TaskStatusSuccess, nil
	}

	ss := e.task.Deployment.GetStrategy()
	step, idx := func() (models.DeploySteps, int) {
		replicas := versionMap[e.task.Deployment.VersionID]
		raito := float64(replicas) / float64(total)
		i := 0
		for j, s := range ss {
			if raito < s.CanaryRatio || replicas < s.BatchSize {
				return s, j
			}
			i = j
		}
		s := ss[i]
		return s, i
	}()

	var unhealthy bool
	var lastUpdate time.Time
	for _, status := range appStatus {
		if status.Version == e.task.Deployment.ID {
			lastUpdate = status.Updated
			if status.Healthy.Level < 80 {
				unhealthy = true
				break
			}
		}
	}
	if unhealthy && step.AutoRollback {
		deployment.Status = models.DeploymentStatusRolledBack
		_ = e.deploymentRepo.Update(ctx, deployment)
		return models.TaskStatusRolledBack, nil
	}

	if time.Since(lastUpdate).Seconds() < float64(step.BatchInterval) ||
		step.ManualApprovalStatus != nil && !*step.ManualApprovalStatus {
		return models.TaskStatusRunning, nil
	}

	if idx == len(ss)-1 {
		// 没有下一个阶段, 等待全量结束
		// 判断 live = expect
		return models.TaskStatusRunning, nil
	}

	ns := ss[idx+1]
	if ns.BatchSize > 0 {
		pkg.Replicas = ns.BatchSize
	} else {
		pkg.Replicas = int(float64(total) * ns.CanaryRatio)
	}
	_, err = e.client.Apply(ctx, e.task.AppID, e.task.Deployment.ID, pkg)
	if err != nil {
		return models.TaskStatusFailed, err
	}

	return models.TaskStatusRunning, nil
}
