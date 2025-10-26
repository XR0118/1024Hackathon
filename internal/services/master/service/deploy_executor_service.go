package service

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/boreas/internal/interfaces"
	"github.com/boreas/internal/pkg/logger"
	"github.com/boreas/internal/pkg/models"
	"github.com/boreas/internal/services/master/mock"
	"go.uber.org/zap"
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

func NewSimpleDeployExecutor(task models.Task,
	deploymentRepo interfaces.DeploymentRepository, versionRepo interfaces.VersionRepository) *SimpleDeployExecutor {
	return &SimpleDeployExecutor{
		task:           task,
		deploymentRepo: deploymentRepo,
		client:         mock.MockAgent,
		versionRepo:    versionRepo,
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
		if s.Replicas == 0 {
			continue
		}
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
		_, err = e.client.Apply(ctx, e.task.AppID, e.task.Deployment.VersionID, pkg)
		if err != nil {
			return models.TaskStatusFailed, nil
		}
		return models.TaskStatusRunning, nil
	}

	var versions []*models.Version
	for vid := range versionMap {
		ver, _ := e.versionRepo.GetByID(ctx, vid)
		versions = append(versions, ver)
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreatedAt.Before(versions[j].CreatedAt)
	})

	// 没有比自己更旧的版本则说明发布成功
	if versions[0].ID == e.task.Deployment.VersionID {
		return models.TaskStatusSuccess, nil
	}

	ss := e.task.Deployment.GetStrategy()
	step, idx := func() (models.DeploySteps, int) {
		replicas := versionMap[e.task.Deployment.VersionID]
		raito := float64(replicas) / float64(total)
		for j, s := range ss {
			if raito < s.CanaryRatio || replicas < s.BatchSize {
				return s, j
			}
		}
		// 返回最后一个 step 只用配置
		return ss[len(ss)-1], len(ss)
	}()

	var unhealthy bool
	var lastUpdate time.Time
	for _, status := range appStatus {
		if status.Version == e.task.Deployment.VersionID {
			lastUpdate = status.Updated
			if status.Healthy.Level < 80 {
				unhealthy = true
				break
			}
		}
	}
	if unhealthy && step.AutoRollback {
		logger.GetLogger().Warn("deployment unhealthy, auto rollback", zap.String("deployment_id", deployment.ID))
		deployment.Status = models.DeploymentStatusRolledBack
		_ = e.deploymentRepo.Update(ctx, deployment)
		return models.TaskStatusRolledBack, nil
	}

	if time.Since(lastUpdate).Seconds() < float64(step.BatchInterval) ||
		step.ManualApprovalStatus != nil && !*step.ManualApprovalStatus {

		logger.GetLogger().Warn(
			"deployment unhealthy, wait for next batch",
			zap.String("deployment_id", deployment.ID),
			zap.Duration("batch_interval", time.Duration(step.BatchInterval)*time.Second),
			zap.Any("ManualApprovalStatus", step.ManualApprovalStatus),
		)
		return models.TaskStatusRunning, nil
	}

	if idx == len(ss) {
		// 没有下一个阶段, 开始全量
		logger.GetLogger().Info("========full deployment",
			zap.String("deployment_id", deployment.ID),
			zap.Any("version_status", versionMap),
		)
		increase := 0
		idx_ := 0
		for i, version := range versions {
			if version.ID == e.task.Deployment.VersionID {
				idx_ = i
				break
			}
			increase += versionMap[version.ID]
		}
		if idx_ == 0 || increase == 0 {
			return models.TaskStatusRunning, nil
		}

		pkg.Replicas = versionMap[e.task.Deployment.VersionID] + increase
		_, err = e.client.Apply(ctx, e.task.AppID, e.task.Deployment.VersionID, pkg)
		if err != nil {
			return models.TaskStatusFailed, err
		}
		for j := 0; j < idx; j++ {
			version := versions[j]
			pkg.Replicas = 0
			_, err = e.client.Apply(ctx, e.task.AppID, version.ID, pkg)
			if err != nil {
				return models.TaskStatusFailed, err
			}
		}

		return models.TaskStatusRunning, nil
	}

	ns := ss[idx]
	if ns.BatchSize > 0 {
		pkg.Replicas = ns.BatchSize
	} else {
		pkg.Replicas = int(float64(total) * ns.CanaryRatio)
	}

	increaseNum := pkg.Replicas - versionMap[e.task.Deployment.VersionID]
	var decreaseMap = map[string]int{}
	for _, version := range versions {
		if version.ID == e.task.Deployment.VersionID {
			break
		}

		if increaseNum < versionMap[version.ID] {
			decreaseMap[version.ID] = increaseNum
			increaseNum = 0
		} else {
			decreaseMap[version.ID] = versionMap[version.ID]
			increaseNum -= versionMap[version.ID]
		}

		if increaseNum <= 0 {
			break
		}
	}

	logger.GetLogger().Info("========step deployment",
		zap.Any("step", ns),
		zap.Any("version_status", versionMap),
		zap.Any("decrease", decreaseMap),
	)

	pkg.Replicas -= increaseNum
	_, err = e.client.Apply(ctx, e.task.AppID, e.task.Deployment.VersionID, pkg)
	if err != nil {
		return models.TaskStatusFailed, err
	}
	for vid, replicas := range decreaseMap {
		pkg.Replicas = versionMap[vid] - replicas
		_, err = e.client.Apply(ctx, e.task.AppID, vid, pkg)
		if err != nil {
			return models.TaskStatusFailed, err
		}
	}

	return models.TaskStatusRunning, nil
}
